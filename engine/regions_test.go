package engine

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestRegionalMathSelectedAndRestores(t *testing.T) {
	rng := rand.New(rand.NewSource(20261016))
	seed := make([]byte, 8192)
	if _, err := rng.Read(seed); err != nil {
		t.Fatal(err)
	}
	data := make([]byte, 3*1024*1024)
	for p := 0; p < len(data); p += 8192 {
		mask := byte((p / 8192) * 13)
		for j := 0; j < 8192 && p+j < len(data); j++ {
			data[p+j] = seed[j] ^ mask
		}
	}
	archive, _, err := compress(data)
	if err != nil {
		t.Fatal(err)
	}
	if archive[4] != modeMixedCross && archive[4] != modeMixed {
		t.Fatalf("regional matching not reached in default encoder: mode=%d", archive[4])
	}
	if len(archive) >= len(data)/2 {
		t.Fatalf("regional matching did not produce substantial reduction: %d vs %d", len(archive), len(data))
	}
	decoded, err := decompress(archive, defaultLimit)
	if err != nil || !bytes.Equal(decoded, data) {
		t.Fatalf("regional decoding failed: %v", err)
	}
	duplicate, _, err := compress(data)
	if err != nil || !bytes.Equal(duplicate, archive) {
		t.Fatalf("nondeterministic archive: %v", err)
	}
	corrupted := bytes.Clone(archive)
	corrupted[len(corrupted)-1] ^= 0x17
	if _, err := decompress(corrupted, defaultLimit); err == nil {
		t.Fatal("corrupted region archive accepted")
	}
}

func TestSkewedNativeHuffmanNoRawFallback(t *testing.T) {
	data := bytes.Repeat([]byte{'a'}, 262144)
	for b := 0; b < 256; b++ {
		data[b*1009] = byte(b)
	}
	lengths, ok := nativeLengths(data)
	if !ok {
		t.Fatal("skewed natural distribution still rejected by Huffman coder")
	}
	for _, l := range lengths {
		if l > 15 {
			t.Fatalf("overlong Huffman code: %d", l)
		}
	}
	pack := packStream(data)
	if len(pack) >= len(data)/5 {
		t.Fatalf("length-limited Huffman ineffective: %d", len(pack))
	}
	decoded, err := nativeUnpackBytes(pack, len(data))
	if err != nil || !bytes.Equal(decoded, data) {
		t.Fatalf("Huffman restoration failed: %v", err)
	}
}
