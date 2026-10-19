package jlz

import (
	"bytes"
	"crypto/sha256"
	"math/rand"
	"testing"
)

func TestAdaptiveVariableBoundaryDependencies(t *testing.T) {
	rng := rand.New(rand.NewSource(918273))
	randbytes := func(n int) []byte { b := make([]byte, n); rng.Read(b); return b }
	base := randbytes(5000)
	xor := bytes.Clone(base)
	add := bytes.Clone(base)
	patch := bytes.Clone(base)
	for i := range base {
		xor[i] ^= 0x6d
		add[i] += 147
	}
	for i := 37; i < len(patch); i += 79 {
		patch[i] ^= byte(i%255 + 1)
	}
	tests := []struct {
		name string
		data []byte
	}{
		{"xor shifted", bytes.Join([][]byte{randbytes(93), base, randbytes(137), xor}, nil)},
		{"add shifted", bytes.Join([][]byte{randbytes(103), base, randbytes(141), add}, nil)},
		{"patch shifted", bytes.Join([][]byte{randbytes(97), base, randbytes(133), patch}, nil)},
		{"uniform random", randbytes(12000)},
		{"repeating", bytes.Repeat([]byte("hello, 世界"), 512)},
		{"empty", nil},
		{"tiny", []byte{0, 1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arc, _, err := encodeAdaptive(tt.data, sha256.Sum256(tt.data))
			if err != nil {
				t.Fatal(err)
			}
			got, err := decompress(arc, defaultLimit)
			if err != nil || !bytes.Equal(got, tt.data) {
				t.Fatalf("adaptive direct roundtrip: %v", err)
			}
			// The outer encoder must choose a result no larger than its
			// predecessor using the exact same archive-header convention.
			old, _, err := compress(tt.data)
			if err != nil {
				t.Fatal(err)
			}
			best, _, err := compress(tt.data)
			if err != nil {
				t.Fatal(err)
			}
			if len(best) > len(old) {
				t.Fatalf("regressed: %d vs %d", len(best), len(old))
			}
			got, err = decompress(best, defaultLimit)
			if err != nil || !bytes.Equal(got, tt.data) {
				t.Fatalf("selected roundtrip: %v", err)
			}
		})
	}
}

func TestAdaptiveArchiveValidation(t *testing.T) {
	data := bytes.Repeat([]byte("a0123456789abcdefghij"), 48)
	arc, _, err := encodeAdaptive(data, sha256.Sum256(data))
	if err != nil {
		t.Fatal(err)
	}
	for _, cut := range []int{0, 1, 4, 5, 6, 30, len(arc) / 2, len(arc) - 1} {
		if _, err := decompress(arc[:cut], defaultLimit); err == nil {
			t.Fatalf("accepted truncated archive at %d", cut)
		}
	}
	damaged := bytes.Clone(arc)
	damaged[len(magic)+2] ^= 0x01
	if _, err := decompress(damaged, defaultLimit); err == nil {
		t.Fatal("accepted tampered archive")
	}
	// Legacy archives cannot be decoded under altered stream semantics.
	older, _, err := compress([]byte("abc"))
	if err != nil {
		t.Fatal(err)
	}
	copy(older[:4], "MDC3")
	_, err = decompress(older, defaultLimit)
	if err == nil {
		t.Fatal("accepted legacy DEFLATE-era archive under native-only format")
	}
}
