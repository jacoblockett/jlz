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

**Current accepted baseline:** **Round 29** (October 18, 2026). Native, automatically selected 16-bit modular additive-word residuals with signed zigzag encoding compress complete Silesia (**211,938,580 original B**) to **62,644,663 B**, compared with Round 25's **63,104,117 B** (0.728% smaller). Ten **candidate-only** whole-corpus medians are **7.516364 s encode** and **2.821363 s decode**, versus the retained previous Round 25 medians **7.549920 s encode** and **2.839905 s decode**; those sessions are not paired and fractional speed changes are uncertain. Historical external Silesia archives: Zstd-3 **66,216,569 B** and XZ-6 **49,233,340 B**, with non-comparable CLI times. New internal block type 14 is backward-decoded by Round 29, but the older Round 25 decoder rejects Round 29's x-ray archive. See docs/reports/2026-10-18-round29.md for the ten individual timings, Silesia audit, compatibility and exact mathematical transform.

**New agent? Read [AGENTS.md](AGENTS.md), then [current state](docs/STATE.md) before changing code.** The accepted source is on `main`. Keep unapproved experiments in disposable local directories; persistent edits land directly on `main`. [Research index](docs/RESEARCH.md) records prior approaches and decisions; [engineering contract](docs/ENGINEERING.md) records immutable constraints.

## Releases

The JLS-compatible lifecycle has reusable Build, ephemeral Test trigger, rolling Nightly, and manual stable Release workflows. Release archives contain compiled executables; the Git repository contains **Go source, tests, required fixtures, workflows and documentation**, not compiled program binaries. See [workflow guide](docs/WORKFLOWS.md).

Licensed under the repository's MIT LICENSE.
