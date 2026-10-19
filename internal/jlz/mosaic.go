package jlz

// Mosaic mode selects reversible predictors and entropy codecs per region.
// It jointly optimizes variable-length (4-64 KiB) boundaries by actual
// emitted archive byte sizes, rather than trusting estimated token costs.
// This is deliberately bounded exploration; other modes compete globally.

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

const modeMosaic byte = 8
const mosaicQuantum = 4096

type mosaicChoice struct {
	frame  []byte
	method string
	size   int
}

func compressMosaicSegment(data []byte) (mosaicChoice, error) {
	// Every frame is self-delimiting and self-contained. The decoder
	// reverses exactly one selected transform after entropy decoding.
	best := mosaicChoice{}
	add := func(typ, codec byte, lag int, payload []byte) {
		var buf bytes.Buffer
		putVar(&buf, uint64(len(data)))
		buf.WriteByte(typ)
		putVar(&buf, uint64(lag))
		buf.WriteByte(codec)
		putVar(&buf, uint64(len(payload)))
		buf.Write(payload)
		if best.frame == nil || buf.Len() < len(best.frame) {
			best = mosaicChoice{buf.Bytes(), fmt.Sprintf("%d/%d/%d", typ, lag, codec), len(data)}
		}
	}
	// Raw is always legal. codec=2 is exclusively raw/no transform.
	add(predNone, 2, 0, data)
	transforms := []struct {
		typ byte
		lag int
	}{{predNone, 0}, {predXor, 1}, {predSub, 1}, {predXor, 4}, {predSub, 4}}
	for _, t := range transforms {
		if t.lag >= len(data) {
			continue
		}
		transformed := transform(data, t.typ, t.lag)
		out, err := nativeBytes(transformed)
		if err != nil {
			return best, err
		}
		add(t.typ, codecNative, t.lag, out)
		// Entropy coder generally benefits when data is predictable at a byte
		// boundary. Don't spend excessive CPU trying it on random-looking
		// transforms that compress worse even with the native byte coder.
		if len(out) < len(data)*95/100 || t.typ == predNone {
			out = arithmeticEncode(transformed)
			add(t.typ, codecArithmetic, t.lag, out)
		}
	}
	return best, nil
}
func encodeMosaic(input []byte, digest [32]byte) ([]byte, string, error) {
	if len(input) < mosaicQuantum || len(input) > 192*1024 {
		return nil, "", nil
	}
	count := (len(input) + mosaicQuantum - 1) / mosaicQuantum
	dp := make([]int, count+1)
	picks := make([]mosaicChoice, count)
	steps := make([]int, count)
	types := []int{1, 2, 4, 8, 16}
	for at := count - 1; at >= 0; at-- {
		best := int(^uint(0) >> 1)
		for _, step := range types {
			end := at + step
			if end > count {
				end = count
			}
			if end <= at {
				continue
			}
			// Don't redundantly evaluate two candidate lengths that
			// end at the same final, partial block.
			startByte := at * mosaicQuantum
			endByte := end * mosaicQuantum
			if endByte > len(input) {
				endByte = len(input)
			}
			c, err := compressMosaicSegment(input[startByte:endByte])
			if err != nil {
				return nil, "", err
			}
			total := len(c.frame) + dp[end]
			if total < best {
				best = total
				steps[at] = end
				picks[at] = c
			}
			if end == count {
				break
			}
		}
		dp[at] = best
	}
	var archive bytes.Buffer
	archive.WriteString(magic)
	archive.WriteByte(modeMosaic)
	putVar(&archive, uint64(len(input)))
	archive.Write(digest[:])
	stats := make(map[string]int)
	countRegions := 0
	for at := 0; at < count; at = steps[at] {
		c := picks[at]
		archive.Write(c.frame)
		stats[c.method]++
		countRegions++
	}
	return archive.Bytes(), fmt.Sprintf("mosaic, %d regions, coding methods %v", countRegions, stats), nil
}
func decodeMosaic(reader *bytes.Reader, n uint64) ([]byte, error) {
	out := make([]byte, 0, int(n))
	for uint64(len(out)) < n {
		chunkLen, e := readVar(reader)
		if e != nil {
			return nil, e
		}
		if chunkLen == 0 || chunkLen > n-uint64(len(out)) {
			return nil, errors.New("invalid mosaic region length")
		}
		typ, e := reader.ReadByte()
		if e != nil {
			return nil, e
		}
		lag, e := readVar(reader)
		if e != nil {
			return nil, e
		}
		codec, e := reader.ReadByte()
		if e != nil {
			return nil, e
		}
		size, e := readVar(reader)
		if e != nil {
			return nil, e
		}
		if size > uint64(reader.Len()) || size > n*2+4096 {
			return nil, errors.New("invalid mosaic payload length")
		}
		if typ > predSub || (typ == predNone && lag != 0) || (typ != predNone && (lag == 0 || lag >= chunkLen)) {
			return nil, errors.New("invalid mosaic transform")
		}
		if codec > 2 || (codec == 2 && typ != predNone) {
			return nil, errors.New("invalid mosaic codec")
		}
		payload := make([]byte, int(size))
		if _, e = io.ReadFull(reader, payload); e != nil {
			return nil, e
		}
		var block []byte
		switch codec {
		case 0:
			block, e = unnativeBytes(payload, int(chunkLen))
		case 1:
			block, e = arithmeticDecode(payload, int(chunkLen))
		case 2:
			if uint64(len(payload)) != chunkLen {
				return nil, errors.New("invalid raw mosaic size")
			}
			block = payload
		}
		if e != nil {
			return nil, e
		}
		out = append(out, reverseTransform(block, typ, int(lag))...)
	}
	if reader.Len() != 0 {
		return nil, errors.New("extra mosaic archive bytes")
	}
	return out, nil
}
