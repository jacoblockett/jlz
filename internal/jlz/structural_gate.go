package jlz

// Cheap, distributed probe for transformations that could be useful beyond
// a former 32 KiB matching window. Unlike the old mixed gate, this can admit
// files whose winning blocks are match-coded, not predictor-coded. The probe
// cannot authorize a mode or claim savings; it only authorizes an actual
// whole-archive comparison. All counts and alignments are encoder-only.
func stage3StructuralSampleUseful(input []byte) bool {
	if len(input) < 2*stage3BlockSize {
		return false
	}
	for _, at := range [...]int{0, (len(input) / 2 / stage3BlockSize) * stage3BlockSize, ((len(input) - 1) / stage3BlockSize) * stage3BlockSize} {
		end := min(len(input), at+stage3BlockSize)
		if stage3StructuralEvidence(input[at:end]) {
			return true
		}
	}
	return false
}

func stage3StructuralEvidence(data []byte) bool {
	const stride = 32
	const distance = 32 * 1024
	const minSavingsEvidence = 20 // strict percentage; weak signals waste full mixed passes
	if len(data) < 2*distance {
		return false
	}
	type probeKey struct {
		kind  byte
		value uint32
	}
	seen := make(map[probeKey]int, len(data)/stride*3)
	matches, samples := 0, 0
	for i := 0; i+16 <= len(data); i += stride {
		samples++
		matched := false
		for kind := byte(0); kind < 3; kind++ {
			k := probeKey{kind, matchFingerprint(data, i, kind)}
			prev, ok := seen[k]
			seen[k] = i
			if !ok || i-prev < distance {
				continue
			}
			mask := byte(0)
			if kind == 1 {
				mask = data[i] ^ data[prev]
			}
			if kind == 2 {
				mask = data[i] - data[prev]
			}
			valid := true
			for j := 0; j < 16; j++ {
				value := data[prev+j]
				if kind == 1 {
					value ^= mask
				} else if kind == 2 {
					value += mask
				}
				if value != data[i+j] {
					valid = false
					break
				}
			}
			if valid {
				matched = true
				break
			}
		}
		if matched {
			matches++
		}
	}
	return matches*100 >= samples*minSavingsEvidence
}
