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

## Verdict and next experiment

**DO NOT PROMOTE.** The proposed early-exit gate is exact under its stated nonnegative-length model and potentially eliminates expensive suffix sorts, but real workload occurrence and end-to-end benefit remain unmeasured.

When a licensed full corpus and executable isolated baseline checkout are available: implement this guard in a disposable local copy, compare it against untouched Round 21 on Silesia plus heterogeneous legally reusable data, record per-file eligibility/early-reject counts, conduct ten paired full-corpus encode runs and relevant decoding/compatibility/test/vet/race checks, and only then consider promotion onto `main`. Continue avoiding file-specific exceptions and unnecessary suffix-enumerator replacement.
