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

**Current accepted baseline:** **Round 25** (October 18, 2026). Native reversible byte-lane streams are automatically selected when complete archive bytes beat the existing mathematical/native/BWT representations. On the 12-file Silesia corpus (**211,938,580 original B**), ten paired whole-corpus runs against unchanged Round 23 measured **63,658,799 → 63,104,117 archive B (0.87% smaller)**, encoding median **7.622 → 7.550 s (0.95% faster)** and decoding median **2.910 → 2.840 s (2.42% faster)**. Eight unrelated real-world files showed 0.71% smaller archives and 4.38% faster decoding, at a 0.75% aggregate encoding penalty. Historical Zstd-3 Silesia size **66,216,569 B** and XZ-6 **49,233,340 B** remain separate, input/toolchain-limited references. See [Round 25 canonical research report](docs/reports/2026-10-18-round25.md).

**New agent? Read [AGENTS.md](AGENTS.md), then [current state](docs/STATE.md) before changing code.** The accepted source is on `main`. Keep unapproved experiments in disposable local directories; persistent edits land directly on `main`. [Research index](docs/RESEARCH.md) records prior approaches and decisions; [engineering contract](docs/ENGINEERING.md) records immutable constraints.

## Releases

The JLS-compatible lifecycle has reusable Build, ephemeral Test trigger, rolling Nightly, and manual stable Release workflows. Release archives contain compiled executables; the Git repository contains **Go source, tests, required fixtures, workflows and documentation**, not compiled program binaries. See [workflow guide](docs/WORKFLOWS.md).

Licensed under the repository's MIT LICENSE.
