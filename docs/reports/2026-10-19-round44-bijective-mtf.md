# Round 44 — Native bijective binary MTF zero-run stream (ACCEPT H1)

**Date:** 2026-10-19. **Decision: ACCEPT** as current native standalone JLZ codec, superseding R42. Encoding uses a reversible alternative representation **inside the already selected BWT blocks**. It does not replace BWT, introduce a third-party codec, or change the two-command interface. All lasting code is committed directly on main; GitHub Actions were not used.

## Rationale, theory and actual implementation

Round 36–43 exhausted several marginal encoder gating, match hot-loop, reciprocal and table optimizations; instead this round tests a genuine, complementary representation of an **existing winning stream**. BWT's move-to-front output has very long zero-rank streaks. Legacy `bwt.MTFEncode` emits each streak as byte zero followed by an unsigned varint length, then native Huffman/raw coding. The alternative treats each positive zero streak length as a **bijective base-two integer** with digits 0 and 1 (for length `r`, emit digit `(r-1)&1`, update `r=(r-1)>>1` until zero). Nonzero MTF ranks 1..253 map to bytes 2..254, while ranks 254 and 255 escape as `255,0` and `255,1`. This is a known RUNA/RUNB-like idea independently implemented in Go, **not a claim of a newly discovered compression theorem or proprietary invented coding principle**. Inversion accumulates `(digit+1)*2^position`, emits that many copies of current MTF-front symbol, then applies the ordinary MTF-list move for escaped/nonzero ranks. Limits and malformed escapes are checked.

Added:
- `bwt/zero_bijective.go` native reversible stream encoder/decoder and `bwt/zero_bijective_test.go` deterministic and corruption tests.
- `engine/bwt.go` encodes non-sample winning BWT blocks with the new stream, using existing `packStream`, existing BWT transform, and unchanged preflight. Existing legacy BWT mode **12** and its decoder remain fully readable. New **block mode 15** uses the bijective stream.
- `engine/cross.go` admits and decodes mode 15 inside the existing `JLZ1` mixed-cross framing. `engine/bwt_bijective_test.go` tests both legacy and new modes and truncated payload rejection.

**Compatibility:** Round44 decodes older JLZ1 archives through unchanged legacy paths, including old BWT mode 12. Earlier Round42 decoders do **not** understand the new mode 15 when present; record that forward-compatibility limitation instead of implying both directions. No change to checksum/framing, original-byte limit, CLI, standalone dependency policy, or worker quotas. BWT sampling/selection **remains unchanged** to isolate encoding format effects; the three already-encoded sample blocks in each winning BWT file are reused in old mode12 rather than rebuilt. The other BWT blocks use new mode15, and a raw fallback remains for blocks where encoding expands beyond the original block size. No file-name condition exists in production.

## Component viability: selected actual full source streams

Used **only** original 12-file Silesia (211,938,580 B); its 12 source SHA-256s match the canonical Round23 manifest. For the six files known to actually select BWT, re-encoded every complete 256 KiB BWT block with the existing and new *native packed stream*, verifying the new MTF decode against the exact BWT transform. Other six full files were accounted for as selecting other representations, not repurposed as a benchmark. Across **412 selected BWT blocks**, new native-packed MTF stream cost **20,533,358 B**, old cost **21,900,853 B**; **412/412 new streams are smaller** and theoretical block-selection opportunity is **1,367,495 B** before the identical surrounding archive framing. This is a component comparison, **not a complete archive result**. Eighteen prior old-format preflight samples (three per BWT file) are intentionally reused in the actual integrated codec, explaining why realized archive savings are smaller.

## Candidate-only complete original Silesia benchmark (ten runs)

Every counted run uses the **complete original twelve source files** (211,938,580 B), loads source data before timing, sums all twelve in-memory encode and decode durations, serializes each complete JLZ1 archive including metadata and SHA-256, then decompresses and checks every reconstructed original with `bytes.Equal`. All **120/120 full-file reconstructions passed** and every run yielded the same compressed total **61,418,211 B**. One bulk shell wrapper terminated after run 4 had *already finished and printed its PASS/total*; that fully completed run is included. The ten valid logs were collected in complete independent Go test processes with four effective Go workers, Go1.23.2 linux/amd64; no incomplete pass entered the median. Full-corpus input identities were separately checked using `sha256sum -c` on the canonical twelve-file manifest.

