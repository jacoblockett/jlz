package engine

// Portable reversible relative-operand normalization. The encoder transforms
// bytes after E8/E9 markers and skips the four operand bytes. The decoder
// follows identical boundaries on reconstructed transformed data, so arbitrary
// non-executable inputs remain exactly reversible.
import "encoding/binary"

const blockBranchBWT byte = 19

func shouldBranchNormalize(input []byte) bool {
	if len(input) < 64*1024 { return false }
	n := 0
	for i := 0; i+5 <= len(input); i++ {
		if input[i] == 0xe8 || input[i] == 0xe9 {
			n++
			i += 4
		}
	}
	return n*100 >= 2*len(input)
}

func normalizeBranchRelative(input []byte, reverse bool) []byte {
	out := append([]byte(nil), input...)
	for i := 0; i+5 <= len(out); i++ {
		if out[i] != 0xe8 && out[i] != 0xe9 { continue }
		v := binary.LittleEndian.Uint32(out[i+1:])
		offset := uint32(i+5)
		if reverse { v -= offset } else { v += offset }
		binary.LittleEndian.PutUint32(out[i+1:], v)
		i += 4
	}
	return out
}
