# Round 54 — Fast dual-anchor exact parsing with algebraic guard (REJECT H1, production rolled back)

**Date:** 2026-10-19. **Final decision: REJECT H1, explicitly rejected by the user because compressed archive size regressed.** An earlier provisional ACCEPT decision was overturned. The entire Round54 fast-parser source was removed from main and the exact prior R53 `engine/region_fast.go` blob restored. **Round53 is the current approved production baseline.** Native stand-alone Go only; no Zstd/DEFLATE/XZ/Brotli/LZMA backend, imported compressor, copied compressor source, extra CLI verbs, flags or new wire subtype.

> **Final disposition (user correction, 2026-10-19): REJECT.** Previous provisional acceptance of Round54 is void. Ratio regression is not permitted regardless of speed gain. The R54 parser and test files were deleted from production and R53 source restored; the figures below document a rejected experiment, not an approved codec.

## Hypothesis and actual implementation

The Round52 Zstd-3 architectural audit identified a structural alternative to expensive per-byte match probing, not another minor fingerprint/bucket adjustment. Implemented from scratch an exact-byte fast parser with **two independent hashes**, one 4-byte anchor and one 8-byte anchor, each holding one recent position in a 65,536-entry table. Both candidates are checked for at least five exact bytes before emitting a match. Long matches compare eight bytes at a time; backward extension consumes pending literals; recent-distance opcodes use JLZ's already existing three-distance scheme. Miss-dependent bounded position advancement (1..16 bytes, rising over long misses) saves index/update work. At most the existing 128 KiB region history is searched and the existing native match token format is emitted. No learned model, metadata or decoder change.

**Automatic hybrid selection rather than wholesale replacement:** Before parsing a mathematical region, reuse the existing deterministic `likelyTransformedMatches` structural probe. If it finds a verified XOR/ADD transformation, retain JLZ's full previous mathematical matcher; otherwise call the new exact-only dual-anchor parser. This affects the actual mathematical regional encoder and history-backed groups only, with no user switch. The already approved BWT/finite-state/additive streams and all legacy codec readers remain intact. Neither path calls external compressors.

**Source:** `engine/match_dual_anchor.go` and the two `tokenizeRegionHybrid` calls in `engine/region_fast.go`. Focused exact-token and history-aware regression tests are in `engine/match_dual_anchor_test.go`. Original candidate implemented in disposable local source and temporarily committed to `main` after validation; all R54-only Go source and tests were subsequently reverted after the user rejected the compression loss. Research measurements are retained for provenance.

## Candidate screening, all on twelve complete canonical Silesia files

Each initial screen compressed and decompressed **all 12 complete original Silesia files, 211,938,580 B** (full JLZ archive metadata), with exact source restoration. These are **one complete-run screens**, not medians. Timings are noisy and should not be used to rank near ties.

| Experimental structural scope | Actual complete archived bytes | Encode s, one full screen | Decode s, one full screen |
|---|---:|---:|---:|
| All mathematical regions use fast exact-only dual-anchor parser | 62,243,308 | 4.963739 | 2.325917 |
| Fast exact parser only in cross-history blocks | 61,887,279 | 5.264079 | 2.454822 |
| Fast exact parser only in local/no-history blocks | 61,246,614 | 5.642692 | 2.217879 |
| Exact-only parser with costly mathematical fallback if its stream is near raw | 62,238,309 | 4.730844 | 2.306762 |
| **Selected hybrid: fast exact unless transformed-relationship probe succeeds** | **61,866,415** | **5.054135** | **2.191824** |

**Why the candidate was selected:** Full fast-only incurs +1.405 MB against R53. Cross-only retains most of the growth for less speed benefit; local-only keeps sizes closer to R53 but shows no compelling encode gain in its one screen; post-packing fallback adds work and recovers too few bytes. The hybrid offers the strongest evidence of an actual, scalable speed-versus-size tradeoff without losing the algebraic capabilities where the probe finds transformed structure. No 10-run medians are attributed to rejected alternatives. Distinct threshold tuning and multiple further scopes were not pursued.

## Complete R54 ten-run candidate-only benchmark

All ten independent complete Silesia encode/decode passes used the same precompiled Go candidate. Each pass visited **every byte** of each original file, charged complete archived sizes, then restored all twelve files exactly, checksums included. Four batches of 3, 3, 2, 2 passes. No incomplete runs included. Every run emitted exactly **61,866,415 B**.

| Run | Archive bytes | Encode s | Decode s |
|---:|---:|---:|---:|
| 1 | 61,866,415 | 5.039710098 | 2.237432441 |
| 2 | 61,866,415 | 5.266992085 | 2.321942860 |
| 3 | 61,866,415 | 5.240587551 | 2.375735899 |
| 4 | 61,866,415 | 5.041347206 | 2.272801247 |
| 5 | 61,866,415 | 5.034969048 | 2.260773088 |
| 6 | 61,866,415 | 5.077787797 | 2.247981709 |
| 7 | 61,866,415 | 5.987941910 | 2.538016707 |
| 8 | 61,866,415 | 5.711238229 | 2.443201754 |
| 9 | 61,866,415 | 5.225700291 | 2.414589569 |
| 10 | 61,866,415 | 5.485542955 | 2.380904344 |
| **Median (separate columns)** | **61,866,415** | **5.233143921** | **2.348839380** |

