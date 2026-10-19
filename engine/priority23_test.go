package engine

import (
	"bytes"
	"crypto/sha256"
	"math/rand"
	"testing"
)

func TestPackedStreamsRoundtripAndCorruption(t *testing.T) {
	rng := rand.New(rand.NewSource(123))
	samples := [][]byte{nil, []byte{1}, bytes.Repeat([]byte{42}, 4096), bytes.Repeat([]byte("abracadabra"), 1024)}
	random := make([]byte, 65536)
	rng.Read(random)
	samples = append(samples, random)
	for _, data := range samples {
		packed := packStream(data)
		r := bytes.NewReader(packed)
		restored, err := unpackStream(r, 65536)
		if err != nil || r.Len() != 0 || !bytes.Equal(restored, data) {
			t.Fatalf("packed roundtrip %d: %v", len(data), err)
		}
		for _, candidate := range [][]byte{packed[:len(packed)-1], append([]byte(nil), packed...)} {
			if len(candidate) == len(packed) {
				// Find the codec byte after two variable-length fields.
				index := bytes.NewReader(candidate)
				_, _ = readVar(index)
				_, _ = readVar(index)
				candidate[len(candidate)-index.Len()] = 0xff
			}
			_, err = unpackStream(bytes.NewReader(candidate), 65536)
			if err == nil {
				t.Fatal("accepted corrupt packed stream")
			}
		}
	}
}

func TestPackedMatchEncoderVariantsAndIntegrity(t *testing.T) {
	rng := rand.New(rand.NewSource(44))
	seed := make([]byte, 8192)
	rng.Read(seed)
	input := append([]byte(nil), seed...)
	for i := 0; i < 4; i++ {
		extra := append([]byte(nil), seed...)
		for j := range extra {
			if i%2 == 0 {
				extra[j] ^= 0x39
			} else {
				extra[j] += 0x51
			}
		}
		input = append(input, extra...)
	}
	tokens, err := tokenizeMatchV2(input)
	if err != nil || len(tokens) == 0 {
		t.Fatalf("matcher: %v", err)
	}
	for _, predict := range []bool{false, true} {
		archive, err := encodeMatchPackedVariant(input, sha256.Sum256(input), tokens, predict)
		if err != nil {
			t.Fatal(err)
		}
		decoded, err := decompress(archive, defaultLimit)
		if err != nil || !bytes.Equal(decoded, input) {
			t.Fatalf("packed variant %v: %v", predict, err)
		}
		corrupted := append([]byte(nil), archive...)
		corrupted[len(corrupted)-1] ^= 0x1
		if _, err := decompress(corrupted, defaultLimit); err == nil {
			t.Fatal("corrupt packed archive accepted")
		}
	}
}

func TestMixedBlockSharedDictionaryAndModes(t *testing.T) {
	rng := rand.New(rand.NewSource(98))
	first := make([]byte, stage3BlockSize)
	rng.Read(first)
	second := make([]byte, stage3BlockSize)
	rng.Read(second)
	copy(second[:16384], first[len(first)-16384:])
	input := append(first, second...)
	archive, err := encodeMixedBlocks(input, sha256.Sum256(input), len(input)*2)
	if err != nil {
		t.Fatal(err)
	}
	recovered, err := decompress(archive, defaultLimit)
	if err != nil || !bytes.Equal(recovered, input) {
		t.Fatalf("dictionary decode %v", err)
	}
	// Native-only archives must never contain the old external preset dictionary mode.
	r := bytes.NewReader(archive[5:])
	_, err = readVar(r)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = r.Seek(32, 1); err != nil {
		t.Fatal(err)
	}
	for r.Len() > 0 {
		mode, e := r.ReadByte()
		if e != nil {
			t.Fatal(e)
		}
		if mode == blockOldDictionary {
			t.Fatal("external DEFLATE dictionary survived native-only rewrite")
		}
		_, e = readVar(r)
		if e != nil {
			t.Fatal(e)
		}
		size, e := readVar(r)
		if e != nil {
			t.Fatal(e)
		}
		if size > uint64(r.Len()) {
			t.Fatal("malformed native block")
		}
		if _, e = r.Seek(int64(size), 1); e != nil {
			t.Fatal(e)
		}
	}
}
