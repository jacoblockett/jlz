# JLZ engineering contract: current and candidate

**Current** is the accepted, native Go compressor source on GitHub `main`; **candidate** means a disposable, unapproved experiment. Do not create new numbered research rounds or name the active compressor after old numbered reports. Read the [current state](STATE.md), [strategy ledger](research/strategy-ledger.md), [architecture](ARCHITECTURE.md) and [local workflows](WORKFLOWS.md) before edits. Historic dated reports retain their archival identifiers as evidence only.

## Fundamental restrictions

- Exact original input bytes must reconstruct, including checkable archive integrity, with appropriate corruption/truncation rejection and documented backwards-decoding support. The current `JLZ1` format preserves earlier modes; historical MDC/MDN formats are not silently accepted.
- Native, standalone Go codec. No DEFLATE, Zstd, XZ/LZMA, LZFSE, Brotli, gzip, zlib or other mature compressor backend, wrapper, nested final encoder, library fallback or imported implementation. Study their architectures; implement original native mathematical ideas. Go standard library for ordinary utilities is allowed (including generic `index/suffixarray`), not Go's preexisting compression codecs in JLZ production.
- One automatic compressor, exactly `compress` and `decompress` CLI verbs. Internal automatically selected modes are permitted with evidence; no levels, extra user switches, undocumented model dependencies or ad-hoc experimental archive syntax.
- YAGNI: only validated, necessary source changes. Favor small additive components; wholesale representation replacement needs extraordinary complete-corpus evidence. Preserve decoding of historical `JLZ1` modes where practical. Explain new decoder boundaries and any forward-incompatible format if needed.
- All durable code and documentation edits land directly on `main` only after validation. Work in disposable local folders for uncertain ideas. GitHub Actions are **disabled** and must not be dispatched, restored or triggered; there are no test/release branches under current workflow.
- Workers adapt to CPU quotas and available blocks; bound memory, avoid uncontrolled goroutine/allocation growth and data races. Comment algorithmic invariants, bitstream layouts and dangerous overflow/buffer bounds rather than tautological code.

## Strict measurable promotion gate

**The user explicitly requires a candidate to be non-regressing on ALL THREE performance dimensions; no unilateral tradeoffs:**

1. **Full compressed archive bytes**, including every stream model, frame, checksum, side-information header and fallback decision across the complete original twelve-file Silesia corpus (**211,938,580 source bytes**).
2. **Median compression wall time** across **ten candidate-only complete twelve-file runs**, compared to the retained previous approved complete-corpus median.
3. **Median decompression wall time** across the same ten complete runs, against its retained accepted reference.

No tradeoff between compressed size and either speed, or between encoding and decoding speed, is approvable without the user's explicit permission. A candidate preserving all three can still be rejected if it introduces disproportionate code, memory, forward incompatibility, per-file regressions or security complexity. Reject early from obviously unfavorable component/one-pass screens rather than run unnecessary complete ten-pass benchmarks. Numerical differences in independent, unpaired timing sessions may be noise; do not claim substantial speed improvements without supporting repeated evidence.

**Current approved ten-run numbers:** full original Silesia JLZ archive **60,244,257 B**, compression **5.081535442 s**, decompression **1.9046838445 s**; previous accepted historical 60,244,257 B / 5.322679488 / 2.003659634 s (not rerun); Zstd-3 **66,216,569 B / 1.2561 s / 0.8551 s** and XZ-6 **49,233,340 B / 116.9892 s / 2.3079 s** (historical external CLI with different timing convention). **Do not rerun current, Zstd or XZ benchmarks** just to measure a candidate. Verify original source SHA-256 against the canonical Silesia manifest.

Only the complete original twelve-file Silesia is authorized for codec performance measurement. Unit/synthetic fixtures can test reversibility, malformed input and regression correctness, **not claim benchmark superiority**. Selected regions from complete originals are diagnostic components only.

## Validation before a production change

- `go test -count=1 ./...`, `go vet ./...`, affected `go test -race -count=1` targeted to concurrent/data-sensitive paths; state any tests not run.
- Source identity: all twelve Silesia canonical SHA-256 checks; complete candidate compression and decompression exact bytes; include embedded archive integrity checks and corrupt/truncated decode tests.
- Ten complete **candidate-only** source-verified runs if plausible for promotion, report independently computed medians and per-file byte totals, worst regression, memory and decoder dependencies.
- At least one existing prior `JLZ1` archive decoded with the candidate. New format/wire subtype is not readable by earlier binaries unless proven; never imply symmetric forward/backward compatibility.
- `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build` for Windows cross-platform compilation; do not describe a successful cross-build as proof of native Windows runtime.
- Any candidate that fails size or either speed metric remains disposable. Preserve its specific outcomes and reproducible method in repository research docs, **not rejected Go code in production**.

## Current decoder/wire invariants

