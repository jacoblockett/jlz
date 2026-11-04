# Oscillator phase and periodic-field innovation screen (2026-11-03)

**Status: SCREENED, promising natural-source component; NO integrated JLZ candidate, NO whole-codec benchmark, and NO production source change.**

## Inspiration and prior-work collision check

Atomic-clock resonant interrogation and phase/frequency servo feedback (NIST, https://www.nist.gov/atomic-clocks/how-do-atomic-clocks-work; https://tf.nist.gov/ofm/smallclock/LO.html) suggest identifying a stable cyclic phase and encoding its *innovations*, not treating arbitrary bytes as physical pendulums. Discrete oscillator prediction has the form `x[t] ≈ 2 cos(ωΔt) x[t-1] − x[t-2]`. Also examined changing phase (phase slips), phase-lock tracking, and harmonics.

The nearest JLZ work is **already accepted** original-source 28-byte record-aligned mode13 field lanes in `sao` and the four accepted 16/48-byte `mozilla` lane regions, together with **rejected** 28-byte previous-record XOR/SUB before `sao` lanes, 8-byte fixed-lag XOR before 8-way lanes in other files, and older spatial/recurrence transforms. See `docs/research/layer-placement-matrix.md` sections on record28, native-lane lag8, and field-relative deltas. This screen's novel **placement and data** is a first-order modular/bitwise innovation *before the current selected 16/48-way Mozilla record lanes*, not another period search or a repetition of the rejected sao test.

## Exact proposed transform and placement

For a fixed phase-related byte lag `p`:
- Modular velocity/innovation: `d[i] = (x[i] - x[i-p]) mod 256` for `i>=p`; prefix bytes stay literal. Inverse `x[i] = (d[i] + x[i-p]) mod 256`.
- Also screened XOR innovation `x[i] XOR x[i-p]` and discrete acceleration `(x[i] - 2*x[i-p] + x[i-2p]) mod 256`, each with a causal exact inverse.
- **Supplement** the current `mode13` original-byte lane codec by computing innovations *first*, then splitting the innovation bytes across the **already chosen** lane width. Preserve the original incumbent byte-lane payload as fallback; select only using full native packed bytes after a cheap preflight. For the width48 Mozilla region, predictor lag16 is independently plausible, not necessarily lag48. New decoding would require an explicitly tagged mode/transform field, and older archives must remain readable.

All figures below are **ideal zero-order empirical entropy of independent byte lanes, expressed in bytes**, *not* actual JLZ `packStream` output, headers, codebooks, or benchmarked archive sizes.

## Natural-original-source component evidence

The supplied complete original Silesia ZIP totals **211,938,580 bytes**; all **12 SHA-256** values matched `docs/reports/2026-10-18-round23-silesia-sha256.txt`.

On four **currently selected** original Mozilla 256-KiB blocks:
- Block 58, width16/lag16: original lanes **32,794 B**; modular SUB **13,094 B** (optimistic **19,700 B** less); XOR 14,252 B; second difference 16,378 B.
- Block 59, width16/lag16: original **43,812 B**; SUB **28,879 B** (**14,933 B** less); XOR 29,194 B; second 32,911 B.
- Block 60, width16/lag16: original **87,608 B**; XOR **72,203 B** (**15,406 B** less); SUB 72,797 B; second 76,095 B.
- Block 85, width48 with predictor lag16: original **47,379 B**; XOR **31,953 B** (**15,425 B** less); SUB 32,111 B; second 34,893 B.

**Per-block best-of-XOR/SUB entropy bound: 211,593 B old lane entropy vs 146,129 B proposed, an optimistic 65,464-B difference, BEFORE actual entropy-code lengths, additional mode metadata, decode CPU or automatic selector costs. This is not a predicted JLZ archive saving.** `sao` block10, width28/lag28, runs the opposite way: original 174,187 B, SUB 182,074 B (+7,887 B worse), XOR 182,323 B. That agrees with the earlier rejected sao-relative-delta placement: selection is essential.

**Reversibility:** an independently compiled standalone original Go diagnostic round-tripped all tested modulo-SUB, XOR and second-difference transforms exactly across **36 actual source-block/lag combinations × 3 transforms** (four Mozilla blocks, two lags each; all 28 complete original sao blocks, lag28). This is not an integrated JLZ codec test.

## Phase relocking separate screen

A 4-KiB-window phase map with costed phase-switch penalty (128 bits/change) and two rounds of per-lane likelihood fitting found **zero phase switches and zero ideal entropy benefit** on sao block10 and the four actually selected Mozilla field-lane blocks at their current widths. Removing the penalty generated at most tens of optimistic bytes with many changes, a losing trade after phase metadata. A **synthetic** insertion of 63 bytes into a 256-KiB sao source block caused fixed phase28 empirical lane entropy to rise to 240,978 B; an oracle phase map lowered it to 174,248 B, demonstrating the mechanism but **not** proving a natural-source gain.

A separate causal three-lag (p−1,p,p+1) phase/frequency-tracking byte predictor did not improve the selected strongly periodic blocks; its lag selector stayed on the incumbent lag essentially 100% of the time. Simple second-difference residuals also lost to first differences after the existing lane separation. Neither phase tracking nor all conceivable oscillator models are exhausted.

## Decision and remaining validation

**PROMISING COMPONENT, NOT YET A BENCHMARKED CANDIDATE.** Do not promote or append it to `FINDINGS.md`'s candidate benchmark table. Current accepted JLZ remains **58,767,739 B, 4.7873709325 s encode, 1.8721679705 s decode** on the complete original 12-file Silesia corpus; these are historical retained ten-run medians, not rerun.

Next: with a runnable copy of current repo and its root `default.pgo`, add a *disposable* native-Go mode13-adjacent differential option only after the existing lane has been selected; gate it on cheap sample statistics, test actual existing `packStream` framing, exact inverse, backward decoding and memory/speed. Check whether candidate affects only four Mozilla blocks and whether actual packed savings survive. **Only if the integrated result is competitive** run 10 complete 12-file candidate-only trials and standard correctness/security tests. Record per-file regressions and speed. Do not commit rejected experimental source.

**Reproducible exploratory diagnostic:** a standalone original Go source file `jlz_oscillator_component_screen.go` was generated separately for the user to download. It reads only the Silesia ZIP, prints the five selected natural-block entropy measurements above, and checks 36×3 exact inverses. It is NOT a JLZ production dependency.
