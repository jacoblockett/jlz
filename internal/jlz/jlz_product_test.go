package jlz

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestJLZFormatAndFixedEncoder(t *testing.T) {
	samples := [][]byte{nil, []byte("JLZ"), bytes.Repeat([]byte("red blue green blue red\n"), 3000)}
	for _, original := range samples {
		encoded, _, err := compress(original)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(encoded[:4], []byte("JLZ1")) {
			t.Fatalf("wrong JLZ header %q", encoded[:4])
		}
		duplicate, _, err := compress(original)
		if err != nil || !bytes.Equal(encoded, duplicate) {
			t.Fatal("encoder not deterministic")
		}
		decoded, err := decompress(encoded, defaultLimit)
		if err != nil || !bytes.Equal(original, decoded) {
			t.Fatalf("invalid reconstruction: %v", err)
		}
		old := bytes.Clone(encoded)
		copy(old[:4], "MDN1")
		if _, err = decompress(old, defaultLimit); err == nil {
			t.Fatal("accepted older MDN1 format")
		}
	}
}

func TestJLZBundledFixtureRoundTrips(t *testing.T) {
	paths, err := filepath.Glob(filepath.Join("testdata", "fixtures", "*.bin"))
	if err != nil {
		t.Fatal(err)
	}
	paths = append(paths, filepath.Join("testdata", "fixtures", "21_services.txt"), filepath.Join("testdata", "fixtures", "22_records.json"), filepath.Join("testdata", "fixtures", "23_uncompressed_audio.wav"))
	if len(paths) < 10 {
		t.Fatal("fixtures missing")
	}
	for _, path := range paths {
		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		arc, _, err := compress(src)
		if err != nil {
			t.Fatalf("%s encode: %v", path, err)
		}
		restored, err := decompress(arc, defaultLimit)
		if err != nil || !bytes.Equal(restored, src) {
			t.Fatalf("%s decode: %v", path, err)
		}
	}
}
