package engine

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestBijectiveBWTBlockCodecAndLegacy(t *testing.T) {
	samples := [][]byte{[]byte("banana"), bytes.Repeat([]byte("bwt suffix repeated xxyyzzyyzz"), 4000), bytes.Repeat([]byte{0, 255, 15, 15, 1, 1}, 4096)}
	random := make([]byte, 20000)
	rand.New(rand.NewSource(44)).Read(random)
	samples = append(samples, random)
	for _, src := range samples {
		encoded, err := encodeBWTBijective(src)
		if err != nil {
			t.Fatal(err)
		}
		restored, err := decodeBWTBijective(encoded, len(src))
		if err != nil || !bytes.Equal(src, restored) {
			t.Fatalf("bijective source=%d: %v", len(src), err)
		}
		legacy, err := encodeBWTNative(src)
		if err != nil {
			t.Fatal(err)
		}
		old, err := decodeBWTNative(legacy, len(src))
		if err != nil || !bytes.Equal(src, old) {
			t.Fatalf("legacy source=%d: %v", len(src), err)
		}
		if len(encoded) > 2 {
			if _, err := decodeBWTBijective(encoded[:len(encoded)-1], len(src)); err == nil {
				t.Fatalf("accepted truncation")
			}
		}
	}
}
