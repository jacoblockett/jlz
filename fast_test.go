package main

import (
	"bytes"
	"crypto/sha256"
	"math/rand"
	"testing"
)

func TestFastArchiveRoundTrip(t *testing.T) {
	r := rand.New(rand.NewSource(1967))
	src := make([]byte, 8192)
	r.Read(src)
	base := make([]byte, 4096)
	r.Read(base)
	xor := append([]byte(nil), base...)
	add := append([]byte(nil), base...)
	for i := range base {
		xor[i] ^= 0x6d
		add[i] += 73
	}
	cases := [][]byte{
		nil, []byte("a"), bytes.Repeat([]byte("abca"), 0),
		bytes.Repeat([]byte("foobar\n"), 8192), src,
		bytes.Join([][]byte{src[:17], base, src[17:49], xor, src[49:73], add}, nil),
	}
	for i, data := range cases {
		arc, _, err := compress(data)
		if err != nil {
			t.Fatalf("case %d compress: %v", i, err)
		}
		out, err := decompress(arc, defaultLimit)
		if err != nil || !bytes.Equal(data, out) {
			t.Fatalf("case %d roundtrip: %v", i, err)
		}
		arc[len(arc)-1] ^= 1
		if _, err := decompress(arc, defaultLimit); err == nil {
			t.Fatalf("case %d corruption undetected", i)
		}
	}
}

func TestFastAdaptiveTransformedMatches(t *testing.T) {
	r := rand.New(rand.NewSource(2026))
	base := make([]byte, 4096)
	r.Read(base)
	for _, typ := range []byte{0, 1, 2} {
		copy2 := append([]byte(nil), base...)
		for i := range copy2 {
			switch typ {
			case 1:
				copy2[i] ^= 0x77
			case 2:
				copy2[i] += 37
			}
		}
		src := bytes.Join([][]byte{base, []byte("shifted insert"), copy2}, nil)
		arc, err := encodeMatchV2(src, sha256.Sum256(src))
		if err != nil || len(arc) == 0 {
			t.Fatalf("transform %d was not matched: %v", typ, err)
		}
		out, err := decompress(arc, defaultLimit)
		if err != nil || !bytes.Equal(src, out) {
			t.Fatalf("transform %d decode: %v", typ, err)
		}
		if len(arc) > len(src)/2+100 {
			t.Fatalf("transform %d poor match: %d", typ, len(arc))
		}
	}
}
