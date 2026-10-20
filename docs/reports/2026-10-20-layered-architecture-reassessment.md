# Layered placement reassessment: concrete native JLZ combinations, source and benchmark evidence

**Recorded:** 2026-10-20 (America/New_York). **STATUS: tested several previously omitted additive placements; no production acceptance.** Current accepted JLZ remains unchanged on `main`. No numbered round naming. User directed correction of research process: never declare a concept exhausted merely because an **isolated alternative** or **replacement implementation** performed badly; test whether the idea can complement current original-byte transformations, symbols, residuals, or other established streams before closing a *specific placement*. Prioritize real reversible experiments, but do not claim every conceivable combination can be affordably tested.

**Canonical source:** complete original 12-file Silesia, **211,938,580 B**, all **12/12 canonical SHA-256 PASS**. Source-block component screens used first/middle/last blocks of the actual complete originals, and the BCJ ooffice pairing examined all 24 natural source blocks; no unrelated benchmark corpus used. Reversible exact-data tests and actual native packed stream sizes, with all known per-alternative selector bytes charged. **Component packing comparisons are not whole JLZ archive results.**

Prior research cross-check: [older *actual* complementary integration audit](../research/complementary-integration-audit.md) already established **36 additive-coordinate residual block tests**, selected match-length 16-rule grammar, and selected BWT symbol transforms; many promising sample results were *already additive*, not stand-alone. Previous [strategy ledger](../research/strategy-ledger.md) and [research history](../research/attempted-ideas.md) explain why some synthetic or preexisting-DEFLATE/Zstd results cannot count as JLZ wins. New [placement matrix](../research/layer-placement-matrix.md) explicitly tracks both historical tests and this work with **TESTED**, **SCREENED**, **NOT TESTED**, or **NOT ELIGIBLE** scope.

## Layered concept 1: LZP prediction feeding current entropy streams (two placements, SCREENED negative)

Inspired by libbsc's LZP-before-block-sorting pipeline, implemented an original native Go decoder-synchronous 4-byte contextual predictor. For each byte position, a 65,536-entry local index records the latest prior occurrence sharing the exact four decoded context bytes, predicting the symbol at that earlier occurrence. Encoder XORs this prediction with the input byte; decoder rebuilds its own index **from already restored bytes**, regenerates the same prediction, XORs the residual and restores the original. No external dictionary/model metadata; **one mode byte** charged for a candidate, plus the existing fully packed native stream.

Two independent placements, both actually packed using JLZ's **current** native algorithms and individually inverse-verified:
- **LZP→original-byte BWT→bijective MTF→native 1/2/4-region finite-state coding:** 21 first/middle/final selected BWT source blocks from seven complete original Silesia files. Native incumbent stream **957,151 B**; LZP residual+exact same BWT coder with mode metadata **1,658,244 B**, **+701,093 B**. **0/21 winning blocks**. Exact restore 21/21. Not integrated; BWT suffix-context ordering was harmed by preconditioning.
- **LZP applied ONLY to the actual mathematical tokenizer's literal column**, preserving match commands, lengths, distances, masks and current literal packing: first/middle/final original mozilla/ooffice/samba blocks, 9 stream-level tests. Existing native literal code **600,874 B** vs contextual LZP residual fed to same literal coder **623,726 B**, **+22,852 B**. **0/9 wins**, all exact. Not a whole-file archive or a final selected-math-stream oracle (these blocks were independently tokenized without cross-region history).

**REJECT only those two concrete placements**, not every conceivable context model nor all LZP variants. Stop changing hash size/predictor context lengths solely to tune these losses.

## Layered concept 2: GF(2) algebraic permutations and recurrence INSIDE current mathematical literal streams (two new placements, SCREENED negative)

Historical work tried GF(2) linear relations on raw data, synthetic parity, and some additive residual bitplanes. This work instead operates **after JLZ's mathematical tokenizer has extracted literal bytes**, feeding only that literal column into the *existing* codeword/stream encoder. The math opcode, transformed matches, reference-cache state, decoder grammar and entropy alternatives remain logically preserved (a new decoder-mode flag would be required for full integration).

