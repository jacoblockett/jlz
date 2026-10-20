# Source-level competitor architecture audit: Zstd, XZ/LZMA2 and Brotli

**Date:** 2026-10-19 (America/New_York). **Scope: primary source inspection and comparison to JLZ `main`.** No JLZ candidate implemented or benchmarked in this audit; current remains unchanged at **60,585,717 B / 5.493744 s encode / 2.194428 s decode**, the previously established ten-complete-original-Silesia-run medians. No candidate table or false speed claims. **No external compressor implementation was copied, modified or used as a JLZ backend.**

## Objective and method

User explicitly requested systematic investigation of *actual industry-leading compressors* rather than disconnected mathematical tuning. Inspected up-to-date upstream source paths for Zstd high-level parsing, match finders, LDM, entropy/sequence handling; XZ/LZMA2 normal and fast parsers, match finders, probability prices, block resets and executable branch filters; Brotli high-quality Zopfli shortest path, model-price updates, command/distance/literal representation, contextual histograms and static dictionary transformations. These are source-level architectural analyses, **not claims of having replicated competitor benchmarks or exhaustively audited every function**. Reconciled with the current live JLZ Go `engine/match.go`, `engine/region_fast.go`, `engine/packed.go`, `engine/slots.go`, `engine/lengths.go`, `engine/bwt.go`, and [strategy ledger](../research/strategy-ledger.md).

### Primary upstream source references

