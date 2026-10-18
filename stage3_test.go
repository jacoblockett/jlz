package main

import (
	"bytes"
	"crypto/sha256"
	"math/rand"
	"testing"
)

func TestStage3MatchedLiteralAndOldNative(t *testing.T) {
	rng := rand.New(rand.NewSource(169))
	seed := make([]byte, 8192)
	rng.Read(seed)
	sample := append([]byte(nil), seed...)
	// Large transformed references punctuated by unmatched bytes.
	for variant := 0; variant < 4; variant++ {
		block := append([]byte(nil), seed...)
		for i := range block {
			if variant%2 == 0 {
				block[i] ^= 0x6f
			} else {
				block[i] += 0x35
			}
		}
		for _, at := range []int{71, 173, 492, 1047, 2050, 3303, 5900} {
			block[at] ^= byte(variant + 1)
		}
		sample = append(sample, block...)
	}
	tokens, err := tokenizeMatchV2(sample)
	if err != nil || len(tokens) == 0 {
		t.Fatalf("no candidate: %v", err)
	}
	for _, predict := range []bool{false, true} {
		archive, err := encodeMatchNativeVariant(sample, sha256.Sum256(sample), tokens, predict)
		if err != nil {
			t.Fatal(err)
		}
		restored, err := decompress(archive, defaultLimit)
		if err != nil || !bytes.Equal(restored, sample) {
			t.Fatalf("native pred %t: %v", predict, err)
		}
	}
}

func TestStage3MixedRoundtripAndTampering(t *testing.T) {
	rng := rand.New(rand.NewSource(90))
	a := bytes.Repeat([]byte("stage 3 text file line with structure\n"), 6000)
	b := make([]byte, stage3BlockSize)
	rng.Read(b)
	// Deliberately mixes entropy and structured segments.
	data := append(append([]byte(nil), a...), b...)
	data = append(data, bytes.Repeat([]byte{0x44}, stage3BlockSize)...)
	arc, err := encodeMixedBlocks(data, sha256.Sum256(data), len(data)*2)
	if err != nil || len(arc) == 0 {
		t.Fatalf("mixed encode %v", err)
	}
	if arc[4] != modeMixed {
		t.Fatal("not mixed mode")
	}
	restored, err := decompress(arc, defaultLimit)
	if err != nil || !bytes.Equal(data, restored) {
		t.Fatalf("mixed round trip %v", err)
	}
	for _, mutate := range []func([]byte) []byte{
		func(a []byte) []byte { return a[:len(a)-1] },
		func(a []byte) []byte { a[len(a)-1] ^= 0x08; return a },
		func(a []byte) []byte { a[39] = 0xff; return a },
	} {
		corrupted := mutate(append([]byte(nil), arc...))
		if _, err = decompress(corrupted, defaultLimit); err == nil {
			t.Fatal("corrupt mixed archive accepted")
		}
	}
}

func TestStage3NoRegressionOnSmallInput(t *testing.T) {
	samples := [][]byte{nil, []byte("z"), bytes.Repeat([]byte("abbreviated text"), 800), bytes.Repeat([]byte("XYZ"), 15000)}
	for _, src := range samples {
		archive, _, err := compress(src)
		if err != nil {
			t.Fatal(err)
		}
		result, err := decompress(archive, defaultLimit)
		if err != nil || !bytes.Equal(src, result) {
			t.Fatalf("input size %d: %v", len(src), err)
		}
	}
}

func TestStage3MixedCanSelectDistinctPredictionBlocks(t *testing.T) {
	rng := rand.New(rand.NewSource(20261))
	var input []byte
	for _, lag := range []int{1, 1, 16, 16} {
		part := make([]byte, stage3BlockSize)
		for i := range part {
			if i < lag {
				part[i] = byte(rng.Intn(256))
			} else {
				part[i] = part[i-lag] + byte(rng.Intn(3)-1)
			}
		}
		input = append(input, part...)
	}

	archive, s, err := compress(input)
	if err != nil {
		t.Fatal(err)
	}
	if archive[4] != modeMixed {
		t.Fatalf("expected profitable block representation; got %s (%d bytes)", s.method, len(archive))
	}
	recovered, err := decompress(archive, defaultLimit)
	if err != nil || !bytes.Equal(input, recovered) {
		t.Fatalf("block prediction roundtrip: %v", err)
	}
}
