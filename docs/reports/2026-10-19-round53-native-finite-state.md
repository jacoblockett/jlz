# Round 53 — Native normalized finite-state coding on selected JLZ streams (ACCEPT)

**Date:** 2026-10-19. **Decision: ACCEPT H1.** Approved production is now **Round 53**, replacing Round 51. No external compressor backend, copied compressor code, hosted CI, alternate user-selectable level, or additional CLI verb. The design is an independently written **byte-renormalized rANS-style finite-state entropy coder**, **not** a port/reproduction of Zstd's FSE/tANS implementation. Arithmetic normalization and reversible state equations are conventional entropy-coding mathematics, not claimed novel.

## Research question and implementation

After auditing Zstd level 3 in Round 52, test whether cheap finite-state coding supplements JLZ's *already selected* mathematical metadata and BWT bijective MTF/zero-run symbol streams while retaining the successful baseline codec. Implemented independent streams with fixed total frequency **4096 (12 bits)**, a lower state bound **2^23**, a deterministic sparse normalized count model, reverse-order state transitions, and forward exact restoration. State renormalizes through ordinary bytes, with a final four-byte initial decoder state and a sparse per-symbol normalized-frequency table. All model, stream framing, and exact reconstruction data are carried in the archive. Codec subtype **5** is a new native packed-stream mode; existing modes **0–4** are decoded unchanged.

