# JLZ

JLZ is an experimental, **standalone, byte-exact lossless compressor**, written in Go 1.23 using the standard library. One automatic algorithm, two operations, no levels or experimental CLI switches. It does not call Zstd, XZ, DEFLATE, or any other mature compressor internally.

```sh
go test -count=1 ./...
go vet ./...
go build -o jlz ./cmd/jlz
./jlz compress input.bin
./jlz decompress input.bin.jlz
```

Custom output paths are supported as a third argument. JLZ writes `JLZ1` containers; older MDN/MDC containers require their respective historical decoders. This remains a bounded-memory research prototype with a 512 MiB input/output limit, not a hardened service.

**Current accepted baseline:** **Round 29** (October 18, 2026). Native, automatically selected 16-bit modular additive-word residuals with signed zigzag encoding compress complete Silesia (**211,938,580 original B**) to **62,644,663 B**, compared with Round 25's **63,104,117 B** (0.728% smaller). Ten **candidate-only** whole-corpus medians are **7.516364 s encode** and **2.821363 s decode**, versus the retained previous Round 25 medians **7.549920 s encode** and **2.839905 s decode**; those sessions are not paired and fractional speed changes are uncertain. Historical external Silesia archives: Zstd-3 **66,216,569 B** and XZ-6 **49,233,340 B**, with non-comparable CLI times. New internal block type 14 is backward-decoded by Round 29, but the older Round 25 decoder rejects Round 29's x-ray archive. See docs/reports/2026-10-18-round29.md for the ten individual timings, Silesia audit, compatibility and exact mathematical transform.

**New agent? Read [AGENTS.md](AGENTS.md), then [current state](docs/STATE.md) before changing code.** The accepted source is on `main`. Keep unapproved experiments in disposable local directories; persistent edits land directly on `main`. [Research index](docs/RESEARCH.md) records prior approaches and decisions; [engineering contract](docs/ENGINEERING.md) records immutable constraints.

## Local development and builds

GitHub Actions are disabled. Go tests and candidate-only Silesia benchmark runs are executed locally. Go 1.23+ on Linux supports Windows/amd64 cross-compilation; native Windows execution requires a Windows runtime. The previous workflow definitions are archived as inactive reference files. See [local workflow](docs/WORKFLOWS.md).

## Repository layout

- `cmd/jlz/`: CLI entry point for the existing `compress` and `decompress` commands.
- `internal/jlz/`: native compression/decoding engine, codecs, and package-level regression tests.
- `internal/jlz/testdata/fixtures/`: versioned correctness inputs (not performance benchmarks).
- `docs/`: engineering notes, active research, historical reports, and inactive Actions definitions.

Go organizes code by **package boundary**, not one directory per algorithm. The codec components intentionally share a single internal package while their dependencies remain tightly coupled; splitting them into artificial subpackages would risk cyclic imports and complicate testing. See [architecture](docs/ARCHITECTURE.md).

To run without building: `go run ./cmd/jlz compress input.bin`. Windows: `go build -o jlz.exe ./cmd/jlz`.

Licensed under the repository's MIT LICENSE.
