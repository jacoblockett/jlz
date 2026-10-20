package engine

import (
	"bytes"
	"crypto/sha256"
	"testing"
)

// Extended BWT blocks must remain independent and correctly positioned even
// when more than four 512 KiB blocks span multiple decoding worker groups.
func TestExtendedBWTGroupOffsets(t *testing.T) {
	plain := bytes.Repeat([]byte("a long and varied book with repeated passages. "), 77000)
	plain = plain[:5*2*stage3BlockSize+123]
	digest := sha256.Sum256(plain)
	arc, err := encodeParallelBWTRegions(plain, digest, len(plain)+4096, mathRegionSamples{})
	if err != nil || len(arc) == 0 {
		t.Fatalf("encode extended BWT: %v", err)
	}
	restored, err := Decompress(arc, MaxInputSize)
	if err != nil || !bytes.Equal(restored, plain) {
		t.Fatalf("extended BWT exact restoration: %v", err)
	}
	bad := append([]byte(nil), arc...)
	bad[4] = 255
	if _, err := Decompress(bad, MaxInputSize); err == nil {
		t.Fatal("decoder accepted invalid archive mode")
	}
	// A 512 KiB BWT block must never be misinterpreted as a legal
	// 512 KiB cross-history math block, whose limit stays 256 KiB.
	firstModeAt := 5 + varSize(uint64(len(plain))) + 32
	firstMode := arc[firstModeAt]
	if firstMode == blockBWTBijective || firstMode == blockBWTSubregions {
		other := append([]byte(nil), arc...)
		other[firstModeAt] = blockCrossPacked
		if _, err := Decompress(other, MaxInputSize); err == nil {
			t.Fatal("decoder accepted overlong cross-history block")
		}
	}
}
