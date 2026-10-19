package main

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// The integer field is a little-endian 16-bit word. Model each original word
// as the previous reconstructed word plus a modular difference in Z/(2^16).
// Signed zigzag maps small positive/negative differences to low nonnegative
// values; the independent low/high streams reuse JLZ's native entropy coder.
// A payload always starts with byte 4, then two framed packed streams, then
// the optional unmatched last byte. The metadata byte rejects other variants
// experimented with locally but never admitted to JLZ's stable bitstream.
func encodeAdditiveZigzagWords(data []byte) []byte {
	n := len(data) / 2
	lo := make([]byte, n)
	hi := make([]byte, n)
	var prev uint16
	for i := 0; i < n; i++ {
		x := binary.LittleEndian.Uint16(data[2*i:])
		d := x
		if i > 0 {
			d -= prev
		}
		z := (d << 1) ^ uint16(int16(d)>>15)
		lo[i], hi[i] = byte(z), byte(z>>8)
		prev = x
	}
	var out bytes.Buffer
	out.WriteByte(4)
	out.Write(packStream(lo))
	out.Write(packStream(hi))
	if len(data)%2 != 0 {
		out.WriteByte(data[len(data)-1])
	}
	return out.Bytes()
}

func decodeAdditiveWords(data []byte, size int) ([]byte, error) {
	if len(data) == 0 || data[0] != 4 {
		return nil, errors.New("invalid additive block header")
	}
	r := bytes.NewReader(data[1:])
	n := size / 2
	lo, e := unpackStream(r, n)
	if e != nil {
		return nil, e
	}
	hi, e := unpackStream(r, n)
	if e != nil {
		return nil, e
	}
	out := make([]byte, size)
	var prev uint16
	for i := 0; i < n; i++ {
		z := uint16(lo[i]) | uint16(hi[i])<<8
		d := (z >> 1) ^ uint16(-int16(z&1))
		v := d
		if i > 0 {
			v += prev
		}
		binary.LittleEndian.PutUint16(out[2*i:], v)
		prev = v
	}
	if size%2 != 0 {
		v, e := r.ReadByte()
		if e != nil {
			return nil, e
		}
		out[size-1] = v
	}
	if r.Len() != 0 {
		return nil, errors.New("additive block trailing bytes")
	}
	return out, nil
}
