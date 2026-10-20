# Encoder work-elimination experiments: transformed match word compares and finite-state early abort

**Date: 2026-10-20.** **Outcomes: REJECT both current implementations for promotion; no source changes on main.** This is original native-Go research based on [industry hot-path source audit](2026-10-20-industry-encoder-hotpaths-audit.md). No new JLZ release identifier, backend, copied existing compression implementation, or user-visible codec mode. Current complete-original-Silesia reference remains **60,585,717 archive B, historical ten-run median 5.493743735 s encode and 2.194427893 s decode**, which was **not rerun**. This report uses **screening runs**, not newly established candidate ten-run performance medians.

## Source provenance / environment

The experimental source was independently copied from the locally verified approved current Go implementation; changed only the relevant internal code in disposable directories, not the GitHub `main` branch:

- `/mnt/data/jlz_wordmatch_trial/engine/match.go` and `engine/word_match_equivalence_test.go`.
- `/mnt/data/jlz_fse_early_abort_trial/engine/finite_state.go` and diagnostic `engine/zz_fse_diagnostic_test.go`.

Main source `engine/match.go` Git blob **ff344090beddf2578a27af1354d51d90cac24b71**, `engine/bwt.go` **3f4af2d234c9382b2a3f9c20fd649974e0871e09**, and `engine/cross.go` **59794f52ac58b0712996a706b6c39b25d18ed7ca** remain unchanged. All **twelve original Silesia files** from `/mnt/data/corpus_silesia`, exactly **211,938,580 B**, matched canonical source SHA-256 manifest **12/12 PASS**. Tested with local Go1.23.2 linux/amd64; Windows cross-compiled only.

## H1: eight-byte XOR/ADD transformed match extension, output-equivalent (REJECT concrete implementation)

**Prior failed strategies contrasted:** Rejected fast exact-only dual-anchor parsing discarded transformed candidates and harmed ratio; rejected adaptive pricing increased encoder CPU; approved existing 64-bit *exact* comparisons were already implemented. This trial **preserved all exact/XOR/ADD candidates, token scoring, recent-distance lookup, hash insertion order, lazy matching and decoder bytes**, changing *only* verification of a transformed match's continued length after its existing five-byte anchor.

**Actual implementation:** Load 8 source bytes and 8 destination bytes with `binary.LittleEndian.Uint64`. For XOR, distribute the existing single-byte mask across eight lanes and XOR. For modular ADD, isolate even and odd byte lanes using `0x00FF00FF00FF00FF`, apply the same mask independently to four eight-bit lanes in each half without leaking carries across byte boundaries, rejoin lanes, and compare with the actual eight decoded bytes. Find the first differing byte with `bits.TrailingZeros64(diff)/8`, then retain exact original bytewise tail comparison. Does not extend past the original transformed max of 4096 bytes or block bounds.

**Correctness proof and component tests:** New `TestWordMatchAllMasksMismatches` verified **every 256-byte mask** for both transforms, mismatch at every one of 67 positions (including no difference through the tested region), and an additional **250,000 seeded randomized source/match cases** against the old scalar loop. All passed. Uncached full Go `go test -count=1 ./...` PASS, `go vet ./...` PASS, targeted Go race test PASS, CGO-disabled Windows/amd64 cross-build PASS.

**One original full-Silesia candidate screen and two complete repetitions** each compressed/decompressed all twelve original Silesia files exactly; every complete archive aggregate was **60,585,717 B**, and each of the twelve individual archive byte sizes matched approved current. The three separately measured candidate-only complete Silesia aggregate encode/decode times:
- **6.152612 s encoding / 2.374226 s decoding**.
- **6.468791 s / 2.546510 s**.
- **6.643046 s / 2.474894 s**.

These are **three one-run screens, NOT ten-run medians**; running another baseline was deliberately avoided. The independently timed approved reference was **5.493744/2.194428 s**, so candidate screens show **no reliable tri-metric speed win**. Because source files and archive bytes remain unchanged, a potential incremental encoder benefit cannot be claimed absent repeated robust data. **Do not promote or tune 8/16-word widths, masks, or test thresholds**; word compares may shift work from scalar comparisons into lane arithmetic and memory reads and can be counterproductive for short transformed matches. Keep the general principle of avoiding repeated work, but the concrete candidate is not an accepted optimization.

