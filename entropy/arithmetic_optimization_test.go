package entropy

import "testing"

// Reference division widths and bit indexing match the historical JLZ1 coder.
func TestBoundedProbabilityExact(t *testing.T) {
	model := new(bitModel)
	for i := 0; i < 100000; i++ {
		prev := byte(i * 73 >> 3)
		prefix := uint16(1 + (i*37)%255)
		z0, o0 := normalizedCount(model.base[prefix])
		z1, o1 := normalizedCount(model.prev[int(prev)*256+int(prefix)])
		p0 := (uint64(z0) * 65536) / (uint64(z0) + uint64(o0))
		p1 := (uint64(z1) * 65536) / (uint64(z1) + uint64(o1))
		weight := uint64(1)
		if z1+o1 >= 10 {
			weight = 5
		}
		want := (p0*(8-weight) + p1*weight) / 8
		if want < 1 {
			want = 1
		}
		if want > 65535 {
			want = 65535
		}
		if got := model.probability(prev, prefix); got != want {
			t.Fatalf("i=%d got=%d want=%d", i, got, want)
		}
		model.update(prev, prefix, byte((i>>3)&1))
	}
}

func TestBufferedBitReaderExact(t *testing.T) {
	data := []byte{0, 0x96, 0xff, 0x81, 0x55}
	r := bitReader{data: data}
	for p := 0; p < len(data)*8+65; p++ {
		var want uint64
		if p/8 < len(data) {
			want = uint64((data[p/8] >> (7 - uint(p%8))) & 1)
		}
		if got := r.bit(); got != want {
			t.Fatalf("bit %d got=%d want=%d", p, got, want)
		}
	}
}
