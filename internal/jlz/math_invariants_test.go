package jlz

import (
	"bytes"
	"crypto/sha256"
	"math/rand"
	"testing"
)

func TestArithmeticCodecLossless(t *testing.T) {
	rng := rand.New(rand.NewSource(20261015))
	for _, n := range []int{0, 1, 2, 3, 7, 8, 31, 255, 1024, 4096, 16384} {
		sources := [][]byte{make([]byte, n), bytes.Repeat([]byte("a"), n)}
		r := make([]byte, n)
		if _, e := rng.Read(r); e != nil {
			t.Fatal(e)
		}
		sources = append(sources, r)
		for k, src := range sources {
			packed := arithmeticEncode(src)
			recovered, e := arithmeticDecode(packed, len(src))
			if e != nil || !bytes.Equal(src, recovered) {
				t.Fatalf("codec length %d case %d: %v", n, k, e)
			}
		}
	}
}

func TestJLZPredictiveAndMosaic(t *testing.T) {
	rng := rand.New(rand.NewSource(1229))
	random := make([]byte, 8192)
	rng.Read(random)
	sources := [][]byte{random, bytes.Repeat([]byte("this is an example; abcdefgh\n"), 300)}
	combined := append([]byte(nil), random...)
	for _, v := range random {
		combined = append(combined, v+17)
	}
	sources = append(sources, combined)
	for _, src := range sources {
		h := sha256.Sum256(src)
		for _, which := range []string{"predictive", "mosaic", "adaptive"} {
			var arc []byte
			var e error
			switch which {
			case "predictive":
				arc, _, e = encodePredictive(src, h)
			case "mosaic":
				arc, _, e = encodeMosaic(src, h)
			case "adaptive":
				arc, _, e = encodeAdaptive(src, h)
			}
			if e != nil {
				t.Fatalf("%s compress: %v", which, e)
			}
			if arc == nil {
				t.Fatalf("%s unexpectedly skipped", which)
			}
			out, e := decompress(arc, defaultLimit)
			if e != nil || !bytes.Equal(src, out) {
				t.Fatalf("%s roundtrip: %v", which, e)
			}
			arc[len(arc)-1] ^= 0x40
			if _, e := decompress(arc, defaultLimit); e == nil {
				t.Errorf("%s failed to detect corruption", which)
			}
		}
	}
}
func TestMosaicMixedMethods(t *testing.T) {
	rng := rand.New(rand.NewSource(888))
	src := make([]byte, 16*1024)
	rng.Read(src[:8192])
	for i := 8192; i < len(src); i++ {
		src[i] = src[i-8192] + 1
	}
	h := sha256.Sum256(src)
	arc, _, e := encodeMosaic(src, h)
	if e != nil {
		t.Fatal(e)
	}
	out, e := decompress(arc, defaultLimit)
	if e != nil || !bytes.Equal(src, out) {
		t.Fatalf("mixed region loss: %v", e)
	}
}
