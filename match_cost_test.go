package main

import (
	"bytes"
	"crypto/sha256"
	"os"
	"path/filepath"
	"testing"
)

// Optional corpus-level verification. With SILESIA set, exercise the revised
// parser and both unchanged token decoders on three distant regions per file.
func TestCostedMatchSilesiaRoundTrips(t *testing.T) {
	root := os.Getenv("SILESIA")
	if root == "" {
		t.Skip("SILESIA not set")
	}
	names := []string{"dickens", "mozilla", "mr", "nci", "ooffice", "osdb", "reymont", "samba", "sao", "webster", "x-ray", "xml"}
	checked := 0
	for _, name := range names {
		data, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}
		const n = 256 * 1024
		for _, at := range []int{0, (len(data) - n) / 2, len(data) - n} {
			region := data[at : at+n]
			tokens, err := tokenizeMatchV2(region)
			if err != nil {
				t.Fatal(err)
			}
			if len(tokens) == 0 {
				continue
			}
			hash := sha256.Sum256(region)
			for _, encode := range []func([]byte, [32]byte, []byte) ([]byte, error){encodeMatchV2Tokens, encodeMatchNativeTokens} {
				arc, err := encode(region, hash, tokens)
				if err != nil {
					t.Fatal(err)
				}
				got, err := decompress(arc, defaultLimit)
				if err != nil || !bytes.Equal(got, region) {
					t.Fatalf("%s at %d failed: %v", name, at, err)
				}
				checked++
			}
		}
	}
	if checked < 20 {
		t.Fatalf("unexpectedly few reconstructed match samples: %d", checked)
	}
	t.Logf("verified %d native and legacy match archives", checked)
}
