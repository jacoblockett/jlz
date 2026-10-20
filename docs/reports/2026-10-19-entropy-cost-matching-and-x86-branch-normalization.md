# Entropy-cost match decisions and reversible relative-branch normalization

**Recorded:** 2026-10-19 America/New_York. **Final decisions: REJECT both concrete integrated implementations; keep current unchanged.** No new numbered development rounds or release identifiers. These were original, reversible native-Go experiments, inspired by audited **Zstd high-quality parse**, **XZ/LZMA2 price-aware parsing and BCJ preprocessing**, and **Brotli optimized reference paths**. No compressor library, ported codec, copied external source, additional user command, or training table.

## Approved reference and strict rule

Current original **twelve-file Silesia** (211,938,580 B), fixed accepted **60,585,717 compressed B / 5.493743735 s median encode / 2.194427893 s median decode** (retained earlier ten-run values, NOT remeasured). External fixed benchmarks Zstd-3 **66,216,569 B / 1.2561 s / 0.8551 s**, XZ-6 **49,233,340 B / 116.9892 s / 2.3079 s** (external CLI speed conventions differ from in-memory JLZ).

Without the user's explicit tradeoff approval, no candidate can increase archive size, encoding time, **or** decoding time. Favor cheap component/one-run rejection and no production source changes when a strictly nonregressing ten-run outcome is implausible.

## Strategy 1: model-priced choices among JLZ's exact/XOR/ADD reference candidates

**Closest rejected approaches:** exact-only replacement dual-hash matcher damaged compression, added-second-full-token-encoding cost too much CPU, merely lengthening history and changing lazy match score cutoffs caused diminishing returns. **This implementation was materially different:** it retained the original JLZ 3-fingerprint exact/XOR/ADD finder, recent-distance code cache, existing one-byte lazy search, native token grammar, and decoder, but fed the previously **already discovered** legal candidate choices through a *source-frequency-informed price* covering expected literal symbol bits, opcode, varint length, logarithmic distance slot + extra bits, transformed mask and cache-repeat state. No additional candidate search depth, entropy model, or new decoder mode was introduced.

