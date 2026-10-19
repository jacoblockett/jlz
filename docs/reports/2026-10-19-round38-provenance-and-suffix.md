# Round 38 — Historical size provenance and direct native BWT speed prototypes

**Date:** 2026-10-19. **Decision: REJECT H1.** No native codec edits were promoted. Archived code on `main` remains the accepted implementation. Tests and builds after removing all disposable prototypes succeeded locally; GitHub Actions were not invoked.

## 1. Resolve the 150,268-byte reference discrepancy to specific source archives

The current production source from the user-supplied `jlz-main.zip` (GitHub `708f864`), with the already-promoted fixture-path-only correction `28865fb`, was used for one **untimed correctness-and-provenance audit**, not a new baseline timing benchmark. The independently supplied complete original 12-file Silesia (211,938,580 B) matches the Round 23 SHA-256 manifest for all twelve files. All twelve original files were reconstructed byte-exactly and all twelve current archive SHA-256 identities were captured. Total source archive output is **62,494,395 B**.

Historical Round 29 total: **62,644,663 B**. The entire **150,268 B** delta comes from four source archives:

| File | Historical Round 29 source size (retained Round 23/25 size) | Current unchanged source size | Difference |
|---|---:|---:|---:|
| mozilla | 18,623,280 | 18,521,909 | -101,371 |
| ooffice | 3,037,854 | 3,023,969 | -13,885 |
| osdb | 3,472,353 | 3,459,257 | -13,096 |
| samba | 5,077,310 | 5,055,394 | -21,916 |
| **Sum** | **30,210,797** | **30,060,529** | **-150,268** |

The current output of each of these four files exactly matches the **historical Round 21** published per-file size; the later Round 23 report recorded the larger values as its unchanged archive sizes. The remaining eight files, including Round 29's additive x-ray archive at **4,266,987 B**, match the Round 29 expected sizes. Side-by-side connected GitHub blob inspection verified that `fast.go` (now `engine/compress.go`), `region_first.go`, `math_region_fast.go`, `math_regions.go` and `native_packed.go` are **byte-identical to accepted Round 29 after replacing the package declaration**. The Go package restructure and subsequent fixes do not touch these five encoder routines.

**Interpretation:** the numerical discrepancy is now *fully localized*, not a Round 36/37/38 improvement. The ultimate original-run/source-state cause is **not proven**; the Round 23 record itself notes earlier historical Silesia discrepancies despite manifest-hash checks. Preserve both the historical approved ten-run timing/size reference and the newer measured current-source archive SHA digests. Do not claim an unexplained 150,268-byte algorithm win or silently rewrite accepted measurements.

## 2. Candidate experiments

### A. Direct suffix array, two-key radix prefix-doubling

Implemented independent stable two-key radix sorting with doubling equivalence classes to replace the standard-library suffixarray `Write` serialization/parsing in `bwt.Transform`. Verified suffix ordering exactly against a naive lexicographic reference for uniform, periodic, consecutive and pseudorandom correctness vectors (1–4096 bytes). An untimed whole-Silesia executable audit exceeded a **160-second** execution deadline before completing; no complete full-corpus candidate measurement or archive proof. **Reject immediately**: superlinear construction is impractical against the existing standard-library suffix sorter. No tuning cycle or promotion.

### B. Floating-point exact arithmetic probability quotients

Replaced two hot-loop divisions by equivalent floating-point quotients followed by integer conversion in the native adaptive context coder. For every numerator 0..denominator and every denominator 1..2048, the candidate's quotient matched the original exact `uint64(numerator*65536/denominator)`: approximately 2.1 million combinations. All normal targeted Go tests passed. A complete untimed whole-Silesia source/byte audit yielded **62,494,395 B**, with twelve byte-identical archive SHA-256 digests and exact reconstruction. However the audit took **33.53 s** wall clock compared with **17.04 s** in the preceding stock-source correctness audit, independent runs and not a controlled pair. No indication of speed benefit. **Reject**, without the unnecessary ten-run promotional benchmark.

### C. Direct SA-IS induced suffix sorting

Implemented a native linear-time suffix sorter using L/S classification, LMS substring naming, reduced recursion and stable bucket induction; bypasses `index/suffixarray.Index.Write` and private serialized index parsing. Verified exact lexical suffix indices against standard reference sorting on the same correctness vectors. One complete Silesia source/archive audit passed: **62,494,395 B** and exactly the same per-file archive SHA-256 hashes as unmodified source, with all originals restored. Candidate-only, **one** complete-corpus timed screen (summed 12 complete input encodes and decodes):

- **Encode 13.474531979 s**, decode **4.662483447 s**, total output **62,494,395 B**.
- These are single-run screening times, not ten-run medians; comparing against historical Round 29 medians on a different session/host is not a controlled effect estimate.
- No substantial measured gain. Independent Go SA-IS implementation and supporting array allocation add significant complexity to a solved library routine; **reject** rather than tuning.

## 3. Required full-corpus comparison (screening only, not a candidate promotion)

| Model | Compressed size | Reduction from original | Encode time | Decode time |
|---|---:|---:|---:|---:|
| Zstd-3, historical external CLI | 66,216,569 B | 68.76% | 1.2561 s | 0.8551 s |
| XZ-6, historical external CLI | 49,233,340 B | 76.77% | 116.9892 s | 2.3079 s |
| Approved JLZ Round 29, **historical** ten-run median | 62,644,663 B | 70.44% | 7.516364 s | 2.821363 s |
| Candidate direct SA-IS, **one new screening pass** | 62,494,395 B | 70.51% | 13.474532 s | 4.662483 s |

| Compression-size delta | Compression-time delta | Decompression-time delta |
|---:|---:|---:|
| **0.0000% vs actual current source**; -0.2399% vs historical is *not attributable* | +79.27% nominal vs independent historical median, **not paired or promotional** | +65.26% nominal vs independent historical median, **not paired or promotional** |

All historical codec benchmarks were reused, not rerun. The Silesia-only benchmark policy was respected; other data was used only to test suffix-array correctness. A single candidate screen cannot be promoted.

## 4. Validation and next priority

After restoring unmodified production files and deleting disposable tests/prototypes, `go test -count=1 ./...`, `go vet ./...`, and `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build` all **passed locally**. No production source changed on main. No candidate-only ten-run benchmark was justified by negative screening evidence. No Actions usage.

**REJECT H1.** Avoid more variants of standalone suffix sorting and arithmetic division unless profiling on *actually winning* full-corpus encoded paths identifies concrete headroom. The only defensible next step is a detailed selected-mode CPU allocation/profiling audit (encoding-only, complete original Silesia) followed by one substantial, bounded native algorithm redesign with measured whole-corpus benefit, rather than repeat microtuning or competing against already high-performance standard-library SA-IS.
