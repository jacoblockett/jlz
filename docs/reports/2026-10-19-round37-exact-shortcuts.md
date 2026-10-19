# Round 37: Exact-output performance shortcuts (2026-10-19)

**Verdict: REJECT H1 (no promotable encode-speed gain).** All production experiment changes were discarded locally; no Go codec sources were changed on main. Original Silesia corpus only. All 12 original bytes and source hashes match the canonical Round 23 Silesia manifest.

## Scope and controlled constraints

Investigated the five shortcut families requested after the Round 36 audit: (1) arithmetic stream early abort, (2) parallel losing-region cancellation, (3) native Huffman mathematical lower-bound rejection, (4) reuse of expensive regional/additive samples, and (5) avoid redundant stream copies and internal SHA-256. No external codecs, new wire modes, CLI flags, or GitHub Actions were used. Candidate source derived from the user-uploaded `jlz-main.zip`, source commit `708f864`; applied only the already accepted fixture path correction from `28865fb`. GitHub comparison verifies changes from that snapshot through pre-experiment main `b1d450d` are documentation and this fixture-path test fix, **no production Go source changes**.

All candidate runs were complete original 12-file Silesia (211,938,580 B) and 12/12 byte-exact reconstructions, except where a test was interrupted before completion; partial attempts never counted. Individual SHA-256 source identities match `docs/reports/2026-10-18-round23-silesia-sha256.txt`. Experimental results are in-memory Go tests. Separate single-process and twelve-independent-process screens have different lifecycle/GC behavior and MUST NOT be mixed into one median.

## Implementations and measured screening observations

1. **Arithmetic early-stop (exact bound):** Extended the independent native adaptive arithmetic encoder with a bounded-output function. After each 256 input bytes, if *already emitted* arithmetic bytes reach the complete alternative stream length, abort and retain the cheaper raw/RLE/Huffman stream. There is no chance for previously written bits to shrink. Combined with removing one redundant raw copy: **62,494,395 B**, **7.249863 s encode / 2.623309 s decode** (one whole-Silesia screen). Exact full restoration; no substantial speed gain.
2. **Parallel region cancellation (exact bound):** Added atomically tracked payload-byte totals to BWT, lane and mathematical-region worker pools; when the sum of completed nonnegative region payload sizes alone reaches/exceeds the winning complete archive size, unsent/undone worker tasks are abandoned. The existing final framing-size check remains. **62,494,395 B**, **7.248878 s / 2.399562 s** (one 12-independent-process complete-corpus screen), 12/12 restoration. A separate *correctness-only* audit proves all 12 archives byte-identical to the unmodified source. No convincing encode gain.
3. **Native Huffman collision-entropy lower bound:** Reused a precomputed byte histogram and tried rejecting Huffman construction when the collision entropy (a theoretical lower bound on Shannon entropy, hence on Huffman average bit length) plus table costs cannot beat the existing stream. **62,494,433 B**, **7.371727 s / 2.383866 s** (one complete corpus in 12 independent processes). This differed **+38 bytes** from the unmodified source across two files and **FAILED the strict archive-byte preservation requirement**. The prototype's implementation/bounds need investigation before any reuse. Rejected with no promotion.
4. **Distributed sample reuse:** Cached original mathematical regions from the initial early-regional preflight, retaining them for any later size comparison, and cached three actual additive sample encodings to avoid recomputation in the complete additive pass. **62,494,395 B**, **8.519755 s / 2.751267 s** (one 12-independent-process Silesia screen). No meaningful gain. A monolithic per-corpus invocation stalled partway through; the independent-process run completed all twelve and passed exact reconstruction.
5. **Avoid redundant copies/hashes:** Removed the unnecessary `append([]byte(nil), data...)` in the raw-stream candidate because `bytes.Buffer.Write` already copies the eventual payload and all compressed selections hold their own bytes. The independent single complete-corpus screen was **62,494,395 B**, **7.122553 s / 2.621038 s**. Also replaced discarded internal block SHA-256 calculations with zero digests where the temporary outer archive header is immediately stripped (the digest is not serialized to the enclosing block): combined with the copy change, one screen was **62,494,395 B**, **7.597983 s / 2.713942 s**. The copy-only variant had **12/12 byte-identical archive SHA-256** to the unmodified source, and passed Go tests. Hash elimination did not improve the measured screen.

