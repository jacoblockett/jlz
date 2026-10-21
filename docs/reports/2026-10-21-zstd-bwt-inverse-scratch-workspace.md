# Zstd bounded workspace: native BWT inverse scratch reuse and group output placement

**Date:** 2026-10-21. **Status:** original native-Go complete-codec SCREEN; significant allocation reduction, **no credible end-to-end speed improvement qualified**; source left disposable. No production or PGO change.

## Mechanism and novelty

Inspired by bounded reusable workspace in [Meta Zstd `zstd_cwksp.h`](https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_cwksp.h) and the output lifetime decisions of [Zstd's decoder](https://github.com/facebook/zstd/blob/dev/lib/decompress/zstd_decompress_block.c). Historical JLZ Round 42 tested (and discarded) `sync.Pool` **for match-index buckets**, including generation stamps. This screen instead reuses **decoder-side BWT inverse LF arrays** (`column []uint16`, `lf []int32`, `pair []uint64`), which prior failed table-pooling experiments did **not** target; source/codec/decoder path is otherwise identical. A separate control variant writes completed decoded groups directly into the final output region instead of growing a redundant temporary group buffer, preserving disjointness among concurrent groups.

Local accepted production source Git blob hashes matched live `main` `4ea0ae3`: `bwt/transform.go` `422a24d8c5f537d240846c0e38d624e2d90d8fb2`, `engine/cross.go` `5f27fde8830b9099f6bb24769f7c041968f3ce50`. Root `default.pgo` untouched, Go1.23.2 linux/amd64. All 12 original corpus SHA-256s verified against canonical manifest. Every measured full corpus is complete 211,938,580 B→58,767,739 B and restores all twelve sources byte-exact. Pool trial preserved all twelve **individual archive SHA-256s** against independently generated unchanged controls.

The actual Go CPU/allocation profile (both encode and decode) on the unmodified full corpus showed 6.69 GiB allocated in aggregate (pprof rounded) and **1.76 GiB flat to `bwt.Inverse`**. Scratch-pool trial `sync.Pool` owns private per-call working slices for all LF preparations and returns them only after inverse completes; returned output is independently allocated. All rows and cells overwritten per invocation. No format or match grammar change. The independent output-placement trial directs decoded data into disjoint final group spans; source/decoder output unchanged.

## Diagnostic actual complete original corpus measurements

| Screen | Fully framed archived B | Complete 12-file encode/decode exact? | Allocation profile total (rounded) | Full original 12-file local timing |
|---|---:|---|---:|---|
| Current, unchanged source | 58,767,739 | PASS | **6.69 GiB** | Paired control encode **5.26724 /5.41012 s**, decode **1.94515 /1.94428 s** |
| BWT inverse scratch pool only | 58,767,739; all 12 SHA identical | PASS | **5.22 GiB**, roughly **22% less allocation** | In paired ABBA candidate encode **5.85515 /5.24067 s**, decode **2.03835 /1.79744 s** |
| Direct final-group destination only | 58,767,739 | PASS | **6.62 GiB**, only ~0.07 GiB less | One complete trial, timing not sufficiently repeatable for verdict |

`bwt.Inverse` flat allocations fell from **~1.76 GiB to ~0.175 GiB** in the pooled variant. This is **total allocated bytes in one instrumented test run**, *not* proof of peak RSS reduction, memory bound, fewer live objects, or time savings. `sync.Pool` may drop objects at GC and retains variable-size buffers; a production version would require a maximum pooled-region size or per-worker bounded ownership to satisfy JLZ's controlled-memory rule.

**Two interleaved ABBA controls (one 4-run group):** baseline encode median **5.33867869s**, pool median **5.5479074565s, +3.92% (worse)**. Pool decode two-run median **1.917893878s** versus baseline **1.944712371s**, **-1.38% nominal (better)**, but considerable host drift. Four further separately completed full original candidate-only encode/decode runs gave encode **5.089752798,5.125271192,4.948071298,5.020586876s**, decode **1.810030661,1.892920026,1.840227053,1.797104716s**, all 12 sources exact and all per-file archive hashes unchanged. Their 4-run medians **5.055169837s encode /1.825128857s decode** are **NOT a ten-run candidate median**, and unpaired historical timing comparisons would mislead. A subsequent fifth longer candidate process stalled/timed out and is **excluded**.

A diagnostic *decoder-only* whole-corpus loop (each pass includes all twelve original accepted files, but **not** a new encode on every pass) completed eight controls and eight pooled loops in separate processes: medians control **1.912874539s**, pool **1.898759028s**, just **0.738% nominal faster (better), inconclusive**. Another decoder-only process stalled after three loops and was excluded. Timing noise, intermittent stalled batch commands and lack of repeated paired advantage preclude speed promotion.

Full candidate `go test -count=1 ./...` and `go vet ./...` **PASS**; targeted Go race tests on engine/BWT **PASS**; Linux-host **Windows/amd64 CGO_ENABLED=0 cross-build PASS**, runtime not verified on Windows. Production archive backwards decoding/corrupt validation covered by full test suite; no format was changed. The separate output-placement variant also passed engine/bwt uncached tests, vet, and Windows cross-build.

## Decision

**Do not promote either disposable prototype.** Pool reuse demonstrates a real ~1.5 GiB saved allocation *per profile*, a worthwhile future architectural hypothesis, but **no qualifying ten-run integrated timing evidence**, and bounded-retention accounting remains outstanding. Direct output placement is less persuasive on total allocations. This explicitly does **not** close Zstd's bounded scratch architecture, per-worker scratch hoisting, or zero-copy block materialization. Next meaningful refinement should avoid arbitrary cached scratch growth and quantify **peak live heap, decoder speed under a controlled paired schedule, and GC churn** before a ten-run promotion candidate.

Historical retained accepted JLZ **58,767,739 B / 4.7873709325 s / 1.8721679705 s**; historical Zstd-3 **66,216,569 B / 1.2561s /0.8551s**; XZ-6 **49,233,340 B /116.9892s /2.3079s**. References not rerun or silently replaced; timings of external CLIs are not methodologically identical.
