package main

import (
	"math/rand"
	"testing"
)

// Match costs must remain identical after the hot-path scoring rewrite.
func TestOptimizedMatchScoreEquivalence(t *testing.T) {
	rng := rand.New(rand.NewSource(20261018))
	for i := 0; i < 100000; i++ {
		c := matchChoice{
			length:   5 + rng.Intn(matchMaximum-4),
			distance: 1 + rng.Intn(matchMaxDistance),
			typ:      byte(rng.Intn(3)),
			repeat:   -1,
		}
		recent := [3]int{rng.Intn(matchMaxDistance + 1), rng.Intn(matchMaxDistance + 1), rng.Intn(matchMaxDistance + 1)}
		if i%4 == 0 {
			recent[i%3] = c.distance
		}
		if i%7 == 0 {
			recent = [3]int{c.distance, c.distance, 0}
		}
		cost, rep := matchCost(c, recent)
		bits := matchLiteralValueBits
		if c.typ == 0 && c.distance <= matchLocalCopyWindow {
			bits = matchLocalCopyValueBits
		}
		expected := c.length*bits - cost*8
		got := scoreMatchChoice(c, recent)
		if got.score != expected || got.repeat != rep {
			t.Fatalf("iteration %d: got %+v; expected score %d repeat %d", i, got, expected, rep)
		}
	}
}
