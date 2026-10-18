package main

// Expands already validated backward references. Exact copies use Go's bulk
// copy primitive while honoring overlapping LZ semantics. XOR/ADD copies
// retain forward byte recursion, which is required when source overlaps.
func expandMatch(out []byte, at, distance, length int, kind, mask byte) {
	if kind == 0 {
		if length <= distance {
			copy(out[at:at+length], out[at-distance:at-distance+length])
			return
		}
		copied := copy(out[at:at+length], out[at-distance:at])
		for copied < length {
			run := min(copied, length-copied)
			copy(out[at+copied:at+copied+run], out[at:at+run])
			copied += run
		}
		return
	}
	if kind == 1 {
		for i := 0; i < length; i++ {
			out[at+i] = out[at+i-distance] ^ mask
		}
	} else {
		for i := 0; i < length; i++ {
			out[at+i] = out[at+i-distance] + mask
		}
	}
}
