# Round 42 — Two-entry native mathematical match buckets (ACCEPT H1)

**Date:** 2026-10-19. **Approved successor to Round 39.** The production encoder changes only match hash-bucket depth from four stored candidate offsets to the two most recent offsets, including ring mask and lookup loop. The unchanged decoder and `JLZ1` framing still reconstruct any accepted historical archive. No external codec/backend or user flags. The existing match-search reference regression test now uses the deliberately changed two-entry candidate budget.

## 1. Selected-path structural diagnostic

The first stage of Round 42 profiled the **complete original Silesia corpus only** (twelve independently sourced files totaling 211,938,580 bytes, SHA-256s all matched the Round23 manifest), with additional temporary counters in the unmodified Round39 matcher; diagnostic timing is **not** a new accepted baseline. Actual source files that choose mathematical regions dominate matcher work:

| Selected math source | Match index builds | Transformed indexes | Indexed bytes (MiB, approximate) | History reindexed (MiB, approximate) | `find` calls |
|---|---:|---:|---:|---:|---:|
| mozilla | 194 | 26 | 63.91 | 15.75 | 20,474,162 |
| ooffice | 25 | 14 | 7.46 | 1.50 | 3,689,910 |
| osdb | 41 | 0 | 13.61 | 3.62 | 3,692,468 |
| samba | 85 | 39 | 28.30 | 7.38 | 5,377,814 |
| **Selected math subtotal** | **345** | **79** | **113.28** | **28.25** | **33,234,354** |

The other eight files account for eighteen more index builds and one additional transformed index, used mainly in preflight math probes on BWT winners (no match indices on sao or x-ray). Across all 12, **363 table builds** and **80 transformed-table builds** were observed. `pprof` allocation-space data across a full encode **and decode** diagnostic attributed about **0.65 GiB** of cumulative allocation to `engine.newMatchIndexTransforms`, and CPU sampling showed `(*matchIndex).find` **4.11 s flat CPU** plus `addBefore` **1.64 s flat CPU**; these are independent cumulative CPU samples over concurrent workers, not total elapsed times.

## 2. Tested hypotheses and decisions

### A. Reuse cleared hash tables (REJECT)

Added bounded `sync.Pool` ownership of exact/transform match tables and cleared them on reuse. Every complete source file reconstructed and all **12 archives matched the unmodified reference SHA-256**; normal uncached Go suite and vet passed. One complete Silesia screening pass: **62,494,395 B / 6.677546 s encode / 2.512938 s decode**. Not persuasive against accepted 6.213451 s encode historical median. No 10-run promotional test; implementation discarded.

### B. Generation-stamped bucket reuse (REJECT)

Used explicit per-bucket generation stamps to avoid clearing all slots, including deterministic wrap-around reset; recycled privately owned index tables via `sync.Pool`. Same twelve byte-identical archive hashes, complete exact reconstructions. Uncached Go suite, vet, Windows cross-build and targeted race suite passed. One full Silesia screen: **62,494,395 B / 6.565329 s encode / 2.387977 s decode**. Extra generation checks and larger bucket structures do not show speed advantage. No ten-run promotion; discarded.

### C. Two-entry buckets (ACCEPT)

Replaced `matchBucket.positions [4]uint32` with `[2]uint32`, `bucket.next&3` with `&1` and restricted the `find` bucket search from four to two entries. The candidate is a **genuine compression-speed tradeoff**, not an exact compressed-output optimization: shorter bucket history can remove available references, changing chosen match tokens and size. Unlike an artificial hardcoded corpus gate, it is general to all native mathematical matching. Each bucket shrinks from **20 to 12 bytes (40% less)**, reducing table memory traffic and reference probes, with no additional abstraction or wire features. The optimized-vs-reference search regression test's deliberate algorithmic oracle depth was updated from four to two; the rest of the test suite was preserved.

Single full-Silesia screen: **62,741,534 B / 5.878322 s encode / 2.404287 s decode**. This warranted a 10-run candidate-only promotion benchmark.

## 3. Ten candidate-only complete original Silesia runs

Every accepted run processes and encodes **all twelve original source files**, independently serializing full JLZ1 archives including checksums and metadata, decoding them all and verifying `bytes.Equal` to original source. Every run's 12 original source identities match the canonical original Silesia manifest. Input reads occur before the timed operations. Four initial standalone process runs and six additional complete standalone process runs were recorded; incomplete/hung attempts were discarded, never included. Every **120/120** source reconstructs byte-exactly. Per run, the total archived size is identically **62,741,534 B**. This is one consistent single-process-per-complete-corpus methodology under `GOMAXPROCS=4`, Go 1.23.2, linux/amd64; no per-file-process totals mixed in.

