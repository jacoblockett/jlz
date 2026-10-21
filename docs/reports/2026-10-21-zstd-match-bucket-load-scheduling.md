# Meta Zstd hash-bucket load scheduling in native match search

**Date:** 2026-10-21. **Disposition:** REJECT these portable Go independent-load implementations as demonstrated improvements. Main native source, decoder, archive format and root PGO unchanged. Zstd cache-local matched prefetch *architectural family* remains open for genuinely different pipelined designs, not further cosmetic variants.

## Source and novelty gate

The Meta Zstd source audit in [zstd architecture coverage](../research/zstd-architecture-coverage.md) and [industry hotpath audit](2026-10-20-industry-encoder-hotpaths-audit.md) identified prefetch / independent hash lookup work in `zstd_fast.c` and `zstd_double_fast.c`. Native JLZ already uses a **two-offset packed 64-bit hash bucket** per exact/XOR/ADD transform. Round 42 already accepted reducing depth four to two for throughput. Previously failed hash-tag packing preserves two candidates but has inconsistent measured speed. The current screen is **neither a deeper bucket, new candidate hash, nor stored tag**: issue all independent exact/XOR/ADD bucket loads **before the first hash-candidate source-anchor comparison**, then use cached immutable buckets in the original exact, XOR, ADD priority and scorer. No `unsafe`, native assembly, prebuilt entropy engine, external codec or new archive mode.

Trial A stores three preloaded buckets in a local `[3]matchBucket`; B stores explicit `exactBucket`, `xorBucket`, `addBucket` locals and dispatches at the original loop. Both preserve index-building order, newest/second-newest offsets, repeat-distance selection, transformed mask verification, lazy parse and original decoded contents. Go does not expose a portable nonfaulting software prefetch primitive here; **preload is not proof of actual CPU cache-miss latency overlap**. This test does not claim to implement the upstream hardware instruction.

## Experimental method and source identity

Production source snapshot `8566c8a` matched **live `main` `ece182c`** exactly for Git blobs: `engine/match.go` `6c9714114a056bf49f8124f81826f0eded68a5d8`, `match/hash.go` `a640f464bfc08ee17ced6de97b78819196ce69d9`, `match/gate.go` `c63ca969296b48e3e87eacd2da6d2d7ed62f5956`, and `default.pgo` `382b3c59a75bfa11bc60968f71cdde67dbe8d634`. Intervening live commits were research documentation only.

The **twelve complete original Silesia files** have **211,938,580 input B**, every `sha256sum -c docs/reports/2026-10-18-round23-silesia-sha256.txt` PASS. Disposable separate Go builds on Go1.23.2 linux/amd64, native `compressFast` and `decompress`, `GOMAXPROCS=4`, root PGO present. **All twelve candidate archive lengths and SHA-256 digests matched control exactly** in explicit independent per-file checks (not merely aggregated bytes), and every original reconstructed `bytes.Equal`. Complete archive total **58,767,739 B**, no change to any per-file size or digest.

## Complete original-corpus measurements

**Candidate B** qualified for measurement through ten complete original-Silesia encode+decode runs. Each of the ten independently initiated candidate-only runs encoded all 12 with full existing JLZ framing/checksums and decoded every original correctly; all reported 58,767,739 B. These are real candidate measurements, not a source-component benchmark.

| Run | Candidate B encode seconds | Candidate B decode seconds |
|---:|---:|---:|
| 1 | 5.438354802 | 1.947070576 |
| 2 | 5.380085768 | 1.908715182 |
| 3 | 5.187258079 | 1.852685926 |
| 4 | 5.325327827 | 1.974569993 |
| 5 | 5.268709388 | 1.904298149 |
| 6 | 5.889505591 | 2.110616077 |
| 7 | 5.421674084 | 1.971154812 |
| 8 | 5.403073832 | 1.905145353 |
| 9 | 5.043929782 | 1.878832717 |
| 10 | 5.119291777 | 1.826390048 |
| **Median** | **5.3527067975** | **1.9069302675** |

The retained historical ten-run **accepted** JLZ reference remains **58,767,739 B / 4.7873709325s encode / 1.8721679705s decode** and was NOT replaced with a new local anchor. Historical external Zstd-3 **66,216,569 B / 1.2561s /0.8551s** and XZ-6 **49,233,340 B /116.9892s /2.3079s** were NOT rerun and have different CLI vs in-memory methodology. Nominal candidate B vs older accepted JLZ medians +11.81% encode /+1.86% decode **(both nominally worse, unpaired and not causal)**. No archive size delta.

To resolve major host/time drift, ran actual separate unchanged control binaries during the same experimental period. Six earlier completed control encodes [6.111342921,5.645058688,5.702649685,5.927912465,5.848688798,5.420669984]s had a median 5.7756692415s, misleadingly 7.32% slower than candidate B median. However three more controls later measured **5.090649654 / 5.381996117 / 5.432461280** seconds. This shows environment drift comparable to or greater than the alleged speedup.

Two independently ordered **ABBA four-run full-corpus** paired diagnostics, every archive exact:
- Control **5.251460424**, B **5.246419733**, B **5.048573646**, control **4.937964969** seconds. Control median **5.0947126965**, candidate median **5.1474966895**: **B +1.04% slower (worse)**.
- B **6.126207138**, control **5.806980705**, control **5.948959929**, B **5.282744552** seconds. Control median **5.877970317**, candidate median **5.704475845**: **B -2.95% faster (better)**.

The opposite signs, drift, and unrelated unchanged decoder timing fluctuation defeat a credible speed-promotion claim. **Trial A** passed exact complete archive screens but displayed the same within-session timing overlap (e.g. baseline 5.927912465 vs A 5.929686071). It did not warrant ten candidate runs. A few extended batch commands stalled/timed out; they are excluded and **not** counted as completed full passes.

## Validation and conclusion

Complete accepted native test suite `go test -count=1 ./...`: **PASS**. `go vet ./...`: **PASS**. `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build`: **PASS**, build only. Targeted race `TestMatchV2OverlappingRepeatAndTransforms`, `TestParallelNativeIntegrity`, `TestExistingArchiveCompatibility`, and `TestCorruptionDetected`: **PASS**, after separating race tests into bounded groups; a broad race command with `TestShortCircuitMatchSearchEquivalence` timed out, therefore no claim of that full race group. Existing old-archive, corruption and truncation coverage passed in unchanged test suite. Decoder not modified. All twelve per-file archive hashes and all 120 candidate-only roundtrip checks passed.

**Reject A/B for production:** exact candidate parity alone does not establish that early Go bucket loads produce a throughput win, and this two-candidate design has very few serial loads to overlap. Do not retune array-vs-variable naming, hash bit widths, prefetch offsets or bucket depth as new research. Future Zstd cache-miss scheduling should prove a different algorithmic dependency overlap in an actual CPU trace or a batched sequence store, rather than another small loop rewrite. Continue Zstd inspiration, especially open adaptive candidate work budgeting, native sequence fusion and genuine *precomputed match-graph* reuse when a bounded cost model exists.
