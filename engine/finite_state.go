package engine

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// The lower bound and 12-bit state quotient allow byte renormalization with
// uint32 arithmetic. The model is independently serialized per chosen stream.
const finiteStateScale = 12
const finiteStateTotal = 1 << finiteStateScale
const finiteStateLower = 1 << 23

type finiteStateModel struct {
	freq   [256]uint32
	cum    [256]uint32
	lookup [finiteStateTotal]byte
	count  int
}

func finiteStateNormalize(counts [256]uint64) (m finiteStateModel, err error) {
	var total uint64
	for _, v := range counts {
		total += v
	}
	if total == 0 {
		return m, errors.New("empty model")
	}
	sum := 0
	var rem [256]uint64
	for s, n := range counts {
		if n == 0 {
			continue
		}
		q := int((n * finiteStateTotal) / total)
		if q == 0 {
			q = 1
		}
		m.freq[s] = uint32(q)
		rem[s] = (n * finiteStateTotal) % total
		sum += q
		m.count++
	}
	for sum < finiteStateTotal {
		best := -1
		for s, n := range counts {
			if n == 0 {
				continue
			}
			if best < 0 || rem[s] > rem[best] || (rem[s] == rem[best] && counts[s] > counts[best]) {
				best = s
			}
		}
		if best < 0 {
			return m, errors.New("missing symbol")
		}
		m.freq[best]++
		rem[best] = 0
		sum++
	}
	for sum > finiteStateTotal {
		best := -1
		for s, n := range m.freq {
			if n > 1 && (best < 0 || n > m.freq[best]) {
				best = s
			}
		}
		if best < 0 {
			return m, errors.New("frequency overflow")
		}
		m.freq[best]--
		sum--
	}
	cum := uint32(0)
	for s, f := range m.freq {
		m.cum[s] = cum
		for j := uint32(0); j < f; j++ {
			m.lookup[int(cum+j)] = byte(s)
		}
		cum += f
	}
	if cum != finiteStateTotal {
		return m, errors.New("bad model sum")
	}
	return m, nil
}
func finiteStateModelFor(data []byte) (finiteStateModel, error) {
	var f [256]uint64
	for _, v := range data {
		f[v]++
	}
	return finiteStateNormalize(f)
}
func finiteStateTable(m *finiteStateModel) []byte {
	b := make([]byte, 0, 1+m.count*3)
	b = binary.AppendUvarint(b, uint64(m.count))
	for s, f := range m.freq {
		if f != 0 {
			b = append(b, byte(s))
			b = binary.AppendUvarint(b, uint64(f))
		}
	}
	return b
}
func finiteStateReadTable(buf []byte) (m finiteStateModel, n int, err error) {
	count, k := binary.Uvarint(buf)
	if k <= 0 || count < 1 || count > 256 {
		return m, 0, errors.New("bad count")
	}
	p := k
	for i := uint64(0); i < count; i++ {
		if p >= len(buf) {
			return m, 0, errors.New("table trunc")
		}
		s := buf[p]
		p++
		f, j := binary.Uvarint(buf[p:])
		if j <= 0 || f > finiteStateTotal || f == 0 || m.freq[s] != 0 {
			return m, 0, errors.New("bad frequency")
		}
		p += j
		m.freq[s] = uint32(f)
	}
	cum := uint32(0)
	for s, f := range m.freq {
		m.cum[s] = cum
		for j := uint32(0); j < f; j++ {
			if cum+j >= finiteStateTotal {
				return m, 0, errors.New("bad sum")
			}
			m.lookup[int(cum+j)] = byte(s)
		}
		cum += f
		if f > 0 {
			m.count++
		}
	}
	if cum != finiteStateTotal {
		return m, 0, errors.New("sum mismatch")
	}
	return m, p, nil
}
func finiteStateEncode(data []byte, m *finiteStateModel) ([]byte, error) {
	x := uint32(finiteStateLower)
	tmp := make([]byte, 0, len(data))
	for i := len(data) - 1; i >= 0; i-- {
		s := data[i]
		f := m.freq[s]
		if f == 0 {
			return nil, errors.New("unmodeled symbol")
		}
		limit := f << 19
		for x >= limit {
			tmp = append(tmp, byte(x))
			x >>= 8
		}
		x = ((x / f) << finiteStateScale) + (x % f) + m.cum[s]
	}
	b := make([]byte, 4, len(tmp)+4)
	binary.LittleEndian.PutUint32(b, x)
	for i := len(tmp) - 1; i >= 0; i-- {
		b = append(b, tmp[i])
	}
	return b, nil
}
func finiteStateDecode(buf []byte, m *finiteStateModel, length int) ([]byte, error) {
	if len(buf) < 4 {
		return nil, errors.New("short input")
	}
	x := binary.LittleEndian.Uint32(buf)
	if x < finiteStateLower {
		return nil, errors.New("invalid initial state")
	}
	at := 4
	out := make([]byte, length)
	for i := 0; i < length; i++ {
		v := x & (finiteStateTotal - 1)
		s := m.lookup[v]
		out[i] = s
		f := m.freq[s]
		x = f*(x>>finiteStateScale) + v - m.cum[s]
		for x < finiteStateLower {
			if at >= len(buf) {
				return nil, errors.New("ran out of bits")
			}
			x = (x << 8) | uint32(buf[at])
			at++
		}
	}
	if at != len(buf) || x != finiteStateLower {
		return nil, fmt.Errorf("trailing state at=%d len=%d x=%d", at, len(buf), x)
	}
	return out, nil
}

// packFiniteState supplements only selected JLZ streams. The size comparison
// charges both the normalized model and the complete native stream frame.
// Other paths retain the historical raw/RLE/Huffman encoders.
func packFiniteState(data []byte) []byte {
	best := packStream(data)
	if len(data) < 64 {
		return best
	}
	m, err := finiteStateModelFor(data)
	if err != nil {
		return best
	}
	body, err := finiteStateEncode(data, &m)
	if err != nil {
		return best
	}
	table := finiteStateTable(&m)
	payloadLen := len(table) + len(body)
	candidateLen := varSize(uint64(len(data))) + varSize(uint64(payloadLen)) + 1 + payloadLen
	if candidateLen >= len(best) {
		return best
	}
	out := make([]byte, 0, candidateLen)
	out = binary.AppendUvarint(out, uint64(len(data)))
	out = binary.AppendUvarint(out, uint64(payloadLen))
	out = append(out, 5)
	out = append(out, table...)
	out = append(out, body...)
	return out
}
