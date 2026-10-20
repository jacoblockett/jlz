# Source-syntax partition before native BWT: webster screen

**2026-10-20. Status: SCREENED and REJECTED for this placement.** Current accepted native engine and PGO unchanged; no integrated candidate or full-corpus ten-run benchmark was warranted. This is a real native reversible pretransform performed on source blocks of the **original full Silesia `webster` file**, not a theoretical rate calculation.

## Source diagnosis and placement

The original `webster` contains frequent `<...>` markup mixed into otherwise highly repetitive dictionary text; the accepted current pipeline chooses native BWT/MTF/FSE for its 256KiB blocks, archive **9,835,416 B** over the complete 41,458,703-byte source. Previous attempts at a short lexical word dictionary *before BWT* expanded selected BWT block outputs; this experiment is structurally different: **separate syntax classes into two existing BWT encodings without a positional mask**, not a per-word replacement or a contiguous two/four-part entropy segment.

**Reversible native Go prototype:** parse the original source deterministically with a two-state grammar. Outside markup, append bytes to the text stream including an opening `<`; after seeing `<`, append subsequent bytes to tag stream through the first closing `>`; then resume normal text, including cases of unfinished markup at EOF. The decoder reconstructs the original byte ordering using only already decoded source delimiters and the original output length; no stored bitmap or positions. Encode text and tag streams separately with the **existing `encodeBWTSubregions`** native path. Charge the two selected BWT mode bytes, the recovered text original size, two BWT payload sizes, and a new option byte; check both native BWT inverses and exact source reintegration. This experiment targets input-domain markup mixing rather than a new entropy model.

## Three natural original webster source blocks

| Block position | Actual full block source B | Markup source B | Existing BWT B | Split text+tag B including framing | Delta B |
|---|---:|---:|---:|---:|---:|
| First (0) | 262,144 | 59,677 | 64,496 | 66,101 | +1,605 |
| Middle (79) | 262,144 | 63,810 | 65,937 | 67,543 | +1,606 |
| Last (158) | 39,951 | 6,581 | 13,895 | 14,059 | +164 |
| **Totals** | **564,239** | **130,068** | **144,328** | **147,703** | **+3,375** |

All **3/3 complete natural sample blocks** restored byte-for-byte after both BWT codec inverses and deterministic stream reintegration. The BWT path selected its established mode15 for all measured components. The split representation expanded **2.339%** relative to incumbent BWT component output. No integration or ten-run candidate follows a clear source-stage encoded-size loss of this magnitude plus an extra BWT model and decode workload. No source code or new format committed.

## Decision

**REJECT this concrete markup/text separation before independently coded BWT for `webster`.** The accepted BWT likely already groups much of this syntax context without paying to describe it twice; the trial establishes expansion, not a proof that markup-aware methods cannot succeed on other substreams. Stop nearby alternate delimiter/tag boundary rule tuning; move to different source-aware redundancy in math-selected `mozilla`, or a genuinely new low-cost structural model. Underlying scientific grammar concepts remain OPEN for other distinct placements.
