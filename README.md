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

**Current accepted baseline:** **Round 21** (October 18, 2026). Complete 12-file Silesia corpus: **211,938,580 original bytes → 63,770,766 archive bytes**, down from Round 16's 66,519,205. Ten paired full-corpus runs against unchanged Round 16: median **12.151 → 10.562 s encoding (13.08% faster)** and **4.844 → 4.160 s decoding (14.11% faster)**. Current JLZ is 2,445,803 bytes smaller than the historical **Zstd-3** archive (66,216,569 B); **XZ-6** remains smaller (49,233,340 B). External CLI times are not directly comparable with JLZ's in-memory runs. See [Round 21 canonical report](docs/reports/2026-10-18-round21.md).

**New agent? Read [AGENTS.md](AGENTS.md), then [current state](docs/STATE.md) before changing code.** The accepted source is on `main`. Keep unapproved experiments in disposable local directories; persistent edits land directly on `main`. [Research index](docs/RESEARCH.md) records prior approaches and decisions; [engineering contract](docs/ENGINEERING.md) records immutable constraints.

## Releases

The JLS-compatible lifecycle has reusable Build, ephemeral Test trigger, rolling Nightly, and manual stable Release workflows. Release archives contain compiled executables; the Git repository contains **Go source, tests, required fixtures, workflows and documentation**, not compiled program binaries. See [workflow guide](docs/WORKFLOWS.md).

Licensed under the repository's MIT LICENSE.
