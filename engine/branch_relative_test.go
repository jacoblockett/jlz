package engine

import (
	"bytes"
	"crypto/sha256"
	"testing"
)

func TestBranchRelativeBWTMalformedAndExact(t *testing.T) {
	src := make([]byte, stage3BlockSize)
	for i := 0; i+5 <= len(src); i+=19 {
		src[i] = 0xe8
		src[i+1] = byte(i)
		src[i+2] = byte(i>>8)
		src[i+3] = 11
		src[i+4] = 42
	}
	if !shouldBranchNormalize(src) { t.Fatal("fixture must be branch-dense") }
	encoded, subtype, err := encodeBWTSubregions(normalizeBranchRelative(src, false))
	if err != nil { t.Fatal(err) }
	payload := append([]byte{subtype}, encoded...)
	wrap := func(body []byte) []byte {
		var b bytes.Buffer
		b.WriteString(magic)
		b.WriteByte(modeMixedCross)
		putVar(&b, uint64(len(src)))
		h := sha256.Sum256(src)
		b.Write(h[:])
		b.WriteByte(blockBranchBWT)
		putVar(&b, uint64(len(src)))
		putVar(&b, uint64(len(body)))
		b.Write(body)
		return b.Bytes()
	}
	decoded, err := Decompress(wrap(payload), MaxInputSize)
	if err != nil || !bytes.Equal(decoded, src) { t.Fatalf("branch BWT roundtrip: %v", err) }
	for _, bad := range [][]byte{nil, {0}, append([]byte{blockCrossPacked}, encoded...), payload[:len(payload)-1]} {
		if out, e := Decompress(wrap(bad), MaxInputSize); e == nil { t.Fatalf("accepted malformed payload (input=%d output=%d)", len(bad), len(out)) }
	}
}