**Zstd (Meta):**
- [`zstd_compress.c` strategy dispatch / block and matching configuration](https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_compress.c)
- [`zstd_lazy.c` greedy/lazy family and match candidates](https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_lazy.c)
- [`zstd_opt.c` optimal parser, statistical sequence prices](https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_opt.c)
- [`zstd_ldm.c` separate long-distance matcher and sequence injection](https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_ldm.c)
- [`zstd_compress_sequences.c` encoded match-sequence columns](https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_compress_sequences.c)
- [Zstd API: LDM parameters, compression strategy classes](https://github.com/facebook/zstd/blob/dev/lib/zstd.h)

**XZ / liblzma (Tukaani):**
- [`lzma_encoder_optimum_normal.c` normal-mode state-aware cost/optimal path](https://github.com/tukaani-project/xz/blob/master/src/liblzma/lzma/lzma_encoder_optimum_normal.c)
- [`lzma_encoder_optimum_fast.c` faster lookahead and nice-match threshold](https://github.com/tukaani-project/xz/blob/master/src/liblzma/lzma/lzma_encoder_optimum_fast.c)
- [`lz_encoder_mf.c` bounded hash-chain and binary-tree match discovery](https://github.com/tukaani-project/xz/blob/master/src/liblzma/lz/lz_encoder_mf.c)
- [`lzma2_encoder.c` chunk framing, uncompressed fallbacks and model/dictionary resets](https://github.com/tukaani-project/xz/blob/master/src/liblzma/lzma/lzma2_encoder.c)
- [`simple/x86.c` x86 relative-branch address preprocessing](https://github.com/tukaani-project/xz/blob/master/src/liblzma/simple/x86.c)
- [liblzma public match finder/fast-v-normal API](https://github.com/tukaani-project/xz/blob/master/src/liblzma/api/lzma/lzma.h)

**Brotli (Google):**
- [`backward_references_hq.c` Zopfli-like shortest path and iterative cost model](https://github.com/google/brotli/blob/master/c/enc/backward_references_hq.c)
- [`metablock.c` split literal/command/distance histograms, context maps/clustering](https://github.com/google/brotli/blob/master/c/enc/metablock.c)
- [`block_splitter.c` block segmentation](https://github.com/google/brotli/blob/master/c/enc/block_splitter.c)
- [`static_dict.c` transformed static dictionary matches](https://github.com/google/brotli/blob/master/c/enc/static_dict.c)
- [`c/common/context.c` decoder-known literal context classification](https://github.com/google/brotli/blob/master/c/common/context.c)

## Competitor findings and precise JLZ gaps

### 1. Zstd high levels: separate match discovery from sequence-path economics

Zstd uses a strategy family: fast/double-fast; greedy; one- and two-step lazy; binary-tree-assisted lazy; and high-effort binary-tree optimal parsing (btopt/btultra/btultra2). These are **different encoder match-generation and parse-cost budgets**, not different lossless decode formats. Optimal code enumerates match candidates, estimates costs using literal/length/offset distributions and recent-offset state, and finds a lower-price command path over a bounded region. In particular btultra2 has extra initial statistics work (explicitly noted in its source); this *costly* approach is unsuitable for JLZ's strict encoding Pareto gate in unmodified form. Zstd LDM independently finds long-distance sequences and injects them into normal block parsing, rather than merely enlarging every local match probe. Its production implementation is sophisticated and carries memory/decoder-window costs.

**Already present in JLZ:** separate opcode/length/distance/transform/literal streams; three repeat offsets; one-byte lazy parsing; validated exact/XOR/ADD references; native finite-state alternative on selected streams; region-level winner selection; basic local entropy bound. Prior exact-only 4/8-anchor fast parser degraded compression (rejected), and previous long-distance content anchors produced gross source coverage without net archive evidence at high scan cost.

**Still structurally missing:** explicit *sequence-path* cost search based on JLZ's actual emitted stream price estimates, with **transform type and repeat-offset-cache state** included. Current `scoreMatchChoice` substitutes fixed **5 or 6 bits per literal-covered byte** and one-/two-byte varint surrogates, while final native entropy uses nontrivial op/length/distance distributions and separate models. A greedy locally highest-score match can preclude a better combined set of future matches. A narrowly bounded, **cache-aware parse decision** using matches already discovered by JLZ's exact/XOR/ADD matcher is different from replacing it with a cheap exact-only parser. It can retain the decoder and format. Need quantify actual opportunity before any expensive integration.

**Anti-spiral restriction:** Do not copy Zstd's graph algorithms or run global optimal paths over all bytes. Test on already selected mathematical blocks, do not enumerate second full bucket/matcher passes or continuously tweak lazy thresholds. Advance only if a reversible bounded-work original Go prototype demonstrates meaningful changes after exact native entropy and acceptable work.

### 2. XZ/LZMA2: context-state costs, alternate match finders and transform filters

LZMA fast and normal compressors explicitly differ in parser work. Fast parser exploits repeat distances and acceptably long “nice” matches; normal parser uses dynamic price tables conditioned on LZMA coding state and byte position, including literal/match/rep flags, match length, distance slots and low/alignment bits. Its match-finder classes (HC3/HC4 hash chains; BT2/3/4 binary trees) trade speed, RAM and candidate diversity, with explicit search-depth caps. LZMA2 wrapper groups data into bounded chunks, may emit an uncompressed chunk when it cannot beat the raw size, and controls dictionary/state resets **separately**.

**Already present in JLZ:** repeat distances, native log-distance slots/low-bit bytes, 256 KiB regions, raw/Huffman fallback, bounded match buckets, independent decode groups, several predictors and transforms. JLZ's local matcher already estimates candidate length and distance but not the entropy cost of the *resulting command sequence*.

**Missing with plausible distinct scope:**
- A **locally priced parser** that reflects actual native opcode/length/distance/transform costs, rather than a copy of LZMA's adaptive arithmetic range codec (JLZ has its own). This overlaps Zstd's high-level principle and must be **one shared hypothesis**, not three nominal research rounds.
- **Transform-domain executable normalization**: the x86 BCJ prefilter converts selected relative branch/call displacements to more repetitive address-domain values while leaving a reversible byte stream. This is materially different from JLZ's accepted byte-wise XOR/ADD fixed-mask matches and its signed modular 16-bit word residual mode. It can be tested independently on the **complete original Silesia files** as a real reversible original native Go pretransform, charging selection bit, inverse pass and full archive framing. Do not port liblzma x86.c; derive own instruction-aware bounded transform, verify arbitrary non-executable bytes and original reconstruction, avoid applying blindly to all files. The full encoder+decoder CPU cost likely makes strict Pareto difficult; start with a cheap semantic relevance screen.
- LZMA's large BT dictionaries and expensive normal-mode parser are **LOW priority** under existing rejected larger-history and slow exact-match experiments. Merely widening hash depth/history is a known spiral and must not be tried as a new idea.

### 3. Brotli high quality: iterative modeled paths, context clustering and dictionaries

Brotli's high-quality source uses a Zopfli-style shortest-path graph over possible literal insertion, match length and distance choices. It builds histograms for literal/command/distance symbols, estimates Shannon-like code costs, finds a command path, and can **refit costs from chosen commands** for another iteration. The source's `ComputeMinimumCopyLength` and best-start queue implement pruning to avoid considering hopeless future copies. This is an expensive *compression-focused* strategy and is inappropriate to transplant wholesale into JLZ.

Brotli's meta-block encoder treats literal, command and distance streams separately, builds context maps and clusters histograms, and includes transformed static dictionary matches. All those aids are **model/computation decisions accompanied by metadata**, not free compression.

**Already present in JLZ:** split entropy streams, repeat distances, BWT symbol entropy local model selection, selected metadata fSE/legacy modes, raw fallback and several transform/candidate frontends. Rejected history includes role-conditioned rANS and three-table predecessor-MTF context coding with only modest byte savings and much more CPU, plus fixed segmented math metadata that expanded net bytes. Splitting histograms or adding more contexts is therefore **not a new idea** by itself.

**Remaining real gap:** cost-model recalibration and *dominance pruning* within match-path decisions, with reuse of existing candidate matches. Brotli and XZ/Zstd support the **same shared parsing hypothesis** rather than three independent architectures to test repeatedly. Static English dictionary lookup, transform-table side information and locale-specific built-ins would add a major permanent decoder asset, bias toward text files, and duplicate/complement existing matching only uncertainly; **deprioritize** absent compelling sample feasibility.

## Shared evidence map: hypotheses ranked by possible Pareto impact

**A. FIRST research priority: bounded-work, entropy-informed parser decisions on already available native match candidates.**
- Distinct from rejected fast exact-only parser: **keep original JLZ exact/XOR/ADD candidate generator**, improve how competing candidates are chosen.
- Distinct from costly full re-encode fallback: no second complete entropy pass, no new format, no two simultaneous token streams.
- Distinct from previous lazy-score threshold sweeps: path cost uses measured/estimated current actual native metadata/literal frequencies and repeat-cache effects, and prunes suboptimal options with a bounded lookahead rather than changing one threshold.
- First experiment: instrument currently *winning* mathematical regions of the original Silesia files; measure how often current winning matches conflict with cheaper two-choice continuations under a model calibrated from the **already emitted** native command streams. Implement reversible token-path tests with unchanged framing and original input restoration. Derive a pessimistic component upper bound on net bytes and added candidate-search work. **Stop immediately** if net native bits/pricing or CPU benefit is too small. Only then attempt complete JLZ integration and candidate-only ten original Silesia runs.
- Risk: choosing a new parse can alter symbol histograms, invalidate the provisional prices, and cost more than it saves; this is a hypothesis, not a reported win.

**B. SECOND distinct priority: x86 relative-offset reversible normalization, selected by source evidence.**
- A distinct *original-domain transform*, not entropy-mode microtuning, prior BWT segmentation, or an expanded history search.
- Screen full original executable-like Silesia bytes for plausible relative branch patterns. Verify native forward/inverse exactly; measure a real transformed candidate stream, restoration and frame/decoder costs before any whole-archive attempt.
- Risk: many false branch-opcode coincidences; adds full transform and inverse passes; cannot be accepted merely for better bytes if decoding slows. Revisit only if a novel fast selective gate makes tri-metric improvement credible.

**C. CONDITIONAL/LOW priority: native cheap match-candidate index with bounded depth and reuse.**
- XZ HC/BT and Zstd LDM prove alternatives exist but larger history, dual hashes and sparse absolute-file anchors have already failed for JLZ because of performance or lack of incremental coded-byte evidence.
- Must show a *different match-source representation* rather than a new hash width, chain depth, threshold, file-size gate or window.
- Risk: poor encoding CPU and complicated decoder history. Do not promote general binary-tree matcher or long-range state without a strong strictly better model-cost premise.

**D. LOW priority: Brotli transformed word dictionary and context-map clustering.**
- Potentially useful for natural-language text, but standalone decoder tables, model metadata and extra passes raise complexity, dependency and decode costs. The current selected BWT text files already compress well. Past JLZ context/grammar implementations yielded inadequate cost/benefit.
- Need a sample-only compelling and truly different reversible component use case before more coding; never port/import a dictionary codec.

**E. Stop gate: multiple “new” versions of the same optimal-parsing graph are NOT separate theories.** Zstd, XZ and Brotli all contribute to one state-aware cost-parsing hypothesis. Don't implement three similar optimizers or tune graph width repeatedly. If a bounded graph fails, record the concrete outcome and move to BCJ or another genuinely different domain.

## Design experiment boundaries and approval

No candidate has been integrated or benchmarked from this source audit. Do not report a completed compressor improvement, ten-run timing, percentage-delta comparison table or acceptance. The current codec and source on `main` remain unchanged.

A promotable candidate must be strictly nonregressing in **all three** complete original twelve-file Silesia metrics: complete packed archive bytes, candidate ten-run median encode and candidate ten-run median decode, compared with retained previous current medians. No unapproved ratio-for-speed or speed-for-ratio trade. Any prototype must preserve exact original input bytes, all metadata, backwards decoding, automatic two-command UX and appropriate Go tests/race/vet/crosscompile. This remains true even when a competitor's higher compression level deliberately spends more CPU to save bytes.

**Audit verdict:** sufficient genuine outside architecture remains, but the largest gap is **decision quality at existing JLZ match candidates**, not adding another entropy mode or window size. Before permanent code, demand actual measured candidate-position evidence, not a broad theory story.
