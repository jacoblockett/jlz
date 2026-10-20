package engine

// Cheap source-level probe for a 16-byte record layout. Unlike exploring
// every possible period, this requires both substantial same-field similarity
// and significantly weaker eight-byte similarity, so ordinary repeats are
// unlikely to trigger an additional full native entropy encoding.
func stronglyAligned16(block []byte) bool {
	if len(block) < 16*1024 {
		return false
	}
	var hit16, hit8, pairs int
	for _, start := range [3]int{0, len(block)/2 - 1024, len(block) - 2048} {
		s := block[start : start+2048]
		for i := 16; i < len(s); i++ {
			if s[i] == s[i-16] {
				hit16++
			}
			if s[i] == s[i-8] {
				hit8++
			}
			pairs++
		}
	}
	return hit16*100 > 65*pairs && (hit16-hit8)*100 > 17*pairs
}

// Reuses only established mode13 framing; no new archive mode is needed.
// The current regional result remains the fallback and exactly priced winner.
func betterPeriodicLane(block []byte, mode byte, current []byte) (byte, []byte) {
	if (mode != blockNativeBytes && mode != blockPredictive) || !stronglyAligned16(block) {
		return mode, current
	}
	alternative := encodeLaneBlock(block, 16)
	if len(alternative) >= len(current) {
		return mode, current
	}
	mode, current = blockLanePacked, alternative
	// A three-record field grouping sometimes yields a better native codebook.
	// Test it only when field splitting has already beaten the incumbent.
	alternative = encodeLaneBlock(block, 48)
	if len(alternative) < len(current) {
		current = alternative
	}
	return mode, current
}
