# Theoretical research backlog: ACTIVE

**Directive (2026-10-18):** The user's mathematical/scientific ideas catalog is a **source of experiments to perform when the current optimization route stops producing meaningful improvements**. It is not just historical reading, and past failures are not blanket disqualifications. Distinguish the **idea or theorem** from each tested **finite implementation**. Try new implementations when materially different and log the outcome. The retrospective prototype record `attempted-ideas.md` is for deduplication and interpretation, **not to veto experiments**.

The current practical threshold is holistic archive size, encode/decode times, memory, exact original-byte restoration and YAGNI. There is no hard numerical improvement minimum. **Only original, complete Silesia (12 files, 211,938,580 B)** is authorized for performance benchmarking. Reuse the recorded current JLZ and Zstd/XZ reference measurements, rather than rerunning them. Cheap samples are allowed as preliminary screening, **not** to claim full-corpus performance.

## Active queue

| Idea family | Investigation status | Next *concrete* implementable test |
|---|---|---|
| Prime-index chunk and irreducible-prime factor strategies | **PARTIAL SCREEN 2026-10-18:** prime-position Huffman segregation failed badly against Round 25 archive totals; **full family OPEN** | Try a structural prime-factor/base description or related recurrence against selected original bytes; not just prime vs nonprime rearrangement |
| Higher-order partial information decomposition | **PARTIAL SCREEN 2026-10-18:** held-out pair-XOR diagnostic and reversible lag-XOR residual test; **full PID OPEN** | Find nontrivial multivariable *conditional* structure and evaluate held-out adaptive coding cost, not merely static XOR |
| Polynomial Freiman–Ruzsa / additive combinatorics | **Round 29 accepted narrow rank-one modular codec. Round 30 rank-two screens rejected:** implemented 2D mixed differences, explicit two-generator affine grid plus exact residuals, and two-generator additive walk with exact exceptions; all 36 Silesia sample blocks roundtripped, but no rank-two packed sample beat both native BWT and rank-one. Exact bounded additive doubling showed some structured regions (e.g. x-ray middle 3.41), but not a new compressed-archive win. **Full PFR theory OPEN.** | Revisit only if the next approach exploits low-doubling sets through an independent reconstructible model with fully charged metadata; meanwhile progress to CAT(0)-inspired boundary-distance reconstruction |
| CAT(0) geometric boundary-distance reconstruction | **ROUND 31 FINITE PROTOTYPES REJECTED:** wrote and ran fully reversible 2D bitplane quadtree square labels, horizontal run boundary distances, and vertically integrated edge-incidence boundaries. All 36 Silesia regions restored exactly; optimistic geometric packed 6,404,465 B vs best existing native BWT/additive 2,676,525 B, 0/36 wins. **General CAT(0) theory remains OPEN**, including nontrivial geometric embeddings/inverse boundary methods. | Move to distinct Tietze/presentation grammar transformation next; reconsider geometry only with a materially different exploitable structure |
| Whitehead/Nielsen transformations | **SIMPLIFIED historical trials only; full OPEN** | Encode candidate free-group word mappings with fully charged inverse metadata and native entropy coding |
| Tietze-style grammar/presentation minimization | **ROUND 32: three fully reversible finite Go prototypes rejected.** Nested two-symbol generator rules with charged definitions; block-sorted grammar tokens; rule-introduction/elimination using best complete serialized presentation among 8 checkpoints. On 36 original Silesia regions, best grammar 3,234,359 B vs best native BWT/additive 2,676,525 B, 0/36 wins. **General Tietze theory remains OPEN.** | Next independent active theory: sparse GF(2) syzygy with fully encoded basis/coefficient/exception costs; return to grammar only with a genuinely different presentation architecture |
| Syzygy/Gröbner-like dependency minimization | **SIMPLIFIED historical GF(2) basis screens only; full OPEN** | Find nontrivial sparse dependencies with complete recovery coefficients and bit-cost checks |
| Prime-index irreducible basis vs shortest description | **OPEN** | Bound candidate search and charge generator/program bytes, no manufactured synthetic success |
| General shortest-program/algorithmic complexity approximations | **UNTESTED FULL SEARCH; scoped practical variants OPEN** | Bound grammar/recurrence program description and decode resources; avoid claiming undecidable ideal optimal coding |
| Neural Mamba/RWKV/transformer predictive coding | **NOT ATTEMPTED as a full standalone codec** | Only consider if learned model weight cost, exact bitstream coding and CPU/memory are economically credible |
| Quantum computation | **NOT ATTEMPTED; requires actual quantum infrastructure to test genuinely** | Avoid mislabeling classical simulation as a quantum implementation; if feasible, plan a testable device-backed hypothesis |

