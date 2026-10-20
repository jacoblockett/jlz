package engine

import (
	"bytes"
	"errors"
	"sync"

	"github.com/jacoblockett/jlz/bwt"
)

const (
	blockBWTPacked    byte = 12
	blockBWTBijective byte = 15 // zero-run bijective binary, legacy mode 12 still readable
)

func bwtNativeTransform(data []byte) ([]byte, int, error)       { return bwt.Transform(data) }
func inverseBWTNative(data []byte, primary int) ([]byte, error) { return bwt.Inverse(data, primary) }
func mtfZeroEncode(data []byte) []byte                          { return bwt.MTFEncode(data) }
func mtfZeroDecode(data []byte, expected int) ([]byte, error)   { return bwt.MTFDecode(data, expected) }

func encodeBWTNative(data []byte) ([]byte, error) {
	bwt, primary, err := bwtNativeTransform(data)
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	putVar(&out, uint64(primary))
	out.Write(packStream(mtfZeroEncode(bwt)))
	return out.Bytes(), nil
}

// encodeBWTBijective preserves the original BWT transform while replacing
// only its MTF zero-run representation with a bijective binary digit stream.
func encodeBWTBijective(data []byte) ([]byte, error) {
	transformed, primary, err := bwtNativeTransform(data)
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	putVar(&out, uint64(primary))
	out.Write(packFiniteState(bwt.BijectiveZeroEncode(transformed)))
	return out.Bytes(), nil
}

func decodeBWTBijective(data []byte, expected int) ([]byte, error) {
	r := bytes.NewReader(data)
	primary, err := readVar(r)
	if err != nil {
		return nil, err
	}
	if expected < 1 || primary > uint64(expected) {
		return nil, errors.New("bijective BWT primary out of bounds")
	}
	encoded, err := unpackStream(r, expected*2+20)
	if err != nil {
		return nil, err
	}
	if r.Len() != 0 {
		return nil, errors.New("bijective BWT trailing payload")
	}
	transformed, err := bwt.BijectiveZeroDecode(encoded, expected)
	if err != nil {
		return nil, err
	}
	return inverseBWTNative(transformed, int(primary))
}

func decodeBWTNative(data []byte, expected int) ([]byte, error) {
	r := bytes.NewReader(data)
	primary, err := readVar(r)
	if err != nil {
		return nil, err
	}
	if expected < 1 || primary > uint64(expected) {
		return nil, errors.New("BWT primary out of bounds")
	}
	encoded, err := unpackStream(r, expected*2+20)
	if err != nil {
		return nil, err
	}
	if r.Len() != 0 {
		return nil, errors.New("BWT trailing payload")
	}
	bwt, err := mtfZeroDecode(encoded, expected)
	if err != nil {
		return nil, err
	}
	return inverseBWTNative(bwt, int(primary))
}
func encodeParallelBWTRegions(input []byte, digest [32]byte, rawSize int, samples mathRegionSamples) ([]byte, error) {
	count := (len(input) + stage3BlockSize - 1) / stage3BlockSize
	results := make([]mathRegionResult, count)
	for i, idx := range samples.positions {
		results[idx] = samples.results[i]
	}
	jobs := make(chan int)
	var wg sync.WaitGroup
	for worker := 0; worker < availableWorkers(count); worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				at := idx * stage3BlockSize
				block := input[at:min(len(input), at+stage3BlockSize)]
				enc, mode, err := encodeBWTSubregions(block)
				if len(enc) >= len(block) {
					mode = stage3Raw
					enc = block
				}
				results[idx] = mathRegionResult{mode: mode, payload: enc, err: err}
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
	for idx, r := range results {
		if r.err != nil {
			return nil, r.err
		}
		if len(r.payload) == 0 {
			return nil, errors.New("BWT region has no payload")
		}
		at := idx * stage3BlockSize
		n := min(len(input)-at, stage3BlockSize)
		out.WriteByte(r.mode)
		putVar(&out, uint64(n))
		putVar(&out, uint64(len(r.payload)))
		out.Write(r.payload)
		if out.Len() >= rawSize {
			return nil, nil
		}
	}
	return out.Bytes(), nil
}

func sampleBWTPreference(input []byte, regional mathRegionSamples) (bool, mathRegionSamples, error) {
	var samples mathRegionSamples
	samples.positions = regional.positions
	baseline, total := 0, 0
	wins := 0
	for _, result := range regional.results {
		baseline += len(result.payload)
	}
	for i, idx := range regional.positions {
		at := idx * stage3BlockSize
		block := input[at:min(len(input), at+stage3BlockSize)]
		// Evaluate the exact bijective representation emitted by full-region BWT.
		// A legacy-MTF sample can reject a genuinely smaller new-format archive.
		data, e := encodeBWTBijective(block)
		if e != nil {
			return false, samples, e
		}
		samples.results[i] = mathRegionResult{mode: blockBWTBijective, payload: data}
		total += len(data)
		if len(data) < len(regional.results[i].payload) {
			wins++
		}
		// Remaining nonnegative payloads cannot recover a lost vote or
		// undo an already-exceeded total-size margin. Never reject a winner.
		remaining := 2 - i
		if wins+remaining < 2 || total*100 > baseline*98 {
			return false, samples, nil
		}
	}
	// Require two of three independent position wins and a 2% total margin.
	// The full archive still uses measured sizes, not the sample estimate.
	return wins >= 2 && total*100 <= baseline*98, samples, nil
}
