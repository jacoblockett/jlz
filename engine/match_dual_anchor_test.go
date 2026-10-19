package engine

import (
	"bytes"
	"crypto/sha256"
	"testing"
)

func TestDualAnchorExact(t *testing.T) {
	blocks := [][]byte{
		bytes.Repeat([]byte("this phrase repeats, this phrase repeats;"), 100),
		bytes.Repeat([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}, 300),
		append(bytes.Repeat([]byte("abcdefghijklmnop"), 128), bytes.Repeat([]byte("zyxwvutsrqponmlk"), 256)...),
	}
	for i, b := range blocks {
		tokens, e := tokenizeDualAnchor(b, 0)
		if e != nil || len(tokens) == 0 {
			t.Fatalf("parse %d %v", i, e)
		}
		arc, e := encodeMatchV2Tokens(b, sha256.Sum256(b), tokens)
		if e != nil {
			t.Fatal(e)
		}
		got, e := decompress(arc, defaultLimit)
		if e != nil || !bytes.Equal(got, b) {
			t.Fatalf("whole %d %v", i, e)
		}
		history := append([]byte(nil), b[:len(b)/3]...)
		block := b[len(b)/3:]
		whole := append(append([]byte(nil), history...), block...)
		tokens, e = tokenizeDualAnchor(whole, len(history))
		if e != nil || len(tokens) == 0 {
			t.Fatalf("cross parse %d %v", i, e)
		}
		payload, e := packCrossRegionTokens(block, tokens)
		if e != nil {
			t.Fatal(e)
		}
		got, e = decodeMatchPackedHistory(bytes.NewReader(payload), uint64(len(block)), false, history)
		if e != nil || !bytes.Equal(got, block) {
			t.Fatalf("cross decoded %d %v", i, e)
		}
	}
}
