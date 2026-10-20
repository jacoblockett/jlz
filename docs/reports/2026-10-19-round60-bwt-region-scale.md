# Round 60: Larger BWT context horizons; ratio gain but unacceptable time regression (FINAL REJECT H1)

**Date:** 2026-10-19. **FINAL: REJECTED by the user and rolled back, restoring R57.** The original provisional acceptance was erroneous. R60 is historical candidate research only, not the production baseline. Only the complete original twelve-file Silesia corpus (211,938,580 B) and its selected original regions were analyzed. Original corpus SHA-256 manifest 12/12 PASS. No external compressor, copied compression algorithm, new command, learned per-file table or GitHub Actions.

## Hypothesis A: decoder-synchronized symbol remapping (REJECT and STOP)

To move beyond R59's costly three-context native finite-state coder, implemented actual reversible previous-symbol subtraction, XOR, deterministic condition-dependent byte permutations, and context-specific up-to-four transposition mappings of high-frequency bijective rank codes. Every sampled and selected BWT stream was restored exactly; context followed the previously decoded original byte. The transposition tables were costed with a favorable two-byte descriptor per swap, not actually integrated into a full archive. These are **optimistic component estimates**, not a measured complete compressor.

First/middle/end **21 selected original Silesia BWT blocks**: baseline single-table symbols with 1-byte mode estimate **965,521 B**; previous-byte subtraction **1,156,780 B**, XOR **1,073,184 B**. Best fixed peak-remap **964,766 B**, saving just 755 B. On **all 451 actually BWT-selected source blocks**, the independently framed native 1/2/4-table incumbent component oracle measured **23,192,014 B**. Even a freely chosen best-of-nine remapping per block measured **23,182,941 B**, a maximum **9,073 B (0.039%)** gain on only 44 blocks, BEFORE actual per-block selectors/model serialization. This is not meaningful and was **REJECTED**. Close the entire adjacent role/byte remap family rather than repeating context threshold or transposition sweeps.

## Hypothesis B: larger original BWT region, unchanged BWT/MTF/finite-state mathematics

An orthogonal mechanism: BWT only sorts suffix contexts inside its input region. The accepted 256 KiB limit cut related contexts at each original input block boundary. Increase only **already selected BWT regions**, leaving 256 KiB math/cross-history regions and all automatic codec selection unchanged.

On **19 distributed aligned original-Silesia pairs**, exact native original 2×256 KiB BWT/bijective/entropy packing used **2,250,636 B**, one 512 KiB BWT/entropy used **2,139,702 B**, 19/19 wins, **-110,934 B (-4.93%)** sample-only. BWT+native entropy sampling elapsed **713.1 ms vs 717.0 ms**, near equal. All sampled BWT transforms inverted exactly.

One completely integrated 12-file encode/decode screen each (NOT medians):
- **512 KiB** BWT input regions: **59,329,033 B / 5.9414 s compress / 2.3879 s decompress**.
- **1 MiB**: **58,364,486 B / 6.2655 s / 3.0795 s**.
- **2 MiB**: **57,627,775 B / 6.5741 s / 5.6384 s**.

All twelve original files recovered exactly in each screen, and no individual archive size regressed. The 1 MiB and 2 MiB configurations increase decoding working set and wall time disproportionally relative to extra compression, so both were **REJECTED** as production configurations. The region-size sweep **ends at these three meaningful sizes**.

## Actual approved implementation: 512 KiB only

The selected BWT encoder in **engine/bwt.go** uses 512 KiB source chunks and the existing mode15/16 bijective+finite-state coding with existing shorter-of fallback. Existing three 256 KiB preflight samples continue as routing evidence, but cannot be reused as 512 KiB final payloads. Math/cross-reference encoding and 128 KiB history never change. No new entropy subtype.

The existing **engine/cross.go** decoder permits 512 KiB block lengths only for BWT mode15/16 or raw fallback, while mathematical/cross-history blocks stay restricted to 256 KiB. Worker output offsets are derived from validated per-block prefix sums rather than multiplying a uniform block size. A BWT/raw-only archive decodes each independent block in its own parallel job; history-dependent archives retain sequential 4-block dependency groups. **engine/bwt_extended_test.go** adds multiple-group offset recovery and explicit malformed mode/oversized-cross-history rejection.

The archive format uses the same mode IDs and framing but **changes accepted block-size bounds**. New R60 decoder restores old archives; pre-R60 decoders reject newly emitted 512 KiB BWT blocks. This forward-incompatibility is expected and explicit.

## Final promotional benchmark: TEN complete original 12-file candidate-only runs

First 512 KiB prototype without independent decode jobs: 10 complete passes all **59,329,033 B**; median **5.888722004 s encode / 2.420466870 s decode**. A final, separately benchmarked version with independent BWT decode jobs produced the same bytes and these **new candidate-only ten complete pass** results:

| Run | Complete archive B | Encode s | Decode s |
|---:|---:|---:|---:|
| 1 | 59,329,033 | 5.942740033 | 2.395683112 |
| 2 | 59,329,033 | 5.795084544 | 2.232461854 |
| 3 | 59,329,033 | 5.917361075 | 2.417032148 |
| 4 | 59,329,033 | 6.010076659 | 2.376140834 |
| 5 | 59,329,033 | 6.229778026 | 2.364702526 |
| 6 | 59,329,033 | 5.888616525 | 2.324494035 |
| 7 | 59,329,033 | 6.018740688 | 2.312011237 |
| 8 | 59,329,033 | 6.019547945 | 2.323717953 |
| 9 | 59,329,033 | 6.391317801 | 2.460701197 |
| 10 | 59,329,033 | 6.091428312 | 2.302365560 |
| **Median** | **59,329,033** | **6.014408674** | **2.344598280** |

