# Current-engine complementary layers: geometry, reverse references, distance roles, sparse-copy corrections

**Date:** 2026-10-20. **Outcome: four distinct concrete placement screens rejected or deferred; NO benchmarked integrated candidate.** Active \`current\` on GitHub main unchanged: the retained complete original twelve-file Silesia archive is **60,585,717 B**, historical 10-run medians **5.493743735 s encode and 2.194427893 s decode**. These timings were **not rerun**. No comparison table should be presented to the user because **none of these component placements became a full-codec benchmarked candidate**.

**Dataset:** exclusively the **twelve complete original Silesia source files**, 211,938,580 B. SHA-256 verification against \`docs/reports/2026-10-18-round23-silesia-sha256.txt\`: **12/12 PASS**. All screens use full original source files or natural 256 KiB blocks drawn from them. Individual block and stream packing sizes are **not complete JLZ archive sizes**, and do not include a new complete archive format or all inter-block model benefits. No external compression backend or third-party codec imported.

## 1. Two-dimensional 16-bit geometry within current additive coordinate coding — SCREENED negative

**Historical difference:** JLZ's accepted additive coder predicts each 16-bit little-endian integer using a modular difference from the immediately preceding 16-bit word, zigzags the residual, and independently codes its low and high bytes through existing native raw/RLE/Huffman. Older 2D grids used widths 16/256 and sparse boundaries on selected additive residuals. This trial first measured actual original \`x-ray\` spatial recurrence: within a fixed data-supported source interval, neighboring 16-bit words at stride **1,900** exhibit considerably more near-equal values than nearby strides (peak around 1,900), suggesting a meaningful 2D image row. This is not an arbitrarily tuned historical grid width.

**Actual native-Go reversible component:** \`/mnt/data/jlz_geometry_1900/engine/zz_geometry_screen_test.go\`. Applied vertical (previous row), second-order gradient (left + above − upper-left, modulo 2^16), and Paeth (choose left/above/upper-left nearest the predicted gradient) to each full natural 256 KiB region. Zigzag signed differences, re-use the **current native low/high stream encoders**, charge **three extra bytes per region** as a favorable width/mode selector estimate, decode both native packed streams, invert predictor using **already restored source values** and compare every original byte. PASS on every tested block from complete original files x-ray, mozilla, ooffice, nci.

**Actual selected-domain component sizes, per original file:** 
- **x-ray**, relevant because current globally selects additive coding: approved additive independent source blocks **4,266,715 B**; vertical **4,365,048 B**; gradient **4,392,314 B**; Paeth **4,272,192 B** (**+5,477 B** overall), despite **19** isolated Paeth block wins. Unrealistically favorable per-block *choose-current-or-any-new-model* oracle **4,261,531 B**, only **5,184 B** saving **before an actual decoder mode, selector complexity, and added per-word work**. **Reject practical integration**.
- **mozilla**, globally math-selected, not additive: additive original **44,158,003 B**, vertical **43,752,256 B**, gradient **48,442,807 B**, Paeth **44,701,419 B**, selective oracle **43,234,656 B**. Even its *best additive oracle* is **more than twice** the actual accepted complete-file mathematical archive **18,609,139 B**. No impact on the selected path.
- **ooffice**, globally math-selected: additive original **5,768,611 B**, vertical **5,955,531 B**, gradient **5,997,999 B**, Paeth **5,824,882 B**, oracle **5,745,017 B**; still inferior to math selected archive.
- **nci**, globally BWT-selected: additive original **17,672,793 B**, vertical **19,848,879 B**, gradient **25,041,971 B**, Paeth **19,100,544 B**, oracle identical to incumbent additive.

**Conclusion:** this actual 2D image predictor inside the already approved additive stream is valid/reversible but cannot justify a new decoder mode, especially with obligatory per-word CPU. **CLOSE exactly this row-geometry placement**, not geometry/boundaries/contours in general. Do not sweep guessed row-width values around 1,900; source-backed width was selected for the diagnostic. Extra global width detection is **not included** in the optimistic component sizes and would further weaken it.

## 2. Backward-facing mirrored copy references alongside existing math references — SCREENED negative

**Historical difference:** existing JLZ reference tokens are forward exact-copy, constant XOR/ADD and cached distances; earlier mirror/reverse references were not implemented as an original-Silesia candidate token. This new diagnostic searches for a **reverse-direction reference ending at a prior source position** within 128 KiB. It is not a bit reversal of a byte or a compressor backend.

**Actual Go diagnostic:** \`/mnt/data/jlz_reverse_match_screen.go\`. One last-seen hash index of backward-facing four-byte anchors at previous source positions; forward current four-byte prefix lookup; byte-exact reverse extension up to 4 KiB. Tested the **complete original files**, including mozilla, ooffice, samba, x-ray, nci. A naive gross count vastly overstates savings because symmetric zero/constant runs produce millions of overlapping mirrored matches that the current engine already compresses. Applied explicit filters: count only nonoverlapping prospective source spans of at least 16 bytes, require at least four different byte values within the first 64 matched bytes, and exclude candidates where the same source range can be copied **normally forward** without mirroring. This is still only a diagnostic; other existing JLZ references can explain additional matches.

- **mozilla:** **114** distinct qualifying spans, **3,291 source bytes** gross filtered, compared with ~**3.54 s** original standalone scan on the complete source.
- **ooffice:** **one**, **16 B**.
- **samba:** **282**, **18,616 B**.
- **x-ray:** **zero**.
- **nci:** **five**, **140 B** (current BWT-selected, not relevant to original math).
  
**Outcome:** only ~**21.9 KB** nontrivial math-path source coverage before new reference opcode, distance, length, backward read, exceptions, or metadata, and the initial Go index scan is far too expensive. **Do not integrate or sweep reverse anchor lengths**. This implementation does **not** disprove mathematically different algebraic reference bonds or mirrored compressed-plane models.

## 3. Opcode-conditioned split of native mathematical distance metadata — SCREENED weak

**Historical difference:** the old two-/three-way match-length-conditioned distance entropy mode expanded selected streams and was rejected. This diagnostic conditions on an **existing decoded operation type** rather than match length: the present stream mixes unsigned-varint explicit match distances and single-byte repeat-offset indices for \`RepXor/RepAdd\`. A repeat index zero can force \`packVarintStream\`'s slot candidate to fail, leaving a possibly poorly coded mixed byte stream. Decoder already knows the opcodes, so no new type symbols need be sent. This is a structurally different split, not a distance-window expansion.

**Actual native Go component test:** \`/mnt/data/jlz_geometry_1900/engine/zz_distance_layer_test.go\` on **every complete-file natural source block** of original mozilla, ooffice, samba, mr, x-ray; each block tokenized with the current original exact/XOR/ADD matcher. Split actual distance byte stream by opcode into explicit varints plus repeat-index bytes, pack explicit values through existing \`packVarintStream\`, repeat indices through existing \`packStream\`, count **one extra selector byte**, and prove exact reconstruction of original mixed distance byte stream by opcode replay. PASS on every tested block.

**Current encoded distance stream B → role-separated total B → favorable choose-per-block oracle B**:
- mozilla: **2,321,309 → 2,320,848 → 2,320,103** (oracle save **1,206 B**).
- ooffice: **347,082 → 344,638 → 344,554** (oracle save **2,528 B**).
- samba: **1,448,099 → 1,439,605 → 1,439,309** (oracle save **8,790 B**).
- mr: **896,711 → 896,867 → 896,711** (no wins).
- x-ray: **66,281 → 66,413 → 66,281** (no wins).

The three math-selected complete-file *component oracle* saving totals **12,524 B** even before adding an actual decoder mode, the extra decoded stream handling and any full pipeline match-history interactions. **Reject this concrete role-separated distance coder** as not worth its CPU/format cost, even though the existing mixed-distance stream has a recognizable structural imperfection. A truly different compact metadata representation remains open, not disproven by this weak split.

## 4. Sparse exceptions layered onto already selected exact match tokens — full-source native packed component SCREEN, negative

**Distinct placement:** Previous sparse-defect coding experimented on signed additive residual planes. Here we add exceptions **inside the existing mathematical match grammar**, immediately after a real exact-copy decision with length >=32. No new long-range indexes, second match search or deleted exact/XOR/ADD candidate families. An actual source reference can be extended up to 256 additional bytes if at least 64 extend and at most four byte substitutions are required. Proposed decoder copies earlier bytes, substituting the specified original bytes at monotonically increasing token-relative offsets.

**Gross opportunity diagnostic:** current exact-match tokenizer output was decoded/parsed for every original natural 256KiB source block of complete mathematical source files. Existing selected exact copies potentially extend through <=4 defects within 256 bytes:
- mozilla **1,216** qualifying continuations covering **310,792** extra source bytes.
- samba **1,103** / **282,368 bytes**.
- ooffice **10** / **2,560 bytes**.
- mr **8** / **2,048 bytes** (BWT-selected current path).
This is **gross coverage of overlapping future literal/match commands**, not net compressed-byte savings or a bound on actual JLZ improvements.

**Actual native reversible component implementation:** \`/mnt/data/jlz_sparse_match_trial/engine/match.go\` modified only in a **disposable** source copy, exposing \`tokenizeMatchV2RangeStrideSparse(...,true)\` beside the incumbent function. It reuses current matcher choices and emits an experimental nibble opcode **10** to hold original exact backreference distance, length, correction count and original replacement bytes at precise token-relative varint offsets. A separately written decoder in \`engine/zz_sparse_match_layer_test.go\` replays all original exact/XOR/ADD/rep tokens plus sparse corrections, verifies range, monotonic correction positions, backward reference bounds, complete coverage and **byte-for-byte original restoration for every tested block**. Existing native five-column opcode, length, distance, transform (including correction metadata), and literal stream encoders produced **actual fully packed component byte counts**, including an experimental one-byte mode selector. The mode10 is an **experimental token opcode**, not a committed archive block mode; old decoder has not been changed and no final JLZ archive was emitted by this component. A first trial mistakenly reused opcode 9, which is already assigned to \`matchRepAdd\`; corrected to unused opcode **10** before final evidence. That initial malformed trial was **not used** for size conclusions.

**Actual original complete-source block component B**:
- mozilla: incumbent **18,811,124 B**, sparse-token layered **18,872,265 B** (**+61,141 B**), **0** individually winning blocks.
- ooffice: **3,063,783 → 3,065,738 B** (**+1,955 B**), **0** wins.
- samba: **5,301,701 → 5,323,876 B** (**+22,175 B**), **0** wins.
- mr: **3,485,397 → 3,485,492 B** (**+95 B**), **0** wins.

**Why it loses:** current existing tokenizer already turns most of those apparently interesting future source bytes into cheaper nearby literals/repeated references. Charging correction offsets/values and the new opcode column makes the actual candidate more expensive on EVERY block despite hundreds of KB in gross reference coverage. **REJECT this actual supplementary sparse-copy design**. Do not test additional exception-count/length threshold sweeps; the measured structural redundancy is already substantially accounted for. The theorem that sparse corrections can encode other structures remains open at other placements.

## Production and verification

No experimental codec Go files were committed to \`main\` and no candidate was promoted. Only disposable local Go code + research reporting. The geometry and sparse-token checks ran against the accepted local Go code base. Exact packed native geometry streams were decoded, inverted and compared with originals on all tested blocks; reverse references were exact source comparisons; role-separated distances were reconstructed exactly from their decoded operation classes; sparse-match token output was replayed by an independent original Go decoder across complete-source natural blocks, confirming exact output and structural bounds. The original twelve Silesia SHA-256 manifest passed **12/12**. The sparse layer \`go test -count=1 ./engine -run '^TestSparseLayerScreen$'\` passed; the test file was run with original test fixtures in a disposable source copy. These are **component correctness tests**, not full uncached Go suite/vet/race/Windows results for an integrated candidate, and no such claim is made.

**ANTI-SPIRAL STOP:** Closed four **concrete additive placements** after original-data screens with no credible big fully framed improvement. Do not conflate the negative results with failure of geometry, sparse coding, p-adic number theory, reverse algebraic bonds, or causal graph compression as whole theories. User's strict three-metric nonregression still applies. Before the next experiment, compare actual source and failed placements in the [strategy ledger](../research/strategy-ledger.md) and [layer-placement matrix](../research/layer-placement-matrix.md), then choose an orthogonal, substantial mechanism rather than parameter tweaks.
