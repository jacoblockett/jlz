package jlz

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestExistingArchiveCompatibility(t *testing.T) {
	root := os.Getenv("JLZ_CORPUS")
	old := os.Getenv("JLZ_REFERENCE_ARCHIVES")
	if root == "" || old == "" {
		t.Skip("reference corpus required")
	}
	names := []string{"dickens", "mozilla", "mr", "nci", "ooffice", "osdb", "reymont", "samba", "sao", "webster", "x-ray", "xml"}
	for _, name := range names {
		input, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}
		reference, err := os.ReadFile(filepath.Join(old, name+".jlz"))
		if err != nil {
			t.Fatal(err)
		}
		decoded, err := decompress(reference, 1<<30)
		if err != nil || !bytes.Equal(decoded, input) {
			t.Fatalf("old archive decode %s: %v", name, err)
		}
		encoded, _, err := compressFast(input)
		if err != nil {
			t.Fatal(err)
		}
		decoded, err = decompress(encoded, 1<<30)
		if err != nil || !bytes.Equal(decoded, input) {
			t.Fatalf("new archive decode %s: %v", name, err)
		}
		if !bytes.Equal(encoded, reference) {
			t.Fatalf("archive changed with adaptive concurrency: %s", name)
		}
	}
}
