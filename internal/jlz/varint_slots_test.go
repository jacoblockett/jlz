package jlz

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestVarintSlotsExact(t *testing.T) {
	r := rand.New(rand.NewSource(78281))
	for _, count := range []int{0, 1, 20, 100, 200, 900, 10000} {
		for typ := 0; typ < 7; typ++ {
			var input []byte
			for i := 0; i < count; i++ {
				var v uint64
				switch typ {
				case 0:
					v = uint64(r.Intn(3) + 1)
				case 1:
					v = uint64(r.Intn(400) + 1)
				case 2:
					v = uint64(r.Intn(150000) + 1)
				case 3:
					v = uint64(r.Uint32()) + 1
				case 4:
					v = uint64((i % 17) + 1)
				case 5:
					v = uint64((i % 4096) + 1)
				default:
					v = uint64(r.Intn(20000000) + 1)
				}
				input = appendUvar(input, v)
			}
			enc := packVarintStream(input)
			dec, e := nativeUnpackBytes(enc, len(input)+20)
			if e != nil || !bytes.Equal(input, dec) {
				t.Fatalf("count %d type %d: %v", count, typ, e)
			}
		}
	}
}
func TestVarintSlotsCorruption(t *testing.T) {
	var input []byte
	for i := 0; i < 2000; i++ {
		input = appendUvar(input, uint64(1+(i*191)%200000))
	}
	packed := packVarintStream(input)
	r := bytes.NewReader(packed)
	_, _ = readVar(r)
	_, _ = readVar(r)
	codec, _ := r.ReadByte()
	if codec != 4 {
		t.Fatal("test data did not select slots")
	}
	for _, data := range [][]byte{packed[:len(packed)-1], append(append([]byte{}, packed...), 1)} {
		_, e := nativeUnpackBytes(data, len(input)+20)
		if e == nil {
			t.Fatal("accepted damaged stream")
		}
	}
}
