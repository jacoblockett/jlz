package jlz

import (
	"bytes"
	"encoding/binary"
	"sync"
)

const blockAdditivePacked byte = 14

// Screen for local arithmetic progressions in 16-bit modular integer
// positions. A majority of signed differences fitting in one zigzag byte
// is necessary (not sufficient) for this representation to pay. This
// limits costly Huffman sample work on irrelevant nonnumeric inputs.
func likelyAdditiveStructure(input []byte, positions [3]int) bool {
	plausible := 0
	for _, idx := range positions {
		at := idx * stage3BlockSize
		b := input[at:min(len(input), at+stage3BlockSize)]
		words := len(b) / 2
		if words < 2 {
			continue
		}
		prev := binary.LittleEndian.Uint16(b)
		small := 0
		for i := 1; i < words; i++ {
			v := binary.LittleEndian.Uint16(b[i*2:])
			d := int16(v - prev)
			if d >= -128 && d <= 127 {
				small++
			}
			prev = v
		}
		if small*4 >= (words-1)*3 {
			plausible++
		}
	}
	return plausible >= 2
}

// Compete only where a different, already supported archive family has
// demonstrated BWT wins. Distributed actual sizes (including Huffman
// codebooks) must independently favor modular affine-word residuals in at
// least two samples and beat aggregate BWT sample size by >2%.
func additiveWordPreference(input []byte, bwt mathRegionSamples) bool {
	if !likelyAdditiveStructure(input, bwt.positions) {
		return false
	}
	wins, original, predicted := 0, 0, 0
	for i, idx := range bwt.positions {
		at := idx * stage3BlockSize
		block := input[at:min(len(input), at+stage3BlockSize)]
		enc := encodeAdditiveZigzagWords(block)
		n := len(enc)
		m := len(bwt.results[i].payload)
		if n < m {
			wins++
		}
		original += m
		predicted += n
	}
	return wins >= 2 && predicted*100 <= original*98
}

// Encode whole original regions independently. The existing modeMixedCross
// decoder handles framing and SHA-256. The new block mode uses order-one
// modular 16-bit additive residues in two separate entropy streams.
func encodeParallelAdditiveArchive(input []byte, digest [32]byte) ([]byte, error) {
	count := (len(input) + stage3BlockSize - 1) / stage3BlockSize
	blocks := make([]mathRegionResult, count)
	jobs := make(chan int)
	var wg sync.WaitGroup
	for w := 0; w < availableWorkers(count); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				at := idx * stage3BlockSize
				block := input[at:min(len(input), at+stage3BlockSize)]
				enc := encodeAdditiveZigzagWords(block)
				if len(enc) >= len(block) {
					blocks[idx] = mathRegionResult{mode: stage3Raw, payload: block}
				} else {
					blocks[idx] = mathRegionResult{mode: blockAdditivePacked, payload: enc}
				}
			}
		}()
	}
	for i := range blocks {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	var out bytes.Buffer
	out.WriteString(magic)
	out.WriteByte(modeMixedCross)
	putVar(&out, uint64(len(input)))
	out.Write(digest[:])
	for idx, r := range blocks {
		at := idx * stage3BlockSize
		n := min(stage3BlockSize, len(input)-at)
		out.WriteByte(r.mode)
		putVar(&out, uint64(n))
		putVar(&out, uint64(len(r.payload)))
		out.Write(r.payload)
	}
	return out.Bytes(), nil
}
