package engine

import (
	"bytes"
	"testing"

	"github.com/jacoblockett/jlz/bwt"
)

func TestBWTSubregionsExact(t *testing.T) {
	for _, plain := range [][]byte{bytes.Repeat([]byte("The long book was open on page one. "), 250), bytes.Repeat([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 1000), bytes.Repeat([]byte("abbaabbaabba"), 1200)} {
		packed, mode, e := encodeBWTSubregions(plain)
		if e != nil {
			t.Fatal(e)
		}
		var decoded []byte
		if mode == blockBWTSubregions {
			decoded, e = decodeBWTSubregions(packed, len(plain))
		} else {
			decoded, e = decodeBWTBijective(packed, len(plain))
		}
		if e != nil || !bytes.Equal(decoded, plain) {
			t.Fatalf("roundtrip mode=%d %v", mode, e)
		}
	}
}

// Explicitly exercise the on-wire decoder for both partition counts even
// when a deterministic encoder's size chooser would prefer the old format.
func TestBWTSubregionsWireAndCorruption(t *testing.T) {
	plain := bytes.Repeat([]byte("this is a repeated but varying sample: 0123456789 xyz. "), 300)
	transformed, primary, err := bwtNativeTransform(plain)
	if err != nil {
		t.Fatal(err)
	}
	syms := bwt.BijectiveZeroEncode(transformed)
	for _, parts := range []int{2, 4} {
		var out bytes.Buffer
		putVar(&out, uint64(primary))
		out.WriteByte(byte(parts))
		for j := 0; j < parts; j++ {
			a := len(syms) * j / parts
			z := len(syms) * (j + 1) / parts
			out.Write(packFiniteState(syms[a:z]))
		}
		full := out.Bytes()
		decoded, e := decodeBWTSubregions(full, len(plain))
		if e != nil || !bytes.Equal(plain, decoded) {
			t.Fatalf("mode parts=%d: %v", parts, e)
		}
		for _, cut := range []int{0, 1, 2, len(full) - 1} {
			_, e := decodeBWTSubregions(full[:cut], len(plain))
			if e == nil {
				t.Fatalf("accepted truncated mode=%d cut=%d", parts, cut)
			}
		}
		modified := append([]byte(nil), full...)
		offset := varSize(uint64(primary))
		modified[offset] = 3
		if _, e := decodeBWTSubregions(modified, len(plain)); e == nil {
			t.Fatal("accepted invalid count")
		}
	}
}
