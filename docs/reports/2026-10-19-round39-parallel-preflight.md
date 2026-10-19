# Round 39 — Parallel deterministic early-regional preflight (ACCEPT)

**Date:** 2026-10-19. **Decision: ACCEPT H1.** Native JLZ1 format unchanged; output byte-identical for all original Silesia source files. Only the independent *first/middle/last* region probes are scheduled concurrently. The original eligibility tests, sample order, winner-selection thresholds, codecs, fallback behavior and decoder are unchanged. The single-worker path explicitly preserves the old immediate sample rejection instead of computing unnecessary samples.

## Why this work was targeted

Encoding-only CPU profiling (one **diagnostic**, untimed-for-reference complete original Silesia pass on current unmodified source) confirmed where CPU is spent; it intentionally did not repeat accepted JLZ/Zstd/XZ benchmark measurements. The profile measured **23.91 seconds sampled aggregate CPU across worker threads** with approximately **7.53 seconds elapsed**; notable cumulative hotspots included mathematical token parsing (`tokenizeMatchV2RangeStride`, 8.09 s CPU), BWT transform (6.51 s), native adaptive arithmetic encode (5.92 s), match-index find (5.00 s) and match-index updates. These CPU cumulates overlap with parent call stacks and must not be added together as independent work. The actual encoded Silesia archives comprised 6 files entirely BWT blocks (dickens/mr/nci/reymont/webster/xml), 4 mathematical-region files (mozilla/ooffice/osdb/samba), 1 additive-word file (x-ray), and 1 lane file (sao).

The confirmed avoidable wall-time critical path was **serial execution of three independently encoded 256 KiB mathematical/native sample pairs** inside `tryEarlyRegional` before the winning full-file representation can start. Unlike previous approximate gate/cutoff attempts, the same mathematical and native block encodings are still computed and their results are compared in original sample order; with two or three workers they can run concurrently. The existing native `workers.Available` respects Go `GOMAXPROCS` and the Linux container CPU quota (four cores available in the experiment). With one available worker, keep the exact old short-circuit on the first failing sample. No new compression mode or external codec.

## Measured candidate-only complete original Silesia

Inputs: **all 12 original Silesia files**, 211,938,580 original bytes, loaded before timed encode/decode operations; no alternative corpora were benchmarked. Tested on Go 1.23.2 linux/amd64, four effective CPU workers. A complete run sums **twelve encodes** and **twelve decodes**, validates `bytes.Equal` against every original source and checks the entire archive framing/embedded SHA-256 during decode. Each measured archive output is **62,494,395 bytes**. All candidate source archive SHA-256 hashes observed in full-corpus correctness audits match the **unmodified production source** byte-for-byte on 12/12 files, including the three kinds of selected representations and all metadata. The size is identical to the actual pre-Round39 production source; it differs from the *historical* Round29 benchmark number for the already-investigated, still-unresolved four-file provenance reason (Round38).

The first three complete measured runs were performed as independent single-process corpus passes. A subsequent monolithic invocation included one stalled/incomplete process; it was terminated and **excluded**. Seven further **completed independent single-process whole-corpus runs** provided runs 4–10. All ten successful runs are counted, with no partial or incomplete run included. The last single-worker-short-circuit cleanup changes only the `availableWorkers==1` branch, which is not taken on the four-worker candidate benchmark machine; full Go tests and targeted race checks were rerun after that cleanup.

| Run | Candidate complete archive B | Encode s | Decode s |
|---:|---:|---:|---:|
| 1 | 62,494,395 | 6.271518192 | 2.511425328 |
| 2 | 62,494,395 | 6.284540834 | 2.501863237 |
| 3 | 62,494,395 | 6.342707769 | 2.380185082 |
| 4 | 62,494,395 | 6.521590725 | 2.622749697 |
| 5 | 62,494,395 | 6.339756756 | 2.485531101 |
| 6 | 62,494,395 | 6.106162828 | 2.419315913 |
| 7 | 62,494,395 | 6.113280182 | 2.354914248 |
| 8 | 62,494,395 | 6.085760030 | 2.293306429 |
| 9 | 62,494,395 | 5.958031614 | 2.265365185 |
| 10 | 62,494,395 | 6.155383447 | 2.406023656 |
| **Median of whole-corpus totals** | **62,494,395** | **6.213450820** | **2.412669785** |

