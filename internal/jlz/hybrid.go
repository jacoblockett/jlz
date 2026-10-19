package jlz

// Experimental relationship graph retained as in-project component, not an optional CLI encoder. with per-node transforms,
// local predictors, and selectable literal/native-byte/bit-grammar payloads.
// This is an experiment, not an efficient or hardened compression format.

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"sort"
)

const (
	hSubResidual  byte = 6
	hBitRotate    byte = 7
	hLocalXor     byte = 8
	hLocalSub     byte = 9
	hLiteral      byte = 10
	hCodecRaw     byte = 0
	hCodecNative  byte = 1
	hCodecGrammar byte = 2
)

type hybridStats struct {
	counts               [11]int
	grammar, native, raw int
}

func bytePayload(raw []byte, allowGrammar bool) []byte {
	// The decoder knows the uncompressed block size from the outer framing.
	best := append([]byte{hCodecRaw}, raw...)
	var zipped bytes.Buffer
	w, err := newNativeWriter(&zipped)
	if err == nil {
		_, e1 := w.Write(raw)
		e2 := w.Close()
		if e1 == nil && e2 == nil {
			var b bytes.Buffer
			b.WriteByte(hCodecNative)
			putVar(&b, uint64(zipped.Len()))
			b.Write(zipped.Bytes())
			if b.Len() < len(best) {
				best = b.Bytes()
			}
		}
	}
	return best
}

func frame(kind byte, args []uint64, mask []byte, payload []byte) []byte {
	var b bytes.Buffer
	b.WriteByte(kind)
	for _, a := range args {
		putVar(&b, a)
	}
	if kind == blockXorConstant || kind == hBitRotate {
		b.Write(mask)
	}
	if kind == blockXorPattern {
		putVar(&b, uint64(len(mask)))
		b.Write(mask)
	}
	if payload != nil {
		b.Write(payload)
	}
	return b.Bytes()
}

func compressedResidual(current, ref []byte, kind byte) ([]byte, int) {
	delta := make([]byte, len(current))
	var hist [256]int
	for i := range current {
		if kind == hSubResidual {
			delta[i] = current[i] - ref[i]
		} else {
			delta[i] = current[i] ^ ref[i]
		}
		hist[delta[i]]++
	}
	max := 0
	for _, n := range hist {
		if n > max {
			max = n
		}
	}
	return delta, max
}

func rotateLeftBits(in []byte, k byte) []byte {
	out := make([]byte, len(in))
	if len(in) == 0 {
		return out
	}
	for i := range in {
		out[i] = (in[i] << k) | (in[(i+1)%len(in)] >> (8 - k))
	}
	return out
}

func chooseHybridBlocks(input []byte, width int, useGrammar bool) ([]byte, hybridStats) {
	count := (len(input) + width - 1) / width
	chunks := make([][]byte, count)
	for i := 0; i < count; i++ {
		end := (i + 1) * width
		if end > len(input) {
			end = len(input)
		}
		chunks[i] = input[i*width : end]
	}
	var result bytes.Buffer
	var st hybridStats
	sigs := make(map[uint64][]int, count)
	norms := make(map[uint64][]int, count)
	for i, cur := range chunks {
		best := frame(hLiteral, nil, nil, bytePayload(cur, useGrammar))
		bestKind := hLiteral
		offer := func(kind byte, args []uint64, mask []byte, payload []byte) {
			c := frame(kind, args, mask, payload)
			if len(c) < len(best) {
				best, bestKind = c, kind
			}
		}
		// Equality/NOT/constants, looking across the entire past, indexed by
		// a content signature invariant under constant-byte XOR.
		for _, j := range norms[normHash(cur)] {
			if len(chunks[j]) != len(cur) {
				continue
			}
			if mask, ok := xorConstant(cur, chunks[j]); ok {
				if mask == 0 {
					offer(blockCopy, []uint64{uint64(j)}, nil, nil)
				} else {
					offer(blockXorConstant, []uint64{uint64(j)}, []byte{mask}, nil)
				}
			}
		}
		// Discover C=A XOR B, including cases where A/B are derived blocks.
		if len(cur) >= 8 {
			for j := i - 1; j >= 0; j-- {
				if len(chunks[j]) != len(cur) {
					continue
				}
				wanted := blockSig(cur) ^ blockSig(chunks[j])
				for _, k := range sigs[wanted] {
					if k >= j || len(chunks[k]) != len(cur) {
						continue
					}
					if xorPairEqual(chunks[j], chunks[k], cur) {
						offer(blockXorPair, []uint64{uint64(j), uint64(k)}, nil, nil)
					}
				}
			}
		}
		// Short periodic XOR masks and cheap exact bit-rotation relations.
		from := i - 24
		if from < 0 {
			from = 0
		}
		for j := i - 1; j >= from; j-- {
			ref := chunks[j]
			if len(ref) != len(cur) {
				continue
			}
			if p, ok := xorPattern(cur, ref); ok {
				offer(blockXorPattern, []uint64{uint64(j)}, p, nil)
			}
			for k := byte(1); k <= 7; k++ {
				if bytes.Equal(cur, rotateLeftBits(ref, k)) {
					offer(hBitRotate, []uint64{uint64(j)}, []byte{k}, nil)
					break
				}
			}
			// Only materialize residuals which have a dominant byte value.
			// Random residuals have none, and paying for them would be wasteful.
			for _, kind := range []byte{blockXorResidual, hSubResidual} {
				delta, peak := compressedResidual(cur, ref, kind)
				if peak*100 < len(cur)*12 {
					continue
				}
				payload := bytePayload(delta, false)
				offer(kind, []uint64{uint64(j)}, nil, payload)
			}
		}
		// Within-node periodic/lag predictors, independent of graph references.
		// The output is still lossless because initial lag bytes are literals.
		if len(cur) >= 64 {
			lags := []int{1, 2, 3, 4, 8, 16, 32, 64, 128, 256, 512}
			for _, lag := range lags {
				if lag >= len(cur)/2 {
					break
				}
				for _, kind := range []byte{hLocalXor, hLocalSub} {
					delta := make([]byte, len(cur))
					copy(delta[:lag], cur[:lag])
					var hist [256]int
					for p := lag; p < len(cur); p++ {
						if kind == hLocalXor {
							delta[p] = cur[p] ^ cur[p-lag]
						} else {
							delta[p] = cur[p] - cur[p-lag]
						}
						hist[delta[p]]++
					}
					peak := 0
					for _, n := range hist {
						if n > peak {
							peak = n
						}
					}
					if peak*100 < (len(cur)-lag)*20 {
						continue
					}
					offer(kind, []uint64{uint64(lag)}, nil, bytePayload(delta, false))
				}
			}
		}
		result.Write(best)
		st.counts[bestKind]++
		// Non-raw payload statistics available for quick auditing.
		if bestKind == hLiteral {
			if len(best) > 1 {
				switch best[1] {
				case hCodecRaw:
					st.raw++
				case hCodecNative:
					st.native++
				case hCodecGrammar:
					st.grammar++
				}
			}
		}
		sigs[blockSig(cur)] = append(sigs[blockSig(cur)], i)
		norms[normHash(cur)] = append(norms[normHash(cur)], i)
	}
	return result.Bytes(), st
}

