package jlz

import (
	"bytes"
	"crypto/sha256"
	"testing"
)

// Reference the original, bit-at-a-time Huffman body implementation.
func referencePackStream(data []byte) []byte {
	best := append([]byte(nil), data...)
	codec := byte(0)
	if len(data) > 2 {
		all := true
		for _, v := range data[1:] {
			if v != data[0] {
				all = false
				break
			}
		}
		if all {
			best = []byte{data[0]}
			codec = 1
		}
	}
	if len(data) > 16 && codec == 0 {
		lengths, ok := nativeLengths(data)
		if ok {
			codes, _, _, _, _, err := nativeCodes(lengths)
			if err == nil {
				cnt, bits := 0, 0
				for _, l := range lengths {
					if l > 0 {
						cnt++
					}
				}
				for _, v := range data {
					bits += int(lengths[v])
				}
				header := varSize(uint64(cnt)) + 2*cnt
				if header+(bits+7)/8 < len(best) {
					var buf bytes.Buffer
					putVar(&buf, uint64(cnt))
					for v, l := range lengths {
						if l > 0 {
							buf.WriteByte(byte(v))
							buf.WriteByte(l)
						}
					}
					payload := make([]byte, (bits+7)/8)
					at := 0
					for _, v := range data {
						l := int(lengths[v])
						c := codes[v]
						for k := l - 1; k >= 0; k-- {
							if (c>>k)&1 != 0 {
								payload[at/8] |= 1 << uint(7-at%8)
							}
							at++
						}
					}
					buf.Write(payload)
					if buf.Len() < len(best) {
						best = buf.Bytes()
						codec = 2
					}
				}
			}
		}
	}
	var result bytes.Buffer
	putVar(&result, uint64(len(data)))
	putVar(&result, uint64(len(best)))
	result.WriteByte(codec)
	result.Write(best)
	return result.Bytes()
}

func TestFastPackingByteExactWithBitwiseReference(t *testing.T) {
	var seed uint64 = 0x192e123caa445566
	next := func() byte { seed ^= seed << 13; seed ^= seed >> 7; seed ^= seed << 17; return byte(seed >> 17) }
	for _, n := range []int{0, 1, 2, 17, 31, 127, 1024, 8192, 65536} {
		for typ := 0; typ < 5; typ++ {
			d := make([]byte, n)
			for i := range d {
				switch typ {
				case 0:
					d[i] = next()
				case 1:
					d[i] = byte(i % 3)
				case 2:
					d[i] = byte(i % 23)
				case 3:
					d[i] = 0
				case 4:
					d[i] = byte(i % 113)
				}
			}
			got, want := packStream(d), referencePackStream(d)
			if !bytes.Equal(got, want) {
				t.Fatalf("n=%d type=%d", n, typ)
			}
		}
	}
}
func TestSharedPackedCandidateByteExact(t *testing.T) {
	for _, input := range [][]byte{
		bytes.Repeat([]byte("Lorem ipsum dolor sit amet. "), 6000),
		bytes.Repeat([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 16000),
	} {
		tokens, err := tokenizeMatchV2(input)
		if err != nil {
			t.Fatal(err)
		}
		if len(tokens) == 0 {
			continue
		}
		d := sha256.Sum256(input)
		a, err := encodeMatchPackedVariant(input, d, tokens, false)
		if err != nil {
			t.Fatal(err)
		}
		b, err := encodeMatchPackedVariant(input, d, tokens, true)
		if err != nil {
			t.Fatal(err)
		}
		a, e := innerPayload(a, modeMatchPacked, len(input))
		if e != nil {
			t.Fatal(e)
		}
		b, e = innerPayload(b, modeMatchPackedPred, len(input))
		if e != nil {
			t.Fatal(e)
		}
		m, want := blockMatchPacked, a
		if len(b) < len(a) {
			m, want = blockMatchPackedPred, b
		}
		gotmode, got, err := packedRegionVariantsShared(input, tokens)
		if err != nil {
			t.Fatal(err)
		}
		if gotmode != m || !bytes.Equal(got, want) {
			t.Fatalf("packed variants differ")
		}
	}
}
