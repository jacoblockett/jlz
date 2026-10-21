# Zstd source pre-splitting applied before JLZ BWT, complete selected-block screens

**Date:** 2026-10-20 (US/Eastern). **Status:** SOURCE AUDIT + EXECUTED NATIVE-GO COMPONENT SCREENS. **Rejected concrete fixed-half and Zstd-style border-heuristic source split**, but **general content-defined source splitting remains OPEN**, not claimed exhausted. No integrated JLZ new mode, no ten-run candidate benchmark, no source, archive, or PGO promotion.

## Origin and novelty check

Read Meta's actual [Zstd preSplit implementation](https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_preSplit.c), Git blob `d820c20ac24debf5e2877e217e143339862ea74f`. Zstd compares block-border byte-frequency fingerprints and optionally 8KiB-chunk fingerprints, then splits **source** at candidate boundaries, *before* match discovery and sequence/literal entropy encoding. Its fast default uses 512-byte fingerprints from each end and the center to choose quarter/half/three-quarter boundaries when divergence is sufficiently strong; more expensive modes sample two-byte hashes with 8KiB chunks.

**Closest JLZ prior work, distinct placements:** post-BWT *symbol-stream* segmentation with 2/4 parts (accepted mode16), rejected 8/16 parts, and rejected variable change-point search after bijective MTF; increasing BWT input horizon to 512KiB (rejected); and syntax-class source partition before BWT (rejected on webster). The tested mechanism here instead makes **multiple smaller independent BWT transforms** of consecutive source input sections, with their own sentinel and primary index and native MTF/rANS frames. This changes suffix-array work and transform statistics; post-BWT model splitting cannot emulate it. The prior Zstd level-3 architecture audits covered source paths but **did not run this exact native pre-BWT source split**.

**Architectural limitation:** Zstd's rationale concerns changed local match/literal statistics in LZ blocks; independent smaller suffix-sorting BWT blocks can actually lose large cross-boundary sort context, so Zstd source preSplit's heuristic may not transfer. Independent block-mode costs of two subblocks are optimistically counted as a fixed **8 B extra** beyond their exact existing `encodeBWTSubregions` outputs. This is not an integrated legal JLZ wire mode, and 8 B is a favorable provisional charge, so all size savings below are **component oracle numbers**, not complete archive improvements. More complete metadata would reduce gains. No Zstd codec, parser, implementation, or backend was copied into JLZ.

## Corpus, scope and exact reconstruction

User-provided original source archive `silesia.zip`: all **12 unmodified original Silesia members, 211,938,580 B**, 12/12 full source SHA-256 digests reconfirmed against committed manifest `docs/reports/2026-10-18-round23-silesia-sha256.txt`. Selected actual BWT-original source files: **dickens, mr, nci, osdb, reymont, webster, xml**, covering **451 original 256KiB BWT blocks** (including a final tiny block). **450** blocks at least 32KiB supported this two-section experimental screen; one small final partial block was left at incumbent. `ooffice` additionally uses **branch normalization before BWT**; this initial placement deliberately left that already accepted transform unchanged and does not constitute its screening. All tested `encodeBWTSubregions` legacy-mode15/mode16 frames were decoded using their **unchanged existing JLZ decoders** and compared byte-for-byte with the corresponding entire input segments. No actual on-disk candidate mode created; no ten-run timings.

Verified disposable original Go source against production blobs `engine/bwt_subregions.go` `f80f6c417b4beb58b0a7ad5d525094b6b80b4b18`, `bwt/transform.go` `422a24d8c5f537d240846c0e38d624e2d90d8fb2`; these remained unchanged on live `main` during the trial. The test code is only in disposable `/mnt/data/jlz_zstd_split/engine/zz_zstd_*_test.go`, not committed. Root `default.pgo` unchanged.

## Full all-selected-file component measurements

Compared unsplit source BWT to two separately packed source chunks split at 1/4, 1/2, 3/4 of each 256KiB source block. Oracle chooses the smaller of existing unsplit and the three split sizes at block-level, **without counting its selection search CPU**. Fixed-half simply emits halves for every eligible block.

| Source | Existing unsplit native BWT B | Forced half B | Best-of-three per-block oracle B | Oracle B saved |
|---|---:|---:|---:|---:|
| dickens | 3,152,600 | 3,330,558 | 3,152,600 | 0 |
| mr | 2,722,419 | 2,778,509 | 2,713,778 | 8,641 |
| nci | 2,162,302 | 2,361,143 | 2,162,302 | 0 |
| osdb | 3,424,772 | 3,792,127 | 3,424,772 | 0 |
| reymont | 1,455,011 | 1,527,090 | 1,446,883 | 8,128 |
| webster | 9,834,263 | 10,391,086 | 9,834,263 | 0 |
| xml | 515,676 | 573,014 | 514,943 | 733 |
| **Total** | **23,267,043** | **24,753,527** | **23,249,541** | **17,502** |

