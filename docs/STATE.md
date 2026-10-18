# State and next steps

**Last updated:** 2026-10-18. **Status:** Round 25 automatically sampled native byte-lane representation **accepted**, source updated on `main`. **Branch policy:** all persistent source and document edits directly on `main`; temporary `jls-test/*` branches are solely for CI dispatch and deleted automatically. **Approved source:** JLZ Round 25.

## Current benchmark policy and Round 27 decision (2026-10-18)

**Effective immediately by user instruction:** **Performance research/benchmark only on the complete 12-file Silesia corpus, 211,938,580 original bytes.** Older external-file reports are history only; do not use external corpora again. **Benchmark candidate only**, retaining fixed approved JLZ, Zstd-3 and XZ-6 results rather than spending time rerunning baselines. Historic Zstd/XZ CLI speed data cannot be directly compared with in-memory JLZ and their decompression timings are unavailable. Future responses must use a **four-row table (Zstd-3, XZ-6, current JLZ, candidate JLZ)** with size, reduction, encode and decode times, a **three-column candidate-vs-current percentage delta table**, and a **prominent green/red H1 accept/reject recommendation**. Reusing historical JLZ medians does not create new paired measurements; note that limitation. The revised contract is in `AGENTS.md` and `docs/ENGINEERING.md`.

**[Round 27](reports/2026-10-18-round27.md): REJECT.** Two isolated native entropy alternatives to the selected BWT/MTF/zero-run stream: adaptive arithmetic was 60,672,443 B but took 10.718 s encode / 5.854 s decode in one screening run; preceding-zero-context dual Huffman streams produced **62,103,195 B** versus Round 25's **63,104,117 B**. Candidate-only ten-run full-Silesia medians: **7.833325 s encoding and 2.946025 s decoding**, compared against the **previously recorded** approved Round 25 medians **7.549920 / 2.839905 s**. The BWT context candidate saves 1.586% bytes but appears 3.754% slower encode and 3.737% slower decode; timings are independent session results, not paired. Both candidates were rejected and **no codec source changed**. New code types were not promoted. Raw candidate runs are in [Round 27 timing CSV](reports/2026-10-18-round27-candidate10.csv).

**Next structural priority:** investigate match literal and distance streams on Silesia only, seeking a substantially larger gain rather than repeated low-yield BWT/lane tweaks.

## Current approved Round 25 result (2026-10-18)

Native, automatically selected residue-class byte-lane entropy coding for interleaved original-region structures, with a single sampled histogram shared across strides 2/4/8. One internal block mode `13`, no user options or third-party compressors; complete byte-serialized archive sizes determine promotion. Canonical [Round 25 research report](reports/2026-10-18-round25.md) and [ten-pair raw measurements](reports/2026-10-18-round25-paired10.csv).

**12-file Silesia:** 211,938,580 input B; unchanged paired Round 23 **63,658,799 B** → Round 25 **63,104,117 B** (0.87% smaller), ten-pair encode median **7.622122 → 7.549920 s** (0.95% faster), decode **2.910318 → 2.839905 s** (2.42% faster). **8-file external corpus:** **32,316,619 → 32,085,884 B** (0.71% smaller), encode **3.053475 → 3.076390 s** (0.75% slower), decode **0.773359 → 0.739511 s** (4.38% faster). Only Silesia `sao` and external SoundFont change archive bytes, saving 554,682 B and 230,735 B respectively. Their per-file encodes slow roughly 51% and 27.5%, while decodes accelerate strongly; other files remain identical.

**Compatibility:** 27/27 old Round 23 archives decoded by the new decoder, and 27/27 new archives roundtrip. 26/27 new archives remain byte-identical and decode on Round 23; its decoder correctly rejects the new `sao` mode 13 archive. Tests, vet and targeted race passed locally. The pre-existing historical Round 21 benchmark discrepancy remains unresolved and must be described separately.

**Research next:** seek larger structural gains in remaining match/literal metadata, diversify the public/licensed corpora, and avoid endless selection-threshold microtweaks. The earlier interleaved-stream experiments in the attempted-ideas inventory used external backends or diagnostic measurements; Round 25 is a newly tested standalone native representation.

## Round 26 rejected research (2026-10-18)

