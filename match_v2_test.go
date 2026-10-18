package main

import (
	"bytes"
	"crypto/sha256"
	"math/rand"
	"testing"
)

func TestMatchV2OverlappingRepeatAndTransforms(t *testing.T) {
	var tok bytes.Buffer
	var want []byte
	literal := []byte("ABCD")
	tok.WriteByte(matchLiteral)
	putVar(&tok, uint64(len(literal)))
	tok.Write(literal)
	want = append(want, literal...)
	// Overlap a 4-byte source with a 40-byte target.
	tok.WriteByte(matchCopy)
	putVar(&tok, 4)
	putVar(&tok, 40)
	for i := 0; i < 40; i++ {
		want = append(want, want[len(want)-4])
	}
	tok.WriteByte(matchRep0)
	putVar(&tok, 12)
	for i := 0; i < 12; i++ {
		want = append(want, want[len(want)-4])
	}
	// Shifted transformed references can overlap as well.
	tok.WriteByte(matchXor)
	putVar(&tok, 4)
	putVar(&tok, 9)
	tok.WriteByte(0x32)
	for i := 0; i < 9; i++ {
		want = append(want, want[len(want)-4]^0x32)
	}
	tok.WriteByte(matchRepAdd)
	tok.WriteByte(0)
	putVar(&tok, 11)
	tok.WriteByte(13)
	for i := 0; i < 11; i++ {
		want = append(want, want[len(want)-4]+13)
	}
	digest := sha256.Sum256(want)
	var arc bytes.Buffer
	arc.WriteString(magic)
	arc.WriteByte(modeMatchV2)
	putVar(&arc, uint64(len(want)))
	arc.Write(digest[:])
	arc.WriteByte(0)
	putVar(&arc, uint64(tok.Len()))
	arc.Write(tok.Bytes())
	got, err := decompress(arc.Bytes(), defaultLimit)
	if err != nil || !bytes.Equal(got, want) {
		t.Fatalf("bad v2 reconstruction: %v", err)
	}
	damaged := append([]byte(nil), arc.Bytes()...)
	damaged[len(damaged)-1] ^= 1
	if _, err := decompress(damaged, defaultLimit); err == nil {
		t.Fatal("damaged archive accepted")
	}
}

func TestMatchV2FullSizeAndLegacyPreserved(t *testing.T) {
	rng := rand.New(rand.NewSource(12550))
	prefix := make([]byte, 650*1024)
	rng.Read(prefix)
	copyBlock := append([]byte(nil), prefix[199123:199123+32768]...)
	data := append(prefix, copyBlock...)
	archive, err := encodeMatchV2(data, sha256.Sum256(data))
	if err != nil || len(archive) == 0 {
		t.Fatalf("no v2 archive: %v", err)
	}
	if len(archive) >= len(data)-20000 {
		t.Fatalf("missed distant match: %d", len(archive))
	}
	decoded, err := decompress(archive, defaultLimit)
	if err != nil || !bytes.Equal(data, decoded) {
		t.Fatalf("v2 large mismatch: %v", err)
	}
	// An old decoder mode remains accepted by the current decoder.
	old, err := fastAdaptive(data[:256*1024], sha256.Sum256(data[:256*1024]))
	if err != nil {
		t.Fatal(err)
	}
	if len(old) > 0 {
		legacy, err := decompress(old, defaultLimit)
		if err != nil || !bytes.Equal(legacy, data[:256*1024]) {
			t.Fatalf("legacy mismatch: %v", err)
		}
	}
}
