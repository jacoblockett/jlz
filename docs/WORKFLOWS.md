# Local-only development, testing and release workflow

**Effective 2026-10-19, GitHub Actions are off due to billing. GitHub is source control only.** Do not dispatch/retry Actions, run scheduled nightlies, create `jls-test/*` branches, or rely on hosted runners. The five earlier YAML definitions are preserved under `docs/archived-actions/`, outside GitHub's executable workflow path.

## Local research and validation
1. Obtain a local source checkout/snapshot from the now-public GitHub repository. The GitHub connector is not itself a local checkout; if the execution environment blocks outbound networking, transfer source explicitly before claiming local tests have run. The local system already has Go 1.23.2 linux/amd64 and the full user-supplied original 12-file Silesia corpus, 211,938,580 B; all canonical Silesia SHA-256 identities were checked.
2. Experiment in disposable local directories. Keep the accepted baseline intact; all durable approved changes go directly to `main` through the GitHub source connector after validation. Reject unbenchmarked experiments without committing candidate code.
3. Run `go test -count=1 ./...`, `go vet ./...`, affected `go test -race -count=1`, and exact byte-for-byte roundtrips including checksum. For substantial candidates, use ten candidate-only **complete original Silesia** encode/decode passes, report median whole-corpus per-run totals plus worst per-file regressions; reuse stored approved JLZ, Zstd and XZ reference measurements without rerunning them. Historical external CLI timings are not directly comparable with local in-memory timings. Unit fixtures can be used only for correctness tests, not benchmark corpora.
4. Prioritize substantial compression throughput gains with near-Zstd-3 archive size. Draw inspiration from open-source compressors without copying their code or employing their codecs as JLZ backends. Record each experiment and decision in `docs/reports/`, `docs/STATE.md`, and `docs/RESEARCH.md`.

## Windows first
Linux local Go tests run normally. Cross-compile Windows/amd64 using `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o /tmp/jlz.exe ./cmd/jlz` and verify that the artifact is produced. This does **not** test native Windows runtime or prove Windows throughput. Native Windows tests require an actual Windows machine, not GitHub Actions.

No hosted CI, nightly, Test or stable release execution. Keep binaries, corpora, and temporary benchmark files out of Git. Local archives/releases can be produced on explicit request; do not reinstate Actions without explicit authorization.

## Current package layout
Run all commands from the repository root. The executable is `./cmd/jlz`; production encoder/decoder and adjacent package tests live in `./internal/jlz`; correctness fixtures live in `./internal/jlz/testdata/fixtures`. Use `JLZ_CORPUS=/path/to/silesia` only for explicitly enabled corpus research/audit tests. `go test -count=1 ./...` and `go vet ./...` operate on both Go packages. No GitHub Actions are permitted.
