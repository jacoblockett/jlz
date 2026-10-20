# JLZ research index: strategies and evidence, not numbered sessions

**Current production:** `current` on `main`, complete original twelve-file Silesia archive **60,585,717 B**, prior ten-run medians **5.493743735 s encoding / 2.194427893 s decoding**. **Benchmarked/integrated candidate: none active.** Strict three-metric no-regression without express user approval. [Full current state](STATE.md).

## Canonical research records

- **[Strategy ledger](research/strategy-ledger.md): READ BEFORE EACH PROPOSAL.** Concrete accepted/failed/scoped mechanisms, actual measured outcomes, reasons for failure, overlap/novelty test, external architecture avenues and anti-spiral stop rule. Update this ledger whenever an implementation is tried or a family is closed.
- **[Theory backlog](research/theory-backlog.md): ACTIVE.** Catalog of mathematical/scientific structures and executable next hypotheses, including genuinely distinct variants of previously failed ideas. A single rejected implementation is not a theorem disproval.
- **[Attempted ideas](research/attempted-ideas.md): retrospective.** Original research families, implementations, negative results and what was *not* actually tested. Distinguish theoretical inspiration, component probes, and integrated complete-codec findings.
- **[Complementary integration audit](research/complementary-integration-audit.md):** How existing native streams and mathematical transforms can be combined additively; reject replacement on inadequate evidence.
- **[Dated reports](reports/):** Full methods, executed prototypes, exact archive byte counts, benchmark provenance, compatibility, regression details and accepted/rejected decisions. Historical reports may contain numbered filenames from the project's former naming convention. **Do not create new numbered reports or refer to the production implementation by a round number.**

## Experiment ledger discipline

1. Describe a hypothesis by **mechanism and encoded domain**. Cite the closest already rejected concrete implementation and say exactly what new mechanism changes. Avoid repeated hash widths, segment counts, boundary positions or file-size gates; call it a spiral and move on.
2. For meaningful prototypes, execute a **reversible standalone/native Go implementation** on actual original Silesia inputs. Exact data, model serialization, framing, metadata, decoder dependencies and throughput all matter. Partial samples are explicitly just component diagnostics; only a complete candidate measures codec improvements.
3. **Candidate-only** ten full original twelve-file Silesia runs for any plausible promotion, with unchanged previously approved and external reference measurements. Report separate medians, every file's exact archive size, full reconstruction, compatible old archives, Go uncached suite/vet/race and Windows crosscompile; do not fake certainty from independently measured timing sessions.
4. **Only accept a tri-metric Pareto result** absent the user's express tradeoff permission. Neither a smaller archive with slower speeds nor a faster archive with larger bytes qualifies. At very small differences, noise and complexity can disqualify a candidate even if observed measurements nominally move favorably.
5. **After every substantive screen, record the exact tested implementation, dataset scope, what was not tested, costs, conclusion and whether the underlying theory remains open**. Reject the concrete candidate rather than incorrectly closing the entire mathematics family.
6. **Report in chat.** Only provide the user a full-codec comparison table (rows Zstd, XZ, current, candidate; four required measured columns plus three delta columns) when a real complete integrated candidate was actually benchmarked. Do not produce an otherwise identical current-only or sample-only table.

## Current priorities

**Source code and broader architecture:** investigate a faster exact suffix-order extraction without serializing/reparsing Go's suffix-array index, and assess the bounded-work feasibility of candidate-cost parsing in mathematical streams. General inspiration from original architectures of LZFSE (sequences, finite-state entropy), LZ4/Snappy (low-cost match/copy and block independence), liblzma/LZMA (hash-chain vs binary-tree match finding), libdeflate (minimum-cost parse), Brotli (contexts and transformed static dictionaries), Zstd (strategy-specific modeling), and non-code mathematical theories. These are source references and hypotheses, never imported compressors or ports.

**Most recent component diagnostic:** [sparse content-defined far-reference anchors](reports/2026-10-19-sparse-far-anchor-screen.md) were implemented and reconstructed all original Silesia files, but no full JLZ archive and no demonstrated incremental compressed bytes; the initial scan was too costly. Status SCREENED; avoid further anchor-mask/window tweaking.

