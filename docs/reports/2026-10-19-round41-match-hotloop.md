# Round 41 — Exact-output mathematical matcher hot-loop trials (REJECT H1)

**Date:** 2026-10-19. **Decision: REJECT H1, no production source changes.**

## Hypothesis and why these experiments were selected

The previously captured *encoding-only* Silesia CPU profile identified `engine.(*matchIndex).find` as a high-cost selected-path matcher function (roughly 3.25 CPU seconds flat in that earlier profile, and `(*matchIndex).addBefore` at roughly 1.38 CPU seconds flat). Four selected Silesia archives rely predominantly on mathematical matching (mozilla, ooffice, osdb, samba), while six select BWT, one additive word, and one lane. Unlike repeated general preflight gating or BWT suffix sorting, these trials specifically change tight loops in the winning native mathematical match index. Neither experiment alters archive syntax, parser decisions, encoded tokens, or decoder. Implemented only in disposable local copies of the approved Round 39 source.

### A. Eight-byte transformed match verification

Replaced the byte-at-a-time comparison for established XOR and modular-additive reference candidates with an eight-byte check using `binary.LittleEndian.Uint64`. XOR uses a repeated-byte word mask; additive matching splits even and odd byte lanes before modulo-256 addition to prohibit cross-lane carry. Falls back to the identical scalar tail at the first nonmatching eight-byte unit. The comparator only sees matches already verified by the original fingerprint and anchor. It therefore preserves candidate type, length, repeat-offset scoring, command syntax, and archive bytes.

One complete original Silesia screen: **62,494,395 B**, encode **6.548066457 s**, decode **2.508143396 s**, all 12 exact roundtrips and all twelve archive SHA-256 hashes matched the retained unmodified source. Full uncached Go suite and `go vet ./...` passed; Windows/amd64 cross-build passed. This is **+5.39% encoding-time nominal versus the historical Round39 median**, not a paired comparison; no demonstrated benefit, so stopped without ten-run promotion. Do not assume the eight-byte verification is always faster: most attempted matches are short or mismatching, so extra word operations can dominate.

### B. Four-lane transformed-match fingerprint arithmetic

The native match hash's five-byte XOR fingerprint originally computed four per-byte XORs and shifted them into a `uint32`. Replaced with one unaligned little-endian `Uint32` load XOR a replicated first byte. For modular subtraction, used guarded even/odd 16-bit-lane subtraction and masking; the guard prevents borrows between output bytes, exactly preserving modulo-256 fingerprint values. No index-size, collision-handling, reference selection, or wire change.

One complete original Silesia screen: **62,494,395 B**, encode **6.266165928 s**, decode **2.421734942 s**, all 12 exact roundtrips and **12/12 archives byte-identical to the retained unmodified source** (verified against archived per-file SHA-256 logs). Full uncached Go suite and `go vet ./...` passed; Windows/amd64 cross-build passed. Nominal encoding-time **+0.85%** versus retained unpaired Round 39 median; effectively neutral. No ten-run promotion justified.

Both are fundamentally *screening-only* measurements, **one full-corpus screen per candidate**, not ten medians. One candidate cannot be called faster or slower with high confidence relative to historical independent sessions; no material gain was demonstrated. The run harness was the previously validated `TestRound39Profile` with complete 12 original Silesia inputs read before timing and each archive decoded and checked using the embedded digest and exact original bytes. All source archives total exactly 62,494,395 B; no source baseline or Zstd/XZ benchmarking was repeated.

## Required comparative table

| Model | Full compressed size | Reduction from original 211,938,580 B | Compression time | Decompression time |
|---|---:|---:|---:|---:|
| Zstd-3, **historical CLI** | 66,216,569 B | 68.76% | 1.2561 s | 0.8551 s |
| XZ-6, **historical CLI** | 49,233,340 B | 76.77% | 116.9892 s | 2.3079 s |
| **Current accepted JLZ Round 39, historical ten-run median** | **62,494,395 B** | **70.51%** | **6.213451 s** | **2.412670 s** |
| **Candidate JLZ fingerprint rewrite, one full-corpus screening pass** | **62,494,395 B** | **70.51%** | **6.266166 s** | **2.421735 s** |

| Archive-size change vs Round 39 | Compression-time change vs Round 39 | Decompression-time change vs Round 39 |
|---:|---:|---:|
| **0.00%** (every archive hash identical) | **+0.85% nominal** (unpaired, inconclusive) | **+0.38% nominal** (unpaired, unchanged decoder) |

Historical Zstd-3/XZ-6 measurements are external CLI timings from a different runtime methodology. Round39 JLZ is a ten-run median from another session, not a concurrently benchmarked control. Never claim a controlled slowdown from these nominal deltas.

## Decision and next research direction

**REJECT H1.** Both independently reversible and source-level exact match hot-loop optimizations reproduce every complete Silesia archive byte but do not demonstrate worthwhile compression-speed benefits. No ten-run promotion was attempted; neither performance screen justified that cost. All trial source code left only in disposable local work folders, no Go code promoted to `main`. Current approved **Round39** remains unchanged. No unsafe or backward-format behavior introduced.

The useful architectural finding is that vectorizing short/mostly failing reference comparisons does not establish a meaningful speed improvement. Stop pursuing nearby bitwise microtuning. Next project should profile actual selected mathematical-region token stream *work distributions*, especially index rebuilds and cross-region history traffic, before attempting a bounded architectural matcher change. For transformative end-to-end gains, reconsider avoiding redundant whole-method computation rather than assuming individual arithmetic/bit operations are the speed limit.
