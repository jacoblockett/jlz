package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestGateStudy(t *testing.T) {
	root := os.Getenv("SILESIA")
	if root == "" {
		t.Skip("SILESIA missing")
	}
	for _, name := range []string{"dickens", "mozilla", "mr", "nci", "ooffice", "osdb", "reymont", "samba", "sao", "webster", "x-ray", "xml"} {
		data, e := os.ReadFile(filepath.Join(root, name))
		if e != nil {
			t.Fatal(e)
		}
		sample := data[:min(len(data), 24*1024)]
		digest := sha256.Sum256(sample)
		a, e := encodeMatchV2(sample, digest)
		if e != nil {
			t.Fatal(e)
		}
		b, e := encodeFastNative(sample, digest)
		if e != nil {
			t.Fatal(e)
		}
		eligible, e := matchSampleUseful(data)
		if e != nil {
			t.Fatal(e)
		}
		fmt.Printf("GATE %s match %d baseline %d ratio %.3f pass %v\n", name, len(a), len(b), float64(len(a))/float64(len(b)), eligible)
	}
}
