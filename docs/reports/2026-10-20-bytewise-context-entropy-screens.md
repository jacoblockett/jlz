# Layered fast bytewise entropy trial: three reversible placements rejected

**2026-10-20. Outcome: REJECT the specific tested implementations; approved JLZ unchanged, no integrated candidate.** Historical accepted complete original 12-file Silesia: 211,938,580 source B, 60,244,257 archive B, accepted ten-run 5.322679488 s compression and 2.003659634 s decompression with committed `default.pgo`. Those references were not rebenchmarked.

## Architectural motivation and closest earlier experiments

Fresh CPU diagnosis found native binary order-1 arithmetic decoding accounted for 47.56% *cumulative sampled* decompression CPU. The current `entropy.DecodeArithmetic` performs 8 sequential prefix-bit arithmetic decisions, probability estimates and two model updates per literal byte. Studied external source-level approaches: Zstd's `lib/decompress/zstd_decompress_block.c` table-driven FSE/Huffman, FiniteStateEntropy `lib/entropy_common.c`, and libsais induced sorting as a distinct future BWT topic. No external code/backend used.

Existing rejected native context rANS was chiefly on BWT symbols and math opcode/length streams, while a prior new binary EMA estimator was fully integrated into math literals but failed the speed gate despite decreasing size. This trial targets a **new encoded domain/representation pairing**: actual V2 mathematical *literal bytes*, framed by the current native packed-stream API. It does not repeat the binary EMA predictor or replace matching/BWT/BCJ. Its failures do not refute all context entropy methods.

## Actual reversible native prototypes and measured components

For each full original Silesia source file `mozilla`, `samba` and `ooffice`, selected first, middle and final natural 256 KiB source blocks; generated JLZ V2 math token streams using `tokenizeMatchV2` and extracted literal columns using `splitMatchTokens`. These are **nine real original-source math-token diagnostics**, not claims that those particular token streams were selected in the accepted final archives, especially `ooffice` which selects BCJ→BWT extensively. Compared each representation to actual `packLiteralStream` including all candidate's mode, payload length and model-table bytes. Independent inverse decoded all nine streams exactly.

### H1A: static context-conditioned rANS

A single native 32-bit rANS state uses eight 4096-normalized tables selected by the top three bits of the **previous decoded literal byte**. The encoder knows preceding bytes when encoding in reverse, and the decoder reconstructs them in forward order; all eight sparse per-context model tables and the final state are transmitted. Reuses existing JLZ native finite-state model normalization and table layout, not another codec.

Nine sample current framed literal payloads **586,168 B**; alternative **635,203 B (+49,035)**. One same-process component timing diagnostic: current encode/decode **184.12/177.00 ms** versus new **11.95/7.77 ms**. **9/9 exact, 0/9 compressed-byte wins. REJECT placement** without an integrated benchmark or context-size tuning.

### H1B: adaptive byte-symbol arithmetic with exact previous-byte contexts

A native Go 32-bit interval coder encodes an entire next byte per decision instead of eight successive binary-prefix decisions. Each of 256 previous-byte contexts owns 256 Laplace-initialized symbol counts and a 257-entry Fenwick tree. The decoder determines its target symbol via 8-level cumulative lookup, follows identical arithmetic renormalization and updates exactly the same adaptive frequencies; models rescale at 4096 counts. No learned model sent in archive. This is genuinely byte-symbol adaptive coding, not another triangular reciprocal table on the existing bit model.

Nine framed components current **586,168 B**, candidate **597,212 B (+11,044)**. Component timings in one pass: old **232.36/201.57 ms** encode/decode, candidate **105.23/88.13 ms**. All nine expanded, all nine exactly inverse-restored. **REJECT placement.**

### H1C: decoder-synchronous hierarchical prior for first-use contexts

A distinct backoff mechanism seeds each new predecessor context's Fenwick histogram once from the global byte histogram built only from already decoded symbols, scaled conservatively; subsequent coding is unchanged. This adds no model metadata and preserves exact inverse. Nine framed components current **586,168 B**, new **598,202 B (+12,034)**; same-pass old **190.53/178.29 ms** versus new **93.50/98.48 ms**. Two `ooffice` blocks gained only 19 B and 289 B, but aggregate loss exceeded the unseeded coder. **REJECT this hierarchical initialization**; sweeping prior weights is inappropriate.