| Run | Complete archive bytes | Encode s | Decode s |
|---:|---:|---:|---:|
| 1 | 62,741,534 | 6.299059194 | 2.464180442 |
| 2 | 62,741,534 | 5.980258698 | 2.412555147 |
| 3 | 62,741,534 | 6.125378099 | 2.341882043 |
| 4 | 62,741,534 | 5.855361024 | 2.400102500 |
| 5 | 62,741,534 | 5.993453311 | 2.433004774 |
| 6 | 62,741,534 | 5.899011299 | 2.423659672 |
| 7 | 62,741,534 | 6.511389219 | 2.531297202 |
| 8 | 62,741,534 | 5.931391553 | 2.400726832 |
| 9 | 62,741,534 | 5.472310170 | 2.192798487 |
| 10 | 62,741,534 | 5.855034594 | 2.390408476 |
| **Median, complete-corpus totals** | **62,741,534** | **5.955825126** | **2.406640990** |

### Full archive size regressions

Only four mathematical-region winners changed their archive bytes. Comparison with actual Round 39 source's known-size original Silesia archives: **mozilla +112,945 B (+0.610%), ooffice +13,492 B (+0.446%), osdb +28,122 B (+0.813%), samba +92,580 B (+1.831%; worst relative per-file regression)**. All other eight archive lengths unchanged, and complete per-file original-byte reconstruction passed. Across all files **+247,139 B (+0.39546%)**, worsening original-byte reduction by 0.11661 percentage points. Unlike Round36–39's historical provenance discrepancy, these **247,139 bytes genuinely come from the Round42 candidate**.

## 4. Required comparison and limitations

| Model | Compressed size | Reduction vs 211,938,580 B | Time to compress | Time to decompress |
|---|---:|---:|---:|---:|
| Zstd-3 (**historical external CLI**) | 66,216,569 B | 68.76% | 1.2561 s | 0.8551 s |
| XZ-6 (**historical external CLI**) | 49,233,340 B | 76.77% | 116.9892 s | 2.3079 s |
| Current approved JLZ R39 (**retained ten-run median**) | 62,494,395 B | 70.513% | 6.213451 s | 2.412670 s |
| Candidate JLZ R42 (**new ten-run median**) | 62,741,534 B | 70.396% | 5.955825 s | 2.406641 s |

| Candidate size delta vs current | Candidate encode time delta vs current | Candidate decode time delta vs current |
|---:|---:|---:|
| **+0.39546%** | **-4.14626%** nominal faster | **-0.24988%** nominal, unchanged decoder |

**Timing caveat:** Round39 reference was measured in an independent earlier session; these times are not a paired control experiment and do not establish a high-precision causal speedup. Zstd and XZ are historical external CLI numbers, *not directly comparable* with in-memory JLZ. The decoder has not changed and its tiny nominal time difference is measurement noise. Current Round39 actual-size archive baseline (62,494,395 B) is used throughout; never attribute the older historical Round29 four-file 150,268 B discrepancy to this candidate.

## 5. Final checks and decision

- Normal `go test -count=1 ./...` **PASS**.
- `go vet ./...` **PASS**.
- Targeted race tests `TestMatchV2OverlappingRepeatAndTransforms|TestParallelNativeIntegrity` **PASS**; full reference-equivalence oracle under race instrumentation exceeded the local runtime budget, while that same oracle passes uninstrumented. No claim of full-race-suite execution.
- `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build` **PASS**, build only, no Windows runtime benchmark.
- All 12 canonical original Silesia input SHA-256s verified and 120/120 original files reconstructed correctly across candidate's ten runs.
- Backward decoding regression suite `go test -count=1 ./...` **PASS**. Only encoder candidate-depth changes; decoder and archive framing untouched.
- Final production blobs are exactly those compiled and tested locally (Git blob hashes checked for source and adapted oracle). No private benchmark harness or binary committed.

**ACCEPT H1.** A native, broadly applicable two-entry hash candidate budget materially improves the measured 10-run encoding median by a nominal 4.15% with only 0.395% aggregate archive-size sacrifice and 40% less per-bucket memory. It is a favorable speed-first ratio/complexity tradeoff, albeit more modest than Round39 and based on unpaired source-reference sessions. The smallest/most readable implementation requires only a four-line matcher diff and a two-line intentional test oracle update. No fixed acceptance threshold. Do not continue tuning to depth=1/3 without new reason; the next research round should revisit selected mathematical match indexing architecture or identify a distinct information-theoretic complementary technique, not chase nearby one-value candidate-budget parameters.