| Run | Complete archive B | Encode s | Decode s |
|---:|---:|---:|---:|
| 1 | 61,418,211 | 6.471087658 | 2.661320523 |
| 2 | 61,418,211 | 5.645186051 | 2.405479358 |
| 3 | 61,418,211 | 5.931221744 | 2.435392705 |
| 4 | 61,418,211 | 5.786176276 | 2.383360231 |
| 5 | 61,418,211 | 6.033435052 | 2.471567012 |
| 6 | 61,418,211 | 5.989858665 | 2.491078925 |
| 7 | 61,418,211 | 6.125051372 | 2.544262213 |
| 8 | 61,418,211 | 5.959930905 | 2.352605837 |
| 9 | 61,418,211 | 5.808888775 | 2.455148068 |
| 10 | 61,418,211 | 5.919944089 | 2.349189107 |
| **Median of complete corpus totals** | **61,418,211** | **5.945576325** | **2.445270387** |

### Required comparison

| Codec | Complete compressed size | Original data reduction | Encode time | Decode time |
|---|---:|---:|---:|---:|
| Zstd-3 (**retained historical external CLI**) | 66,216,569 B | 68.76% | 1.2561 s | 0.8551 s |
| XZ-6 (**retained historical external CLI**) | 49,233,340 B | 76.77% | 116.9892 s | 2.3079 s |
| **Current accepted JLZ R42 (retained ten-run median)** | **62,741,534 B** | **70.39636%** | **5.955825126 s** | **2.406640990 s** |
| **Candidate JLZ R44 (new ten-run median)** | **61,418,211 B** | **71.02075%** | **5.945576325 s** | **2.445270387 s** |

| Candidate compressed size delta vs R42 | Encode-time delta vs R42 | Decode-time delta vs R42 |
|---:|---:|---:|
| **-2.10917%** (-1,323,323 B) | **-0.17208% nominal** (neutral/unpaired) | **+1.60512% nominal** (unpaired, new decoder mode) |

Reference measurements are from different sessions and not a controlled paired speed comparison; the 0.17% nominal encoding change is **not** a demonstrated speed win, and the 1.61% decode difference is likewise only directional. External Zstd/XZ CLI times are not directly comparable to in-memory JLZ timing. Unlike R29's known source/provenance discrepancy, the 1,323,323 B reduction comes from actual R44 code relative to the correct current R42 archive total.

### Per-file compression regressions

Exactly six **selected BWT** archives shrink: `dickens` **-164,441 B (-4.74%)**, `mr` **-39,024 B (-1.35%)**, `nci` **-262,528 B (-10.69%)**, `reymont` **-103,461 B (-6.39%)**, `webster` **-713,237 B (-6.55%)**, and `xml` **-40,632 B (-7.04%)**. The other six archives `mozilla`, `ooffice`, `osdb`, `samba`, `sao`, `x-ray` are **byte-identical** to R42, as verified by comparison with captured source archive SHA-256 logs. **Worst relative size regression is zero**, no archive grows; largest absolute gain is webster. Mixed-cross archive metadata costs and new mode byte are included.

## Validation, disposition and next priority

- `sha256sum -c docs/reports/2026-10-18-round23-silesia-sha256.txt` **12/12 pass**.
- `go test -count=1 ./...` **PASS**; `go vet ./...` **PASS**.
- Targeted `go test -race -count=1 ./bwt ./engine -run 'TestBijectiveZeroRoundtrip|TestBijectiveBWTBlockCodecAndLegacy|TestBWTNativeRoundtrip|TestEarlyRegionalWorkerEquivalence'` **PASS**.
- `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build` **PASS**, build-only, no Windows runtime timing claim.
- Full ten-run Silesia decode **120/120 PASS**; candidate corrupt-stream/unit tests and legacy BWT block roundtrips passed.
- New source and regression tests committed with exact Git blob SHA equivalence to locally compiled/tested files. Temporary full-Silesia instrumentation/tests and candidate build binaries were **not** promoted.
- Candidate per-stream wins cover only BWT-selected Silesia data; new run encoding may lose on other unseen input distributions. Native BWT remains automatically gated as before and retains raw-block fallback. The original legacy encoder/decoder path still exists for preflight/backward archives.

**ACCEPT H1.** The new stream representation saves **2.11% total compressed bytes** at essentially unchanged measured median encoding time, a useful practical increase in ratio without another hot-loop microtuning cycle. Decode is directionally ~1.6% slower, but the new method is reversible, bounded, simple and a genuine additional native stream form. Current approved Round44 is the base for future tests. Next: instead of inventing extra BWT zero-run parameters, investigate whether this bijective-digit principle can **complement chosen arithmetic literal residuals** or other metadata streams, using cheap profitability screens first. Continue looking for transformative whole-corpus compression-speed improvements. 
