# Zstd mathematical-sequence superblock: native shared rANS model screen

**Date:** 2026-10-20. **Disposition:** SOURCE-SELECTED COMPONENT SCREEN, NOT PROMOTED. Current JLZ production codec, PGO and archive modes unchanged.

## Novelty, source, method

Studied Meta Zstd `lib/compress/zstd_compress_superblock.c` and `zstd_compress_sequences.c`. Unlike prior rejected cross-group table reuse and BWT-mode16 shared models (0/91 winners), this experiment shares **one native normalized rANS model per semantically matched mathematical metadata stream across existing four-block independently decodable history groups**.

Used accepted `main` HEAD `8566c8a8b28e3161ed941a98ec6ed4fc69906e47` with the original 12-file 211,938,580 B Silesia corpus, whose SHA-256s matched `docs/reports/2026-10-18-round23-silesia-sha256.txt`. Encoded every original file using current `compressFast`. Parsed *actually selected* mixed-cross payloads (not hypothetical math selections) via existing `unpackStream`: **279 selected mathematical blocks**. Collected real opcode, combined length, separate literal/reference length, distance, transform and literal streams. For each group/stream with at least two nontrivial members, built one `finiteStateModelFor`, serialized the complete model, encoded each member independently via `finiteStateEncode` and **verified every trial body with `finiteStateDecode` against exactly the original column bytes**. Compared real existing encoded-column sizes including their framing against a hypothetical group envelope that counts the model, one group flag, and every per-member raw-length, body-length, codec and byte-renormalized state payload.

| Actual stream category | Groups screened | Group wins | **Optimistic selected-field savings** | Max in one group |
|---|---:|---:|---:|---:|
| Opcode nibbles | 71 | 21 | 1,163 B | 98 B |
| Combined length varints | 67 | 36 | 9,819 B | 652 B |
| Separate literal-length varints | 15 | 5 | 655 B | 315 B |
| Separate reference-length varints | 15 | 4 | 334 B | 168 B |
| **Distance varints** | **71** | **38** | **30,188 B** | **1,739 B** |
| Transform masks | 3 | 1 | 3 B | 3 B |
| Literal bytes | 71 | 0 | 0 B | 0 B |
| **Optimistic sum of individually selected group/field wins** | | | **42,162 B** | |

**42,162 B is an optimistic component-level upper envelope, NOT a fully framed integrated archive improvement.** It amounts to just 0.07174% of the accepted 58,767,739 B archive *before* new group/decoder mode cost, extra group subtype selection and worker scheduling. Unconditionally sharing a table often loses, particularly for literals. Distance sharing does exhibit a distinct model-amortization effect: testing independent per-block native rANS instead of the existing selected `packVarintStream` provides only **5,465 B potential savings on 104 winning blocks**, while rANS would expand other blocks by 160,914 B without selection.

## Decision and scope

**DEFER/REJECT this particular new group-model archive mode** as uncompetitive in savings versus format, decode dependencies and worker complexity. Did not integrate or run ten candidate-only JLZ passes because there is no credible performance/complexity purchase. Do not call the entire mathematical-superblock family exhausted: a genuinely distinct *within one parsed math sequence* partition or cheaper predefined distribution remains open subject to separate upper-bound evidence, not minor shared-group parameter tuning.

Every tested native shared rANS stream restored exactly. The unchanged accepted source also passed `go test -count=1 ./...`, `go vet ./...`, targeted race tests and Linux-host Windows/amd64 cross-build. No new decoder, old-archive format, production source or wire changes. Historical accepted anchor **58,767,739 B / 4.7873709325 s encode / 1.8721679705 s decode**; Zstd-3 **66,216,569 B / 1.2561 / 0.8551 s** and XZ-6 **49,233,340 B / 116.9892 / 2.3079 s** retained without rerun (external CLI timing conventions differ).
