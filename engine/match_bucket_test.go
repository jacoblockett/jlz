package engine

import "testing"

func TestCompactMatchBucketPreservesOrder(t *testing.T) {
	var b matchBucket
	var pos [2]uint32
	var next uint32
	for i := uint32(0); i < 2048; i++ {
		if b.latest(0) != pos[(next-1)&1] && next != 0 {
			t.Fatalf("latest mismatch at %d", i)
		}
		if next == 0 && (b.latest(0) != 0 || b.latest(1) != 0) {
			t.Fatal("empty bucket not empty")
		}
		if next >= 2 && b.latest(1) != pos[(next-2)&1] {
			t.Fatalf("older mismatch at %d", i)
		}
		val := i + 1
		pos[next&1] = val
		next++
		b = matchBucket((uint64(b) << 32) | uint64(val))
	}
	if b.latest(0) != 2048 || b.latest(1) != 2047 {
		t.Fatal("last two offsets incorrect")
	}
}