A combined prototype (arithmetic abort + copy removal + sample caching + region cancellation + Huffman gate) had one full-corpus screening measurement **62,494,395 B / 7.545327 s encode / 2.549820 s decode** and all 12 correct restorations. One later Huffman-only variant changed the archive by 38 B, so its exact-output implementation is rejected separately. All metrics above are *single-run screening only*, NOT ten-run promotional evidence. Comparing individual run times across independently executed processes is noisy.

### Best candidate: copy-only additional repeats

Whole-corpus **single-process-per-run** screens, complete and exact:
- Run 1: 62,494,395 B, 7.575849 s encode, 2.612104 s decode.
- Run 2: 62,494,395 B, 7.433439 s encode, 2.645072 s decode.
- Run 3: 62,494,395 B, 7.495602 s encode, 2.592768 s decode.
- Three-run screening medians: **7.495602 s encode, 2.612104 s decode**.

A second set using **twelve fresh processes per corpus pass** completed three further, distinct screening-only repetitions (each 12/12 exact): encode **7.649016 / 7.499712 / 8.145337 s**, decode **2.524711 / 2.616683 / 2.706201 s**, size 62,494,395 B each. These are *not* combined with the three whole-process medians due to different execution methods. Several extended monolithic test attempts hung and were abandoned, not counted. Neither screening series establishes a material encode gain. No ten-run promotional benchmark completed; stopped early under the user's anti-diminishing-returns directive.

### Critical provenance discovery: current source vs recorded Round 29 size

An **untimed, correctness-only 12-file SHA-256 archive comparison** between the *unmodified uploaded source* and the copy-only candidate proved 12/12 **bit-for-bit identical archives**. The independent cancellation candidate likewise matched all twelve original archive hashes. The actual unmodified source thus already generates **62,494,395 B** on canonical complete Silesia. This is **150,268 B smaller than the historical Round 29 accepted total of 62,644,663 B**. The apparent size improvement is **NOT caused by any Round 37 shortcut**. Both figures must remain provenance-labeled until the historical source/benchmark discrepancy is resolved. The uploaded source is the exact production source of pre-Round37 main except a fixture-path test fix and docs. Do not claim a new accepted 150,268 B ratio improvement, and do not silently overwrite the accepted historic Round 29 reference.

## Required historical-versus-screen comparison (not promotional)

| Model | Full compressed archive | Reduction from original | Encode | Decode |
|---|---:|---:|---:|---:|
| Zstd-3 (historical CLI) | 66,216,569 B | 68.76% | 1.2561 s | 0.8551 s |
| XZ-6 (historical CLI) | 49,233,340 B | 76.77% | 116.9892 s | 2.3079 s |
| Current approved JLZ Round 29 (**historical** ten-run record) | 62,644,663 B | 70.44% | 7.516364 s | 2.821363 s |
| Candidate copy-only (**new source**, 3 whole-file-pass screens) | 62,494,395 B | 70.51% | 7.495602 s | 2.612104 s |

| Archive bytes change vs historical JLZ | Encode time change vs historical JLZ | Decode time change vs historical JLZ |
|---:|---:|---:|
| -0.2399% **not attributable (source already identical)** | -0.2762% nominal, **not demonstrated** | -7.42% nominal, **same decoder, not attributable** |

Historical Zstd/XZ CLI runs have different host/methodology. The prior JLZ reference was a ten-run median and these are only three whole-process candidate screens, across independent sessions. Do not interpret apparent speed deltas as causally established improvements. The source snapshot discrepancy is open.

## Validation and decision

All original complete Silesia files restored byte-for-byte in the completed screens; 12/12 current-source-vs-copy-only archive SHA-256 match and 12/12 current-source-vs-cancellation archive match. Copy-only changes pass `go test -count=1 ./...`, `go vet ./...`, targeted Go race tests, and Windows/amd64 cross compilation; no runtime Windows benchmark. Histograms/Huffman variant fails the required archival byte identity check (+38 B). Only local temporary prototype directories were modified. No Go source will be promoted.

**REJECT H1:** The shortcuts are mostly conceptually valid but either already partly implemented, too late to cancel significant work, or not demonstrably substantial on Silesia. Preserve these negative findings. Before further speed tuning, reconcile the current source's 62,494,395 B with the published accepted Round 29 historical 62,644,663 B. Next serious targets should address *winning* BWT suffix construction and actual match parsing rather than endless gating/threshold tweaks.
