package bwt

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestTwoHopInverseRoundtrip(t *testing.T) {
	cases := [][]byte{{0}, {255}, {0, 0}, {255, 0}, {0, 255}, []byte("banana"), bytes.Repeat([]byte{0}, 4096), bytes.Repeat([]byte{255, 0, 128, 0, 255}, 3500)}
	for _, n := range []int{3, 8, 31, 256, 4096, 65536, 256 * 1024} {
		b := make([]byte, n)
		rand.New(rand.NewSource(int64(n))).Read(b)
		cases = append(cases, b)
	}
	for _, v := range cases {
		coded, primary, e := Transform(v)
		if e != nil {
			t.Fatal(e)
		}
		decoded, e := Inverse(coded, primary)
		if e != nil || !bytes.Equal(decoded, v) {
			t.Fatalf("two-hop inverse size=%d primary=%d: %v", len(v), primary, e)
		}
	}
	for _, tc := range []struct {
		b []byte
		p int
	}{{nil, 0}, {[]byte{1}, 0}, {[]byte{0}, 2}} {
		if _, e := Inverse(tc.b, tc.p); e == nil {
			t.Errorf("accepted invalid BWT primary %d", tc.p)
		}
	}
}
