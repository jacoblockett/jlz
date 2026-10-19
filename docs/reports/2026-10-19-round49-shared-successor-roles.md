# Round 49 — Fast interleaved Huffman and decoder-synchronized shared successor states (REJECT)

**Date:** 2026-10-19. **Decision: REJECT H1** for the tested complete hybrid. Approved production baseline remains **Round 44**, source unchanged. Original full twelve-file Silesia only (211,938,580 B); all canonical source SHA-256 checks passed. No GitHub Actions or external compressor backend. Local disposable source was reconstructed from user-provided Round44 snapshot and the Round48 in-progress prototype ZIP; candidate implementations were not committed.

## Hypotheses and actual implementations

**H1A, conditional coder throughput:** Round48's single interleaved role-Huffman encoder emitted each canonical code one bit at a time; decoder scanned up to fifteen depths one bit at a time. Implemented a byte-buffered canonical-bit encoder, an 11-bit 2,048-entry prefix lookup for every context with exact canonical fallback on rare long codes and tails, and bypassed unnecessary online-model updates for fixed previous-byte low-bit roles. No wire-format change within the Round48 experimental codec-5. Uncached Go suite, vet, targeted races, Windows cross-build, and a complete 12-file exact screen passed; **60,821,703 B**, exactly the previous experimental Huffman candidate's archive hashes, **6.564944 s encode / 2.606809 s decode in ONE full-corpus screening pass** versus previous independent Round48 full screen **7.093449 / 2.681969 s**. Directionally faster but still slower than approved R44; *not a ten-run median* and not proof of causal improvement from unpaired timing.

**H1B, artificial meaning through learned shared structural state:** Implemented a different decoder-synchronized construction, not a partition into Huffman-role lanes. For each previous-byte context maintain a deterministic LRU of the last 2 or 4 *distinct successors*, optionally sharing one state across all predecessor bytes with the same high nibble (16 states). Before encoding a new symbol, perform a sequence of small exact transpositions that send cached successors into early byte-symbol positions 0..k-1; code the resulting exact residual-symbol bytes with JLZ's existing native raw/RLE/Huffman `packStreamCore`. The decoder applies transpositions in reverse order and updates its cache only after reconstructing the original symbol. No learned tables or side information. This is an invertible context-adaptive permutation, **not** an information-theoretically new source of entropy reduction; real coding improvements depend on the existing Huffman estimator and frequency concentration.

Reversibility verified for zero-length, repeated, patterned and mixed byte data for **four** models: exact-predecessor LRU2, exact-predecessor LRU4, shared-high-nibble LRU2, and shared-high-nibble LRU4. On all **412 selected complete R44 BWT/bijective symbol blocks**, old direct native packed size **20,533,358 B**. The four fully framed trial costs were respectively **21,446,428 B**, **22,167,328 B**, **21,384,082 B**, **22,033,855 B**: all lose, so **REJECT** shared-state swaps as BWT complements (just 1,186 B optimistic best-of-four selectable improvement). On **676 actual mathematical opcode/length streams extracted from fully roundtripped experimental archives**, direct native repacking cost **4,516,419 B** and tested full transformed costs LRU2 **4,294,840 B** (-221,579), LRU4 **4,221,917 B** (-294,502), shared-nibble LRU2 **4,507,637 B** (-8,782), shared-nibble LRU4 **4,452,209 B** (-64,210). Best-of-four oracle **4,133,629 B**, saving **382,790 B** as a component screen. These model selection and repacking totals are **not** JLZ archive savings, and the native repacked reference is not the precise original archived metadata size. All output streams were decoded byte-exact.

**H1C, additive full integration:** Implemented a new disposable experimental native packed-stream **codec 6** for automatic shorter-of-original-versus-LRU4 transformed and nested native-packed mathematical opcode/length streams, preserving codec 0–5 readers. Kept codec-5 fast interleaved Huffman for BWT-selected symbol streams. Only selected streams acquire new candidates, and exact fully framed payloads win strictly by actual size. The legacy R44 encoder/decoder source and public CLI are not changed. New candidate subtype needs a newer decoder; old R44 decoder does not understand codecs 5 or 6. Candidate decoder restored a preexisting archived R44-produced native CLI fixture (126,479 B source) byte-exactly. A math-only full-corpus *single-screen control* with native BWT entropy instead of codec-5 was inferior: **61,041,142 B / 6.555003 s encode / 2.757400 s decode**. The combined fast-Huffman-BWT + LRU4-math candidate outperformed that control.

## Ten complete original-Silesia candidate-only runs

