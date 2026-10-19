package main

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"testing"
)

func TestLengthGrammarStreamRoundTrip(t *testing.T) {
	rng := rand.New(rand.NewSource(19))
	samples := [][]byte{nil, {1}, bytes.Repeat([]byte{1, 2, 3, 1, 2, 3}, 10000)}
	for _, n := range []int{2, 2048, 8192, 65536, 262144} {
		x := make([]byte, n)
		_, _ = rng.Read(x)
		samples = append(samples, x)
	}
	for _, input := range samples {
		encoded := packLengthStream(input)
		reader := bytes.NewReader(encoded)
		out, err := unpackStream(reader, len(input))
		if err != nil || reader.Len() != 0 || !bytes.Equal(out, input) {
			t.Fatalf("length=%d: inverse %v", len(input), err)
		}
	}
	// Force a grammar-coded stream independently of its production gate.
	input := bytes.Repeat([]byte{1, 2, 3, 1, 2, 3}, 10000)
	encoded := packLengthGrammar(input, len(input)*4)
	if len(encoded) == 0 {
		t.Fatal("failed to build grammar")
	}
	for i := 0; i < len(encoded); i++ {
		reader := bytes.NewReader(encoded[:i])
		got, err := unpackStream(reader, len(input))
		if err == nil && bytes.Equal(got, input) {
			t.Fatalf("accepted truncated grammar at offset %d", i)
		}
	}
	// Guard against both forward generator references and recursive grammar.
	at := lenVarAt(encoded)
	if at < 0 {
		t.Fatal("invalid test fixture")
	}
	bad := append([]byte(nil), encoded...)
	bad[at+1] = 255
	reader := bytes.NewReader(bad)
	if _, err := unpackStream(reader, len(input)); err == nil {
		t.Fatal("accepted invalid generator")
	}
}
func lenVarAt(stream []byte) int {
	_, n := binary.Uvarint(stream)
	if n <= 0 { return -1 }
	_, m := binary.Uvarint(stream[n:])
	if m <= 0 || stream[n+m] != lengthGrammarCodec { return -1 }
	return n+m+1+1+1
}
