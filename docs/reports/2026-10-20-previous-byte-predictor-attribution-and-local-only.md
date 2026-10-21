# Previous-byte predictive information attribution and native local-only E3 coder

**Date:** 2026-10-20. **Disposition:** REJECT the concrete decoder-compatible-new-subtype native **local-only adaptive arithmetic** integration due to small archive gain and large nominal encoding regression; **SCREEN/REJECT** static 256-way previous-byte Huffman conditional table projection before implementation. **No accepted production code or root PGO changed.** This research fulfills the dependency-level architecture audit demanded after the rejected E3 normalization scheduling trial.

## Precise novelty and nearest earlier research

Reviewed current main `AGENTS.md`, `FINDINGS.md`, `docs/STATE.md`, strategy ledger and placement matrix, and verified existing native `entropy/arithmetic.go`. Previous negatives already include: binary EMA replacing the existing adaptive counts, reciprocal probability LUT, original E3 deferred normalization, byte-normalized range coder with same probabilities, direct 256-context *byte-symbol adaptive Fenwick* coder, predecessor-high-bit static rANS and 8-way conditional bit-rANS, ZPAQ indirect state map, and ordinary static native 12-bit rANS on selected mathematical literal streams. No other experiment in the ledger was a previous-byte-only **existing prefix-bit-tree adaptive probability model that removes the current global-model blend and one full set of count updates**. That is this trial's structural change; it is **not** a novel industry entropy algorithm.

Current E3 coder selects an 8-bit-context `prev[previousByte*256+currentPrefix]` counter, **and** a 256-node global-prefix count, divides both to construct probabilities, and blends them at weights 1:7 or 5:3 depending on local count age. It updates **both** count tables after each of eight bit decisions. The new coder keeps the same E3 interval, bit I/O, exact previous-byte/prefix conditioning, local count initialization/rescaling, and eight bits per literal; probability is only `local.zero / (local.zero+local.one)`, with existing p0 clamp, and only the local counter is updated. Thus at least one integer division and one mutable update/cache access are eliminated per bit. Unlike the earlier adaptive byte-symbol Fenwick coder, the prefix-tree probabilities and arithmetic decisions are retained. The candidate has its own native literal subtype6; candidate decoder retains exact existing subtype3 for old archives.

A separate source-level context gap was assessed using a **static exact predecessor byte → symbol histogram**. Brotli-style static contextual symbol tables and already tried JLZ conditional static rANS are near precedents; here the measured question is the specific full 256-context model's ideal Huffman/data/header feasibility, not a repeat 8-class conditional rANS implementation. This second model was **a histogram/canonical-code length cost projection only, not a functional encoder/decoder**.

## Real original-source selection and model-information ablation

Source ZIP `jlz-6a388ba91682.zip` matches accepted production code `entropy/arithmetic.go` blob `08549011d624f06b67d93f863945e119d804e955`, `engine/stream.go` blob `572681332382bc1323224f4ec34ebb0c46884735`; live `main` had later documentation commits but unchanged production. The **original complete twelve-file Silesia**, 211,938,580 B, already verified against canonical source SHA-256 manifest and decoded bytewise. From *actual accepted compressed archives* of Mozilla, ooffice, and samba, extracted **all 262 currently chosen adaptive arithmetic literal streams**, **15,193,831 original literal bytes**, after fully decoding their original codec3 frames. All proposed new interval encodings were independently decoded to exactly the original bytes.

For each true selected stream, separately reset the current native global-prefix and previous-byte-prefix counters, stepped those counts in the **same actual byte/bit order** as production, and computed exact-model pre-quantization binary log-loss for current blend, global-only, and previous-byte-local-only. This is a model-dependency attribution, not the cost of an implemented archive mode:

| Selected mathematical literal stream group | Raw literal bytes | Bytewise zero-order H0 ideal B | H1 empirical previous-byte conditional entropy B | Actual E3 framed B |
| --- | ---: | ---: | ---: | ---: |
| Mozilla, 179 | 12,390,914 | 10,777,288 | 8,393,446 | 9,736,327 |
| Ooffice, 4 | 286,468 | 244,042 | 198,019 | 218,358 |
| Samba, 79 | 2,516,449 | 1,896,502 | 1,566,027 | 1,723,509 |
| **Total, 262** | **15,193,831** | **12,917,832** | **10,157,492** | **11,678,194** |