- **Reversible 8×8 bit-matrix transposition** (an involution, `T(T(data))=data`, no per-block learned matrix, one selector byte) across nine actual original Silesia math block literal columns. Existing packed literals **600,874 B**; bit-transposed + native literal coding **708,012 B**, **+107,138 B** with **1/9** isolated block wins. Every transposition and packed-stream inverse exact.
- **Reversible fixed three-lag GF(2) predictor** on those same extracted literal columns: `residual[i]=literal[i] XOR literal[i-1] XOR literal[i-2] XOR literal[i-3]`; decoder computes the inverse from already decoded prior literals. Existing **600,874 B**; new preconditioner+same literal code **715,562 B**, **+114,688 B**, **0/9** winning blocks. All exact.

**REJECT these two placements**, not general multi-parent causal modeling, GF(2) basis factorization, polynomial relations or the independent additive residual bitplane results already documented historically. No redundant parameter sweeps.

## Layered concept 3: BCJ branch-address normalization before current BWT and predictive paths (full 24-block original ooffice screen)

Earlier native x86 relative branch-address normalization was tested **before mathematical matching**, saving **154,167 B** on the full original Silesia but nominally slowing both compress and decompress; it had been rejected. The new complementary question is whether its transformed **source bytes** improve *JLZ's existing BWT* or *native predictive coder* more than they improve mathematical matching.

Original reversible BCJ input pass: recognize `0xE8`/`0xE9` plus four bytes, convert the signed-wrapping little-endian 32-bit relative displacement to absolute position with `uint32(opcode_position+5)`, then skip the four operand bytes. Decoder reverses every operation and advances identically. This is mathematical address normalization, not an existing-codec backend or an imported xz BCJ implementation.

All **24 original ooffice 256-KiB/last-short natural source blocks**, exact transform inverse and **actual native packed payloads**:
- Current independently encoded math block representation, no group history: **3,036,867 B**.
- Original direct BWT: **3,040,062 B**.
- **BCJ→BWT→existing bijective/rANS native encoding** with submode selector: **2,701,751 B**, **22/24 blocks** smaller than current math.
- BCJ→native raw/RLE/Huffman: **4,925,795 B**, **0/24 wins**.
- BCJ→existing mathematical code (previous positive placement, checked for context): **2,752,577 B**, **21/24 wins**.
- **BCJ→current native predictive coder**: **6,068,483 B** when nil-returning predictor is represented as the raw transformed fallback; the predictor was available on only 1/24 blocks; **0/24 wins**. Every source block transform and non-nil encoded predictive result decoded exactly.

**Interpretation:** BWT and BCJ interact beneficially because normalizing positional operand bytes exposes source-domain reorderability to the existing BWT suffix ordering. This is a genuine *layer synergy*, rather than a new freestanding codec. It **does not prove** savings against JLZ's winning cross-history mathematical blocks until whole-codec integration.

## Fully integrated BCJ→BWT automatic candidate: TEN COMPLETE original Silesia passes (REJECT under strict performance contract)

Built in **disposable original native-Go full JLZ source** `/mnt/data/jlz_layer_bcjbwt`, based on a preserved verified current engine snapshot. When the exact existing BCJ-density gate signals branch-dense original 256-KiB math block, preprocess that block with BCJ and **feed the result to current `encodeBWTSubregions`**, with no duplicate match-entropy search on this selected path. Other blocks and global mathematical/BWT automatic selection remain current; existing 256 KiB frames stay. New EXPERIMENTAL native internal block mode **19** wraps one byte specifying the existing BWT bijective mode15 or 2-/4-segment mode16, followed by its normal packed payload. Decoder uses the **existing BWT inverse** then BCJ inverse; frame length, integrity digest and output bounds unchanged. This is an **additive branch-before-BWT path**, not replacement of JLZ's compressor and not a new third-party codec.

The outer mathematical-region path can still use its original 128KiB cross-region history where not using this alternate; selected transformed BWT blocks become independent. Full candidate integrated 12 original source files all restored exactly and embedded SHA passed, with candidate aggregate archive size **60,244,257 B** on ALL ten passes. Individual archive sizes:
- dickens **3,155,856 B** (unchanged),
- mozilla **18,609,139 B** (unchanged),
- mr **2,729,261 B** (unchanged),
- nci **2,163,199 B** (unchanged),
- **ooffice 2,692,252 B**, compared with accepted **3,033,712 B** (**-341,460 B**),
- osdb **3,425,086 B** (unchanged),
- reymont **1,455,234 B** (unchanged),
- samba **5,134,262 B** (unchanged),
- sao **6,261,713 B** (unchanged),
- webster **9,835,416 B** (unchanged),
- xml **515,852 B** (unchanged),
- x-ray **4,266,987 B** (unchanged).