Original Go disposable source, local directory /mnt/data/jlz_price_parser_trial. Diagnostic \`engine/zz_cost_parse_research_test.go\` (not production) initially compared alternative native actual packed streams and decoded each sample exactly:
- **Nine 256 KiB-or-shorter source regions** (first/middle/final from original Silesia mozilla/ooffice/samba), baseline packed native match streams total **754,913 B**.
- A path that called existing matcher with a per-candidate enumeration callback and one-step discounted choice, using an original raw-byte-derived symbol histogram, measured **748,864 B (-6,049 B)** but the tokenizer generally took **1.5–2x baseline time**.
- An *optimistic* matched-literal-frequency model sampled from **previously fully tokenized** baseline streams measured **748,981 B (-5,932 B)** under that two-choice path. This requires a redundant prior parse and was not viable as a real candidate.
- **In-place price selection** integrated into the original matcher, avoiding the separate callback candidate enumeration, still retained all source candidate types. Using raw input histogram plus actual baseline metadata prices, measured **752,446 B (-2,467 B)**; one sample (samba middle) regressed in size. This model, too, depended on source-local previously encoded metadata tables and was screening only.
- A **cheap production-plausible, untrained, raw-only MDL price proxy** used the current source's literal histogram and fixed generic op/varint/distance code estimates; no initial match pass or trained model. On the nine full source region screens it yielded **751,552 B (-3,361 B)**, with one samba middle-region regression (+492 B), and usually slightly more parser work than the original.

The last cheap native proxy was then **fully integrated** into \`encodeMathRegionBlockHistory\` in a disposable source copy, using one byte-exact existing native token stream per region and the **same format/decoder**. A SINGLE entire original twelve-file Silesia encode/decode screening pass produced:
- **60,435,710 B** complete archive bytes, saving **150,007 B (-0.24759% previous complete compressed size)**.
- Encode **6.667317086 s** and decode **2.311218779 s** (ONE screen only; not a candidate ten-run median). Versus prior current historic medians nominal **+21.36% encode / +5.32% decode**, unpaired.
- Only mathematically selected files changed: **mozilla 18,609,139→18,521,240 B (-87,899)**, **ooffice 3,033,712→3,011,183 B (-22,529)**, **samba 5,134,262→5,094,683 B (-39,579)**. Nine other complete Silesia archive sizes byte-identical to current.
- All twelve complete sources reconstructed original bytes byte-for-byte and passed checksum. Full uncached \`go test -count=1 ./...\`, \`go vet ./...\`, targeted \`go test -race\` and Windows/amd64 CGO-disabled crossbuild PASS. Canonical original 12 Silesia SHA-256 manifest identities PASS.

**REJECT this exact parser as a production candidate** due gross failure of both timing requirements. No ten-run benchmark warranted, no source promoted. The underlying fact that approximate local match scores differ from actual entropy costs remains valid; **do not re-run nearby opcode-price constant, one-byte lookahead or histogram-bucket parameter sweeps**. To revisit require a *structurally different* virtually free exact cost-recovery mechanism with credible three-metric performance.

## Strategy 2: reversible x86-relative branch normalization before existing math compression

**Architectural origin:** the XZ/liblzma BCJ family transforms relative executable branch operands into a more redundant address domain, then inverses it after decompression. Original implementation does **not port liblzma x86.c**. Mechanism is distinct from fixed-byte XOR/ADD matches, signed 16-bit modular differences, post-BWT contexts and score tweaks.

Implemented a portable self-synchronizing reversible 32-bit relative operand transform:
- Scan each independent original input block for bytes \`0xE8\` and \`0xE9\` with four following bytes; interpret the next four bytes little-endian, add \`uint32(opcode_index + 5)\` modulo 2^32; **skip** all four operand bytes before scanning for the next instruction marker. Inverse scans identically and **subtracts** the offset, reconstructing arbitrary bytes without requiring actual executable correctness, instruction decoding, global file offsets, sign-extension metadata, or trained tables.
- Original Go prototype tested complete native packed mathematical payloads on 12 first/middle/final original Silesia source blocks from mozilla, ooffice, samba, and x-ray. **All 12 forward/inverse transforms exact.** Full native transformed payload total **1,101,050 B including a nominal 1-byte mode selector per block**, versus baseline **1,116,760 B**; primarily due to ooffice. These are diagnostic blocks, NOT integrated archives.
- **Every one of ooffice's 24 complete natural 256 KiB blocks** then checked with actual native math packing and byte-exact inversion: unchanged original math component stream total **3,036,867 B**, all-BCJ **2,752,577 B**, conservative **once-chosen 2%-opcode-density gate (plus 64KiB minimum block) 2,754,921 B**, a *component-level* nominal **-281,946 B** on those source blocks. 21 individual transformed source blocks saved bytes, 3 lost. The 20 gated full-size blocks ALL won; some source block losses happened at low instruction-like density. These screens were performed without selected cross-region history or the actual global file/model routing, so component savings must not be added to complete archive benefits.
- Timed each original and transformed native math codec separately in the 24-block diagnostic; original packed math stage **1.304 s aggregate** vs transformed **1.182 s aggregate**, unpaired consecutive measurements and **not** whole JLZ wall time. BCJ forward scanning itself was generally around **0.1–0.5 ms per source block**.

**Fully integrated real candidate:** original Go \`engine/branch_relative.go\` in disposable local /mnt/data/jlz_bcj_candidate and minimal changes to \`engine/region_fast.go\` and \`engine/cross.go\`; current \`main\` source unchanged. On selected BCJ-dense mathematical input blocks only, preflight **chooses transform before encoding**, avoiding a second full match/entropy encode; uses the same original mathematical matching/token/entropy backend and wraps its existing payload with one inner block-mode byte under a new **experimental internal block mode 18**. Current raw fallback and all other automatic paths remain. Decoder validates wrapped native mode, restores the prior original-byte mode, then applies the byte-exact inverse transform. No precompressed backend, user option, file-name branch or learned Silesia identity. The density gate is a heuristic, not a universally proven coding win.

A first **ONE full entire original 12-file Silesia screen** produced **60,431,550 B / 5.790166351 s encode / 2.300569554 s decode**, all twelve restored exactly. A separate compiled candidate test binary then performed **TEN complete candidate-only twelve-file encode/decode passes**, all verifying full original bytes and embedded SHA. Long combined wrapper terminated after **FOUR completed runs**, which were retained, followed by three independently completed batches of TWO runs each; incomplete work not counted. All ten produced exactly **60,431,550 B**:

| Full candidate run | Complete archive B | Encode s | Decode s |
|---:|---:|---:|---:|
| 1 | 60,431,550 | 5.703103748 | 2.218337720 |
| 2 | 60,431,550 | 5.574155822 | 2.229494705 |
| 3 | 60,431,550 | 5.801827650 | 2.273892770 |
| 4 | 60,431,550 | 5.993373254 | 2.486765634 |
| 5 | 60,431,550 | 5.720715439 | 2.262779435 |
| 6 | 60,431,550 | 5.723128785 | 2.237708802 |
| 7 | 60,431,550 | 5.710773435 | 2.276987862 |
| 8 | 60,431,550 | 5.561223717 | 2.270586276 |
| 9 | 60,431,550 | 5.917244089 | 2.296418397 |
| 10 | 60,431,550 | 5.703572289 | 2.224205777 |
| **Median (each time series)** | **60,431,550** | **5.715744437** | **2.266682856** |

Compared with retained accepted current **60,585,717 B / 5.493743735 / 2.194427893** historical ten-run medians, BCJ complete archives genuinely saved **154,167 B (-0.254461% of complete compressed archive bytes, +0.072741 percentage points reduction of original)**, but encoding was nominal **+4.04097%** and decoding **+3.29266%** slower. The timing samples are **independent, unpaired sessions**, and no causal precision is implied. Even allowing for noise, no credible demonstrated performance non-regression exists in both measures.

**Complete original per-file size:** only **ooffice: current 3,033,712 B → candidate 2,879,545 B (-154,167 B)**; all eleven other complete files **byte-identical**. This reflects cross-reference and global preflight constraints not represented by the larger isolated 24-block office component savings.

**Correctness and compatibility:** full uncached Go package suite PASS, vet PASS, targeted race tests PASS, Windows/amd64 cross-build PASS (not native Windows execution); all 12 original Silesia canonical SHA-256 PASS, **120/120** full-file exact candidate reconstructions. Candidate CLI successfully decoded a real archive produced by the prior approved current CLI (entire original \`mr\`, source SHA verified), so backward decode is preserved. Approved older binary rejects new mode-18 ooffice archive with “cross block mode” as expected, while new candidate decoder restores it byte-exactly. **New mode is forward-incompatible**, a further cost against adoption.

## Required comparison: ONLY benchmarked integrated BCJ candidate

| Model | Complete original Silesia compressed size | Original reduction | Encode time | Decode time |
|---|---:|---:|---:|---:|
| Zstd-3 (historical external CLI) | 66,216,569 B | 68.76% | 1.2561 s | 0.8551 s |
| XZ-6 (historical external CLI) | 49,233,340 B | 76.77% | 116.9892 s | 2.3079 s |
| **Current JLZ (retained historical 10-run median)** | **60,585,717 B** | **71.41355%** | **5.493744 s** | **2.194428 s** |
| **BCJ candidate (NEW 10-run median)** | **60,431,550 B** | **71.48629%** | **5.715744 s** | **2.266683 s** |

| BCJ vs current archive size delta | BCJ vs current encode-time delta | BCJ vs current decode-time delta |
|---:|---:|---:|
| **-0.254461% (-154,167 B)** | **+4.04097% nominal** | **+3.29266% nominal** |

The external fixed Zstd/XZ CLI benchmark timing conventions are not the same as JLZ in-memory benchmark timings.

## Final decision and anti-spiral check

**REJECT BOTH implementations for production.** The parser has real native packed-stream benefit, but a failed full-corpus screen on CPU; don't tweak its score coefficients or fallback match thresholds. The BCJ branch-relative domain transform is reversible and found concentrated legitimate source-file structure, but its ten-run candidate performance worsened both encoder and decoder and adds new decoder mode18; **do not promote it merely for 154 KB saving**. Source remains unchanged at current; the disposable trial directories preserve local prototypes for verification but no rejected code is committed.

The **theory** of transform-domain executable normalization is not disproved. A materially different mechanism (for example, decoder-synchronous normalization that is fused into an already necessary input or output scan, needing no extra pass and no second format, not a slightly different opcode density gate) would be needed for renewed experimentation. Do not sweep opcode thresholds, source widths or repeatedly repackage this same work as a new idea.

**Future strategy pivot:** avoid repeating BWT entropy splitting, history-size expansion, dual matcher, small parser cost adjustments and free-standing BCJ preprocessing. Look for a structurally new *zero-extra-pass* transform or a native match token/framing improvement with direct decoder speed gains. Continue architectural learning from leading compressors, treating their expensive high-compression modes as source-level design evidence, never as an excuse to trade away compression/decompression speed.
