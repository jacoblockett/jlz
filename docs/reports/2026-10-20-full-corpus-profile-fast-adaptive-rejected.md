# Full-corpus CPU diagnosis and fast native arithmetic model: layered integrated candidate rejected

**2026-10-20. Status: REJECT candidate H1; engine unchanged.** Research uses the supplied **complete original 12-file Silesia corpus** (211,938,580 B). No mature compression backend, no numbered round, no GitHub Actions. All experiments used disposable local copies of the verified commit `684b5e3` with its exact checked-in root `default.pgo`.

## Existing architecture and source bottleneck profile

Baseline source aligned with GitHub `main` engine tree and accepted full Silesia archive **60,244,257 B**. All 12 original files passed the canonical 12 SHA-256 checks (verified on the supplied corpus previously), and profiling run reproduced precisely the accepted 12 archive byte lengths and exact complete decompression of each file.

A single separately instrumented full-corpus Go CPU profile (`runtime/pprof`, not a new timed benchmark) recorded **22.03s sampled CPU** during about 6.74s profiled encoding wall time and **8.41s sampled CPU** during 2.62s profiled decoding wall time. Totals exceed elapsed time because JLZ compresses independent groups concurrently. Main cumulative hotspots:

| Current encoder function | Cumulative samples | Share of total sampled CPU |
|---|---:|---:|
| `bwt.Transform` (Go suffixarray + serialization + parsing) | 7.17s | 32.55% |
| `tokenizeMatchV2RangeStride` | 5.53s | 25.10% |
| `entropy.EncodeArithmetic` | 4.86s | 22.06% |

| Current decoder function | Cumulative samples | Share |
|---|---:|---:|
| `entropy.DecodeArithmetic` | 4.00s | 47.56% |
| `bwt.Inverse` | 1.96s | 23.31% |
| `bwt.BijectiveZeroDecode` | 0.83s | 9.87% |
| `finiteStateDecode` | 0.43s | 5.11% |

Cumulative profiles can overlap; never sum these percentages or subtract them directly from wall times. This identifies bottlenecks in two essentially different subpipelines: mathematically referenced data uses costly adaptive bit coding and match generation, while other data uses suffix sorting and LF inverse.

Compared source-level design to upstream **Zstd** block and sequence design, **libdeflate** table-based fast Huffman decode, **liblzma/XZ** stateful coding/parsers and BCJ, and **Kanzi** ROLZ/word transformations. Distinguish their available architectural ideas from JLZ's already implemented match tokens, history, separate entropy columns, rANS, backwards-compatible block formats and executable normalization. Libdeflate-style Huffman lookup cannot plausibly solve the dominant 47.6% arithmetic decoder hotspot. This is a **JLZ code-first profile**, not a claim to have benchmarked other compressors.

## H1: low-work decoder-synchronous probability model supplementing current math literal packing

Previous concrete negatives: exact reciprocal and triangular probability LUT attempts tried to preserve the old frequency quotient and generally worsened speed; context-conditioned BWT rANS improved some compressed symbols but slowed decoding; prior LZP and bitplane placements on math literals expanded size. This experiment is structurally different: it uses **two bounded exponentially moving binary probabilities per prefix**, one global and one indexed by previously decoded literal byte, rather than two divisions of current normalized frequency counters. Both encoder and decoder update the same 256 prefix models and 65,536 (previous byte, prefix) models, with fixed integer updates, 16-bit clamped probabilities, and a tiny saturating context age determining the blend. It layers into `packLiteralStream` using the **existing original 32-bit binary arithmetic interval encoder/decoder**, retaining raw/RLE/Huffman alternatives and introducing a new packed-stream **codec 6**. Old codec 3 remains decodable. No external backend.

Concrete prototype files (kept out of production): `entropy/arithmetic_fast.go`, `engine/stream.go`, targeted exact roundtrip tests. New literal entropy stream is evaluated inside actual JLZ math block mode selection; branch-normalized BWT, mathematical exact/XOR/ADD references and the two-command automatic interface remain unchanged. No compressor-side competing double-encode of old and new arithmetic models is performed.

