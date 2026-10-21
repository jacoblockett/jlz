# Deferred E3 normalization on currently selected mathematical literal streams

**Date:** 2026-10-20. **Disposition:** REJECT tested byte-batched E3 mode; no production source, root PGO, archive subtype, or accepted benchmark changed. This is a genuinely different **renormalization scheduling** strategy, NOT a new entropy predictor.

## Prior-work/novelty verification

Before code, reviewed `docs/research/strategy-ledger.md`, `docs/research/layer-placement-matrix.md`, `docs/research/theory-backlog.md`, and the measured full selected-stream cost report [2026-10-20-selected-stream-cost-gap-and-native-fse-screen.md](2026-10-20-selected-stream-cost-gap-and-native-fse-screen.md). The nearest attempts include **byte-normalized arithmetic/range coding with the exact existing binary probability model** (historically one integrated ten-run trial rejected), exact reciprocal approximation for current E3 arithmetic, adaptive EMA probability replacement, static conditional rANS, adaptive bytewise Fenwick, full serialized bit-model ICM, and a FSE-only supplementary source model. None is a test of **retaining the existing E3 bit output format while deferring its E1/E2/E3 normalization loop across bits in a byte**.

The accepted `entropy/arithmetic.go` coder computes an exact current-model probability then normalizes its E3 arithmetic interval **after each of eight bit decisions per literal byte**. The experimental coder retains the **exact same order-0/order-1 bit-count tables, normalized p0 divisions, byte/prefix context, per-bit model updates, and 32-bit E3 interval**. It skips normalizing between adjacent bits of a byte when the interval still has sufficient precision, normalizes once after all eight bits, and performs an early emergency E3 normalization whenever the current interval width falls below **2^18**. The latter guards rounding loss: after normalization, even a 1/65536 rare event leaves a positive interval for the next decision. Unlike the previously rejected byte-output range coder, the experimental entropy output still uses bitwise E3 renormalization and the original bitWriter/bitReader; only the **schedule** changes. Delaying normalization changes E3 arithmetic rounding and hence the bitstream even though the binary model is unchanged, so an explicit independent new codec ID is required.

Architecture reference: existing actual JLZ code, together with the E3 normalization formulation in the repository's prior binary-range studies. **No established compressor backend, copied codec implementation, or external build dependency used.**

## Exact original data and verification

Local root repository snapshot from `jlz-6a388ba91682.zip`, with accepted production `entropy/arithmetic.go` Git blob **08549011d624f06b67d93f863945e119d804e955** and `engine/stream.go` blob **572681332382bc1323224f4ec34ebb0c46884735**, both verified **identical to live `main` at research start/end**. Entire original 12-file Silesia source, **211,938,580 B**, from user-supplied `silesia.zip`, with source manifest already verified 12/12; root `default.pgo` SHA-256 **3b35cda74832360a252b63c1805e025aa1f4ed7a4fe4124d9a55929a6e86289f**, unchanged.

Original independent native-Go `entropy.EncodeArithmeticBatch` and `DecodeArithmeticBatch` built and tested in a **disposable copy** of accepted JLZ. The component tests included seeded random bytes, repetitive, patterned, low-alphabet, empty, 1- and 2-byte, and up to 32768-byte fixtures, then *all 262 arithmetic-selected mathematical literal streams* from current JLZ's **actual complete original-Silesia compressed archives**. Every stream reconstructed **byte-exactly**. The original selected literal payload itself was re-encoded with accepted `EncodeArithmetic` and confirmed bit-for-bit equal to its stored codec3 payload; the alternate model was encoded/decode verified independently. There is no unverified synthetic replacement for the actual selected-stream screen.

## Selected real literal-stream component comparison

The sample spans **all 262 actually selected adaptive arithmetic streams** in Mozilla (179), ooffice (4), and samba (79), not arbitrary slices of original input or different match-token parsing. Compare native stream **bodies** only; framing is equal when both models use the same length and codec-byte format. The candidate codec would use a newly identified subtype 6, whereas incumbent uses subtype 3.

| Actual selected literal group | Existing E3 body B | Batched E3 body B | Difference | Winning streams |
| --- | ---: | ---: | ---: | ---: |
| Mozilla, 179 streams | 9,735,087 | 9,735,088 | +1 B (worse) | 3 |
| Ooffice, 4 streams | 218,332 | 218,332 | 0 | 0 |
| Samba, 79 streams | 1,723,010 | 1,723,008 | -2 B (better) | 2 |
| **Total 262** | **11,676,429** | **11,676,428** | **-1 B (better)** | **5** |

Measured in one combined component pass, sequentially encoding and decoding every exact stream through both implementations: **accepted E3 encode 3.782106 s, deferred E3 encode 3.027411 s (-19.96%, better)**; **accepted E3 decode 3.702171 s, deferred E3 decode 3.377503 s (-8.77%, better)**. This is a stage-level, order-dependent **single diagnostic pass** and not evidence of whole-codec speed gains. Both perform the same expensive eight adaptive probability calculations per byte; the gain is from fewer checks of E3 normalization conditions, while byte-frequency/rANS alternatives previously lost the order-1 information.

