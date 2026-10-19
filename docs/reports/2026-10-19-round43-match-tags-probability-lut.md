# Round 43 — Tagged native match-index entries and exact arithmetic lookup (REJECT H1)

**Date:** 2026-10-19. **Production decision: REJECT both candidates, no Go source changes.** Round 42 remains approved. All measurements used **only the original complete 12-file Silesia corpus**: 211,938,580 uncompressed B; all twelve SHA-256 identities checked against `docs/reports/2026-10-18-round23-silesia-sha256.txt`. Native Go 1.23.2/linux/amd64. Prior approved JLZ/Zstd/XZ results were not rerun. Tests and all prototype sources ran locally with files supplied by the user; GitHub Actions were not used.

## Profile-led diagnosis

Reconstructed the **exact R42 production source** from the user-uploaded `jlz-main.zip` plus the previously committed R39 and R42 changes, verified against Git blob SHA-1 for `engine/match.go` (`ff344090...`), match-oracle test (`ebded82...`), R39 preflight (`0c3b8658...`) and historical fixture-path correction (`80c8f197...`). Normal uncached Go tests passed before any changes. An encoding-only CPU profile over complete original Silesia (profiling diagnostics, **not** a new accepted baseline benchmark) yielded 24.29 s of aggregate sampled CPU over 7.25 s wall clock, showing `matchIndex.find` **2.94 s flat / 4.01 s cumulative**, `matchIndex.addBefore` **1.28 s flat**, `entropy.bitModel.probability` **1.64 s flat**, `entropy.acEncoder.bit` **1.69 s flat**, and `entropy.EncodeArithmetic` **6.42 s cumulative**. Parent/call-stack cumulative times overlap and cannot be added. These figures motivated two independent experiments.

Instrumented `matchIndex.find` with disposable per-index counters. Across complete Silesia, **34,885,363 find calls**, **49,762,021 exact bucket probes** (21,751,224 passed anchor verification), **11,206,861 XOR probes** (180,010 passed), and **11,202,597 ADD probes** (212,967 passed). Thus **22,409,458 transformed bucket probes** but only **392,977 (1.75%)** passed their initial anchor; this is an exact profile-led pruning opportunity rather than an assumed distribution. This diagnostic counter instrumentation was discarded and never promoted.

## Candidate A — Pack fingerprint tags into unused high bits of stored positions

Each mathematical-region source slice is normally <2^19 bytes; current match buckets store an offset as a `uint32` even when the absolute offset needs only 19 bits. Prototyped a *same-footprint*, dynamically sized representation: compute `positionBits = bits.Len(uint(len(input)))`; write the original positive offset `(at+1)` into low bits and a high-bit tag of the exact original four-byte match fingerprint above it. Matching first five anchor bytes (plain equality or consistent XOR/add mask) **necessarily implies matching the same fingerprint**. A mismatched fingerprint tag therefore proves that the candidate's anchor cannot match. Read offsets only after the tag matches, preserving exact match-bucket ordering, candidate depth two, and all score/tie rules. For larger sources with no spare bits, the original representation and verification remain. No archive field or decoder behavior changes.

Unit test `TestMatchFingerprintTagExactness`, written as a **disposable regression oracle**, compared tagged and original untagged indexes over deterministic random source arrays with inserted exact and XOR match regions, with/without transformed tables, repeated updates, and realistic recent distances; passed normally and under targeted Go race testing. Unchanged full Go suite and vet passed, as did Windows/amd64 cross-build. Untimed production-source-vs-candidate identity comparison verified **all twelve archive SHA-256s identical**. **Ten candidate-only complete original Silesia runs**, each re-encoding and decoding all 12 complete input files and asserting `bytes.Equal`, yielded 120/120 identical archive hashes and 120/120 byte-exact restores:

