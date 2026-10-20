# Byte-normalized arithmetic range coding: complete integrated candidate benchmark

**2026-10-20. Final status: REJECT concrete integrated H2 for automatic promotion.** The earlier **nine-stream, +36-byte component screen** was prematurely treated as final; that screening decision was corrected by implementing, validating and benchmarking the entire existing JLZ pipeline. The numerical size penalty on the component alone was **not** an economically sufficient reason to stop; the real complete-codec result matters.

## Motivation and architecture

- Accepted source snapshot: root Go 1.23.2 native repository commit `684b5e3df1e3cfbd9adfb85d4750aaee5a2b8681`, with identical checked-in `default.pgo` SHA-256 `3b35cda74832360a252b63c1805e025aa1f4ed7a4fe4124d9a55929a6e86289f`. Current GitHub `main` had subsequent documentation-only commits; no production Go edits.
- JLZ's math-literal coder spends substantial CPU in order-1 bitwise adaptive arithmetic. Previous distinct byte-symbol and conditional-rANS replacements were faster on nine original Silesia math-literal stream components but enlarged their coded outputs much more. This trial **preserves the exact existing JLZ `bitModel.probability` and `bitModel.update` state sequence**, while replacing only the 32-bit E3 bit-pending interval output with 32-bit byte-normalized range coding. Architecture inspired by the publicly available LZMA range-output construction, implemented independently in native Go, not linking or using liblzma or any established compressor as a backend.
- New experimental packed-stream codec ID **6** encodes JLZ's original order-1 binary model using a 32-bit byte-normalized arithmetic interval with carry caching and canonical five-byte finalization. `packLiteralStream` compares the fast range-coded candidate against the existing raw/RLE/Huffman native options; the decoder handles both old codec3 and new codec6. Mathematical exact/XOR/ADD matching, accepted BCJ→BWT path, file SHA envelopes, automatic selection and two CLI commands are unchanged.
- The original nine-stream experiment (one component diagnostic timing pass) gave `packLiteralStream` **586,168 B** versus new range **586,204 B (+36 B, +4 B each stream)** and timed coder work about **163.07→86.51 ms encode**, **161.42→116.39 ms decode**. These are **NOT full-corpus speedups**, and even the 36 bytes are not the completed archive delta. This important earlier methodological mistake was explicitly corrected by full integration.

## Complete original twelve-file Silesia benchmark

- Total original **211,938,580 bytes**, 12/12 SHA-256 against `docs/reports/2026-10-18-round23-silesia-sha256.txt` **PASS**. No partial corpus, benchmark substitutes or external codecs.
- Candidate output **60,245,285 bytes**, exact original Silesia reconstructed **12/12 per pass**, total **120/120** complete byte-exact reconstructions, with embedded SHA integrity check.
- 10 **completed** candidate-only runs: an original 10-run invocation was interrupted after four *completed* passes. Those four measurements were retained; two separately completed 3-pass invocations supplied the other six. One interrupted run is excluded. All three runs used the identical candidate source/PGO and 12-file corpus.
- Accepted current median and Zstd/XZ reference values are **historical only and were NOT rerun**. External tool wall-time convention is different.

| Complete original Silesia file | Historical accepted size B | Candidate size B | Delta B |
|---|---:|---:|---:|
| dickens | 3,155,856 | 3,155,856 | 0 |
| mozilla | 18,609,139 | 18,609,844 | +705 |
| mr | 2,729,261 | 2,729,261 | 0 |
| nci | 2,163,199 | 2,163,199 | 0 |
| ooffice | 2,692,252 | 2,692,268 | +16 |
| osdb | 3,425,086 | 3,425,086 | 0 |
| reymont | 1,455,234 | 1,455,234 | 0 |
| samba | 5,134,262 | 5,134,569 | +307 |
| sao | 6,261,713 | 6,261,713 | 0 |
| webster | 9,835,416 | 9,835,416 | 0 |
| xml | 515,852 | 515,852 | 0 |
| x-ray | 4,266,987 | 4,266,987 | 0 |
| **TOTAL** | **60,244,257** | **60,245,285** | **+1,028** |

