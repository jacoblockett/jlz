package bwt

import (
	"errors"
)

// BijectiveZeroEncode replaces the MTF rank-zero streak with bijective
// base-two digits in symbols 0 and 1. Nonzero ranks map to 2..254,
// and ranks 254..255 are escaped as 255,0 and 255,1 respectively.
// This is an experimental standalone reversible byte stream; it does not
// read or write any JLZ archive mode.
func BijectiveZeroEncode(data []byte) []byte {
	var sym [256]byte
	var position [256]byte
	for i := range sym {
		sym[i] = byte(i)
		position[i] = byte(i)
	}
	out := make([]byte, 0, len(data))
	run := 0
	flush := func() {
		for run > 0 {
			run--
			out = append(out, byte(run&1))
			run >>= 1
		}
	}
	for _, v := range data {
		rank := int(position[v])
		if rank == 0 {
			run++
			continue
		}
		flush()
		if rank >= 254 {
			out = append(out, 255, byte(rank-254))
		} else {
			out = append(out, byte(rank+1))
		}
		for j := rank; j > 0; j-- {
			old := sym[j-1]
			sym[j] = old
			position[old] = byte(j)
		}
		sym[0] = v
		position[v] = 0
	}
	flush()
	return out
}

// BijectiveZeroDecode reconstructs exactly the original BWT byte sequence,
// rejecting truncated escapes, impossible run lengths and trailing excess.
func BijectiveZeroDecode(data []byte, expected int) ([]byte, error) {
	if expected < 0 {
		return nil, errors.New("negative BWT byte count")
	}
	var sym [256]byte
	for i := range sym {
		sym[i] = byte(i)
	}
	out := make([]byte, 0, expected)
	var run, weight uint64
	weight = 1
	flush := func() error {
		if run > uint64(expected-len(out)) {
			return errors.New("BWT zero-run overflows output")
		}
		for j := uint64(0); j < run; j++ {
			out = append(out, sym[0])
		}
		run = 0
		weight = 1
		return nil
	}
	for i := 0; i < len(data); i++ {
		code := data[i]
		if code <= 1 {
			// Cap run and weight before 64-bit wrap, even on malicious inputs.
			remaining := uint64(expected - len(out))
			if weight > remaining || (uint64(code)+1)*weight > remaining-run {
				return nil, errors.New("BWT zero-run exceeds declared size")
			}
			run += (uint64(code) + 1) * weight
			if weight <= uint64(expected) {
				weight *= 2
			}
			continue
		}
		if err := flush(); err != nil {
			return nil, err
		}
		rank := int(code) - 1
		if code == 255 {
			i++
			if i >= len(data) || data[i] > 1 {
				return nil, errors.New("invalid BWT rank escape")
			}
			rank = 254 + int(data[i])
		}
		if len(out) >= expected {
			return nil, errors.New("BWT output length exceeded")
		}
		value := sym[rank]
		out = append(out, value)
		copy(sym[1:rank+1], sym[:rank])
		sym[0] = value
	}
	if err := flush(); err != nil {
		return nil, err
	}
	if len(out) != expected {
		return nil, errors.New("BWT output length mismatch")
	}
	return out, nil
}
