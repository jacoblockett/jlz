package main

import (
	"math/rand"
	"testing"
)

func TestMathGateDistinguishesTransformedCopies(t *testing.T) {
	rng := rand.New(rand.NewSource(20261017))
	seed := make([]byte, 64*1024)
	if _, err := rng.Read(seed); err != nil {
		t.Fatal(err)
	}
	random := make([]byte, 256*1024)
	if _, err := rng.Read(random); err != nil {
		t.Fatal(err)
	}
	if likelyTransformedMatches(random) {
		t.Fatal("random-only block needlessly enables mathematical matching")
	}
	transformed := append([]byte(nil), seed...)
	for _, mask := range []byte{29, 77, 131} {
		block := append([]byte(nil), seed...)
		for j := range block {
			block[j] ^= mask
		}
		transformed = append(transformed, block...)
	}
	if !likelyTransformedMatches(transformed) {
		t.Fatal("failed to enable matching for distant transformed copies")
	}
}