The **empirical H1** is an in-sample retrospective ideal, not achievable without learning/serializing conditional models and not evidence of a 1.52MB savings opportunity for free. Relative *pre-quantization bit log loss* on the same 262 streams: **current blended probability 11,676,279 B**, **global prefix-only 12,781,541 B**, **local predecessor+prefix-only 11,638,206 B**. Removing the previous-byte context loses **~1,105,262 B in this ideal diagnostic**, while removing the global blend **improves ideal diagnostic by ~38,073 B**. Previous-byte context carries nearly all useful conditioning in this selected domain. This ablation is separate from the full archive and no actual output was generated by it.

## Native adaptive local-only E3 component: functional inverse and speed

The experiment implemented native-Go `EncodeArithmeticLocal` and `DecodeArithmeticLocal` as a minimal variant of the accepted coder with only previous-byte prefix state. The full native 32-bit E3 interval, corruption-envelope rules and per-bit update order were retained. **262/262 selected actual-source streams inverse exactly**, current `EncodeArithmetic` also inverse-exact, including deterministic Mozilla/ooffice/Samba block selection:

| Source | Incumbent codec3 arithmetic body B | Local-only arithmetic body B | Local-only difference | Winning streams |
| --- | ---: | ---: | ---: | ---: |
| Mozilla | 9,735,087 | 9,685,498 | **-49,589 B (better)** | 105/179 |
| Ooffice | 218,332 | 220,381 | **+2,049 B (worse)** | 1/4 |
| Samba | 1,723,010 | 1,732,475 | **+9,465 B (worse)** | 59/79 |
| **Total** | **11,676,429** | **11,638,354** | **-38,075 B (better)** | **165/262** |

All measurements below are **single sequential component pass, order-sensitive and not ten-run JLZ timings**: current adaptive E3 encode **3.357052540 s**, local-only encode **2.624664212 s (-21.82%, better)**; current decode **3.370897460 s**, local-only decode **2.374687362 s (-29.55%, better)**. One entire 256-node global model update and global probability division were removed at each bit, plausibly explaining the local stage change. The two coders were measured in fixed incumbent-then-candidate order; independent profiling/paired repeats were not done.

## Complete native JLZ candidate, ten full original-Silesia runs

Integrated the **local-only model as the only full arithmetic-encoding attempt** in `packLiteralStream` instead of the incumbent's full mixed predictor, competing with incumbent `packStream` raw/RLE/Huffman envelope on actual byte size. New packed-literal subtype6 decodes with its independent synchronized model while accepting all existing subtype0–5 modes. Mathematical match/XOR/ADD search, BWT, BCJ, native field lanes, archive headers, worker model and cross-block histories are unchanged. This avoids a second whole-stream entropy pass.

**Every one of ten complete candidate-only runs compressed and decompressed all 12 original full Silesia files exactly**, totaling **120/120 full source restores, byte and SHA exact**, with deterministic total **58,720,441 B**. The 10 newly completed runs were collected in three independent test processes containing three, three and four repetitions, respectively, to avoid execution timeouts. Accepted reference measurements were **not rerun**.

| Model | Full original 12-file archive B | Median compression s (10) | Median decompression s (10) | Size versus current | Encode versus current historical | Decode versus current historical |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Accepted, historical ten-run | **58,767,739** | **4.7873709325** | **1.8721679705** | anchor | anchor | anchor |
| Native previous-byte-only adaptive E3, **new ten runs** | **58,720,441** | **5.507955719** | **1.802442080** | **-47,298 B (-0.08048%, better)** | **+15.052% (worse nominal)** | **-3.724% (better nominal, unpaired)** |

Entire archive compression reduction improves from **72.2713349311% to 72.2936517740%**, **+0.02231684 percentage points**. On a single detailed full pass, Mozilla archive **18,452,630→18,394,444 B (-58,186 B, better)**, ooffice **2,692,252→2,694,301 (+2,049 B, worse)**, Samba **5,134,262→5,143,101 (+8,839 B, worse)**, all other nine original files unchanged. This sums exactly to **-47,298 B**. Thus per-file regressions exist even though aggregate improves.

**All ten (encode, decode), seconds:** (5.656451636,1.848384326), (5.992395673,1.891811508), (5.556793261,1.822263087), (5.410289898,1.782552933), (5.489952169,1.808308265), (5.282833372,1.740369849), (5.381413919,1.718034636), (5.298274285,1.708184044), (5.525959269,1.861325628), (5.600979556,1.796575895).

