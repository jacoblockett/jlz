# Zstd cross-history length work-elimination: five-variant batch and ten-run native qualification

**Date:** 2026-10-21. **Status:** five candidate selection/source layouts screened, **one additional ten-complete-original-Silesia candidate-only encode/decode qualification**, reproducible compression savings but no reliable same-host speed improvement. **No production promotion**. Accepted baseline, Go source and root `default.pgo` unchanged.

## Novelty/source gate

This is the second-stage investigation of the accepted [Zstd sequence-symbol encoding choices](https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_compress_sequences.c) and [Zstd superblock sequence store](https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_compress_superblock.c). Specifically audited `ZSTD_selectEncodingType`'s cheap histogram, repeat/default/dynamic table-cost selection and the requirement that individual literal-length/match-length sources remain independently reconstructible. JLZ **already accepted** separate literal/reference length coding in non-cross modes in historical Round16, and the [previous batch](2026-10-21-zstd-batched-superblock-and-cross-lengths.md) introduced an independently decodable **cross-history type20 mode** with exact framed incumbent-vs-split selection, +0.134849pp reduction. [Ten-run initial qualification](2026-10-21-zstd-cross-length-qualification.md) demonstrated the exact-size improvement, but current-host speed was inconclusive. This batch focuses on **removing redundant source/entropy passes**, not rediscovering separate-length encoding or the historically ineffective H0 Shannon shortcut. Same old match/XOR/ADD tokenizer, native codebooks, 128KiB group-local history and format type20.

The local accepted snapshot and both experimental Go sources retained the committed root `default.pgo` byte-identically. The source snapshot had the same production Go blobs as live docs-only main (per earlier report). The original unmodified twelve Silesia members **211,938,580 bytes** all passed the committed SHA-256 manifest in their corpus root (12/12). No external compressors/codecs, generated datasets or alternate corpora.

## Five concrete variants batched before full qualification

1. **CHOOSER control** (already tested): split original commands into native five columns; encode merged-length and two typed length streams; choose the actual smaller framed payload. Preserves `blockCrossPacked` fallback. On four selected mathematical files mozilla/samba/ooffice/osdb, compressed total **29,418,432 B**.
2. **ALWAYS-SPLIT**: if both typed streams are nonempty, encode only them and *never* encode the merged-length alternative; retain old mode for degenerate one-sided commands. Does not retokenize. **29,418,596 B** on four selected files, exactly **164 B** larger than chooser. Full original 12-file archive **58,482,105 B (72.406106996% reduction, Δ+0.134772065 pp better vs accepted JLZ)**. Only Mozilla **+51 B** and Samba **+113 B** relative to chooser; remaining ten file sizes exactly equal. No file-size regression vs accepted baseline.
3. **SMALL-STREAM FALLBACK:** ALWAYS-SPLIT except encode the merged alternative when combined varint length stream is shorter than 1,024 B. Four selected-file archive **29,418,596 B**, same as ALWAYS-SPLIT. The 164 B of actual small size losses therefore do **not** fall below this crude byte-count threshold. Its first matched screen was slower; **reject this specific uncalibrated cutoff**, no subsequent long sweep.
4. **FUSED TYPED PARSER:** keep original exact smaller framed stream decision but extend only the cross-history `splitMatchTokens` path to append per-type varints while reading source commands *once*. Other ordinary JLZ modes still call unmodified `splitMatchTokens` wrapper, with no typed work. Avoid the second `splitOperationLengths` command/varint walk; preserve all typed and merged stream bytes and all current mode selections. Four selected-file archive **29,418,432 B**, full original all12 **58,481,941 B**, *byte-identical per-file to the prior CHOOSER*.
5. **FUSED ALWAYS-SPLIT:** combine variants 2 and 4; four selected archive **29,418,596 B** and one complete original archive **58,482,105 B**, no additional size gain, unstable/slower one-pass encode. Do not promote the extra mutation.

For all cases, complete framed archive roundtrips from the original source were exact. Every source archive in all 10 FUSED qualification passes decoded and restored byte-for-byte; each pass independently compared the SHA-256 for every original file (120/120 restores). **No full ten-run median was executed for ALWAYS or FUSED ALWAYS**; timings for those are strictly screen/paired diagnostics.

## Proper ten-run complete original Silesia: FUSED TYPED

| Metric | Accepted JLZ original historical | FUSED ten complete runs |
|---|---:|---:|
| Compression reduction | 72.271334931% | **72.406184377%**, **+0.134849446 pp better** |
| Compressed bytes | 58,767,739 | **58,481,941**, **0.486318% smaller better** |
| Encode median | 4.7873709325 s | **10.451408820 s in current host** |
| Decode median | 1.8721679705 s | **3.914058119 s in current host** |
| Exact restored originals | previous accepted corpus | **120/120 exact; every archive same size** |