[Round 26 predicted byte-lane report](reports/2026-10-18-round26.md): Four-byte residue streams were optionally coded as XOR/subtraction from the previous same-lane byte before native Huffman coding. This new native transform was tested **only in isolated local source copies**, never promoted to `main`. Ten paired full-Silesia runs saved **124,052 B (0.20%)** but slowed encoding **4.26%**; ten paired external-file runs saved **145,307 B (0.45%)** with essentially unchanged encoding and 0.65% slower decoding. An optimized histogram preflight retained archive bytes but remained **2.16% slower** in six-pair screening. **Rejected for cost/complexity; Round 25 remains approved and code/format unchanged.** Three raw timing CSVs are linked from the report.

**Next experiment:** seek larger structural compression improvements to independently packed mathematical-match literals and distance streams, reviewing the Round 18–20 diagnostics and historical attempted ideas before implementation. Avoid repeatedly tuning the lane-based pipeline for isolated small gains.

## Historical approved Round 23 result (2026-10-18)

Complete uploaded Silesia corpus, 12 files, 211,938,580 B. **Paired ten-run local comparison** on Go 1.23.2 linux/amd64: unchanged Round 21 source **64,505,708 B** → Round 23 **63,658,799 B** (**1.31% smaller**); encode median **6.712879 → 6.351571 s (5.38% faster)**; decode median **2.515717 → 2.330026 s (7.38% faster)**. Expanded mathematical/native and mathematical/raw eligibility limits to 80%, with the existing 2/3-sample and 5% BWT selection requirements unchanged. Mathematically safe early abort skips impossible BWT sample gates. `mr` saves 584,674 B, `x-ray` saves 262,235 B, and 25/27 audit archives are byte-identical. The main per-file tradeoff is `x-ray` encoding **93.66% slower**; overall CPU improves.

**Historical reproducibility warning:** The locally rebuilt unchanged Round 21 compressor yielded 64,505,708 B on these user-supplied files, **not** its separately recorded historical 63,770,766 B. Input hashes are now canonicalized in [Round 23 corpus SHA-256 manifest](reports/2026-10-18-round23-silesia-sha256.txt), but the original research corpus hashes/toolchain were not available to resolve this discrepancy. Never mix the prior-host Round 21 timings with current Round 23 paired medians. See [Round 23 report](reports/2026-10-18-round23.md) and [ten paired whole-corpus runs](reports/2026-10-18-round23-paired10.csv) for details. All 27/27 candidate and reference archives decoded byte-exactly in both Round 21 and Round 23; tests, vet and targeted race tests passed.

## Historical Round 21 performance

Complete Silesia corpus: **12 files, 211,938,580 original bytes**. Approved **Round 21: 63,770,766 complete archive bytes**, saving **2,748,439 B (4.13%)** vs approved Round 16 **66,519,205 B**. Verified final-source **ten paired full-corpus runs**, alternating order: Round 16 vs Round 21 encoding median **12.150943 s → 10.562019 s (13.08% faster)**, decoding **4.843944 s → 4.160373 s (14.11% faster)**. Ten-pair aggregate and 120 per-file encode records, ten-pair decode records and 27-file archive audit are stored beside `docs/reports/2026-10-18-round21.md`.

Historical external reference archives: **Zstd-3 66,216,569 B**, **XZ-6 49,233,340 B**. JLZ is **2,445,803 B smaller than historical Zstd-3 on Silesia only**; historical CLI times (**1.256 s**, **116.989 s**) must not be compared directly with JLZ's in-memory median timings.

## Shipped architecture

One fixed, native Go lossless compressor; `compress` and `decompress` only, no optional levels. Grouped 256 KiB regions with bounded 128 KiB history and independent parallel decode groups. Exact, XOR, and modular-add matches; adaptive worker limits from CPU availability; logarithmic distance slots; native entropy coding. Round 16 added automatically selected operation-conditioned lengths. **Round 21 adds a sampled, automatically selected native BWT + move-to-front + zero-run + Huffman encoding for original 256 KiB regions**, reusing three representative probes and avoiding another full mathematical match pass when BWT wins by at least 5% across two of three sample blocks. The decoder restores the unique BWT sentinel and then inverts MTF and BWT. No external DEFLATE, Zstd, XZ, Brotli, or LZMA backends.

