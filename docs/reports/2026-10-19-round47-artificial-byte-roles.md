# Round 47 — Learned artificial predictive byte roles, isolated and complementary tests

**Date:** 2026-10-19. **Decision: REJECT H1 for the tested role-split coding architecture; preserve predictive-role theory OPEN for materially different approaches.** Current accepted native standalone JLZ is **Round 44**; no production compressor or decoder source changed on main.

## Hypothesis and concrete reversible constructions

User proposed inventing *artificial byte meaning* by binding otherwise unrelated values that play similar roles. Implemented seven fully reversible **Go** experiments using deterministic behavioral groupings, not linguistic/semantic interpretation. In each independently coded block, count transitions `previous byte -> current byte`. Assign each predecessor a **modal successor**; rank modal successors by the number of observed predecessor transitions; choose 3 or 7 frequent successors and group predecessor byte identities with the same modal successor into 4 or 8 shared predictive roles, respectively (remaining predecessors assigned role 0). The map and learned predictors are **serialized and charged**, never assumed known by the decoder. All decoders reconstruct in original byte order with only previously decoded symbols and the model carried in the payload.

Tested:
- **split4/split8**: route each current symbol into a role-conditioned subsequence selected by the *previous decoded* byte; each role stream is independently native raw/RLE/Huffman packed; interleave exactly on decode.
- **xor4/xor8**: encode current byte XOR its role's modal successor, then native-pack exact residuals; invert using decoded previous byte.
- **swap4/swap8**: swap each role's modal successor with symbol 0 before native packing, and undo the permutation during decode; this is a bijection and preserves exact identities.
- **direct256**: unclustered per-previous-byte modal successor XOR control (256 predictors). This tests whether the clustering itself is beneficial.

Each block carries its decoded length, representation variant/group count, an independently native-packed 256-byte role map (except direct256), predictor-table data where needed, and independently native-packed role or residual streams. Initial predecessor = 0. No external compression codec. A dedicated unit test exercised nil, constant, periodic and patterned byte sources for all seven formats and passed.

These are established *context grouping/conditional source coding* concepts combined in a fresh JLZ experiment; no claim of a new information-theoretic law. The chain rule still applies: learned class grouping cannot make information disappear and can be worthwhile only through better coding/reduced model overhead.

## Isolated original-byte streams: all twelve complete Silesia sources

Only canonical original Silesia was used: **12 full input files / 211,938,580 original bytes**, with 12/12 source SHA-256 hashes checked against the Round23 manifest. For controlled modeling, divided each whole file into complete consecutive 256 KiB independently decodable regions (including final partial regions); **no original input bytes were omitted**. Every variant decoded every region exactly. Totals below include every model table/role map, block length and variant header, and nested native stream framing. They **do not include whole-file JLZ1 archive frames/digests**, so they are **component totals, not competing JLZ file archives**. Native control is `packStream` directly on each original source region; current Round44 whole-file architecture is much stronger.

| Isolated region representation | Sum of actually packed region bytes | Relative to unmodified native packStream |
|---|---:|---:|
| Native packStream baseline | 133,318,073 B | — |
| Learned shared predecessor roles, four groups | 122,118,809 B | -8.40% |
| Learned shared predecessor roles, eight groups | **115,781,631 B** | **-13.15%** |
| Role-conditioned XOR, four groups | 136,018,520 B | +2.03% |
| Role-conditioned XOR, eight groups | 134,730,888 B | +1.06% |
| Swap modal successor with zero, four groups | 129,873,014 B | -2.58% |
| Swap modal successor with zero, eight groups | 125,933,677 B | -5.54% |

The per-byte per-previous-context modal XOR control also roundtripped; across complete input files the model frequently improves the native raw streams, but no whole-file integrated advantage is inferred from it. Initial monolithic test process halted after eleven fully completed files; the remaining complete `xml` file was run in a separately bounded process. This does not affect the sum-of-complete-blocks correctness and **none of these are candidate promotional timing medians**.

## Supplemental selected-stream component tests

