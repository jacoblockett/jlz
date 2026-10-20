# JLZ current state and next research

**Status: `current` accepted on `main`. `candidate`: NONE.** New JLZ states and research reports do **not** use numbered rounds. Historical numbered filenames remain only as archive references; strategies and outcomes are indexed in [the strategy ledger](research/strategy-ledger.md).

## Current complete-Silesia benchmark

- **Original input:** the exact original 12-file Silesia corpus, **211,938,580 B** total; manifest at `docs/reports/2026-10-18-round23-silesia-sha256.txt`, all 12 SHA-256 verified.
- **Complete archived size:** **60,585,717 B**. Reduction from original ~**71.41355%**.
- **Previously established candidate-only 10 complete-run medians:** **5.493743735 seconds compress; 2.194427893 seconds decompress**. These are retained accepted timings, **not rerun**.
- **Historical external reference:** Zstd-3 **66,216,569 B**, **1.2561 s encode / 0.8551 s decode**; XZ-6 **49,233,340 B**, **116.9892 s encode / 2.3079 s decode**. External CLI elapsed times are not directly comparable to in-memory JLZ timings.

**Binding promotion rule:** unless the user explicitly approves a particular trade, the **candidate** must not worsen **any** of compressed bytes, median compression seconds, or median decompression seconds on the entire original twelve-file Silesia corpus. Exact size includes all headers, table costs, metadata and checksums. Candidate-only ten complete runs before promotion; clear unfavorable component or one-run screens may be rejected earlier. Disclose per-file regressions, compatibility, repeatability and memory/format costs. Tiny unpaired timing differences are inconclusive, not evidence of significant performance gain.

## Current implementation

JLZ is a standalone native exact-byte Go encoder/decoder with exactly two commands, `compress` and `decompress`, one automatic selection policy, and no third-party compressor backend. The encoder combines existing native mathematical exact/XOR/ADD references, prediction, arithmetic/model-coded entropy and an optional BWT/MTF/bijective zero-run + finite-state representation. Mathematical and BWT source regions are **256 KiB**, with a bounded **128 KiB** cross-block reference history and independently decodable four-block history groups. The BWT winner can encode one entropy model or an MDL-predicted 2-/4-part local model (internal mode16). Rigid decoder bounds, prior archive modes and backward decoding are retained.

**What is NOT in production:** wider 512 KiB BWT blocks, enlarged 512 KiB mathematical search history, dual-anchor exact-only parser, high-cost conditional multi-model BWT entropy, expensive variable entropy boundaries, alternative packed LF inverse, and the new sparse-far-reference diagnostic. Their tested outcomes are preserved by **strategy**, not hidden as "past rounds".

Source control is GitHub `jacoblockett/jlz` `main` only. No hosted GitHub Actions or temporary workflow branches. Use disposable local experimental directories and local Go tests.

## Latest evidence, no benchmarked candidate

**Sparse content-defined far-match diagnostic**, all twelve original source files, byte-exact restoration of a separate literal-gap/past-offset copy prototype: found **988,266 B** of far-copy-covered source in mozilla, **33,812 B** in ooffice and **163,490 B** in samba (math-selected); plus **1,763,120 B** in BWT-selected mr and **6,181 B** in BWT-selected xml. These are **source byte spans, not JLZ compressed-byte savings**, with no complete JLZ archive integration. The prototype consumed ~2 seconds just scanning mozilla, and absolute far-file references would compromise independent group decoding unless solved differently. Status **SCREENED, NOT A CANDIDATE**; no benchmark comparison table warranted. Full [diagnostic report](reports/2026-10-19-sparse-far-anchor-screen.md).

A prior widened conventional mathematical matcher improved archived bytes but slowed encoding and decoding; a separate aligned macrochunk screen found no identical 16 KiB aligned chunks. These were reviewed before pursuing content-defined anchors. Adjacent history sizes, anchor masks, or record lengths are NOT the next research queue.

## Next research: different architecture, not redundant tuning

