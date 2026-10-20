# Local-only development, testing and release workflow

**Effective 2026-10-19, GitHub Actions are off due to billing. GitHub is source control only.** Do not dispatch/retry Actions, run scheduled nightlies, create `jls-test/*` branches, or rely on hosted runners. The five earlier YAML definitions are preserved under `docs/archived-actions/`, outside GitHub's executable workflow path.

## Local research and validation
1. Obtain a local source checkout/snapshot from the now-public GitHub repository. The GitHub connector is not itself a local checkout; if the execution environment blocks outbound networking, transfer source explicitly before claiming local tests have run. The local system already has Go 1.23.2 linux/amd64 and the full user-supplied original 12-file Silesia corpus, 211,938,580 B; all canonical Silesia SHA-256 identities were checked.
2. Experiment in disposable local directories. Keep the accepted baseline intact; all durable approved changes go directly to `main` through the GitHub source connector after validation. Reject unbenchmarked experiments without committing candidate code.
3. Run `go test -count=1 ./...`, `go vet ./...`, affected `go test -race -count=1`, and exact byte-for-byte roundtrips including checksum. For substantial candidates, use ten candidate-only **complete original Silesia** encode/decode passes, report median whole-corpus per-run totals plus worst per-file regressions; reuse stored approved JLZ, Zstd and XZ reference measurements without rerunning them. Historical external CLI timings are not directly comparable with local in-memory timings. Unit fixtures can be used only for correctness tests, not benchmark corpora.
4. Pursue strict improvements in archived bytes, encoding time and decoding time simultaneously, not ratio/speed tradeoffs, unless the user explicitly authorizes one. Only promote a complete-corpus candidate after a valid candidate-only ten-run median. Review the strategy ledger and distinctness gate before starting. Prioritize substantial compression throughput gains with near-Zstd-3 archive size. Draw inspiration from open-source compressors without copying their code or employing their codecs as JLZ backends. Record each experiment and decision by mechanism in `docs/reports/YYYY-MM-DD-<strategy-name>.md`, `docs/research/strategy-ledger.md`, and `docs/STATE.md`. Label production `current`, integrated trial `candidate`; do not create numbered research rounds. Only include comparative whole-codec tables for actually benchmarked integrated candidates.

## Windows first
Linux local Go tests run normally. Cross-compile Windows/amd64 using `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o /tmp/jlz.exe .` and verify that the artifact is produced. This does **not** test native Windows runtime or prove Windows throughput. Native Windows tests require an actual Windows machine, not GitHub Actions.

No hosted CI, nightly, Test or stable release execution. Keep binaries, corpora, and temporary benchmark files out of Git. Local archives/releases can be produced on explicit request; do not reinstate Actions without explicit authorization.

## Current package layout
Run all commands from the repository root. The executable source is root `main.go`; production encoder/decoder and adjacent package tests live in `./engine`; correctness fixtures live in `./engine/testdata/fixtures`. Use `JLZ_CORPUS=/path/to/silesia` only for explicitly enabled corpus research/audit tests. `go test -count=1 ./...` and `go vet ./...` operate on all Go packages. No GitHub Actions are permitted.

Feature source packages are rooted at `additive/`, `bwt/`, `match/`, `entropy/`, `lane/`, `predictive/`, `archive/`, and `workers/`. They are built and tested with `go test ./...`; no hosted CI is permitted.


## Placement-aware research gate

Before implementing or rejecting a research strategy, inspect the canonical [layer-placement matrix](research/layer-placement-matrix.md), [strategy ledger](research/strategy-ledger.md), and [historical complementary audit](research/complementary-integration-audit.md). Identify the **specific domain and integration placement** previously tested: original-source transform, native matching, transformed residual, match length/opcode/distance/literal metadata, BWT/MTF symbols, additive residual coordinates, or model/decoder state. Test economically plausible complementary placements **inside existing JLZ encoders/decoders** rather than only standalone replacements. A negative sample screen may close an implementation **at that placement**, not a mathematical theory. Avoid exhaustive pairwise composition of previously CPU-regressing modes. Charge every layer's actual side-information, bounds, model costs, full encoder and decoder work. Only full original Silesia integrated candidate-only ten-run medians can qualify for promotion; strict no-regression on all three metrics remains binding. New exact work and untested placements are recorded at [2026-10-20 layered reassessment](reports/2026-10-20-layered-architecture-reassessment.md).