Timing comparisons are **unpaired** with the retained accepted historical ten-run medians and exposed to variable machine and worker load. The real stage speedup does **not** imply a verified whole-codec Pareto gain, and the nominal 15% full encoder slowdown cannot justify promoting an additional persistent format for only 47 KB saved. A full-pipeline paired baseline ten-run comparison was deliberately not performed because JLZ instructions say to preserve and reuse the historical accepted reference, not repeatedly rerun it. The nominal decoder difference is encouraging but not a confirmed independent full-codec result.

Fixed external refs, **not rerun**: Zstd-3 **66,216,569 B / 1.2561 s encode / 0.8551 s decode**; XZ-6 **49,233,340 B / 116.9892 s / 2.3079 s**, with differing CLI/in-memory timing conventions.

## Additional conditional static-symbol feasibility, component PROJECTION only

The native codec's model dependency diagnostics motivate static exact previous-byte Markov tables as a way to process one symbol at a time. Before implementing a third coder, calculated actual **per-predecessor 256-symbol Huffman optimum lengths** from each of the 262 selected real literal streams. Proposed one 256-bit context-occupancy bitmap per stream, each enabled context's sparse (symbol,length) Huffman description (two bytes per nonzero), one context selection marker, seven bytes of optimistic stream framing, and existing literal symbols encoded under the chosen previous-byte context. Contexts for which Huffman bits + their description cannot beat fixed raw 8-bit symbols are optimistically assumed to fall back to raw without extra decoder-cost framing. The bit cost and tables are a **feasibility projection only, no decoder/bitstream generated**.

Across all 262 selected streams, conditional byte-symbol Huffman with this explicitly counted model table and per-context raw fallback projects to **13,736,680 B** compared to the **existing arithmetic bodies' 11,676,429 B**, or **+2,060,251 B (worse)**, despite empirical H1 being much lower. This is driven by high per-context side information and context-specific raw-fallback inefficiency. The projection is **optimistic**: decoder table reconstruction/switching, complete external archive subtype/edge framing and real Huffman bit finalization were not implemented or timed. No codec candidate/ten-run result is attributed to this path. The actual fully functional earlier byte-symbol adaptive Fenwick model and coarser static context-rANS are separate prior failures. **SCREEN/REJECT this particular static 256-way Huffman description**, not all possible compressed context-tree model representations or possible rANS/context coder designs.

## Validation, disposition, next novelty gate

Full native-Go local-only candidate `go test -count=1 ./...`, `go vet ./...`, targeted `go test -race -count=1` (current entropy/engine synthetic compatibility/corruption fixtures), and `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build` all **PASS**. Existing tests referencing literal codec3 were updated *in the disposable candidate only* to accept the new codec6; the fixture formerly requiring mathematical prediction mode15 instead selected a valid mode14 with raw literal codec0, so the test retained roundtrip and mutation checks without falsely mandating a different automatic winner. A separate first unit test still forces a new selected literal codec6. Two **real existing accepted** complete `mozilla` and `sao` archives decoded exactly using the candidate decoder (2/2), synthetic truncations at 0/1/7/17/half/last-byte and four bit mutations rejected or failed to restore exact original. Windows cross-compiled only; no Windows runtime, exhaustive new-mode fuzz, or race run over complete original Silesia. Unpromoted source exists only in a disposable working directory, no production .go files uploaded.

**REJECT this concrete local-only new wire mode**: -47,298 B and a promising isolated arithmetic speedup are outweighed by no demonstrated overall encoder gain, a nominal historical +15.1% encoder regression, archive expansion on ooffice/samba, and the decoder subtype/maintenance cost. No fixed acceptance threshold was used; the total evidence and format cost drive the decision. Do not blindly tune the local/global blend weight, first-bit special cases, or prefix model cap; such tweaks are not new architectures. The key robust discovery for subsequent source-level research is **the selected arithmetic streams rely heavily on exact preceding-byte context but not necessarily on the global prefix blend**. Any new lower-serial-work representation needs to preserve previous-byte conditional information and charge transmitted models fully.

**Accepted production unchanged:** original JLZ **58,767,739 B / 4.7873709325 s / 1.8721679705 s**, root profile SHA-256 **3b35cda74832360a252b63c1805e025aa1f4ed7a4fe4124d9a55929a6e86289f**.