| Candidate complete run | Archive bytes | Encode seconds | Decode seconds |
|---:|---:|---:|---:|
| 1 | 62,741,534 | 6.100073247 | 2.564636830 |
| 2 | 62,741,534 | 5.869207911 | 2.378964363 |
| 3 | 62,741,534 | 6.014424369 | 2.494595442 |
| 4 | 62,741,534 | 6.063916020 | 2.488789231 |
| 5 | 62,741,534 | 6.124333921 | 2.442955856 |
| 6 | 62,741,534 | 5.925970327 | 2.450337568 |
| 7 | 62,741,534 | 5.670412653 | 2.435893563 |
| 8 | 62,741,534 | 5.944081650 | 2.468677163 |
| 9 | 62,741,534 | 5.915449581 | 2.584940118 |
| 10 | 62,741,534 | 5.840130450 | 2.471436991 |
| **Ten-run whole-corpus median** | **62,741,534** | **5.935025989** | **2.470057077** |

The first complete-corpus screen was 5.619059 s encoding and 2.371646 s decoding; this result did **not** persist as a substantial encoding-speed improvement in repeat tests, and must not be substituted for the median.

### Required comparison

| Compressor | Complete archive bytes | Reduction from original | Encode seconds | Decode seconds |
|---|---:|---:|---:|---:|
| Zstd-3 (historical external CLI) | 66,216,569 | 68.76% | 1.2561 | 0.8551 |
| XZ-6 (historical external CLI) | 49,233,340 | 76.77% | 116.9892 | 2.3079 |
| Current accepted JLZ R42 (**retained ten-run median**) | 62,741,534 | 70.396% | 5.955825126 | 2.406640990 |
| Candidate JLZ tagged matcher (**new ten-run median**) | 62,741,534 | 70.396% | 5.935025989 | 2.470057077 |

| Compressed-size delta candidate vs current | Encoding-time delta | Decoding-time delta |
|---:|---:|---:|
| **0.0000%**, exact archive hashes | **-0.3492%**, nominal, unpaired | **+2.6350%**, nominal, decoder unchanged |

Historical R42 results were obtained in another session, not a paired A/B; minor runtime differences should not be called causal. External Zstd/XZ CLI timings have different I/O/host conventions. No source/archive format changes and no input per-file size regression. **REJECT**: dynamic tag masks/branches/position bit packing are not justified by a nominal 0.35% encode-time improvement.

## Candidate B — Exact lookup table for adaptive probability quotients

Built a native static triangular `uint16` lookup table for all normalized frequency pairs with sums <= 2048, avoiding the two integer divisions in `entropy.bitModel.probability`. The table is approximately 4 MiB. For larger sums an exact integer fallback preserves the original value. Reverted tagged matcher first; this is an independent R42-based prototype. Full Go tests passed. One complete Silesia screen generated **62,741,534 B**, all 12 same archive SHA-256s and byte-exact reconstructions; **6.594074 s encode / 3.131522 s decode**. No sign of material speed gain on either side; extra lookup memory/initialization/caching cost is unfavorable. **REJECT** after this complete screening run, without a wasteful ten-run repetition.

## Final rollback, checks, and next priorities

All temporary Go production and diagnostic test files were removed from the trial source; restored hashes verified identical to the approved committed R42 source for the four canonical touched files (match source/oracle, regional preflight, fixture test). After rollback, `go test -count=1 ./...`, `go vet ./...`, and `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build` passed. The tagged candidate also passed a targeted `go test -race` exactness and existing transformed matching test. Complete Silesia inputs were verified with `sha256sum -c` against all twelve canonical digests. GitHub source branch `main` receives **documentation only**.

**REJECT H1.** The profile identifies substantial *genuine* selected-path work in match searching and arithmetic, but storing collision tags and replacing probability divisions with a large table do not yield substantive practical wins. Avoid further fingerprint/tag/mask/LUT variations. Next investigations should consider a structurally different coding or parsing architecture that reduces *necessary* per-byte work, or integrate a higher-ranked mathematical research transform without discarding current winners. Preserve byte-exact decoder compatibility and benchmark complete original Silesia candidates only.