func encodeHybrid(input []byte, digest [32]byte, width int) ([]byte, string, error) {
	return encodeHybridInternal(input, digest, width, true)
}
func encodeHybridInternal(input []byte, digest [32]byte, width int, enableGrammar bool) ([]byte, string, error) {
	if width <= 0 {
		return nil, "", errors.New("invalid width")
	}
	body, st := chooseHybridBlocks(input, width, enableGrammar && width <= 1024)
	var out bytes.Buffer
	out.WriteString(magic)
	out.WriteByte(modeHybrid)
	putVar(&out, uint64(len(input)))
	out.Write(digest[:])
	putVar(&out, uint64(width))
	out.Write(body)
	labels := []string{"copy", "xor-constant", "xor-period", "xor-pair", "xor-residual", "sub-residual", "rotate", "local-xor", "local-sub", "literal"}
	var parts []string
	for k := 1; k < len(st.counts); k++ {
		if st.counts[k] > 0 {
			parts = append(parts, fmt.Sprintf("%s=%d", labels[k-1], st.counts[k]))
		}
	}
	sort.Strings(parts)
	return out.Bytes(), fmt.Sprintf("block=%d; %s; codecs: raw=%d native=%d grammar=%d", width, joinLabels(parts), st.raw, st.native, st.grammar), nil
}
func joinLabels(a []string) string {
	if len(a) == 0 {
		return "none"
	}
	var b bytes.Buffer
	for i, s := range a {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(s)
	}
	return b.String()
}
func bestHybrid(input []byte, digest [32]byte, enableGrammar bool) ([]byte, string, error) {
	var best []byte
	var detail string
	for _, w := range []int{256, 512, 1024, 2048, 4096, 8192} {
		if w > len(input) || (len(input)+w-1)/w > 512 {
			continue
		}
		arc, desc, err := encodeHybridInternal(input, digest, w, enableGrammar)
		if err != nil {
			return nil, "", err
		}
		if best == nil || len(arc) < len(best) {
			best, detail = arc, desc
		}
	}
	return best, detail, nil
}

