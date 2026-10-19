package engine

import (
	"bytes"
	"crypto/sha256"
	"math/rand"
	"testing"
)

// Demonstrate that default compression discovers, selects, and decodes a real
// winning mixed archive due to the previously disconnected legacy match mode.
// This is a structural control, NOT a real-world corpus result.
func TestDefaultSelectsMatchDrivenMixedArchive(t *testing.T) {
	r := rand.New(rand.NewSource(900210))
	var input []byte
	for b := 0; b < 4; b++ {
		seed := make([]byte, 64*1024)
		if _, err := r.Read(seed); err != nil {
			t.Fatal(err)
		}
		block := append([]byte(nil), seed...)
		for variant := 0; variant < 3; variant++ {
			transformed := append([]byte(nil), seed...)
			for i := range transformed {
				if variant%2 == 0 {
					transformed[i] ^= byte(19 + b + variant)
				} else {
					transformed[i] += byte(101 + b)
				}
			}
			block = append(block, transformed...)
		}
		input = append(input, block...)
	}
	if !stage3StructuralSampleUseful(input) {
		t.Fatal("structural gate did not detect transformed repeats")
	}
	result, s, err := compress(input)
	if err != nil {
		t.Fatal(err)
	}
	if result[4] != modeMixed {
		t.Fatalf("match not connected to default selection: method=%s, candidate=%d", s.method, len(result))
	}
	restored, err := decompress(result, defaultLimit)
	if err != nil || !bytes.Equal(input, restored) {
		t.Fatalf("mixed roundtrip: %v", err)
	}
	// The mixed archive must actually CONTAIN mode 5, not merely happen to
	// select mixed prediction on an unrelated block.
	rdr := bytes.NewReader(result[5:])
	_, err = readVar(rdr)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = rdr.Seek(32, 1); err != nil {
		t.Fatal(err)
	}
	legacy := 0
	for rdr.Len() > 0 {
		m, err := rdr.ReadByte()
		if err != nil {
			t.Fatal(err)
		}
		_, err = readVar(rdr)
		if err != nil {
			t.Fatal(err)
		}
		plen, err := readVar(rdr)
		if err != nil {
			t.Fatal(err)
		}
		if plen > uint64(rdr.Len()) {
			t.Fatal("bad length")
		}
		if m == blockLegacyMatch || m == blockNative || m == blockMatchedLiteral || m == blockMatchPacked || m == blockMatchPackedPred {
			legacy++
		}
		if _, err = rdr.Seek(int64(plen), 1); err != nil {
			t.Fatal(err)
		}
	}
	if legacy == 0 {
		t.Fatalf("mixed archive lacked native mathematical matches, archive size %d", len(result))
	}
	t.Logf("mixed=%d native match blocks=%d", len(result), legacy)
}

// A legacy block has no independent digest: integrity comes from the outer
// archive's SHA-256. Confirm truncated, altered and structurally malformed
// combinations are rejected by the real decoder.
func TestMixedLegacyCorruptionAndDeterminism(t *testing.T) {
	rng := rand.New(rand.NewSource(909))
	seed := make([]byte, 64*1024)
	if _, err := rng.Read(seed); err != nil {
		t.Fatal(err)
	}
	data := append([]byte(nil), seed...)
	for k := 0; k < 3; k++ {
		part := append([]byte(nil), seed...)
		for i := range part {
			part[i] ^= byte(31 + k)
		}
		data = append(data, part...)
	}
	digest := sha256.Sum256(data)
	archive, err := encodeMixedBlocks(data, digest, len(data)*2)
	if err != nil {
		t.Fatal(err)
	}
	if len(archive) == 0 {
		t.Fatal("no mixed archive")
	}
	cr := bytes.NewReader(archive[5:])
	if _, err = readVar(cr); err != nil {
		t.Fatal(err)
	}
	if _, err = cr.Seek(32, 1); err != nil {
		t.Fatal(err)
	}
	modeIndex := len(archive) - cr.Len()
	if archive[modeIndex] != blockLegacyMatch && archive[modeIndex] != blockNative && archive[modeIndex] != blockMatchedLiteral && archive[modeIndex] != blockMatchPacked && archive[modeIndex] != blockMatchPackedPred {
		t.Fatalf("expected a native match block, mode %d", archive[modeIndex])
	}
	for _, mutant := range [][]byte{
		archive[:len(archive)-1],
		func() []byte { c := append([]byte(nil), archive...); c[len(c)-2] ^= 0x17; return c }(),
		func() []byte { c := append([]byte(nil), archive...); c[modeIndex] = 0xff; return c }(),
		func() []byte { c := append([]byte(nil), archive...); c[6] ^= 0x01; return c }(),
	} {
		if _, err = decompress(mutant, defaultLimit); err == nil {
			t.Fatal("corruption accepted")
		}
	}
	second, err := encodeMixedBlocks(data, digest, len(data)*2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(second, archive) {
		t.Fatal("mixed encoder is not deterministic")
	}
	got, err := decompress(second, defaultLimit)
	if err != nil || !bytes.Equal(got, data) {
		t.Fatalf("mixed roundtrip: %v", err)
	}
}
