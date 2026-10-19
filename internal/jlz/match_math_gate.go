package jlz

// Sparse probes across the whole region detect nonzero constant XOR/ADD
// relationships even when sources and targets are far apart. If none is
// detected, avoid allocating and searching the two transformed indexes.
// This is a speed heuristic, never an exact proof of absence.
func likelyTransformedMatches(data []byte) bool {
	if len(data) < 16*1024 {
		return true
	}
	const tableBits = 16
	const minVerified = 16
	var last [2][1 << tableBits]uint32
	for pos := 0; pos+minVerified <= len(data); pos += 8 {
		for typ := byte(1); typ <= 2; typ++ {
			h := matchHash(matchFingerprint(data, pos, typ))
			past := last[typ-1][h]
			last[typ-1][h] = uint32(pos + 1)
			if past == 0 {
				continue
			}
			from := int(past) - 1
			if pos-from < minVerified {
				continue
			}
			var mask byte
			if typ == 1 {
				mask = data[pos] ^ data[from]
			} else {
				mask = data[pos] - data[from]
			}
			if mask == 0 {
				continue
			}
			found := true
			for k := 1; k < minVerified; k++ {
				v := data[from+k]
				if typ == 1 {
					v ^= mask
				} else {
					v += mask
				}
				if v != data[pos+k] {
					found = false
					break
				}
			}
			if found {
				return true
			}
		}
	}
	return false
}
