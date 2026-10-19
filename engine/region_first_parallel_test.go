package engine

import (
	"bytes"
	"crypto/sha256"
	"runtime"
	"testing"
)

// Early regional probes must reach the same archive decision regardless of
// worker parallelism. Both paths compare samples in original source order.
func TestEarlyRegionalWorkerEquivalence(t *testing.T) {
	// At least 12 fixed regions are required to enter the early preflight.
	source := bytes.Repeat([]byte("a repeated long patterned phrase: one two three four; "), (13*stage3BlockSize)/51+1)
	source = source[:13*stage3BlockSize]
	digest := sha256.Sum256(source)
	rawSize := len(magic) + 1 + varSize(uint64(len(source))) + len(digest) + len(source)
	old := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(old)
	serial, serialOK, err := tryEarlyRegional(source, digest, rawSize)
	runtime.GOMAXPROCS(4)
	parallel, parallelOK, parallelErr := tryEarlyRegional(source, digest, rawSize)
	if err != nil || parallelErr != nil {
		t.Fatalf("regional probe errors: serial=%v parallel=%v", err, parallelErr)
	}
	if serialOK != parallelOK || !bytes.Equal(serial, parallel) {
		t.Fatalf("different preflight selection or archive: serial=%v (%d) parallel=%v (%d)", serialOK, len(serial), parallelOK, len(parallel))
	}
	if serialOK {
		restored, e := decompress(serial, defaultLimit)
		if e != nil || !bytes.Equal(restored, source) {
			t.Fatalf("roundtrip: %v", e)
		}
	}
}
