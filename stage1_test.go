package main

import (
	"bytes"
	"crypto/sha256"
	"io"
	"math/rand"
	"sync"
	"testing"
)

func TestStage1NativePackingDeterminism(t *testing.T) {
	rng := rand.New(rand.NewSource(2081))
	data := make([]byte, 128<<10)
	rng.Read(data)
	for _, source := range [][]byte{data, bytes.Repeat([]byte("abcdefgh12345678"), 8192), data[:200]} {
		packed, e := nativePackBytes(source)
		if e != nil {
			t.Fatal(e)
		}
		for k := 0; k < 3; k++ {
			again, e := nativePackBytes(source)
			if e != nil {
				t.Fatal(e)
			}
			if !bytes.Equal(packed, again) {
				t.Fatal("non-deterministic native encoding")
			}
			restored, e := nativeUnpackBytes(again, len(source))
			if e != nil || !bytes.Equal(restored, source) {
				t.Fatalf("native roundtrip: %v", e)
			}
		}
	}
}

func TestStage1RawFallback(t *testing.T) {
	data := make([]byte, 64<<10)
	rand.New(rand.NewSource(117)).Read(data)
	arc, summary, err := compressFast(data)
	if err != nil {
		t.Fatal(err)
	}
	if !summary.raw || arc[len(magic)] != modeRaw {
		t.Fatalf("expected raw fallback, got %s", summary.method)
	}
	restored, err := decompress(arc, defaultLimit)
	if err != nil || !bytes.Equal(restored, data) {
		t.Fatalf("raw roundtrip: %v", err)
	}
}

func TestStage1NativeConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	errs := make(chan error, 16)
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			original := bytes.Repeat([]byte{byte(i), byte(i + 1), byte(i + 5)}, 32768)
			archive, err := nativePackBytes(original)
			if err != nil {
				errs <- err
				return
			}
			reader := newNativeReader(bytes.NewReader(archive))
			result, err := io.ReadAll(reader)
			reader.Close()
			if err != nil || sha256.Sum256(result) != sha256.Sum256(original) {
				errs <- io.ErrUnexpectedEOF
			}
		}(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatal(err)
	}
}