Encode run times, s: **10.291468456, 9.267999379, 9.655135344, 11.137736231, 11.637799020, 11.405196172, 12.275015549, 10.611349185, 9.567689435, 9.573274631**. Decode run times: **3.552341657, 3.743698611, 3.793335362, 4.287056599, 4.555413324, 4.377712543, 5.370121431, 4.034780876, 3.605533318, 3.503001994**. Every run encoded/decompressed all twelve full original files with full JLZ framing, integrity and source-byte validation. **Do not compute a causal candidate-versus-historical speed delta**: original accepted reference was recorded in another host session. Historical Zstd-3 (66,216,569 B /1.2561s /0.8551s) and XZ-6 (49,233,340 B /116.9892s /2.3079s) were not rerun; external CLI timing conventions differ.

### Paired current-host separate-process full-original diagnostics

| ABBA diagnostic, all12 complete+exact | Unchanged JLZ controls, s | New variant, s | Paired median change |
|---|---|---|---|
| FUSED encode | 9.593807664; 9.927585096 | 10.084871227; 10.198567432 | **+3.904% slower (worse)** |
| FUSED decode | 3.659927774; 3.706692442 | 3.736162638; 3.957938592 | **+4.445% slower (worse)** |
| ALWAYS-SPLIT encode | 11.333178389; 9.811130837 | 10.362649814; 10.446771896 | **-1.584% faster nominally (better), uncontrolled drift** |
| ALWAYS-SPLIT decode | 4.791575807; 3.428482785 | 3.767844914; 3.774510773 | **-8.245% faster nominally (better), uncontrolled drift** |

Every diagnostic is a complete, byte-exact all12 original encode/decode; the controls are real unchanged native JLZ. The **ALWAYS controls alone drifted 1.52s encode and 1.36s decode**, much greater than apparent signal. In contrast FUSED's full ABBA was unfavorable despite four selected-file screens with apparently faster FUSED and ALWAYS medians. No current session speed gain deserves promotion; failed/stalled earlier run attempts not counted.

Four-file selected-run screens (diagnostic only): an ABBA **CHOOSER 6.299495/6.131285s** vs **ALWAYS 5.817923/5.726988s**, and independent FUSED **5.697565/5.878614s**. Current selected-source results did not predict the full twelve-file CPU direction, so do not use them as acceptance.

## CPU attribution and practical bounds

Collected a real `go test -cpuprofile` during complete original mozilla/samba/ooffice/osdb CHOOSER encoding **and decoding**, 89,064,756 source bytes, 30.48 CPU seconds of samples across worker goroutines. Entire native `engine.packFiniteState` inclusive accounted for **0.69 CPU seconds (2.26% of combined encode/decode samples)**. `engine.packCrossRegionTokensPreferred` inclusive was 5.84s (19.16%), but includes all expensive native literals/other streams; **do not confuse it with duplicate merged-length entropy work**. Source profiling also attributes ~8.14s (26.71%) to `entropy.EncodeArithmetic` and 6.69s (21.95%) to `DecodeArithmetic`; those are major bottlenecks, but previously screened native arithmetic model replacements and raw/Huffman alternatives are not novel without different Zstd-motivated architecture. Removing **one** merged-length pack cannot explain a robust many-percent whole-codec improvement; the finite-state total is an empirical loose upper bound for this four-file workload, *not* a global proof.

## Correctness/compatibility/decision

**FUSED** and **ALWAYS** each passed `go test -count=1 ./...`, `go vet ./...`, targeted `go test -race -count=1 ./engine` for existing archives, corrupt frames, concurrent decode (and FUSED's new split-frame fixture), and `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build`. The new FUSED targeted fixture writes >2,000 literal/reference command pairs ensuring type20 selection, checks exact inverse, rejects truncated stream and verifies legacy unsplit decoding. Other small variants passed targeted matcher/cross regression tests but did not warrant a full multi-platform promotion gate. Old JLZ archives remain readable in new proposed decoder; old production decoders cannot read *newly encoded* type20 frames. Both new test sources and all research harnesses are disposable, **not committed as production source**.

**Conclusion:** Cross-history type20 compression gain remains a real potentially valuable **candidate**, but simplifying the selector with ALWAYS loses 164 B and lacks trustworthy end-to-end speed evidence; FUSED preserves bytes and passed a full ten-run candidate qualification but its controlled 4-pass result is **3.9% slower encoding /4.4% slower decoding**, against which added source complexity is unjustified. **No production promotion.** Reject crude size threshold, defer simple always-split, and reject FUSED as demonstrated performance upgrade (not the broader architecture). Stop repeatedly chasing this low-CPU-share entropy selector with nearby thresholds, wrappers or stream-copy loops. The next high-impact Zstd-only architecture must address dominant work, such as amortized native literal sequence coding without loss of context, or shared parser/entropy cost decisions, with a new explicit source-backed novelty test.

**Local retained prototypes:** `/mnt/data/jlz_optbatch_20261021/{always,guarded,fused,fusedalways}`; old chooser `/mnt/data/jlz_qual_20261021/candidate`. Paths are not portable guarantees. No root PGO/production Go change and no workflow/branch was made.