## Engineering checks and scope

- Original nine stream restores and incumbent native unpack: exact **9/9** for each of three variants.
- Unit controls included empty, constant, patterned, full-byte-cycle and randomized data; exact inverse for every case.
- Disposable trial `go test -count=1 ./...` PASS; `go vet ./...` PASS; targeted `go test -race -count=1 ./entropy -run '^TestByteAdaptiveExact$'` PASS; `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build` PASS (cross-build, not Windows runtime).
- Accepted root `default.pgo` remained byte-identical, SHA-256 `3b35cda74832360a252b63c1805e025aa1f4ed7a4fe4124d9a55929a6e86289f`.
- Component encoding/decode timings are **single diagnostic measurements**, not new ten-run whole-corpus medians. No complete candidate integration or complete-corpus candidate benchmark was run. None of the expanding candidate representations qualifies for addition under strict all-three-metric acceptance.
- Rejected Go implementation and diagnostic tests remain disposable; this report is research evidence, not new production source.

## Next mechanistically distinct opportunity

The performance cost is real, but the current bit-prefix probability model wins compressed bytes. Stop contextual rANS class-size sweeps, arithmetic per-context frequency cutoff/prior tuning and a redundant dual-encode fallback. To remove meaningful work while preserving bytes, investigate a substantively different source-stage mechanism: cache-local induced suffix ordering/BWT output **beyond Go's already-used SA-IS** (libsais is an architectural reference, never a compressor backend), or novel bounded-work discovery of mathematical exact/XOR/ADD references. Demand actual stage speed and bitstream equivalence before a complete integrated candidate. The previous prefix-doubling suffix builder and packed-LF inverse are already measured failures.

## H2: byte-renormalized binary range coding with original JLZ probability model

**Newly executed after H1; concrete native reversible output coder only, SCREENED and REJECTED.** Inspired by the explicit LZMA 32-bit bytewise renormalization described in `tukaani-project/xz/src/liblzma/rangecoder/range_encoder.h` and `range_decoder.h`, but implemented in original native Go, with JLZ's `bitModel.probability` and `bitModel.update` **unchanged**. Unlike H1B and the rejected EMA model, this does NOT change probabilities, contexts, bit order or match-token streams. It replaces the E3 bit-pending interval-output machinery with a 32-bit byte-normalized width/low/carry-cache arithmetic range coder, flushing the state at the end and validating exact consumed bytes and terminal state.

All **nine same original-source math-literal streams** restored exactly, independent of current native unpack. Using a coarse range cut `(width >> 16) * probability`, aggregate fully framed payloads **586,168→586,224 B (+56)**; a mathematically more accurate cut `(uint64(width)*probability)>>16` lowered overhead to **586,204 B (+36)**, **exactly +4 bytes in each of 9 tested streams**. One same-process sample with exact cut: old literal coder **163.07 ms encode / 161.42 ms decode**, new **86.51 ms / 116.39 ms** (timing diagnostic only). No experimental sample won the byte-cost gate, despite 1.9× encode and 1.4× decode component gains.

The 4-byte overhead is attributable to this particular native range coder's five-byte canonical termination. Explicit attempts to shorten that fixed termination while adding implicit decoder zero-padding failed exact byte restoration on original Silesia data. The successful five-byte-flush version retained precise truncation checks, and separate empty/constant/patterned/random small-and-large roundtrips and truncated payload rejection PASS. Further variations of finalization without a *proved and tested* shorter canonical representation risk weakened corruption handling; **do not force a noncanonical flush solely to meet a size target**.

The trade is now isolated: preserving current model quality, while altering arithmetic **output/renormalization only**, speeds the component but does not improve packed bytes. No full JLZ integration or ten-run candidate was justified under strict byte nonregression. No trial source committed; previous accepted decoder and profile unchanged. This H2 result is distinct from H1 static contextual rANS and byte-symbol Fenwick probability modeling. Investigate another actual algorithmic bottleneck, not probability-scale or flush-count tuning.