## Full JLZ integration and ten completed candidate passes

Supplemented the existing **native mathematical literal** selection: encoded the alternate using `EncodeArithmeticBatch` in place of the existing codec3 candidate, compared the same full framed payload against `packStream` as before, emitted new stream **codec6** only where it won, left all other current byte/RLE/Huffman/FSE/BCJ/lane/math BWT modes untouched; added a bounded decoder case6 retaining old codec3. This replaces, rather than duplicates, the existing expensive literal arithmetic encoding pass.

Ten **fully completed candidate-only encode AND decode runs**, split into 3+4+3 short self-contained test processes to avoid tool timeout. **All 12 complete original source files restored byte-for-byte and source-hash-exact on all 10 runs = 120/120.** All ten emitted the **deterministic, full-frame 58,767,738 B archive**, exactly **1 B** less than accepted JLZ. Per-file only changes: Mozilla **18,452,630→18,452,631 B** (1 B worse) and Samba **5,134,262→5,134,260 B** (2 B better). All others unchanged. Candidate compression ratio on 211,938,580 B improves only **0.000000472 percentage points**.

| Model | Complete original Silesia archive B | Median encode s | Median decode s | Size versus accepted | Encode versus accepted | Decode versus accepted |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Accepted JLZ, historical ten-run median | **58,767,739** | **4.7873709325** | **1.8721679705** | anchor | anchor | anchor |
| Native batched E3 codec6, **new ten full runs** | **58,767,738** | **5.120182884** | **1.856871727** | -1 B (0.00000170% better) | **+6.9519% (worse nominal)** | **-0.8170% (better nominal, unpaired/inconclusive)** |

**Actual ten complete candidate (encode, decode) elapsed seconds:**
1. 5.244602316, 1.935269426
2. 5.113328583, 1.843626015
3. 5.133207726, 1.825750391
4. 5.267737508, 1.902104409
5. 5.119101411, 1.866923934
6. 5.121264357, 1.862443448
7. 5.042356654, 1.819622070
8. 5.347606305, 1.898292017
9. 4.973746374, 1.851300006
10. 4.916852909, 1.817704840

This **unpaired historical** comparison is susceptible to machine/session load and does not prove the change causes slower encoding or faster decoding. Nonetheless, **it does not establish a sufficient practical improvement** to justify a new permanently maintained format and inverse. No fixed numerical acceptance threshold was assumed; implementation/compatibility costs clearly exceed the proven archive value. A proposed new subtype cannot be decoded by older JLZ binaries even though new binaries retain correct backward decoding. The apparent ~0.8% full decode win is especially weak given unpaired measurements and much faster isolated stage gains.

Historical unchanged external references, not rerun: **Zstd-3 66,216,569 B / 1.2561s encode / 0.8551s decode**; **XZ-6 49,233,340 B / 116.9892s encode / 2.3079s decode**. These external CLI timing methods differ from JLZ's in-memory timings.

## Additional compatibility, integrity and build verification

- Native synthetic and 262 existing selected source-stream bit-perfect decode tests PASS.
- Ten complete candidate full Silesia exact restores PASS, 120/120. New accepted candidate original archives all source-exact.
- A separate original accepted production encoder generated complete **Mozilla 18,452,630 B and sao 4,941,704 B** archives. Candidate decoder decoded **both actual old archives exactly (2/2)**. The accepted (older) decoder was separately tested to **reject a genuine new candidate Mozilla archive** carrying codec6.
- A genuine candidate archive was truncated at sizes 0, 1, 7, 17, halfway, and last-byte; and had bytes flipped at header, early body, middle and near tail. **All tested malformed archives rejected**. Raw arithmetic payloads historically allow implicit zero padding; archive length/bounds/checksum remain the final guard, so no claim of exhaustive arithmetic-only truncation handling.
- Full candidate `go test -count=1 ./...` **PASS** after updating two existing test assertions expecting codec3 to expect the new candidate's codec6; no original assertions removed. `go vet ./...` **PASS**. Targeted `go test -race -count=1 ./entropy ./engine -run '^TestBatchDeferredRoundtrip$|^TestBatchBackwardAndCorruption$'` **PASS**. `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ./` **PASS**. Native Windows runtime and exhaustive fuzzing not performed.
- Root `default.pgo` preserved. Candidate source and temporary test/Windows executable remain in disposable local test paths only; they are **NOT** promoted or committed.

**Disposition:** **REJECT specific deferred-byte-renormalized E3 new subtype** after a complete ten-run full original-Silesia integrated experiment. The hypothesized fast multi-symbol arithmetic decoding is **not** globally impossible: this measured stage speed improvement is real in a sampled context, but by itself does not constitute a valuable JLZ-wide improvement. Do not iterate byte group size, emergency interval threshold, or E3 normalization micro-timing without a materially distinct architecture. The production single-binary, native-Go exact lossless codec and all existing decoding remain unchanged. Future work must first identify a structurally different predictor representation or parallelism approach that preserves the ~1.24 MB of measured context gain over ideal H0 without adding serial per-bit overhead.