**TEN new candidate-only complete source runs**, each doing all twelve complete file encode/decode/compare. An initial one-pass screen (not counted in the ten-run sample) was also byte-exact at **5.458162294 s compress / 2.059992528 s decompress**, and was sufficient to justify the promotional measurement. All ten archive totals identical:

| Complete candidate run | Full source archive bytes | Encoding seconds | Decoding seconds |
|---:|---:|---:|---:|
| 1 | 60,244,257 | 5.547330414 | 2.072063770 |
| 2 | 60,244,257 | 5.662599750 | 2.151494101 |
| 3 | 60,244,257 | 5.590220905 | 2.127602480 |
| 4 | 60,244,257 | 5.627089052 | 2.059555270 |
| 5 | 60,244,257 | 5.568362539 | 2.195283429 |
| 6 | 60,244,257 | 5.756946187 | 2.153083612 |
| 7 | 60,244,257 | 5.496503511 | 2.163248014 |
| 8 | 60,244,257 | 5.531726473 | 2.089835416 |
| 9 | 60,244,257 | 5.515130994 | 2.069127816 |
| 10 | 60,244,257 | 5.534716237 | 2.083566654 |
| **Median** | **60,244,257** | **5.557846477** | **2.108718948** |

Previously accepted `current`, **retained historical ten complete run medians**, NOT rerun for comparison: **60,585,717 B / 5.493743735 s encode / 2.194427893 s decode**. The candidate yields **-341,460 B (-0.563598% of compressed size)**, corresponding to **+0.161113 percentage points source-byte reduction**; candidate encoding is **+1.166832% slower nominal**, decompression is **-3.905754% faster nominal**, compared with independently timed historical values. These independent-session runtime comparisons are **not paired-controlled**. Compression-size saving is exact; time data **do not demonstrate a nonregressing compression median**. User's strict acceptance criterion forbids promotion even though two of three metrics improve. **REJECT, no production source committed.**

Mandatory full comparison for the actual ten-run integrated candidate:

| Compressor | Total archived size B | Reduction from 211,938,580 B original | Encode s | Decode s |
|---|---:|---:|---:|---:|
| Zstd-3 (retained historical external CLI) | 66,216,569 | 68.7567% | 1.2561 | 0.8551 |
| XZ-6 (retained historical external CLI) | 49,233,340 | 76.7700% | 116.9892 | 2.3079 |
| **Current JLZ (retained historical ten-run median)** | **60,585,717** | **71.413550%** | **5.493744** | **2.194428** |
| **Candidate BCJ→BWT (NEW ten-run median)** | **60,244,257** | **71.574662%** | **5.557846** | **2.108719** |

| Candidate compressed-size delta | Candidate compression-time delta | Candidate decompression-time delta |
|---:|---:|---:|
| **-0.563598%** | **+1.166832% nominal** | **-3.905754% nominal** |

External CLI measurements not methodologically identical to JLZ in-memory measurements. Earlier measured current/Zstd/XZ speed references were reused unchanged.

**Validation** for disposable candidate:
- **12/12** original Silesia SHA-256 PASS, 120/120 complete source reconstructions exact and full file-embedded checksum verified.
- `go test -count=1 ./...` PASS, `go vet ./...` PASS, targeted race-enabled `TestLayeredBCJBWTMalformedAndExact` PASS, Windows/amd64 CGO-disabled cross-compilation PASS (not native Windows runtime).
- Explicit candidate-mode test built a valid BCJ→BWT mixed archive, restored it exactly, and rejected unknown wrapped BWT subtype, truncated stream and missing payload.
- Candidate CLI decoded an **actual previously approved `mr` archive** to byte-identical original source. Old approved CLI correctly rejects a new ooffice mode19 archive (expected **forward incompatibility**). Candidate CLI's ooffice archive size confirmed **2,692,252 B** and exact restore.
- All tested implementation files remain **disposable**. The retained local experiment's `engine/bwt.go` and `engine/cross.go` text was based on a local approved-code snapshot whose source blobs differ in nonessential structure from live `main`; before ANY eventual promotion, verify/port against exact live main and repeat correctness and ten-run candidate medians. No promotion occurred.

**One economical follow-up:** Reused `regionalMathUseful`'s existing distributed preflight mode choices to skip BCJ scanning on full files without two of three positive BCJ→BWT samples. The three actual math-selected complete files had preflight outcomes **mozilla modes 10/10/10**, **ooffice 19/19/11**, **samba 10/10/10**. This file-level gate generated **exactly the same 60,244,257 B** complete candidate archive and reconstructed all files in a SINGLE screen, but took **5.716874 s encoding / 2.145389 s decoding**; no demonstrated improvement or reason to tune thresholds. **Stop file-size, opcode-density, gate-depth microtuning**; preserve main unchanged.

