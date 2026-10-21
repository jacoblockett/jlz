# Zstd post-parse, within-column entropy partition screen

**Date:** 2026-10-21. **Status:** original-Go selected-stream COMPONENT SCREEN; no fully integrated candidate and no new format, production code, or PGO. Zstd-only research ongoing.

## Architectural source and novelty

Inspiration: [Meta Zstd `zstd_compress_superblock.c`](https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_compress_superblock.c) separates sequence discovery from entropy partition decisions. Previous JLZ native experiments either **split source before BWT** (451 blocks; coarse ideal -17,502 B) or **shared one rANS model across separate mathematical blocks** (optimistic -42,162 B) or tested **post-BWT symbol entropy segmentation** (already integrated historically). This is a different, still-untested placement: partition **the actual selected output columns of an already parsed mathematical block** into independently framed entropy chunks, without re-searching source matches or changing token selection. This particular trial divides each column by its *byte* positions, not by synchronously aligned command/sequence boundaries. Therefore it is **not** a complete Zstd sequence-store superblock experiment; semantically aligned interleaved partitions remain untested.

## Actual original-source method

Matched actual source to live `main` `4ea0ae3e824faa18412d4c1547888778262d3f55`; local `engine/packed.go` SHA1 Git blob `f25a6daf862e554c8f60b62a048eb2213e2866fe`. All twelve canonical original Silesia SHA-256s pass `sha256sum -c`; accepted corpus original 211,938,580 B. Ran native `compressFast` and `decompress` on **all 12 complete originals**; 58,767,739 B total, every source byte restored. Parsed all **279 actually selected mathematical blocks** from original accepted archives, using actual `unpackStream` for current emitted columns and actual emitted packed-column byte sizes as the baseline.

In one disposable native-Go test, split each eligible real output column at 1/2 or into equal quarters; each partition was encoded with existing original Go codec family: `packFiniteState` for opcode/length streams, `packVarintStream` for distance, `packStream` for masks, and `packLiteralStream` (unless the original native literal gate requires `packStream`) for literals. Charged **2 B extra per decomposed stream**, on top of full independent normal stream frames (including all models and renormalized states). Every partition was independently reversed using production `unpackStream`; concatenated bytes compared exactly against the real source column. The only oracle was choosing whichever *entire exact stream* changed favorably, not assuming every stream benefits. This is a component-stage selective optimistic envelope, not a new decoder format.

| Actually selected column | 2 equal parts, oracle-positive gain B | 4 equal parts, oracle-positive gain B |
|---|---:|---:|
| Opcode nibble stream | 13,738 | 23,346 |
| Joint command lengths | 15,352 | 27,203 |
| Separate literal lengths | 2,816 | 3,660 |
| Separate match lengths | 8,136 | 13,926 |
| Distance varints | 5,260 | 4,440 |
| Transform masks | 0 | 0 |
| Matched literal bytes | 40,848 | 66,674 |
| **Upper selected-field total** | **86,150** | **139,249** |

Second original-Go screen tested **two-part quarter, half and three-quarter cut positions**, plus four equal parts, on every eligible *metadata* field; kept the previous first trial's literal numbers instead of repeatedly running the more expensive adaptive arithmetic codec. Metadata field win ceilings by choosing the most favorable of those four layouts for each stream were opcode **26,015 B**, joint lengths **30,496 B**, separate literal lengths **5,085 B**, separate match lengths **15,500 B**, distances **10,683 B**, masks **0 B**: **87,779 B** total on metadata. Adding only the previous fixed-quarter literal upper estimate (66,674 B) yields a **154,453 B / 0.263% optimistic upper envelope**, before charging required partition-selection mode bits, metadata, decoder complexity, or extra encoding passes. The 4-way metadata-only fixed test gave less than 87,779 B, and direct full-size differences would remain far smaller than source content already efficiently coded.

The first adaptive experiment attempted extra repeated full arithmetic literal trials and **timed out**; it does not count as a successful completed screen. Rerunning adaptive boundaries with literals explicitly excluded completed in **12.50s**. The initial equal-half/quarter screen completed in **29.09s**. Each eligible partition inverse was exact. No full-corpus 10-run candidate because this upper bound and additional coding/format work do not justify a competing archive mode.

## Scoped conclusion

**DEFER / DO NOT IMPLEMENT this independent column-chunk wire mode.** The 0.263% oracle is not a hard information-theoretic bound on smarter boundaries, per-command model context, or synchronized sequences. Still OPEN: true Zstd-style command-aligned partition and shared sequence costs, or a cheap *no-duplicate-pass* predictor with demonstrated improvement large enough to buy new decoder complexity. Do **not** re-label two/quarter/four raw column boundaries as untested inventions, do not mistake a component upper bound for a fully encoded result, and do not modify the accepted encoder/decoder to chase this small modeled gain. Historical accepted JLZ: **58,767,739 B /4.7873709325s encode /1.8721679705s decode**; external historical Zstd-3 **66,216,569 B /1.2561 /0.8551s**, XZ-6 **49,233,340 B /116.9892 /2.3079s**, not rerun.
