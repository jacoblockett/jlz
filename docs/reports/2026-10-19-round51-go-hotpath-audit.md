# Round 51 — Go implementation hot-path audit: bounded divisions and buffered arithmetic bits (ACCEPT H1)

**Date:** 2026-10-19. **Approved production:** Round 51 supersedes Round 44 for these strictly format-preserving internal speed changes. **Compression ratio and archive bytes are unchanged.** The observed small timing changes are directionally favorable, but the comparison is across independently measured sessions and is not proof of precise speedup.

## Question and source audit

User requested an actual Go/programming-practice review rather than another mathematical compression theory. Started from the supplied production `b41ce26e5432` Round44 source on `main`, read AGENTS and core engineering docs, and used the original canonical twelve-file Silesia corpus (211,938,580 B, SHA-256 12/12 PASS).

Go CPU profiling of a **complete original-Silesia encode plus the harness's first-run decode**, using Go1.23.2, identified flat CPU sample hotspots: `entropy.(*bitModel).probability` **3.03 s**, `engine.(*matchIndex).find` **2.45 s**, `entropy.DecodeArithmetic` **2.36 s**, `entropy.addCount` **1.90 s**, `bwt.Inverse` **1.57 s**. BWT transform `6.02 s` *cumulative* CPU and arithmetic `EncodeArithmetic` `5.87 s` cumulative; these multithreaded/profile totals are **NOT the elapsed time of the complete encoder**, nor causal speed contributions. Code survey included `compressFast`, regional preflight/worker scheduling, native Huffman, suffix-array BWT, mathematical hashing/scoring and arithmetic code. Existing Round37–43 history already rejected several tag/LUT/pool/fingerprint variants; avoid repeating them.

## Implemented exact Go changes

1. **Bounded 32-bit model math:** `entropy/arithmetic.go` counts rescale when total hits 2,048. Change `normalizedCount` and the two adaptive probability quotient calculations/weighted arithmetic from `uint64` to `uint32`; preserve the `uint64` probability API for the existing arithmetic range coder. Max multiplication from the real reachable count state is `2048 * 65536`, and max weighted probability sum is `8 * 65536`, both fit unsigned 32-bit exactly. No float math, probability approximation or wire change.
2. **Buffered bit reader:** replace repeated bit-index division/modulo and bit masking per decoded arithmetic bit with a current-byte field and one byte load every eight bits. Preserve original MSB-first bit order, reader advancement and implicit zeros after payload end, including malformed-input behavior. No new allocations per bit, no new codec mode, flags, worker changes or runtime architecture.

Independent temporary candidates `u32`, `reader`, `both` passed the targeted Go tests. One complete original Silesia screen per configuration, **not a median**, produced the exact approved R44 archive sum **61,418,211 B** in every case: `u32` encode **6.799607 s**, `reader` encode **5.682486 s**, `both` encode **5.152725 s**. These independently sampled timings vary substantially and must **NOT** be interpreted as reliable controlled speed percentages; retain only the fully validated ten-run combined candidate below.

Added `entropy/arithmetic_optimization_test.go` with a 64-bit historical probability reference applied during 100,000 deterministic model updates, and a byte-by-byte buffered-reader equivalence check including 65 implicit zero-padding bits after EOF. Both pass including targeted race test. The Go feature-package layout remains unchanged.

## Ten candidate-only complete original Silesia passes

The exact same compiled combined candidate was used for all runs, with ten entire 12-file archive encodes/decompresses, timers around individual full-file encode/decode calls, data fully preloaded, all 120 decoded outputs compared byte-exact. First four successfully completed within one process, whose wrapper later timed out; six were collected in three subsequent two-run invocations of the same test implementation. All **10** valid complete runs produced **61,418,211 B**, with **zero per-file size regressions**.

