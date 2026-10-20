# JLZ

JLZ is an experimental, **standalone, byte-exact lossless compressor**, written in Go 1.23 using the standard library. One automatic algorithm, two operations, no levels or experimental CLI switches. It does not call Zstd, XZ, DEFLATE, or any other mature compressor internally.

```sh
go test -count=1 ./...
go vet ./...
go build -o jlz .
./jlz compress input.bin
./jlz decompress input.bin.jlz
```

Custom output paths are supported as a third argument. JLZ writes `JLZ1` containers; older MDN/MDC containers require their respective historical decoders. This remains a bounded-memory research prototype with a 512 MiB input/output limit, not a hardened service.

**Current implementation (`main`):** Native automatic mathematical, predictive, additive and selectively enabled BWT/bijective/finite-state compression, with additive reversible branch-relative normalization on selected source blocks before the existing BWT pipeline. The complete original 12-file Silesia corpus (**211,938,580 original B**) compresses to **60,244,257 B** (**71.5747% reduction**). Current newly measured 10-complete-run medians with the committed root **default.pgo**: **5.322679 s compression** and **2.003660 s decompression**. The previous accepted reference was 60,585,717 B / 5.493744 s / 2.194428 s. The checked-in Go profile is applied by ordinary Go 1.23+ builds. See [accepted benchmark](docs/reports/2026-10-20-native-bcj-bwt-pgo-accepted.md). Historical external archive sizes: Zstd-3 **66,216,569 B**, XZ-6 **49,233,340 B**; external CLI benchmark timings use a different convention. See [current state](docs/STATE.md) and the [strategy ledger](docs/research/strategy-ledger.md) for fully framed benchmark provenance, prior accepted mechanisms, rejected candidate implementations, and open research.

**Candidate:** none integrated or benchmarked. JLZ no longer assigns numbered round identifiers to its accepted state or new experiments. Historical numbered filenames remain accessible for provenance. No change is promoted without explicit user permission unless it preserves or improves **all three** whole-Silesia metrics: compressed bytes, compression time, and decompression time.

**New agent? Read [AGENTS.md](AGENTS.md), then [current state](docs/STATE.md) before changing code.** The accepted source is on `main`. Keep unapproved experiments in disposable local directories; persistent edits land directly on `main`. [Strategy ledger](docs/research/strategy-ledger.md) records prior approaches and decisions; [research index](docs/RESEARCH.md) links evidence; [engineering contract](docs/ENGINEERING.md) records immutable constraints.

## Local development and builds

GitHub Actions are disabled. Go tests and candidate-only Silesia benchmark runs are executed locally. Go 1.23+ on Linux supports Windows/amd64 cross-compilation; native Windows execution requires a Windows runtime. The previous workflow definitions are archived as inactive reference files. See [local workflow](docs/WORKFLOWS.md).

## Repository layout

`main.go` is the source entry point. Root-level feature packages (`additive/`, `bwt/`, `match/`, `entropy/`, `lane/`, `predictive/`, `archive/`, `workers/`) implement independent algorithms; `engine/` coordinates JLZ1 archives and legacy compatibility. Correctness fixtures live under `engine/testdata/fixtures/`. See [architecture](docs/ARCHITECTURE.md).

To run without building: `go run . compress input.bin`. Windows: `go build -o jlz.exe .`.

Licensed under the repository's MIT LICENSE.
