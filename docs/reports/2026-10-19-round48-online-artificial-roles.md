# Round 48 — Artificial byte meaning: decoder-synchronous roles, temporal relations, Huffman, arithmetic and native rANS

**Date:** 2026-10-19. **Verdict: REJECT H1 for all four fully integrated codec configurations; real component-level predictive redundancy confirmed and multiple next-generation ideas remain OPEN.** No production Go files promoted. **Round 44** remains approved. No external compressor or GitHub Actions; temporary trials were built from byte-verified approved R44 or the discarded R47 prototype. The prior R47 static learned-map result is retained as a comparison, not re-benchmarked. All measurement input consists exclusively of **all twelve complete original Silesia files** (211,938,580 B; 12/12 canonical SHA-256 matches), divided into actual source and already selected native streams as expressly identified below.

## 1. What “artificial meaning” can and cannot do

A byte role is a decoder-reconstructible state/class `R_i=f(X_{<i},model)` that predicts or routes the next exact byte/symbol. Examples are neighboring-value relationships, online learned modal successors, shared clusters of predecessors, or a context-conditioned codebook. The encoder can use source-dependent training **only if it serializes the necessary model**, or learns deterministically and synchronously at both ends from already reconstructed bytes. The exact symbol identity is never lost.

`H(X|h)=H(R|h)+H(X|R,h)` for a class `R` determined by `X`; arbitrary relabeling cannot manufacture Shannon information. An advantage must come from estimated/model redundancy, state sharing, cheaper codebooks, conditioning and finite-coder overhead, with all framing/time counted. In particular, hard partitions of `X_{i-1}` are not necessarily semantically novel and may outperform a costly learned grouping without disproving predictive equivalence.