More catalog families may be appended as they are identified. No finite negative result closes the parent theory.

## Concrete experiments already recorded

| Date | Candidate | Evidence | Decision |
|---|---|---|---|
| 2026-10-18 | Prime-index positions split into two native pack streams, exact deterministic reconstruction | Three 256 KiB Silesia samples per file saved 7,380 B over ordinary native packing; full 12-file **prime-only archive *estimate*** 133,180,437 B versus 63,104,117 B approved JLZ | **REJECT this representation**, full prime strategy family remains open |
| 2026-10-18 | Pair-lag XOR parity as a limited higher-order-information proxy | `x-ray` holdout entropy difference favored two-lag parity, but after actual XOR residual/native packing three 256 KiB `x-ray` blocks were 778,895 B versus native BWT 429,274 B; reversible exact on all 12 samples | **REJECT this XOR residual**, full PID family remains open |

Detailed methodology and qualifications: [Round 28 research](../reports/2026-10-18-round28-theory-screens.md). Earlier isolated screens that did not implement the complete theory must never be cited as the theory being definitively disproven.

| 2026-10-18 | Rank-one additive progression in Z/(2^16), signed-zigzag and two JLZ native entropy-packed word coordinates, automatically selected on original Silesia | Ten candidate-only full-Silesia medians 62,644,663 B and 7.516364/2.821363 s versus retained Round 25 63,104,117 B and 7.549920/2.839905 s; 27 prior-archive and new-archive validations | **ACCEPT concrete rank-one coder as Round 29**; full PFR / higher-rank family remains OPEN |


| 2026-10-18 | Round 30: three complete *component* Go encoder/decoder prototypes for rank-two additive grids/walks; bounded exact 16-bit group doubling diagnostic | Across 36 original Silesia regions (7,728,404 B), optimistic best rank-two packed payload 5,716,493 B versus 2,676,525 B best existing rank-one/BWT choice; no sample wins. Tests, vet and targeted race passed. No full JLZ archive candidate / 10-run metrics. | **REJECT these three models; keep PFR open**; move to genuinely different CAT(0)-inspired geometric encoding |

| 2026-10-18 | Round 31: 2D CAT(0)-inspired exact bitplane cubical quadtrees, horizontal boundary-distance rows, vertical edge-incidence contours | 36 representative Silesia regions, 7,728,404 original B: best actually packed geometry 6,404,465 B, best existing BWT/additive 2,676,525 B; **0/36 geometric sample wins**, all reconstructed exactly. Tests/vet/race passed; no integrated JLZ candidate. | **REJECT these specific cubical boundary encodings**, general CAT(0) open |

| 2026-10-18 | Round 32: executed nested acyclic Tietze-style two-symbol grammar (24/64 rules), sorted grammar token streams, and explicit generator elimination/selection across 8 encoded presentation checkpoints | 36 original Silesia regions 7,728,404 B; smallest grammar payload 3,234,359 B vs incumbent BWT/additive payload 2,676,525 B, zero improvements. Byte-exact Go decoder, tests/vet/race passed. No complete JLZ archive candidate. | **REJECT three finite grammar representations; full Tietze theory OPEN** |

## Research workflow for next sessions

1. Select the highest-value *untested* catalog hypothesis once codec optimization is diminishing, consulting the historic inventory only to avoid **identical** previous implementations.
2. Devise a bounded and genuinely **reconstructible** encoding hypothesis, preferably with a simple held-out Silesia viability probe.
3. If the result has realistic scope to beat approved archive bytes, implement a standalone native JLZ candidate in a disposable local copy; verify exact original bytes.
4. Run ten candidate-only complete Silesia encode/decode benchmarks, reuse published approved/reference results, and report the mandated four model rows plus three percentage deltas and colored H1 accept/reject recommendation. Disclose independent-run timing uncertainty.
5. Record **all** measured attempts, successful or not, under `docs/reports/`, then update this backlog, `docs/STATE.md` and `docs/RESEARCH.md`. Promote code only when justified and only on `main`.
