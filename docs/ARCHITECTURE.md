# Code organization

JLZ is a single executable with one automatic native compressor, a backwards-compatible `JLZ1` decoder, and two CLI verbs. This layout follows Go package boundaries without inventing artificial subsystems.

```text
cmd/jlz/main.go                    CLI parsing, file operations and reporting
internal/jlz/api.go                Small command-to-engine interface
internal/jlz/jlz.go                Container format, shared data structures, decoder
internal/jlz/*.go                  Native compression strategies and stream primitives
internal/jlz/*_test.go             White-box tests next to the implementation
internal/jlz/testdata/fixtures/    Tracked byte-exact correctness fixtures
docs/research/                     Active mathematical research and negative results
docs/reports/                      Historical measurement records
docs/archived-actions/             Inactive hosted-CI definitions, never executable
```

## Why one engine package?

Go directories define separate packages. The current native BWT, mathematical matchers, entropy stream coders, and region selectors intentionally call shared unexported primitives. Splitting each filename into its own package now would require broad interface exposure and risk import cycles. Keep `internal/jlz` cohesive until a measured, stable dependency boundary justifies its own package. Unit tests retain private symbol access.

## Stable behavior

This is a source-layout refactor, **not** a compressor algorithm change. It preserves the `JLZ1` archive modes, checksum, encoded token representations, automatic codec selection, and the two existing CLI verbs. `cmd/jlz` delegates to `internal/jlz` through a small, exported API. The accepted Silesia measurements in prior reports are historical; no new speed or compression claim follows from moving files.

All development and validation are local. Build from the repository root:

```sh
go test -count=1 ./...
go vet ./...
go build -o jlz ./cmd/jlz
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o jlz.exe ./cmd/jlz
```

Package-level correctness fixtures belong under `testdata/fixtures`, which Go tooling excludes from package discovery. For a full-Silesia research run, use a separately acquired verified original corpus outside the Git repository.