## Accepted layered BCJ→BWT plus Go PGO (2026-10-20)

**Current** uses reversible per-block E8/E9 32-bit little-endian relative operand normalization before existing native 256KiB BWT/MTF/finite-state coding on automatically selected math regions. Mixed-cross mode 19 holds existing BWT subtype 15 or 16. The decoder reconstructs via existing BWT inverse then branch inverse; old archive modes remain backward-decodable, while old binaries reject newly written mode 19. A root checked-in **default.pgo** profile (16,804 compressed bytes; SHA-256 3b35cda74832360a252b63c1805e025aa1f4ed7a4fe4124d9a55929a6e86289f) is automatically read by ordinary Go 1.23+ compilation. Without the measured profile the same codec was historically REJECTED for encoding regression; **do not claim the approved timings for an unprofiled build**.

Current **12-file original Silesia** accepted complete archive: **60,244,257 B**; candidate-only ten-run medians **5.322679488 s encode / 2.003659634 s decode**, 120 exact reconstructions. Previous accepted historical (not rerun): 60,585,717 B / 5.493743735 s / 2.194427893 s. Net -341,460 B, nominal -3.1138% encode, -8.6933% decode (independently timed sessions). [Actual ten-run data, method and backward-compatibility checks](reports/2026-10-20-native-bcj-bwt-pgo-accepted.md). No pending candidate and no numbered release label.

The checked-in profile is build data, not a mature codec backend. Include it when reproducing benchmarks. Profiling improves the compiler's generated code; it is host/compiler-sensitive and not a cross-platform timing guarantee.


## Updated approved measurement (2026-10-20)

The current accepted native Go engine is source-aligned with GitHub `main` commit `e75bb094e518c8f0a4c74731f0d517143cf66989`. The full original 12-file 211,938,580 B Silesia corpus compresses to **60,244,257 B**, with ten complete candidate-only medians **5.081535442 s encode / 1.9046838445 s decode**. Previous reference **60,244,257 B / 5.322679488 s / 2.003659634 s** is now historical, not to be rebenchmarked automatically. Root `default.pgo` required for comparability; the archive format and old decoding remain unchanged. Pair-LF inverse needs extra scratch memory. See [ten-run source validation](reports/2026-10-20-compact-index-two-hop-inverse-accepted.md).


## New current accepted benchmark (2026-10-20)

Accepted `main` native automatic record-aligned field-lane + Huffman decoder improvement (source commit `37f02082063fed3bb42aa05817dd0afabd0790ae`): complete original 12-file Silesia **58,924,248 B** from **211,938,580 B original**, ten candidate-only complete-run medians **4.8896856215 s compression and 1.876473360 s decompression**. Prior historical accepted **60,244,257 B / 5.081535442 / 1.9046838445**, not rerun. Automatic mode13 source-period lane width 12..64 may require this newer decoder; old archives decode. Original committed `default.pgo` unchanged. [Full report](reports/2026-10-20-record-aligned-native-lanes-accepted.md). Do not rerun older references during new candidate tests.


## Accepted reference after efficient math-region field gating (2026-10-20)

The new current accepted engine `main` source `e632bc02872b983c050785d88297425457bc3e6f` compresses complete original twelve-file Silesia 211,938,580 B into **58,767,739 B**, ten complete candidate-only median **4.7873709325 s encode /1.8721679705 s decode**. Previous accepted historical reference **58,924,248 B /4.8896856215 /1.876473360** must **not** be rerun. [Validation/benchmark report](reports/2026-10-20-cheap-aligned-field-math-regions-accepted.md). `default.pgo` unchanged, no new on-disk formats; small nominal -0.23% decode delta is inconclusive.
