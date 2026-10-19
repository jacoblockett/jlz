# Round 46 — Exact DNA-like templates with insertion/deletion and reverse references

**Date:** 2026-10-19. **Decision: REJECT H1 for the tested integrated-block representation; preserve template-edit theory as OPEN for materially different implementations.** Round 44 remains the approved native JLZ codec. No production source or decoder changes, no GitHub Actions, no branch.

## Hypothesis and implementation

Round45's complete-original-Silesia raw mirror-reference scanner identified nontrivial orientation-inverted bytes, especially in `mr`, but offered no demonstrated incremental compression. User approved making DNA-like mutations (template + insertion/deletion/substitution) into a concrete executable prototype **complementary to the actual winning JLZ streams**, and investigating mirrors there too.

Implemented a disposable native Go encoder **and exact decoder** with a 2-entry, 131,072-bucket eight-byte anchor index and a maximum 2,048-byte reference. It emits a byte-tagged stream of seven token kinds: raw literal run, exact forward template, one-byte substitution, one-byte insertion, one-byte deletion, byte-reversed copy, byte-reversed bitwise-complement copy. Each token includes an unsigned-varint target length and source start when applicable; mutation types additionally include a varint edit position and optional inserted/replacement literal byte. Templates are constrained entirely to **already reconstructed bytes**; the decoder verifies source positions, consumed lengths, output bounds, opcodes and exact target length. An initial deletion-length accounting bug was found by original-byte verification and fixed **before** counting any experimental results. This stream then goes through JLZ's existing standalone native `packStream`, charging its stream header and Huffman/raw/RLE selection costs. There is no third-party compression backend. Complete encode/decode restoration was checked for every relevant whole Silesia region and intermediate symbol stream, not merely expected by theory.

Three controls compared: current reference stream; the new exact-copy-only token representation; and the fully edit/mirror-enabled token representation. This distinguishes whether insertions/deletions/substitutions/reversed templates improve *the new scheme* versus merely comparing against all-literal output.

## Only original complete Silesia, selected-path component comparisons

Read **all twelve** original canonical Silesia files (**211,938,580 bytes**) and verified 12/12 SHA-256 identities against `docs/reports/2026-10-18-round23-silesia-sha256.txt`. On six files whose whole-file winner uses BWT (`dickens,mr,nci,reymont,webster,xml`), evaluated every **412** full original-source BWT region: constructed the exact native BWT + bijective-zero MTF symbol stream, then compared `packStream(symbols)` to `packStream(templateEditEncode(symbols))`; decoded every token stream byte-exactly. These are genuine BWT intermediate streams, but the production R44 encoder reuses 18 already encoded *legacy* BWT preflight blocks, and their precise archive payloads can differ from the all-bijective component baseline here.

On four files normally selecting mathematical regions (`mozilla,ooffice,osdb,samba`), evaluated **342** independent full input regions: compared `encodeMathRegionBlockFast(originalBytes)` with `packStream(templateEditEncode(originalBytes))`, with exact reconstruction. **Limitation:** this mathematical comparison deliberately uses each isolated region's fast mode; the production encoder sometimes chooses a cross-region-history mode, so these are *component* comparisons rather than a reconstructed production archive. `sao` and `x-ray` were read and included in corpus identity validation but use neither selected BWT nor independent math-region mode; no alternative block codec was tested there.

| Original file | Baseline packed component B | New packed component B | Difference B | Blocks where edit candidate won | Sum of isolated winning-block savings B |
|---|---:|---:|---:|---:|---:|
| dickens | 3,294,743 | 3,297,478 | +2,735 | 0 | 0 |
| mozilla* | 18,693,710 | 31,853,403 | +13,159,693 | 2 | 225 |
| mr | 2,856,176 | 2,856,399 | +223 | 0 | 0 |
| nci | 2,187,050 | 2,187,518 | +468 | 8 | 2,752 |
| ooffice* | 3,040,496 | 4,628,425 | +1,587,929 | 0 | 0 |
| osdb* | 3,566,085 | 7,242,587 | +3,676,502 | 0 | 0 |
| reymont | 1,504,496 | 1,507,301 | +2,805 | 1 | 20 |
| samba* | 5,233,121 | 9,344,846 | +4,111,725 | 0 | 0 |
| webster | 10,158,516 | 10,168,121 | +9,605 | 0 | 0 |
| xml | 532,377 | 534,264 | +1,887 | 0 | 0 |
| **Total selected regions** | **51,066,770** | **73,620,342** | **+22,553,572** | **11/754** | **2,997** |

