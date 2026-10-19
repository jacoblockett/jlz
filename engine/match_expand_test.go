package engine

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestExpandMatchAllOverlapAndTransforms(t *testing.T) {
	rnd := rand.New(rand.NewSource(31))
	for _, kind := range []byte{0, 1, 2} {
		for dist := 1; dist < 70; dist++ {
			for n := 1; n < 200; n++ {
				at := 70
				mask := byte(7)
				before := make([]byte, at+n)
				rnd.Read(before[:at])
				expected := bytes.Clone(before)
				for j := 0; j < n; j++ {
					v := expected[at+j-dist]
					if kind == 1 {
						v ^= mask
					} else if kind == 2 {
						v += mask
					}
					expected[at+j] = v
				}
				actual := bytes.Clone(before)
				expandMatch(actual, at, dist, n, kind, mask)
				if !bytes.Equal(actual, expected) {
					t.Fatalf("kind=%d distance=%d length=%d", kind, dist, n)
				}
			}
		}
	}
}