- Native packed-stream entropy supports raw, RLE, canonical Huffman and selected independently framed normalized finite-state (12-bit frequency total **4096**) modes; validate every table sum, nonzero frequency, initial/final state and exact decoded stream sizes.
- BWT-selected original source blocks are **256 KiB**. Current byte BWT uses bijective MTF zero-runs with independently normalized entropy and an optional MDL-predicted 2-/4-part neighborhood representation (block mode **16**, with its own independently framed entropy segments). Earlier BWT MTF mode **12** and bijective mode **15** remain decodable; packed-stream finite-state subtype **5** remains native.
- **No production 512 KiB BWT blocks.** Wider BWT source frames were expressly rejected for slowing both encode and decode; current mixed-cross decoder restricts original block size to 256 KiB.
- Mathematical match references support exact bytes and XOR/ADD transformed relationships with three repeat-distance cache entries. Normal cross-group input regions are **256 KiB** with bounded lookback **128 KiB**, and sequential references within independent four-block groups. More history was experimentally slower and is **not** an accepted decoder behavior.
- BWT preflight must compare actual **bijective** block representation, not an obsolete MTF proxy; the current gate requires two of three independent representative wins and an approximately 2% sample margin, with actual payload sizes deciding final output.
- The Go BWT forward suffix-order extractor currently uses the ordinary `index/suffixarray` library and parses its serialized `Index.Write` index; maintain the existing validation of that Go format until any demonstrably faster, equivalent and portable original extraction mechanism is accepted. A tested radix-prefix-doubling constructor was correct but slower, and is not production code.

## Research reporting

Read the [strategy ledger](research/strategy-ledger.md) before proposing a change. Explain novelty **relative to the closest tested failed implementation** and distinguish an abstract hypothesis from a concrete coding method. Detect and stop parameter-tuning spirals. New results are documented under `docs/reports/YYYY-MM-DD-<strategy-name>.md`, with ledger and state updated after each substantive experiment. No new numbered round identifiers.

User-facing response: concise method, implementation, actual reconstruction, result, decision, next structurally distinct opportunity. **Only if a real integrated candidate was benchmarked** use the full comparison table (Zstd, XZ, current, candidate; archived bytes, original reduction, encode/decode seconds), followed by a three-column candidate-versus-current percentage-delta table. No repetitive baseline-only or diagnostic-only tables. Historical figures and new run counts must always be distinguished.


## Placement-sensitive validation (2026-10-20)

The supplementary layer is the default design hypothesis; a previously failed standalone transform must not be dismissed as generally useless until its interaction with current separately entropy-coded native streams is examined for mathematical compatibility, work and metadata cost. Keep a concrete [source-to-substream placement matrix](research/layer-placement-matrix.md). Trial a reversible component with the existing stage's actual native coder, model metadata and inverse; if it produces a clear net cost loss, close only that concrete placement. If profitable, integrate **within existing JLZ automatic selection and all original source files**, accounting for cross-block history and decoding, before claiming whole-codec gain. A real ten-run BCJ-before-existing-BWT candidate tested 2026-10-20 saved 341,460 compressed bytes and sped decoding but worsened encoding median **5.493744→5.557846 s** nominally; **REJECT, current unchanged**. See [actual run data](reports/2026-10-20-layered-architecture-reassessment.md). Do not accept a two-of-three performance gain without explicit user's tradeoff permission. Failed work and genuinely untested placements must remain distinguishable in the ledger.


## Accepted layered BCJ→BWT plus Go PGO (2026-10-20)

**Current** uses reversible per-block E8/E9 32-bit little-endian relative operand normalization before existing native 256KiB BWT/MTF/finite-state coding on automatically selected math regions. Mixed-cross mode 19 holds existing BWT subtype 15 or 16. The decoder reconstructs via existing BWT inverse then branch inverse; old archive modes remain backward-decodable, while old binaries reject newly written mode 19. A root checked-in **default.pgo** profile (16,804 compressed bytes; SHA-256 3b35cda74832360a252b63c1805e025aa1f4ed7a4fe4124d9a55929a6e86289f) is automatically read by ordinary Go 1.23+ compilation. Without the measured profile the same codec was historically REJECTED for encoding regression; **do not claim the approved timings for an unprofiled build**.

Current **12-file original Silesia** accepted complete archive: **60,244,257 B**; candidate-only ten-run medians **5.322679488 s encode / 2.003659634 s decode**, 120 exact reconstructions. Previous accepted historical (not rerun): 60,585,717 B / 5.493743735 s / 2.194427893 s. Net -341,460 B, nominal -3.1138% encode, -8.6933% decode (independently timed sessions). [Actual ten-run data, method and backward-compatibility checks](reports/2026-10-20-native-bcj-bwt-pgo-accepted.md). No pending candidate and no numbered release label.


## Accepted native speed-only layers (2026-10-20)

The source-aligned combined candidate was promoted with unchanged complete-original-Silesia archive content **60,244,257 B**, ten complete candidate-only medians **5.081535442 s encode / 1.9046838445 s decode**, previous accepted medians historical **5.322679488 / 2.003659634 s**. A packed low32/high32 two-entry math match bucket replaces its 12-byte ring structure without changing selection or tokens. A two-hop BWT LF inverse constructs pair symbols plus next-two-row index in temporary 64-bit entries, reducing dependent decoder traversals at the cost of ~2MiB more memory per full 256KiB BWT block. All old archive modes remain readable, no new format subtype, and root PGO unchanged. Full evidence and all tests: [accepted report](reports/2026-10-20-compact-index-two-hop-inverse-accepted.md). This new section supersedes any former "current" timing in historical sections above. Keep strict all-three nonregression promotion gate unless user explicitly authorizes a specific trade.