1. **Profile/research direct suffix-array ordering extraction in the current BWT implementation.** It presently constructs the standard Go `index/suffixarray`, serializes it with `Index.Write`, then parses varints to recover suffix order; a truly faster, native and independently reversible suffix-array construction/extraction mechanism might keep *identical archives* and decode time while reducing compression time. Screen on actual selected original BWT inputs against the existing full-transform time before attempting integration. Avoid copying existing compressor code, Go unsafe reflection, or compiler-dependent internal links.
2. **Cost-aware candidate parsing on the math side**, drawing structurally on libdeflate's minimum-cost parsing and liblzma match-finder classes, but only if a bounded-work approach complements JLZ's XOR/ADD references and can plausibly beat existing 5.494 s encoder. The previously rejected fast exact-only replacement grew archives, so a parser-only swap is not a fresh mechanism.
3. **Re-open genuinely different experimental mathematical/scientific ideas** from `docs/research/theory-backlog.md` when not merely a disguised rerun of an existing token/context/boundary strategy. A theory's failure in one concrete codec is not a universal impossibility.

**Anti-spiral gate:** before code, name the closest failed strategy in [strategy-ledger.md](research/strategy-ledger.md), state what structurally differs, and estimate real metadata/decode/CPU costs. If it is only a cutoff, model count, hash width or span tweak, stop and move on.

## Where to look

- [Strategy ledger: accepted/failed/scoped/untested ideas](research/strategy-ledger.md)
- [Theory backlog, active experimental queue](research/theory-backlog.md)
- [Historical concrete approaches](research/attempted-ideas.md)
- [Conceptual research index](RESEARCH.md)
- [Engineering](ENGINEERING.md), [local workflow](WORKFLOWS.md)
- Historical dated experimental reports remain under `docs/reports/`, including older numerical filenames retained solely for traceability.


## Latest suffix-order diagnostics: no benchmarked candidate

**Current remains unchanged; candidate remains NONE.** The current BWT forward transform's standard-library suffix-array serialization layer was measured on **all 451 selected source blocks** of seven original Silesia files: building suffix indexes consumed ~**4.088 s** cumulative stage time; serializing their offsets ~**0.519 s**; validating/parsing serialized offsets ~**0.364 s**. Intermediate serialized data totaled ~**461.7 MB** processed block by block, **not** embedded in the compressed archives. These figures are a **single component timing pass**, not an integrated JLZ compression benchmark or a promised ~0.88s speedup. An independently written radix-prefix-doubling suffix constructor produced exactly the same BWT bytes and primary indexes on **21/21 distributed original Silesia block samples** with exact inverse reconstruction, but was **~2–7× slower** than the current complete BWT forward transform. **REJECT this concrete builder and stop nearby radix/step-size tuning**; lower-algorithmic-complexity direct construction remains possible but unproven and may not justify implementation complexity. [Evidence](reports/2026-10-19-suffix-order-extraction-screen.md). The current Go source was not modified.

**Priority shift:** Having screened macrochunk dedup, sparse far anchors, and generic prefix-doubling suffix building without a tri-metric opportunity, favor **distinct bounded-work mathematical match-parsing models** or scientifically new transformations from the active theorem backlog rather than adjacent BWT or far-anchor tuning. Identify exact historical overlap before coding. No ten-run whole-corpus candidate benchmark was run, hence no candidate comparison table is warranted.


## Source-audited architecture opportunities, no candidate yet

Inspected primary current implementation paths in higher-level Zstd (optimal parser and LDM), XZ/LZMA2 (normal/fast costed parser, HC/BT match finders, reset/raw chunks, x86 branch prefilters) and Brotli (Zopfli-like optimized command paths, context/histogram clustering and dictionary transforms). Compared to JLZ's actual match scorer and present encoded streams and explicitly cross-referenced earlier rejected exact-only, history, context and duplicate-encoding prototypes. The **highest-value distinct gap** is bounded **entropy-informed match *sequence selection*** from **existing** exact/XOR/ADD candidates, without more expensive search or a decoder change. The second independently distinct gap is native reversible executable-relative-address preprocessing, but format/inverse runtime risk is high. [Full primary-source audit](reports/2026-10-19-industry-architecture-audit.md) and [strategy ledger](research/strategy-ledger.md). **No benchmarked candidate; `current` unchanged, no tri-metric improvement claimed.**


## Mathematical match price diagnostic, candidate still NONE

