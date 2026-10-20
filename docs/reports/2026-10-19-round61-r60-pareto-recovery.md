# Round 61 — User-directed reversal of R60; strict three-metric acceptance; recovery experiments (REJECT H1)

**Date:** 2026-10-19. **FINAL: REJECT R60 production, restore approved R57.** This supersedes the Round60 report's erroneous provisional ACCEPT decision. User specifically objected that R60's improvement from 71.41355% to 72.00650% source reduction (about **+0.593 percentage points**) was not worth **+0.520665 s encoding** and **+0.150170 s decoding** based on retained ten-run medians. That judgment was correct; no agent may unilaterally accept a speed/ratio tradeoff, even when archive bytes decrease. **Current approved benchmark: R57, 60,585,717 B / 5.493743735 s encoding / 2.194427893 s decoding.**

## Strict acceptance contract, effective immediately

Any candidate to be promoted without explicit user permission must improve or preserve **all three complete original 12-file Silesia metrics**: total fully framed archive bytes, ten-run median compression time, and ten-run median decompression time. A change that gains size but loses either speed metric, or gains speed but grows the archive, is **REJECTED**. Small unpaired timing differences within noise are insufficient for a claim of speed improvement; do not label a regression acceptable by asserting the ratio change is "substantial." Tests, original 12-file exact roundtrips, metadata overhead, per-file regression disclosure and backwards decoding are still required. Don't cherry-pick easy files or rely on thresholds tuned to Silesia without broader rationale. When no true Pareto winner exists, leave production source unchanged and keep evidence in research docs.

## Actual rollback on main

Removed R60-only engine/bwt_extended_test.go. Restored exact approved pre-R60 Git source blobs:
- engine/bwt.go: 3f4af2d234c9382b2a3f9c20fd649974e0871e09
- engine/cross.go: 59794f52ac58b0712996a706b6c39b25d18ed7ca

These produce the accepted R57 256 KiB independent BWT regions, native MDL two/four-part entropy choices, original 256 KiB per-block decoder bound, old cross-block dependencies and existing backwards-readable archive format. No R60 512 KiB frame is now emitted by main. R60 is kept as a **rejected research candidate**, not promoted production.

## Revisit R60 without wholesale abandonment

Three actual native-Go disposable prototypes attempted to recover the 512 KiB size gain while controlling speed. No source code was promoted.

### Candidate A: packed-LF inverse BWT for 512 KiB regions

The existing BWT inverse touched two separate arrays per row: a uint16 byte/sentinel column and an int32 LF-index column. Rewrote the inverse to use a **single uint32 entry per row** storing 23-bit next-row index and 9-bit byte-plus-sentinel, with exact, validated byte restoration. No archive-format modification; 23-bit bounds work for R60's maximum 512 KiB input size. Local experiment /mnt/data/jlz_r61_lfpacked: Go BWT and engine tests PASS, independent whole 12-file screens with all original bytes exact:
- **59,329,033 B**, 6.081519 s encode, 2.172160 s decode (ONE screening pass).
- Follow-up complete screens 6.555175/2.314988 s, 6.049722/2.277116 s, 6.059376/2.185206 s encode/decode.
The decoding improvement is inconsistent; compression time remains well above the required R57 **5.493744 s** in every screen. **REJECT current 512 KiB + packed-LF candidate**, no 10-run median and no promotion.

### Candidate B: bounded 512 KiB regions only where total source is at most 16 MiB

Experimentally retain the original 256 KiB BWT region size and 3 preflight sample payload reuse for sources larger than 16 MiB. For smaller BWT-selected sources allow 512 KiB regions and use the packed-LF inverse. This conservative memory/parallel-work bound preserves most source's historical format, without file-specific names. It is still a threshold experiment and must not be tuned to individual Silesia files. One complete all-12-file candidate-only exact screen:
- **60,022,327 B**, 6.298676 s encode / 2.233001 s decode.
- Compression remains **563,390 B smaller** than R57, but **+0.804932 s encoding / +0.038573 s decoding** versus historical retained medians, nominal unpaired.
**REJECT**. This is not a speed gain and cannot be promoted; do not parameter-sweep neighboring source-length thresholds.

