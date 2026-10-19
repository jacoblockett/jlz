# Go feature-package reorganization, 2026-10-19

**Scope:** repository layout and package/API boundaries, not compression research or a new codec revision. Current performance baseline stays Round 29.

## What changed

- Restored source entry point `main.go` at repository root. It remains source only, not a tracked executable.
- Removed redundant `cmd/jlz/` and `internal/jlz/` paths.
- Moved shared archive selection, dispatch, cross-region coordination, historical decoder formats, and white-box integration tests into `engine/`.
- Extracted the reversible modular word transform into `additive/words.go`, the BWT/MTF/zero-run transforms into `bwt/transform.go`, arithmetic coder into `entropy/arithmetic.go`, independently framed lane coding into `lane/block.go`, matching fingerprint/probe/copy utilities into `match/`, lag transforms into `predictive/transform.go`, portable worker policy into `workers/`, and varint framing into `archive/varint.go`.
- Changed duplicated engine implementations into thin internal forwarding functions to preserve established wire-format entry points, internal tests, and codec selection. The leaf packages import no engine dependencies, avoiding import cycles.
- Moved all 15 byte-exact fixtures to `engine/testdata/fixtures/` and updated documentation/build instructions. Root `go.mod` retains `github.com/jacoblockett/jlz`. Compiled binaries remain excluded by `.gitignore`.

## Validation evidence and limitation

Static Git-tree checks confirm root `main.go`, 15 intact fixture blobs, no `cmd/jlz` or `internal/jlz` paths, no tracked built binaries, and consistent dependency direction from `engine` to feature packages. Reviewed the moved/modified Go sources for missing imports and cross-package identifiers; corrected use of BWT varints, transformed-match hashing constants, and predictive transform constants.

**Go build, gofmt, go vet, unit/race tests, and actual original-byte Silesia roundtrips have NOT been executed for this refactor.** The local runtime cannot reach github.com (DNS and direct connection refused) and GitHub connector source is not mounted into it; hosted Actions are explicitly forbidden. This is a critical outstanding validation step. Do not claim that archive-size/time equivalence has been experimentally verified. Use local source checkout and run `gofmt -w`, `go test -count=1 ./...`, `go vet ./...`, Windows cross-build, and full Silesia byte-by-byte reconstruction. Any reported failures must be fixed before the source is described as validated.

No new compression-size or throughput results are reported. Historical JLZ R29 and Zstd/XZ measurements remain unchanged.
