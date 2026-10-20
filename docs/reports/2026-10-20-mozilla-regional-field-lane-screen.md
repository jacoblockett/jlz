# Mozilla structured-field lanes within selected mathematical regions: integrated tradeoff rejected

**2026-10-20. Status: REJECT H1 and H2 for production.** The approved JLZ on `main` is unchanged at **58,924,248 B** on the **complete original 12-file Silesia corpus**, historical 10-run encode median **4.8896856215 s** and decode median **1.876473360 s**. Root `default.pgo` remains unchanged. All new code described below was tested **only in disposable copies** of the source. This is not a standalone compressor or a new compression backend.

## Research basis and competitor crosswalk

Analyzed current JLZ region encoder and selected mode13 lanes, confirmed `mozilla` is a **51,220,480 B POSIX TAR archive with 843 entries** including numerous ELF `.so` files; `mozilla` archived output **18,609,139 B** is the largest single-file cost. Source analysis confirms XZ's reversible x86 relative-branch filter and Zstd's long-distance rolling-reference subsystem have distinct placement ideas; native JLZ already implements E8/E9→BWT branch normalization, 128KiB math history, and previously rejected widened 512KiB math history / sparse long anchors. No XZ or Zstd encoder was used or rebenchmarked.

### New H0: accepted BCJ normalization before existing native math tokens

Implemented actual `normalizeBranchRelative` → existing `tokenizeMatchV2` → existing `packedRegionVariantsShared` comparison on **21 distributed genuine 256KiB source blocks** from original `mozilla`, with exact native token decode and inverse BCJ reconstruction. Original current math-token components **2,890,565 B**; BCJ-pretransformed components plus a one-byte extra mode **2,941,760 B** (**+51,195 B, 0/21 wins**). This directly tested a different *placement* than the already accepted BCJ→BWT stage, but confirmed that new math-token placement is counterproductive. There is no corresponding integrated candidate or ten-run result. **CLOSED the concrete mozilla BCJ→math pretransform placement**, not the abstract x86 filter.

### New H1: selectively apply established native record-lane codec to math-region losers

A genuine supplementary mode-13 opportunity exists *inside* `mozilla`'s four local regions where math-match packing currently loses to native bytes or predictive bytes. Parsed the actual accepted complete mozilla archive and measured **every one of its 196 selected original-source 256KiB blocks** against independently packed native 2/4/8/12/16/20/24/28/32/40/48/64 byte-aligned lanes, including every lane's existing packed entropy stream and framing. **4/196 winners, aggregate 156,509 B** better:
- Original block58: selected mode2 **46,223 B**, lane16 saves **12,767 B**.
- Block59: selected mode1 **118,684 B**, lane16 saves **71,822 B**.
- Block60: mode2 **101,422 B**, lane16 saves **686 B**.
- Block85: mode1 **124,730 B**, lane48 saves **71,234 B**.
The specific field structure had a 16-byte detected period and a three-record grouping at width48 was beneficial in the fourth block. Existing `lane.DecodeBlock` already accepts these widths, so candidate would require **no new decoder subtype or archive-format change**, and the original mathematical parser is preserved. Every winning lane stream was inverse-tested byte-for-byte.

Screened all independently selected genuine original source blocks of **samba** (83), **osdb** (39), and **ooffice** (24) with the same width set against their incumbent native modes: **0 additional component winners**, so there is no evidence of a useful cross-file field-lane opportunity on those archives.

The first integrated implementation called a cheap `recordAlignedStride` on three distributed 4096-byte samples of each locally math-losing block and compared existing lane widths `p,2p,3p` (at most 64) by exact packed byte costs. The new selector was called during **all** early math sampler passes. It created a genuine **preflight interaction regression**: `sao` switched to an early cross-regional candidate **5,529,896 B**, bypassing the existing globally better 28-lane `sao` archive **4,941,704 B**. This was **a regression solely in the experimental source**; a separate clean re-run of committed accepted source reproduced `sao=4,941,704 B` and **58,924,248 B** total, all 12 exact. Earlier in-chat warning about an accepted-tree problem was promptly corrected by this clean source verification. A temporary further preflight adjustment restored the global winner but added work.