| Completed full candidate run | Size B | Encode seconds | Decode seconds |
|---:|---:|---:|---:|
| 1 | 60,245,285 | 5.354068846 | 2.064659089 |
| 2 | 60,245,285 | 5.437781152 | 2.134644204 |
| 3 | 60,245,285 | 5.330814128 | 1.972602600 |
| 4 | 60,245,285 | 4.867557138 | 2.016721944 |
| 5 | 60,245,285 | 5.426957382 | 2.220981002 |
| 6 | 60,245,285 | 5.382774648 | 2.084463105 |
| 7 | 60,245,285 | 5.413442430 | 2.133393709 |
| 8 | 60,245,285 | 5.583681814 | 1.991180167 |
| 9 | 60,245,285 | 5.130865045 | 1.978868658 |
| 10 | 60,245,285 | 4.905288422 | 1.958629355 |
| **Median, each independent series** | **60,245,285** | **5.368421747** | **2.0406905165** |

**Retained accepted historical complete-corpus medians:** **60,244,257 B, 5.322679488 s encode, 2.003659634 s decode**. Candidate relative **+0.00170639% compressed bytes**, **+0.859384% nominal encoding wall**, **+1.848162% nominal decoding wall**. These are independent, unpaired historical and candidate sessions; sub-2% runtime shifts may be system noise, so do **not claim proof of universal slower runtime**, but no significant positive end-to-end speed improvement was demonstrated. The component's ~1.9x/1.4x speedup does **not** multiply into the complete pipeline's throughput.

Historical external Zstd-3 **66,216,569 B / 1.2561 / 0.8551 s** and XZ-6 **49,233,340 B / 116.9892 / 2.3079 s** are unchanged comparison references only.

## Correctness, security and backward decoding

- New range coder original fixture suite covered 0/1/2/30/256/4096/65536 byte payloads under constant, periodic, random and perturbed distributions, all roundtripped exactly; truncation rejection tested for every nontrivial payload.
- The two existing codec3-selection assertions were adjusted to require codec6 **only in the disposable candidate**, preserving their decoded output, corruption rejection and archive integrity checks; no production test was edited.
- **Uncached `go test -count=1 ./...` PASS**, **`go vet ./...` PASS**, targeted `go test -race -count=1 ./entropy -run '^TestRangeByteTrialExact$'` PASS. **CGO=0 Windows/amd64 build PASS** (cross-build only).
- An actual earlier accepted CLI-generated archive from the fixed `04_sparse_mutations.bin` original fixture was decompressed using the experimental codec6-capable decoder, with exact output comparison and SHA integrity: **backward decoding PASS**. Existing binaries are not expected to decode newly emitted codec6 (forward incompatibility).
- All experimental production code remained in the disposable directory `/mnt/data/jlz_range_full` and **was not committed**. No third-party compressor is linked or used for production.

## Decision and next architecture gate

**REJECT this exact integrated range-output candidate under the existing all-three-metric acceptance rule.** Not because **36 bytes** mattered, but because the actual end-to-end ten-run evidence shows **+1,028 B** and **no demonstrated overall speed benefit**. Encoding +0.86% and decoding +1.85% nominal are too noisy to overinterpret, but do not substantiate the substantial speedup observed on extracted streams. Even user approval for a tiny size trade would not automatically make this integrated implementation compelling without a real whole-codec speed win.

Do not erase the positive component insight or treat generic native byte-range coding as theoretically exhausted. A further attempt needs a *structurally different* way to make faster entropy paths matter end-to-end: e.g., eliminate expensive duplicate model work or alter whole-codec workload routing without a large byte regression. Do not loop over byte terminators, context counts or precision shifts. Main production stays the accepted native BCJ→BWT plus PGO engine.
