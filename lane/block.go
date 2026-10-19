package lane

import (
 "bytes"
 "errors"
 "math"
)

func EncodeBlock(data []byte, stride int, pack func([]byte) []byte) []byte {
	var out bytes.Buffer
	out.WriteByte(byte(stride))
	for lane := 0; lane < stride; lane++ {
		seg := make([]byte, 0, (len(data)+stride-1)/stride)
		for i := lane; i < len(data); i += stride {
			seg = append(seg, data[i])
		}
		out.Write(pack(seg))
	}
	return out.Bytes()
}
func DecodeBlock(data []byte, n int, unpack func(*bytes.Reader, int) ([]byte,error)) ([]byte, error) {
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
		b, e := unpack(r, sz)
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
func EntropyEstimates(block []byte) [3]float64 {
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
