# ACCEPT: field-period-aware native lane coding with faster Huffman decode

**2026-10-20. Status: ACCEPTED on `main` source commit `37f02082063fed3bb42aa05817dd0afabd0790ae`.** Exactly two existing user commands; one automatic native-Go compressor; no foreign codec backend. Original root `default.pgo` retained unchanged (SHA-256 `3b35cda74832360a252b63c1805e025aa1f4ed7a4fe4124d9a55929a6e86289f`).

## Source-first diagnosis and distinct hypothesis

The accepted 12-file Silesia archive (211,938,580 B input) previously contained **6,261,713 B** for `sao`, a floating/binary-record file automatically encoded in independently packed byte lanes. Its actual content displays repeated fields at a **28-byte source period**, with approximately **20.9% equality at a 28-byte lag** on the original complete file; the generic accepted two/four/eight-lane encoder conflated these highly nonuniform field positions. On real complete `sao`, field entropy averaged **5.533 bits per input byte** when stratified by position mod28, versus **6.926 bits** by mod8. These measurements expose substantial unused structural redundancy, not a generic claim that all source files are periodic.

Previously tested negative techniques were **spatial prediction of x-ray 16-bit words (width 1900)**, **XOR lag8 before sao byte lanes** (expanded), **context-conditioned opcode/distance models**, and **fixed enlarged BWT block windows**. They did **not** select a short inferred source-record period and feed a field-aligned layout directly into JLZ's *existing independently Huffman-packed lanes*. This candidate changes the meaning of field group, not just a cutoff of any failed predictor; the decoder uses already existing lane interleaving and native framed entropy streams.

Independent upstream architectural inspiration: Zstd's distinct literal/sequence encoding, LZFSE's separate symbol classes, and libdeflate's fast Huffman code-table decode, but the concrete encoder/detector/decoder are independent JLZ Go source and do not call any other codec.

## Actual implementation and automatic selection

- `engine/lane.go`: `recordAlignedStride` checks short lags 12–64 on three distributed source samples, requiring a conspicuous repeated-byte coincidence and preferring the shortest period near the strongest signal to avoid unnecessarily duplicating entropy tables. It runs only after JLZ's existing cheap initial evidence makes a lane archive plausible. For a detected record period, **actually encode those three samples using the existing native lane stream packer**; choose it only when its fully framed sampled output beats already tested 2/4/8 lanes. The **complete original archive still competes with the existing mathematical, BWT, and other modes using exact on-disk bytes**. Never selects a user-controlled level, hardcoded filename, or uncharged ideal entropy.
- `lane/block.go`: existing mode13 extended to admit stride values 12..64 in addition to backward-compatible 2/4/8. Every lane remains independently framed; old lane streams decode unchanged and each decoded lane length is derived from the already bounded original block size. Actual `sao` uses **28**; older binaries reject the new stride (forward compatibility not promised), current decoder retains all previous JLZ1 modes.
- `engine/huffman.go`: preserve the exact existing canonical Huffman stream bytes while allocating the code lookup for each actual maximum code length instead of a fixed 15-bit table; use a single 32-bit big-endian word read for common interior code lookups, retaining strict tail bounds, final padding and invalid-prefix rejection. This is a supplementary decoder acceleration, not a different Huffman codec or new archive mode.
- `engine/lane_record_test.go`: arbitrary short/full-sized exact inverse, truncation and unsupported stride rejection; 28-byte synthetic record detection and random false-positive controls. Accepted compact match buckets, two-hop BWT inverse, mathematical tokens, BCJ→BWT and root PGO are unchanged.

## Original component and complete integrated results

Complete original 28-block `sao` components: **existing 8-lane encoding 6,317,645 B**; new 28-lane native `packStream` payloads **4,941,467 B**, saving **1,376,178 B / 21.78% at this specific component comparison**, before whole archive metadata. An additional native FSE alternative produced **4,904,394 B** but was not promoted as it entails substantial additional encoding/model work for only 37KB extra on an already large ratio gain; no FSE full-corpus performance claim is made. Every original source block was independently reconstructable with the current native unpacker.

**Complete JLZ archive improvement is different and authoritative:** `sao` prior **6,261,713 B**, current **4,941,704 B** (**-1,320,009 B**). The other eleven original Silesia compressed file lengths remained exactly the previous accepted values.

| Silesia original | Previous accepted archive B | Candidate archive B | Delta B |
|---|---:|---:|---:|
| dickens | 3,155,856 | 3,155,856 | 0 |
| mozilla | 18,609,139 | 18,609,139 | 0 |
| mr | 2,729,261 | 2,729,261 | 0 |
| nci | 2,163,199 | 2,163,199 | 0 |
| ooffice | 2,692,252 | 2,692,252 | 0 |
| osdb | 3,425,086 | 3,425,086 | 0 |
| reymont | 1,455,234 | 1,455,234 | 0 |
| samba | 5,134,262 | 5,134,262 | 0 |
| sao | 6,261,713 | 4,941,704 | **-1,320,009** |
| webster | 9,835,416 | 9,835,416 | 0 |
| xml | 515,852 | 515,852 | 0 |
| x-ray | 4,266,987 | 4,266,987 | 0 |
| **TOTAL** | **60,244,257** | **58,924,248** | **-1,320,009** |

