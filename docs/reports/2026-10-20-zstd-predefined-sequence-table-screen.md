# Zstd predefined entropy tables for selected mathematical sequence streams

**Date:** 2026-10-20. **Status:** screened and rejected as practical new coding mode; the exact predefined-distribution implementation, not Zstd's complete sequence-entropy family. No production changes or fully integrated candidate.

## Why this is distinct

Meta Zstd's `zstd_compress_sequences.c` can select predefined distribution tables for literal length, match length and offset codes, avoiding dynamic table headers at low parse counts. JLZ already has dynamic normalized-state entropy and Huffman/raw/RLE/slot coding. Its hypothetical grouped table sharing showed at most 42,162 B of upper-bound savings and required a new group state. This trial instead uses **zero transmitted model bytes**, a independently reconstructible deterministic native Go rANS distribution for **each individually selected native mathematical stream**. There is no cross-block state.

## Candidate model families and real source test

Four original native normalized-to-4096 frequency presets, all 256 bytes supported, were generated from fixed analytical rules with no external corpus-trained assets: (A) nibble-pair opcode prior favoring literal/exact, (B) inverse-byte-rank heavy-tail varint prior, (C) shallower grouped-rank prior, and (D) low-seven-bit leading/continuation mixture. Tested every true existing mathematical block selected by the accepted original 12-member 211,938,580 B Silesia compressor: **279 mathematical blocks**. Baseline was the **actual chosen stream payload** (raw/RLE/Huffman/rANS/varint-slot), not just another full dynamic-rANS model. Each native fixed-model body counted raw length varint, coded-body length varint, one codec byte, and all state bytes; no serialized model header. The preset identity / modified decoder dispatch *not* integrated and would cost additional bytes. Every winning candidate's payload was independently decoded using its respective preset via the unchanged `finiteStateDecode` and compared byte-exact with the true selected original stream.

| Selected stream | Number available | Best gain over incumbent, **oracle chosen among 4 models** | Best model's winner count |
|---|---:|---:|---:|
| Packed opcodes | 279 | 0 B | 0 |
| Combined lengths | 191 | 1,340 B | 9 |
| Separate literal lengths | 88 | 1,205 B | 8 |
| Separate reference lengths | 88 | 0 B | 0 |
| Distances | 279 | 205 B | 4 |
| Transform masks | 64 | 149 B | 3 |
| Literal bytes | 279 | 0 B | 0 |
| **Optimistic model chosen by stream category** | | **2,899 B** | |

The model-by-model win totals for categories combined lengths were **22/880/1,340/623 B**, separate literal lengths **714/1,128/1,205/978 B**, distances **4/0/20/205 B**, transforms **103/140/149/100 B**; all other categories zero. These are **sum-of-positive oracle estimates**, not a complete frame. Even a free model selector would save only 2,899 B (0.00493% of 58,767,739 B), before transmitting preset IDs or new codec types. Dynamic JLZ stream selection is already effective, and model mismatch usually exceeds header savings.

**Decision:** no plausible justifiable new preset wire format or decoder, no ten-run candidate. Do not optimize four model constants or add dozens of hand-tuned tables. The related fundamentally different structural targets, including cheaper execution of existing normalized entropy and a one-pass already-parsed superblock scheme, remain open when an actual whole-codec work budget justifies them.

Reference **58,767,739 B /4.7873709325 s compression /1.8721679705 s decompression** (historical ten-run), Zstd-3 **66,216,569 B /1.2561 /0.8551 s**, XZ-6 **49,233,340 B /116.9892 /2.3079 s**, not rerun. Trial inverse checks and normal accepted `go test -count=1 ./...` passed. The existing vet/race/Windows cross-build were completed before this experiment against unchanged accepted source.
