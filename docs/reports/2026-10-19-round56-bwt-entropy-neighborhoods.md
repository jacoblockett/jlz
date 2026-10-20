# Round 56 — BWT-local entropy regimes versus uniform whole-block coding (REJECT production promotion; positive component and complete candidate evidence)

**Date:** 2026-10-19. **Production baseline remains accepted Round55.** All experiments used original complete twelve-file Silesia (211,938,580 B, 12/12 canonical source SHA-256 passed), native Go and no outside compression backend or source copying. This is a new practical test of **nonstationary post-BWT entropy neighborhoods**, inspired by research comparing MTF, distance coding, incremental-frequency postprocessing and QLFC. It is not a Zstd, libbsc or Brotli implementation.

## Tested implementation

The accepted BWT uses native suffix transform → bijective MTF rank-zero coding → one independently modeled finite-state stream per 256 KiB source block. A BWT naturally sorts contexts into nearby neighborhoods, so its emitted symbol histogram need not be stationary across a complete block.

Implemented a new reversible experimental block subtype 16. One BWT transform and one existing bijective rank transform feed a comparison among:
- existing unchanged subtype 15: one fully framed native entropy stream;
- subtype 16: two or four contiguous, independently framed finite-state streams, with a one-byte partition-count header.
A distinct expensive exploratory configuration also compared eight parts. Selection is **strictly shorter fully encoded block payload**, including all model tables and framing. Decoder reads old subtype 15 unchanged, and for 16 reads the count, validates each stream, joins exact symbols, performs existing bijective MTF inverse, and then BWT inverse. The three sampled BWT preflight blocks stay unchanged/byte-identical to accepted Round55; only the remaining BWT-winning blocks use the supplemental choices. No new CLI options and no mathematical winner paths replaced. The new subtype **would require newer decoders**; this is counted against promotion.

Code implemented and tested in local disposable directories, not promoted to GitHub:
- /mnt/data/jlz_r56_trial (2,4,8 alternatives)
- /mnt/data/jlz_r56_24 (2,4 alternatives)
- /mnt/data/jlz_r56_4 (4-only screen)
Only the actual research *findings* are retained in the repository, not a disabled experimental decoder.

## Actual selected BWT component screens

Full original source coverage of the **seven BWT-selected** Silesia files (dickens, mr, nci, reymont, webster, xml, osdb), 451 source blocks / 117,233,324 source bytes, each subjected to the actual BWT+bijective symbol transform and actual native packed-stream compression. Baseline 23,415,904 B. Independently framed equal-length 2/4/8/16-part alternatives (charged 1 selector byte):
- fixed 2 parts: 23,299,606 B, **-116,298 B**; 432/451 individually win;
- fixed 4 parts: 23,215,283 B, **-200,621 B**; 304/451 individually win;
- fixed 8 parts: 23,302,598 B, **-113,306 B**; 175/451 individually win;
- fixed 16 parts: 23,520,592 B, **+104,688 B**; 109/451 individually win.
The ideal select-per-block among these choices yields a **269,091 B** upper bound of isolated savings; this is NOT a whole-archive result. Models are complete native serialized models, not entropy estimates.

A separate, genuinely reversible and decoder-validated proposed *second-rank and nonzero-rank repeat escape* (unused 255,2 marker in the accepted bijective MTF alphabet) was tested on all the same 451 blocks. All fixed candidates **expanded** the selected symbol streams by 29,955–33,107 B; even per-block smaller-of-old/new saved **at most 2,296 B**. **REJECT that specific rank1/rank2/nonzero-rank run family** as not worth a new subtype.

## Fully integrated, complete original twelve-file candidate results

Both candidates were independently integrated into JLZ's actual selected BWT winners, not merely packer components; every full archive includes all native metadata, checksums and unaffected math/lane families.

