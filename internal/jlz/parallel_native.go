package jlz

// Parallel encoding of independent blocks using the project's own native
// raw/RLE/Huffman coder. No external codec and no preset DEFLATE dictionary.
import (
	"bytes"
	"errors"
	"io"
	"math"
	"sync"
)

const (
	modeParallelNative      byte = 16 // v1 identifier is isolated by the new JLZ magic
	parallelNativeBlockSize      = 2 * 1024 * 1024
	parallelNativeMinSize        = 4 * 1024 * 1024
)

type parallelBlock struct {
	raw  bool
	data []byte
	err  error
}

func encodeParallelNative(input []byte, digest [32]byte, workers int) ([]byte, error) {
	count := (len(input) + parallelNativeBlockSize - 1) / parallelNativeBlockSize
	if count == 0 {
		return nil, errors.New("empty input")
	}
	if workers <= 0 {
		workers = availableWorkers(count)
	} else {
		workers = min(workers, availableWorkers(count))
	}
	results := make([]parallelBlock, count)
	jobs := make(chan int)
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				start := idx * parallelNativeBlockSize
				end := min(start+parallelNativeBlockSize, len(input))
				src := input[start:end]
				p, e := nativePackBytes(src)
				if e != nil {
					results[idx].err = e
					continue
				}
				if len(p) >= len(src) {
					results[idx] = parallelBlock{raw: true, data: src}
				} else {
					results[idx] = parallelBlock{data: p}
				}
			}
		}()
	}
	for i := 0; i < count; i++ {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	var out bytes.Buffer
	out.WriteString(magic)
	out.WriteByte(modeParallelNative)
	putVar(&out, uint64(len(input)))
	out.Write(digest[:])
	for i, v := range results {
		if v.err != nil {
			return nil, v.err
		}
		start := i * parallelNativeBlockSize
		orig := min(parallelNativeBlockSize, len(input)-start)
		if v.raw {
			out.WriteByte(0)
		} else {
			out.WriteByte(1)
		}
		putVar(&out, uint64(orig))
		putVar(&out, uint64(len(v.data)))
		out.Write(v.data)
	}
	return out.Bytes(), nil
}
func decodeParallelNative(r *bytes.Reader, n uint64) ([]byte, error) {
	if n > uint64(math.MaxInt) {
		return nil, errors.New("invalid size")
	}
	out := make([]byte, 0, int(n))
	for uint64(len(out)) < n {
		kind, e := r.ReadByte()
		if e != nil {
			return nil, e
		}
		orig, e := readVar(r)
		if e != nil {
			return nil, e
		}
		size, e := readVar(r)
		if e != nil {
			return nil, e
		}
		if orig == 0 || orig > parallelNativeBlockSize || orig > n-uint64(len(out)) || size > uint64(r.Len()) || size > orig+65536 {
			return nil, errors.New("invalid native parallel block lengths")
		}
		b := make([]byte, int(size))
		if _, e = io.ReadFull(r, b); e != nil {
			return nil, e
		}
		var restored []byte
		if kind == 0 {
			if size != orig {
				return nil, errors.New("invalid raw length")
			}
			restored = b
		} else if kind == 1 {
			restored, e = nativeUnpackBytes(b, int(orig))
			if e != nil {
				return nil, e
			}
		} else {
			return nil, errors.New("invalid block kind")
		}
		if uint64(len(restored)) != orig {
			return nil, errors.New("native block decoded length")
		}
		out = append(out, restored...)
	}
	if r.Len() != 0 {
		return nil, errors.New("trailing parallel bytes")
	}
	return out, nil
}
