# Round 52 — Zstd level-3 source architecture audit (RESEARCH ONLY)

**Date:** 2026-10-19. **Production baseline: Round 51, unchanged.** This is a source-level comparison and research prioritization, **not** a JLZ candidate implementation, integrated benchmark, or proof of predicted gains. No codec source was modified; the existing Silesia full-corpus baseline measurements are retained and not rerun. No Zstd implementation or source was copied into JLZ.

## Source provenance and actual level-3 strategy

Examined upstream Meta/Facebook zstd `dev` primary sources and format specification (links below). Zstd compression levels are parameter presets and the preset depends on source size/dictionary settings. For the **>256 KiB** default row, level 3 uses **ZSTD_dfast** (double-fast), parameters `windowLog=21, chainLog=16, hashLog=17, searchLog=1, minMatch=5, targetLength=0`. Thus the normal history window cap is about **2 MiB**, hash tables about **2^16 short and 2^17 long slots**. Do not confuse this with the more expensive greedy/lazy/binary-tree/optimal parsers at higher levels. Under short-input rows the exact parameter values differ. Source: https://github.com/facebook/zstd/blob/dev/lib/compress/clevels.h and https://github.com/facebook/zstd/blob/dev/lib/zstd.h.

The `zstd_double_fast.c` no-dictionary block parser maintains **two simple recency hash tables at different anchor lengths** (short configured minimum and 8-byte long), probes recent-offset references directly, scans/advances more aggressively after unsuccessful searches, and tests a long candidate around the following position when a short match is found. It uses backward match extension, opportunistic hash insertion adjacent to matches and quick repeated-offset sequences. These are conventional hash-based LZ techniques, not unique algorithms. Source: https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_double_fast.c.

Once a block's LZ parse is generated, the block format contains a **literal section** and a **sequence section** with literal lengths, match lengths and offsets. Literals use raw/RLE/Huffman with a stored or **reused** table; larger literals may be separated into four Huffman bitstreams. Three sequence symbol alphabets use finite-state entropy coding (FSE/tANS), falling back to predefined/default distributions, RLE, or prior-table reuse as selected; escape/extra bits carry exact parameters. The **level-3 fast-strategy branch uses cheap heuristics** when selecting predefined, repeat or freshly trained entropy tables rather than always paying for exact trial encodes. A block may remain uncompressed when net gain is too small. These choices preserve byte-exactness and make modeling overhead predictable. Sources:
- https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_compress_sequences.c
- https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_compress_literals.c
- https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_compress_internal.h
- https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md
- https://www.rfc-editor.org/rfc/rfc8878.html

## JLZ source-specific gap audit

Directly inspected current native Go snapshot, especially `engine/compress.go`, `engine/match.go`, `engine/region_fast.go`, `engine/cross.go`, `engine/mixed.go`, `engine/stream.go`, `engine/slots.go`, `engine/packed.go`, `engine/huffman.go`, and `entropy/arithmetic.go`, with prior reports R35–51.

