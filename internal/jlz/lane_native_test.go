package jlz

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/rand"
	"testing"
)

func TestLaneNativeRoundtrip(t *testing.T) {
	rng := rand.New(rand.NewSource(25))
	for _, n := range []int{1, 2, 3, 5, 7, 16, 31, 1024, stage3BlockSize} {
		b := make([]byte, n)
		rng.Read(b)
		for _, stride := range []int{2, 4, 8} {
			encoded := encodeLaneBlock(b, stride)
			d, e := decodeLaneBlock(encoded, len(b))
			if e != nil || !bytes.Equal(d, b) {
				t.Fatalf("lane size=%d stride=%d: %v", n, stride, e)
			}
			for _, damaged := range [][]byte{nil, {1}, {3}, {9}} {
				if _, e := decodeLaneBlock(damaged, len(b)); e == nil {
					t.Fatalf("accepted corrupt lane metadata")
				}
			}
		}
	}
}

// A complete mixed-cross archive must preserve exact bytes, original digest,
// and original backward-block framing while decoding the new lane blocks.
func TestLaneArchiveIntegrity(t *testing.T) {
	data := make([]byte, 3*stage3BlockSize)
	rng := rand.New(rand.NewSource(125))
	rng.Read(data)
	for i := 0; i < len(data); i += 4 {
		data[i] = 0
		data[i+1] = 0x5a
	}
	archive, e := encodeLaneArchive(data, sha256.Sum256(data), len(data)*2, 4)
	if e != nil || len(archive) == 0 {
		t.Fatalf("failed to build lane archive %v", e)
	}
	restored, e := decompress(archive, defaultLimit)
	if e != nil || !bytes.Equal(restored, data) {
		t.Fatalf("lane archive decode %v", e)
	}
	bad := append([]byte(nil), archive...)
	bad[len(bad)-1] ^= 0x40
	if _, e := decompress(bad, defaultLimit); e == nil {
		t.Fatal("accepted corrupt lane archive")
	}
}

func TestLaneSharedHistogramMatchesIndependentCounts(t *testing.T) {
	for _, n := range []int{1, 7, 512, 65536, stage3BlockSize} {
		block := make([]byte, n)
		rand.New(rand.NewSource(int64(n))).Read(block)
		got := laneEntropyEstimates(block)
		for j, stride := range []int{2, 4, 8} {
			var freq [8][256]int
			var totals [8]int
			for i, b := range block {
				freq[i%stride][b]++
				totals[i%stride]++
			}
			want := 0.0
			for lane := 0; lane < stride; lane++ {
				total := float64(totals[lane])
				for _, c := range freq[lane] {
					if c != 0 {
						want += float64(c) * math.Log2(total/float64(c))
					}
				}
			}
			if math.Abs(got[j]-want) > 1e-6 {
				t.Fatalf("n=%d stride=%d got %.10f want %.10f", n, stride, got[j], want)
			}
		}
	}
}
