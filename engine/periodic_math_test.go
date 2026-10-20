package engine

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestStronglyAligned16(t *testing.T) {
	if stronglyAligned16(nil) || stronglyAligned16(make([]byte, 4096)) {
		t.Fatal("accepted undersized source")
	}
	r := rand.New(rand.NewSource(938))
	raw := make([]byte, 256*1024)
	if _, err := r.Read(raw); err != nil {
		t.Fatal(err)
	}
	if stronglyAligned16(raw) {
		t.Fatal("random bytes should not be periodic")
	}
	record := append([]byte(nil), raw...)
	for i := range record {
		if i%16 < 13 {
			record[i] = byte((i % 16) + 128)
		}
	}
	if !stronglyAligned16(record) {
		t.Fatal("missed strong 16-byte field signature")
	}
	old := append([]byte(nil), record...)
	mode, coded := betterPeriodicLane(record, blockNativeBytes, old)
	if mode != blockLanePacked {
		t.Fatal("record fields not encoded")
	}
	decoded, err := decodeLaneBlock(coded, len(record))
	if err != nil || !bytes.Equal(decoded, record) {
		t.Fatalf("record bytes not preserved: %v", err)
	}
	mode, coded = betterPeriodicLane(raw, blockNativeBytes, raw)
	if mode != blockNativeBytes || !bytes.Equal(coded, raw) {
		t.Fatal("random source altered")
	}
}
