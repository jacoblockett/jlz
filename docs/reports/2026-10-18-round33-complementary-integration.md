# Round 33: Complementary integration of catalog theories (candidate rejected; open paths retained)

**Date:** 2026-10-18. **Approved source remains Round 29.** This round responds to the user's directive to revisit every catalog theory as a potentially *additive* representation in JLZ, not automatically an original-byte replacement. A 40-row source inventory and explicit integration map is recorded in `docs/research/complementary-integration-audit.md`; this is an **audit**, not a false claim that all 40 families have been implemented. It identifies where each can touch the current match, BWT, and additive streams, and distinguishes actually tested combinations from still-open work.

## Current JLZ integration surfaces actually examined

1. The **already computed** 16-bit modular/additive residual low and high coordinate streams, block type 14, from Round 29. Candidate subcodecs are invoked only if they beat native entropy payload size, accounting for per-stream tag bytes.
2. BWT's **MTF + zero-run symbol stream**, while retaining the exact existing suffix sort and inverse. Native packed stream metadata and a per-stream selector counted.
3. **Previously selected packed match blocks** from the fixed Round 25 Silesia archive files unchanged by Round 29. The archive was parsed into its five existing packed streams; no new baseline or established external compressor benchmark. Actual stored arithmetic/distance-slot payload lengths were measured, rather than recompressing references.

Only original Silesia (12 files, 211,938,580 B total) was used as a performance source. Sampled tests are explicitly distinct from full archive results. All the tested transformations were *executable, exact-inverse* Go prototypes and included their own native packed metadata. No outside compressor is in the pipeline.

## Residual coordinate component test

Three distributed original-data 256 KiB (or shorter last) regions per Silesia source = **36 blocks / 7,728,404 source B**. For each, the selected candidate only changes the entropy representation **of already computed additive low/high residual streams**, with baseline **5,966,439 B** of native packed lane payloads. Figures are optimistic minima independently chosen per region and per coordinate, including a byte tag per replaced stream; the unselected files' improvements must not be imputed to actual Round 29 archives.

| Complementary representation | Selection wins / 36 | Optimistic aggregate packed B | Saving against native-additive-lane components |
|---|---:|---:|---:|
| 16-rule Tietze-style nested pair grammar | 16 | 5,669,126 | 297,313 |
| Zero-default sparse-defect positions + values | 10 | 5,870,142 | 96,297 |
| Horizontal boundary distances | 8 | 5,885,322 | 81,117 |
| Vertical integrated boundary differences | 7 | 5,904,768 | 61,671 |
| Bitplane split with eight individually packed planes | 5 | 5,939,601 | 26,838 |
| Two-dimensional mixed difference, grid width 16 | 1 | 5,963,329 | 3,110 |
| Grid width 256 | 1 | 5,966,367 | 72 |
| Prime-index partitions | 2 | 5,966,428 | 11 |
| XOR residual lags (1,2), (1,4), (7,11); high XOR low | 0 | 5,966,439 | 0 |

The x-ray sample is the critical one because x-ray is the sole complete original-file archive currently selecting additive. On its three sampled source regions the lane payload baseline totaled **310,270 B**; sparse defect minima **275,635 B** (34,635 B saved), and grammar16 **276,610 B** (33,660 B saved). This motivated a true additive-integrated sparse-defect candidate.

## BWT symbol-stream complementary test

Six original Silesia files (dickens, mozilla, nci, reymont, webster, xml), first and middle region each, 12 exact 256 KiB input regions. Each transformed symbol stream is created using existing JLZ BWT+MTF-zero-run; alternative entire packed native streams are decoded back to the exact MTF bytes, and the original input block through inverse MTF/BWT. **This does not assert all 12 blocks are selected as BWT in final archives.**

Native BWT entropy stream total: **742,585 B**. Prime partitioning, independently Huffman packing prime/nonprime symbols, wins **6/12** and saves **1,580 B** when selected only on winning streams. Sparse defects save **81 B** on one of 12. Grammar16, XOR, grid, and boundaries win **0/12**. These are real component bytes but do not justify modifying BWT framing; it would also require screening CPU and tags for small full-archive benefit.

## Truly selected match metadata streams

Read immutable approved Round 25 archives for ten Silesia files unchanged in Round 29 (excluded files where a different whole-file representation was selected); **338 packed match blocks** appeared. From nine distributed selected blocks, unpacked *their existing* packed length, distance and literal streams and compared fully serialized alternative payloads to the original stored codec cost. All alternatives reconstructed stream bytes exactly; there was no rematching or baseline compression.