Each run compressed and decompressed **every original byte of all 12 files**, including native archive metadata and digests; 120/120 restored source files match exactly. The same archive sum **60,747,358 B** was produced in every pass. Inputs were loaded before the encode/decode timers; times are sums of whole-file calls within each full-corpus pass. Four of the early invocations used the Go test wrapper and remaining runs used a compiled equivalent test executable after command-wrapper timeout; the test function, code and corpus are unchanged.

| Run | Archive bytes | Encode seconds | Decode seconds |
|---:|---:|---:|---:|
| 1 | 60,747,358 | 6.313014892 | 2.491038020 |
| 2 | 60,747,358 | 6.226914110 | 2.700564870 |
| 3 | 60,747,358 | 6.486430812 | 2.606448807 |
| 4 | 60,747,358 | 6.471364163 | 2.499907518 |
| 5 | 60,747,358 | 6.543614081 | 2.659751892 |
| 6 | 60,747,358 | 6.019975089 | 2.558627816 |
| 7 | 60,747,358 | 6.715626112 | 2.573043909 |
| 8 | 60,747,358 | 6.259856773 | 2.697295754 |
| 9 | 60,747,358 | 6.317087420 | 2.550796469 |
| 10 | 60,747,358 | 6.285349561 | 2.580292615 |
| **Median** | **60,747,358** | **6.315051156** | **2.576668262** |

Relative to original input, R49 reduces bytes **71.33728%**. Versus approved historical R44 **61,418,211 B / 5.945576 s / 2.445270 s** (ten-run median retained, not rerun), net archive change is **-670,853 B (-1.09227%)**, encode change **+6.21429% nominal**, decode **+5.37357% nominal**. These timings were taken in different measurement sessions, *not paired*, and exact performance attribution is uncertain. The earlier rejected R47 static model had **60,694,806 B / 6.186603 s / 2.579713 s** historically, meaning R49 is 52,552 B larger and nominally slower to encode even relative to a rejected prototype.

Full candidate archives on all twelve inputs are smaller than or equal to the approved R44 archive sizes. The six BWT-selected files have the same archive bytes as the prior Round48 interleaved screen; four math-selected files improve further by a combined **74,345 B** (mozilla 44,996; ooffice 1,582; osdb 10,334; samba 17,433). `sao` and `x-ray` remain byte-identical to approved R44. Worst per-file archive size regression versus R44: **0 B**. The model changes alter experimental candidate file hashes, but not previously accepted files.

## References and verification

| Compressor | Complete archive bytes | Original-byte reduction | Compress s | Decompress s |
|---|---:|---:|---:|---:|
| Zstd-3, historical external CLI | 66,216,569 | 68.76% | 1.2561 | 0.8551 |
| XZ-6, historical external CLI | 49,233,340 | 76.77% | 116.9892 | 2.3079 |
| Current accepted JLZ R44, historical 10-run median | **61,418,211** | **71.02075%** | **5.945576** | **2.445270** |
| Candidate JLZ R49, new 10-run median | **60,747,358** | **71.33728%** | **6.315051** | **2.576668** |

| Candidate size delta | Candidate encode delta | Candidate decode delta |
|---:|---:|---:|
| -1.09227% | +6.21429% nominal | +5.37357% nominal |

- Original corpus matched canonical 12/12 SHA-256 manifest.
- All 120 full-file reconstructions pass; component forward/inverse exactness checks pass.
- `go test -count=1 ./...`, `go vet ./...`, targeted `go test -race -count=1` and Windows/amd64 cross-build pass on the candidate. Windows binary was compiled, not run on Windows.
- Pre-R49 archive read through candidate CLI: approved R44 native fixture compressed with original R44 CLI and restored by experimental R49 CLI, checksum and exact comparison **PASS**. Ordinary project compatibility fixtures also pass in the unchanged suite.
- No accepted JLZ, Zstd or XZ reference compression benchmark was repeated. External CLI times are historical and use different I/O conventions than JLZ in-memory times.

## Decision and next direction

**REJECT H1**: the candidate is a meaningful *actual* archive reduction, but its >6% nominal encode slowdown and >5% nominal decode slowdown, new wire subtype and added state transitions are not justified by ~1.09% fewer bytes in a speed-first compressor, particularly given the independently measured rejected R47 model. Do not promote research codec 5/6 or its temporary code to production. **Round44 remains approved.**

**PRESERVE**: LRU state switching produced an actual ~294 KB selected-math component gain despite failing badly on BWT. Fixed-table Huffman can be sped up with a prefix decoder, but full candidate speed remains problematic. **Next structurally different tests:** compact joint models shared across math opcode/length and transform channels, and cross-block decoder-visible predictor state with sparse exceptions and hard MDL frame limits. Require a smaller *measured* full-stream cost without repeated candidate double encoding and test memory/time before new full-corpus promotion. Do not revisit hard low-bit group-count tuning or mislabel permutation as creation of information.
