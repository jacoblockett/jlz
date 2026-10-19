package engine

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	cases := [][]byte{
		nil,
		{0},
		{255},
		{0, 255, 0, 255},
		bytes.Repeat([]byte{0}, 128),
		bytes.Repeat([]byte{0xAA, 0x55}, 500),
		bytes.Repeat([]byte("ABCD-hello"), 350),
		bytes.Repeat([]byte("a long-distance repeated chunk and another block!"), 100),
	}
	rng := rand.New(rand.NewSource(12345))
	for _, n := range []int{1, 2, 3, 7, 16, 64, 128, 1024, 4096} {
		b := make([]byte, n)
		_, _ = rng.Read(b)
		cases = append(cases, b)
	}

	for i, in := range cases {
		compressed, _, err := compress(in)
		if err != nil {
			t.Fatalf("case %d compress: %v", i, err)
		}
		got, err := decompress(compressed, defaultLimit)
		if err != nil {
			t.Fatalf("case %d decompress: %v", i, err)
		}
		if !bytes.Equal(in, got) {
			t.Fatalf("case %d mismatch", i)
		}
	}
}

func TestRewriteNotAndXor(t *testing.T) {
	rules := []rule{
		{opConcat, 0, 0}, // id 2 = 00
		{opConcat, 1, 1}, // id 3 = 11 -> NOT(id2)
		{opConcat, 0, 1}, // id 4 = 01
		{opConcat, 1, 0}, // id 5 = 10 -> NOT(id4)
		{opConcat, 2, 3}, // id 6 = 0011
		{opConcat, 4, 4}, // id 7 = 0101
		{opConcat, 4, 5}, // id 8 = 0110 -> XOR(id6,id7)
	}
	original := make([]model, 2, len(rules)+2)
	original[0], original[1] = model{1, 0, true}, model{1, 1, true}
	for _, r := range rules {
		original = append(original, ruleModel(r, original))
	}
	rewrite(rules)
	if rules[1].op != opNot {
		t.Fatalf("NOT was not discovered: %+v", rules[1])
	}
	if rules[6].op != opXor {
		t.Fatalf("XOR was not discovered: %+v", rules[6])
	}
	revised := make([]model, 2, len(rules)+2)
	revised[0], revised[1] = model{1, 0, true}, model{1, 1, true}
	for i, r := range rules {
		revised = append(revised, ruleModel(r, revised))
		if revised[i+2] != original[i+2] {
			t.Fatalf("rule %d changed its expansion", i+2)
		}
	}
}

func TestGrammarChosenForRepetition(t *testing.T) {
	src := bytes.Repeat([]byte{0xaa, 0x55, 0xf0, 0x0f}, 1000)
	arc, s, err := compress(src)
	if err != nil {
		t.Fatal(err)
	}
	if s.raw || len(arc) >= len(src) {
		t.Fatalf("expected successful grammar compression: %+v", s)
	}
	got, err := decompress(arc, defaultLimit)
	if err != nil || !bytes.Equal(got, src) {
		t.Fatalf("bad grammar decode: %v", err)
	}
}

func TestCorruptionDetected(t *testing.T) {
	src := []byte("corrupt")
	arc, _, err := compress(src)
	if err != nil {
		t.Fatal(err)
	}
	arc[len(arc)-1] ^= 0x80
	if _, err = decompress(arc, defaultLimit); err == nil {
		t.Fatal("corruption went undetected")
	}
}
