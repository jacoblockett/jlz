# Industry-leading compressor source audit: hot loops, early cost gates, and concrete gaps in JLZ

**Date:** 2026-10-20. **Type: SOURCE-LEVEL AUDIT ONLY.** No new JLZ codec implementation, no benchmarked candidate, no accepted change. Current unchanged on `main`: complete original twelve-file Silesia **60,585,717 B**, retained ten-run medians **5.493744 s encoding and 2.194428 s decoding**. No new numbered research rounds. This is a **second, specifically fast-path/code-level audit**, not a repackaging of the earlier high-compression parser study.

## What was actually examined

Read concrete upstream C/C++ encoder source of **six established compressors**, not just marketing descriptions, format PDFs or general algorithm names. GitHub main/dev/master source snapshots recorded with content SHAs for reproducibility:

- **Zstd**: `facebook/zstd` `lib/compress/zstd_fast.c` (blob `ee25bcbac8d8eb3881ff33b17a30aef85a6b961e`), `zstd_compress_internal.h` (row tags, lazy skipping), `zstd_compress_superblock.c` (blob `a1b4192cc3a6ea3e36ed6073f80b079b1d0e44f8`), `zstd_ldm.c` (blob `7bf36cd8655f5870b466e93f900d5563f925e799`), `zstd_compress.c` (blob `38c3fc0d4cde208fbdb12408e6ef2008ab254a67`).
- **LZFSE**: `lzfse/lzfse` `src/lzfse_encode_base.c` (blob `367d7a23bb0817e48eca0801e7cb9886202652b8`), and decoder/header organization in `src/lzfse_internal.h`.
- **libdeflate**: `ebiggers/libdeflate` `lib/deflate_compress.c` (blob `b24087c20c2dca727785c5a71aed6c5fc696a73e`): fastest versus lazy/near-optimal policies, sequence store limits, 8 literal + 2 match observation classes, 512-observation block checks and optimization cost structures.
- **LZ4 HC**: `lz4/lz4` `lib/lz4hc.c` (blob `3234df369393f53cc7f70e3f707982e6ac0e521c`): insertion and wider-match search, quick mismatch tests, repetitive-pattern counting and backward extension, decode-speed bias.
- **Snappy**: `google/snappy` `snappy.cc` (blob `add72f876980c58007f196ed897c575b5972580c`): pipelined hash lookup, direct emit after match, subsequent-index updates, match length counting and overlapping copy cost.
- **libbsc**: `IlyaGrebnov/libbsc` `libbsc/lzp/lzp.cpp` (source 47 KB), `libbsc/coder/qlfc/qlfc.cpp` (111 KB), `libbsc/libbsc.h`: LZP prediction before BWT/block sorting and alternative quantized local-frequency coding. Source pipeline is **not** synonymous with JLZ's plain BWT→MTF→entropy. Broader BWT-adjacent experiments previously had large time penalties; novelty exists but requires exceptional proof.

Primary source links:
- https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_fast.c
- https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_compress_superblock.c
- https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_compress_internal.h
- https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_ldm.c
- https://github.com/lzfse/lzfse/blob/master/src/lzfse_encode_base.c
- https://github.com/lzfse/lzfse/blob/master/src/lzfse_internal.h
- https://github.com/ebiggers/libdeflate/blob/master/lib/deflate_compress.c
- https://github.com/lz4/lz4/blob/dev/lib/lz4hc.c
- https://github.com/google/snappy/blob/main/snappy.cc
- https://github.com/IlyaGrebnov/libbsc/blob/master/libbsc/lzp/lzp.cpp
- https://github.com/IlyaGrebnov/libbsc/blob/master/libbsc/coder/qlfc/qlfc.cpp

Inspected current JLZ accepted Go source `engine/match.go`, `match/expand.go`, `engine/finite_state.go`, `engine/stream.go`, `engine/packed.go`, `engine/region_fast.go`, `engine/slots.go`, and compared [strategy ledger](../research/strategy-ledger.md). No foreign source was copied or used as an encoder backend.

## Specific mechanisms, exact code-level mismatches, and tri-metric feasibility

### Opportunity A — fast EXACTLY EQUIVALENT transformed-match verification (LZFSE + LZ4HC)

**Upstream mechanics:** LZFSE `lzfse_encode_base` extends equality runs with repeated `load8` reads and XOR, then uses `__builtin_ctzll(diff)>>3` to find the first mismatching byte within a 64-bit machine word. LZ4HC likewise uses machine-word equality checks and common-byte counting for exact and repeated-pattern discovery. These are data-path optimization techniques, not changes to the archive grammar. This audit **does not claim LZFSE supports JLZ's constant-mask XOR/ADD transformed references**; that extension is our own mathematical proposal.

**Exact JLZ mismatch:** `engine/match.go` `matchIndex.find` already compares **exact** candidate tails eight bytes at a time using `binary.LittleEndian.Uint64`, but **transformed XOR and ADD** references still extend candidate lengths with a **single-byte loop, up to 4,096 bytes per candidate**. The existing `match/expand.go` already uses Go bulk `copy` for exact decode, so there is no novel copying optimization to propose there. This is specifically the *forward transformed match verification* hot loop.