| Candidate, all 12 complete Silesia files | Archive bytes | Encoding seconds | Decoding seconds | Evidence |
|---|---:|---:|---:|---|
| Full 2/4/8 candidate | 60,473,492 B | 5.980229 | 2.211946 | NEW ten candidate-only complete-corpus medians |
| Four-only | 60,522,589 B | 6.026967 | 2.194522 | One complete screen, not a median |
| **Reduced-work 2/4 candidate** | **60,513,252 B** | **5.685114** | **2.203662** | NEW ten candidate-only complete-corpus medians |

The full 2/4/8 candidate saved 252,941 B over accepted R55, but encoding was nominally **+8.73%** slower. The 2/4 reduced candidate saved 213,181 B (**-0.35105% of compressed archive bytes**) while encoding was nominally **+3.3607%** slower and decoding nominally **+0.2374%** slower. Previous R55 times were independently measured and never rerun: 60,726,433 B, encode 5.500265326 s, decode 2.198443653 s. Those nominal timing comparisons are **unpaired**, so minor percentages are not causal proof.

### Reduced 2/4 candidate ten-run record (full 12 files compressed AND decoded, all exact)

| Run | Entire archived bytes | Encode seconds | Decode seconds |
|---:|---:|---:|---:|
| 1 | 60,513,252 | 5.955413488 | 2.290250609 |
| 2 | 60,513,252 | 5.689999959 | 2.173121856 |
| 3 | 60,513,252 | 5.680227363 | 2.132985481 |
| 4 | 60,513,252 | 5.669594919 | 2.166364166 |
| 5 | 60,513,252 | 5.730832570 | 2.220863104 |
| 6 | 60,513,252 | 5.600612814 | 2.238929400 |
| 7 | 60,513,252 | 5.627445827 | 2.186460485 |
| 8 | 60,513,252 | 5.583245734 | 2.123979543 |
| 9 | 60,513,252 | 5.804384359 | 2.302155130 |
| 10 | 60,513,252 | 5.811118505 | 2.250308695 |
| **Median** | **60,513,252** | **5.685113661** | **2.203661795** |

Every one of ten runs recreated the same bytes and decoded all 12 source files exactly, 120/120 complete original-byte roundtrips.

Per-file change R56 2/4 vs accepted R55:
- dickens **-42,252 B**
- mr **-98,046 B**
- nci **-7,222 B**
- osdb **-1,793 B**
- reymont **-9,126 B**
- webster **-48,977 B**
- xml **-5,765 B**
- mozilla, ooffice, samba, sao, x-ray **byte-identical**.
- **Worst file-level size regression: 0 B**. These actual complete file archive sizes sum to 60,513,252 B.

## Distinct math metadata segmentation feasibility screen

Used full Round55 correctly selected archives and extracted **598 actual selected mathematical opcode/length streams** from mozilla (384), ooffice (48), samba (166); osdb became BWT-selected under Round55 and contributes zero actual mathematical streams. Original native re-packed cost **4,031,551 B**. Independent 2/4/8 equal-length partitions (each charged one selector byte, separate table + full packed stream framing) cost **4,033,147 B** (+1,596), **4,068,423 B** (+36,872), and **4,155,052 B** (+123,501). Best-per-stream oracle saved 74,917 B across 195/598 selected streams, before integration selection costs. All candidate pieces were truly packed, unpacked and joined exactly. **REJECT fixed math segmenting**, preserve per-stream optional nonstationary modeling OPEN only if a low-cost model selector is demonstrated.

## Method and validation