**Compatibility:** The Round 21 decoder reads all approved Round 16 JLZ1 archives and all new BWT archives. Older decoders reject the new BWT block type 12 as expected: **21/27** audited archives remain byte-identical and decode with old versions; **6/27** new BWT files require Round 21. New archives still carry the original SHA-256 integrity check.

**Round 22 historical screening (not promoted):** [BWT gate early rejection](reports/2026-10-18-round22-gate-screen.md) was proven selection-equivalent over 262,144 synthetic length combinations and saved 13.1% of hypothetical BWT sample encodes in that grid. At that point, only 11/12 independent Silesia source files were accessible; the complete uploaded ZIP later supplied all twelve and informed approved Round 23. Independent suffix-index/serialization medians on `osdb`'s first/middle/last sampled regions totaled 29.070 ms, a component diagnostic only, with no end-to-end JLZ claim. Measurements are linked from the Round 22 report.

## Next research direction

**Rounds 21 and 23 are documented historical baselines; Round 25 is approved.** Full source: `bwt_native.go`, `bwt_native_test.go`, `region_first.go`, `cross_region.go`. Canonical current comparison: `docs/reports/2026-10-18-round25.md`; keep the older Round 21/23 paired reports unchanged. All permanent edits go directly to `main`, except automatically deleted `jls-test/*` CI trigger branches.

**Known costs:** in historical Round 21, `osdb` had a **13.9% slower** per-file median encoding time, with an identical archive; `ooffice`, `samba`, `mozilla`, and `sao` also have smaller timing regressions despite improving aggregate CPU use. Round 23's worst measured per-file slowdown was **`x-ray` (93.66% slower)** despite improving aggregate size and speed. These are measured targets, not excuses to overfit to Silesia. Test wider licensed corpora before stable-releasing new heuristics.

**Engineering caveat:** Standard-library `index/suffixarray.Index.Write` does not publicly expose sorted suffix offsets; the new encoder parses its implementation-specific, 10-byte-prefixed serialization, with explicit header and bounds validation. Unit tests and CI toolchain checks must catch format changes. The compressed *JLZ1* wire format is independently decoded, but only replace the suffix enumerator if a simpler safe alternative is demonstrably needed. Maintain YAGNI, comments on non-obvious internals and no external compression backends.

Prior negative screens remain canonical in Round 17–20 reports; see `docs/RESEARCH.md` and `docs/research/attempted-ideas.md` before new theory proposals. Round 24 [heterogeneous source screen](reports/2026-10-18-round24-corpus-screen.md) verified 8/8 roundtrips on non-Silesia data. Forced native BWT was **1,078,889 B (3.34%) larger** than Round 23 automatic selection in aggregate; SoundFont was a single smaller but costlier BWT exception. **No codec promotion in Round 24.** Build workflow now includes source/toolchain/version-keyed package caching while rerunning tests, vet, race and integrity validation for every workflow invocation. Future research should prioritize new structural representations, not marginal BWT gates.

## Warnings and handoff

- Historical MDC benchmarks were affected by DEFLATE; never claim transform-plus-external-codec success as standalone JLZ.
- Round 11 eight-candidate matching: 0.88% smaller but 11.01% slower. Rejected.
- Round 13 sparse long-range matching: 0.63% smaller but 5.75% slower. Rejected.
- Round 14 literal byte-pair grammar: 0.12% smaller but roughly 10% slower in screening. Rejected.
- Round 15 conditional skipping of redundant coding passes: 5.30% faster with unchanged archive bytes. Accepted.
- Round 16 operation-conditioned length encoding: 0.42% smaller, 4.21% faster encoding and 3.31% faster decoding. Accepted with backwards-decoder incompatibility noted.
- **Round 21** automatically gated native BWT/MTF/zero-run/Huffman original-region coding: 4.13% smaller, 13.08% faster encoding, 14.11% faster decoding in complete paired tests. Accepted; six new-format Silesia archives need Round 21 decoder.

Before ending any experiment, canonicalize reports and raw results in `docs/reports/`, update this file and `docs/RESEARCH.md`, and preserve previous approved benchmarks. Follow YAGNI, comment non-obvious intent and invariants, and never commit executables or redistribute copyrighted corpora.
