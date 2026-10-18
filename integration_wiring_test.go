package main

import (
	"bytes"
	"crypto/sha256"
	"os"
	"path/filepath"
	"testing"
)

func TestMixedBlockLegacyCanActuallyWin(t *testing.T) {
	root := os.Getenv("SILESIA")
	if root == "" {
		t.Skip("SILESIA not set")
	}
	for _, name := range []string{"mozilla", "samba", "mr"} {
		data, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}
		block := data[:stage3BlockSize]
		mode, payload, err := encodeStage3Block(block)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%s first block: selected=%d payload=%d matchPrefilter=%t", name, mode, len(payload), stage3MatchPotential(block))
		// Compose a complete valid mixed frame, proving the final dispatcher can decode mode 5.
		archive, err := encodeMixedBlocks(block, sha256.Sum256(block), len(block)*2)
		if err != nil || len(archive) == 0 {
			t.Fatalf("archive: %v", err)
		}
		got, err := decompress(archive, defaultLimit)
		if err != nil || !bytes.Equal(got, block) {
			t.Fatalf("%s: %v", name, err)
		}
		if name == "mozilla" && mode != blockLegacyMatch {
			t.Fatalf("expected legacy block win, got %d", mode)
		}
	}
}