**Previously accepted subsystems and rejected alternatives:** refer to the single [strategy ledger](research/strategy-ledger.md) rather than a numerical succession of "rounds." The pre-2026-10-19 historical development record remains intact in repository Git history and original dated reports. The current engine's historical benchmark provenance is documented under the original archival report `docs/reports/2026-10-19-round57-mdl-bwt-regimes.md`; its filename is legacy-only.

## Current production references

- Complete original source bytes: **211,938,580 B** (12 files), [Silesia SHA-256 manifest](reports/2026-10-18-round23-silesia-sha256.txt).
- `current`: **60,585,717 B**, **5.493743735 s encode / 2.194427893 s decode** (previously measured medians of ten complete runs).
- Zstd-3 external CLI: **66,216,569 B / 1.2561 s / 0.8551 s** (historical; not rerun).
- XZ-6 external CLI: **49,233,340 B / 116.9892 s / 2.3079 s** (historical; not rerun).
- Never imply external CLI timings and in-memory JLZ timings share I/O measurement conventions.

Last documentation reorganization: 2026-10-19 America/New_York. **No new version/round identifiers are to be assigned.**


**Latest research outcome (not an integrated candidate):** Original reversible prefix-doubling suffix-array construction exactly reproduced current BWT bytes on 21 original-source diagnostic blocks but ran ~2–7× slower than the standard-library forward transform; serializing and re-parsing current suffix indexes accounted for ~0.882 seconds of cumulative stage work across all 451 selected BWT source blocks in a separate one-pass component audit. Current source is unchanged. [Suffix-array research report](reports/2026-10-19-suffix-order-extraction-screen.md), [strategy ledger](research/strategy-ledger.md). Do not call these figures full-archive benchmark savings.


[Primary-source Zstd/XZ/Brotli architecture audit](reports/2026-10-19-industry-architecture-audit.md) (2026-10-19): match generation vs priced sequences, full vs bounded optimal parsing, x86 branch-address reversible transforms, entropy/context model overhead, explicit crosswalk to current Go and failed strategy ledger. This is **source-level research only, no measured candidate**. Top candidate hypothesis is bounded, state-aware match-choice pricing within existing native exact/XOR/ADD candidates. Second distinct idea is original branch-address prefilter. No release/round number assigned.


**Native match price feasibility diagnostic:** [Measured nine exact reversible math-selected source-region token streams](reports/2026-10-19-native-match-price-screen.md) after inspecting industry high-quality parser sources. Current fixed 5/6-bit match scorer diverges from the observed 4.47–7.95 bits/byte zero-order entropy of left-behind literal streams on selected Silesia source regions; this is not a proven marginal entropy price or a full archive candidate gain. Keep one **bounded entropy-priced native parse selection** hypothesis open, not a threshold-sweep, with unchanged decoder and transformed-match candidate generator. No integrated candidate or ten-run benchmark.


**Benchmark-verified rejection of two industry-audit hypotheses:** [original cost-aware parser and x86-relative branch transform](reports/2026-10-19-entropy-cost-matching-and-x86-branch-normalization.md). Exact/XOR/ADD candidate repricing found **150,007 B** compressed savings on one complete Silesia screen but slowed compression to **6.667 s** and decompression to **2.311 s**, so no ten-run candidate was warranted. Executable-address normalization materially improved ooffice, **10 candidate-only original twelve-file runs 60,431,550 B / 5.715744 s encoding / 2.266683 s decoding**, 120 exact restorations; **-154,167 B archive, nominal +4.04%/+3.29% slower** than historical current. Candidate-only Go/vet/race/Windows tests pass, backward decoding of earlier archive passes, old binaries reject new mode18. Both strategies **REJECTED** for failing strict tri-metric nonregression. Do not relabel as accepted merely because they improve compression, and do not repeat cost-score or x86 density threshold sweeps. Current remains untouched.