### Candidate C: same packed-LF inverse applied only to approved 256 KiB R57

Original R57 engine checkout copied without format changes; only its BWT Inverse was substituted with the native packed-LF indexing above. Full BWT+engine uncached Go tests PASS, exact selected archives on the twelve original files still **60,585,717 B**. One complete source screen 5.788924 s encode/2.128431 s decode; **four additional completed candidate-only full 12-file runs** produced respectively
- 5.881639 / 2.109165 s
- 5.641655 / 2.105094 s
- 5.553917 / 2.053380 s
- 5.799253 / 2.123743 s

Every candidate archive was exactly **60,585,717 B** with all original files reconstructed byte-exact. The sixth through tenth attempted run wrapper timed out and did **not** produce a complete benchmark; no missing runs were interpolated or included. Decode appears directionally better, but all completed encoding times are **above** the independently measured R57 median **5.493744 s**; no demonstrated all-three-metric improvement. **REJECT promotion and stop rather than grind another ten rounds of nearby tuning**. No spurious candidate 10-run medians reported.

## Required complete-corpus reference table

| Compressor | Fully framed complete Silesia archive B | Source reduction | Compress s | Decompress s |
|---|---:|---:|---:|---:|
| Zstd-3, historical external CLI | 66,216,569 | 68.76% | 1.2561 | 0.8551 |
| XZ-6, historical external CLI | 49,233,340 | 76.77% | 116.9892 | 2.3079 |
| **Current approved JLZ R57, retained previous ten-run median** | **60,585,717** | **71.41355%** | **5.493744** | **2.194428** |
| **Rejected R60 512 KiB candidate, historical ten-run median** | 59,329,033 | 72.00650% | 6.014409 | 2.344598 |

| R60 compressed size delta vs R57 | R60 compression time delta | R60 decompression time delta |
|---:|---:|---:|
| **-2.074225%** | **+9.4774% nominal** | **+6.8433% nominal** |

Original old/new timing sessions were unpaired, but the size/time tradeoff was large enough to violate the user's approval criteria. R61 experiments are mostly one-pass screens; they are **not** alternative accepted ten-run results. Historical CLI and in-memory measurement conventions differ.

## External compressor architecture remains an open research source

The assistant improperly narrowed to repeatedly refining JLZ BWT/MTF model boundaries, despite earlier successful inspiration from Zstd's simpler parser and native finite-state entropy. External designs remain a source of **abstract architectural inspiration**, not code or backend:
- **LZFSE:** match-sequence generation coupled to finite-state entropy and compact decoding; explore how efficient sequence/literal separation can complement JLZ's existing mathematical paths, without copying or importing an external codec.
- **LZ4:** deliberately cheap matching, byte-oriented literal/match decoding and block independence. The precise existing fast dual-hash shortcut R54 was rejected for ratio loss, so do NOT revive it without a *distinct* additive mechanism which preserves all metrics.
- **Zstd:** multi-table frame/model reuse, fast candidate selection, and exact match-sequence entropy mechanisms; R52/53 already adapted a native unconditional finite-state approach, so do not repeat it as new.
- **Block-sorting architectures:** larger windows are genuinely useful but require an equally efficient sorting and inversion pipeline before any promotion; R60 showed current implementation too slow, so do not repeat adjacent window sweeps.
Research designs, not imported code: https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_double_fast.c ; https://github.com/lz4/lz4/blob/dev/doc/lz4_Block_format.md ; https://developer.apple.com/documentation/compression/compression_lzfse .

**Next independent experimental work**, not automatically accepted: a genuinely low-work mathematical source representation that improves ratio and throughput together, or direct suffix-array construction with fewer intermediate allocations but preserving exact bitstream and strict decoding time; require complete original-Silesia candidate-only ten-run validation before promotion. Avoid new rounds of negligible BWT context parameter tuning.

## Final disposition

**REJECT H1; production R57 restored on main.** This reversal corrects earlier unilateral acceptance. Round60 remains documented and its 512 KiB experiment reproducible in disposable local source, not a canonical production codec. No feature rollout, new archive format, or external codec dependency from Round61.