**All 120 original source-file reconstructions byte-exact; archive size identical across all ten runs.** Retained prior approved R57 (independent ten-run historic medians) **60,585,717 B / 5.493743735 s encode / 2.194427893 s decode**, not rerun for speed. **True -1,256,684 B (-2.074225%) compressed archive size**, nominal **+9.4774% encoding / +6.8433% decoding**. The timing measurements come from independent sessions and cannot be interpreted as precise causal estimates. This extra bounded CPU is justified by a much larger measured compression gain than previous context micro-experiments.

### Required full Silesia comparison

| Model | Archive bytes | Reduction of original 211,938,580 B | Compression time | Decompression time |
|---|---:|---:|---:|---:|
| Zstd-3, retained historical external CLI | 66,216,569 | 68.76% | 1.2561 s | 0.8551 s |
| XZ-6, retained historical external CLI | 49,233,340 | 76.77% | 116.9892 s | 2.3079 s |
| **Approved R57, retained earlier ten-run median** | **60,585,717** | **71.41355%** | **5.493744 s** | **2.194428 s** |
| **Candidate R60, NEW ten-run median** | **59,329,033** | **72.00650%** | **6.014409 s** | **2.344598 s** |

| Compressed-size delta vs R57 | Encode-time delta vs R57 | Decode-time delta vs R57 |
|---:|---:|---:|
| **-2.074225%** | **+9.4774% nominal** | **+6.8433% nominal** |

Zstd/XZ external CLI methods differ from JLZ in-memory timing conventions. Reuse references; do not rerun Zstd/XZ/current accepted codec timings in every future round.

### Complete per-file compressed sizes

| File | R57 historical B | R60 newly measured B | Delta B |
|---|---:|---:|---:|
| dickens | 3,155,856 | 3,004,645 | -151,211 |
| mozilla | 18,609,139 | 18,609,139 | 0 |
| mr | 2,729,261 | 2,673,111 | -56,150 |
| nci | 2,163,199 | 2,006,947 | -156,252 |
| ooffice | 3,033,712 | 3,033,712 | 0 |
| osdb | 3,425,086 | 3,181,804 | -243,282 |
| reymont | 1,455,234 | 1,377,274 | -77,960 |
| samba | 5,134,262 | 5,134,262 | 0 |
| sao | 6,261,713 | 6,261,713 | 0 |
| webster | 9,835,416 | 9,298,374 | -537,042 |
| xml | 515,852 | 481,065 | -34,787 |
| x-ray | 4,266,987 | 4,266,987 | 0 |
| **TOTAL** | **60,585,717** | **59,329,033** | **-1,256,684** |

**Seven files smaller, five unchanged, worst file-level archive regression zero.**

## Validation and compatibility

- All **12/12** canonical original Silesia source SHA-256 match and source bytes sum 211,938,580.
- Final ten-run candidate: **120/120** exact original restored files, SHA and archive framing verified.
- **go test -count=1 ./... PASS**, **go vet ./... PASS**, targeted **go test -race -count=1 ./engine** on extended groups, old cross history, mode16 and native entropy PASS; Windows/amd64 CGO-disabled cross-build PASS (compile only).
- Reconstructed old R57 CLI encoded the original complete **mr** corpus file to exactly **2,729,261 B**, matching recorded approved R57. New R60 CLI restored its bytes/digest exactly. New R60 CLI encoded mr to **2,673,111 B**. Old R57 decoder correctly rejected the new 512 KiB frame with cross block bounds error. Older archives are backward-readable, newer extended BWT frame archives not forward-readable by old binaries.
- No rejected R59 three-table entropy path or R60 remapping code committed. Only narrow production BWT and mixed-cross boundary/decode changes, plus tests.

## Original provisional conclusion (superseded by user's rejection)

**Correction, 2026-10-19:** The user rejected this performance trade despite the genuine 2.074% archive reduction. Encoding worsened by approximately 0.521 s and decoding by 0.150 s relative to the historical ten-run R57 measurements, equivalent to 0.593 percentage points of additional reduction from the original corpus. The **entire R60 source change was reverted** and **R57 is restored as approved.** New binding requirement: absent explicit user permission, a candidate cannot be accepted unless complete Silesia compressed size, encoding time and decoding time all preserve or improve the approved baseline. Additional R61 attempts to salvage the improvement with a packed LF inverse and bounded BWT windows were unsuccessful; see [Round61 recovery report](2026-10-19-round61-r60-pareto-recovery.md).

## Historical conclusion, void

**ACCEPT H1, Round60.** The 2.074% compressed-size improvement is materially larger than previous 0.2–0.35% modeling gains, affects only selected BWT input regions, and has no per-file size regressions. Nominal ~9.5% encoding and ~6.8% decoding costs, bounded memory and newer-decoder requirement are documented. The user's strict compression-priority rule is preserved, while overhead is proportionate to this substantive size gain.

**Anti-spiral boundary:** Mark the R59–60 predecessor-state remapping family exhausted at this measured design scale (best 9 KB oracle); do not run another ten rounds of map/threshold swaps. Also stop doubling BWT size after testing 512 KiB/1 MiB/2 MiB: larger variants trade too much decoding CPU. Next independent theory should target the remaining mathematically selected source streams in mozilla/ooffice/samba or genuinely different cross-field reversible mechanisms, not further BWT micro-model tweaking.