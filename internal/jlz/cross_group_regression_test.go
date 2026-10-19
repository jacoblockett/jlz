package jlz

import (
	"bytes"
	"crypto/sha256"
	"math/rand"
	"testing"
)

func TestCrossGroupReferenceAndCorruption(t *testing.T) {
	rng := rand.New(rand.NewSource(17))
	seed := make([]byte, 128*1024)
	if _, err := rng.Read(seed); err != nil {
		t.Fatal(err)
	}
	data := make([]byte, 5*stage3BlockSize)
	for i := 0; i < len(data); i += len(seed) {
		copy(data[i:], seed)
	}
	// Confirm that a cross-region reference can be emitted and decoded using
	// reconstructed prefix bytes, without transmitting a dictionary.
	history := data[:crossHistorySize]
	current := data[crossHistorySize : crossHistorySize+stage3BlockSize]
	tokenInput := append(bytes.Clone(history), current...)
	tokens, err := tokenizeMatchV2Range(tokenInput, len(history))
	if err != nil {
		t.Fatal(err)
	}
	payload, err := packCrossRegionTokens(current, tokens)
	if err != nil {
		t.Fatal(err)
	}
	restored, err := decodeMatchPackedHistory(bytes.NewReader(payload), uint64(len(current)), false, history)
	if err != nil || !bytes.Equal(restored, current) {
		t.Fatalf("cross match restoration: %v", err)
	}
	if len(payload) >= len(current)/2 {
		t.Fatalf("cross references ineffective: %d", len(payload))
	}
	// The cross parser must not accept a truncated command stream.
	if _, err := decodeMatchPackedHistory(bytes.NewReader(payload[:len(payload)/2]), uint64(len(current)), false, history); err == nil {
		t.Fatal("truncated cross reference accepted")
	}
}

func TestCrossGroupOldArchiveDecoding(t *testing.T) {
	rng := rand.New(rand.NewSource(2026))
	data := make([]byte, 600000)
	rng.Read(data)
	copy(data[400000:], data[100000:300000])
	digest := sha256.Sum256(data)
	arc, err := encodeMatchV2(data, digest)
	if err != nil {
		t.Fatal(err)
	}
	if len(arc) == 0 {
		t.Skip("matcher declined fixture")
	}
	dec, err := decompress(arc, defaultLimit)
	if err != nil || !bytes.Equal(dec, data) {
		t.Fatalf("v2 archive compatibility: %v", err)
	}
}