## Layered concept 4: previously rejected BCJ + previously rejected cost-aware math parser (full integration, ONE screen, REJECT)

Built a distinct fully reversible **combined candidate** on a current JLZ copy: the prior original native **BCJ→math** option plus earlier **source-histogram priced selection among existing exact/XOR/ADD matches** in the same path, with no third-party backend. This is NOT the new BCJ→BWT branch.

One complete original 12-file Silesia encode/decode screen produced **60,273,819 B / 6.387206 s encode / 2.234396 s decode**, all original files exact. Versus historical `current`, size **-311,898 B** while **both** measured speeds were worse (one unpaired complete screen), as expected given the previous individual time losses. For comparison, the integrated BCJ→BWT ten-run candidate was smaller *and* materially faster, so the combination was **dominated**. **REJECT without ten-run suite; no source promoted.** This is concrete evidence rather than speculative claims about combining individually rejected methods.

## Layered concept 5: short reversible grammar on currently extracted math match-length streams (positive tiny component signal, no candidate)

Historical complementary audit measured a deeper sixteen-rule dictionary providing **15,138 B** optimistic selected math-length-stream savings across nine separately selected archived math blocks, but did not establish a final current-end-to-end tri-metric win. To verify the role of **current native entropy packers**, implemented an original reversible short byte-pair substitution grammar with up to four rules, using only otherwise absent byte markers. Each rule `marker→pair` has a three-byte transmitted descriptor; charged a selector byte, rule-count byte and original length varint before feeding the reduced byte sequence into existing `packFiniteState` (which itself includes every native model/table cost).

On first/middle/final independently tokenized original `mozilla`, `ooffice`, `samba` blocks, all nine packed grammar streams inverse-restored the exact original length-byte stream. Current native packed length streams **73,861 B**. Best fully charged **1–4-rule** alternative per block **71,074 B**, saving **2,787 B**, with **6/9** individual wins; mozilla saves 1,151 B, ooffice 56 B, samba 1,580 B. **This is an upper-bound per-block selection component result**, with full rule and stream metadata but not a whole-file codec mode or complete Silesia archive benchmark. The effect is small and extra table search/model recompressions add encoder CPU; **do not integrate or promote without a materially faster one-pass deterministic rule detector and a strong full-model benefit**. Neither the general theory of Tietze/grammar nor the deeper historical candidate is marked impossible.

## Final decisions and next honest research queue

**All new concrete implemented layered placements recorded; no promotion.** The added BCJ→BWT branch is the **first explicitly benchmarked ten-run end-to-end layered candidate in this reassessment** and materially outperforms prior BCJ→math in compressed bytes and median decoding speed, yet still has nominal +1.17% encoding regression. The mandatory three-metric contract forbids accepting it. No other screens have ten complete integrated candidate runs, and they are not presented as such.

The ledger and new [placement matrix](../research/layer-placement-matrix.md) separate:
- **Already integrated supplements** (existing exact/XOR/ADD, BWT MDL 2/4 parts, matched-literal residuals, prior sparse defect candidate);
- **New exact layered component screens** (LZP→BWT/literals, GF2 bit transpose/recurrence on math literals, BCJ→BWT/native/predictive, grammar→math lengths);
- **New integrated whole-codec results** (BCJ→BWT ten runs; BCJ+priced math one screen);
- **Not yet tested** placements (nontrivial graphs/syzygy on actual current selected residuals, algorithmically distinct cheap geometry/boundary on current additive x-ray, proper latent-variable models, reference transformations, metadata-specific grammar and constraint enumerators);
- **Ineligible without new format/CPU premise** combinations (external reference universe, large learned decoder, replay of DEFLATE/Zstd-backed historical transformers, cross-file references).

**Do not perform brute-force combinatorial layering**: two individually CPU-expensive candidates can reduce size and still be a Pareto failure, as the BCJ+priced-math integration demonstrates. Future tests require a concrete encoder- and decoder-work elimination argument, rigorous exact metadata cost, and comparison against the closest failed *placement*, not just the abstract name of a theorem. The user explicitly requested **layered evaluation by default**, with observed CPU costs sufficient to reject a concrete placement without denying the broader theory.

**No approved JLZ source change.**
