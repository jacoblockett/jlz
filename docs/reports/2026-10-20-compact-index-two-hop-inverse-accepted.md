# Accepted byte-identical performance: compact matcher buckets plus two-hop BWT inverse

**2026-10-20. ACCEPT H1.** Accepted production engine on `main` commit `e75bb094e518c8f0a4c74731f0d517143cf66989`. Both changes preserve all twelve original-Silesia archive bytes exactly and all historical JLZ1 decoder modes. No new mode, compression switches, external codec, GitHub Actions, or unprofiled build. Committed original Go1.23 root `default.pgo` retained byte-for-byte.

## Architectural analysis and novelty

The prior complete-Silesia Go CPU profile (one diagnostic) found cumulative encoder CPU in native BWT forward transform ~32.55%, math tokenizer ~25.10%, arithmetic encoder ~22.06%; decoder arithmetic ~47.56%, BWT inverse ~23.31%. Current matcher has three 65,536-bucket tables, storing each bucket as two `uint32` positions plus a `uint32` ring counter (12 bytes). The two-entry candidate limit itself was already accepted historically; failed adjacent strategies tried increasing/decreasing depth, generation stamps, fingerprint tags, lookup micro-arithmetic and word-wide match checking. Here the **same two newest candidates are represented in a single 64-bit bucket** (newest low 32 bits, second-newest high 32 bits). Its exact candidate positions, candidate order, score, parser, distance-repeat cache and encoded token stream stay identical, while each table becomes 8 rather than 12 bytes per bucket. With three tables this removes **786,432 bytes** of matcher index storage per concurrently active transformed-matcher instance. Source changes: `engine/match.go`, independent oracle adaptation in `engine/match_search_equivalence_test.go`, and `engine/match_bucket_test.go`.

Previous packed LF inverse (historical larger 512KiB BWT experiment) represented **one symbol and one LF transition per entry**, with no complete-codec speed benefit. The new **two-hop inverse** is not another one-hop repacking: it computes *two* successor symbols and the row after *two* LF transitions in a 64-bit entry. Each dependent random table fetch in the final inverse loop restores two source bytes, halving serialized dependent row traversals. No new wire metadata, encoder transform or codec dependency. Source change: `bwt/transform.go`, exact random/edge verification `bwt/two_hop_test.go`. Tradeoff: two-hop table adds 8 B per BWT row during decoding, about **2.0 MiB temporary additional allocation per full 256KiB block** over the original single-hop arrays. Decoding still bounds original archive region sizes at 256KiB and respects existing worker limits. This CPU-for-temporary-memory trade is disclosed, not treated as free.

The layout simplification is consistent with fast native index work in Zstd/LZFSE; the two-hop LF traversal is JLZ's own reversible arrangement. No foreign compressor source was used as a backend or copied.

## Separate feasibility checks

- Matcher **alone** integrated into full JLZ with unchanged decoder: ten complete original-Silesia candidate-only medians **60,244,257 B / 5.0814215205 s encode / 2.054880984 s decode**, all 120 roundtrips exact, source archive per-file sizes identical. Nominal decoder +2.56% versus *unpaired* historical reference despite unchanged decoder program, illustrating noise and why the improvement was **retained as a research candidate, not prematurely rejected or promoted**.
- BWT inverse two-hop component on **all 451 actually BWT-selected 256KiB original source blocks** from seven complete source files (**117,233,324 original bytes**) reconstructed exactly. One component-only stage measurement: original `bwt.Inverse` **1593.358 ms cumulative**, new two-hop **1306.444 ms** (about -18.01%). This is a single component diagnostic, not a whole-codec speedup estimate.
- The additive combined full JLZ candidate's first screening pass **60,244,257 B / 5.151875999 s encode / 1.997370634 s decode**, exact 12/12. The subsequent **ten complete** candidate-only passes below determine promotion.

## Full original 12-file Silesia acceptance evidence

