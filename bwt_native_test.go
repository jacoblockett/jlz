package main

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestBWTNativeRoundtrip(t *testing.T) {
	samples := [][]byte{[]byte("banana"), []byte("hello world"), bytes.Repeat([]byte("abcdef"), 1024), bytes.Repeat([]byte{0, 255, 1, 1, 3}, 8192)}
	random := make([]byte, 16000)
	rand.New(rand.NewSource(21)).Read(random)
	samples = append(samples, random)
	for _, original := range samples {
		arc, e := encodeBWTNative(original)
		if e != nil {
			t.Fatal(e)
		}
		restored, e := decodeBWTNative(arc, len(original))
		if e != nil || !bytes.Equal(restored, original) {
			t.Fatalf("BWT roundtrip %d bytes err=%v", len(original), e)
		}
		for _, damage := range [][]byte{nil, {0}, {255, 255, 255, 255, 255, 255, 255}} {
			if _, err := decodeBWTNative(damage, len(original)); err == nil {
				t.Fatalf("accepted invalid BWT metadata")
			}
		}
	}
}
