# BWT suffix-order extraction without serialization: cost audit and reversible prototype

**Date:** 2026-10-19 America/New_York. **Strategy:** suffix-order extraction/construction. **Status:** baseline-cost DIAGNOSED; concrete native prefix-doubling replacement REJECTED in a component screen. **No integrated JLZ candidate was created, no complete archive benchmark exists, and `current` was unchanged.** This report uses no numbered round identifier.

## Historical overlap / novelty

Existing JLZ BWT stages obtain suffix order using the ordinary Go standard-library `index/suffixarray.New` and its `Index.Write` serialization. BWT then re-parses length-prefixed serialized varint suffix offsets, because the standard-library index does not export its internal suffix array. The current implementation is valid and Go tests pass, but incurs extra serialization and parsing work. Prior attempts to widen BWT blocks improved ratio but slowed both encoder and decoder and were explicitly rejected. Prior packed LF-index attempts targeted **BWT inverse** and never established a tri-metric win; this is distinct, addressing **BWT forward suffix-array ordering** only.

Inspired by general suffix-index construction, not by an existing compression codec or source-code port.

## Full selected-source cost diagnosis

Standalone Go diagnostic used **all 451 current BWT-selected 256 KiB regions** spanning the seven complete original Silesia files (dickens, mr, nci, osdb, reymont, webster, xml), **117,233,324 source bytes**. It independently timed:
- Go `suffixarray.New`: **4.088358535 s** cumulative stage time.
- `Index.Write` standard-library suffix order serialization: **0.518594153 s** cumulative stage time.
- Parsing all index records/varint suffix offsets: **0.363556024 s** cumulative stage time.
- Total these stages: **4.970508712 s cumulative**; serialized intermediate bytes **461,712,816 B**, reused block by block rather than all simultaneously.

All serialized index records were parsed completely and validated that their suffix count equaled the input block size; no original file bytes were altered. **These are one-pass component stage totals, not compression/decompression benchmark times and not a prediction of real JLZ wall-time savings.** In particular, Go's current whole-compressor parallel worker scheduling, entropy/matching costs and concurrent memory pressure prevent simply subtracting 0.882 s from the historical JLZ encode median. Nevertheless, the measured serialization and parsing overhead represents real process work.

The seven selected file sources come unchanged from the original complete twelve-file Silesia corpus; the canonical twelve source hashes verified in this session before testing.

## Independently written reversible alternative: radix prefix-doubling suffix ordering

Implemented a disposable native Go radix-counting-sort **prefix-doubling suffix-array** constructor. It sorts suffix positions by progressively doubled paired ranks with an explicit end-of-suffix rank, then generates JLZ's **exact same** sentinel-aware BWT output without writing or parsing a standard-library index. Unlike a direct high-performance induced-sorting construction, this conventional algorithm may require O(n log n) passes and several arrays.

Compared against the **current Go BWT transform** on 21 first/middle/last samples from the same seven selected complete Silesia originals, including short natural final blocks. Each of the 21 results had:
- An identical primary index.
- **Byte-for-byte identical BWT transformed bytes**.
- Successful inverse restoration of the entire original selected source block.
- Both old and new forward-transform wall-time recorded for that sample.

**All 21 new constructors were slower** than the current Go suffix-array constructor PLUS its serialization/parsing/BWT. Representative samples:
- dickens first 256 KiB: current 15.705 ms, radix-doubling 37.293 ms.
- mr first 256 KiB: current 11.641 ms, radix-doubling 47.749 ms.
- nci first 256 KiB: current 7.414 ms, radix-doubling 47.152 ms.
- webster first 256 KiB: current 15.318 ms, radix-doubling 31.800 ms.
- xml first 256 KiB: current 7.908 ms, radix-doubling 55.533 ms.

The new builder is roughly **2–7 times slower** on those samples (with intermediate smaller block differences also observed). That fails the strict no-encode-time-regression rule decisively before any complete archive candidate needs to be integrated or benchmarked. These are **one-run component screens**, NOT ten complete JLZ runs.

## Outcome and stop condition

**REJECT the concrete prefix-doubling replacement.** No production change, no archive format change, no disabled alternative engine remains in current. Performance diagnostic supports that avoiding `Index.Write` serialization could be valuable *only if suffix order can be constructed equally efficiently without it*, e.g. a fundamentally different induced-order algorithm. It does not support recoding a modern linear-like Go suffix constructor as a generic O(n log n) ranking algorithm.

**No more nearby radix-pass optimization or doubling-factor sweeps**: the construction complexity, not a tiny threshold, dominates. A further attempt would require an independently derived/direct native suffix-order algorithm with credible O(n) behavior and significantly better measured stages, a compact implementation, exact equality on real original blocks and no unsafe reliance on Go internal layouts. If this is not economically justified, move to a genuinely different mathematical-match architecture. Preserve older BWT inverse and larger-block rejections distinctly in the strategy ledger.

Experimental standalone files were deliberately kept disposable. There is **no full codec candidate performance table** and no candidate approval claim.
