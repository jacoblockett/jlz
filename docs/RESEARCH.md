# Canonical research index

**Approved baseline:** Round 16. The latest dated report is `docs/reports/2026-10-18-round16.md`; Round 15 remains recorded in `docs/reports/2026-10-18-round15.md`. Earlier figures are historical and should not be treated as simultaneously controlled benchmarks.

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

See `docs/research/attempted-ideas.md` for the extensive **actual historical prototype inventory** covering algebra, parity, hyperedges, geometry, number systems, recurrences, model sharing and latent masks. Synthetic 97–99% reductions and transform-plus-Zstd results are *not* evidence of standalone natural-file compression. That inventory is a map of prior implementations and limits, not a theorem blacklist.

## Next investigations

Round 16 demonstrates useful joint token/operation structure without deeper searching. Round 17 disconfirmed raw input zero-order entropy as a cheap reliable proxy for post-match literal cost, and rejected a 512-count triangular arithmetic probability table due to no measured benefit. Full report: `docs/reports/2026-10-18-round17.md`.

Next investigate a unified, statistically cost-aware match/length/distance model implemented in a single inexpensive pass, avoiding speculative extra stream partitions or costly extra grammar/matching passes. Run cheap rejection screens first, then ten-run medians against immutable approved source. Record both positive and negative findings.
