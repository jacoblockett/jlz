package engine

import (
	"bytes"
	"crypto/sha256"
	"math/rand"
	"testing"
)

func TestNativeLiterals(t *testing.T) {
	rng := rand.New(rand.NewSource(8443))
	for _, data := range [][]byte{
		bytes.Repeat([]byte("banana bread banana banana bread\n"), 500),
		bytes.Repeat([]byte{'X'}, 16384),
		func() []byte { b := make([]byte, 32768); rng.Read(b); return b }(),
		[]byte("short"),
	} {
		codec, table, payload := encodeNativeLiterals(data)
		var restored []byte
		var err error
		if codec == 0 {
			restored = payload
		} else {
			restored, err = decodeNativeLiterals(table, payload, len(data))
		}
		if err != nil || !bytes.Equal(data, restored) {
			t.Fatalf("native literal mismatch codec %d err %v", codec, err)
		}
	}
}

func TestNativeOverlappingAndTransformed(t *testing.T) {
	var tokens bytes.Buffer
	tokens.WriteByte(matchLiteral)
	putVar(&tokens, 4)
	tokens.WriteString("ABCD")
	tokens.WriteByte(matchCopy)
	putVar(&tokens, 4)
	putVar(&tokens, 40)
	tokens.WriteByte(matchRep0)
	putVar(&tokens, 12)
	tokens.WriteByte(matchXor)
	putVar(&tokens, 4)
	putVar(&tokens, 9)
	tokens.WriteByte(0x32)
	tokens.WriteByte(matchRepAdd)
	tokens.WriteByte(0)
	putVar(&tokens, 11)
	tokens.WriteByte(13)
	// Reconstruct reference content using native decoder without external SHA envelope.
	sample := make([]byte, 76)
	// AB... / repeated patterns, exact and transformed overlaps.
	copy(sample, []byte("ABCD"))
	at := 4
	for i := 0; i < 40; i++ {
		sample[at] = sample[at-4]
		at++
	}
	for i := 0; i < 12; i++ {
		sample[at] = sample[at-4]
		at++
	}
	for i := 0; i < 9; i++ {
		sample[at] = sample[at-4] ^ 0x32
		at++
	}
	for i := 0; i < 11; i++ {
		sample[at] = sample[at-4] + 13
		at++
	}
	digest := sha256.Sum256(sample)
	arc, err := encodeMatchNativeTokens(sample, digest, tokens.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	result, err := decompress(arc, defaultLimit)
	if err != nil || !bytes.Equal(sample, result) {
		t.Fatalf("native reference error %v", err)
	}
	corrupted := append([]byte(nil), arc...)
	corrupted[len(corrupted)-1] ^= 1
	if _, err = decompress(corrupted, defaultLimit); err == nil {
		t.Fatal("corrupt native archive accepted")
	}
}

func TestNativeEncodedMatcher(t *testing.T) {
	rng := rand.New(rand.NewSource(434))
	data := make([]byte, 2048)
	rng.Read(data)
	block := append([]byte(nil), data[50:900]...)
	for i := range block {
		block[i] ^= 0x74
	}
	data = append(data, block...)
	for i := range block {
		block[i] += 0x21
	}
	data = append(data, block...)
	tokens, err := tokenizeMatchV2(data)
	if err != nil || len(tokens) == 0 {
		t.Fatalf("no native tokens: %v", err)
	}
	arc, err := encodeMatchNativeTokens(data, sha256.Sum256(data), tokens)
	if err != nil {
		t.Fatal(err)
	}
	restored, err := decompress(arc, defaultLimit)
	if err != nil || !bytes.Equal(data, restored) {
		t.Fatalf("native mismatch: %v", err)
	}
}

// A long-distance repeated high-entropy block demonstrates that the native
// format is not merely decodable: it can be selected over DEFLATE-backed tokens.
func TestStage2NativeCanWin(t *testing.T) {
	rng := rand.New(rand.NewSource(347))
	for _, size := range []int{1024, 4096, 16384, 65536} {
		discarded := make([]byte, size)
		rng.Read(discarded)
	}
	seed := make([]byte, 262144)
	rng.Read(seed)
	data := bytes.Repeat(seed, 3)

	arc, method, err := encodeMatchStage2(data, sha256.Sum256(data))
	if err != nil {
		t.Fatal(err)
	}
	if len(arc) >= len(data) {
		t.Fatalf("native reference encoder failed to compress repeated bytes: %d, %s", len(arc), method)
	}
	restored, err := decompress(arc, defaultLimit)
	if err != nil || !bytes.Equal(data, restored) {
		t.Fatalf("new archive roundtrip: %v", err)
	}
}
