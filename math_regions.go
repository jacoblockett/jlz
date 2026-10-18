package main

// Fixed-work regional mathematical references for large files. Every block is
// independently reconstructible, allowing bounded parallelism without
// storing external dictionaries or requiring third-party compression engines.
// The wire representation is the existing JLZ mixed-block encoding.

import (
	"bytes"
	"sync"
)

type mathRegionResult struct {
	mode    byte
	payload []byte
	err     error
}

// Keep the three expensive block encodings already produced during screening.
// They are reused verbatim when full regional compression is selected.
type mathRegionSamples struct {
	positions [3]int
	results   [3]mathRegionResult
}

func regionalMathUseful(input []byte, bestSize int) (bool, mathRegionSamples) {
	var samples mathRegionSamples
	if len(input) < 3*stage3BlockSize {
		return false, samples
	}
	count := (len(input) + stage3BlockSize - 1) / stage3BlockSize
	positions := [3]int{0, count / 2, count - 1}
	samples.positions = positions
	good := 0
	substantial := 0
	// Compare the representative block ratio with the *actual* best whole-file
	// rate, counting the mode and varints that the block needs on disk.
	baseline := float64(bestSize) / float64(len(input))
	for i, idx := range positions {
		at := idx * stage3BlockSize
		end := min(len(input), at+stage3BlockSize)
		block := input[at:end]
		mode, encoded, err := encodeMathRegionBlockFast(block)
		if err != nil {
			return false, mathRegionSamples{}
		}
		samples.results[i] = mathRegionResult{mode: mode, payload: encoded}
		ratio := float64(len(encoded)+1+2*10) / float64(len(block))
		if ratio < baseline*0.98 {
			good++
		}
		if ratio < baseline*0.75 {
			substantial++
		}
	}
	return good >= 2 || substantial >= 2, samples
}

func encodeParallelMathRegions(input []byte, digest [32]byte, bestSize int, samples mathRegionSamples) ([]byte, error) {
	count := (len(input) + stage3BlockSize - 1) / stage3BlockSize
	results := make([]mathRegionResult, count)
	for i, idx := range samples.positions {
		if idx%crossGroupRegions == 0 || !crossPotential(input[max((idx/crossGroupRegions)*crossGroupRegions*stage3BlockSize, idx*stage3BlockSize-crossHistorySize):idx*stage3BlockSize], input[idx*stage3BlockSize:min(len(input), (idx+1)*stage3BlockSize)]) {
			results[idx] = samples.results[i]
		}
	}
	jobs := make(chan int)
	var wg sync.WaitGroup
	for w := 0; w < availableWorkers(count); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				at := idx * stage3BlockSize
				end := min(len(input), at+stage3BlockSize)
				var mode byte
				var p []byte
				var e error
				lo := max((idx/crossGroupRegions)*crossGroupRegions*stage3BlockSize, at-crossHistorySize)
				if idx%crossGroupRegions != 0 && crossPotential(input[lo:at], input[at:end]) {
					mode, p, e = encodeMathRegionBlockHistory(input[at:end], input[lo:at])
				} else {
					mode, p, e = encodeMathRegionBlockFast(input[at:end])
				}
				results[idx] = mathRegionResult{mode, p, e}
			}
		}()
	}
	for i := 0; i < count; i++ {
		if results[i].payload != nil {
			continue
		}
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	var out bytes.Buffer
	out.WriteString(magic)
	out.WriteByte(modeMixedCross)
	putVar(&out, uint64(len(input)))
	out.Write(digest[:])
	for idx, res := range results {
		if res.err != nil {
			return nil, res.err
		}
		at := idx * stage3BlockSize
		n := min(stage3BlockSize, len(input)-at)
		out.WriteByte(res.mode)
		putVar(&out, uint64(n))
		putVar(&out, uint64(len(res.payload)))
		out.Write(res.payload)
		if out.Len() >= bestSize {
			return nil, nil
		}
	}
	return out.Bytes(), nil
}