Nine distinct first/middle/last source-block math literal component screens (original mozilla, samba, ooffice) restored all literal bytes exactly, and native new coder decreased the nine arithmetic body sizes **9,475 B** versus old order-1 coder: mozilla 207,063→204,527; samba 155,718→152,442; ooffice 223,329→219,666. Per-component timings favored the new estimator, **but this is not an integrated full-corpus performance claim**.

## Complete integrated JLZ candidate, new ten-run measurements

The candidate **fully encoded and decoded all twelve original Silesia files on each of ten completed runs**, rechecking byte equality and SHA-256 on every restoration: **120/120 PASS**. Because the execution window interrupted a fifth run in an original batch, only four completed runs were retained, followed by two separate 3-run batches. No incomplete run is counted. The size was identically **59,909,752 B** across every completed run, and no original file's compressed length increased. Differences vs current: mozilla **18,609,139→18,320,567 (-288,572 B)**; ooffice **2,692,252→2,674,623 (-17,629 B)**; samba **5,134,262→5,105,958 (-28,304 B)**; other nine unchanged.

| Completed candidate run | Archive bytes | Encode s | Decode s |
|---|---:|---:|---:|
| 1 | 59,909,752 | 6.055532592 | 2.169659951 |
| 2 | 59,909,752 | 5.932058158 | 2.159490193 |
| 3 | 59,909,752 | 5.889681136 | 2.181327101 |
| 4 | 59,909,752 | 6.132568554 | 2.294137567 |
| 5 | 59,909,752 | 6.268531781 | 2.230045919 |
| 6 | 59,909,752 | 6.243065910 | 2.224250411 |
| 7 | 59,909,752 | 6.097533884 | 2.237408085 |
| 8 | 59,909,752 | 6.275214902 | 2.271781448 |
| 9 | 59,909,752 | 6.224181548 | 2.279779973 |
| 10 | 59,909,752 | 5.721230736 | 2.126430047 |
| **Independent medians** | **59,909,752** | **6.115051219** | **2.227148165** |

Historical accepted current, not rerun as a comparison benchmark: **60,244,257 B / 5.322679488 s / 2.003659634 s**. Relative candidate/current **-334,505 archive bytes (-0.555248%), +14.886708% encode wall, +11.154017% decode wall**. Because sessions are unpaired, do not infer causal exact percentage of overhead solely from wall times. Performance regressions are much larger than the nominal noise margin.

A separate candidate full-corpus CPU profile showed encoder arithmetic sampled CPU fell from 4.86s to about 4.31s and decoder arithmetic from 4.00s to 3.51s, but overall CPU and scheduling still dominate and **end-to-end timing regressed**. This is not a replacement for ten-run numbers.

Targeted new codec roundtrips and full-corpus exact reconstruction passed. **Unmodified existing tests asserting codec 3 specifically failed as expected when candidate emitted codec 6**. An uncached all-package test PASS must not be claimed for the rejected candidate. No new Windows cross-build, race, or backwards-existing-archive test was required for promotion because the candidate fails the strict wall-time gate decisively. The historical decoder code was explicitly left present but decoding a pre-existing binary archive with the trial was not independently performed.

## Additional distinct float/numeric layer feasibility screen

An independent reversible 8-byte-lag XOR *layer applied before JLZ's existing 8-lane native byte entropy representation*, exact inverse tested across entire individual original Silesia source blocks: sao 28 blocks **6,317,645→6,980,911 (+663,266)**, x-ray 33 **5,893,903→4,898,897 (-995,006)**, osdb 39 **8,369,594→9,317,179 (+947,585)**, mozilla 196 **33,438,716→38,553,087 (+5,114,371)**. Although it improves the particular x-ray lane component, actual accepted x-ray JLZ math/additive archive is only **4,266,987 B**, so the new layer is **larger than JLZ's selected representation** and adds encode/decode passes. This does not justify integration or a ten-run codec candidate. Different from 16-bit spatial geometry and fixed 3-lag math-literal GF2 screens, but the observed ratio still loses to actual current selection.

## Conclusion

