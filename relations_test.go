package main

import (
	"bytes"
	"crypto/sha256"
	"math/rand"
	"testing"
)

func randomBlock(seed int64, n int) []byte {
	b := make([]byte, n)
	rand.New(rand.NewSource(seed)).Read(b)
	return b
}
func xor(a, b []byte) []byte {
	out := make([]byte, len(a))
	for i := range a {
		out[i] = a[i] ^ b[i]
	}
	return out
}
func testGraphCase(t *testing.T, blocks [][]byte, expected byte) {
	t.Helper()
	src := bytes.Join(blocks, nil)
	ops, _ := selectGraphOps(src, len(blocks[0]))
	present := false
	for _, op := range ops {
		if op.kind == expected {
			present = true
		}
	}
	if !present {
		t.Fatalf("didn't find op %d: %+v", expected, ops)
	}
	arc, _, err := encodeGraph(src, sha256.Sum256(src), len(blocks[0]))
	if err != nil {
		t.Fatal(err)
	}
	restored, err := decompress(arc, defaultLimit)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(src, restored) {
		t.Fatal("graph mismatch")
	}
}
func TestGraphExact(t *testing.T) {
	a := randomBlock(123, 1024)
	b := randomBlock(456, 1024)
	testGraphCase(t, [][]byte{a, b, a, b, a}, blockCopy)
}
func TestGraphXorConstant(t *testing.T) {
	a := randomBlock(123, 1024)
	b := make([]byte, len(a))
	for i := range a {
		b[i] = a[i] ^ 0xA3
	}
	testGraphCase(t, [][]byte{a, b}, blockXorConstant)
}
func TestGraphXorPeriodic(t *testing.T) {
	a := randomBlock(123, 1024)
	b := make([]byte, len(a))
	for i := range a {
		b[i] = a[i] ^ []byte{0x11, 0x25, 0x68}[i%3]
	}
	testGraphCase(t, [][]byte{a, b}, blockXorPattern)
}
func TestGraphXorPair(t *testing.T) {
	a := randomBlock(123, 1024)
	b := randomBlock(456, 1024)
	testGraphCase(t, [][]byte{a, b, xor(a, b)}, blockXorPair)
}
func TestGraphXorResidual(t *testing.T) {
	a := randomBlock(123, 1024)
	b := bytes.Clone(a)
	for i := 1; i < len(b); i += 79 {
		b[i] ^= 0xF3
	}
	testGraphCase(t, [][]byte{a, b}, blockXorResidual)
}
func TestLagBothOperations(t *testing.T) {
	rng := rand.New(rand.NewSource(76854))
	for _, op := range []byte{lagXOR, lagSubtract} {
		src := make([]byte, 10243)
		rng.Read(src[:127])
		for i := 127; i < len(src); i++ {
			if op == lagXOR {
				src[i] = src[i-127] ^ 0xA4
			} else {
				src[i] = src[i-127] + byte(i%3)
			}
		}
		res, _ := lagResidual(src, 127, op)
		arc, err := encodeLag(src, sha256.Sum256(src), 127, op, res)
		if err != nil {
			t.Fatal(err)
		}
		restored, err := decompress(arc, defaultLimit)
		if err != nil {
			t.Fatalf("op=%d: %v", op, err)
		}
		if !bytes.Equal(src, restored) {
			t.Fatalf("wrong lag op %d", op)
		}
	}
}
func TestNativeAndLegacyRejection(t *testing.T) {
	src := bytes.Repeat([]byte("the quick brown fox"), 100)
	arc, err := encodeNativeBytes(src, sha256.Sum256(src))
	if err != nil {
		t.Fatal(err)
	}
	out, err := decompress(arc, defaultLimit)
	if err != nil || !bytes.Equal(src, out) {
		t.Fatalf("native mode: %v", err)
	}
	// Legacy raw archive from v1 must continue to decompress correctly.
	var raw bytes.Buffer
	raw.WriteString("MDC1")
	raw.WriteByte(modeRaw)
	putVar(&raw, uint64(len(src)))
	sum := sha256.Sum256(src)
	raw.Write(sum[:])
	raw.Write(src)
	_, err = decompress(raw.Bytes(), defaultLimit)
	if err == nil {
		t.Fatal("legacy container silently accepted")
	}
}
func TestRelationsRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(20918))
	for i := 0; i < 30; i++ {
		n := 256 + rng.Intn(7000)
		src := make([]byte, n)
		rng.Read(src)
		for _, w := range []int{128, 256, 384, 1024} {
			arc, _, err := encodeGraph(src, sha256.Sum256(src), w)
			if err != nil {
				t.Fatal(err)
			}
			if len(arc) == 0 {
				continue
			}
			out, err := decompress(arc, defaultLimit)
			if err != nil || !bytes.Equal(src, out) {
				t.Fatalf("graph fuzz n=%d width=%d err=%v", n, w, err)
			}
		}
		for _, op := range []byte{lagXOR, lagSubtract} {
			lag := 1 + rng.Intn(n/2)
			r, _ := lagResidual(src, lag, op)
			arc, err := encodeLag(src, sha256.Sum256(src), lag, op, r)
			if err != nil {
				t.Fatal(err)
			}
			out, err := decompress(arc, defaultLimit)
			if err != nil || !bytes.Equal(src, out) {
				t.Fatalf("lag fuzz n=%d lag=%d op=%d err=%v", n, lag, op, err)
			}
		}
	}
}
func TestHashDetectsMutations(t *testing.T) {
	src := bytes.Repeat([]byte{7, 3, 12, 8}, 300)
	arc, _, err := compress(src)
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range []int{5, len(arc) / 2, len(arc) - 1} {
		damaged := bytes.Clone(arc)
		damaged[i] ^= 128
		out, err := decompress(damaged, defaultLimit)
		if err == nil && !bytes.Equal(out, src) {
			t.Fatal("corrupt data accepted")
		}
	}
}
