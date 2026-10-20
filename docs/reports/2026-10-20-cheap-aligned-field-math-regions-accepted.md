# ACCEPT: inexpensive aligned-field probe in selected mathematical regions

**2026-10-20. Current approved integrated JLZ on main; no active candidate.** Native Go source commit **`e632bc02872b983c050785d88297425457bc3e6f`**. Combined full original twelve-file Silesia, 211,938,580 bytes source; archive **58,767,739 bytes**; ten candidate-only medians **4.7873709325 s encode**, **1.8721679705 s decode**. Root `default.pgo` remains unchanged, Git blob `382b3c59a75bfa11bc60968f71cdde67dbe8d634`.

## Source-first hypothesis and prior different placement

Prior research measured genuine structured ELF-like regions in `mozilla`, whose full original 51,220,480-byte tar archive used native math-region modes with **4 out of 196 blocks** benefiting from the already accepted lane decoder: original blocks 58,59,60 at stride16 and 85 at stride48, for **156,509 B** total candidate savings. The first two integrated earlier experiments added per-block all-period 12..64 correlations and multiple full lane encodings, which showed:
- **H1 early probing:** 58,767,739 B but encode median **5.2219454765 s**, versus historical 4.8896856215 s, a 6.795% slowdown. A separate initial early-sampler variant wrongly short-circuited superior whole-file `sao` route; that failure was *only* experimental. Corrected early routing was still slow.
- **H2 deferred region probing:** 58,767,739 B but encode median **5.1387043605 s**, versus historical 4.8896856215 s (+5.093%), and decode median 1.8934515405 s (slightly worse than 1.876473360). **REJECTED** on real complete runs. [Detailed earlier experiments](2026-10-20-mozilla-regional-field-lane-screen.md).

This accepted H3 is **a low-work algorithmic placement**, not nearby threshold optimization of the discarded exhaustive search. Source analysis of original mozilla shows profitable 16-byte records have 16-lag equalities of approximately **0.91, 0.86, 0.72 and 0.84** respectively, whereas 8-lag equality is substantially lower. Previous period searches ran 53 lag comparisons on three 4KiB samples per candidate, and encoded `period,2period,3period` whenever a coarse period was detected. H3 substitutes a **single O(3×2048) source equality test** comparing lags **16 versus 8**, after the existing native math block has actually selected **blockNativeBytes or blockPredictive** (not on regions already won by costly matching). A large same-field signal and separation from 8-lag are required. Once selected, **encode native stride16 exactly once** using JLZ's existing `encodeLaneBlock`, compare true fully framed encoded bytes against current result, and **only if lane16 wins** compare one additional existing stride48 representation. No future mode identifier or extra entropy decoder is needed. Original mathematical matching and source stream bytes are preserved for other blocks, and the regional sampled early selector is unmodified.

Architectural inspirations from established open-source compressors: specialized literal/sequence streams in Zstd, independent field/byte classes in LZFSE, and separation of transform domain from the match parser. None was used as an encoder or decoder backend, and no third-party code was copied.

New permanent code confined to **`engine/periodic_math.go`**, one call in **`engine/regions.go`**, and **`engine/periodic_math_test.go`**. The source selection mechanism adds no new stream format and reuses original `blockLanePacked=13`. It does not widen source match history, add data copies, or change previous ordinary source modes. New field recognition may marginally miss unusual lengths/sources due to the conservative gate; this is a CPU/coverage choice with actual-size verification for any positive selection, never a decompression correctness tradeoff.

## Ten completed full original Silesia candidate-only runs

Each completed independent test process (compiled with Go1.23.2 and retained root `default.pgo`) compressed all **twelve complete original files**, decompressed all twelve, compared bytes and SHA-256 to the original, and charged every stream, archive header, codebook, and checksum. Complete original source total verified by the test as exactly 211,938,580 B. **120/120 full reconstructions PASS**. One preliminary full-source screen before the ten-run series gave **58,767,739 B / 4.770185348 s compression /1.864519000 s decompression**, also exact; this single screen is NOT included in ten medians.

