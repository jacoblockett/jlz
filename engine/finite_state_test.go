package engine

import (
	"bytes"
	"testing"
)

// Validate model serialization, the full native packed stream framing,
// byte-renormalized state restoration, and truncated/corrupt input handling.
func TestFiniteStateExact(t *testing.T) {
	for _, n := range []int{1, 2, 3, 9, 17, 255, 4096, 30000} {
		for mode := 0; mode < 5; mode++ {
			b := make([]byte, n)
			for i := range b {
				switch mode {
				case 0:
					b[i] = 0
				case 1:
					b[i] = byte(i & 1)
				case 2:
					b[i] = byte((i*i + 73*i) >> 4)
				case 3:
					b[i] = byte(i)
				case 4:
					b[i] = byte((i * 7919) ^ (i >> 2))
				}
			}
			m, e := finiteStateModelFor(b)
			if e != nil {
				t.Fatal(e)
			}
			wire := finiteStateTable(&m)
			restoredModel, at, e := finiteStateReadTable(wire)
			if e != nil || at != len(wire) {
				t.Fatalf("table n=%d mode=%d: %v", n, mode, e)
			}
			coded, e := finiteStateEncode(b, &m)
			if e != nil {
				t.Fatal(e)
			}
			out, e := finiteStateDecode(coded, &restoredModel, len(b))
			if e != nil || !bytes.Equal(out, b) {
				t.Fatalf("coder n=%d mode=%d: %v", n, mode, e)
			}
			packed := packFiniteState(b)
			out, e = unpackStream(bytes.NewReader(packed), len(b))
			if e != nil || !bytes.Equal(out, b) {
				t.Fatalf("packed n=%d mode=%d: %v", n, mode, e)
			}
			if len(packed) > len(packStream(b)) {
				t.Fatalf("size n=%d mode=%d", n, mode)
			}
			if len(coded) > 4 {
				if _, e = finiteStateDecode(coded[:len(coded)-1], &m, len(b)); e == nil {
					t.Fatalf("truncation n=%d mode=%d", n, mode)
				}
			}
		}
	}
}

func TestFiniteStateRejectCorruptTable(t *testing.T) {
	for _, table := range [][]byte{{}, {0}, {1, 2}, {1, 2, 0}, {2, 3, 255, 3, 255}, {2, 0, 128, 32, 0, 128, 32}} {
		if _, _, e := finiteStateReadTable(table); e == nil {
			t.Fatalf("accepted invalid table %v", table)
		}
	}
	m, e := finiteStateModelFor([]byte("mississippi"))
	if e != nil {
		t.Fatal(e)
	}
	if _, e = finiteStateDecode([]byte{0, 0, 0, 0}, &m, 1); e == nil {
		t.Fatal("accepted invalid state")
	}
}
