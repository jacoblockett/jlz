package main

// Native mixed-block codec with self-contained, backward-only references.
// Four-region groups remain independently decodable and parallelizable.
import (
	"bytes"
	"errors"
	"io"
	"sync"
)

const modeMixedCross byte = 17
const (
	blockCrossPacked          byte = 9
	blockMatchPackedSplit     byte = 10
	blockMatchPackedSplitPred byte = 11
)
const crossHistorySize = 128 * 1024

type crossBlock struct {
	typ     byte
	size    int
	payload []byte
}

// As with mixed blocks, framing includes a file-wide digest. Cross references
// may address only bytes from already reconstructed blocks.
const crossGroupRegions = 4

// Four-region independently decodable groups restore parallel decoding while
// allowing bounded references across adjacent region boundaries within a group.
func decodeMixedCross(r *bytes.Reader, n uint64) ([]byte, error) {
	if n > uint64(defaultLimit) {
		return nil, errors.New("cross decoded size exceeds limit")
	}
	var blocks []crossBlock
	produced := uint64(0)
	for produced < n {
		mode, e := r.ReadByte()
		if e != nil {
			return nil, e
		}
		size, e := readVar(r)
		if e != nil {
			return nil, e
		}
		plen, e := readVar(r)
		if e != nil {
			return nil, e
		}
		if size == 0 || size > stage3BlockSize || size > n-produced || plen > uint64(r.Len()) || plen > size*8+1024*1024 || (size < stage3BlockSize && size < n-produced) {
			return nil, errors.New("cross block bounds")
		}
		if mode != stage3Raw && mode != blockCrossPacked && (mode == blockOldDictionary || mode > blockBWTPacked) {
			return nil, errors.New("cross block mode")
		}
		if mode == stage3Raw && plen != size {
			return nil, errors.New("cross raw block size")
		}
		payload := make([]byte, int(plen))
		if _, e = io.ReadFull(r, payload); e != nil {
			return nil, e
		}
		blocks = append(blocks, crossBlock{mode, int(size), payload})
		produced += size
	}
	if r.Len() != 0 {
		return nil, errors.New("cross archive trailing bytes")
	}
	output := make([]byte, int(n))
	groupCount := (len(blocks) + crossGroupRegions - 1) / crossGroupRegions
	errs := make([]error, groupCount)
	jobs := make(chan int)
	var wg sync.WaitGroup
	for w := 0; w < availableWorkers(groupCount); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for group := range jobs {
				lo := group * crossGroupRegions
				hi := min(len(blocks), lo+crossGroupRegions)
				restored := make([]byte, 0, crossGroupRegions*stage3BlockSize)
				for i := lo; i < hi; i++ {
					b := blocks[i]
					var dec []byte
					var err error
					if b.typ == blockBWTPacked {
						dec, err = decodeBWTNative(b.payload, b.size)
					} else if b.typ == stage3Raw {
						dec = b.payload
					} else if b.typ == blockCrossPacked {
						at := max(0, len(restored)-crossHistorySize)
						dec, err = decodeMatchPackedHistory(bytes.NewReader(b.payload), uint64(b.size), false, restored[at:])
					} else {
						dec, err = decodeOneMixedBlock(b.typ, b.payload, b.size)
					}
					if err != nil {
						errs[group] = err
						break
					}
					if len(dec) != b.size {
						errs[group] = errors.New("cross reconstructed block size")
						break
					}
					restored = append(restored, dec...)
				}
				if errs[group] == nil {
					dst := lo * stage3BlockSize
					if dst+len(restored) > len(output) {
						errs[group] = errors.New("cross group bounds")
					} else {
						copy(output[dst:dst+len(restored)], restored)
					}
				}
			}
		}()
	}
	for i := 0; i < groupCount; i++ {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	for _, e := range errs {
		if e != nil {
			return nil, e
		}
	}
	return output, nil
}

// A small sample of 8-byte prefixes shows whether history would help.
// This only limits expensive work; final mode selection uses actual bytes.
func crossPotential(prev, block []byte) bool {
	if len(prev) < 1024 || len(block) < 1024 {
		return false
	}
	seen := make(map[uint64]struct{}, len(prev)/8)
	for i := 0; i+8 <= len(prev); i += 8 {
		seen[binary8(prev[i:])] = struct{}{}
	}
	hits, n := 0, 0
	for i := 0; i+16 <= len(block); i += 64 {
		n++
		if _, yes := seen[binary8(block[i:])]; yes {
			hits++
		}
	}
	return hits >= 10 && hits*100 >= n*4
}
func binary8(b []byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func packCrossRegionTokens(block, tokens []byte) ([]byte, error) {
	st, err := splitMatchTokens(tokens, len(block))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	putVar(&buf, uint64(st.ncommands))
	for i, part := range [][]byte{st.op, st.lengths, st.distances, st.transforms} {
		if i == 2 {
			buf.Write(packVarintStream(part))
		} else {
			buf.Write(packStream(part))
		}
	}
	if len(st.literals)*4 > len(block)*3 {
		buf.Write(packStream(st.literals))
	} else {
		buf.Write(packLiteralStream(st.literals))
	}
	return buf.Bytes(), nil
}
