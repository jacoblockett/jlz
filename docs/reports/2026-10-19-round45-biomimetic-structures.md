# Round 45 — Sponge/skin/DNA structural hypotheses and reversible component screens

**Date:** 2026-10-19. **Decision: REJECT H1 (two finite standalone/component compression candidates); retain biomimetic structural research queue OPEN.** Approved Round 44 native JLZ encoder/decoder is unchanged. No GitHub Actions, branch or production codec edits. The user specifically requested adding mathematically motivated, currently unimplemented research (sponge mechanics, skin recoil, DNA structure), in parallel with reducing winning-path compute.

## Scientific mechanism vs compression hypothesis

**Bath sponges:** many marine bath-spongin skeletons contain an open 3D interconnected network of collagen-like spongin fibers; connectivity, flexibility and water-mediated mechanics support structural recovery after compressive loading. Their mechanical state is constrained by **graph/topological scaffold, local stiffness and energy minimization**. DNA is involved in biosynthesis and organismal development, but is not a runtime byte-for-byte position store used during passive sponge recoil. Reviews: [Marine Spongin: Naturally Prefabricated 3D Scaffold-Based Biomaterial](https://pmc.ncbi.nlm.nih.gov/articles/PMC5867632/); [3D Skeletal Scaffolds of Marine Keratosan Demosponges](https://pmc.ncbi.nlm.nih.gov/articles/PMC12938604/).

**Skin:** elastin-rich fibers provide large, near-reversible elastic recoil, collagen limits overstretch and provides tensile strength, extracellular-matrix crosslinks and water shape the mechanics. Restorative force reflects a mix of molecular conformational **entropic elasticity**, energy and mechanical constraints, not DNA rereading a complete desired surface mesh after every stretch. Review: [Mechanical Properties and Functions of Elastin](https://pmc.ncbi.nlm.nih.gov/articles/PMC10046833/) and [Clinical Relevance of Elastin in Skin](https://pmc.ncbi.nlm.nih.gov/articles/PMC8239663/).

**DNA:** complementary base pairing (A/T, C/G) reconstructs the complementary strand deterministically; gene expression builds proteins and regulatory structures, not direct coordinates of every cell. Repeated sequence families, substitutions, insertions/deletions and strand inversions inspire **template+edits** coding and reversible orientation/complement operators. NHGRI: [DNA fact sheet](https://www.genome.gov/about-genomics/fact-sheets/Deoxyribonucleic-Acid-Fact-Sheet).

**Compression translation:** represent original bytes `X = F(scaffold,parameters) XOR residual` (or modular-add residual) and charge the entire scaffold/model/parameter/residual/selection/framing cost. Recovery must be deterministic and *exact*, not merely convergence toward a similar shape. This cannot violate entropy or information-theoretic bounds. Domain-specific useful structure must be demonstrated, not inferred from natural elasticity.

## A. Reversible elastic-lattice predictor on original Silesia (finite test REJECTED)

Implemented a small independent Go source-to-residual/decoder diagnostic. For each complete 256 KiB block with horizontal grid width `w` in `{16,64,256,1024,2048}`, predict each input byte as the modular additive sum of previous horizontal and previous vertical byte minus the previous diagonal byte:

`prediction[i] = (x[i-1]+x[i-w]-x[i-w-1]) mod 256`; for first row use prior horizontal byte, for row beginning use vertical byte, first byte predicts zero. The residual is `(x[i]-prediction[i]) mod 256`. The decoder reconstructs sequentially and verifies equality of *all original bytes*; the identity follows since prediction depends only on already restored earlier positions. The scheme is a 2D cross-difference predictor, not a claim of a newly discovered physics theorem. The source file data are treated as candidate 2D lattice coordinates, not as actual images.

For every block, independently compute empirical Shannon entropy of each of the five residual variants, then choose the minimum. The sum of blockwise zero-order entropy bounds is an **optimistic minimum payload bound** for a per-block zero-order/Huffman representation: **179,014,479 B (84.47% of original), before codebook, mode, stride and metadata**. The approved Round44 JLZ entire archives use only **61,418,211 B**. This finite scheme cannot compete as a standalone native *zero-order packed* replacement for JLZ. All 12 files and all five widths were checked for exact original-byte inversion. This bound **does not rule out** other entropy coders, layouts or complementary residual operations; those remain separate hypotheses.

## B. Elastic-lattice predictor on *actually selected* Round44 BWT streams (finite test REJECTED)

Implemented a temporary Go test using JLZ's actual current native `bwt.Transform`, `bwt.BijectiveZeroEncode` and `packStream`, iterating all original 256 KiB regions of the **six files that truly select BWT**: dickens, mr, nci, reymont, webster, xml. Test completed in 6.52 s (diagnostic wall clock, **not a JLZ benchmark**). For each of the **412** selected full-source regions, compared its *actual existing native-packed BWT symbol stream* against the **most optimistic zero-order entropy lower bound** of the five possible elastic-lattice residual streams. No model, codebook, mode-byte or stride metadata was charged to the residual—a strongly favorable, unrealistic lower bound.

| Selected BWT file | Existing actual native packed stream B | Best lattice-residual entropy bound B | Bound below existing: regions |
|---|---:|---:|---:|
| dickens | 3,294,743 | 4,427,964 | 0/39 |
| mr | 2,856,176 | 3,327,715 | 0/39 |
| nci | 2,187,050 | 3,055,829 | 0/128 |
| reymont | 1,504,496 | 1,985,999 | 0/26 |
| webster | 10,158,516 | 13,947,110 | 0/159 |
| xml | 532,377 | 751,672 | 0/21 |
| **All selected streams** | **20,533,358** | **27,496,289** | **0/412** |

**REJECT this zero-order lattice-residual + native-pack combination**, because even its optimistic Shannon lower bound cannot match the current selected packed BWT stream on any original Silesia region. This is an especially useful *complementary* screen, not a false broad disproof of shape-memory methods. The temporary test file was removed after completion; no experimental production code promoted.

## C. DNA-inspired reverse and bitwise-complement mirrors (reversible component prototype; future integration OPEN)

Implemented independently in Go a bounded scanner and **fully serialized raw-literal+mirrored-copy encoder/decoder** over each of the 12 complete original Silesia files. Uses eight-byte hash anchors, a fixed ~1M-bucket last-seen position table, a 32-byte minimum match, 4096-byte max, requires >=4 distinct symbols in anchor to exclude trivial uniform repeats. Two optional reversible references: `x[i+k] = x[j-k]` (reverse existing bytes) and `x[i+k] = x[j-k] XOR 0xff` (reverse bitwise complement). The latter is **only an arbitrary-byte analog of DNA reverse complement**, not biological A↔T/C↔G mapping. Match start and length are charged as unsigned varints; literal runs and their lengths, mode tag and complete 4-byte format header are serialized; decoder checks positions, bounds, modes and tail and reconstructs every byte from prior output. Greedy selection is limited by hash collisions/one candidate per bucket and can miss many opportunities—**observed coverage is not a global maximum or bound**.

Complete **12/12 Silesia exact reconstructions** and actual standalone serialized cost:

- `mr`: 2,122 reverse spans, **590,048 bytes** covered; standalone archive **9,401,182 B** from 9,970,564 B source (94.29%).
- `samba`: 561 reverse/76 reverse-complement spans, **48,892 bytes** covered; standalone **21,563,354 B** from 21,606,400 B (99.80%).
- `mozilla`: 10 reverse/20 reverse-complement spans, **2,358 bytes** covered; standalone **51,218,429 B** from 51,220,480 B.
- Other nine original Silesia files: no qualifying nontrivial mirrors discovered with this restricted index.
- **All twelve:** 2,789 qualifying reversible reference spans covering **641,298 input bytes**; complete *raw+mirror* serialized bytes **211,324,218 B** vs uncompressed **211,938,580 B**.

This **does not compete** with the existing **61,418,211 B full JLZ Round44 archive** as a standalone replacement, and no pairwise compressed-token integration or complete JLZ candidate was constructed. The positive observation on `mr` is worth preserving as **component-only evidence** that orientation-inverted repetition does occur in one selected BWT file. It does **not establish incremental savings** over existing BWT, match, entropy or cross-block streams. Do not promote a new copy opcode without a measured selected-stream gain.

## Required full-corpus reference table — no integrated JLZ candidate in this round

| Model | Complete compressed B | Reduction from original | Compression time | Decompression time |
|---|---:|---:|---:|---:|
| Zstd-3 (**historical CLI**) | 66,216,569 | 68.76% | 1.2561 s | 0.8551 s |
| XZ-6 (**historical CLI**) | 49,233,340 | 76.77% | 116.9892 s | 2.3079 s |
| **Current JLZ Round44**, retained candidate-only 10-run median | **61,418,211** | **71.02075%** | **5.945576 s** | **2.445270 s** |
| **Candidate JLZ Round45** | **N/A — not integrated** | **N/A** | **N/A** | **N/A** |

| Candidate JLZ compressed-size delta | Candidate JLZ encode-time delta | Candidate JLZ decode-time delta |
|---:|---:|---:|
| **N/A** | **N/A** | **N/A** |

**No candidate JLZ full-archive ten-run promotion was warranted or attempted.** The 211,324,218 B result is a complete raw+mirror *standalone prototype*, not a competitive integrated JLZ compressor and must never be presented as a candidate JLZ win. Both other numbers are lower-bound component screens, **not complete-archive compression results**. Historical Zstd/XZ CLI references were not rerun and are not directly comparable to retained JLZ in-memory results.

## Validation and next experimental targets

The original 12 Silesia SHA-256 identities were independently checked against the canonical manifest; all passed. The stand-alone elastic-lattice transform and raw+mirror prototype restored every complete input byte exactly. After removal of the disposable selected-stream Go test file, the unchanged approved Round44 source passed `go test -count=1 ./...`, `go vet ./...` and `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build`. No additional race suite needed for these unpromoted single-threaded diagnostics; no new production goroutines or persistent tests. No production decoder mode, source, or archive format was added. No GitHub Actions.

**Decision: REJECT the two specific finite elastic-grid applications; keep the broader biological-math and mirror/edits hypothesis OPEN.** Next, only after cheap cost analysis, evaluate orientation-inverted references **inside actually selected `mr` BWT or math intermediate streams**, with explicit token/frame/header costs and speed budget. Separately, a bounded **template + insertion/deletion edit script** (DNA-like mutational motifs / skin rest-shape deviations) is currently unimplemented and materially distinct from JLZ's exact/XOR/ADD reference matches; start with component viability on existing winning streams. Avoid re-running simplistic byte second-differences, generic Huffman on raw residual, or duplicate sparse-defect experiments already rejected in R19/R29/R33.
