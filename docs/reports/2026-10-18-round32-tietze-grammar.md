# Round 32 — Reversible hierarchical grammar presentations (rejected finite experiments)

**Date:** 2026-10-18. **Approved JLZ unchanged:** Round 29. **Decision:** REJECT the three implemented local Tietze-style straight-line-grammar variants on tested Silesia samples, **not** the general algebraic/Tietze theory. No production code or new JLZ1 block mode was promoted; code/tests ran in an isolated disposable Go source copy.

## Scope and implemented mechanisms

Only the **12 original Silesia files** were used for compression research. The experiment reads first/middle/last 256 KiB (or shorter final) original-byte blocks from each: **36 blocks totaling 7,728,404 source bytes**. Original Silesia hashes were checked against the previously recorded archive manifest. Benchmarking did **not** rerun accepted Round 29 or Zstd/XZ. Existing native BWT and additive block encoders were invoked solely as **matched sampled-payload comparators**, not full archive or timed baseline benchmarks.

The idea is to replace recurring data substrings by short generator references and serialize each generator's exact defining relation. This is a practical **acyclic straight-line grammar analogous to limited Tietze introduction/elimination of generators**, not the full general Tietze theorem/presentation isomorphism procedure.

**Prototype A — Hierarchical two-symbol substitution:** Starting with terminals 0–255, count adjacent ordered token pairs; introduce a new nonterminal `G_i := (a_i,b_i)` for the most frequent profitable pair, replacing nonoverlapping occurrences and allowing later rules to refer to earlier generators. This creates nested rules, not just a flat byte dictionary. The on-disk payload explicitly charges original length, number of rules, token count, **four bytes for each rule pair**, and two independently Huffman/RLE/raw-packed token coordinate streams (low/high byte of a 16-bit terminal/nonterminal code). Two limits (24 and 64 generator definitions) were evaluated.

**Prototype B — Grammar-token block sorting:** Reuse the identical reversible generator presentation but encode the token low-byte stream using JLZ's native **BWT→MTF→zero-run→Huffman**, plus a separately native-packed high-byte stream. The BWT length and all generator definitions are explicitly framed. The prototype decoder reconstructs the sorted token stream, validates rule IDs and exact expansion lengths, then reproduces all original bytes.

**Prototype C — Presentation elimination with charged serialized cost:** At generator counts 0, 4, 8, 16, 24, 32, 48, and 64, preserve the intermediate actual grammar presentation, encode its complete payload as above, and select the cheapest. This is an **optimistic per-block oracle** over 8 fully specified grammar versions, not a low-cost automatic production selector. All chosen alternatives include their own rules and stream framing. Every tested checkpoint could be decoded back to the byte-exact original block.

The inverse uses a bounded iterative expansion stack. It rejects references to unknown/later generators (acyclicity), token-length mismatches, excess expansion, malformed headers, and trailing bytes. Generated dictionaries are not assumed to exist in the decoder.

## Actual packed Silesia sample results

| Sampled representation, 36 regions | Bytes |
|---|---:|
| Original Silesia bytes sampled | 7,728,404 |
| Grammar A, 24 nested generator definitions | 4,976,978 |
| Grammar A, 64 nested generator definitions | 4,501,188 |
| Grammar B, 64 generators and native block-sorted token low bytes | 3,236,550 |
| **Grammar C, best intermediate presentation per original region** | **3,234,359** |
| **Existing BWT/additive, best current native representation per same region** | **2,676,525** |

**No grammar variant won on any of 36 sampled regions (0/36).** Even granting cost-free selection of its best of eight presentations separately in every block, grammar was **557,834 bytes (20.84%) larger** than the incumbent's favorable per-block choice over the same Silesia source windows. This is a source-region component comparison only, not a complete JLZ1 archive result. More generator definitions generally helped, but the improvement was insufficient. Do not extrapolate the numbers to an accepted 12-file compressor ratio or time.

The original grammar alone (A) had a much larger deficit; block-sorting the grammar token stream cut it substantially, demonstrating the effect of combining exact symbolic substitution with existing native BWT. It nevertheless did not cross the incumbent archive-size threshold on tested blocks. Any new JLZ1 block headers, selector work, or compatibility costs would have made this proposal less attractive.

## Tests and accuracy

- **Go implementations were actually written and run**, including both encoding and decoding functions, not Shannon bounds or unimplemented pseudocode.
- Initial hierarchical grammar screen encoded 36 real Silesia blocks, byte-exactly decoding the 64-rule presentations (full-screen test passed).
- Native block-sorted grammar screen encoded/decoded the same 36 blocks with exact equality (passed).
- Eight-point rule introduction/elimination checkpoints all passed exact decode and expansion validation wherever that checkpoint had enough rules (passed).
- Independent short, repeated, random-like, empty, and malformed-length fixtures passed as correctness-only controls.
- `go test -count=1 ./...`, `go vet ./...`, and targeted `go test -race -count=1 -run '^TestTietzeGrammar(Correctness|SortedCorrectness)$' ./...` **passed** locally. Silesia-specific loops also independently completed.
- No new complete 12-file JLZ archive, no new archive mode, **no candidate-only 10-run time medians**, and therefore no candidate-vs-current full-archive percentage delta claim. Production Round 29 source was not modified.

## Unchanged references (reused, not rerun)

Approved Round 29 Silesia: **62,644,663 B**, prior candidate-only ten-run median **7.516364 s encoding / 2.821363 s decoding** on 211,938,580 original bytes. Historical Zstd-3 archive **66,216,569 B**, CLI encode **1.2561 s**, decode **0.8551 s**. Historical XZ-6 archive **49,233,340 B**, CLI encode **116.9892 s**, decode **2.3079 s**. Reference CLI and JLZ in-memory timing methodologies differ, and these reference runs were not repeated.

## Decision and theoretical scope

**REJECT these concrete hierarchical straight-line grammar presentations for promotion; keep Round 29 unchanged.** The models were real, reversible and fully charged for generator definitions, yet their sampled serialized payloads lost even to the optimistic current selection. This does not amount to full Tietze transformations nor refute grammar/presentation compression generally. The theoretical catalog is an active queue, not just historical notes; mark this finite family as tested and proceed to a materially different theory instead of endlessly adjusting token codecs.

**Next recommended theory:** a sparse `GF(2)` syzygy / algebraic-dependency representation with a working exact inverse and full coefficient/model costs, tested against original Silesia windows rather than synthetic-only examples. Other Whitehead/Nielsen, prime-structure and higher-order information-theoretic ideas remain active.