The new internal automatic \`packFiniteState\` computes the old \`packStream\` result and the exact fully framed native finite-state alternative, selecting the latter *only* when shorter. It is called for the **selected BWT bijective symbol stream**, and mathematical metadata opcode/length channels in current region/cross-history/split-length packing. Other literals, distances, entropy and source paths retain their previous encoders. Neither block ordering nor worker scheduling was changed; each model belongs to its own stream, retaining current independently decodable regions and region groups. A new decoder is required to read newly encoded codec-5 streams, but it still reads previous JLZ1 archives.

The native Go implementation adds \`engine/finite_state.go\`, focused tests \`engine/finite_state_test.go\`, and limited integration points in \`engine/{bwt,packed,lengths,cross,stream}.go\`. No implementation or codec delegates compression to DEFLATE/Zstd/XZ/Brotli/LZMA, and no external algorithm is imported as a backend.

## Component research and rejected table reuse

**676 actual archived mathematical opcode/length streams** from selected mathematical-region blocks of four complete original Silesia source files. Existing fully framed native packing **4,506,404 B**. Newly written finite-state coding with dynamic per-stream tables and automatic shorter-of-existing selection **4,452,019 B**, **54,385 B** component gain; exact reconstruction on every stream. (An earlier diagnostic used a different extra subtype-marker charge and reported 54,803 B; the 54,385 B figure uses the *actual* native framing definition and is authoritative for this component comparison.)

**412 selected BWT bijective symbol streams** from all six full BWT-selected original Silesia files. Direct native packing **20,533,358 B**. Independently framed new finite-state streams **19,991,247 B**, **542,111 B** saved across 412/412 blocks. Another preliminary fixed-overhead diagnostic charged an extra marker and reported 541,699 B. All selected symbols restored byte-exactly.

**Alternative: reuse previous entropy tables across group-local blocks** instead of serializing each table. Implemented and tested real reversible use of preceding per-plane models within the four-block decoding groups, permitting reuse only when every next symbol has nonzero frequency, and fully accounting for actual framed lengths. It added merely **4,375 B BWT** and **940 B math** savings over dynamic model + old stream selection, on **16/412 BWT** and **35/676 math** reused streams respectively. **REJECT table reuse**: it introduces dependent decoding/model lifecycle and interferes with parallel region work for a negligible **5,315 B** component opportunity, before implementation metadata/coordination costs. The approved coder transmits independent models.

## Complete, fully framed integrations and run count

A *single*, genuinely integrated **BWT-only candidate screen** completed all twelve original Silesia inputs with exact restoration: **60,896,887 B / 5.450074 s encode / 2.278436 s decode**. This is ONE candidate-only whole-corpus screen, NOT a ten-run result, and its timing cannot support a precise paired speed conclusion.

The combined **BWT plus selected mathematical metadata** integrated candidate then produced **60,838,336 B / 5.578464 s / 2.308839 s** in a separate ONE-run full screen. It saves **58,551 B** beyond BWT-only within complete framed files; different from summing isolated component oracles because the actual selected encoder paths and sample blocks differ.

For promotion, built the integrated native Go candidate and executed **TEN** candidate-only full twelve-file Silesia archive compression and decompression passes (one compiled test executable, one complete pass per process, original source bytes loaded outside the timed compression calls). **All 120 restored full files match the originals byte-exactly; every pass produced 60,838,336 B**:

| Run | Complete archive bytes | Encode s | Decode s |
|---:|---:|---:|---:|
| 1 | 60,838,336 | 5.735738045 | 2.258961862 |
| 2 | 60,838,336 | 5.707590649 | 2.307839829 |
| 3 | 60,838,336 | 5.532649875 | 2.367133888 |
| 4 | 60,838,336 | 5.614807516 | 2.274068535 |
| 5 | 60,838,336 | 5.522138486 | 2.262830760 |
| 6 | 60,838,336 | 5.899333927 | 2.352364994 |
| 7 | 60,838,336 | 5.666789740 | 2.341502497 |
| 8 | 60,838,336 | 5.835198384 | 2.395697008 |
| 9 | 60,838,336 | 6.020690915 | 2.371981119 |
| 10 | 60,838,336 | 6.318792039 | 2.532579104 |
| **Median (separate column medians)** | **60,838,336** | **5.721664347** | **2.346933746** |

Run wrapper for the 5–10 batch was interrupted after four complete processes; completed runs 9–10 separately. None of the interrupted/incomplete run output is counted. The ten displayed results are complete and verified. Results are candidate-only; the current/previous reference was NOT benchmarked for speed again.

### Required reference comparison (211,938,580 original bytes)

| Model | Complete archived bytes | Reduction vs original | Compress s | Decompress s |
|---|---:|---:|---:|---:|
| Zstd-3, **historical external CLI** | 66,216,569 | 68.76% | 1.2561 | 0.8551 |
| XZ-6, **historical external CLI** | 49,233,340 | 76.77% | 116.9892 | 2.3079 |
| Current JLZ **R51**, **retained prior ten-run median** | 61,418,211 | 71.02075% | 5.902140113 | 2.414516970 |
| Candidate JLZ **R53**, **NEW ten-run median** | **60,838,336** | **71.29436%** | **5.721664347** | **2.346933746** |

| Candidate size delta vs R51 | Encode-time delta vs R51 | Decode-time delta vs R51 |
|---:|---:|---:|
| **-0.94414%** (**-579,875 B**) | **-3.0578% nominal** | **-2.7990% nominal** |

**Timing provenance:** the previous accepted Round51 ten-run median is an earlier, independently measured session, not a paired control. Encode/decode changes are therefore directional, **not** precise causal estimates; frequency, GC, OS scheduling and memory can affect them. Historical Zstd/XZ external CLI timings use a different I/O convention than JLZ's in-memory calls and are not directly comparable.

**Per-file archive changes versus R51**, all including the original archive metadata/digest:

| File | Previously approved R51 B | New R53 B | Delta B |
|---|---:|---:|---:|
| dickens | 3,307,328 | 3,213,620 | -93,708 |
| mozilla | 18,634,854 | 18,609,139 | -25,715 |
| mr | 2,859,487 | 2,831,212 | -28,275 |
| nci | 2,193,636 | 2,169,196 | -24,440 |
| ooffice | 3,037,461 | 3,033,712 | -3,749 |
| osdb | 3,487,379 | 3,472,004 | -15,375 |
| reymont | 1,514,459 | 1,469,844 | -44,615 |
| samba | 5,147,974 | 5,134,262 | -13,712 |
| sao | 6,261,713 | 6,261,713 | 0 |
| webster | 10,170,379 | 9,851,849 | -318,530 |
| xml | 536,554 | 524,798 | -11,756 |
| x-ray | 4,266,987 | 4,266,987 | 0 |
| **TOTAL** | **61,418,211** | **60,838,336** | **-579,875** |

**Worst per-file size regression 0 B**; ten smaller, two unchanged. Biggest single improvement webster -318,530 B (-3.13%). Reused previous approved R51 per-file archive measurements without rerunning baseline speed.

## Validation and compatibility

- **12/12** canonical source Silesia SHA-256 checks PASS (corpus total 211,938,580 B).
- **120/120** full-corpus complete file reconstructions byte-exact across ten candidate-only runs.
- Native finite-state model-table encode/decode, normalized-state inverse, full packed-stream framing, small/large deterministic fixtures and malformed-model/truncated-input rejection PASS (including targeted race).
- \`go test -count=1 ./...\` PASS and \`go vet ./...\` PASS on final candidate.
- \`go test -race -count=1 ./engine -run '^TestFiniteState'\` PASS.
- Windows \`CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build\` PASS; Windows runtime was not tested.
- **Backward decoding:** compile unchanged R51 CLI, encode bundled original 126,479-byte Go source fixture; R53 CLI decoded it with SHA-256 verification and a byte-for-byte \`cmp\` PASS. The existing project compatibility tests also pass. Earlier decoders **cannot decode new codec-5 subtype streams** (forward wire incompatibility documented). Previously rejected experimental R47 codec 5 is NOT a production format and was never approved.
- No GitHub Actions, no external benchmarking, no third-party compression backend, no external implementation copied.

## Decision

**ACCEPT H1**: a compact additive native finite-state entropy alternative improves **complete original-Silesia archive size by 579,875 B (0.944%)**, with directionally faster candidate encoding and decoding against retained R51 measurements. All files improve or stay constant and full restorations pass. It does introduce a small per-stream normalized model, one new decoder subtype, and memory for a 4,096-byte symbol-lookup table; those costs are justified by the complete-benchmark gain. Unlike previous R47/48 context/role-based entropy schemes, this simple **unconditional, independent** model performs competitively on both ratio and CPU.

**Production R53 supersedes R51** after native tests; retain R51 historical reference, all old decoders in new source, old R44 BWT modes, and the two-command CLI. Keep small table reuse **REJECTED**, and consider the distinct Zstd-inspired cheap exact-only match parser next only with convincing experimental evidence. Do not endlessly tune finite-state scale/group size; this first substantial integrated experiment already passes the acceptance test.
