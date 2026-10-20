# Rounds 62–63: native sequence-reuse and far-history experiments (REJECT H1)

**Date:** 2026-10-19. **Production unchanged: approved R57.** All findings use original full 12-file Silesia only, 211,938,580 B, all 12 canonical SHA-256 verified. No external compressor backend or copied external compression code. No experimental Go source promoted.

## External architectural investigation

Apple LZFSE's match-sequence architecture encodes separate literal-length, match-length and distance columns with finite-state entropy. JLZ already has similar native op/length/distance/transform/literal streams, so copying its separation would simply repeat existing work. LZFSE-inspired opportunity: decode those streams *directly*, avoiding reparsing a legacy synthetic record. LZ4's low-work sequence decoding also reinforces removing intermediate buffers, but our R54 exact-only fast parser was already rejected for harming ratio. Reviewed public primary architecture: https://github.com/lzfse/lzfse/blob/master/src/lzfse_internal.h ; https://github.com/lz4/lz4/blob/dev/doc/lz4_Block_format.md ; https://developer.apple.com/documentation/compression/compression_lzfse .

## R62: exact-format-preserving encoder and decoder optimizations

Implemented four additive changes in a disposable Go copy of approved R57:
1. Reuse the native packed candidate's exact length when the early mathematical-region preflight already encoded it; formerly the same 256 KiB sample could be coded twice merely to compare candidate lengths. If the mathematical huffman lower bound had legitimately skipped native coding, still encode the missing native size.
2. Replace integer divisions in the shared arithmetic probability model with bounded precomputed 32-bit reciprocals and one exact downward correction. Both arithmetic encoder and decoder retain **identical rational integer probabilities and bitstreams**. Counts are limited by the existing 2048-rescale invariant. Exhaustively checked every legal positive (zero, one) pair summing to at most 2048 against the former integer division; all pass.
3. Decode already restored native op/length/distance/transform/literal columns directly, instead of rebuilding, copying and reparsing a synthetic legacy command-format archive. Same verified match expansion primitive, unchanged format and stream checks.
4. Fuse inverse-BWT symbol and next-row LF link in a single 32-bit entry as in the previous R61 reversible prototype, only at the approved original 256 KiB source block size.

**All 12 original file compressed sizes and their reconstructed bytes are unchanged; aggregate 60,585,717 B identical to approved R57.** Full-screen candidate timings were inconsistent and generally worse than the retained approved R57 10-run medians (5.493744 s encode, 2.194428 s decode). One initial complete screen 6.247954/2.438441 s; after reciprocal coding another complete screen 5.897923/2.519706 s. Three complete arithmetic/packed-inverse screens 5.762034/2.300560, 5.333738/2.362747, 5.665494/2.544453 s. Three further complete screens after direct-sequence decoding 6.126883/2.614558, 5.386415/2.293410, 5.761652/2.507035 s. These are **candidate-only screens, NOT ten-run medians**, and do not establish an all-three-metric improvement. Stop rather than spending ten promotional passes on a demonstrated poor trade.

R62 local complete **go test -count=1 ./...**, **go vet ./...**, targeted **go test -race** on math/cross/BWT paths, Windows/amd64 CGO-disabled cross-build and all complete original-source roundtrips passed. Arithmetic exhaustive quotient equivalence PASS. No format or decoder compatibility change needed for R62.

## R63: new structural test inspired by longer-distance dictionary matchers

A distinct reversible native-Go test increased the group-local mathematical cross-region history from **128 KiB to 512 KiB**, maintaining the approved 256 KiB input regions and unchanged BWT-selected blocks. Decoder history was expanded correspondingly (therefore cannot claim old readers support the new large distances). Before coding, inspected exact last-seen previous chunks beyond 128 KiB in each complete Silesia original. Nonconstant aligned 64-byte far duplicates: mozilla ~920,256 B, samba ~428,160 B, ooffice ~11,392 B, nci ~6,827,776 B. These are merely **sample occurrence counts**, not net compression savings; nci already selects BWT, so its long-range opportunities would not enter the current math path. The screen did not require a new token grammar.

**ONE actual complete integrated original 12-file candidate-only Silesia encode/decode run**, exact all 12 original reconstructed files including integrity digests: **60,318,613 B / 5.932930 s encode / 2.485692 s decode**. Compared with R57 60,585,717 B, genuine **-267,104 B (-0.44087% archived size)**, confined to math-selected mozilla **-198,682 B**, ooffice **-20,769 B**, samba **-47,653 B**. Other nine file archive sizes unchanged. Measured encode and decode were slower than R57 retained ten-run medians by **+7.99% and +13.27% nominal**, respectively. THESE ARE ONE-SCREEN unpaired timing deltas, not a ten-run candidate result. The failed tri-metric screen justifies rejecting this concrete 512 KiB-history implementation without varying history-size thresholds.

## Required full original Silesia comparison

| Compressor | All 12 archived bytes | Original-size reduction | Encode s | Decode s |
|---|---:|---:|---:|---:|
| Zstd-3 (historical external CLI) | 66,216,569 | 68.76% | 1.2561 | 0.8551 |
| XZ-6 (historical external CLI) | 49,233,340 | 76.77% | 116.9892 | 2.3079 |
| **Current approved JLZ R57 (retained ten-run median)** | **60,585,717** | **71.41355%** | **5.493744** | **2.194428** |
| **R63 longer mathematical history (ONE new complete screen)** | **60,318,613** | **71.53958%** | **5.932930** | **2.485692** |

| R63 archive delta vs R57 | Encode-time delta vs R57 | Decode-time delta vs R57 |
|---:|---:|---:|
| **-0.44087%** | **+7.99% nominal, ONE run** | **+13.27% nominal, ONE run** |

R62 retained **60,585,717 B** with no size change but failed to show required speed parity. Approved JLZ/Zstd/XZ timing references reused rather than rerun; historical CLI and in-memory benchmarks have different timing conventions.

## Stop decision and forward queue

**REJECT H1 for promotion of R62 and R63.** Preserving size alone is insufficient; all three production metrics must not regress without express user approval. R62 confirms safe mathematical redundant-work reductions but no demonstrated speed win; R63 demonstrates real far-back repetition on the three math-selected files, but gains only 0.44% archived bytes in exchange for much slower encoding/decoding. No 10-run candidate promotional median was warranted. **Stop iterating reciprocal tweaks, column-decoder subvariants and history-length threshold sweeps.**

New original architectural research may investigate *bounded-cost long-distance exact-match indexing* (rather than a larger per-byte search space), or an efficient native suffix-array constructor that avoids expensive Go standard-library index serialization and can preserve existing BWT output bytes. Both are OPEN hypotheses, not accepted gains. Keep R57 production, preserve this negative finding, and return to external architecture literature for genuinely distinct approaches rather than repeating tiny changes.