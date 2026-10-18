package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"index/suffixarray"
	"sync"
)

// Native block-sorting region codec. A recorded sentinel row makes
// suffix sorting reversible even when every byte value occurs in the input.
// Only internal automatic selection is permitted; no additional CLI mode.
const blockBWTPacked byte = 12

func bwtNativeTransform(data []byte) ([]byte, int, error) {
	n := len(data)
	if n == 0 {
		return nil, 0, nil
	}
	var serialized bytes.Buffer
	if err := suffixarray.New(data).Write(&serialized); err != nil {
		return nil, 0, err
	}
	// The standard-library index does not expose its suffix ordering. Its
	// Write representation contains a 10-byte signed length, the original
	// bytes, then length-prefixed varint suffix offsets. This is an
	// implementation-specific Go serialization detail; validate it before
	// interpreting records and preserve regression tests across Go versions.
	b := serialized.Bytes()
	if len(b) < 10+n {
		return nil, 0, errors.New("BWT serialized index truncated")
	}
	declared, parsed := binary.Varint(b[:10])
	if parsed <= 0 || declared != int64(n) || !bytes.Equal(b[10:10+n], data) {
		return nil, 0, errors.New("BWT serialized index format changed")
	}
	at := 10 + n
	full := make([]byte, n+1)
	full[0] = data[n-1]
	count := 1
	primary := -1
	for at < len(b) && count < n+1 {
		if len(b)-at < 10 {
			return nil, 0, errors.New("BWT index record truncated")
		}
		recordStart := at
		size, _ := binary.Varint(b[at:])
		at += 10
		if size < 10 || size > int64(len(b)-recordStart) {
			return nil, 0, errors.New("BWT index record invalid")
		}
		end := recordStart + int(size)
		for at < end {
			value, k := binary.Uvarint(b[at:end])
			if k <= 0 || value >= uint64(n) {
				return nil, 0, errors.New("BWT suffix index invalid")
			}
			at += k
			if value == 0 {
				primary = count
			} else {
				full[count] = data[int(value)-1]
			}
			count++
		}
	}
	if count != n+1 || primary <= 0 {
		return nil, 0, errors.New("BWT incomplete suffix index")
	}
	bwt := make([]byte, n)
	copy(bwt[:primary], full[:primary])
	copy(bwt[primary:], full[primary+1:])
	return bwt, primary, nil
}

func inverseBWTNative(bwt []byte, primary int) ([]byte, error) {
	n := len(bwt)
	if primary <= 0 || primary > n {
		return nil, errors.New("BWT primary index invalid")
	}
	// The sentinel is represented by symbol 0; original bytes use 1..256.
	var freq [257]int
	column := make([]uint16, n+1)
	for i := range column {
		if i == primary {
			column[i] = 0
		} else if i < primary {
			column[i] = uint16(bwt[i]) + 1
		} else {
			column[i] = uint16(bwt[i-1]) + 1
		}
		freq[column[i]]++
	}
	offsets := [257]int{}
	acc := 0
	for i := range offsets {
		offsets[i] = acc
		acc += freq[i]
	}
	lf := make([]int32, n+1)
	for i, v := range column {
		lf[i] = int32(offsets[v])
		offsets[v]++
	}
	out := make([]byte, n)
	row := primary
	for i := n; i >= 0; i-- {
		value := column[row]
		row = int(lf[row])
		if i == n {
			if value != 0 {
				return nil, errors.New("BWT missing sentinel")
			}
			continue
		}
		if value == 0 {
			return nil, errors.New("BWT premature sentinel")
		}
		out[i] = byte(value - 1)
	}
	return out, nil
}

func mtfZeroEncode(data []byte) []byte {
	var sym [256]byte
	var position [256]byte
	for i := range sym {
		sym[i] = byte(i)
		position[i] = byte(i)
	}
	output := make([]byte, 0, len(data))
	run := 0
	for _, v := range data {
		rank := int(position[v])
		if rank == 0 {
			run++
			continue
		}
		if run != 0 {
			output = append(output, 0)
			output = appendUvar(output, uint64(run))
			run = 0
		}
		output = append(output, byte(rank))
		for j := rank; j > 0; j-- {
			other := sym[j-1]
			sym[j] = other
			position[other] = byte(j)
		}
		sym[0] = v
		position[v] = 0
	}
	if run != 0 {
		output = append(output, 0)
		output = appendUvar(output, uint64(run))
	}
	return output
}

func mtfZeroDecode(data []byte, expected int) ([]byte, error) {
	var sym [256]byte
	for i := range sym {
		sym[i] = byte(i)
	}
	out := make([]byte, 0, expected)
	for i := 0; i < len(data); {
		rank := data[i]
		i++
		if rank == 0 {
			run, k := binary.Uvarint(data[i:])
			if k <= 0 || run == 0 || run > uint64(expected-len(out)) {
				return nil, errors.New("BWT MTF run invalid")
			}
			i += k
			for j := uint64(0); j < run; j++ {
				out = append(out, sym[0])
			}
		} else {
			if len(out) >= expected {
				return nil, errors.New("BWT MTF output overflow")
			}
			v := sym[rank]
			out = append(out, v)
			copy(sym[1:int(rank)+1], sym[:rank])
			sym[0] = v
		}
	}
	if len(out) != expected {
		return nil, errors.New("BWT MTF output underflow")
	}
	return out, nil
}
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
				enc, err := encodeBWTNative(block)
				mode := blockBWTPacked
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
	for i, idx := range regional.positions {
		at := idx * stage3BlockSize
		block := input[at:min(len(input), at+stage3BlockSize)]
		data, e := encodeBWTNative(block)
		if e != nil {
			return false, samples, e
		}
		samples.results[i] = mathRegionResult{mode: blockBWTPacked, payload: data}
		baseline += len(regional.results[i].payload)
		total += len(data)
		if len(data) < len(regional.results[i].payload) {
			wins++
		}
	}
	// Require evidence that BWT wins at several independent positions,
	// not merely on one structurally atypical section of the file.
	return wins >= 2 && total*100 <= baseline*95, samples, nil
}
