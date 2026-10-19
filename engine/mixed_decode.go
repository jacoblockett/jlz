package engine

// The mixed format already stores independent blocks. Decode them in parallel
// into disjoint ranges of one output buffer, retaining the original framing,
// corruption checks and checksum validation in the outer JLZ decoder.

import (
	"bytes"
	"errors"
	"io"
	"math"
	"sync"
)

type parallelDecodeJob struct {
	index   int
	offset  int
	size    int
	mode    byte
	payload []byte
}

func decodeMixedBlocksParallel(r *bytes.Reader, n uint64) ([]byte, error) {
	if n > uint64(math.MaxInt) {
		return nil, errors.New("mixed archive too large")
	}
	output := make([]byte, int(n))
	count := (int(n) + stage3BlockSize - 1) / stage3BlockSize
	errs := make([]error, count)
	jobs := make(chan parallelDecodeJob, 8)
	var wg sync.WaitGroup
	for worker := 0; worker < availableWorkers(count); worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				restored, err := decodeOneMixedBlock(job.mode, job.payload, job.size)
				if err == nil && len(restored) != job.size {
					err = errors.New("mixed block length mismatch")
				}
				if err == nil {
					copy(output[job.offset:job.offset+job.size], restored)
				}
				errs[job.index] = err
			}
		}()
	}
	offset, index := 0, 0
	var parseErr error
	for uint64(offset) < n {
		mode, err := r.ReadByte()
		if err != nil {
			parseErr = err
			break
		}
		sz, err := readVar(r)
		if err != nil {
			parseErr = err
			break
		}
		length, err := readVar(r)
		if err != nil {
			parseErr = err
			break
		}
		if sz == 0 || sz > stage3BlockSize || sz > n-uint64(offset) || length > uint64(r.Len()) || length > sz*8+1024*1024 {
			parseErr = errors.New("invalid mixed block size or payload length")
			break
		}
		if mode == stage3Raw {
			if length != sz {
				parseErr = errors.New("wrong mixed raw length")
				break
			}
			if _, err = io.ReadFull(r, output[offset:offset+int(sz)]); err != nil {
				parseErr = err
				break
			}
		} else {
			if mode == blockOldDictionary || mode > blockMatchPackedSplitPred || mode == 0 {
				parseErr = errors.New("unknown or unsupported mixed block mode")
				break
			}
			content := make([]byte, int(length))
			if _, err = io.ReadFull(r, content); err != nil {
				parseErr = err
				break
			}
			jobs <- parallelDecodeJob{index: index, offset: offset, size: int(sz), mode: mode, payload: content}
		}
		offset += int(sz)
		index++
	}
	if parseErr == nil && r.Len() != 0 {
		parseErr = errors.New("trailing mixed archive data")
	}
	close(jobs)
	wg.Wait()
	if parseErr != nil {
		return nil, parseErr
	}
	for _, err := range errs[:index] {
		if err != nil {
			return nil, err
		}
	}
	return output, nil
}

func decodeOneMixedBlock(mode byte, content []byte, size int) ([]byte, error) {
	r := bytes.NewReader(content)
	var restored []byte
	var err error
	switch mode {
	case blockNativeBytes:
		return nativeUnpackBytes(content, size)
	case blockPredictive:
		restored, err = decodePredictive(r, uint64(size))
	case blockLegacyMatch:
		restored, err = decodeMatchV2(r, uint64(size), int64(size))
	case blockNative:
		restored, err = decodeMatchNative(r, uint64(size))
	case blockMatchedLiteral:
		restored, err = decodeMatchNativeVariant(r, uint64(size), true)
	case blockMatchPacked, blockMatchPackedPred:
		restored, err = decodeMatchPacked(r, uint64(size), mode == blockMatchPackedPred)
	case blockMatchPackedSplit, blockMatchPackedSplitPred:
		restored, err = decodePackedWithLengthMode(r, uint64(size), mode == blockMatchPackedSplitPred, nil, true)
	default:
		return nil, errors.New("unknown mixed block mode")
	}
	if err != nil {
		return nil, err
	}
	if r.Len() != 0 {
		return nil, errors.New("mixed payload trailing bytes")
	}
	return restored, nil
}
