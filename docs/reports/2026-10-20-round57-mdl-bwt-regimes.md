# Round 57: MDL selection of BWT entropy neighborhoods (ACCEPT)

**Date:** 2026-10-19. **Final decision: ACCEPT H1** as a bounded, original native enhancement to accepted Round55. Full performance corpus: all twelve complete original Silesia files, **211,938,580 source B**, canonical source SHA-256 **12/12 PASS**. No external compression backend, imported FSE/Zstd/Brotli code, new CLI switch/level or compression-ratio-for-speed trade.

## Research hypothesis and actual implementation

Round56 demonstrated that several native finite-state models over local neighborhoods of *already selected BWT bijective MTF symbols* can improve compression, but evaluating two/four/eight actual encoded alternatives took disproportionate CPU. Round57's hypothesis was to predict the winning regime count *before* actually encoding the alternatives.

Added **one native experimental/internal BWT block subtype 16** with 2 or 4 contiguous independently framed native entropy streams. The decoder validates partition count 2/4 and exact stream bounds, concatenates decoded bijective symbol streams, then uses the unchanged bijective MTF inverse and BWT inverse. Legacy subtype 15 remains readable. New archives containing subtype 16 require new decoders.

The single **minimum-description-length (MDL) selector** measures the zero-order empirical entropy of the whole bijective stream and candidate equal-length 2-/4-part histograms, charging each part a conservative cost of **3 B per distinct symbol + 8 B stream overhead** and 1 B partition count, with an additional 128 B prediction margin for model selection. The encoder:
1. Transforms BWT **once**, builds bijective MTF symbol stream **once**, and encodes incumbent mode15 **once** (as existing Round55).
2. Computes inexpensive symbol-frequency estimates for 1, 2, and 4 neighborhoods without performing additional full stream encodes.
3. Encodes **at most one** predicted 2/4 candidate, with early size abort where possible. The incumbent mode15 still wins all ties and nonprofitable actual packed cases.
4. Reuses the exact unchanged Round55 BWT preference preflight and its sampled mode15 blocks. Mathematical, additive, lane and entropy paths are otherwise untouched. No per-file hardcoded decisions.

The state models are simply native JLZ finite-state frequency models already implemented under Round53, not a third-party compressor. The new code is in engine/bwt_subregions.go with integration in engine/bwt.go and engine/cross.go. Tests explicitly cover both 2/4 decoder modes, byte-exact original reconstruction, truncated payloads and bad partition headers.

## Full original Silesia candidate-only performance

A first **ONE complete 12-file screen** produced **60,585,717 B / 5.392309 s encode / 2.176282 s decode**, all original file bytes exact. Screening time is *not* a promotional ten-run median.

Ten complete original Silesia candidate-only passes were run in a precompiled test binary; two complete runs per independent process. The first long combined wrapper timed out after **four** completed runs and during a fifth incomplete run. That incomplete attempt is **excluded**. Three more independently completed two-run invocations delivered the remaining six. All ten recorded complete runs encoded and decoded **all twelve files** and verified byte-for-byte, **120/120 exact restorations**, every pass emitting the same aggregate **60,585,717 B**.

| Complete candidate run | Archived bytes | Encode seconds | Decode seconds |
|---:|---:|---:|---:|
| 1 | 60,585,717 | 6.137815765 | 2.350623412 |
| 2 | 60,585,717 | 5.713484547 | 2.261364283 |
| 3 | 60,585,717 | 5.416871778 | 2.194108620 |
| 4 | 60,585,717 | 5.360704523 | 2.151197866 |
| 5 | 60,585,717 | 5.405903348 | 2.072980183 |
| 6 | 60,585,717 | 5.455282659 | 2.147712242 |
| 7 | 60,585,717 | 5.481199276 | 2.194747165 |
| 8 | 60,585,717 | 5.704999682 | 2.251413036 |
| 9 | 60,585,717 | 5.506288193 | 2.180609457 |
| 10 | 60,585,717 | 5.625293446 | 2.273740691 |
| **Median** | **60,585,717** | **5.493743735** | **2.194427893** |

Relative to retained historical approved Round55 **60,726,433 B / 5.500265326 s encode / 2.198443653 s decode**, Round57 provides **-140,716 B (-0.231721% archive bytes)**, nominal **-0.11857% encode**, **-0.18266% decode**. Those tiny timing differences are **unpaired between independent sessions and within measurement noise**. Encode/decode speed should be considered *neutral*, **not** a proven speed improvement. This nevertheless recovers **140,716 B** of Round56's predicted neighborhood compression gains with negligible measured CPU overhead.

### Required comparison, retained external references

| Model | Complete Silesia archive bytes | Reduction of 211,938,580 original B | Encode s | Decode s |
|---|---:|---:|---:|---:|
| Zstd-3, historical external CLI | 66,216,569 | 68.76% | 1.2561 | 0.8551 |
| XZ-6, historical external CLI | 49,233,340 | 76.77% | 116.9892 | 2.3079 |
| **Current approved JLZ R55, retained 10-run median** | **60,726,433** | **71.34715%** | **5.500265** | **2.198444** |
| **New candidate JLZ R57, candidate-only 10-run median** | **60,585,717** | **71.41355%** | **5.493744** | **2.194428** |

