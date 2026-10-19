# Round 55: recover the Round54 BWT gain without sacrificing compression (ACCEPT)

**Date:** 2026-10-19. **Verdict: ACCEPT H1. Round55 supersedes Round53.** All results use the complete original 12-file Silesia corpus (211,938,580 B). Original source SHA-256 checks: 12/12 PASS. No external compressor backend, copied Zstd/DEFLATE code, new on-wire subtype, user option, or extra CLI verb.

## Diagnosis and implementations

The user correctly rejected Round54's +1,028,079 B regression in exchange for encoding speed. We revisited its disposable dual-anchor source alongside approved Round53.

**Supplementary exact matcher screen (rejected):** independently encode both the original mathematical stream and the dual-anchor exact-match stream, and choose the smaller fully packed candidate per region. A single COMPLETE 12-file screen produced **60,838,336 B**, exactly equal to R53, but **8.243683 s encoding / 2.894864 s decoding**, much slower. The fast parser produced no additional net archive savings in that integration, so no parser code was promoted.

**Root cause of osdb's Round54 improvement:** archive-block inspection showed R53 selecting math/cross tokens (internal modes 9/10) while R54 selected existing BWT blocks (mostly mode15). R53's BWT preference gate compared old legacy-mode-12 BWT samples against mathematical candidate samples, even though the full BWT encoder emitted the stronger bijective mode15 stream. Its 5% required sample advantage screened out osdb, despite two of three BWT samples individually winning and an actually smaller full BWT archive. The dual-anchor parser had weakened the math sample comparator, *accidentally allowing BWT to win*. It was not responsible for the final codec's smaller data.

**Implemented additive correction:** in engine/bwt.go, measure the actual encodeBWTBijective (mode15) representation in sampleBWTPreference, keep mode15 sampled payloads for reuse in full BWT, and lower the independent two-of-three sample win margin from 5% to 2%. The full archive still uses exact-size selection; this 2% gate is a heuristic supported by the complete-Silesia screen, not a proof of general optimality. A new focused test validates that sampled bytes match the actual encoder and roundtrip exactly. The preexisting mathematical matcher, including all XOR/ADD paths, is entirely unchanged. R54's dual-anchor implementation remains rejected and uncommitted.

## Complete original-Silesia screens, ONE full encode/decode each

| Candidate implementation | All 12 archive bytes | Encode s | Decode s |
|---|---:|---:|---:|
| R53 plus dual-anchor supplemental stream, exact smaller-of-both | 60,838,336 | 8.243683 | 2.894864 |
| Actual bijective samples with previous 5% BWT margin | 60,773,351 | 5.417538 | 2.226846 |
| **R55: actual bijective samples with 2% BWT margin** | **60,726,433** | **5.618694** | **2.223234** |

These are each screening-only, not ten-run medians. No experimental variant above uses an external compressor. The final R55 candidate has **seven smaller files, five unchanged and zero compressed-size regressions**.

## Ten complete original 12-file Silesia candidate-only runs

Each counted run compressed and decompressed the entire original corpus, including full JLZ1 framing/checksums; all 120 original source reconstructions were byte-exact. Same precompiled candidate code and input files. A command wrapper terminated after four complete runs; three subsequent two-run invocations produced the remaining six complete runs, none omitted or partially counted. All archive totals were identical.

| Run | Complete archive bytes | Encode s | Decode s |
|---:|---:|---:|---:|
| 1 | 60,726,433 | 5.383660279 | 2.168268598 |
| 2 | 60,726,433 | 5.451027622 | 2.266837103 |
| 3 | 60,726,433 | 5.412268302 | 2.179289101 |
| 4 | 60,726,433 | 5.338530011 | 2.190641575 |
| 5 | 60,726,433 | 5.664815010 | 2.169265210 |
| 6 | 60,726,433 | 5.881540006 | 2.343478197 |
| 7 | 60,726,433 | 5.604427994 | 2.231321486 |
| 8 | 60,726,433 | 5.442326611 | 2.191676350 |
| 9 | 60,726,433 | 5.549503029 | 2.205210956 |
| 10 | 60,726,433 | 5.761809052 | 2.331187386 |
| **Median** | **60,726,433** | **5.500265326** | **2.198443653** |