## H2: proven-losing finite-state coder early abort (REJECT concrete implementation)

**Industry mechanism:** Zstd estimates entropy-coded subblock costs rather than always completing expensive encodings. Current JLZ `packFiniteState` already packs `packStream` (raw/RLE/Huffman) before building its independent normalized finite-state model and **always encoding that second entire candidate**, even when it will lose.

**Actual bounded native proof-of-concept:** Create the same normalized FSE model and sparse serialized table, compute a rigorous **minimum byte count** for the candidate (known table length + 4-byte final state + codec byte + source-size varint + minimum 1-byte payload length varint + bytes ALREADY emitted during reverse-order rANS). When this minimum is `>= len(existingPackedStream)`, no future renormalization output can make FSE win. Abort its reverse coding after approximately each 1,024 encoded symbols in that situation; otherwise finish **identically**. This is logically exact: emitted bytes accumulate monotonically, and tied archives already favor the existing representation. No entropy approximations, no transmission changes, no probability modification, no decoder change. Optional diagnostic atomics were used ONLY during the separately timed instrumentation pass; they are not a production proposal.

**Empirical trigger prevalence on ONE complete-original-Silesia instrumentation screen:**
- FSE alternatives attempted: **1,642**.
- Emission-aborted proven losers: **81** (**4.9%**).
- Complete FSE bodies still encoded: **1,561**.
- FSE candidate selected as winner: **1,552** (**94.5% of attempts**).
- Remaining completely encoded loser attempts: **9**.
  
All twelve original files restored byte-for-byte, aggregate complete archived size **60,585,717 B**. A separate **uninstrumented ONE full-corpus screen** encoded in **6.828008 s**, decoded in **2.523180 s**; the instrumented screen encoded in **6.220535 s**, decoded in **2.540676 s**, but **instrumentation changes CPU costs, so these times are diagnostic and not 10-run promotional timing**. Compared to independently retained historical current median 5.493744/2.194428 s, neither screen supports a tri-metric win.

Targeted native engine/finite-state tests PASS, `go vet ./...` PASS, CGO-disabled Windows/amd64 build PASS. **Uncached full Go test suite for this instrumented trial did not finish within its attempted command's timeout**; do not report it as a pass. A previous targeted Go test including full original Silesia restoration passed. No need to burn ten candidate-only benchmark runs on an early shortcut which triggers only for ~5% of streams while imposing extra branches in almost every FSE encode.

**Do not tune polling intervals or histogram thresholds**. The intrinsic reason is structural: the finite-state coder already wins virtually every time it is attempted on JLZ's actual chosen data, so there is little redundant output work to eliminate. This is different from previously rejected cross-block FSE table reuse or multi-model segmentation; nonetheless this concrete safe early-abort branch is economically unhelpful.

## Conclusion and strategy ledger update

**Both implementations: REJECT promotion.** This is a failure to demonstrate speed nonregression, not a claim that word-wise transformed verification is wrong or that guaranteed early abort cannot be useful on other data. Identical archive bytes and original reconstruction are proved on the authorized original Silesia corpus; neither trial is a valid benchmarked ten-run replacement. Current and historical external benchmarks were not rerun; no corpus comparator/delta table is appropriate in user-facing communication. No experimental source committed.

**Next independent source-grounded direction:** Since the FSE coder already wins on **1,552/1,642** real attempts, further FSE alternative selection heuristics have poor room to save work. Stronger nonredundant opportunities would arise from avoiding unnecessary *model initialization* or the encoder-side decoder-only lookup-table construction, but these need a new profile showing sufficient CPU share—not blind parameter tuning. Alternatively study novel low-work match candidate access or codec architectures beyond the already failed parsing/BWT variants. Continue to audit established compressor implementations and use actual source-level differences to justify an implementation.
