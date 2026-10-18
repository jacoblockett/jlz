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

**Current accepted baseline:** **Round 23** (October 18, 2026). Full 12-file Silesia input: 211,938,580 B. Local ten paired whole-corpus runs against unchanged Round 21 source: **64,505,708 → 63,658,799 archive bytes (1.31% smaller)**, encoding median **6.712879 → 6.351571 s (5.38% faster)** and decoding median **2.515717 → 2.330026 s (7.38% faster)**. The older Round 21 run reported 63,770,766 B on a prior environment, which was not reproduced with the newly uploaded Silesia corpus and Go 1.23.2; no cross-host speed comparison is implied. Historical independent Silesia archive references: Zstd-3 **66,216,569 B**, XZ-6 **49,233,340 B**; their CLI timings are not directly comparable. See [Round 23 report](docs/reports/2026-10-18-round23.md) and its exact corpus SHA-256 manifest.

**New agent? Read [AGENTS.md](AGENTS.md), then [current state](docs/STATE.md) before changing code.** The accepted source is on `main`. Keep unapproved experiments in disposable local directories; persistent edits land directly on `main`. [Research index](docs/RESEARCH.md) records prior approaches and decisions; [engineering contract](docs/ENGINEERING.md) records immutable constraints.

## Releases

The JLS-compatible lifecycle has reusable Build, ephemeral Test trigger, rolling Nightly, and manual stable Release workflows. Release archives contain compiled executables; the Git repository contains **Go source, tests, required fixtures, workflows and documentation**, not compiled program binaries. See [workflow guide](docs/WORKFLOWS.md).

Licensed under the repository's MIT LICENSE.
