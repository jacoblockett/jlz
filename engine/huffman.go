package engine

// Canonical Huffman decoder with a lookup table sized to the longest code
// present in each stream, and a bulk word fast path for interior bytes.
// Both preserve the existing bitstream and bound code lengths to 15 bits.
// Invalid prefixes, truncation, padding, and trailing data remain errors.

import (
	"encoding/binary"
	"errors"
	"io"
)

func decodeCanonicalHuffman(lengths [256]byte, payload []byte, count int) ([]byte, error) {
	codes, _, _, _, _, err := nativeCodes(lengths)
	if err != nil {
		return nil, err
	}
	maxBits := 0
	for _, length := range lengths {
		if int(length) > maxBits {
			maxBits = int(length)
		}
	}
	if maxBits == 0 {
		return nil, errors.New("empty canonical Huffman table")
	}
	lookup := make([]uint16, 1<<uint(maxBits)) // length in high byte, symbol in low byte
	for sym, nbits := range lengths {
		if nbits == 0 {
			continue
		}
		span := 1 << uint(maxBits-int(nbits))
		start := int(codes[sym]) * span
		entry := uint16(nbits)<<8 | uint16(sym)
		for i := 0; i < span; i++ {
			lookup[start+i] = entry
		}
	}
	output := make([]byte, count)
	pos := 0
	total := len(payload) * 8
	for i := range output {
		if pos >= total {
			return nil, io.ErrUnexpectedEOF
		}
		at := pos >> 3
		var bits uint32
		if at+4 <= len(payload) {
			bits = binary.BigEndian.Uint32(payload[at:]) >> 8
		} else {
			bits = uint32(payload[at]) << 16
			if at+1 < len(payload) {
				bits |= uint32(payload[at+1]) << 8
			}
			if at+2 < len(payload) {
				bits |= uint32(payload[at+2])
			}
		}
		key := (bits >> uint(24-maxBits-(pos&7))) & uint32((1<<uint(maxBits))-1)
		entry := lookup[key]
		length := int(entry >> 8)
		if length == 0 {
			return nil, errors.New("invalid canonical Huffman prefix")
		}
		if pos+length > total {
			return nil, io.ErrUnexpectedEOF
		}
		output[i] = byte(entry)
		pos += length
	}
	if (pos+7)/8 != len(payload) {
		return nil, errors.New("canonical Huffman trailing bytes")
	}
	if pos%8 != 0 && payload[len(payload)-1]&byte((1<<uint(8-pos%8))-1) != 0 {
		return nil, errors.New("nonzero canonical Huffman padding")
	}
	return output, nil
}
