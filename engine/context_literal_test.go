package engine

import (
	"bytes"
	"crypto/sha256"
	"io"
	"math/rand"
	"os"
	"testing"
)

func TestNativeContextLiteralCodec(t *testing.T) {
	sources := [][]byte{
		bytes.Repeat([]byte("a long symbolic sequence at a distance: "), 500),
		bytes.Repeat([]byte{0, 1, 5, 2, 0, 1, 5, 2, 0, 1, 7}, 700),
		make([]byte, 16384),
	}
	rng := rand.New(rand.NewSource(42))
	rng.Read(sources[2])
	for i, data := range sources {
		packed := packLiteralStream(data)
		got, err := unpackStream(bytes.NewReader(packed), len(data))
		if err != nil || !bytes.Equal(data, got) {
			t.Fatalf("literal %d: %v", i, err)
		}
	}
	// Make the arithmetic path mandatory and verify both the decoder and
	// existing integrity envelope reject truncated or modified payloads.
	data := bytes.Repeat([]byte("algorithmic conditional coding context!"), 1000)
	st := packLiteralStream(data)
	r := bytes.NewReader(st)
	_, _ = readVar(r)
	_, _ = readVar(r)
	kind, _ := r.ReadByte()
	if kind != 3 {
		t.Fatalf("arithmetic codec not selected, got %d", kind)
	}
	original, err := unpackStream(bytes.NewReader(st), len(data))
	if err != nil || !bytes.Equal(data, original) {
		t.Fatal(err)
	}
	input := bytes.Repeat([]byte("ABCABCABCABCABC 1234567890 xyz "), 1024)
	toks, err := tokenizeMatchV2(input)
	if err != nil || len(toks) == 0 {
		t.Fatal(err)
	}
	for _, pred := range []bool{false, true} {
		arc, err := encodeMatchPackedVariant(input, sha256.Sum256(input), toks, pred)
		if err != nil {
			t.Fatal(err)
		}
		decoded, err := decompress(arc, defaultLimit)
		if err != nil || !bytes.Equal(input, decoded) {
			t.Fatalf("roundtrip %v %v", pred, err)
		}
	}
}

func TestContextLiteralActualArchiveAndCorruption(t *testing.T) {
	input, err := os.ReadFile("fixtures/04_sparse_mutations.bin")
	if err != nil {
		t.Fatal(err)
	}
	compressed, _, err := compress(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(compressed) < 40 || compressed[4] != modeMatchPackedPred {
		t.Fatalf("expected winning packed prediction mode, got %d", compressed[4])
	}
	// Confirm the *selected final archive* actually uses arithmetic literals.
	r := bytes.NewReader(compressed[5:])
	_, err = readVar(r)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = r.Seek(32, io.SeekCurrent); err != nil {
		t.Fatal(err)
	}
	_, err = readVar(r)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 4; i++ {
		if _, err = unpackStream(r, len(input)*8+1024); err != nil {
			t.Fatal(err)
		}
	}
	_, err = readVar(r)
	if err != nil {
		t.Fatal(err)
	}
	_, err = readVar(r)
	if err != nil {
		t.Fatal(err)
	}
	codec, err := r.ReadByte()
	if err != nil {
		t.Fatal(err)
	}
	if codec != 3 {
		t.Fatalf("winning archive literal codec=%d, not context arithmetic", codec)
	}
	restored, err := decompress(compressed, defaultLimit)
	if err != nil || !bytes.Equal(restored, input) {
		t.Fatalf("selected mode roundtrip %v", err)
	}
	mutant := bytes.Clone(compressed)
	mutant[len(mutant)-20] ^= 0x40
	restored, err = decompress(mutant, defaultLimit)
	if err == nil || bytes.Equal(restored, input) {
		t.Fatal("corruption undetected")
	}
	truncated := compressed[:len(compressed)-20]
	if _, err = decompress(truncated, defaultLimit); err == nil {
		t.Fatal("truncation undetected")
	}
}
