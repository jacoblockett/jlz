package match

import "encoding/binary"

// Fixed-width fingerprint mixing for native exact and transformed matching.
const matchHashBits = 16
func Fingerprint(data []byte, at int, typ byte) uint32 {
	if typ == 0 {
		return binary.LittleEndian.Uint32(data[at:]) ^ uint32(data[at+4])*0x9e3779b1
	}
	a := data[at]
	b, c, d, e := data[at+1], data[at+2], data[at+3], data[at+4]
	if typ == 1 {
		return uint32(b^a) | uint32(c^a)<<8 | uint32(d^a)<<16 | uint32(e^a)<<24
	}
	return uint32(b-a) | uint32(c-a)<<8 | uint32(d-a)<<16 | uint32(e-a)<<24
}

func Hash(key uint32) uint32 {
	key ^= key >> 15
	key *= 0x2c1b3c6d
	key ^= key >> 12
	return key >> (32 - matchHashBits)
}

