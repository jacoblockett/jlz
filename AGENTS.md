# JLZ Agent Instructions

JLZ is a standalone, native, exact-byte lossless compressor written in Go. Your objective is to surpass established compressors in practical compression ratio and performance through original research and rigorous experimentation.

### Development

Preserve simplicity, modularity, backward decoding, and the existing two-command (`compress`, `decompress`) automatic compressor. **Never** use existing compression codecs as JLZ backends, fallbacks, or entropy engines. Keep independent bounded blocks, robust archive validation, controlled memory/workers, and YAGNI. Preserve the committed root `default.pgo` when reproducing accepted performance.

Always commit and push progress directly to `main` in `jacoblockett/jlz`. Never create branches without explicit permission and justification. Use disposable local directories for prototypes; commit validated production changes and research findings, but not rejected experimental source. GitHub Actions are disabled: run tests locally and do not dispatch workflows.

### Repository Guidance

Read current `main` and previous handoff material before work. Consult `docs/STATE.md` for the accepted baseline, `docs/ENGINEERING.md` and `docs/WORKFLOWS.md` for technical and testing constraints, and `docs/RESEARCH.md`, `docs/research/strategy-ledger.md`, `docs/research/layer-placement-matrix.md`, and `docs/research/theory-backlog.md` for research history and open placements. Use `docs/research/attempted-ideas.md` and older reports when checking prior failures. Verify that any local source ZIP matches live `main`.

**These instructions supersede conflicting procedural rules in older documents**, especially the former strict all-three-metrics non-regression gate. Historical reports remain evidence, not current policy. Record new research by mechanism and date, never as numbered rounds; distinguish `current`, `candidate`, `accepted`, `rejected`, and `screened`.

### Research

Maintain a repository history of successful and failed research. **Prioritize studying leading open-source compressors** for architectures that could benefit JLZ, comparing their actual mechanisms against JLZ's implemented capabilities and past experiments. When no promising architectural inspiration remains, optimize the existing architecture through improvements ranging from fundamental algorithms to minor implementation details.

Do not repeat failed parameter tweaks without a genuinely different mechanism or placement. Distinguish a failed implementation from an untested theory, and examine plausible supplementary placements before declaring a strategy exhausted.

Maximize improvements while minimizing regressions. **No fixed acceptance threshold applies.** Evaluate compressed size, encoding/decoding speed, memory, complexity, and compatibility together; ask when an important tradeoff or candidate decision is uncertain.

### Candidates, Integration, and Benchmarking

Always test candidates at component level first, then integrate promising results **as supplements to JLZ's existing architecture**, not isolated replacements. Consider replacement only if isolated experiments demonstrate substantial compression and speed improvements. Even then, test supplementary integration.

Use **only the complete original 12-file Silesia corpus** available in project sources for codec benchmarks (211,938,580 source bytes; verify against the committed SHA-256 manifest). Source-file regions may be used for component screens, and synthetic fixtures for correctness, but neither is a complete JLZ benchmark.

Benchmark integrated candidates with medians of **10 complete candidate-only runs** of all 12 files, including full archive framing/metadata, compression and decompression, and exact byte/hash restoration. Reuse established **ZSTD-3, XZ-6, and latest accepted JLZ** reference benchmarks without rerunning them. Distinguish historical/reference timings, component screens, and full ten-run results. Report per-file regressions and never infer end-to-end gains from component tests alone.

Before promoting code, run `go test -count=1 ./...`, `go vet ./...`, targeted race tests, Windows/amd64 cross-build, backward archive decoding, and corrupt/truncated input checks as applicable. State any unperformed checks and any new decoder or format incompatibility.

### Documentation and Deliverable

Record meaningful successful and failed findings in `docs/reports/YYYY-MM-DD-<strategy>.md`; update `docs/research/strategy-ledger.md` and `docs/STATE.md`. Report **directly in chat** using the repository-root `FINDINGS.md` template, including its benchmark table for measured integrated candidates. Always show reference and candidate deltas clearly, labeling changes **(better)** or **(worse)**; never fabricate missing benchmark results.

Interpret short or nonsensical replies (`k`, `asdf`, etc.) as approval of findings and candidate selection, and instructions to continue your recommended direction, unless the user explicitly objects, questions, or redirects. If an important clarification remains unanswered, ask again directly rather than assuming approval.
