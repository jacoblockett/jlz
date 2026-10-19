package engine

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestJointTokenLengthsRoundTrip(t *testing.T) {
	// A deliberately heterogeneous command stream with long literal runs,
	// short copies, and repeating history offsets exercises both length models.
	input := make([]byte, 0, 160*1024)
	rng := rand.New(rand.NewSource(17))
	base := make([]byte, 1024)
	_, _ = rng.Read(base)
	for i := 0; i < 140; i++ {
		variation := make([]byte, 64+i%53)
		_, _ = rng.Read(variation)
		input = append(input, variation...)
		input = append(input, base...)
	}

	tokens, err := tokenizeMatchV2(input)
	if err != nil {
		t.Fatal(err)
	}
	st, err := splitMatchTokens(tokens, len(input))
	if err != nil {
		t.Fatal(err)
	}
	op, dist, trans := packStream(st.op), packVarintStream(st.distances), packStream(st.transforms)
	header, ok := packSplitLengthHeader(st, op, dist, trans)
	if !ok {
		t.Fatal("no eligible split length streams")
	}
	packed := append(header, packLiteralStream(st.literals)...)
	out, err := decodePackedWithLengthMode(bytes.NewReader(packed), uint64(len(input)), false, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(input, out) {
		t.Fatal("split format corrupts original bytes")
	}
	for _, truncated := range [][]byte{packed[:1], packed[:len(packed)/2]} {
		if _, err := decodePackedWithLengthMode(bytes.NewReader(truncated), uint64(len(input)), false, nil, true); err == nil {
			t.Fatal("accepted truncated split stream")
		}
	}
}