Across **nine length streams**, Tietze16 with exact generator definitions saves **15,138 B** on six streams, relative to current actual packed length bytes totaling **100,709 B**. The sampled literal streams (**474,478 B** stored) and distance metadata (**96,048 B** stored) show zero improvement from the same grammar. Prime partition saves only **162 B** and bitplane split **26 B** on one length stream each. Sparse, XOR, grid and boundary representations each won **0/27** of the nine lengths/nine distances/nine literals tested. Thus the additive grammar hypothesis has a **specific target** in actual successful JLZ block encodings: match-length streams, not general literals.

## Fully integrated sparse-defect candidate and benchmark

A candidate-only disposable Go version augments Round 29's additive block type 14. It preserves the approved mode-4 payload exactly for all regions where sparse streams do not win, and adds payload subtype 5 containing two per-lane selector bytes and native packed positions/exceptions (source zero-default). It retains raw-region fallback, BWT/additive sampling gates, checksum and standard archive framing; the old mode 4 decoder is retained alongside subtype 5. Importantly, an initial version redundantly recomputed the old native streams; a cleaned candidate computes them once and shares the result with sparse selection. No source was promoted.

Each of **ten complete original-Silesia candidate-only archive runs**, from independently spawned precompiled Go 1.23.2 linux/amd64 binary processes, produced the exact same aggregate **62,147,274 B** archive bytes and restored all 12 source files exactly. Before-candidate approved Round 29 **62,644,663 B** was reused, not rerun. Ten candidate whole-corpus medians **7.9441650235 s encode** and **3.005318434 s decode** versus retained historical Round 29 **7.516364 / 2.821363 s**. The reference and candidate sessions are **not paired**, so percentage CPU changes are directional and may include environment load.

| Model | Complete Silesia archive bytes | Reduction from 211,938,580 B | Encode seconds | Decode seconds |
|---|---:|---:|---:|---:|
| Zstd-3 fixed historical CLI | 66,216,569 | 68.7567% | 1.2561 | 0.8551 |
| XZ-6 fixed historical CLI | 49,233,340 | 76.7700% | 116.9892 | 2.3079 |
| Approved JLZ Round 29 (accepted historic 10-run median) | 62,644,663 | 70.4421% | 7.516364 | 2.821363 |
| **Round 33 integrated candidate** (current 10-run median) | **62,147,274** | **70.6768%** | **7.944165** | **3.005318** |

| Candidate-vs-current archive-size delta | Encode-time delta | Decode-time delta |
|---:|---:|---:|
| **-0.79398%** (497,389 B saved) | **+5.6916%** (slower) | **+6.5201%** (slower) |

Zstd and XZ use historical CLI timings including startup/I/O rather than JLZ's memory benchmark; do not compare CPU speeds as if measured identically. No benchmark of Zstd/XZ or accepted Round 29 was performed in this round.

## Verification and exact limitations

- Complete 12/12 original Silesia roundtrips in **each of 10 candidate runs** with full JLZ1 SHA-256 authenticity. All ten candidate archives have identical aggregate size.
- Exact inverse Go tests for variable-sized sparse streams, random/periodic/adversarial 16-bit words, zero runs, malformed headers and truncation; matching packed stream variants, BWT/MFT/inverse and geometric/grammar screens were tested with exact bytes.
- Targeted `go test -count=1` for additive/sparse tests, `go vet ./...`, and targeted `go test -race -count=1`: **PASS** in the disposable candidate source.
- A broad uncached `go test -count=1 ./...` in the disposable local tree was interrupted while running multiple unrelated long-running research tests; this was **not** counted as a passing full-suite run. A new 27-archive cross-version compatibility matrix and CI Test were **not** run because the size/speed candidate is rejected.
- No permanent production Go source, user-facing CLI option, or new block syntax is committed. Local prototype source is not a promoted codec. Historical references and the active theory backlog remain valid.

## Recommendation and next integration

**REJECT the integrated sparse-defect subtype:** the real 0.794% complete-archive saving comes with ~5.69% encode and ~6.52% decode regression versus stored accepted medians. That is an unfavorable trade for additional archive parsing/format variants. Maintain approved Round 29, including its current additive representation.

**Do not reject the sparse-defect, grammar, prime, geometric or algebraic theories generally.** These companion screens show positive results in *specific intermediate streams*, most strongly 16-rule grammar on selected packed match **lengths**. The next candidate should test a **low-CPU, gated native grammar/dictionary representation just for the measured winner match-length streams** (not mandatory across all blocks), using exact framing/model costs and a bounded preprocessing budget. Other untested theory×stream combinations explicitly remain OPEN in the 40-family complementary audit. The preferred development direction is *additive to JLZ*, not replacing it without transformative complete end-to-end evidence.