### Required comparison, historical references explicitly retained

| Model | Complete compressed archive | Reduction of original | Encode s | Decode s |
|---|---:|---:|---:|---:|
| Zstd-3, historical external CLI | 66,216,569 B | 68.76% | 1.2561 | 0.8551 |
| XZ-6, historical external CLI | 49,233,340 B | 76.77% | 116.9892 | 2.3079 |
| Approved JLZ Round 29, **historical** ten-run median | 62,644,663 B | 70.44% | 7.516364 | 2.821363 |
| Candidate JLZ Round 39, **new** ten-run median | 62,494,395 B | 70.51% | 6.213451 | 2.412670 |

| Archive-size delta vs **historical** Round 29 | Encoding-time delta vs historical JLZ | Decoding-time delta vs historical JLZ |
|---:|---:|---:|
| -0.2399%, **not attributable**; **0% vs actual current source** | **-17.33%**, unpaired historical-to-new | -14.49%, **not attributable**: decoder unchanged and runs unpaired |

The historical Round 29 medians were obtained independently and are not a controlled same-session paired comparison, and the older Zstd/XZ CLI measurements include different I/O and host conditions. The unmodified-source **diagnostic** whole-corpus encoding pass was ~7.53 s with profiling overhead; the candidate ten-run median of 6.213 s corroborates a substantial four-worker encoding gain without establishing a precise paired percentage. **Do not present the decode delta as a new decoder improvement.**

**File-level regression accounting:** all 12 archived file sizes, SHA-256 archive contents, original-file reconstructions, metadata and selected internal block modes are unchanged versus the actual pre-Round39 source. Consequently the worst archive-size regression is **zero bytes on every file**, and there are no wire-format compatibility changes. Per-file encoding-time deltas are not accepted as paired causal estimates from older benchmark sessions; negative one-file speed behavior on fewer available cores remains a possible environment-specific tradeoff, mitigated by preserving the serial early-abort path. No extra benchmark corpus was used.

## Validation

- Full uncached `go test -count=1 ./...` **PASS** after final cleanup, plus `go vet ./...` **PASS**.
- Targeted race checks `go test -race -count=1 ./engine -run 'TestEarlyRegionalWorkerEquivalence|TestParallel|TestBWT|TestAdditive'` **PASS**.
- `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build` **PASS** (build only; no Windows runtime claim).
- New `TestEarlyRegionalWorkerEquivalence` compares single-worker and four-worker sample selection/complete archive bytes, and checks decoding when the early archive wins. Entire original Silesia confirms identical archives in practice and 10/10 whole-corpus roundtrips.
- Existing production `JLZ1` decoder/format unchanged. Existing Go tests cover previously encoded archives. No GitHub Actions or branch experiments.

## Decision and next target

**ACCEPT H1:** The changed code and regression test were verified by exact Git blob hashes against the locally tested files before promotion. The output is provably unchanged for the complete Silesia benchmark, and encoding benefits materially from parallelizing work already required for winning regional candidates on the tested four-worker machine. Avoid multiplying expensive losing samples on low-core devices; the single-worker path keeps the old short-circuit.

Next: profile remaining serial *BWT* sample eligibility where a winning BWT candidate actually needs every sample, and distinguish the cost of preparation versus successful full-region encoding. However, **do not** commit to another gate-optimization spiral; prioritize genuine architectural speed opportunities, and stop early if a one-pass screen has no clear win. Archive-size baseline discrepancy from R29 remains a historical provenance issue, not R39 savings.
