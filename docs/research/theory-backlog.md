# Theoretical research backlog: ACTIVE

**Directive (2026-10-18):** The user's mathematical/scientific ideas catalog is a **source of experiments to perform when the current optimization route stops producing meaningful improvements**. It is not just historical reading, and past failures are not blanket disqualifications. Distinguish the **idea or theorem** from each tested **finite implementation**. Try new implementations when materially different and log the outcome. The retrospective prototype record `attempted-ideas.md` is for deduplication and interpretation, **not to veto experiments**.

The current practical threshold is holistic archive size, encode/decode times, memory, exact original-byte restoration and YAGNI. There is no hard numerical improvement minimum. **Only original, complete Silesia (12 files, 211,938,580 B)** is authorized for performance benchmarking. Reuse the recorded current JLZ and Zstd/XZ reference measurements, rather than rerunning them. Cheap samples are allowed as preliminary screening, **not** to claim full-corpus performance.

## Active queue

| Idea family | Investigation status | Next *concrete* implementable test |
|---|---|---|
| Prime-index chunk and irreducible-prime factor strategies | **PARTIAL SCREEN 2026-10-18:** prime-position Huffman segregation failed badly against Round 25 archive totals; **full family OPEN** | Try a structural prime-factor/base description or related recurrence against selected original bytes; not just prime vs nonprime rearrangement |
| Higher-order partial information decomposition | **PARTIAL SCREEN 2026-10-18:** held-out pair-XOR diagnostic and reversible lag-XOR residual test; **full PID OPEN** | Find nontrivial multivariable *conditional* structure and evaluate held-out adaptive coding cost, not merely static XOR |
| Polynomial Freiman–Ruzsa / additive combinatorics | **PARTIAL PRACTICAL SUCCESS, Round 29 ACCEPTED:** implemented complete rank-one modular 16-bit additive-word progression encoder/decoder and automatic gate, saving 459,454 B (0.728%) on full Silesia with near-neutral accepted-time comparison. **General PFR and low-doubling higher-rank discovery remain OPEN.** | Build/test a genuine rank-two generalized arithmetic progression / low-doubling set model with exact inverse and fully charged parameters/exceptions on Silesia; do not simply keep tuning the accepted one-dimensional codec |
| CAT(0) geometric boundary-distance reconstruction | **UNTESTED FULL THEORY** | Define an exact and economical binary labeling/boundary representation and cheaply test reconstructed description lengths |
| Whitehead/Nielsen transformations | **SIMPLIFIED historical trials only; full OPEN** | Encode candidate free-group word mappings with fully charged inverse metadata and native entropy coding |
| Tietze-style grammar/presentation minimization | **SIMPLIFIED historical trials only; full OPEN** | Search stronger shared-rule transformations with actual decoding and rule amortization across large Silesia regions |
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


## Research workflow for next sessions

1. Select the highest-value *untested* catalog hypothesis once codec optimization is diminishing, consulting the historic inventory only to avoid **identical** previous implementations.
2. Devise a bounded and genuinely **reconstructible** encoding hypothesis, preferably with a simple held-out Silesia viability probe.
3. If the result has realistic scope to beat approved archive bytes, implement a standalone native JLZ candidate in a disposable local copy; verify exact original bytes.
4. Run ten candidate-only complete Silesia encode/decode benchmarks, reuse published approved/reference results, and report the mandated four model rows plus three percentage deltas and colored H1 accept/reject recommendation. Disclose independent-run timing uncertainty.
5. Record **all** measured attempts, successful or not, under `docs/reports/`, then update this backlog, `docs/STATE.md` and `docs/RESEARCH.md`. Promote code only when justified and only on `main`.
