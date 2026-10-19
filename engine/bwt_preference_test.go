package engine

import (
	"bytes"
	"testing"
)

// A BWT preference probe must measure and retain the very same representation
// the full BWT path emits. The earlier legacy MTF sampler understated its gain
// and could reject an archive that was actually smaller.
func TestBWTPreferenceUsesActualBijectiveEncoding(t *testing.T) {
	input := make([]byte, 3*stage3BlockSize)
	for i := range input {
		input[i] = byte("the quick brown fox jumps over the lazy dog. "[i%45])
	}
	var regional mathRegionSamples
	regional.positions = [3]int{0, 1, 2}
	for i := range regional.results {
		regional.results[i] = mathRegionResult{mode: stage3Raw, payload: input[i*stage3BlockSize : (i+1)*stage3BlockSize]}
	}
	win, samples, err := sampleBWTPreference(input, regional)
	if err != nil || !win {
		t.Fatalf("BWT probe: win=%v err=%v", win, err)
	}
	for i, r := range samples.results {
		if r.mode != blockBWTBijective {
			t.Fatalf("sample %d: mode %d", i, r.mode)
		}
		source := input[i*stage3BlockSize : (i+1)*stage3BlockSize]
		expected, err := encodeBWTBijective(source)
		if err != nil || !bytes.Equal(expected, r.payload) {
			t.Fatalf("sample %d not reused bijective payload: %v", i, err)
		}
		restored, err := decodeBWTBijective(r.payload, len(source))
		if err != nil || !bytes.Equal(source, restored) {
			t.Fatalf("sample %d not exact: %v", i, err)
		}
	}
}