**H1 repaired initial variant**, fully integrated with early per-block field detection plus a guard against prematurely terminating before global lane selection: complete original 12-file archive **58,767,739 B**. Full ten completed candidate-only passes (each exact 12/12) recorded:
- Encoder seconds: **5.223747880, 5.301344786, 5.198887361, 5.494137095, 5.473000133, 4.957743233, 5.148877389, 5.220143073, 5.302299354, 5.189613460**. Median **5.2219454765 s**.
- Decoder seconds: **1.890481923, 1.872354583, 1.888792788, 1.840130007, 1.819808679, 1.832031055, 1.837971082, 1.862372851, 1.850112350, 1.967930795**. Median **1.8562426005 s**.
Relative to historical accepted: **-0.2656105% bytes (better), +6.795117% encode time (worse), -1.078127% decode time (better)**. **REJECT H1** as compression slowdown is not negligible for a ~0.27% archive saving.

### H2: defer field checks until actual final selected math-region encoding

In a second source variant, *reverted both early selector edits entirely* and moved `betterPeriodicLane` into the job loop of `encodeParallelMathRegions`, after the original `encodeMathRegionBlockFast/History` has selected its final regional winner. This preserves all original preflight decisions (including whole-file `sao` acceptance), avoids extra field work on regional samples, and still picks the identical four `mozilla` lanes. One complete candidate screen reconstructed 12/12 and compressed to **58,767,739 B**.

**Ten complete candidate-only original 12-file Silesia passes for H2, all 120 byte/SHA exact:**

| Run | Complete archive B | Encoding s | Decoding s |
|---:|---:|---:|---:|
| 1 | 58,767,739 | 5.423689387 | 1.939083443 |
| 2 | 58,767,739 | 5.143183153 | 1.879192285 |
| 3 | 58,767,739 | 5.246233710 | 1.863499990 |
| 4 | 58,767,739 | 5.129123104 | 1.888843728 |
| 5 | 58,767,739 | 5.285566393 | 1.898059353 |
| 6 | 58,767,739 | 5.144053991 | 1.955843477 |
| 7 | 58,767,739 | 5.134225568 | 1.826433273 |
| 8 | 58,767,739 | 4.787997739 | 1.910896697 |
| 9 | 58,767,739 | 4.769006575 | 1.779580585 |
| 10 | 58,767,739 | 4.693029490 | 1.919169287 |
| **Median** | **58,767,739** | **5.1387043605** | **1.8934515405** |

Accepted historical current (not rerun for benchmarking): **58,924,248 B**, **4.8896856215 s** encoder, **1.876473360 s** decoder. H2 deltas **-156,509 B (-0.2656105%; better)**, **+5.092735% encoder time (worse)**, **+0.904792% decoder time (worse)**. Unpaired session medians can contain host noise; neither H1 nor H2 shows a credible complete-pipeline Pareto speed win. **REJECT H2 production promotion.**

The two new source variants do not change decoder code or accepted archive interpretation. Disposable `go test -count=1 ./...`, `go vet ./...`, Windows/amd64 CGO=0 crossbuild, and focused `go test -race` for the new `betterPeriodicLane` inverse tests PASS. One broader targeted race command stalled and was superseded by the narrower successful check; **do not claim the stalled broad race suite passed**. Production root PGO and old archive decoder paths remained unchanged.

## Conclusions and future placement scope

**REJECT H1 and H2** as integrated automatic compressor changes; do not merge extra periodic math source or the temporary preflight bypass. The source-specific **156,509 B** component opportunity is real and could be useful as an *additive* option if paired with a materially cheaper selector/decoder improvement in a future, different architecture. Do not iterate on adjacent field widths, 27/29 period thresholds or multiply 16 indefinitely. The arithmetic math-literal coder and ELF inter-library references remain more important sources of total mozilla archive bytes.

**Approved `main` remains: 58,924,248 B total (72.1974885% original reduction), 4.8896856215 s encode, 1.876473360 s decode, ten-run historical reference.** This report is evidence, not production source. Neither Zstd nor XZ was rerun.
