package bwt

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestBijectiveZeroRoundtrip(t *testing.T) {
	cases := [][]byte{nil, {0}, {255}, bytes.Repeat([]byte{42}, 100000), bytes.Repeat([]byte("banana!"), 30000), bytes.Repeat([]byte{0, 1, 2, 254, 255}, 5000)}
	for _, n := range []int{1, 12, 259, 4096, 65536} {
		v := make([]byte, n)
		rand.New(rand.NewSource(int64(n))).Read(v)
		cases = append(cases, v)
	}
	for _, v := range cases {
		out := BijectiveZeroEncode(v)
		got, e := BijectiveZeroDecode(out, len(v))
		if e != nil || !bytes.Equal(got, v) {
			t.Fatalf("len=%d out=%d error=%v", len(v), len(out), e)
		}
	}
	for _, v := range [][]byte{{255}, {255, 2}, {1, 1, 1, 1, 1, 1}, {2, 2}, {0, 0, 0, 0}} {
		if _, e := BijectiveZeroDecode(v, 1); e == nil {
			t.Fatalf("accepted corrupt %v", v)
		}
	}
}
