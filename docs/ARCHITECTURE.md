# Source organization

JLZ has one standalone executable, with `main.go` at the repository root. Go packages are organized by actual dependency boundaries, not by filename prefix. The repository does not track compiled executables.

```text
main.go                    two-command CLI (source, not a binary)
engine/                    automatic selection and JLZ1 archive integration
archive/varint.go          shared wire framing primitive
additive/words.go          reversible 16-bit modular differences
bwt/transform.go          block sorting and move-to-front/zero-run
entropy/arithmetic.go     native adaptive context arithmetic coding
lane/block.go             byte-lane transforms and entropy estimates
match/                    fingerprints, transformed-match probes, LZ expansion
predictive/transform.go   reversible lag predictors
workers/                  CPU quota-aware bounded parallelism
engine/testdata/fixtures/ byte-exact correctness fixtures
docs/                     research, reports, engineering and inactive workflows
```

Each directory is a real Go package. `engine/` owns high-level archive construction, encoder selection, legacy-format decoding and interoperability between formats. Feature packages provide independent transforms or policy without importing `engine`. Existing unit/regression tests that exercise internal engine symbols remain alongside `engine` implementations. Wire formats, the two CLI verbs and existing selection logic are unchanged.

Build locally from the repository root:

```sh
go test -count=1 ./...
go vet ./...
go build -o jlz .
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o jlz.exe .
```

The build output is ignored by Git. Do not place binaries in source control, use GitHub Actions, or claim unrun tests as passing. Full candidate throughput benchmarking continues to use only the original 12-file Silesia corpus and stored reference measurements.


Round53 adds `engine/finite_state.go` for a bounded, self-contained, byte-renormalized 12-bit native finite-state coder with exact packed-stream subtype 5. The `engine/` orchestrator owns selection and decoding because this codec is only called for selected intermediate BWT and mathematical streams. Existing `entropy/` arithmetic and raw/RLE/Huffman formats remain independently implemented and backwards-readable; do not substitute an external compression codec backend.

**Rejected architecture note:** Round54 briefly added a dual-anchor fast matcher but was rejected because it increased whole-Silesia archive bytes by 1.68985% despite faster encoding. The experimental `engine/match_dual_anchor.go` and associated test were removed, and R53 mathematical encoding code restored. **No R54 parser remains in production; current approved JLZ is Round53.** See `docs/reports/2026-10-19-round54-dual-anchor-speed.md`. New speed optimizations must not sacrifice compression ratio.

Round55 corrects the BWT preference sampling in existing engine/bwt.go to score **actual bijective mode15 payloads**, retain their bytes, and use a 2% / two-of-three position evidence margin. No new codec package, format subtype, or decoder work is added. Rejected R54 dual hash remains absent. See docs/reports/2026-10-19-round55-bwt-preflight-recovery.md.


Round57 adds `engine/bwt_subregions.go`, a bounded reversible BWT model-region selector and decoder for new internal block mode16, and `engine/bwt_subregions_test.go`. Region-level math/BWT preflights stay unchanged. Subregions are only attempted for BWT-winning non-sample blocks, and actual byte lengths control fallback to existing mode15. Older JLZ1 archives remain readable; old binaries cannot read new mode16. Source-backed research: `docs/reports/2026-10-19-round57-mdl-bwt-regimes.md`.

Round60 expands the selected native BWT block horizon to **512 KiB** without changing its BWT/bijective/entropy transform or the 256 KiB mathematical region pipeline. The mixed-cross archive decoder supports validated variable BWT/raw block lengths, prefix-sum output positioning and per-block BWT parallelism; history-dependent blocks retain four-region grouping. No new mode ID, but pre-R60 decoders reject new larger frames. Research and format evidence: docs/reports/2026-10-19-round60-bwt-region-scale.md. The reversible prior-rank-label remapping family is closed as too small a component gain.
