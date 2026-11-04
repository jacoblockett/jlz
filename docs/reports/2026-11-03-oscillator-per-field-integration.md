# Oscillator-derived field-state supplement: measured integration (2026-11-03)

**Decision: BENCHMARKED CANDIDATE, NOT PROMOTED.** The accepted production encoder/decoder and root `default.pgo` on `main` are unchanged. This trial uses a new mode 21 in disposable native-Go source only. [Earlier discovery/novelty and phase screen](2026-11-03-oscillator-field-innovation-screen.md).

## Mechanism and novelty

The user-suggested oscillator/atomic-clock analogy led to periodic **innovations**: decode a residual against the same field one 16-byte record earlier. A fixed harmonic second difference and costed phase relocking were separately screened in the initial report and did not win. The existing JLZ mode13 already separates field lanes. The genuinely new **placement** here is a first-order periodic innovation **inside the four actually selected Mozilla mode13 winners**, with a per-field choice of identity, modular subtraction, or XOR. This is distinct from the previously rejected fixed lag28 XOR/SUB before *sao* field lanes, fixed lag8 before *x-ray* lanes, and raw-source periodic record detection.

For each byte `i >= 16`, the encoder chooses, for a full byte-position field, one of `r[i]=x[i]`, `r[i]=x[i]-x[i-16] (mod 256)`, or `r[i]=x[i] XOR x[i-16]`. The decoder restores original bytes sequentially using previously restored `x[i-16]`. Initial 16 bytes stay unchanged. Under internal block mode 21 the payload holds predictor kind, fixed lag16, an existing native lane frame (width 16 or 64), and, for mixed selectors, a two-bit-per-field map at the end of the frame. Existing block independence, four-region history groups, complete checksum, and earlier decoding paths remain intact. Other methods remain eligible and the original encoded mode13 payload is preserved as the fallback.

Novelty gate: this screen is not another period-width threshold sweep. Fixed-period lanes and sao lag28 deltas had already been tested. The differentiated architecture is **independently choosing native reversible predictor roles per field and preserving the incumbent coder/framing**, followed by a work-elimination assessment. Old source code was not copied from any external codec.

## Actual component and integration steps

- Previous ideal-entropy screen: 211,593 B to 146,129 B across four selected Mozilla blocks, **only an optimistic zero-order entropy signal**, not compressed bytes.
- Stage A: actual existing native lane encoder after one global lag16 XOR/SUB predictor, keeping the prior selected lane width: best 31,689 B total whole-corpus archive saving (single full-corpus screen and 10 complete candidate-only passes, each exact).
- Stage B: allow existing native lane width 64 as a supplement to width16: total whole-corpus saving 44,714 B (single full-corpus screen and 10 complete candidate-only passes, each exact).
- Stage C: encode per-field **identity/SUB/XOR choices** inside width16 or width64, select shortest framed per-field native streams and serialize a 2-bit selector per field. This is the final benchmarked candidate. Its four Mozilla native payloads are approximately: block58 width16 33,456→14,712 B; block59 width16 46,862→35,169 B; block60 width16 100,736→87,290 B; block85 width64 53,496→39,273 B. These component amounts include the new transform tag and per-field map. **Final actual integrated complete Silesia improvement is authoritative: 58,107 B**, including exact whole-archive framing.

All 12 complete *original* Silesia sources (211,938,580 B) were read from the user-supplied archive and passed all twelve committed SHA-256 manifest identities. The source ZIP matched public main at `825826d790dbd0d87301586c5dc35d2f1fba1779`, which has only documentation changes after the accepted source baseline.

## Complete original Silesia benchmark: final candidate only (10 candidate-only runs)

| Run | Compressed B | Encode seconds | Decode seconds |
|---|---:|---:|---:|
| 1 | 58,709,632 | 5.378297760 | 1.998456504 |
| 2 | 58,709,632 | 5.311242177 | 1.986876099 |
| 3 | 58,709,632 | 5.629156388 | 1.955071044 |
| 4 | 58,709,632 | 5.306443127 | 1.977429638 |
| 5 | 58,709,632 | 5.165212048 | 1.948195687 |
| 6 | 58,709,632 | 5.111530742 | 1.893585860 |
| 7 | 58,709,632 | 5.206543564 | 1.923517199 |
| 8 | 58,709,632 | 5.435274219 | 1.990344052 |
| 9 | 58,709,632 | 5.411467295 | 1.948663200 |
| 10 | 58,709,632 | 5.102595752 | 1.891158079 |
| **Median** | **58,709,632** | **5.308842652** | **1.951867122** |

Original bytes: **211,938,580**. Candidate reduction **72.2987518365%**. Retained accepted complete 12-file JLZ: **58,767,739 B / 4.7873709325 s encode / 1.8721679705 s decode** (historical ten-run medians, **NOT rerun**). Final candidate: **-58,107 B, -0.0988757% archive bytes (better), +0.0274169 percentage points reduction (better); +10.89265% encoding (worse) and +4.25705% decoding (worse)** nominally versus historical unpaired timings. **Do NOT interpret these cross-session timing deltas as isolated algorithm slowdowns:** this runtime environment differs from the October acceptance session. External retained references Zstd-3 **66,216,569 B /1.2561 s/0.8551 s** and XZ-6 **49,233,340 B/116.9892 s/2.3079 s** use different CLI timing conventions.

Per-file archives: *mozilla* **18,452,630→18,394,523 B (-58,107 B)**; all eleven other originals stayed byte-size-identical to the accepted baseline. All **120/120 complete source reconstructions exact** and SHA-checked. No additional user options or codec backend. All output sizes are complete final archives, not selected component bytes.

## Local work measurement and validation

A targeted same-process repeated component comparison on the four Mozilla blocks, **20 repetitions of the four-block set**, measured original selected source-lane encoding **0.172190 s aggregate** and **an additional 0.633566 s** for creating/testing the residual candidates, including their separate native field frames. This represents **31.68 ms extra encoder work per four-block set**, before measuring any other JLZ path. The original lane decoder took **0.081370 s**, new mode21 decoder **0.136964 s** aggregate, about **2.78 ms extra decoder work per four-block set**. These are component-level timings, not full-corpus timing changes. The extra encodings allocate temporary field streams; peak heap was not profiled.

All new-mode and legacy tests PASS: `go test -count=1 ./...`, `go vet ./...`, targeted `-race` on oscillator/record tests, `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build`; native randomized and actual-corpus exact inverses, malformed predictor tags, invalid lane-width/class selectors, truncation and corruption. New decoder read a previous baseline JLZ archive exactly. Previous decoder rejected a new-mode Mozilla archive as expected; this is **new archive forward incompatibility with old binaries**, not broken backward decoding. Candidate must not be promoted silently.

## Decision, remaining hypothesis

**Do not promote this concrete mode21 integration:** actual size gain is genuine but small (0.099%) and the added per-field native packing demonstrably costs CPU. Ten candidate-only full benchmark timings did not prove three-metric nonregression against the retained acceptance reference, and adding a decoder mode for this isolated saving merits extra scrutiny. This is **not rejection of periodic innovation theory**. The next genuinely different refinement is **work elimination**, namely source-sample role selection that codes only one chosen per-field residual, potentially reusing current original-lane packet frames rather than recomputing three candidate stream sets. The current candidate has not been established as a worthwhile production speed/ratio trade and remains in a disposable local source patch, not `main`.

No code, PGO, binary, or temporary research harness was committed. Research documentation is committed by strategy/date, not as numbered rounds.
