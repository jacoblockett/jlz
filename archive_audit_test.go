package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestArchiveAudit(t *testing.T) {
	if os.Getenv("JLZ_AUDIT") != "1" {
		t.Skip("audit-only")
	}
	paths, err := filepath.Glob("fixtures/*")
	if err != nil {
		t.Fatal(err)
	}
	silesia := []string{"dickens", "mozilla", "mr", "nci", "ooffice", "osdb", "reymont", "samba", "sao", "webster", "x-ray", "xml"}
	for _, name := range silesia {
		paths = append(paths, filepath.Join("/mnt/data/jlz_round6/corpus", name))
	}
	sort.Strings(paths)
	for _, p := range paths {
		info, e := os.Stat(p)
		if e != nil || info.IsDir() {
			continue
		}
		input, e := os.ReadFile(p)
		if e != nil {
			t.Fatal(e)
		}
		arc, _, e := compress(input)
		if e != nil {
			t.Fatal(p, e)
		}
		decoded, e := decompress(arc, defaultLimit)
		if e != nil || !bytes.Equal(input, decoded) {
			t.Fatalf("roundtrip %s %v", p, e)
		}
		fmt.Printf("AUDIT,%s,%d,%d,%x\n", filepath.Base(p), len(input), len(arc), sha256.Sum256(arc))
	}
}