**Original native opportunity:** use a replicated byte mask to compare eight candidate bytes simultaneously for XOR, then `math/bits.TrailingZeros64` to locate the first mismatch; for modular bytewise ADD, use a *carry-isolated pair of four-byte lanes* (odd and even bytes) to prevent 8-bit arithmetic carries leaking between bytes. Handle partial tails, distance overlaps, source lengths, arbitrary masks and architecture-independent little-endian reads identically to current code. These operations need a correctness proof comparing each transformed candidate's **exact same match length, mask, score, and subsequent token decision** across all actual byte cases, not just roundtrip correctness.

**Expected property if successful:** identical match choices and compressed archives, decoder untouched, potentially faster encoding; no extra entropy passes, window size, metadata or source scans. Strict promotion still requires median compression **nonregression AND** median decompression **nonregression**; unchanged decoder is not evidence that noisy measured decode times always tie. **This is the highest-ranked NEW code-local speed hypothesis.** No speed gain measured yet.

**False novelty to avoid:** LZ4HC and JLZ already extend matches *backward through pending literal runs*; JLZ `tokenizeMatchV2RangeStride` does this. Do not propose it as unimplemented.

### Opportunity B — don't perform losing entropy encodes (Zstd cost-estimate/early-mode decisions)

**Upstream mechanics:** Zstd `ZSTD_estimateSubBlockSize_literal` and `ZSTD_estimateSubBlockSize_sequences` estimate output from symbol frequencies and model tables, including entropy-description size and block framing **before emitting each candidate sub-block**. Zstd separately considers raw/RLE/compressed/repeated representations, and can avoid unnecessary expensive work if a frame cannot profit. This isn't claimed to be a universal exact proof of no win.

**Exact JLZ mismatch:** `engine/finite_state.go` `packFiniteState` **always completely constructs `packStream(data)`** (raw/RLE/Huffman), then for every stream of >=64 bytes **builds the normalized 12-bit finite-state model and encodes a second complete native rANS bitstream** before comparing `candidateLen` to `len(best)`. This includes opcode and length streams in `engine/packed.go` and selected BWT symbol streams. Both full bitstream encodings are currently executed regardless of predicted cost. The existing `huffmanCannotBeat` only skips some **whole-regional native fallback** attempts, not the nested `packFiniteState` race. Prior rejected multi-BWT neighborhood encodes concerned selecting alternate models, not skipping this exact unconditional incumbent-vs-FSE work.

**Original opportunity:** instrument on *actual currently selected* Silesia streams to count how often FSE loses, bytes saved by incumbent, histogram/model costs and encoder CPU. Then derive a sufficiently safe **lower bound of the specific finite-state coder's actual emitted byte length** (normalized symbol frequencies + state/renormalization framing) relative to the known incumbent byte length to skip losing FSE encodes. Conservative full fallback on uncertain cases. **Do not treat zero-order Shannon H0 as a mathematically sufficient lower bound on every particular sequence's coded length**, since a cheaper-than-H0 individual codeword is possible. Test byte-for-byte equivalence of archives, not just original restoration. Avoid an estimator more expensive than the passes it saves. If the actual FSE-loss fraction is small, **STOP**; this would be needless tuning.

**Potential property:** unchanged archive and decoder, reduced encoding work, no extra format, compatible with strict three-metric criteria. This is the second-ranked candidate for **measurement**, distinct from historical table-reuse (model *size*, not encoding *work*) and the failed round of arithmetic reciprocals.

### Opportunity C — hide hash lookup latency without deleting candidate matches (Zstd fast / Snappy)

**Upstream:** Zstd `ZSTD_compressBlock_fast_noDict_generic` explicitly interleaves hash/table/match loads from adjacent positions to break serial read dependencies. Snappy `snappy.cc` similarly stages hash lookups, candidate comparisons and follow-on match checks while amortizing table updates. Zstd separately offers `lazySkipping` to omit some match-index updates over incompressible runs **at a documented small ratio cost**.

**JLZ mismatch:** `matchIndex.addBefore` updates all skipped positions sequentially, up to **three 65,536-bucket tables, two retained positions each**; `find` probes recent distances and up to six transformed/raw hash candidates in a serial path, then occasionally probes next position for a one-byte lazy decision. Potential CPU latency and redundant fingerprint recomputation exist. A state-preserving and candidate-equivalent input-row pipeline might hide memory latency without modifying archive decisions.

**Critical gate:** Zstd's literal strategy of dropping hash insertions is **not permissible** by default under JLZ's strict archived-size nonregression; it can lose good references and the source documents that trade. Interleaving Go lookup/compare phases while preserving **all current match-index update order, candidate positions, transform flags and repeat-cache behavior** is difficult and may not be worthwhile. Require measured per-stage cycle/CPU cost before coding a complex pipeline. **Third priority, conditional** on A/B actual impact.