- Only complete original 12-file Silesia; canonical source SHA-256 manifest **12/12 PASS**, exactly 211,938,580 bytes.
- For **both ten-run candidate configurations**, all 120 complete file archive reconstructions exact. All encoded MTF component tests also decoded exactly.
- For the 2/4 reduced candidate: \`go test -count=1 ./...\`, \`go vet ./...\`, targeted \`go test -race -count=1 ./engine\` on finite-state/BWT selected paths, and Windows/amd64 CGO-disabled cross-compile PASS (Windows executable not run).
- The research-only segment decoder strictly validates primary index, 2/4/8 count, individual framed stream sizes, decoded symbol output bounds, and trailing bytes. No accepted codec changes in main; no older-decoder capability falsely claimed for the experimental mode16.
- Accepted JLZ R55, Zstd-3 and XZ-6 references were **not rerun**.

### Canonical comparison for best *practical* candidate

| Model | Compressed size | Reduction vs 211,938,580 B | Encode seconds | Decode seconds |
|---|---:|---:|---:|---:|
| Zstd-3 (historical external CLI) | 66,216,569 B | 68.76% | 1.2561 | 0.8551 |
| XZ-6 (historical external CLI) | 49,233,340 B | 76.77% | 116.9892 | 2.3079 |
| **Current approved R55 (retained ten-run medians)** | **60,726,433 B** | **71.34715%** | **5.500265** | **2.198444** |
| **Candidate R56 2/4 (new ten-run medians)** | **60,513,252 B** | **71.44774%** | **5.685114** | **2.203662** |

| Candidate archive size delta | Candidate encode time delta | Candidate decode time delta |
|---:|---:|---:|
| **-0.35105%** | **+3.36072% nominal** | **+0.23736% nominal** |

The historical external CLI runtimes are not identical to in-memory JLZ measurement conventions.

## Decision and active next theory queue

**REJECT H1 *as a current production integration*, but retain positive experimental evidence.** The new context-neighborhood model actually reduces entire archives without any per-file size regressions. However, 0.35% archive improvement currently requires nominal ~3.36% extra encoding work and a new decoder block subtype; 2/4/8 is even more disproportionate. Under YAGNI and measured practical overhead, no new decoder complexity is promoted yet, and the approved production baseline stays R55. This is **not** rejection of BWT entropy regime segmentation as a mathematical idea.

**Next distinct test**: an inexpensive minimum-description-length change-point estimator over a small number of BWT histogram windows, picking **one** local model partition without repeatedly fully compressing alternatives. It must prove that the saved encoding work retains material complete-archive byte gains before any promotion. An independently decoded segmented BWT block can use deterministic/few-byte cut boundaries if economically justified; avoid a threshold-parameter spiral.

Other theoretically grounded OPEN avenues for substantially different experiments: 
1. **Brotli-style context IDs** on selected JLZ distance extra-bit/slot streams, conditioned on *already decoded* match length (not R50's concatenation or opcode-conditioned first varint-byte permutation); contextual statistics and model costs must be fully charged.
2. **Incremental/local frequency count** instead of MTF for BWT output, leveraging published alternatives; do not just retest bijective zero-run radix tuning. 
3. **BWT adaptive nonstationary regimes** (measured now) with MDL segmentation and reused sampled BWT buffers. 
4. **Cheap inverse distance/inversion-frequency transforms** on originally selected BWT streams, leveraging the formal equivalence/limits literature; treat prior generic geometry negatives as unrelated, but require genuine native reversible code and selected-stream wins before integrated testing.
5. Higher-order context-mixing and PPM-like models remain mathematically legitimate but computationally riskier for the user’s strong practical performance constraints.

Research inspirations, not code used:
- Travis Gagie and Giovanni Manzini, “Move-to-Front, Distance Coding, and Inversion Frequencies Revisited,” Theoretical Computer Science 411 (2010), https://doi.org/10.1016/j.tcs.2010.04.024
- Jürgen Abel, “Incremental frequency count: A post BWT-stage for the Burrows–Wheeler compression algorithm,” Software: Practice and Experience 37 (2007), https://doi.org/10.1002/spe.763
- RFC 7932 (Brotli, context modeling / block maps), https://www.rfc-editor.org/rfc/rfc7932
- libbsc documentation for QLFC family, https://github.com/IlyaGrebnov/libbsc/blob/master/libbsc/coder/qlfc/qlfc.h

Production engine remains untouched; negative/prospective results remain recorded as actual experiment statuses, not blanket theorem disproof.