package engine

import (
	"bytes"
	"crypto/sha256"
	"math/rand"
	"testing"
)

func TestParallelNativeIntegrity(t *testing.T) {
	r := rand.New(rand.NewSource(78623))
	data := make([]byte, 5*1024*1024+123)
	pattern := make([]byte, 4096)
	if _, err := r.Read(pattern); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(data); i++ {
		data[i] = pattern[i%len(pattern)]
	}
	for i := 0; i < len(data); i += 12347 {
		data[i] ^= byte(i)
	}
	digest := sha256.Sum256(data)
	for _, workers := range []int{1, 2, 4} {
		arc, err := encodeParallelNative(data, digest, workers)
		if err != nil {
			t.Fatal(err)
		}
		if arc[4] != modeParallelNative {
			t.Fatal("wrong mode")
		}
		restored, err := decompress(arc, int64(len(data)))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(data, restored) {
			t.Fatal("round trip mismatch")
		}
		copyArc, err := encodeParallelNative(data, digest, workers)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(arc, copyArc) {
			t.Fatal("parallel result not deterministic")
		}
		bad := bytes.Clone(arc)
		bad[len(bad)-4] ^= 0x80
		if _, err := decompress(bad, int64(len(data))); err == nil {
			t.Fatal("corruption not rejected")
		}
		if _, err := decompress(arc[:len(arc)-2], int64(len(data))); err == nil {
			t.Fatal("truncation not rejected")
		}
	}
}

func TestParallelNativeIncompressible(t *testing.T) {
	r := rand.New(rand.NewSource(19583))
	data := make([]byte, parallelNativeBlockSize+4096)
	if _, err := r.Read(data); err != nil {
		t.Fatal(err)
	}
	arc, err := encodeParallelNative(data, sha256.Sum256(data), 3)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := decompress(arc, int64(len(data)))
	if err != nil || !bytes.Equal(data, decoded) {
		t.Fatalf("random roundtrip: %v", err)
	}
}