## Ten completed full original-corpus candidate-only runs: final integrated code

Built in Go1.23.2 linux/amd64 with unchanged root `default.pgo` explicitly supplied; full 12-file complete original 211,938,580-byte corpus on each pass. The final winning combination of record detection/extended native lane and original format-preserving **sized Huffman lookup plus fast interior 32-bit loads** is compared to the **retained previous accepted medians**; none of the current JLZ reference, Zstd-3 or XZ-6 were rerun. Every completed pass decoded the entire original corpus, with byte equality and source SHA matching: **120/120 exact**. Tool invocation was split into batches of four, three and three isolated test-process passes; no incomplete run was included.

| Candidate completed pass | Whole archived bytes | Encode seconds | Decode seconds |
|---:|---:|---:|---:|
| 1 | 58,924,248 | 5.831335404 | 2.054662042 |
| 2 | 58,924,248 | 5.562619156 | 1.922131808 |
| 3 | 58,924,248 | 4.852588973 | 1.870994665 |
| 4 | 58,924,248 | 5.338811776 | 1.922806426 |
| 5 | 58,924,248 | 4.828231340 | 1.881792860 |
| 6 | 58,924,248 | 4.681010659 | 1.788620338 |
| 7 | 58,924,248 | 4.922573220 | 1.827445142 |
| 8 | 58,924,248 | 4.856798023 | 1.873169938 |
| 9 | 58,924,248 | 4.724745258 | 1.830161833 |
| 10 | 58,924,248 | 4.934886558 | 1.879776782 |
| **Independent medians** | **58,924,248** | **4.8896856215** | **1.876473360** |

Retained prior accepted (historical, **not rerun**): **60,244,257 B / 5.081535442 s encode / 1.9046838445 s decode**. Candidate relative to accepted **-2.191095% compressed bytes**, **-3.775430% encode time**, **-1.481111% decode time**. Original reduction becomes **72.1974885%**, up from **71.574662%** (+0.6228265 percentage points). Retained external references: Zstd-3 **66,216,569 B / 1.2561 / 0.8551 s** and XZ-6 **49,233,340 B / 116.9892 / 2.3079 s**, whose external CLI times are not directly comparable to JLZ in-memory times.

**Repeatability caution:** candidate-only ten-run sessions were independently timed, unpaired to the historical accepted run; a preceding record-only sample 10 had size 58,924,248 B / 4.874988 / 1.898633 s, but a source-guarded pass of ten had 5.091147 /1.918466s, and a subsequent sized-Huffman-only full ten had 5.032386 /1.912411s. These modest differences reflect host/session noise. The **final** 10-pass series above is the authoritative passing source variant, not a claim that either speed improvement will replicate identically on every machine. Separately sampled old vs new 28-lane decoder had modest, variable relative throughput; the final low-overhead Huffman interior-word path measured ~175.6MB/s decoding the selected `sao` 256KiB block vs ~156.3MB/s with prior fixed lookup in one focused diagnostic. Single component timings are not full codec medians.

## Full safety, build and compatibility checks

- **All twelve canonical original Silesia SHA-256 values PASS**, size 211,938,580 B. **120/120 candidate full original-file byte and digest matches** over ten runs, with exact file archive sizes disclosed above.
- **Uncached `go test -count=1 ./...` PASS**, **`go vet ./...` PASS**, targeted **`go test -race -count=1 ./engine ./bwt -run 'TestRecord|Test.*Huffman|TestTwoHopInverse'` PASS**, **CGO=0 Windows/amd64 cross-build PASS** (not native Windows runtime).
- Old previously produced and SHA-valid **mathematical fixture** JLZ archive decoded unchanged, and an actual **complete `mr` BWT archive** encoded by the previous accepted build decoded to the exact original `mr` file. This preserves backward decoding. Old binaries are not guaranteed to understand widened lane strides.
- The four production-source/testing files were committed atomically to GitHub `main` as exact source-aligned Git blobs matching locally tested snapshots. Source formatting and explanatory comments after the measured candidate do not alter executable semantics; uncached tests were rerun after formatting. Root default.pgo was not updated, source/PGo not replaced wholesale, and no experimental FSE or geometry tuner landed.
- Total extra model/side-information costs of the new 28 individually packed lanes are included in the reported final **4,941,704 B `sao` archive**, not hypothetical entropy estimates.

## Decision and next architecture

**ACCEPT H1:** Strong net whole-codec archive decrease with encoding and decoding candidate-only ten-run medians both below the accepted references, exact reconstruction, validated compatibility, and additive preservation of existing JLZ architecture. Next pursue similarly **data-supported structure-aware native representations** in dominant remaining math-selected `mozilla` and BWT-selected `webster`, rather than tuning periods 27/29, switching FSE models or adjusting detector thresholds. Avoid claiming impossible universal 99% reduction: source entropy is real and data-dependent. No performance levels or new command paths introduced.
