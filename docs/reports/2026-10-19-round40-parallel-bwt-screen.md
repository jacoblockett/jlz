# Round 40 — CPU-aware concurrent BWT preflight (REJECT)

**Date:** 2026-10-19. **Verdict: REJECT H1.** Experiment locally in disposable copy of the approved Round 39 source; no codec source promoted. GitHub Actions were not used. The current approved codec remains Round 39 (`72f47d5`). All experiments use native standalone JLZ code only.

## Motivation and exact design

`engine/bwt.go:sampleBWTPreference` compares three first/middle/last 256 KiB independently encoded BWT blocks against the previously constructed mathematical-region sample baseline. It requires at least two samples to win individually and a total BWT payload size at most 95% of the complete mathematical sample payload total. An exact short-circuit exits after any prefix cannot possibly win. A Round 40 diagnostic using the full original 12-file Silesia found this function was reached in eleven files: **seven BWT-preferring cases process all three samples and pass** (`dickens`, `mr`, `nci`, `reymont`, `webster`, `xml`, `x-ray`); **four cases reject**, three after two samples (`mozilla`, `ooffice`, `samba`) and one after three (`osdb`); `sao` never enters this preflight.

Two executable alternatives were made:
- **Two-way lookahead** computes the first two BWT samples concurrently, retains the third sample conditionally so the original early rejection can still avoid it. One uncontended full-Silesia screening run: **62,494,395 B; 6.020049 s encoding; 2.259449 s decoding**. Screening only, not ten-run promotional evidence.
- **Worker-limited parallel lookahead**, `workers = availableWorkers(3)`: one worker retains identical serial early abort; two workers compute samples 0 and 1 concurrently with the third conditional; three workers compute all three independent samples concurrently. Results are consumed and vote/budget comparisons applied **in original sample order**. It may encode an otherwise skippable losing third sample on three-worker systems, so a speed advantage is not guaranteed. No source encoding format, archive bytes, metadata or decoder change.

Do not reuse earlier simultaneous dual-candidate Go test screens; their ~12-second encode totals were distorted by running two multi-worker compressors simultaneously and are explicitly excluded. The three-way variant's first independent uncontended screen was 5.727477 s encoding / 2.374370 s decoding, which also did not establish an effect and motivated ten complete runs before promotion.

## Candidate ten complete original Silesia runs

All 12 files from the complete, original **211,938,580 B** Silesia were used on each run. Input reads excluded from measured in-memory encoding/decoding. `TestRound39Profile` times twelve `compressFast` calls, then twelve decodes and compares exact original bytes. Every run has full metadata and checksum and produces **62,494,395 B**. All **120 per-file archive SHA-256 hashes** match the previously recorded untouched Round 39/source hashes, and every source reconstruction passes. Ten independent single-process corpus runs were completed; no incomplete run enters the medians. The last orchestration shell command hit its wrapper timeout after writing run 10's complete output log; that run has a complete corpus total and all twelve archive hashes, and was included for that reason.

| Complete run | Archive bytes | Encode s | Decode s |
|---:|---:|---:|---:|
| 1 | 62,494,395 | 6.192124364 | 2.545957666 |
| 2 | 62,494,395 | 6.136203326 | 2.373582651 |
| 3 | 62,494,395 | 5.962001063 | 2.335093630 |
| 4 | 62,494,395 | 6.294047720 | 2.468210208 |
| 5 | 62,494,395 | 5.892691654 | 2.305393976 |
| 6 | 62,494,395 | 5.953421342 | 2.375771129 |
| 7 | 62,494,395 | 6.450095628 | 2.498370078 |
| 8 | 62,494,395 | 6.359111939 | 2.366331145 |
| 9 | 62,494,395 | 6.195951337 | 2.472605199 |
| 10 | 62,494,395 | 6.321011897 | 2.433007149 |
| **Median of ten whole-corpus totals** | **62,494,395** | **6.194037851** | **2.404389139** |

## Required comparison

| Model | Complete compressed size | Reduction vs 211,938,580 B | Compression time | Decompression time |
|---|---:|---:|---:|---:|
| Zstd-3, historical external CLI | 66,216,569 B | 68.76% | 1.2561 s | 0.8551 s |
| XZ-6, historical external CLI | 49,233,340 B | 76.77% | 116.9892 s | 2.3079 s |
| **Current JLZ Round 39**, historical ten-run candidate-only median | **62,494,395 B** | **70.51297%** | **6.213451 s** | **2.412670 s** |
| **Candidate JLZ Round 40**, newly measured ten-run median | **62,494,395 B** | **70.51297%** | **6.194038 s** | **2.404389 s** |

| Archive size delta | Encode-time delta | Decode-time delta |
|---:|---:|---:|
| **0.00%** (all 12 exact bytes) | **-0.3124% nominal** (unpaired) | **-0.3432% nominal** (unpaired, decoder untouched) |

Retained Round 39 results were measured in an earlier, independent session; differences this small are not demonstrated causal speedups. Retained Zstd/XZ numbers are external CLI measurements using different lifecycle, hardware and I/O timing methodology and cannot be interpreted as matched-host throughput. Round 29's older 62,644,663 B benchmark record has the already documented four-file provenance discrepancy; do not substitute it as current active size baseline. Worst compressed-size regression: **zero bytes on each of 12 files**. Per-file encode-time regressions were not established by new paired old-vs-candidate timings; a single cross-session profile log is insufficient.

## Validation

- `go test -count=1 ./...` **PASS** on candidate, including all existing tests.
- `go vet ./...` **PASS**.
- `go test -race -count=1 ./engine -run 'TestBWTPreflightWorkerEquivalence|TestEarlyRegionalWorkerEquivalence|TestParallel'` **PASS**.
- `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build` **PASS** (cross-build only).
- New disposable test `TestBWTPreflightWorkerEquivalence` compares serial vs 2/3/4-worker BWT decisions and encoded sample bytes on deterministic data and an explicitly forced early-rejection baseline.
- Full original Silesia exact reconstructions **120/120** and candidate-vs-unchanged Round39 per-file archive digests **120/120**.
- No decoder changes; existing `JLZ1` compatibility remains as in Round 39. No candidate code left on `main`.

## Decision

**REJECT H1:** Parallel BWT preflight adds goroutine scheduling and, on a three-core workload, performs a third BWT transform even when the first two already prove rejection. Ten complete candidate medians show only a nominal **0.31% encode-time difference** against the retained unpaired Round 39 figure, far too small to justify the new work/memory/complexity. Preserve this finding and **do not spiral into BWT sample thresholds or worker-count tuning**.

Next substantive work should target **the winning compressed-region mathematical matcher and its unavoidable tokenization work** (exact indexing/scan redesign with corpus profiling), not further candidate preflight parallelism. No production code commit from this experiment.