Current/previous approved R53 reference was NOT rerun: **60,838,336 B / 5.721664347 s encode / 2.346933746 s decode** (historical independent candidate-only ten-run medians). R54 costs **+1,028,079 B (+1.68985% larger compressed archive)** in exchange for **-8.5381% nominal encode time**, and **+0.0812% nominal decode time (neutral)**. The size regression is genuine; encode/decode times compare independent sessions and cannot establish an exact causal percentage. The 8.54% nominal encoding gap exceeds recent tiny 0.3–1.3% differences but is not a controlled paired benchmark.

### Required reference comparison

| Model | Complete compressed size | Reduction vs 211,938,580 original bytes | Compression time | Decompression time |
|---|---:|---:|---:|---:|
| Zstd-3, historical external CLI | 66,216,569 B | 68.76% | 1.2561 s | 0.8551 s |
| XZ-6, historical external CLI | 49,233,340 B | 76.77% | 116.9892 s | 2.3079 s |
| **Current approved R53, retained 10-run median** | **60,838,336 B** | **71.29436%** | **5.721664 s** | **2.346934 s** |
| **R54 candidate, NEW 10-run median** | **61,866,415 B** | **70.80927%** | **5.233144 s** | **2.348839 s** |

| Compressed-size delta candidate vs R53 | Encode-time delta candidate vs R53 | Decode-time delta candidate vs R53 |
|---:|---:|---:|
| **+1.68985%** (worse ratio) | **-8.53808%** (nominal faster) | **+0.08120%** (neutral) |

External CLI Zstd/XZ timings use different conventions from JLZ in-memory measurements and should not be presented as like-for-like throughput on one measurement host.

### Per-file original-byte exact archive sizes and regressions

| File | R53 retained archive B | R54 archive B | Delta B |
|---|---:|---:|---:|
| dickens | 3,213,620 | 3,213,620 | 0 |
| mozilla | 18,609,139 | 19,457,491 | **+848,352 (+4.56%)** |
| mr | 2,831,212 | 2,831,212 | 0 |
| nci | 2,169,196 | 2,169,196 | 0 |
| ooffice | 3,033,712 | 3,073,787 | +40,075 |
| osdb | 3,472,004 | 3,423,835 | **-48,169** |
| reymont | 1,469,844 | 1,469,844 | 0 |
| samba | 5,134,262 | 5,322,083 | +187,821 |
| sao | 6,261,713 | 6,261,713 | 0 |
| webster | 9,851,849 | 9,851,849 | 0 |
| xml | 524,798 | 524,798 | 0 |
| x-ray | 4,266,987 | 4,266,987 | 0 |
| **TOTAL** | **60,838,336** | **61,866,415** | **+1,028,079** |

This is explicitly a ratio sacrifice and **not** a universal win: three files grow, one shrinks, eight unchanged; worst file-size regression is mozilla **+4.56%**. R54 remains smaller than the historical Zstd-3 aggregate by 4,350,154 B but slower by historical unpaired CLI-vs-in-memory timings. The user rejected exchanging compression ratio for throughput. R53 is restored as production, and R54 is preserved only as negative research evidence.

## Validation, format, compatibility

- Original full Silesia corpus SHA-256 manifest: **12/12 PASS**.
- New candidate: all **120** complete source-file reconstructions exact in the 10-run benchmark.
- `go test -count=1 ./...` PASS; `go vet ./...` PASS; targeted `go test -race -count=1 ./engine -run 'Test(DualAnchorExact|FiniteStateExact|CrossGroup)'` PASS; Windows/amd64 `CGO_ENABLED=0` cross-build PASS (not Windows native execution).
- Dedicated exact local and cross-history token tests PASS, validating accepted on-wire token semantics.
- **Compatibility in both directions:** standalone R53 CLI-encoded full original `xml` was byte-exactly decoded by R54, and standalone R54-encoded full original `mozilla` was byte-exactly decoded by the older R53 CLI, including SHA-256 checks. No new archive codec, mode or decoder code is required.
- No accepted JLZ/Zstd/XZ speed benchmark was repeated. Experiment screens are candidate-only; the retained R53 per-file archive sizes and medians came from Round53. All original BWT/natively selected file archives in the final candidate retain their earlier SHA-256 hashes.
- Current source uses two 65,536-slot `uint32` hash tables per fast matching call (~512 KiB index payload), compared with the previous larger mathematical tables when transformation search is unnecessary. No novel entropy theorem claimed: short/long hash anchors, lazy matching tradeoffs, and skip-after-miss parsing are established concepts applied natively to JLZ.

## Decision and next experiment

**REJECT H1.** Despite nominally faster encoding, Round54 worsens whole-Silesia archive size by **1,028,079 B (+1.68985%)** and grows the `mozilla` file archive by **4.56%**. Under the user's clarified non-negotiable rule, compression cannot be sacrificed for speed. The Round54 trial was fully rolled back from production; **R53 remains the active accepted codec**, and its historical benchmark remains the reference. No ten-run candidate promotion is permitted for any subsequent candidate that increases aggregate original-Silesia archive bytes.

**Next materially distinct research:** investigate low-work mathematical substream prediction using the new finite-state packed metadata, or a bounded algorithmic region classifier that can select a single parser with a stronger byte-cost estimate without encoding both full candidates. Do not immediately repeat stride size, bucket depth, finger-print microtuning or arbitrary threshold sweeps. Better ratio recovery must show actual end-to-end savings that justify candidate runtime.
