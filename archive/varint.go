package archive

import (
 "encoding/binary"
 "io"
)

// Unsigned varint framing shared by all JLZ archive modes.
func PutVar(w io.Writer, n uint64) error {
	var b [binary.MaxVarintLen64]byte
	k := binary.PutUvarint(b[:], n)
	_, err := w.Write(b[:k])
	return err
}
func ReadVar(r io.ByteReader) (uint64, error) { return binary.ReadUvarint(r) }
func VarSize(n uint64) int {
	count := 1
	for n >= 128 {
		n >>= 7
		count++
	}
	return count
}