Forced equal halves **expanded 1,486,484 B (+6.389%)** against incumbent BWT for these source blocks, clearly REJECT. Only **21/450** blocks had a winning 1/4, 1/2 or 3/4 split, predominantly 17 in mr, two in reymont and two in xml. An initial first/middle/last sample already showed the striking isolated source-class boundary in reymont block24 (98,714→91,609 B, a 7,105 B winning three-quarter split), prompting the exhaustive all-selected-block expansion rather than discarding the mechanism from a small sample.

## Actual Zstd-inspired border-fingerprint selector (refinement, not replacement by oracle)

Measured per-block 512-byte beginning/end/center source byte histograms. An integer Manhattan-distance threshold corresponding to Zstd's fast 7/8 fraction of disjoint histograms flags blocks for split; the center histogram selects quarter/half/three-quarter based on resemblance to the borders. This uses no entropy encoding to decide. Evaluated cheap fingerprints over all **450 eligible** natural original 256KiB blocks.

- Flags **32/450**; **11/21** coarse-oracle winners are among the flagged blocks, but correct proposed quarter only **2** of those 11.
- **Following the predicted split unconditionally**, actual native candidate frames total **+58,720 B (worse)** versus the original unsplit BWT frames, with only **2/32 predicted blocks actually winning**.
- Even a predictor-plus-*encode both alternatives and choose* oracle saves only **8,128 B** among the 32 flagged blocks, requiring a redundant BWT pass on the flagged source input. This is a component bound; the real work overhead of this second BWT was not end-to-end benchmarked.
- More permissive border-distance cutoffs flagged 62 blocks at a cost of additional false positives; 20/21 winning blocks were within them, but the predicted split location still matched the best coarse quarter in only two of the flagged wins. These are **diagnostic threshold checks**, not separate integrated codecs.
- No other data-specific gate was accepted or hardcoded. The original 512-byte fingerprint technique is tuned for *Zstd source LZ splits*, not JLZ's suffix-sorting + bijective MTF objective; its failure here does not discredit its native Zstd use.

## Further theory refinement: 8KiB boundary grid on genuinely winning blocks

The 21 coarse-grid original blocks and three negative source controls were re-tested with independent BWT at **every 8KiB boundary** satisfying at least 8KiB on either side. All fragment inverses were exact. **21/21 initial coarse-grid winners** remained winners; **19** improved with non-quarter boundary. The sum of best oracle savings on those 21 blocks rose from **17,502 B to 29,795 B**, including `reymont` block24 **9,025 B** versus the original unsplit 98,714 B, and `xml` block16 **2,002 B**. The three negative controls from dickens/webster/nci still did not win. **This 29,795 B is not an all-451-block fine-grid oracle.** The 429 coarse-grid nonwinners were not exhaustively re-evaluated at every 8KiB boundary, so larger opportunities at other positions remain logically possible. Do not pretend a complete infinite/all-boundaries search was performed.

This refinement demonstrates an important limitation of the original coarse-grid screen, not a reason to hide the positive cases. However these sampled winners' cumulative upper bound amounts to just about **0.051% of current 58,767,739 B total JLZ archive**, even before the new mode and an encoder that can choose these exact boundaries. A content-change classifier specialized to 8KiB granularity would require meaningful additional engineering and full JLZ testing before any acceptance claim.

## Decision, validation, Zstd research continuation

- **Reject fixed 128KiB halves** as an unconditional BWT replacement: huge measured size regression.
- **Reject source border-divergence routing at the tested 512B/center model**, since its forced choices grow bytes, and even exact fallback against old BWT would involve duplicate expensive sorting for only 8,128 B of select-only component gains.
- **Do NOT globally close content-defined source partitions** on JLZ. One may someday use 8KiB-chunk fingerprinting, secondary matched-parser routing, or shared entropy/dictionary state, none of which is disproved by this BWT experiment. But any successor must explain how it solves model selection/duplicate sorting work, independent metadata, and the tiny already demonstrated opportunity.
- **No integrated candidate qualified for ten complete runs** because actual source-stage bytes already fail on the practical predictor. Do not put hypothetical sizes/timings into a full-corpus candidate table. No new archive mode, no source update.
- Baseline accepted `go test -count=1 ./...` and `go vet ./...` PASS. Prototype `go test -count=1 -run TestZstdSourceSplitFeasibility`, `TestZstdBorderPredictedSplit`, and `TestZstdFineSourceSplitRefinement` PASS with exact native fragment inverses. No separate race test/Windows crossbuild/backward archive decode for an unimplemented wire subtype.

**Reference historical complete corpus (NOT rerun as a benchmark):** JLZ accepted **58,767,739 B / 4.7873709325 s encode / 1.8721679705 s decode**, Zstd-3 **66,216,569 B / 1.2561 / 0.8551 s**, XZ-6 **49,233,340 B / 116.9892 / 2.3079 s**. Different in-memory/CLI timing conventions.

**Next:** continue with the dedicated [Zstd architecture coverage ledger](../research/zstd-architecture-coverage.md). This is ONE investigated Zstd family/placement, not an assertion that all Zstd innovations have been tested or that source splitting cannot help mathematical-token streams.