| Candidate run | Full archive bytes | Encode seconds | Decode seconds |
|---:|---:|---:|---:|
| 1 | 61,418,211 | 5.583079668 | 2.338516409 |
| 2 | 61,418,211 | 5.538198457 | 2.366661944 |
| 3 | 61,418,211 | 6.028442550 | 2.556091516 |
| 4 | 61,418,211 | 6.305629406 | 2.560179424 |
| 5 | 61,418,211 | 5.919284555 | 2.481814229 |
| 6 | 61,418,211 | 5.884995670 | 2.403558185 |
| 7 | 61,418,211 | 5.702330350 | 2.470547372 |
| 8 | 61,418,211 | 5.928908492 | 2.425475754 |
| 9 | 61,418,211 | 5.945874418 | 2.342864127 |
| 10 | 61,418,211 | 5.514730810 | 2.262020746 |
| **Median** | **61,418,211** | **5.902140113** | **2.414516970** |

As an *untimed correctness-only archive-identity audit*, encoded each of the twelve complete original Silesia files with the untouched approved R44 source and the combined optimized source, comparing each resulting archive's SHA-256. **12/12 SHA-256 archive hashes exactly identical**; this is not a repeat of R44 speed benchmarks. Original Silesia source SHA-256 12/12 matched canonical manifest separately. Archive compatibility holds **bit for bit** for this corpus; existing JLZ1 decoder semantics and accepted new-BWT mode15 remain unchanged.

### Reused reference comparison

| Model | Complete archive bytes | Reduction of 211,938,580 B original | Compress s | Decompress s |
|---|---:|---:|---:|---:|
| Zstd-3, historical external CLI | 66,216,569 | 68.76% | 1.2561 | 0.8551 |
| XZ-6, historical external CLI | 49,233,340 | 76.77% | 116.9892 | 2.3079 |
| Previously approved JLZ R44, historical 10-run median | 61,418,211 | 71.02075% | 5.945576 | 2.445270 |
| **Candidate JLZ R51, NEW 10-run median** | **61,418,211** | **71.02075%** | **5.902140** | **2.414517** |

| Compression-size delta R51 vs R44 | Encode-time delta R51 vs R44 | Decode-time delta R51 vs R44 |
|---:|---:|---:|
| **0.000% (exact archive SHA-256 equality)** | **-0.731% nominal** | **-1.258% nominal** |

**Timing caution:** R44 was independently timed in a prior session. CPU sharing, process/GC, OS effects and frequency scheduling can exceed this small delta. Do not claim a controlled 0.73% optimization or substantially new Zstd-class throughput. Zstd/XZ historical CLI conventions differ from JLZ in-memory measurements and are not directly comparable.

## Validation and decision

- `go test -count=1 ./...` PASS; `go vet ./...` PASS.
- Targeted `go test -race -count=1 ./entropy ./engine ...` PASS; new arithmetic-optimization tests also pass under `-race`.
- `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build` PASS (binary not natively run on Windows).
- Complete canonical Silesia 12/12 original inputs SHA-256 verified; 120/120 complete candidate reconstructions correct.
- Twelve complete new archives **SHA-256 identical** to untouched approved R44 encoder, zero size regression, no new decoder subtype.
- No external backend, no GitHub Actions, no accepted reference timing reruns, no new CLI options.

**ACCEPT H1 as a bounded implementation cleanup only.** Both edits are small, self-contained, mathematically exact and format-preserving, with directionally favorable but inconclusive sub-1.3% timing differences; they introduce no format/memory/model complexity and cannot worsen archive size. Round51 becomes the current production source and retained speed reference, with historical R44 measurements preserved.

**Next meaningful engineering targets:** profile specifically *encode-only* per selected file; evaluate whether any redundant full-region candidate work can be skipped using an exact lower-bound proof; examine heap traffic in BWT suffix-index serialization without `unsafe` or internal Go dependencies. Do not revive previous rejected probabilistic LUT, fingerprint or index pooling tuning.