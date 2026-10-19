package engine

import (
 "github.com/jacoblockett/jlz/lane"
	"bytes"
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

func encodeLaneBlock(data []byte, stride int) []byte { return lane.EncodeBlock(data,stride,packStream) }
func decodeLaneBlock(data []byte, n int) ([]byte,error) {
 return lane.DecodeBlock(data,n,unpackStream)
}
func laneEntropyEstimates(data []byte) [3]float64 { return lane.EntropyEstimates(data) }

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
