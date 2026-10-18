# State and next steps

**Last updated:** 2026-10-18. **Status:** Round 21 automatically gated native BWT regions **accepted**, source updated on `main`. **Branch policy:** all persistent source and document edits directly on `main`; temporary `jls-test/*` branches are solely for CI dispatch and deleted automatically. **Approved source:** JLZ Round 21.

## Accepted performance

Complete Silesia corpus: **12 files, 211,938,580 original bytes**. Approved **Round 21: 63,770,766 complete archive bytes**, saving **2,748,439 B (4.13%)** vs approved Round 16 **66,519,205 B**. Verified final-source **ten paired full-corpus runs**, alternating order: Round 16 vs Round 21 encoding median **12.150943 s → 10.562019 s (13.08% faster)**, decoding **4.843944 s → 4.160373 s (14.11% faster)**. Ten-pair aggregate and 120 per-file encode records, ten-pair decode records and 27-file archive audit are stored beside `docs/reports/2026-10-18-round21.md`.

Historical external reference archives: **Zstd-3 66,216,569 B**, **XZ-6 49,233,340 B**. JLZ is **2,445,803 B smaller than historical Zstd-3 on Silesia only**; historical CLI times (**1.256 s**, **116.989 s**) must not be compared directly with JLZ's in-memory median timings.

## Shipped architecture

One fixed, native Go lossless compressor; `compress` and `decompress` only, no optional levels. Grouped 256 KiB regions with bounded 128 KiB history and independent parallel decode groups. Exact, XOR, and modular-add matches; adaptive worker limits from CPU availability; logarithmic distance slots; native entropy coding. Round 16 added automatically selected operation-conditioned lengths. **Round 21 adds a sampled, automatically selected native BWT + move-to-front + zero-run + Huffman encoding for original 256 KiB regions**, reusing three representative probes and avoiding another full mathematical match pass when BWT wins by at least 5% across two of three sample blocks. The decoder restores the unique BWT sentinel and then inverts MTF and BWT. No external DEFLATE, Zstd, XZ, Brotli, or LZMA backends.

**Compatibility:** The Round 21 decoder reads all approved Round 16 JLZ1 archives and all new BWT archives. Older decoders reject the new BWT block type 12 as expected: **21/27** audited archives remain byte-identical and decode with old versions; **6/27** new BWT files require Round 21. New archives still carry the original SHA-256 integrity check.

## Next research direction

**Round 21 is approved.** Full source: `bwt_native.go`, `bwt_native_test.go`, `region_first.go`, `cross_region.go`. Canonical signed-off comparison: `docs/reports/2026-10-18-round21.md`. All permanent edits go directly to `main`, except automatically deleted `jls-test/*` CI trigger branches.

**Known costs:** `osdb` has a **13.9% slower** per-file median encoding time, with an identical archive; `ooffice`, `samba`, `mozilla`, and `sao` also have smaller timing regressions despite improving aggregate CPU use. Treat this as a measured gate-overhead target, not an excuse to overfit only to Silesia. Test wider licensed corpora before stable-releasing new heuristics.

**Engineering caveat:** Standard-library `index/suffixarray.Index.Write` does not publicly expose sorted suffix offsets; the new encoder parses its implementation-specific, 10-byte-prefixed serialization, with explicit header and bounds validation. Unit tests and CI toolchain checks must catch format changes. The compressed *JLZ1* wire format is independently decoded, but only replace the suffix enumerator if a simpler safe alternative is demonstrably needed. Maintain YAGNI, comments on non-obvious internals and no external compression backends.

Prior negative screens remain canonical in Round 17–20 reports; see `docs/RESEARCH.md` and `docs/research/attempted-ideas.md` before new theory proposals. Future research should prioritize practical generalization, CPU gate decisions and major remaining literal/match structure, not marginal entropy-model partitions.

## Warnings and handoff

- Historical MDC benchmarks were affected by DEFLATE; never claim transform-plus-external-codec success as standalone JLZ.
- Round 11 eight-candidate matching: 0.88% smaller but 11.01% slower. Rejected.
- Round 13 sparse long-range matching: 0.63% smaller but 5.75% slower. Rejected.
- Round 14 literal byte-pair grammar: 0.12% smaller but roughly 10% slower in screening. Rejected.
- Round 15 conditional skipping of redundant coding passes: 5.30% faster with unchanged archive bytes. Accepted.
- Round 16 operation-conditioned length encoding: 0.42% smaller, 4.21% faster encoding and 3.31% faster decoding. Accepted with backwards-decoder incompatibility noted.
- **Round 21** automatically gated native BWT/MTF/zero-run/Huffman original-region coding: 4.13% smaller, 13.08% faster encoding, 14.11% faster decoding in complete paired tests. Accepted; six new-format Silesia archives need Round 21 decoder.

Before ending any experiment, canonicalize reports and raw results in `docs/reports/`, update this file and `docs/RESEARCH.md`, and preserve previous approved benchmarks. Follow YAGNI, comment non-obvious intent and invariants, and never commit executables or redistribute copyrighted corpora.
