package main

import (
	"bytes"
	"crypto/sha256"
	"math/rand"
	"testing"
)

func TestStructuralMatchGatePreservesLongRangeCandidates(t *testing.T) {
	r := rand.New(rand.NewSource(20261016))
	seed := make([]byte, 128*1024)
	if _, err := r.Read(seed); err != nil {
		t.Fatal(err)
	}
	data := make([]byte, 0, 8*len(seed))
	for pass := 0; pass < 8; pass++ {
		for _, b := range seed {
			data = append(data, b^byte(pass*17))
		}
	}
	if !stage3StructuralSampleUseful(data) {
		t.Fatal("missed repeated distant XOR structure")
	}
	// Ensure the explicit research override remains a complete lossless path.
	archived, _, err := compressFast(data)
	if err != nil {
		t.Fatal(err)
	}
	restored, err := decompress(archived, int64(len(data)))
	if err != nil || !bytes.Equal(restored, data) {
		t.Fatalf("lossless roundtrip: %v", err)
	}
	if sha256.Sum256(restored) != sha256.Sum256(data) {
		t.Fatal("hash mismatch")
	}
}
