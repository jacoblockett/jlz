package jlz

import (
	"bytes"
	"crypto/sha256"
	"math/rand"
	"testing"
)

func makeMixedExperiment() []byte {
	r := rand.New(rand.NewSource(20261015))
	randBlock := func() []byte { b := make([]byte, 1024); r.Read(b); return b }
	var blocks [][]byte
	// Noisy independent base blocks and interleaved high-entropy XOR dependencies.
	for i := 0; i < 8; i++ {
		a, b := randBlock(), randBlock()
		blocks = append(blocks, a, b, xor(a, b))
	}
	// Local recurrence, not a repeating byte substring: the 256-byte initial
	// state is random, and each further 256-byte segment adds 1, 3, or 7.
	for i := 0; i < 12; i++ {
		a := randBlock()
		for j := 256; j < 1024; j++ {
			a[j] = a[j-256] + byte(1+(j/256)*2)
		}
		blocks = append(blocks, a)
	}
	// Long-distance equal, constant XOR, and sparse changed blocks.
	seed := randBlock()
	blocks = append(blocks, seed)
	for i := 0; i < 12; i++ {
		d := bytes.Clone(seed)
		if i%3 == 0 {
			for j := range d {
				d[j] ^= byte(i*9 + 1)
			}
		}
		if i%3 == 1 {
			for j := 0; j < 1024; j += 151 {
				d[j] ^= 0xc7
			}
		}
		blocks = append(blocks, d)
	}
	// Some entirely unrelated blocks, to test the fallback behavior.
	for i := 0; i < 12; i++ {
		blocks = append(blocks, randBlock())
	}
	return bytes.Join(blocks, nil)
}

func TestHybridMixedSelectionAndRoundtrip(t *testing.T) {
	data := makeMixedExperiment()
	arc, detail, err := encodeHybrid(data, sha256.Sum256(data), 1024)
	if err != nil {
		t.Fatal(err)
	}
	restored, err := decompress(arc, defaultLimit)
	if err != nil || !bytes.Equal(restored, data) {
		t.Fatalf("hybrid decode: %v", err)
	}
	if !bytes.Contains([]byte(detail), []byte("local-sub=")) || !bytes.Contains([]byte(detail), []byte("xor-pair=")) {
		t.Fatalf("expected both local predictor and XOR graph in one stream: %s", detail)
	}
	previous, _, err := compress(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(arc) >= len(previous) {
		t.Fatalf("hybrid didn't improve over competing-only MDC2 encodings: hybrid=%d earlier=%d, %s", len(arc), len(previous), detail)
	}
}
func TestHybridDirectModesAndMalformed(t *testing.T) {
	rng := rand.New(rand.NewSource(314159))
	for _, size := range []int{256, 257, 513, 2048, 4096, 11001} {
		b := make([]byte, size)
		rng.Read(b)
		for _, w := range []int{128, 256, 512, 1024} {
			arc, err := testHybridDirect(b, w)
			if err != nil {
				t.Fatal(err)
			}
			got, err := decompress(arc, defaultLimit)
			if err != nil || !bytes.Equal(got, b) {
				t.Fatalf("size=%d width=%d err=%v", size, w, err)
			}
		}
	}
	b := makeMixedExperiment()
	arc, err := testHybridDirect(b, 1024)
	if err != nil {
		t.Fatal(err)
	}
	for _, idx := range []int{0, 4, 6, len(arc) / 2, len(arc) - 1} {
		broken := bytes.Clone(arc)
		broken[idx] ^= 0xff
		if _, err := decompress(broken, defaultLimit); err == nil {
			t.Fatalf("corrupt archive accepted byte %d", idx)
		}
	}
	for _, cut := range []int{0, 5, 40, len(arc) / 2, len(arc) - 1} {
		if _, err := decompress(arc[:cut], defaultLimit); err == nil {
			t.Fatalf("truncated archive accepted at %d", cut)
		}
	}
}