Original source size **211,938,580 B**; canonical source SHA-256 checks **12/12 PASS**. The combined native engine compressed and decompressed each complete original file for every run; **120/120 byte- and SHA-exact full-file reconstructions PASS**. The complete archived size was **60,244,257 B on all ten completed runs**, and *all twelve exact archive byte streams* were independently checked with the previous production executable and the combined candidate to be **byte-for-byte identical** (content verification, not a new baseline timing benchmark). Current `default.pgo` root SHA-256 `3b35cda74832360a252b63c1805e025aa1f4ed7a4fe4124d9a55929a6e86289f` unchanged; executable compiled with Go1.23.2 and explicit `-pgo=./default.pgo` for timed candidate. Normal Go build also incorporates root profile automatically.

Ten completed isolated candidate processes, each checking the entire original 12-file corpus; one separate fifth invocation from an earlier batch stalled and was terminated, so **only completed passes are counted**, and no incomplete observation is included. No established JLZ, Zstd or XZ reference timing was rerun.

| Candidate run | Archived bytes | Encoding seconds | Decoding seconds |
|---:|---:|---:|---:|
| 1 | 60,244,257 | 4.965883334 | 1.877227599 |
| 2 | 60,244,257 | 4.954150782 | 1.976146715 |
| 3 | 60,244,257 | 5.144104691 | 1.874845032 |
| 4 | 60,244,257 | 5.084467556 | 1.861725394 |
| 5 | 60,244,257 | 4.837308184 | 1.890055191 |
| 6 | 60,244,257 | 5.079357323 | 1.919312498 |
| 7 | 60,244,257 | 4.983635720 | 2.136197459 |
| 8 | 60,244,257 | 5.083713561 | 1.943662955 |
| 9 | 60,244,257 | 5.174664785 | 1.972728033 |
| 10 | 60,244,257 | 5.170770415 | 1.883307118 |
| **Median** | **60,244,257** | **5.081535442** | **1.9046838445** |

Historical accepted comparator (not rerun): **60,244,257 B / 5.322679488 s encode / 2.003659634 s decode**. Candidate versus current compressed size **0 B, 0.000%**; encode median **-4.530501%**; decode median **-4.939751%**. Original reduction remains **71.574662%**. Historical other reference rows: Zstd-3 **66,216,569 B / 1.2561 s / 0.8551 s** and XZ-6 **49,233,340 B / 116.9892 s / 2.3079 s**, external CLI work not directly comparable to local JLZ in-memory wall times. Unpaired runs from distinct sessions mean the precise percentage differences are host/environment-specific, not universal speed guarantees.

## Engineering validation and compatibility

- All **12/12 canonical original Silesia source SHA-256** PASS; full candidate archives equal prior current **byte-for-byte for each individual original file**.
- Uncached **`go test -count=1 ./...` PASS**, **`go vet ./...` PASS**, targeted **`go test -race -count=1 ./engine ./bwt`** covering changed unit tests PASS, **`CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build` PASS** (cross-build, not Windows runtime).
- Direct backwards decoding of *both* old accepted **full `mr` BWT archive** and prior **64KiB math fixture archive** PASS by byte equality and embedded SHA integrity.
- Added tests cover compact two-slot ordering against prior ring representation, random/edge 0/255 bytes and lengths up to 256KiB BWT forward/inverse equivalence, invalid BWT primary metadata and existing engine malformed/truncated payload tests. No new format codes or decode metadata.
- Fresh GitHub `main` commit blob hashes for both source files and all three tests independently **match exact locally tested candidate file Git blobs**. Verified unchanged root `default.pgo` blob. The final GitHub tree is source-aligned with the locally executed candidate; the committed tree was also read back from GitHub and verified. GitHub-hosted tests were not run because Actions are disabled.

## Decision and next distinct architecture

**ACCEPT H1:** equality of complete archive bytes, actual 10-run median improvements on both enc/decode times, backward compatibility and confined code justify production. This is a performance optimization **not a compression-ratio breakthrough**; avoid misrepresenting unchanged 71.574662% ratio as progress toward 99% from this change. Prior matched-offset layout tuning and one-hop packed LF negatives remain valid, but do not disprove this separate two-hop method.

Future priority returns to major compressed-byte saving and throughput: target the actual math-selected `mozilla`/ `samba` bytes and BWT-selected `webster` with orthogonal reversible representations. Continue studying established codec architectures but prototype as additive layers inside existing JLZ, benchmark only complete-original Silesia candidates, and require all-three-metric preservation absent explicit trade permission. Avoid repeating previous small-depth/width/cache threshold sweeps.