**Previously approved R53** historical ten-run medians: **60,838,336 B / 5.721664347 s encode / 2.346933746 s decode**, not rerun. **Actual archive change -111,903 B (-0.183935% of compressed bytes); nominal encode -3.86949%, decode -6.32698%.** The timing percentages compare independent sessions, not paired experimental groups; do not treat the precise speed change as conclusively causal.

## Required comparison

| Model | Compressed size on complete 12 files | Reduction of 211,938,580 source B | Compress s | Decompress s |
|---|---:|---:|---:|---:|
| Zstd-3, historical external CLI | 66,216,569 B | 68.76% | 1.2561 | 0.8551 |
| XZ-6, historical external CLI | 49,233,340 B | 76.77% | 116.9892 | 2.3079 |
| Current accepted JLZ R53, retained historical ten-run median | 60,838,336 B | 71.29436% | 5.721664 | 2.346934 |
| **Candidate JLZ R55, newly measured ten-run median** | **60,726,433 B** | **71.34715%** | **5.500265** | **2.198444** |

| Candidate archive-size delta vs R53 | Candidate encode-time delta vs R53 | Candidate decode-time delta vs R53 |
|---:|---:|---:|
| **-0.183935% (-111,903 B)** | **-3.86949% nominal** | **-6.32698% nominal** |

External historical CLI timing measurement conventions differ from in-memory JLZ timings.

## Per-file complete archive sizes

| Source | R53 retained B | R55 B | Delta B |
|---|---:|---:|---:|
| dickens | 3,213,620 | 3,193,294 | -20,326 |
| mozilla | 18,609,139 | 18,609,139 | 0 |
| mr | 2,831,212 | 2,825,153 | -6,059 |
| nci | 2,169,196 | 2,163,199 | -5,997 |
| ooffice | 3,033,712 | 3,033,712 | 0 |
| osdb | 3,472,004 | 3,425,086 | -46,918 |
| reymont | 1,469,844 | 1,456,118 | -13,726 |
| samba | 5,134,262 | 5,134,262 | 0 |
| sao | 6,261,713 | 6,261,713 | 0 |
| webster | 9,851,849 | 9,837,297 | -14,552 |
| xml | 524,798 | 520,473 | -4,325 |
| x-ray | 4,266,987 | 4,266,987 | 0 |
| **TOTAL** | **60,838,336** | **60,726,433** | **-111,903** |

## Validation and compatibility

- Full canonical original Silesia SHA-256 checks 12/12 PASS; ten full candidate-only corpus encode/decode passes, **120/120 exact**.
- All Go package tests PASS via go test -count=1 ./..., go vet ./... PASS, targeted go test -race -count=1 ./engine including BWT preference test PASS.
- Windows/amd64 CGO-disabled cross-build PASS; native Windows execution not asserted.
- The previously compiled R53 CLI decoded R55's osdb archive exactly; new R55 CLI decoded the historical R53 osdb archive exactly. **No new decoder format**.
- Only one production source file changed (engine/bwt.go) and one regression test added (engine/bwt_preference_test.go). No R54 dual-anchor code was restored or copied.
- All accepted JLZ and external compression reference timings reused without repetition.

## Decision

**ACCEPT H1 as Round55**, not as Round54's ratio-regressing fast dual-hash parser. Round54 revealed a *model-selection blind spot* that is now corrected directly. Measured candidate archives are smaller across all twelve original Silesia inputs combined, no individual file grows, and newly measured speed medians are directionally favorable. This satisfies the user's mandatory no-ratio-regression rule.

Future research should look for genuine structural improvements in candidate selection that retain or improve measured full-Silesia archive bytes. Do not repeat dual-hash hash-width/skip tuning or treat a merely faster but larger candidate as acceptable.
