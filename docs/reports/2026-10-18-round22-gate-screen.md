# Round 22: BWT sample-gate early-rejection screen (not promoted)

**Date:** 2026-10-18. **Approved source remains Round 21 on main.**
**Decision:** SCREENING ONLY, no JLZ encoder or decoder source changes. The complete Silesia corpus and a locally executable repository checkout were unavailable in this session; no ten-run whole-corpus or archive measurements were made.

## Hypothesis

The approved `sampleBWTPreference` encodes BWT for three sampled regions even when a partial result proves that the 2-of-3 wins plus 5% aggregate margin can no longer be satisfied. Skip unnecessary suffix sorts without changing any accepted BWT decision or selected archive bytes.

For a fixed three-region mathematical baseline with payload lengths `R[0..2]`, write `Rtotal = sum(R)`. After encoding BWT sample `i`, let `Bpartial` be the nonnegative sum of BWT lengths so far, `wins` the count of BWT lengths below their respective mathematical lengths, and `remaining = 2-i`.

An early rejection is logically sufficient if either condition holds:
- `wins + remaining < 2`. Even every remaining win cannot meet the 2-of-3 rule.
- `100*Bpartial > 95*Rtotal`. The eventual BWT total can only increase, so it cannot satisfy the required 5% total margin.

The original selected-archive semantics are preserved whenever BWT is selected, because no successful gate decision can meet either rejection condition. An unselected BWT candidate need not complete unused sample transformations.

## Isolated analytical screen

A standalone, disposable Go 1.23 prototype replicated the approved final gate predicate and the proposed early-rejection predicate. Exhaustive enumeration over the six positive integer payload sizes `R[0..2], B[0..2]`, each in `1..8`, found:

| Metric | Existing three-sample predicate | Proposed early rejection |
|---|---:|---:|
| Size combinations compared | 262,144 | 262,144 |
| Selection disagreements | Reference | **0** |
| Hypothetical BWT sample encodes | 786,432 | **683,328** |
| Encodes eliminated in this synthetic grid | 0 | **103,104 (13.1%)** |
| Grid cases with early rejection | 0 | **95,040 (36.3%)** |

Of the early-rejection cases, 8,064 stopped after one candidate and 86,976 after two. **These counts are synthetic enumerations of payload lengths, not workload frequencies or a measured JLZ speedup.**

To establish whether avoiding an encode could matter, a separate standalone probe reproduced the approved use of `suffixarray.New(data).Write` and its record parsing. Each locally installed real-world sample was warmed twice and timed ten times, using a fixed 256 KiB prefix where available:

| Locally installed sample | Bytes | Median of 10 independent suffix-index/parse runs |
|---|---:|---:|
| GNU bash executable | 262,144 | 12.590 ms |
| Python 3 executable | 262,144 | 12.489 ms |
| Go standard-library source concatenation | 262,144 | 12.208 ms |
| Boost Asio header | 8,948 | 0.305 ms |

Host-specific times exclude JLZ's MTF, zero-run, entropy packing, mathematical baseline, orchestration and archive framing. The Go 1.23.2 standard-library serialization passed the probe's size/header/suffix checks. No need to replace the suffix enumerator was demonstrated. No third-party compressor was used. The installed files were read locally, not copied into Git.

## Full archive references and limitations

Approved Round 21 remains **63,770,766 B** on 12 Silesia files, with ten paired whole-corpus medians **10.562019 s encoding / 4.160373 s decoding**, as measured previously against unchanged Round 16. Historical independent reference sizes remain **Zstd-3 66,216,569 B** and **XZ-6 49,233,340 B**. None of these benchmarks was rerun here.

**No complete JLZ archive sizes, decoding medians, Silesia per-file effects, or cross-version archive audits exist for this unpromoted screen.** Do not label this a performance improvement to JLZ, infer an `osdb` speedup, or change the approved baseline.

## Follow-up: actual 11/12 Silesia inputs supplied to project sources

The user supplied eleven uncompressed Silesia files as project sources. Their measured input sizes total **206,593,300 B**, matching the Round 21 manifest sizes for those eleven names. The twelfth file, `xml` (**5,345,280 B**), was **not available** in the provided project sources or local mount. The canonical full 12-file corpus remains **211,938,580 B**. This is **not** a full-corpus test.

Ran a standalone Go 1.23.2 local suffix-index/serialization probe on the approved gate's first/middle/last 256 KiB region positions for all eleven available files. Each position was warmed once, then timed independently for ten repetitions. All 33 positions passed original-data serialization header checks. Machine-dependent medians, in milliseconds:

| File | First | Middle | Last | Sum of separate medians |
|---|---:|---:|---:|---:|
| dickens | 13.154 | 12.675 | 11.254 | 37.083 |
| mozilla | 10.525 | 18.813 | 3.506 | 32.844 |
| mr | 8.955 | 8.460 | 0.238 | 17.653 |
| nci | 7.869 | 6.994 | 7.296 | 22.159 |
| ooffice | 13.689 | 14.886 | 3.495 | 32.070 |
| **osdb** | **11.694** | **11.627** | **5.749** | **29.070** |
| reymont | 12.591 | 11.945 | 2.367 | 26.903 |
| samba | 12.736 | 6.043 | 2.734 | 21.513 |
| sao | 13.114 | 12.973 | 8.602 | 34.689 |
| webster | 13.305 | 11.758 | 2.131 | 27.194 |
| x-ray | 12.232 | 12.180 | 3.767 | 28.179 |

Raw measurements: [Round 22 Silesia-11 suffix-sort CSV](2026-10-18-round22-silesia11-sort-probe.csv). Local disposable probe used `suffixarray.New(block).Write` with signed-length/header validation. This measures a **component of** BWT sample-gate CPU overhead, not JLZ's full `encodeBWTNative` path, sample selection, end-to-end encoder duration, archive bytes, or decode time. The `osdb` 29.070 ms sum is not directly comparable with the approved 13.9% whole-file slowdown measured in a different run. The actual gate exit distribution remains unknown without the Round 21 executable and complete BWT and math sample encoding decisions.

**Verdict unchanged: SCREENING ONLY.** No source changes or baseline promotion. `xml` must be supplied before a verified full 12-file benchmark can be claimed. Avoid committing research corpus bytes.

## Verdict and next experiment

**DO NOT PROMOTE.** The proposed early-exit gate is exact under its stated nonnegative-length model and potentially eliminates expensive suffix sorts, but real workload occurrence and end-to-end benefit remain unmeasured.

When a licensed full corpus and executable isolated baseline checkout are available: implement this guard in a disposable local copy, compare it against untouched Round 21 on Silesia plus heterogeneous legally reusable data, record per-file eligibility/early-reject counts, conduct ten paired full-corpus encode runs and relevant decoding/compatibility/test/vet/race checks, and only then consider promotion onto `main`. Continue avoiding file-specific exceptions and unnecessary suffix-enumerator replacement.