\*Math comparator is independent `encodeMathRegionBlockFast`, not a per-block reconstruction of all selected cross-history decisions. Values are already compressed native stream *component* bytes, not total JLZ archive sizes. There is no timing-matched Round44 baseline rerun.

The prototype fully encoded **130,308,185 input stream bytes** as **109,832,728 raw template-token bytes** prior to native packing; this reduction before entropy coding is deceptive. After exact native packing, the edit+mirror candidate was **73,620,342 B** against **74,280,619 B** for the *same* experimental exact-copy-only token scheme (net -660,277 B across components). Thus the edit/mirror operations really helped this representation, but it still lost to the current chosen native streams by 22,553,572 B in aggregate. Across all selected streams: **432,205 exact template spans (19,898,891 source bytes covered)**, **46,730 substitution spans (3,035,468 B)**, **4,742 insertion spans (409,534 B)**, **4,222 deletion spans (353,551 B)**, **3,530 reversed spans (299,648 B)**, and **128 reversed/complemented spans (7,876 B)**. These reflect one greedy bounded search, not an information-theoretic maximum or a promise of these bytes being newly compressible beyond existing JLZ.

The single full-original-Silesia **component screening process** (which recomputes reference BWT transforms and independent math-block compressors to compare costs) took **26.53 s elapsed**. This is **not** end-to-end JLZ compression time and must not be substituted as such. However, computing the edit candidates adds scanning/serialization work and no plausible material whole-file benefit is evident: choosing the alternative only for its **11** isolated winning blocks would save at most **2,997 B** within these measured components before charging any codec-selection overhead/headers, and mathematical eligibility may be illusory due to cross-history selection. No new integrated JLZ codec was built because this cheap fully reversible *selected-stream* screen failed its economic viability test. Therefore **no ten-run full-codec promotion benchmark is appropriate**; do not pretend candidate median timings exist.

## Required compressor comparison (retained references, NO integrated candidate)

| Compressor | Full compressed size | Original-byte reduction | Encode | Decode |
|---|---:|---:|---:|---:|
| Zstd-3, historical external CLI | 66,216,569 B | 68.76% | 1.2561 s | 0.8551 s |
| XZ-6, historical external CLI | 49,233,340 B | 76.77% | 116.9892 s | 2.3079 s |
| **Current accepted JLZ Round44**, historical ten-run median | **61,418,211 B** | **71.02075%** | **5.945576 s** | **2.445270 s** |
| **Candidate JLZ Round46** | **N/A: not integrated** | **N/A** | **N/A** | **N/A** |

| Candidate compression-size delta vs Round44 | Candidate encode-time delta vs Round44 | Candidate decode-time delta vs Round44 |
|---|---|---|
| **N/A** | **N/A** | **N/A** |

No accepted JLZ, Zstd or XZ performance reference was rerun. Their reference timings are retained from earlier independent sessions; CLI and in-memory runs are not directly comparable.

## Validation and future work

- All **754** analyzed independent source or transformed-stream blocks serialized and restored **byte-exactly** with the actual experimental edit decoder. No claim that R46 fully encodes/decodes all 12 files as one integrated codec: `sao` and `x-ray` use existing R44 paths.
- All **12** original Silesia input hashes match the canonical manifest.
- Following experiment file removal, the original R44 production source SHA blobs for `engine/bwt.go`, `engine/cross.go`, `engine/match.go`, and `bwt/zero_bijective.go` were verified against GitHub, then `go test -count=1 ./...`, `go vet ./...`, targeted `go test -race` on existing BWT/match paths, and a Windows/amd64 cross-build all **passed**.
- The serialized edit prototype remains **research only**; no production decoder mode, goroutines, compressor heuristics or compatibility changes.

**REJECT H1** for this bounded greedy template-edit token representation on selected BWT/mathematical streams. Round45's **590,048 B** of raw-file mirrored spans in `mr` do not translate into a viable post-BWT advantage (`mr` selected-stream candidate grows 223 B and has **0** winning blocks). Do **not** chase nearby token minimum lengths, hash seeds or alternative small anchor values: meaningful architecture must produce substantially different cost evidence. **Keep DNA-inspired deeper indel-aware dictionary/grammar models and biology-derived sparse structural scaffold research OPEN**, but require an explicit incremental selected-stream advantage and charged CPU before proposing integration. Round44 remains the current accepted compressor. Next substantial work should focus on speed of necessary winning-region encoding (rather than another near-duplicate copy matcher), or a mathematically distinct high-rank theory from the active queue.
