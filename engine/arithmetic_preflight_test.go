package engine

import (
	"math/rand"
	"testing"
)

func TestMatchedLiteralPreflightDecision(t *testing.T) {
	small := make([]byte, 40000)
	// The histogram prefilter now skips even short equal-distribution
	// alternatives. This changes only encoder work, not restoration.
	if considerMatchedLiteralResidual(small, small) {
		t.Fatal("identical short distributions should not trigger extra encoding")
	}
	shortStrong := make([]byte, 40000)
	shortRandom := make([]byte, 40000)
	rand.New(rand.NewSource(21)).Read(shortRandom)
	if !considerMatchedLiteralResidual(shortRandom, shortStrong) {
		t.Fatal("clearly profitable short residual must be evaluated")
	}
	easy := make([]byte, 65536)
	random := make([]byte, 65536)
	r := rand.New(rand.NewSource(17))
	if _, err := r.Read(random); err != nil {
		t.Fatal(err)
	}
	if considerMatchedLiteralResidual(easy, random) {
		t.Fatal("clearly worse residual should be skipped")
	}
	if !considerMatchedLiteralResidual(random, easy) {
		t.Fatal("stronger residual candidate should be evaluated")
	}
}