func decodePayload(r *bytes.Reader, size int, limit int64) ([]byte, error) {
	codec, e := r.ReadByte()
	if e != nil {
		return nil, e
	}
	if codec == hCodecRaw {
		b := make([]byte, size)
		_, e = io.ReadFull(r, b)
		return b, e
	}
	encN, e := readVar(r)
	if e != nil {
		return nil, e
	}
	if encN > uint64(r.Len()) {
		return nil, errors.New("truncated hybrid payload")
	}
	if encN > uint64(limit)+1024*1024 {
		return nil, errors.New("excessive hybrid payload")
	}
	enc := make([]byte, int(encN))
	if _, e = io.ReadFull(r, enc); e != nil {
		return nil, e
	}
	switch codec {
	case hCodecNative:
		fr := newNativeReader(bytes.NewReader(enc))
		defer fr.Close()
		out := make([]byte, size)
		if _, e = io.ReadFull(fr, out); e != nil {
			return nil, e
		}
		var tail [1]byte
		n, err := fr.Read(tail[:])
		if n != 0 || err != io.EOF {
			return nil, errors.New("hybrid native payload expands beyond chunk")
		}
		return out, nil
	case hCodecGrammar:
		out, err := decompress(enc, int64(size))
		if err != nil {
			return nil, err
		}
		if len(out) != size {
			return nil, errors.New("hybrid grammar decoded incorrect size")
		}
		return out, nil
	default:
		return nil, errors.New("invalid hybrid payload codec")
	}
}
func decodeHybrid(r *bytes.Reader, n uint64, limit int64) ([]byte, error) {
	width64, e := readVar(r)
	if e != nil {
		return nil, e
	}
	if width64 < 1 || width64 > 8192 {
		return nil, errors.New("invalid hybrid chunk size")
	}
	width := int(width64)
	out := make([]byte, int(n))
	count := (len(out) + width - 1) / width
	for i := 0; i < count; i++ {
		off := i * width
		sz := width
		if off+sz > len(out) {
			sz = len(out) - off
		}
		dst := out[off : off+sz]
		kind, e := r.ReadByte()
		if e != nil {
			return nil, e
		}
		refA, refB := 0, 0
		switch kind {
		case blockCopy, blockXorConstant, blockXorPattern, blockXorResidual, hSubResidual, hBitRotate:
			v, err := readVar(r)
			if err != nil {
				return nil, err
			}
			if v >= uint64(i) {
				return nil, errors.New("hybrid forward reference")
			}
			refA = int(v)
		case blockXorPair:
			a, err := readVar(r)
			if err != nil {
				return nil, err
			}
			b, err := readVar(r)
			if err != nil {
				return nil, err
			}
			if a >= uint64(i) || b >= uint64(i) {
				return nil, errors.New("hybrid pair forward reference")
			}
			refA, refB = int(a), int(b)
		case hLocalXor, hLocalSub:
			lag, err := readVar(r)
			if err != nil {
				return nil, err
			}
			if lag == 0 || lag >= uint64(sz) {
				return nil, errors.New("invalid local lag")
			}
			refA = int(lag)
		case hLiteral:
		default:
			return nil, fmt.Errorf("invalid hybrid opcode %d", kind)
		}
		prior := func(j int) ([]byte, error) {
			start := j * width
			end := start + sz
			if end > (j+1)*width || end > len(out) {
				return nil, errors.New("hybrid reference length mismatch")
			}
			return out[start:end], nil
		}
		var srcA, srcB []byte
		if kind >= blockCopy && kind <= blockXorResidual || kind == hSubResidual || kind == hBitRotate {
			srcA, e = prior(refA)
			if e != nil {
				return nil, e
			}
		}
		if kind == blockXorPair {
			srcB, e = prior(refB)
			if e != nil {
				return nil, e
			}
		}
		switch kind {
		case blockCopy:
			copy(dst, srcA)
		case blockXorConstant:
			m, err := r.ReadByte()
			if err != nil {
				return nil, err
			}
			for p := range dst {
				dst[p] = srcA[p] ^ m
			}
		case blockXorPattern:
			k, err := readVar(r)
			if err != nil {
				return nil, err
			}
			if k < 1 || k > 16 {
				return nil, errors.New("invalid periodic mask")
			}
			mask := make([]byte, int(k))
			if _, err = io.ReadFull(r, mask); err != nil {
				return nil, err
			}
			for p := range dst {
				dst[p] = srcA[p] ^ mask[p%len(mask)]
			}
		case blockXorPair:
			for p := range dst {
				dst[p] = srcA[p] ^ srcB[p]
			}
		case hBitRotate:
			shift, err := r.ReadByte()
			if err != nil {
				return nil, err
			}
			if shift < 1 || shift > 7 {
				return nil, errors.New("invalid bit shift")
			}
			copy(dst, rotateLeftBits(srcA, shift))
		case hLiteral, blockXorResidual, hSubResidual, hLocalXor, hLocalSub:
			values, err := decodePayload(r, sz, limit)
			if err != nil {
				return nil, err
			}
			copy(dst, values)
			switch kind {
			case blockXorResidual:
				for p := range dst {
					dst[p] ^= srcA[p]
				}
			case hSubResidual:
				for p := range dst {
					dst[p] += srcA[p]
				}
			case hLocalXor:
				for p := refA; p < len(dst); p++ {
					dst[p] ^= dst[p-refA]
				}
			case hLocalSub:
				for p := refA; p < len(dst); p++ {
					dst[p] += dst[p-refA]
				}
			}
		}
	}
	if r.Len() != 0 {
		return nil, errors.New("extra bytes in hybrid archive")
	}
	return out, nil
}

// For direct experiments without losing the checksummed outer format.
func testHybridDirect(input []byte, width int) ([]byte, error) {
	arc, _, err := encodeHybrid(input, sha256.Sum256(input), width)
	return arc, err
}
