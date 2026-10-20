# JLZ agent guidance: read first

**Current**: accepted native Go compressor on GitHub \`main\`. **Candidate**: no integrated or benchmarked candidate currently active. Do not name the project state using numbered rounds or release identifiers. Earlier dated reports retain legacy filenames solely as historical citations.

**Current complete original-Silesia reference (12 original files; 211,938,580 source bytes):** compressed **60,585,717 B**, historical ten-run medians **5.493743735 s encode**, **2.194427893 s decode**. Fixed external benchmarks, retain without rerun: Zstd-3 **66,216,569 B / 1.2561 s / 0.8551 s**, XZ-6 **49,233,340 B / 116.9892 s / 2.3079 s**. External CLI and JLZ in-memory timings are not directly comparable.

## Mission, architecture, and constraints

- Build **JLZ**, standalone, native, exact-byte lossless compression in Go, pursuing a better practical balance of compressed size, encoding time, and decoding time than mature compressors. No pre-existing compression codec (including DEFLATE, Zstd, XZ/LZMA, Brotli, LZFSE or Snappy) may become a backend, fallback, outsourced entropy engine or hidden production dependency. Study their mathematics and architecture as inspiration; write original JLZ implementations.
- Preserve exactly **two user commands**, \`compress\` and \`decompress\`, with **one automatic compressor** and internal representations selected without user switches. Keep simple Go packages, YAGNI, limited allocation and worker count, independent blocks/groups, legacy archive decoding and robust error bounds. Additive components are preferred; a wholesale replacement needs transformative end-to-end evidence.
- Work on GitHub **\`jacoblockett/jlz\`**, permanent commits **directly to \`main\` only after validation**. Use disposable local directories, not feature branches. GitHub Actions are disabled due billing; **no Actions dispatch/test branches**. For source access read GitHub's live \`main\` and the canonical docs; local snapshots need verification against \`main\` before being treated as production.
- All authoritative user-facing and documentation status labels are **\`current\`**, **\`candidate\`**, **\`accepted\`**, **\`rejected\`**, or **\`screened\`**. **Never introduce a new numbered research round or JLZ round identifier.** Preserve historical reports/files unchanged when possible so their evidence remains discoverable.

## Non-negotiable acceptance gate

**Without the user's explicit permission, promote only changes that preserve or improve ALL THREE measured full-corpus metrics:**
1. **Compressed bytes** of the *complete original 12-file Silesia* corpus, including all headers, models, side-information, and checksums.
2. **Median complete-corpus compression time** across **10 complete candidate-only runs**.
3. **Median complete-corpus decompression time** across **10 complete candidate-only runs**.

A smaller archive does NOT justify a slower compressor **or** decoder. A faster codec does NOT justify larger archives. Ignore claims of improvement based on noisy unpaired sub-percent time differences. Even a three-metric winner may be rejected for disproportionate code, format, memory or complexity costs; report worst per-file regressions and compatibility. If a candidate fails a clear inexpensive screen, **stop before ten full runs** and document precisely what failed; never present screen timings as ten-run medians. Reuse previously measured \`current\` and Zstd/XZ references without rerunning them. Only benchmark using the twelve complete original Silesia files; synthetic small fixtures are for correctness only.

Validate promising candidates with byte-exact restoration including SHA integrity; \`go test -count=1 ./...\`, \`go vet ./...\`, targeted \`go test -race\`, Windows/amd64 cross-build and backwards decoding against old archives. Tests run locally. Forward incompatibility must be explicit.

## Research process and anti-spiral requirement

1. **Read live \`docs/STATE.md\`, \`docs/ENGINEERING.md\`, \`docs/WORKFLOWS.md\`, \`docs/RESEARCH.md\`, \`docs/research/strategy-ledger.md\`, and \`docs/research/theory-backlog.md\`** before implementing. The live repository, not stale chat memory, is authoritative.
2. Treat mathematical/scientific research ideas as an **active experimental queue**. Review \`docs/research/attempted-ideas.md\` and the **strategy ledger** before every new investigation; distinguish the abstract theory from the actual prior code and scope.
3. For each new hypothesis, explain **which failed strategies it resembles, precisely why it differs, and what outcome that difference could change**. Similarity does not automatically forbid a new experiment; changing only a threshold, hash width, context bucket, block/window size or model count without a new mechanism is usually redundant. **Detect spirals yourself, stop, and change architecture.**
4. Build a **real reversible** native proof-of-concept with exact metadata accounting when viable. Distinguish theoretical bound, partial/component screen, one complete-corpus screening pass, and actual ten-run integrated candidate. Do not present an optimistic upper bound as JLZ compression gain.
5. Draw on relevant **existing compressor architecture** continuously; research LZFSE, LZ4, Zstd, liblzma, libdeflate, Brotli and other diverse structures where helpful, but do not copy code or insert their codecs. Identify which components JLZ already implements before proposing anything redundant.
6. Permanently document **positive and negative findings by strategy**. New reports go to \`docs/reports/YYYY-MM-DD-<strategy-slug>.md\`, never numbered rounds. Update \`docs/research/strategy-ledger.md\` and \`docs/STATE.md\`; keep \`docs/RESEARCH.md\` as a concise conceptual index. Rejected source stays out of \`main\`, but research stays in the repository. Historical numbered reports remain accessible as evidence, NOT as new version naming.
7. **Reporting directly in chat**, never primarily as a report download: hypothesis, actual implementation, reversibility, findings, verdict, next distinct mechanism. **Only include comparison and percentage-delta tables if a genuinely integrated candidate was benchmarked**. No redundant baseline-only table after sample screens or theoretical research. For full candidate benchmarks the standard comparison has rows **Zstd, XZ, current, candidate** and columns compressed size, original reduction, encode time, decode time; follow with 3-column candidate/current percentage deltas. Mark historical timings versus new explicitly.
8. Nonsensical or very short messages (e.g. \`k\`) mean **continue autonomous research without clarification**. Recommend changes and treat them as approved only if the strict acceptance gate is actually met; tradeoffs require the user's explicit permission. Work economically, avoid redundant compute, and preserve backward decoding.

## Live links

- [Current implementation, candidate status, and next hypothesis](docs/STATE.md)
- [Strategy outcomes, tested implementations, failures, novelty gate](docs/research/strategy-ledger.md)
- [Scientific/theory queue](docs/research/theory-backlog.md)
- [Older concrete experiments and limitations](docs/research/attempted-ideas.md)
- [Engineering contract](docs/ENGINEERING.md)
- [Local validation workflow](docs/WORKFLOWS.md)
- [Research index and archival reports](docs/RESEARCH.md)