**Existing selected BWT:** Processed **all 412 complete original-data BWT regions** of Round44's six BWT-selected files (dickens, mr, nci, reymont, webster, xml), building the *actual* native BWT → bijective MTF zero-run symbol stream. For each, compared existing `packStream(symbols)` with fully framed learned-role formats and verified exact recovery of every symbol stream. Sum for existing native-packed streams: **20,533,358 B**. The **four-group split** costs **20,305,435 B**, saving **227,923 B**, with **275 of 412 regions individually smaller**. Eight-group split costs **20,336,204 B** (smaller gain); the best previous-mode full-stream total is not the same as the number of actual archive bytes (Round44 reuses 18 legacy-coded sample blocks). Isolated timing over the 412 streams: direct original packer about **219.7 ms**, four-role encoder about **445.4 ms**, and four-role decoder about **629.2 ms**. These are **component timing measurements**, not complete JLZ times and not paired independent compressor benchmarks.

**Actual selected mathematical metadata:** Encoded and decoded complete JLZ Round44 archives for each of the *four mathematical-region winning files* (mozilla, ooffice, osdb, samba) as **untimed selected-stream source diagnostics**, then parsed each actual block's existing first two native-packed streams (operation codes and a length channel; split-length variants represented as encoded). Tested **676 actual packed opcode/length streams** from completed archives, reconstructing each stream exactly. Their actual packed size was **4,506,404 B** versus the alternate four-group role coding **4,067,196 B**, isolated gain **439,208 B**. Note: this counts only the first two actual streams and retains existing distances, transforms and literals; it does **not** assert that every such stream can be altered without changing block-metadata opcode/state framing. These are measured component results, not new full-compressor gains.

## Complete integrated JLZ codec prototype

Implemented a disposable **native codec-5 packed stream subtype** for selected winning BWT symbol streams and mathematical-region metadata streams. It stores an exact, independently recoverable role-model and role streams inside the existing self-described packed-stream framing; legacy codecs 0–4 remain decoded normally. The new variant is automatically selected **only when the fully serialized candidate is shorter** than the existing native stream, with no user switch. Deliberately did **not** enable it for generic `packStream`: doing so changed old byte-exact packing tests and interfered with a nested varint-distance slot grammar. Targeted integration preserved the original native packer's behavior in all other paths and passed the unchanged full Go suite. Newly generated files containing this subtype would **not be forward-decodable by older R44 decoders**, although the candidate decoder retained old modes.

The first full original Silesia pass saved **723,405 B** versus the correct approved R44 archived 61,418,211 B, with complete native framing, archive digest and exact reconstruction. All 12 files got smaller or stayed the same (sao and x-ray unchanged). Largest absolute savings: mozilla **365,446 B**, webster **123,529 B**, dickens **56,768 B**; no size regressions. The integrated savings are slightly greater than the sum of the two isolated component screens because actual candidate streams, wrappers and region selection can differ.

### Candidate-only ten complete original-Silesia encode/decode runs

Each counted run encoded and decoded **all 12 original complete files**, verified byte-for-byte and emitted the identical total **60,694,806 B**. Inputs were loaded before the individually timed encoder/decoder calls. Four complete processes ran first; despite one batch wrapper timing out *after* the fourth process printed `PASS` and a complete 12-file result, that valid process is included. Six additional completed independent process runs were collected separately. No incomplete source pass or failed decoder is included.

| Candidate run | Complete archived bytes | Compression s | Decompression s |
|---:|---:|---:|---:|
| 1 | 60,694,806 | 6.412626243 | 2.584089378 |
| 2 | 60,694,806 | 6.102809285 | 2.531647644 |
| 3 | 60,694,806 | 6.229154371 | 2.633132398 |
| 4 | 60,694,806 | 5.853324061 | 2.620748965 |
| 5 | 60,694,806 | 6.364786832 | 2.610264349 |
| 6 | 60,694,806 | 6.033311002 | 2.532253614 |
| 7 | 60,694,806 | 6.426532897 | 2.623761202 |
| 8 | 60,694,806 | 5.951888537 | 2.471589770 |
| 9 | 60,694,806 | 6.307113556 | 2.548530149 |
| 10 | 60,694,806 | 6.144051970 | 2.575335996 |
| **Median, whole-corpus totals** | **60,694,806** | **6.186603171** | **2.579712687** |

