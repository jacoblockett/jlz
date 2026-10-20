# Sparse content-defined far-reference anchors: reversible diagnostic, not promoted

**Recorded:** 2026-10-19 (America/New_York). **Strategy status: SCREENED; no integrated JLZ candidate; current encoder unchanged.** This document deliberately uses strategy names and experimental outcomes rather than a numbered "round" label.

## Hypothesis and distinction from previous failures

A previous reversible prototype widened the mathematical matcher's continuously searched history from 128 KiB to 512 KiB. That recovered 267,104 archived bytes in one complete-original-Silesia integration but substantially worsened both encode and decode time; it was rejected. The new question is different: **can sparse content-defined anchors find far-away *long* exact references without expanding the expensive per-byte mathematical match search?**

A separate prior experiment tried *aligned* 4/16/64/256 KiB exact chunk duplicates. Scanning every complete original Silesia source file revealed only 4 KiB aligned duplicates of modest coverage and **zero 16 KiB aligned exact duplicate chunks** at the examined 16 KiB grid. Therefore aligned macrochunk-only backreferences have poor potential; the sparse-anchor test was intended to be insensitive to chunk grid alignment and allow variable-length copies.

Architecture inspiration, not code or codec reuse: long-distance LZ dictionaries and match-finder tradeoffs described by the liblzma API; chunk/framing separation in Snappy and LZ4. No external mature compressor ran inside the prototype.

## Actual original standalone implementation

Implemented disposable standalone **Go** diagnostic at /mnt/data/jlz_strategy_trials/sparse_far.go (not committed to production). It loads each of the **complete original twelve Silesia files**, scans all byte offsets for an eight-byte hash-derived content anchor (roughly 1/128 selection in random content), stores the most recent occurrence by exact eight-byte word, and seeks prior occurrences **strictly farther than 128 KiB**. A prospective far match requires **16 verified bytes** followed by exact eight-byte extension, maximum 32 KiB per match, and at least **128 matching bytes**; at most 127 bytes of backwards extension are permitted. Non-overlapping output intervals become a reversible sequence of **literal gaps plus (previous absolute offset, exact length) backreference operations**.

The reference stream was **actually decoded**, using reconstructed earlier output bytes (including overlapping byte-copy semantics), and verified byte-for-byte against each original file. All 12 exact reconstructions passed. Complete original source SHA-256 verified **12/12** against the canonical Silesia manifest; source total **211,938,580 B**.

Important limitations:
- It is *not* a complete JLZ format: there is no transmitted archive header, frequency table, block selection, dictionary, offset-mode byte or checksum, and its literal gaps are not entropy coded. **No JLZ compressed-size claims follow from these counts.**
- It does not subtract redundancy **already encoded by JLZ's approved near-history match/BWT/predictive streams**. "Covered source bytes" are a gross, optimistic upper bound, not extra compressed-byte savings.
- On highly repetitive eight-byte prefixes the scan spends considerable work rechecking false candidates; this is a real efficiency problem, not a proof of the absence of rare long matches.
- The current production encoder uses parallel, independently decodable region/history groups. Absolute far-file references would add reconstruction ordering or dictionary side-information unless a separate independently decodable mechanism is demonstrated.

## Full-original-file diagnostic results (not a candidate performance benchmark)

Actual nonoverlapping exact far-match coverage and reversible operation counts:
- mozilla: **988,266 source bytes**, **1,846** far-copy operations; standalone anchor diagnostic about **2.00 seconds**, a major burden against JLZ's current **5.494-second full-corpus** accepted encode reference.
- samba: **163,490 B**, **393** operations.
- ooffice: **33,812 B**, **23** operations.
- mr: **1,763,120 B**, **6,358** operations, but **mr is BWT-selected** in current JLZ, so most of this is not an accessible improvement to the targeted mathematical pipeline.
- xml: **6,181 B**, **15** operations, also BWT-selected.
- dickens, nci, osdb, reymont, sao, webster, x-ray: **0 qualifying operations** with this concrete anchor/index policy.
- Aggregate **2,954,869 source bytes** in verified reversible far-copy spans; of these, the three currently *math-selected* files (mozilla/ooffice/samba) contribute **1,185,568 source bytes before eliminating overlap with existing compressed representations**.

This was a single component/diagnostic execution on the complete original files, **not** ten complete JLZ encode/decode archive passes. No compressed-size or JLZ encode/decode-time comparator table is warranted.

## Outcome, stop gate, and next architecture work

**SCREENED / DO NOT PROMOTE:** content-defined long-range references are a genuine, reversible and mechanistically distinct alternative to a widened conventional match window. They find real source-level far repeats in math-selected files but at current scan cost and without evidence of *incremental native compressed-byte savings after existing paths and new metadata*. Promoting an absolute-offset decoder would complicate independent group decoding. Do not introduce any codec mode, archive claim or new production branch from this experiment.

**Stop adjacent variants:** no ad-hoc anchor modulus/length/threshold tuning or repeated far-window expansion merely to manufacture gains. Revisit only after a separate structural mechanism demonstrates (a) exact frame-costed native savings *over the already selected math representation*, and (b) far-match lookup and decode costs compatible with the **strict three-metric no-regression rule**.

External architecture survey supporting future distinct investigations:
- liblzma public match-finder API distinguishes hash-chain vs binary-tree matching with differing speed/ratio behavior: https://github.com/tukaani-project/xz/blob/master/src/liblzma/api/lzma/lzma12.h (specific file paths may move); https://github.com/kobolabs/liblzma/blob/master/src/liblzma/api/lzma/lzma.h
- Apple LZFSE describes LZ sequence generation and finite-state coding, much of which JLZ already implements; do **not** rediscover existing op/length/distance separation.
- libdeflate documents lazy and minimum-cost-path parsing; the costly near-optimal path is different from JLZ's current one-byte lazy parser but would need a credible bounded-work plan before experimentation.
- Snappy's framed chunk design motivates grouping/independence and checksums rather than treating large dictionaries as free.

Next priority must be structurally new, not an unprincipled sequence of nearby anchor masks, match-history lengths or post-BWT context splits. Preserve historical results by **strategy and recorded outcome**, linked from the central ledger.