### Opportunity D — variable instruction/entropy table framing, not yet justified (LZFSE / libdeflate / Zstd)

- LZFSE compact block v2 headers bit-pack state fields and encode table frequencies with a small prefix code; JLZ native finite-state header stores each used symbol and its normalized frequency as one byte + uvarint. There may be a **byte-size** opportunity, but it requires a new internal mode or explicit table subtype decoder and incurs new decode work. Previously selected model sharing saved only ~5 KB on component oracle, and previous context schemes generally lost runtime. **LOW priority pending actual current table-byte share measurements**.
- libdeflate's statistics-based block splitting accumulates **8 literal observation classes + 2 match classes**, with periodic observation checks, to decide when a separate Huffman block might be worth modeling. JLZ tried **fixed** entropy splits (math metadata expanded) and **variable BWT change points** (ratio improved but CPU worse). A block-splitting observation method is **mechanistically distinct** from fixed 2/4 segmentation but still adds model metadata and work; **LOW priority**, not permission to restart segment-count sweeps.
- Zstd superblock estimation and shared entropy tables across emitted sub-blocks are not an untouched magic solution: JLZ previously screened reusing selected finite-state tables across blocks and found only ~5,315 B component potential while compromising block independence. **CLOSE simple table-reuse variants**.

### Opportunity E — LZP before BWT + QLFC (libbsc), genuine concept but weak current Pareto case

libbsc's independent `lzp.cpp` matches next bytes predicted by prior short context and its QLFC coder provides an alternative to MTF. Its public API permits LZP preprocessing with BWT or sort transform and QLFC. JLZ currently selects **raw original-source BWT** versus math/predictive paths; it does not stack an LZP residual before BWT in its accepted source. That is conceptually different from increasing a BWT window or changing MTF context buckets. **However:** previous nonzero-rank BWT changes, cumulative-frequency ordering, context-dependent MTF and high-cost transformed BWT trials often lost CPU, and two extra reversible passes would likely worsen encoder AND decoder timing. Therefore **RESEARCH QUEUE ONLY**, low priority until a zero-extra-pass integration is conceivable. Do not install/call libbsc as JLZ backend or translate its C source.

### Other source-derived findings that are ALREADY in JLZ

- LZ4HC finds backward match extension; `engine/match.go` already extends legal references backward through pending literals with exact/XOR/ADD inverse expectations.
- LZ4/Snappy bulk copying of match runs: JLZ `match/expand.go` already handles nonoverlapping exact copies with Go `copy` and replicates overlapping copies in expanding chunks. Transformed references use forward recursion necessarily.
- LZFSE/Zstd separate literal, copy length and distance coding; JLZ already splits native command streams and uses repeat-distance cache.
- Zstd/LZMA high-quality cost-aware optimal parse paths: JLZ **already tested** a source-frequency-informed one-step cost variant that reduced archive 150 KB but caused serious encoding CPU regressions; **CLOSED concrete implementation**, not an invitation to repeat parser lookahead or histogram-price cutoffs.
- XZ BCJ executable branch normalization: original reversible integrated JLZ candidate reduced archive 154 KB but slowed both speeds and added mode18; **REJECTED**, do not sweep density gates.
- Zstd long-distance matching and dictionary length: JLZ's 128→512 KiB conventional history and sparse far-anchor attempts were tested; both concrete variants had excessive CPU or incomplete net-archive evidence.

## Priority and strict test sequence (without pretending a candidate exists)

**Rank 1, encoder match-verification equality:** before a format change, implement a disposable Go isolated function for XOR/ADD **word-wide transformed equality**; on selected regions from the original full 12-file Silesia, compare every resulting candidate length, token byte stream, and native entropy archive **exactly** with current. If work grows or candidate outputs change, stop. If it substantially reduces measured native tokenizer time, only then integrate and run **ten candidate-only complete original 12-file** runs. No candidate exists yet.

**Rank 2, avoid losing second entropy pass:** instrument existing native packed-stream winners and candidates; do not build a new coder. Develop a strict safe no-win test for the existing FSE coder and measure whether it saves work on current selected source streams. If a heuristic changes any complete archive bytes unexpectedly, it must still separately pass all three whole-corpus criteria and exactness; preferred implementation maintains identical archive output.

**Rank 3, hash latency equivalent scheduling** only if profiles show match index is dominant after simpler improvements.

**Not now:** alternate LZ parsers, extra entropy contexts, BWT window doubling, arbitrary LZP→BWT stack, range-code frequency header redesign or match-window threshold sweeps. Explicit nonregression on size, encode and decode remains mandatory for any actual production promotion.

## Experimental status and measurement honesty

This audit inspected **actual competitor source** and **actual live JLZ source**, and identified precise functions and novel differences. It **did not implement or benchmark an integrated JLZ candidate**, did not perform Silesia encode/decode performance runs and **did not claim any saved bytes or seconds**. Future diagnostic component measurements are not a ten-run integrated benchmark. Keep current unchanged; no empty/duplicative four-row comparison table in user-facing research reporting.
