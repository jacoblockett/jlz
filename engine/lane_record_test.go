package engine

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestRecordAlignedLaneRoundTrip(t *testing.T) {
	r := rand.New(rand.NewSource(57))
	for _, n := range []int{1, 27, 28, 29, 55, 56, 4096, 131073, stage3BlockSize} {
		b := make([]byte, n)
		for i := range b {
			switch i % 28 {
			case 0, 1, 2, 7, 15, 20:
				b[i] = byte(i % 28)
			default:
				b[i] = byte(r.Intn(256))
			}
		}
		encoded := encodeLaneBlock(b, 28)
		decoded, err := decodeLaneBlock(encoded, n)
		if err != nil || !bytes.Equal(b, decoded) {
			t.Fatalf("record inverse %d: %v", n, err)
		}
		if len(encoded) > 1 {
			if _, err = decodeLaneBlock(encoded[:len(encoded)-1], n); err == nil {
				t.Fatalf("truncated record accepted n=%d", n)
			}
		}
	}
	encoded := encodeLaneBlock(bytes.Repeat([]byte("abcd"), 80), 28)
	for _, invalid := range []byte{0, 1, 3, 9, 11, 65, 255} {
		damaged := append([]byte(nil), encoded...)
		damaged[0] = invalid
		if _, err := decodeLaneBlock(damaged, 320); err == nil {
			t.Errorf("accepted invalid record stride %d", invalid)
		}
	}
}

func TestRecordPeriodDiscovery(t *testing.T) {
	var samples [3][]byte
	for k := 0; k < 3; k++ {
		b := make([]byte, 4096)
		r := rand.New(rand.NewSource(int64(91 + k)))
		for i := range b {
			switch i % 28 {
			case 0, 4, 12, 19:
				b[i] = byte(i % 28)
			default:
				b[i] = byte(r.Intn(256))
			}
		}
		samples[k] = b
	}
	if got := recordAlignedStride(samples); got != 28 {
		t.Fatalf("expected 28 byte period; got %d", got)
	}
	for k := range samples {
		rand.New(rand.NewSource(int64(101 + k))).Read(samples[k])
	}
	if got := recordAlignedStride(samples); got != 0 {
		t.Fatalf("random data false positive: %d", got)
	}
	samples[0] = []byte{0}
	if got := recordAlignedStride(samples); got != 0 {
		t.Fatalf("short + noise false positive: %d", got)
	}
}