Useful adjacent algorithms: PPM and context-mixing across history functions (Mahoney, 2005: https://repository.fit.edu/ces_faculty/166/), context-tree weighting and context quantization, minimum-description-length model selection, causal/predictive-state merging. These are research inspirations, **not JLZ benchmark claims** or assertions that the underlying mathematics is novel.

## 2. Experiments constructed: 25 online transforms and four coding architectures

**Twenty behavior/role modes were executed with real reversible encoders/decoders**, including:
- Four/eight channels selected by predecessor high/low bits (control, no model).
- Four/eight channels selected by a predecessor's **last seen successor** or **online majority successor** (Boyer–Moore-like bounded saturating counter), split on successor high/low bits. Prediction updates strictly after each reconstructed byte. No model transmitted.
- Predictive roles from predecessor-minus-last-successor in modular byte arithmetic.
- Cross-temporal roles from previous-two-byte modular difference, XOR, mixed affine hash or an explicit repeat/non-repeat state, in four/eight channel layouts.
- Five direct invertible one-stream residuals: last-successor XOR, last-successor subtraction, zero/successor swapping, and two successive extrapolation-based XOR/subtraction predictors.

Every split channel is individually native Huffman/raw/RLE-packed with length and model-mode headers. The decoder interleaves channels using the same previously decoded byte history. No role labels are transmitted per byte. Seven of the modes were also evaluated as **standalone independent original-byte compressors**; remaining variants were screened on all actual selected BWT and mathematical metadata streams, rather than artificially replacing selected archives.

Alternative entropy architectures with exact decoder implementations:
1. **Separate native packed role streams** (4/8 groups), with byte- and time-charged headers.
2. **Single interleaved role-conditioned canonical Huffman code**: 4/8 sparse codebooks transmitted, one bitstream, group chosen from previous decoded bytes; eliminates independently buffered packed channels.
3. **Decoder-synchronized adaptive bitwise arithmetic coder**: role-conditioned prefix counts plus global prefix counts; no transmitted model, actual coded bits measured, but per-bit probability divisions and updates.
4. **Native static role-conditioned 12-bit rANS**: role-conditioned actual block histograms normalized to sum=4096, serialized sparse count/frequency tables, backward native rANS encoder and forward decoder using decoder-reconstructible roles, no external library. This avoids per-bit arithmetic updates and transmits class-conditioned codebooks. Exact model/bitstream/frame sizes measured.

All modes tested are **fully reversible**, not just compressed entropy estimates; targeted `go test -race` exactness tests passed for online splits, interleaved Huffman, role arithmetic and native rANS.

## 3. Isolated on original source bytes: entire original Silesia

Compared native `packStream(raw-region)` to seven actual reversible online-role transforms, **on every consecutive independent original 256 KiB block of all 12 complete files** (including final partial blocks). Each alternative carries its own stream framing, source length and chosen model indicator; there are no global full-file JLZ1 checksum/framing bytes in these **component totals**. The long-running test was divided after nine complete source files; the final three complete source files were executed independently without omitting any original bytes. The byte sum was verified as **211,938,580 B**.

- Baseline native-packed raw regions: **133,318,073 B**.
- **Best of seven online reversible role choices per region:** **120,091,407 B**, a selective **13,226,666 B (9.921%)** saving against *raw native packing*; this is not a complete JLZ archive or direct savings versus R44.
- Historical R47 static learned eight-group raw native variant: **115,781,631 B** (lower still than the new seven-way online selector). Thus removing role-table transmission does not automatically improve compression.
- The first nine original files and independently completed `webster`, `x-ray`, `xml` all supplied exact decoded results. No estimate from an incomplete process was included.

## 4. Supplemental component tests on actual Round44 winners

### 4.1 Selected BWT/bijective MTF symbol streams

Read all six BWT-selected *complete* original Silesia inputs, and independently transformed/encoded every **412** original source regions through native BWT and bijective MTF. Compared the **actual native packed symbol streams**, not original-byte entropy estimates. Baseline **20,533,358 B**.

- Online split **low4**: **20,291,500 B (-241,858)**; last-successor/majority successor high/low and cross-temporal roles generally provide smaller gains; *best of twenty role-split plus five residual modes at each of 412 streams* **20,123,250 B (-410,108)** (selected-only optimistic composite).
- Single interleaved role Huffman **low4**: **20,281,067 B (-252,291)**; four-way explicit repeat-state variant -239,847 B. Best of four interleaved-mode component choices: **20,171,758 B (-361,600)**.
- Native role rANS **repeat4**: **20,195,484 B (-337,874)**, `low4` **20,218,113 B (-315,245)**. Three variants together required 2.014 s role rANS *encoding* and 1.719 s role rANS *decoding* for the 412 streams (component timing only; native BWT transforms and other JLZ paths excluded). Faster than the arithmetic candidate but not automatically faster than the full JLZ encoder.
- Decoder-synchronous role arithmetic **low4**: **19,305,854 B (-1,227,504)**; `major4` 19,306,868 B. Both bitwise variants together spent 13.701 s *encoding* and 14.257 s *decoding* the selected BWT symbol streams; **~6.851 / 7.129 s per configuration**, before any BWT, metadata or other JLZ work. Not practical for speed-first JLZ despite the strong size gain.

### 4.2 Selected mathematical metadata, actual archived stream extraction

Using complete untimed encode/decode of the four currently mathematical-region-winning original Silesia files (`mozilla,ooffice,osdb,samba`), parsed **676 actual packed opcode/length streams** from their existing cross-region archives. Verified each entire source archive restored exactly, and separately that each re-encoded role stream has an exact inverse. Existing packed cost: **4,506,404 B**.

- Online split **majority-successor low8**: **4,285,805 B (-220,599)**; best of 25 reversible candidates per stream **4,063,901 B (-442,503)**.
- Single interleaved Huffman **majority-successor low8**: **4,268,976 B (-237,428)**; best of four variants per stream **4,081,028 B (-425,376)**.
- Native static role rANS **low4**: **4,397,867 B (-108,537)**, with best of three per stream **4,107,380 B (-399,024)**.
- Adaptive role arithmetic **low4**: **4,046,007 B (-460,397)** (the two arithmetic modes together cost 3.061 s encoding/3.147 s decoding these *metadata streams alone*; ~1.53 s per variant).

All component comparisons explicitly count model/stream framing, but are **not complete JLZ1 archive results**, may include BWT sample-block encoder differences, and cannot be directly summed to infer an identical net full codec improvement. Matching old decoder paths requires explicit new stream subtypes.

## 5. Four native JLZ integrations, complete 12-file screens (no ten-run promotion)

For the three new codec families, started from the already validated discarded Round47 integrated native stream subtype. Replaced its static learned-role map with the exact online/single-Huffman/rANS encoders and corresponding native decoders, targeting only **selected BWT symbol** and **actual math opcode/length** streams, with automatic *fully framed* shorter-of-existing-vs-candidate selection. No user flags or existing compressor logic replaced. Legacy codecs remained readable. Successful encoder/decoder candidate archives need a newer decoder for their new packed subtype. **Each of the four configurations was independently compiled and ran a full 12-file Silesia end-to-end encode/decode once, including full archive metadata/digests, verifying original byte equality on all 12 files:**

- **Fixed 4-way low-bit split** on BWT and math: **60,912,207 B**, **6.477577 s compression / 2.757289 s decompression**. vs R44 **-506,004 B (-0.824%)**, nominal **+8.95% encode/+12.76% decode**.
- **Hybrid online split** (`low4` BWT; `majorLow8` mathematical): **60,839,930 B**, **6.437793 / 2.768255 s**. vs R44 **-578,281 B (-0.942%)**, nominal **+8.28%/+13.21%**.
- **Single interleaved role-Huffman** (`low4` BWT; `majorLow8` mathematical): **60,821,703 B**, **7.093449 / 2.681969 s**. vs R44 **-596,508 B (-0.971%)**, nominal **+19.31%/+9.68%**.
- **Role-conditioned native rANS** (`repeat4` BWT; `low4` mathematical): **60,844,977 B**, **6.587936 / 2.821082 s**. vs R44 **-573,234 B (-0.933%)**, nominal **+10.80%/+15.37%**.

All **12 original files reconstructed exactly in every complete integrated screen**, 48/48 file-roundtrips. No source/file archive-size regression in any of the four screen configurations. The two-user-facing compressor operations remain unchanged but newer packed modes cannot be read by the existing R44 decoder; that compatibility cost is another reason not to promote.

**Critical timing caveat:** these are each **ONE newly measured complete original Silesia screening pass**, NOT a candidate ten-run median. Percentages compare against historical R44 ten-run medians from another session and are nominal/unpaired, so they cannot establish exact causal timing penalties. However all four complete encodes/decode screens are materially slower than the retained reference and show **less compression benefit than the previously rejected R47 static role candidate (60,694,806 B / 6.186603 s encode / 2.579713 s decode historical ten-run median)**. No proportional promotion benefit, so ten repeats per rejected mode would waste measurement compute. No Zstd/XZ/approved JLZ reference reruns.

### Required comparison

| Model | Complete compressed archive | Reduction from original | Compression time | Decompression time |
|---|---:|---:|---:|---:|
| Zstd-3 (historical external CLI) | 66,216,569 B | 68.76% | 1.2561 s | 0.8551 s |
| XZ-6 (historical external CLI) | 49,233,340 B | 76.77% | 116.9892 s | 2.3079 s |
| **Approved JLZ Round44 (historical 10-run median)** | **61,418,211 B** | **71.02075%** | **5.945576 s** | **2.445270 s** |
| **Candidate JLZ R48 interleaved Huffman (one full-corpus screen)** | **60,821,703 B** | **71.30218%** | **7.093449 s** | **2.681969 s** |

| Archive-size difference vs R44 | Encode-time difference vs R44 | Decode-time difference vs R44 |
|---:|---:|---:|
| **-0.971%**, true newly emitted bytes | **+19.31% nominal**, unpaired | **+9.68% nominal**, unpaired |

The best-compressing R48 integration is chosen for this comparison; other distinct integrated modes are separately enumerated above with their measured numbers. We deliberately do **not** invent candidate ten-run results or present component savings as integrated wins.

## 6. Validation, rejection rationale and continuing research queue

- Original Silesia 12/12 `sha256sum -c docs/reports/2026-10-18-round23-silesia-sha256.txt` **PASS**.
- All isolated raw regions, 412 BWT streams, 676 selected mathematical metadata streams in tested modes **PASS exact reversible reconstruction**. Complete integrated native archive screens pass exact 48/48 complete-file reconstructions.
- Prototype `go test -race -count=1 ./engine -run 'Test(OnlineRolesExact|InterleavedRoleExact|RoleArithmeticExact|RoleANSExact)'` **PASS**. Three codec split/Huffman candidate `go test -count=1 ./...` + `go vet ./...` **PASS**; native rANS integrated candidate full Go/vet, targeted race and Windows/amd64 cross-build **PASS**.
- **Approved Round44 source untouched:** exact Git hashes of `engine/bwt.go`, `engine/cross.go`, `engine/match.go`, `engine/stream.go`, `bwt/zero_bijective.go` match upstream R44 blobs. Uncached `go test -count=1 ./...`, `go vet ./...` and Windows cross-build **PASS** on that unchanged source. No GitHub Actions, no temporary branches, no promoted candidate files.

**REJECT H1** as a practical source promotion. All new models yield correct compression and role-conditional structure, but no integrated candidate beats the **size/speed/complexity combination** of current accepted R44, and all also underperform the archived R47 rejected static variant. Do not tune isolated group counts/bit hashes further.

**Research directions retained OPEN (distinct architecture, not parameter-churn):**
1. **Fast context-conditional native rANS/range coding with shared models:** conditional coder directly exploits the measured 1.23 MB BWT role arithmetic gain, but needs an *order-of-magnitude* cut in adaptive per-bit modeling and frame overhead. Investigate fixed-point context mixtures, table lookup and multi-symbol ANS/faster inline decode, while charging models and exact residual bytes. Current role rANS gains less compression but substantially lowers coder work.
2. **Bayesian online role mixing / soft overlapping predictive states:** learn mixtures of role-conditioned context experts from already decoded history rather than forcing each predecessor into one class. Must encode with an actual decoder-synchronized native coder and honor per-bit context maintenance and CPU.
3. **Hierarchical role states with suffix/backoff and bounded context-tree weighting:** variable-order model shared across already selected JLZ metadata, not generic standalone neural-model weights. Preserve zero-overhead fallback for weak contexts.
4. **Factorized byte coordinate systems:** small invertible `byte↔(role,exact-index)` or group action under arbitrary context; role mapping must be free to compute from decoded state or explicitly serialized, and conditional residual identity must be accounted for. Static symbol relabeling alone cannot change ideal entropy.
5. **Graph/hypergraph predictive-role merging:** cluster contexts based on conditional outcome distributions or learned interactions, with compact state updates and a strict model-budget ceiling. Distinct from previous spin/graph and offline causal-state screens only if it leverages *actual winning R44 symbols* with charged costs.
6. **Predictor ensembles & online self-organizing automata:** decoder-synchronous mixture of sparse functions of history (lags, operations, byte classes) rather than expensive hard-coded per-byte predictors. Require fast SIMD-/table-friendly evaluation and measured selected-stream byte advantages.
7. **Computational optimization independent of model quality:** encode only the chosen representation rather than materializing both fully packed originals and candidate; pre-estimate packed length safely; avoid extra per-stream allocations, pooling tables, and repeated full frequency sweeps. Real benchmark evidence required, not theoretical speed assumptions.

No part of this finite experimental queue establishes that artificial byte meaning is impossible, or that exact model family is exhausted. Role relationships are *genuinely measurable*; practical near-term blockers are **entropy-code throughput and redundant candidate computation**, not lack of losslessly invertible latent group structures. Next experiments should test a fundamentally cheaper **native conditional coding kernel or model-share strategy**, not another tweak to online group count.