| Candidate full run | Compressed B | Encode seconds | Decode seconds |
|---:|---:|---:|---:|
| 1 | 58,767,739 | 4.794678022 | 1.911817189 |
| 2 | 58,767,739 | 4.955685230 | 1.931324232 |
| 3 | 58,767,739 | 5.055251585 | 1.874293827 |
| 4 | 58,767,739 | 4.680893093 | 1.853276046 |
| 5 | 58,767,739 | 4.727622315 | 1.861448512 |
| 6 | 58,767,739 | 4.749102649 | 1.929576572 |
| 7 | 58,767,739 | 4.860344084 | 1.840169243 |
| 8 | 58,767,739 | 4.747851359 | 1.870042114 |
| 9 | 58,767,739 | 4.874773914 | 1.921634523 |
| 10 | 58,767,739 | 4.780063843 | 1.802065608 |
| **Median** | **58,767,739** | **4.7873709325** | **1.8721679705** |

**Prior accepted JLZ reference (historical, not rerun as timed benchmark):** `58,924,248 B / 4.8896856215 s encode / 1.876473360 s decode` with source-aligned record28+Huffman. **New candidate vs old:** archive **-156,509 B (-0.2656105%, better)**; compression median **-2.0924595% (better)**; decompression median **-0.2294405% (better, functionally neutral/noisy)**. Original reduction **72.27133493%**, versus prior **72.19748854%**. Exact file size changes are confined to **`mozilla` 18,609,139→18,452,630 B**; **other eleven files identical in compressed lengths**, including accepted `sao` 4,941,704 B.

Retained established references, NEVER rerun: Zstd-3 **66,216,569 B / 1.2561 s /0.8551 s**, XZ-6 **49,233,340 B /116.9892 s /2.3079 s**. Different external-CLI wall-time convention than in-memory JLZ.

**Interpreting unpaired timings:** The new encode speed improvement is moderate (~2.1%); decompression delta is **small (~0.23%) and not statistically established** across historical independent sessions. Production decoder and math/BWT/lane stream formats are byte-identical to the accepted baseline, and only four original blocks switch from old native/predictive to already supported lane16/48 mode. This is a net proven space gain at slight additional detect/encode work, but the ten-run timing series did not show a regression. Avoid promising speedup on every machine.

## Validation and backward compatibility

- Full original 12-file Silesia byte-equality and SHA-256 **120/120 PASS**, and identical full-archive total every candidate pass. Confirmed `sao` retains prior accepted record28 size.
- New targeted deterministic fixture tests: short source does not trigger, random input does not trigger, strong 16-byte synthetic source **does trigger and roundtrips** via existing `encodeLaneBlock`/`decodeLaneBlock`, and random source remains unchanged.
- `go test -count=1 ./...` **PASS**, `go vet ./...` **PASS**, `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ./...` **PASS** (not actual Windows runtime), targeted `go test -race -count=1 ./engine -run 'TestStronglyAligned16|TestRecordAlignedLaneRoundTrip|TestRecordPeriodDiscovery|TestTwoHopInverse'` **PASS**.
- A previously generated accepted **full mr JLZ archive** `jlz-previous-mr.jlz` decompressed using the exact new native Go CLI; source bytes and SHA-256 `68637ed52e3e4860174ed2dc0840ac77d5f1a60abbcb13770d5754e3774d53e6` match complete original `mr`, backwards decoder **PASS**. No new decoder formats or token grammar were introduced.
- All three new production/test Go blobs were created on GitHub and **confirmed to match the locally tested Git hash bytes exactly**: `engine/periodic_math.go` `ffb4e43aa3274c86d120e28af869631038d71206`, `engine/regions.go` `6af782d5e8d7c3c5b467b49d81300605d9756341`, and `engine/periodic_math_test.go` `d1ee1f8ac8ff43c8241970fab528a3a2f6f8bc29`. Atomic GitHub source commit `e632bc02872b983c050785d88297425457bc3e6f`; existing PGO unchanged.

## Decision and next distinct work

**ACCEPT H3**, now the **current** native JLZ; despite only 0.266% aggregate archive gain, it reclaims previously demonstrated 156KiB actual redundancy **without the discarded ~5–7% encoder slowdown**, and requires only small additive code. This accepted result updates the current benchmark only. Avoid nearby period parameter tuning.

Next priority should be a materially more powerful source-informed transform or match encoder on the **dominant mathematical region streams of Mozilla/Samba** and BWT-selected webster, ideally cutting multiple megabytes from the complete corpus without increasing work. In particular, investigate avoiding redundant encode passes in mathematical modes, original ELF segment/record relations, and reverse-cost feedback from actual selected tokens. The abstract mathematical structure queue remains open; one rejected placement is not a rejected theorem.
