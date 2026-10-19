package entropy

import "errors"

// Adaptive context arithmetic coder used by JLZ's predictive and literal streams.
const (
	acMax          uint64 = 0xffffffff
	acHalf         uint64 = 0x80000000
	acQuarter      uint64 = 0x40000000
	acThreeQuarter uint64 = 0xc0000000
)

type bitCount struct{ zero, one uint16 }
type bitModel struct {
	base [256]bitCount
	prev [65536]bitCount
}

func normalizedCount(c bitCount) (uint32, uint32) {
	if c.zero == 0 && c.one == 0 {
		return 1, 1
	}
	return uint32(c.zero), uint32(c.one)
}
func (m *bitModel) probability(previous byte, prefix uint16) uint64 {
	global := &m.base[prefix]
	local := &m.prev[int(previous)*256+int(prefix)]
	z0, o0 := normalizedCount(*global)
	z1, o1 := normalizedCount(*local)
	// Counts rescale by 2048, so both quotients and weighted sum fit uint32.
	p0 := (z0 * 65536) / (z0 + o0)
	p1 := (z1 * 65536) / (z1 + o1)
	weight := uint32(1)
	if z1+o1 >= 10 {
		weight = 5
	}
	p := (p0*(8-weight) + p1*weight) / 8
	if p < 1 {
		p = 1
	}
	if p > 65535 {
		p = 65535
	}
	return uint64(p)
}
func addCount(c *bitCount, v byte) {
	if c.zero == 0 && c.one == 0 {
		c.zero = 1
		c.one = 1
	}
	if v == 0 {
		c.zero++
	} else {
		c.one++
	}
	if uint32(c.zero)+uint32(c.one) >= 2048 {
		c.zero = (c.zero + 1) / 2
		c.one = (c.one + 1) / 2
	}
}
func (m *bitModel) update(previous byte, prefix uint16, v byte) {
	addCount(&m.base[prefix], v)
	addCount(&m.prev[int(previous)*256+int(prefix)], v)
}

type bitWriter struct {
	out []byte
	cur byte
	n   uint8
}

func (w *bitWriter) bit(v byte) {
	w.cur = (w.cur << 1) | (v & 1)
	w.n++
	if w.n == 8 {
		w.out = append(w.out, w.cur)
		w.cur = 0
		w.n = 0
	}
}
func (w *bitWriter) finish() []byte {
	if w.n != 0 {
		w.out = append(w.out, w.cur<<(8-w.n))
	}
	return w.out
}

type bitReader struct {
	data    []byte
	offset  int
	current byte
}

func (r *bitReader) bit() uint64 {
	// A full byte is fetched once; the decoder specifies implicit zero padding.
	if r.offset&7 == 0 {
		at := r.offset >> 3
		r.current = 0
		if at < len(r.data) {
			r.current = r.data[at]
		}
	}
	r.offset++
	result := r.current >> 7
	r.current <<= 1
	return uint64(result)
}

type acEncoder struct {
	low, high uint64
	pending   int
	writer    bitWriter
}

func newACEncoder() *acEncoder { return &acEncoder{low: 0, high: acMax} }
func (c *acEncoder) emit(v byte) {
	c.writer.bit(v)
	for c.pending > 0 {
		c.writer.bit(v ^ 1)
		c.pending--
	}
}
func (c *acEncoder) bit(bit byte, p0 uint64) {
	cut := c.low + (((c.high - c.low + 1) * p0) >> 16) - 1
	if bit == 0 {
		c.high = cut
	} else {
		c.low = cut + 1
	}
	for {
		if c.high < acHalf {
			c.emit(0)
		} else if c.low >= acHalf {
			c.emit(1)
			c.low -= acHalf
			c.high -= acHalf
		} else if c.low >= acQuarter && c.high < acThreeQuarter {
			c.pending++
			c.low -= acQuarter
			c.high -= acQuarter
		} else {
			break
		}
		c.low <<= 1
		c.high = (c.high << 1) | 1
	}
}
func (c *acEncoder) finish() []byte {
	c.pending++
	if c.low < acQuarter {
		c.emit(0)
	} else {
		c.emit(1)
	}
	return c.writer.finish()
}
func EncodeArithmetic(input []byte) []byte {
	coder := newACEncoder()
	model := new(bitModel)
	var previous byte
	for _, b := range input {
		prefix := uint16(1)
		for bit := 7; bit >= 0; bit-- {
			v := (b >> uint(bit)) & 1
			coder.bit(v, model.probability(previous, prefix))
			model.update(previous, prefix, v)
			prefix = (prefix << 1) | uint16(v)
		}
		previous = b
	}
	return coder.finish()
}
func DecodeArithmetic(payload []byte, length int) ([]byte, error) {
	if length < 0 || (length > 0 && len(payload) == 0) {
		return nil, errors.New("invalid arithmetic payload")
	}
	reader := bitReader{data: payload}
	lo, hi := uint64(0), acMax
	var value uint64
	for i := 0; i < 32; i++ {
		value = (value << 1) | reader.bit()
	}
	model := new(bitModel)
	out := make([]byte, length)
	var previous byte
	for i := range out {
		prefix := uint16(1)
		var b byte
		for bit := 7; bit >= 0; bit-- {
			p0 := model.probability(previous, prefix)
			cut := lo + (((hi - lo + 1) * p0) >> 16) - 1
			var v byte
			if value <= cut {
				v = 0
				hi = cut
			} else {
				v = 1
				lo = cut + 1
			}
			for {
				if hi < acHalf {
				} else if lo >= acHalf {
					lo -= acHalf
					hi -= acHalf
					value -= acHalf
				} else if lo >= acQuarter && hi < acThreeQuarter {
					lo -= acQuarter
					hi -= acQuarter
					value -= acQuarter
				} else {
					break
				}
				lo <<= 1
				hi = (hi << 1) | 1
				value = (value << 1) | reader.bit()
			}
			model.update(previous, prefix, v)
			b = (b << 1) | v
			prefix = (prefix << 1) | uint16(v)
		}
		out[i] = b
		previous = b
	}
	return out, nil
}
