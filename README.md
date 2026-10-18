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

**Current accepted baseline:** Round 15 (October 18, 2026). Complete 12-file Silesia corpus, 211,938,580 input bytes, 66,797,408 compressed bytes. Median encoding time 6.942 s, based on ten paired in-memory full-corpus runs versus Round 12 at 7.331 s. Historical Zstd-3 total: 66,216,569 bytes; XZ-6: 49,233,340 bytes. External historical CLI timing is **not comparable** to the JLZ in-memory timing. See [canonical report](docs/reports/2026-10-18-round15.md).

**New agent? Read [AGENTS.md](AGENTS.md), then [current state](docs/STATE.md) before changing code.** The accepted source is on `main`. Keep experimental work in a branch until verified. [Research index](docs/RESEARCH.md) records prior approaches and decisions; [engineering contract](docs/ENGINEERING.md) records immutable constraints.

## Releases

The JLS-compatible lifecycle has reusable Build, ephemeral Test trigger, rolling Nightly, and manual stable Release workflows. Release archives contain compiled executables; the Git repository contains **Go source, tests, required fixtures, workflows and documentation**, not compiled program binaries. See [workflow guide](docs/WORKFLOWS.md).

Licensed under the repository's MIT LICENSE.