| Mechanism | Zstd-3 source | Current JLZ / historical coverage | Assessment |
|---|---|---|---|
| Two hashes for **short and long exact references** | Dual separate position tables, fast direct match | JLZ match index holds two latest candidates per hash across up to three *transform types*, probes repeat distances, and often performs one-byte lazy evaluation. Some exact-only/stride scaffolding exists (`tokenizeMatchV2RangeStride`) but not a demonstrated **low-cost two-anchor exact-only independent match encoder on actual winning streams** | **OPEN** as *structural new fast-path experiment*, not bit/fingerprint optimization |
| Adaptive long no-match stepping | Search step grows on misses, quick short/long candidate checks | JLZ matcher normally scans every byte in an index-maintenance pass; stride mode exists but R35 speed experiments/other skip screens were inconclusive | **PARTLY EXPLORED**; do not retry simple stride-count tuning. Reassess only coupled to a distinct verified fast parser |
| Cross-block history | 2 MiB default source-history window; recent offsets across blocks | JLZ has exact self-contained 256 KiB regions, four-region independent decode groups, 128 KiB recent cross-block history with `crossPotential`, and repeat-distance cache | **NOT NEW.** Larger persisted windows might improve ratio but trade speed/parallel decoding/memory; lower priority |
| Separate literal/length/distance streams, log slots | Sequence symbol codes and raw extra offset bits | JLZ already splits op/length/distance/transform/literals; distance log-width slots+verbatim extra bits; some length splits; recurrent distance cache | **ALREADY IMPLEMENTED** |
| **Low-cost static/predefined/repeat FSE sequence tables** | FSE/tANS on LL/ML/OF with predefined/repeat/dynamic, selected heuristically at level3 | JLZ uses raw/RLE/canonical Huffman for metadata; R48 tried **role-conditioned per-block rANS** with slower complete integration, *not* a low-context predefined/reused FSE system on three small selected symbol alphabets | **HIGH PRIORITY OPEN.** Small native independently implemented normalized-state entropy on already-selected metadata, framed and decoded exactly |
| **Reuse one entropy model without concatenating data** | Prior literal Huffman tables and sequence FSE tables can be reused across blocks with minimal framing | R50 tested concatenating shared Huffman streams, which was a loss; that is **not** the same as coding independent streams using previously transmitted tables, selectively when histograms indicate net savings | **MEDIUM PRIORITY OPEN.** Group-local, explicit deterministic table lifecycle preserving independent group decodes |
| 1-vs-4 literal Huffman streams | Four streams for larger literal sections, potentially better decode throughput | JLZ currently uses a single 15-bit-prefix canonical Huffman bitstream with full lookup decoder and whole-group parallelism | **LOWER PRIORITY OPEN.** Test only if literal-Huffman decode CPU is material and framing/jump offsets are priced |
| Preflight and incompressible fallback | Fast small-block thresholds and default-table decisions | JLZ already has sampled math/BWT gates, direct raw/RLE/Huffman comparison, Shannon bounds and previous rejected abort/cancellation screens | **EXHAUSTED nearby variants**. Avoid revisiting familiar skip tuning |
| End-to-end single parse | One fast match parse feeds sequence+literal coding within a block | JLZ's strong automatic choices (BWT, additive, mathematical matching, predictive) can do real competing work, but preflight and sample reuse are already implemented | **Architectural trade-off**, not evidence for replacing a superior native path. Wholesale rewrite requires transformative complete-corpus evidence |

## Ranked experiments for future authorization

1. **Small-alphabet FSE/predefined-table supplement on currently selected mathematical metadata**, not original bytes, with code-length histograms, exact table payload, raw extra bits and no role-state map. Gate on *full serialized* gain and estimate CPU before a full-integrated trial. Explicitly contrast with prior rejected static-role rANS; don't treat all ANS implementations as equivalent. Native original implementation only.
2. **Conditional repeat-table reuse** inside JLZ's *existing independently decodable groups* for selected op/length/BWT metadata. Charge group model header and individual table-reuse flags; retain existing Huffman when cheaper. Distinct from R50 rejected concatenation and R47 learned maps. Under independence requirements, reset at each group boundary.
3. **Simpler dedicated exact-match fast candidate using two distinct anchor lengths**, with cheap match/no-match stride and conventional recent offset, as an *alternative* to transformed math matching where cost evidence warrants. Validate if actual selected math streams permit adequate ratio. Do not remove XOR/ADD production modes, rehash existing bucket depths or blindly increase search budget.
4. **Four-lane Huffman decode for sufficiently large selected native symbol streams**, only after CPU attribution. Prioritize compression-speed wins over a complex standalone decoder optimization.

In descending priority, #1 and #2 have the clearest *untried* complementary mechanics; #3 addresses the speed gap more directly but may give up the math-selected ratio advantage; #4 is speculative. No claim that adopting any Zstd-inspired idea will outperform Zstd or current JLZ. No Zstd code or compressor dependency is authorized.

## Benchmark/provenance (references only; no Round52 candidate)

| Model | Complete Silesia archived bytes | Reduction vs 211,938,580 original bytes | Encode s | Decode s |
|---|---:|---:|---:|---:|
| Zstd-3 (historical external CLI) | 66,216,569 | 68.76% | 1.2561 | 0.8551 |
| XZ-6 (historical external CLI) | 49,233,340 | 76.77% | 116.9892 | 2.3079 |
| Approved JLZ R51 (retained candidate-only ten-run median) | 61,418,211 | 71.02075% | 5.902140 | 2.414517 |
| R52 JLZ candidate | N/A: source audit only | N/A | N/A | N/A |

| Size delta | Encode-time delta | Decode-time delta |
|---|---|---|
| N/A | N/A | N/A |

The external CLI times and in-memory JLZ times have different measurement conventions, so ratios of these wall times do not establish algorithm-level causality.

**Conclusion: ACCEPT H1 as a *research queue*, not a production codec decision.** Zstd level3 has at least two distinct and not-yet-closed stream entropy mechanisms worthy of original native JLZ prototypes. Current approved production remains **Round51**; no source edit, archive, correctness run or new benchmark was performed for this source-only audit.