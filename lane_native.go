package main

import (
	"bytes"
	"errors"
	"math"
	"sync"
)

// Reversible byte-lane decomposition. Independent per-lane native entropy
// streams expose interleaved structural fields that a single byte histogram
// cannot describe. Payload: stride byte (2, 4, or 8), then one independently
// framed native stream per byte-offset lane. Each lane holds indices congruent
// to that offset modulo the stride; its exact decoded length is derived from
// the original region size, so no additional lane lengths are trusted.
const blockLanePacked byte = 13

func encodeLaneBlock(data []byte, stride int) []byte {
	var out bytes.Buffer
	out.WriteByte(byte(stride))
	for lane := 0; lane < stride; lane++ {
		seg := make([]byte, 0, (len(data)+stride-1)/stride)
		for i := lane; i < len(data); i += stride {
			seg = append(seg, data[i])
		}
		out.Write(packStream(seg))
	}
	return out.Bytes()
}
func decodeLaneBlock(data []byte, n int) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("lane payload missing")
	}
	stride := int(data[0])
	if stride != 2 && stride != 4 && stride != 8 {
		return nil, errors.New("invalid lane stride")
	}
	r := bytes.NewReader(data[1:])
	out := make([]byte, n)
	for lane := 0; lane < stride; lane++ {
		sz := (n + stride - 1 - lane) / stride
		b, e := unpackStream(r, sz)
		if e != nil {
			return nil, e
		}
		if len(b) != sz {
			return nil, errors.New("lane length mismatch")
		}
		for k, v := range b {
			out[lane+k*stride] = v
		}
	}
	if r.Len() != 0 {
		return nil, errors.New("lane extra payload")
	}
	return out, nil
}

// One stride-eight histogram also contains exact counts for stride two and
// four: combine positions congruent modulo the requested lane count. This
// avoids rescanning each representative block three times in the gate.
func laneEntropyEstimates(block []byte) [3]float64 {
	var counts [8][256]int
	var totals [8]int
	for i, v := range block {
		lane := i & 7
		counts[lane][v]++
		totals[lane]++
	}
	var entropy [3]float64
	for which, stride := range []int{2, 4, 8} {
		for lane := 0; lane < stride; lane++ {
			var merged [256]int
			n := 0
			for source := lane; source < 8; source += stride {
				n += totals[source]
				for symbol, c := range counts[source] {
					merged[symbol] += c
				}
			}
			if n == 0 {
				continue
			}
			for _, c := range merged {
				if c > 0 {
					entropy[which] += float64(c) * math.Log2(float64(n)/float64(c))
				}
			}
		}
	}
	return entropy
}

// Only compete if distributed zero-order lane entropy predicts a file-wide
// gain beyond the currently selected complete archive. A positive heuristic
// merely authorizes a real full-archive comparison, never determines bytes.
func sampledLaneStride(input []byte, currentBest int) int {
	count := (len(input) + stage3BlockSize - 1) / stage3BlockSize
	if count < 12 {
		return 0
	}
	positions := [3]int{0, count / 2, count - 1}
	bestStride := 0
	bestEstimate := math.MaxFloat64
	sampledBytes := 0
	var sampled [3][]byte
	for i, idx := range positions {
		at := idx * stage3BlockSize
		sampled[i] = input[at:min(len(input), at+stage3BlockSize)]
		sampledBytes += len(sampled[i])
	}
	var estimates [3]float64
	for _, block := range sampled {
		bits := laneEntropyEstimates(block)
		for i := range estimates {
			// Huffman model framing is omitted in this optimistic gate;
			// the winning width is checked again using real archived bytes.
			e := bits[i] / 8
			if e > float64(len(block)) {
				e = float64(len(block))
			}
			estimates[i] += e
		}
	}
	for i, stride := range []int{2, 4, 8} {
		if estimates[i] < bestEstimate {
			bestEstimate = estimates[i]
			bestStride = stride
		}
	}
	// Avoid paying for a second full-file encode if the optimistic entropy
	// bound cannot beat the current complete archive by a useful margin.
	if bestEstimate*100 >= float64(currentBest)*float64(sampledBytes)/float64(len(input))*97 {
		return 0
	}
	// Once worthwhile, select the stride by actual serialized sample bytes,
	// including each lane's framing and codebook; entropy alone can rank
	// nearby strides incorrectly.
	bestPayload := int(^uint(0) >> 1)
	for _, stride := range []int{2, 4, 8} {
		size := 0
		for _, block := range sampled {
			size += len(encodeLaneBlock(block, stride))
		}
		if size < bestPayload {
			bestPayload = size
			bestStride = stride
		}
	}
	// Samples must still plausibly beat the currently selected file rate.
	if bestPayload*100 >= currentBest*sampledBytes/len(input)*98 {
		return 0
	}

	return bestStride
}

// Independent blocks use exactly one automatically selected stride and fall
// back to raw when lane entropy fails to pay for framing; JLZ1 digest and
// existing grouped cross-block decoder remain unchanged.
func encodeLaneArchive(input []byte, digest [32]byte, bestSize int, stride int) ([]byte, error) {
	count := (len(input) + stage3BlockSize - 1) / stage3BlockSize
	results := make([]mathRegionResult, count)
	jobs := make(chan int)
	var wg sync.WaitGroup
	for w := 0; w < availableWorkers(count); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				at := idx * stage3BlockSize
				block := input[at:min(len(input), at+stage3BlockSize)]
				enc := encodeLaneBlock(block, stride)
				if len(enc) >= len(block) {
					results[idx] = mathRegionResult{mode: stage3Raw, payload: block}
				} else {
					results[idx] = mathRegionResult{mode: blockLanePacked, payload: enc}
				}
			}
		}()
	}
	for i := range results {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	var out bytes.Buffer
	out.WriteString(magic)
	out.WriteByte(modeMixedCross)
	putVar(&out, uint64(len(input)))
	out.Write(digest[:])
	for idx, r := range results {
		at := idx * stage3BlockSize
		n := min(len(input)-at, stage3BlockSize)
		out.WriteByte(r.mode)
		putVar(&out, uint64(n))
		putVar(&out, uint64(len(r.payload)))
		out.Write(r.payload)
		if out.Len() >= bestSize {
			return nil, nil
		}
	}
	return out.Bytes(), nil
}