### Additional integrated configuration screens

Compared narrower integrations as fully reversible complete original 12-file archive screens; neither was promoted to ten runs because its initial time/size tradeoff was inferior.

| Integrated scope, one complete screen | Real JLZ1 archive bytes | Compression s | Decompression s |
|---|---:|---:|---:|
| Mathematical metadata only | **60,939,586** | 6.331357 | 2.585098 |
| BWT intermediate symbols only | **61,173,431** | 6.572139 | 2.659631 |
| Both combined (ten-run candidate median above) | **60,694,806** | **6.186603** | **2.579713** |

## Required comparison with retained historical references

| Model | Complete compressed size | Reduction from 211,938,580 B | Compression time | Decompression time |
|---|---:|---:|---:|---:|
| Zstd-3 (historical external CLI) | 66,216,569 B | 68.76% | 1.2561 s | 0.8551 s |
| XZ-6 (historical external CLI) | 49,233,340 B | 76.77% | 116.9892 s | 2.3079 s |
| **Current approved JLZ Round44 (historical candidate-only ten-run median)** | **61,418,211 B** | **71.02075%** | **5.945576325 s** | **2.445270387 s** |
| **Candidate JLZ Round47 (new candidate-only ten-run median)** | **60,694,806 B** | **71.36208%** | **6.186603171 s** | **2.579712687 s** |

| Candidate archive size delta | Candidate encode-time delta | Candidate decode-time delta |
|---:|---:|---:|
| **-1.17783%** (-723,405 B) | **+4.05389%** nominal slower | **+5.49805%** nominal slower |

The Round44 JLZ reference is an earlier independent session on the same method, not a simultaneous controlled A/B; percent timings are **nominal** and cannot establish a precise causal penalty. Established Zstd/XZ historical CLI times involve distinct I/O/runtime conventions and are not directly comparable to JLZ in-memory measurements. No accepted reference compressor was rerun for this study. Comparison is against correctly reconstructed R44 full archive sizes, not the old Round29 provenance discrepancy. Candidate-only source hashes, archive bytes, mode framing and reconstruction were tested in full.

## Validation and decision

- `sha256sum -c docs/reports/2026-10-18-round23-silesia-sha256.txt`: **12/12 canonical inputs PASS**.
- All **120** candidate complete-file reconstructions byte-exact in ten full-corpus runs; 12-file complete archives with SHA-256, no omitted metadata.
- Dedicated reversible role-codec unit tests across all seven models **PASS**.
- Integrated code: `go test -count=1 ./...`, `go vet ./...`, targeted `go test -race ... TestRoleCodecExact`, Windows/amd64 cross-build **PASS**.
- Unmodified accepted R44 separately: `go test -count=1 ./...`, `go vet ./...`, Windows/amd64 cross-build **PASS** after experiments. No GitHub Actions; no new production decoder type.
- The disposable candidate's new packed subtype needs a **new decoder**; only research code was generated and discarded locally rather than committed as accepted source. No source files promoted on main.

**REJECT H1 for this implementation.** It demonstrates a real selected-stream structural compression opportunity and full-JLZ size gain **1.18%**, but makes encoding and decoding directionally **4% / 5.5% slower**, adds nontrivial table training and decoder work, and introduces new packed format compatibility requirements. Under the user's speed-first/practical-simplicity objective, R44 remains preferred. Preserve the concept as **OPEN**, not disproven: next materially different experiment is a *decoder-synchronous online predictive-state merging* codec with no charged per-block learned role map or per-role independent Huffman tables, conditional on cheap selected-stream evidence. Do not spiral into additional group-count, threshold or mapping-key tuning; avoid treating the raw native-pack improvements as full JLZ performance.
