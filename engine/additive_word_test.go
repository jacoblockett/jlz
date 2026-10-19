package engine

import (
	"bytes"
	"crypto/sha256"
	"math/rand"
	"testing"
)

func TestAdditiveWordAllLengths(t *testing.T) {
	for _, n := range []int{0, 1, 2, 3, 7, 1024, 65535, stage3BlockSize, stage3BlockSize - 1} {
		src := make([]byte, n)
		rand.New(rand.NewSource(int64(n + 83))).Read(src)
		for mode := 0; mode < 3; mode++ {
			data := append([]byte(nil), src...)
			for i := 2; i+1 < len(data); i += 2 {
				if mode == 1 {
					data[i] = data[i-2]
					data[i+1] = data[i-1]
				}
				if mode == 2 {
					data[i] = data[i-2] + 4
					data[i+1] = data[i-1]
				}
			}
			packed := encodeAdditiveZigzagWords(data)
			restored, e := decodeAdditiveWords(packed, n)
			if e != nil || !bytes.Equal(data, restored) {
				t.Fatalf("n=%d mode=%d: %v", n, mode, e)
			}
			corrupt := append([]byte(nil), packed...)
			corrupt[0] = 5
			if _, e := decodeAdditiveWords(corrupt, n); e == nil {
				t.Fatal("accepted invalid mode")
			}
		}
	}
}
func TestAdditiveSilesiaArchiveCorruption(t *testing.T) {
	src := make([]byte, 4*stage3BlockSize)
	for i := 0; i+1 < len(src); i += 2 {
		src[i] = byte(i / 2)
		src[i+1] = byte(i / 2048)
	}
	arc, e := encodeParallelAdditiveArchive(src, sha256.Sum256(src))
	if e != nil {
		t.Fatal(e)
	}
	if decoded, e := decompress(arc, defaultLimit); e != nil || !bytes.Equal(src, decoded) {
		t.Fatal(e)
	}
	damaged := append([]byte(nil), arc...)
	damaged[len(damaged)-1] ^= 0x20
	if _, e := decompress(damaged, defaultLimit); e == nil {
		t.Fatal("corrupt archive accepted")
	}
}