**REJECT H1 for integration** despite real size savings: strict required no-regression for size **and both time medians** is not met. Do not commit the new packed codec or leave unaccepted experimental code on `main`. Do not start a shift-rate/probability-age/window tuning spiral or pretend the specific estimator disproves all arithmetic-model alternatives.

The required architectural leap is bigger: (1) significantly reduce the work of the **existing suffix-array-based BWT forward step without a worse replacement**, (2) preserve high-order adaptive modeling at a fraction of **bit-serial arithmetic decode work**, (3) avoid unnecessary expensive mathematical match discovery without sacrificing matched references or selected ratios. Study fast bytewise range/rANS state propagation and direct induced suffix ordering as separate, source-grounded mechanisms only if a plausible measurable decoder/work complexity change is demonstrated. Avoid repeating historical slower prefix doubling, packed LF inverse, reciprocal arithmetic, opcode-distance masks, matcher hash-depth sweeps or duplicate entropy encodings.

99% reduction of this whole Silesia corpus would mean **<=2,119,385 B** of total archive, which is about 28.3x smaller than accepted JLZ. Treat as a research aspiration, not an achieved or empirically implied guarantee. Achieving <1 s each direction on this complete corpus would require >211.9 MB/s effective throughput per direction, well beyond accepted medians on the current host.


## Source archive bit-budget and finite-state lane-placement screen

One diagnostic original full-corpus native JLZ encode (not a new baseline timing benchmark) parsed exact final block modes and payload byte sizes. Major outputs:
- **mozilla** 18,609,139 B (30.89% of total archive). Its largest current payloads: cross-history mathematical packed mode9 **11,018,153 B** (126 blocks), split-length packed mode10 **6,406,815 B** (58 blocks), predictive packed mode11 **389,495 B** (6 blocks). Remaining small original/predictive/match streams fill its size.
- **webster** 9,835,416 B (16.33%), BWT mode15 **9,592,323 B** and segmented mode16 **241,940 B**; 159 blocks.
- **sao** 6,261,713 B (10.39%), **28/28 blocks native lane mode13**, **6,261,476 B** lane payloads.
- **samba** 5,134,262 B (8.52%), mostly math modes9/10 **4,910,975 B**.
- **x-ray** 4,266,987 B (7.08%), all 33 blocks mode14 (additive).
These five account for over 73% of archive bytes. The three largest alone account for **57.6%**, a better allocation of architecture work than optimizing codebook metadata.

**Supplementary finite-state-on-existing-lane component (NEW-SCREEN; CLOSED-PLACEMENT):** JLZ's `lane.EncodeBlock` currently calls `packStream` (raw/RLE/Huffman) on each selected byte lane, although the decoder already supports native finite-state stream codec5. A disposable native Go trial passed `packFiniteState` instead as that same lane payload writer. Existing decoder reconstructed each original block exactly. The *entire original* 256KiB block set for sao (28 blocks) changed **6,317,645→6,293,686 B (-23,959 B)** at component encoding time **0.061→0.146s**; 28/28 component winners. x-ray 33 blocks changed **5,893,903→5,870,322 B (-23,581)**, and osdb 39 blocks **8,369,594→8,338,470 B (-31,124)**. Those are *candidate lane components*, not newly selected complete JLZ archives, especially x-ray and osdb whose present automatic routes already beat their lane counterparts. The extra FSE normalization/encoding adds about 85ms on sao alone for 24KB saved, before accounting for full selection work; and new per-symbol decode state work cannot be assumed faster. **REJECT concrete supplemental lane-FSE placement**, no new codec mode or on-main code.

**Source-grounded next direction:** Zstd's specification documents independent one/four-way Huffman literal bitstreams and FSE for sequence columns; Apple LZFSE uses native table-driven finite-state literal/match/offset decoding; libsais implements highly optimized induced sorting and memory/prefetch locality for BWT. JLZ already has FSE, parallel independent regions and SA-IS through Go's suffixarray. A serious new prototype must change the *dominant* cost: for example direct fast native induced BWT output without Go index serialization while retaining equal ordering, or multi-symbol context-conditioned entropy models that avoid eight serial arithmetic bit updates per literal. Neither external project should be used as a JLZ compression backend. Code and corpus efficiency are first-stage screens, not promises of accepted results.
