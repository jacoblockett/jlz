# Canonical research index

**Approved baseline: Round 21.** The latest accepted full-corpus benchmark is `docs/reports/2026-10-18-round21.md` and four adjacent CSV audit/measurement files. Round 16 remains documented in `docs/reports/2026-10-18-round16.md`; rejected rounds 17–20 are canonicalized in their dated reports. Historical figures vary by host; compare only controlled paired runs.

## Promoted work (historical sequence)

| Round | Accepted change | Evidence / caveat |
|---|---|---|
| Initial | Native-only JLZ1, no levels; separate from DEFLATE-backed MDC | Original refactor preserved MDN1 output; native baseline weak |
| Scoring | More accurate benefit assigned to short exact matches | 96.16 MB to 70.46 MB Silesia; large independent native win |
| 3 | Selective XOR/ADD gate, repeat-distance probing | 69.46 MB and faster compression |
| 4 | Grouped history within independent decode groups | 68.54 MB; encode faster but decode 4.05% slower in that test |
| 6 | Better match fingerprints/scoring and concurrency | 68.37 MB; speed gains validated |
| 7 | Logarithmic distance slots | 67.12 MB; encode/decode effectively unchanged |
| 9 | Skip unlikely redundant arithmetic candidate coding | Archive size unchanged; encode 5.79% faster |
| 10 | Backward match extension and scoring simplification | 66.80 MB; encode 3.90% faster |
| 12 | Early regional-selection sampling | 66.80 MB unchanged; encode 5.13% faster |
| 15 | Avoid redundant prediction and Huffman passes | 66.80 MB unchanged, 5.30% faster; approved |
| **16** | **Operation-conditioned literal-run/copy-length entropy streams** | **66,519,205 B Silesia, 0.42% smaller; 4.21% faster encode, 3.31% faster decode, ten-run paired medians. Decoder now accepts new block modes; old decoder cannot read those new blocks.** |
| **21** | **Sample-gated native BWT/MTF/zero-run/Huffman on original regions** | **63,770,766 B Silesia (4.13% smaller than Round 16); ten-pair medians 13.08% faster encode and 14.11% faster decode. No size regressions. 27/27 verified, six new BWT archives require new decoder; per-file worst encode slowdown: `osdb` 13.9%.** |

Some releases not numbered above were performance checks without promotion. Exact historical run methods vary; compare only matched paired runs when citing improvements.

## Rejected candidates / revisitation warnings

| Round | Hypothesis / result | Verdict |
|---|---|---|
| 5 | rANS saved only 0.23% bytes; preflight/Bloom alternatives had no substantial positive end-to-end result | Reject |
| 8 | Context-2 literal prediction saved ~6 KB at extra CPU; compact Huffman saved 0.25% but <1% speed; longer history costly | Reject |
| 11 | Eight-candidate deep matching saved 589,947 bytes (0.88%) but 11.01% slower | Reject |
| 13 | Sparse long-range grouped history saved 418 KB (0.63%) but 5.75% slower; hash prefilter also slower | Reject |
| 14 | Byte-pair grammar recovered ~82 KB (0.12%) but ~10% slower in screening; altered arithmetic format broke one test. Exact scan 5.63% slower | Reject |
| 15 (other) | Arithmetic-first skip ~2.62% faster, below threshold; combined candidate didn't justify extra complexity | Reject |
| 17 | Raw-source entropy-adaptive match scoring: 46,666 B smaller (0.070%) but much slower in one full-corpus screen; small-count arithmetic probability lookup: identical archives without speed gain. Neither passed screening; **not** ten-run medians. | Reject |
| 18 | Two-context opcode Huffman: **48,772 B smaller (0.0733%)** on Silesia but needs decoder extension. Final-archive audit shows 28.19 MB of selected literal payload and 15.08 MB of distance metadata; fixed-position distance bit model theoretical savings only 53–65 KB. Block-shared arithmetic test exceeded time budget. All results screening/diagnostics only. | Reject |
| 19 | Reconstructed-output delta/XOR literal predictors worse by megabytes in first-order entropy estimate; online order-two literal model +611 KB vs order one; distance-width differences +561 KB packed; exact reciprocal arithmetic math preserved archives but slowed single full-corpus screening. All diagnostic/screening, not 10-run medians. | Reject |
| 20 | Selected literal bit-mixer tuning estimated -75,030 B on 9.62 MB sample; previous-byte group model estimated -55,120 B on separate 5.30 MB sample; literal-run phase gave 943,986 B **in-sample static idealized** savings across 39.56 MB. No new encoded archives, no 10-run times, no promotion. | Reject |

See `docs/research/attempted-ideas.md` for the extensive **actual historical prototype inventory** covering algebra, parity, hyperedges, geometry, number systems, recurrences, model sharing and latent masks. Synthetic 97–99% reductions and transform-plus-Zstd results are *not* evidence of standalone natural-file compression. That inventory is a map of prior implementations and limits, not a theorem blacklist.

## Round 22 gate screening (not promoted)

The [Round 22 gate screen](reports/2026-10-18-round22-gate-screen.md) tested mathematically safe early rejection of BWT selection whenever two sample wins become impossible or the partial BWT size already exceeds the total 5% margin. Exhaustive 262,144-case synthetic length enumeration had zero selection disagreements and avoided 103,104 of 786,432 potential BWT sample encodes (13.1% on that synthetic grid only). Local suffix-sort microbenchmarks were approximately 12–13 ms per 256 KiB real-world sample. **No end-to-end JLZ or Silesia benchmarks were run, no source changes were accepted, and Round 21 is still the baseline.** Next: measure per-file gate exits and ten-paired whole-corpus impact in a disposable checkout before promotion.

## Next investigations

Round 21 demonstrates a substantial structural alternative to repeated match-parser and entropy-tuning tweaks. The accepted encoder uses **three already identified representative regions** for BWT vs mathematical-region sizing, gates on **two of three wins and at least a 5% total margin**, then encodes BWT blocks in parallel without running the full mathematical matcher on those files. The approach is **pure Go/native only**, not bzip2 or any other mature compressor backend. The algorithms BWT, MTF and Huffman have extensive prior art. See `docs/reports/2026-10-18-round21.md` for actual complete-archive numbers, all per-file timing medians, full compatibility, independent Zstd/XZ comparisons, and the important Go suffix-array index serialization dependency.

**Next priority:** test performance and byte-exact safety on diverse legally reusable corpora, identify a *generalizable* cheaper BWT gate that reduces losses on `osdb` and other files (without hard-coded file-type rules), and inspect the implementation-specific suffixarray index enumeration only when needed. Do not compromise YAGNI, modularity, automatic worker limits, backward decoding, or main-only development. Preserve rejected Round 17–20 ideas; only revisit with a materially different cost argument.
