package main

// MDC5: reversible byte-predictor transforms followed by an adaptive,
// order-0/order-1 context arithmetic coder. The coder is standalone and
// has no third-party dependencies. Only an actual size win is retained.

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

const modePredictive byte = 7
const (
	predNone        byte = 0
	predXor         byte = 1
	predSub         byte = 2
	codecNative     byte = 0
	codecArithmetic byte = 1
)

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

func normalizedCount(c bitCount) (uint64, uint64) {
	if c.zero == 0 && c.one == 0 {
		return 1, 1
	}
	return uint64(c.zero), uint64(c.one)
}
func (m *bitModel) probability(previous byte, prefix uint16) uint64 {
	global := &m.base[prefix]
	local := &m.prev[int(previous)*256+int(prefix)]
	z0, o0 := normalizedCount(*global)
	z1, o1 := normalizedCount(*local)
	p0 := (z0 * 65536) / (z0 + o0)
	p1 := (z1 * 65536) / (z1 + o1)
	weight := uint64(1)
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
	return p
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
	data   []byte
	offset int
}

func (r *bitReader) bit() uint64 {
	pos := r.offset
	r.offset++
	if pos/8 >= len(r.data) {
		return 0
	}
	return uint64((r.data[pos/8] >> (7 - uint(pos%8))) & 1)
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
func arithmeticEncode(input []byte) []byte {
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
func arithmeticDecode(payload []byte, length int) ([]byte, error) {
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
func transform(input []byte, typ byte, lag int) []byte {
	out := make([]byte, len(input))
	if typ == predNone {
		copy(out, input)
		return out
	}
	copy(out[:lag], input[:lag])
	for i := lag; i < len(input); i++ {
		if typ == predXor {
			out[i] = input[i] ^ input[i-lag]
		} else {
			out[i] = input[i] - input[i-lag]
		}
	}
	return out
}
func reverseTransform(input []byte, typ byte, lag int) []byte {
	out := append([]byte(nil), input...)
	if typ == predNone {
		return out
	}
	for i := lag; i < len(input); i++ {
		if typ == predXor {
			out[i] ^= out[i-lag]
		} else {
			out[i] += out[i-lag]
		}
	}
	return out
}
func nativeBytes(b []byte) ([]byte, error)          { return nativePackBytes(b) }
func unnativeBytes(b []byte, n int) ([]byte, error) { return nativeUnpackBytes(b, n) }
func encodePredictive(input []byte, digest [32]byte) ([]byte, string, error) {
	bestSize := int(^uint(0) >> 1)
	var best []byte
	bestName := ""
	types := []struct {
		typeID byte
		lag    int
	}{
		{predNone, 0}, {predXor, 1}, {predSub, 1},
		{predXor, 2}, {predSub, 2}, {predXor, 4}, {predSub, 4},
		{predXor, 8}, {predSub, 8},
	}
	for _, t := range types {
		if t.lag >= len(input) {
			continue
		}
		data := transform(input, t.typeID, t.lag)
		// Test both the simple entropy baseline and a context arithmetic coder.
		for _, codec := range []byte{codecArithmetic, codecNative} {
			var encoded []byte
			var e error
			if codec == codecArithmetic {
				encoded = arithmeticEncode(data)
			} else {
				encoded, e = nativeBytes(data)
				if e != nil {
					return nil, "", e
				}
			}
			var b bytes.Buffer
			b.WriteString(magic)
			b.WriteByte(modePredictive)
			putVar(&b, uint64(len(input)))
			b.Write(digest[:])
			b.WriteByte(t.typeID)
			putVar(&b, uint64(t.lag))
			b.WriteByte(codec)
			b.Write(encoded)
			if b.Len() < bestSize {
				best = b.Bytes()
				bestSize = b.Len()
				bestName = fmt.Sprintf("transform=%d lag=%d codec=%d", t.typeID, t.lag, codec)
			}
		}
	}
	return best, bestName, nil
}
func decodePredictive(r *bytes.Reader, n uint64) ([]byte, error) {
	typ, e := r.ReadByte()
	if e != nil {
		return nil, e
	}
	lag, e := readVar(r)
	if e != nil {
		return nil, e
	}
	codec, e := r.ReadByte()
	if e != nil {
		return nil, e
	}
	if typ > predSub || (typ != predNone && (lag == 0 || lag >= n)) || (typ == predNone && lag != 0) {
		return nil, errors.New("bad predictive transform")
	}
	data := make([]byte, r.Len())
	if _, e = io.ReadFull(r, data); e != nil {
		return nil, e
	}
	var decoded []byte
	switch codec {
	case codecArithmetic:
		decoded, e = arithmeticDecode(data, int(n))
	case codecNative:
		decoded, e = unnativeBytes(data, int(n))
	default:
		return nil, errors.New("unknown predictive codec")
	}
	if e != nil {
		return nil, e
	}
	return reverseTransform(decoded, typ, int(lag)), nil
}
