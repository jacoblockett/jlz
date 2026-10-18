# State and next steps

**Last updated:** 2026-10-18. **Status:** Round 17 screening research rejected; approved Round 16 source remains unchanged. **Branch policy:** all persistent source and document edits directly on `main`; temporary `jls-test/*` branches are solely for CI dispatch and deleted automatically. **Approved source:** JLZ Round 16.

## Accepted performance

Complete Silesia corpus: 12 files, 211,938,580 original bytes. Current Round 16: **66,519,205 total archive bytes** versus Round 15 **66,797,408**. Ten paired in-memory full-corpus runs on the same test host: encoding **13.290 s vs 13.874 s** (4.21% faster), decoding **5.311 s vs 5.493 s** (3.31% faster). Historical external archive totals: Zstd-3 **66,216,569**; XZ-6 **49,233,340**. Historical CLI timings are not directly comparable to current JLZ in-memory medians. The complete benchmark and raw paired run records are in `docs/reports/2026-10-18-round16.md` and adjacent CSV files.

## Shipped architecture

One fixed, native Go lossless compressor; `compress` and `decompress` only, no optional levels. Grouped 256 KiB regions with bounded 128 KiB history and independent parallel decode groups. Exact, XOR, and modular-add matches; adaptive worker limits from CPU availability; logarithmic distance slots; native entropy coding. Round 16 adds automatically selected operation-conditioned lengths: literal-run varints and copy varints use separate compressed streams when smaller. No external DEFLATE, Zstd, XZ, Brotli, or LZMA backends.

**Compatibility:** Updated Round 16 decoder reads prior JLZ1 archives. Older JLZ decoders cannot read Round 16 archives containing new split-length block types. This is documented and tested; keep all new archive interpretations explicit.

## Next research direction

Round 17 tested zero-order source-entropy calibrated matching (46,666 B smaller, but a slower single-run screen) and a small-count probability lookup (byte-identical archives, no speed benefit). Both were rejected *before* ten-run acceptance benchmarking. Canonical details: `docs/reports/2026-10-18-round17.md`. Do not misreport screening times as medians or promote either prototype.


Seek a genuinely **unified and economical match/token/entropy cost model derived from actual emitted token streams**, using the operation-conditioned length result as evidence that token structure contains redundancy. Avoid repeatedly parsing/re-encoding with expensive second passes, adding unused format modes, or going deeper into previously rejected match searches. Screen full representative files before ten paired Silesia runs. A gain under 3% still requires skeptical cost/complexity review.

## Warnings and handoff

- Historical MDC benchmarks were affected by DEFLATE; never claim transform-plus-external-codec success as standalone JLZ.
- Round 11 eight-candidate matching: 0.88% smaller but 11.01% slower. Rejected.
- Round 13 sparse long-range matching: 0.63% smaller but 5.75% slower. Rejected.
- Round 14 literal byte-pair grammar: 0.12% smaller but roughly 10% slower in screening. Rejected.
- Round 15 conditional skipping of redundant coding passes: 5.30% faster with unchanged archive bytes. Accepted.
- Round 16 operation-conditioned length encoding: 0.42% smaller, 4.21% faster encoding and 3.31% faster decoding. Accepted with backwards-decoder incompatibility noted.

Before ending any experiment, canonicalize reports and raw results in `docs/reports/`, update this file and `docs/RESEARCH.md`, and preserve previous approved benchmarks. Follow YAGNI, comment non-obvious intent and invariants, and never commit executables or redistribute copyrighted corpora.
