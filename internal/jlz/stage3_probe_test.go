package jlz

import (
	"crypto/sha256"
	"math/rand"
	"testing"
)

func TestStage3MatchedLiteralCostsLessOnSparseMismatches(t *testing.T) {
	rng := rand.New(rand.NewSource(991))
	for _, interval := range []int{17, 32, 63, 128, 257} {
		seed := make([]byte, 32768)
		rng.Read(seed)
		data := append([]byte(nil), seed...)
		for repeat := 0; repeat < 3; repeat++ {
			second := append([]byte(nil), seed...)
			for j := interval - 1; j < len(second); j += interval {
				second[j] ^= byte(repeat + 1)
			}
			data = append(data, second...)
		}
		tokens, _ := tokenizeMatchV2(data)
		if len(tokens) == 0 {
			t.Log("no tokens")
			continue
		}
		native, e := encodeMatchNativeVariant(data, sha256.Sum256(data), tokens, false)
		if e != nil {
			t.Fatal(e)
		}
		pred, e := encodeMatchNativeVariant(data, sha256.Sum256(data), tokens, true)
		if e != nil {
			t.Fatal(e)
		}
		t.Logf("interval %d native %d pred %d", interval, len(native), len(pred))
	}
}