| Compressed-size delta R57 vs R55 | Encode-time delta | Decode-time delta |
|---:|---:|---:|
| **-0.231721%** (-140,716 B) | **-0.11857% nominal, neutral** | **-0.18266% nominal, neutral** |

Historical external Zstd/XZ CLI runtimes and JLZ in-memory runtimes use different measurement conventions, and they must not be interpreted as same-host algorithmic speed ratios.

### Individual complete file archives

| File | R55 B (historical) | R57 B (new) | Delta B |
|---|---:|---:|---:|
| dickens | 3,193,294 | 3,155,856 | -37,438 |
| mozilla | 18,609,139 | 18,609,139 | 0 |
| mr | 2,825,153 | 2,729,261 | -95,892 |
| nci | 2,163,199 | 2,163,199 | 0 |
| ooffice | 3,033,712 | 3,033,712 | 0 |
| osdb | 3,425,086 | 3,425,086 | 0 |
| reymont | 1,456,118 | 1,455,234 | -884 |
| samba | 5,134,262 | 5,134,262 | 0 |
| sao | 6,261,713 | 6,261,713 | 0 |
| webster | 9,837,297 | 9,835,416 | -1,881 |
| xml | 520,473 | 515,852 | -4,621 |
| x-ray | 4,266,987 | 4,266,987 | 0 |
| **TOTAL** | **60,726,433** | **60,585,717** | **-140,716** |

**Five files shrink, seven unchanged, zero per-file compressed-size regressions.**

## Separate genuinely different theory screen: incremental-frequency rank coding (REJECT tested concrete variant)

Independently prototyped a native, exactly reversible **incremental-frequency-count symbol ordering** *instead of MTF* after BWT. At each symbol, rank is determined by a shared state ordering sorted by cumulative occurrence frequency (stable on ties). Encoder transmits bijective zero-run coded rank digits; decoder reconstructs rank and updates its identical deterministic frequency/order state. Each candidate was fully packed through the same JLZ native finite-state encoder, and every sampled original block and BWT transformed byte sequence was restored exactly.

**21 source blocks**, at first/middle/last positions of each of the seven currently BWT-selected complete Silesia files, 3 per file, all 256 KiB except naturally shorter file-end blocks. Original post-BWT MTF/bijective native-packed baseline **966,482 B**; IFC replacement **2,280,951 B including one selector byte per block**, **+1,314,469 B**, **0/21 winners**. It loses materially on every sample. This is an **isolated 21-block screen**, not a complete original Silesia benchmark and not a blanket disproof of all IFC models. It is far too weak for full JLZ integration. The experiment's reversible proof-of-concept was kept out of production code.

## Validation and compatibility

- Full original canonical 12-file Silesia input SHA-256 verification **12/12 PASS**; full candidate source size 211,938,580 B.
- **120/120 candidate full-file exact byte restorations** across ten complete corpus runs, entire archive framing and native SHA-256 digests charged.
- Local source: uncached \`go test -count=1 ./...\` PASS, \`go vet ./...\` PASS, targeted \`go test -race -count=1\` on BWT subregion/BWT preference/native finite-state tests PASS, Windows/amd64 CGO-disabled cross-build PASS (Windows runtime not tested).
- Additional *post-benchmark* pure decoder unit tests explicitly constructed mode16 two-/four-part payloads and verified original byte-for-byte reconstruction, rejection of truncated payloads and invalid counts.
- Backward compatibility: previously built approved R55 CLI compressed original \`mr\` and new R57 CLI decoded that unchanged old archive byte-exact; new R57 standalone CLI compressed and decoded \`mr\` byte-exact. **Old R55 decoders do NOT support new mode16 payloads**, a required forward compatibility caveat. Existing old modes (including 12/15) and native finite-state subtype5 decoding remain untouched.
- Approved R55, Zstd and XZ reference speeds **not rerun**. The old CLI run above is a *correctness-only archive generation*, not a competing baseline benchmark.

## Decision and next theoretical queue

**ACCEPT H1, Round57.** A genuine native post-BWT local-model improvement survives actual complete-archive framing and runs at a ten-run median essentially indistinguishable from approved R55 while saving 140,716 B (0.232% of prior compressed archive) with no individual file regression. This justifies one additional bounded BWT internal subtype because it achieves a favorable actual size/runtime trade rather than the rejected Round56 multi-encode overhead. The strict non-negotiable rule remains: **never trade away compression ratio for speed**.

**Keep open distinct theories:** learned *variable* change-point boundaries rather than fixed 2/4 partitions (must charge boundary locations and avoid multiple full encodes); context-conditioned native entropy for actual match distance extra bits given already-decoded lengths; and other genuinely different symbolic state models. **Reject** the specific cumulative-frequency IFC ordering tested here on 21 complete original-source regions, without labeling all IFC theories impossible. Do not spiral into arbitrary threshold, partition-count or sampling tweaks without new structural mechanism.