An independently run original Go diagnostic on first/middle/final regions of the complete original Silesia source files **mozilla, ooffice, samba** measured actual emitted match literal streams and encoded them as exact native token archives; **9/9 exact SHA-checked roundtrips**. Seven of nine literal H0 rates were **higher than six bits per literal byte** while the live `engine/match.go` match scorer values replacement at fixed 5/6 bits per byte. This is evidence that JLZ's local surrogate does not reflect actual local stream statistics; it **does not establish a compression benefit**, an integrated candidate, or whole-corpus benchmark timings. [Full component diagnostic](reports/2026-10-19-native-match-price-screen.md). Together with the upstream [Zstd/XZ/Brotli source audit](reports/2026-10-19-industry-architecture-audit.md), prioritize **one** bounded cost-aware selection experiment with existing XOR/ADD/exact candidates, and stop if added pricing exceeds any actual packed gain. **Current accepted 60,585,717 B / 5.493744 s encode / 2.194428 s decode historical medians; candidate NONE.**


## Latest integrated candidate outcomes: price-aware parsing and executable-address normalization (both REJECTED)

Following primary-source audits of Zstd high levels, XZ/LZMA2 and Brotli, actually implemented two original native-Go techniques while preserving current on `main`. **Cost-aware parse selection** retained JLZ exact/XOR/ADD matching but priced already discovered alternatives with literal/opcode/length/distance statistics; all nine sampled complete-original Silesia source regions recovered exactly, and an integrated cheap histogram-only model compressed the full twelve originals to **60,435,710 B** in **one complete screen** at **6.667317 s encode / 2.311219 s decode**. This is smaller than current by 150,007 B but slower in BOTH measures, so **REJECT without ten-run timing**. A separate, structurally different original x86-relative-branch source prefilter passed complete actual native stream component screens (especially ooffice) and was fully integrated with a single conservative executable-density gate and decoder internal mode18. **TEN newly measured candidate-only full twelve-file Silesia passes:** **60,431,550 B / 5.715744 s encode / 2.266683 s decode** medians, 120/120 exact, 12 canonical input hashes verified; against current historical **60,585,717 B / 5.493744 / 2.194428**, real **-154,167 B (-0.254461%) size** entirely in ooffice but nominal **+4.041% encode / +3.293% decode** (unpaired). Go suite/vet/race/Windows compile PASS; candidate decodes current-format archive, older binary rejects new candidate mode18 as expected. **REJECT BCJ despite ratio win**, no source promotion. [Methods and raw data](reports/2026-10-19-entropy-cost-matching-and-x86-branch-normalization.md), [strategy ledger](research/strategy-ledger.md). Both concrete families need an entirely different zero-extra-pass mechanism to revisit; do not vary cost-model knobs or opcode thresholds. **Current remains exactly 60,585,717 B / 5.493744 s encode / 2.194428 s decode. Candidate: NONE.**


## Industry hot-path source audit (2026-10-20), no candidate

Further user-directed source-level study examined **Zstd fast and superblock entropy-estimation, LZFSE match verification and headers, libdeflate observations, LZ4 HC pattern search, Snappy match pipeline, and libbsc LZP/QLFC**. Contrary to repeating general slogans about faster compression, this audit maps concrete source instructions to `engine/match.go`, `engine/finite_state.go`, `engine/stream.go`, `match/expand.go` and the measured rejected strategies. **Two genuinely untested candidates for economical research**: (1) preserve exact XOR/ADD matcher decisions but verify 8 transformed bytes per operation rather than scalar byte loops (LZFSE-style mismatch extraction, original modular ADD lanes); (2) eliminate unnecessary complete FSE encodings when a mathematically safe byte-cost gate proves they cannot beat the already fully packed native stream (Zstd-style sub-block estimates). Additional alternatives such as lazy-skipping hash insertions are explicitly ratio-regressive in their original design and do not satisfy user approval. Full [industry source and JLZ gap audit](reports/2026-10-20-industry-encoder-hotpaths-audit.md), [strategy ledger](research/strategy-ledger.md). **Current unchanged: 60,585,717 B / historical 5.493744 s encode / 2.194428 s decode; candidate NONE; no codec benchmark run during this audit.** Further experiments must prove token/bitstream equivalence and tri-metric nonregression, not merely improve a microbenchmark.
