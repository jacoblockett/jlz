# Accepted native layered BCJ→BWT with built-in Go profile-guided optimization

**Recorded 2026-10-20. Status: ACCEPT, current on main after integration. No numbered research round.**

## Architectural diagnosis and hypothesis

JLZ's mathematically selected `ooffice` original contains densely repeated executable-style 0xE8/0xE9 relative displacements; a reversible relative-address normalization can expose additional redundancy to **the existing JLZ BWT/MTF/finite-state pipeline** rather than replacing the engine with a new codec. The original **unprofiled** additive BCJ→BWT integration already established a ten-run full corpus archive improvement but a nominal **1.17% encoding-time regression** and was correctly REJECTED (see [layered placement audit](2026-10-20-layered-architecture-reassessment.md)). The new idea is **architecture-level cooperation**: retain the proven BCJ→BWT native layer while restoring and improving its end-to-end runtime using the Go compiler's built-in PGO, trained against **complete original Silesia native JLZ work**, without touching any mature compressor backend or the decoder's other modes.

An earlier standalone BCJ→math option improved size but slowed encode/decode and remains rejected. The accepted placement is different: **BCJ before already existing BWT in selected math regions**, with original-file sampled automatic selection, and PGO throughout the complete native implementation. An ordinary `go build` uses the checked-in `default.pgo` automatically under Go1.23+. No compressor levels, extra CLI verbs or external compressor. Historical decoders reject new mode 19; new decoder retains all old modes.

## Exact implementation

- `engine/branch_relative.go`: original reversible block-local transform. For 0xE8/0xE9 marker with four operand bytes, encode 32-bit little-endian `operand + uint32(markerOffset+5)` modulo 2^32; inverse subtracts, and both skip the same next four bytes. Works on arbitrary source bytes independent of executable authenticity. Automatic gate requires >=64 KiB block and branch-marker density >=2%, omitting 4-byte operands from the scan. New independently bounded mixed-cross block **mode 19** contains 1 byte for the already valid BWT subtype15/16 and its native payload. Entire prior file-level SHA256 and per-block varint framing continue to apply.
- `engine/region_fast.go`: add the gate as an early alternate before the current prediction / exact-XOR-ADD matching pipeline, and return to the current native BWT/MTF/entropy implementation when selected. No double-encoding of the matched region and no external backend.
- `engine/cross.go`: validate and decode only subtypes 15/16, apply original inverse BWT and inverse relative normalization, preserve old modes and bounds.
- `engine/branch_relative_test.go`: exact synthetic roundtrip and corrupt mode / short payload / truncated stream rejection tests.
- **`default.pgo`** (repo root): committed Go CPU-profile protobuf/gzip binary **16,804 B**, SHA-256 **3b35cda74832360a252b63c1805e025aa1f4ed7a4fe4124d9a55929a6e86289f**, Git blob **382b3c59a75bfa11bc60968f71cdde67dbe8d634**. A deterministic one-in-three CPU-stack-sample version of an original complete-12-file candidate encode/decode profile, removing unused location/function records; it retains a distributed source-workload profile and Go 1.23 compiler reads it without build flags. This is **non-executable optimization data**, NOT a preexisting or shipped compression codec. The original 29KiB full profile and compact 16.8KiB sample are archived locally for research; the **exact committed compact profile** underwent its own full ten-run benchmark, not a claimed prediction from the larger profile. An `ooffice`-only ~5.4KiB PGO profile was actually rejected (ten-run encode median 5.561 s) and was not committed. The full-trace 29KiB PGO sample and 16.8KiB profile both demonstrated three-metric Pareto improvements, but only the latter is production.

## Complete original 12-file Silesia source and archive

Original **211,938,580 B**, verified 12/12 against retained canonical SHA-256 source manifest `docs/reports/2026-10-18-round23-silesia-sha256.txt`.

| Silesia original | Historical current archive B | Accepted candidate archive B | Exact change B |
|---|---:|---:|---:|
| dickens | 3,155,856 | 3,155,856 | 0 |
| mozilla | 18,609,139 | 18,609,139 | 0 |
| mr | 2,729,261 | 2,729,261 | 0 |
| nci | 2,163,199 | 2,163,199 | 0 |
| ooffice | 3,033,712 | 2,692,252 | **-341,460** |
| osdb | 3,425,086 | 3,425,086 | 0 |
| reymont | 1,455,234 | 1,455,234 | 0 |
| samba | 5,134,262 | 5,134,262 | 0 |
| sao | 6,261,713 | 6,261,713 | 0 |
| webster | 9,835,416 | 9,835,416 | 0 |
| xml | 515,852 | 515,852 | 0 |
| x-ray | 4,266,987 | 4,266,987 | 0 |
| **TOTAL** | **60,585,717** | **60,244,257** | **-341,460** |

## Ten complete candidate-only corpus runs, source-aligned Go 1.23 build + exact committed compact profile

Each pass encoded and decoded **every original one of the 12 Silesia files**, compared bytes and verified the original archive digest. The newly compiled candidate used `go test -c -pgo=/mnt/data/bcj_pgo_compact3.prof` on the source-aligned local approved engine + three additive Go source edits. Go's ordinary root `go build` was separately confirmed to select the identical `default.pgo`. No established current JLZ, Zstd or XZ reference was rebenchmarked.

| Candidate run | Full archive B | Compress seconds | Decompress seconds |
|---:|---:|---:|---:|
| 1 | 60,244,257 | 5.310299246 | 2.033128989 |
| 2 | 60,244,257 | 5.592287554 | 2.090681727 |
| 3 | 60,244,257 | 5.191804504 | 1.914323542 |
| 4 | 60,244,257 | 5.244895464 | 1.991144544 |
| 5 | 60,244,257 | 5.449500666 | 1.989730684 |
| 6 | 60,244,257 | 5.391263002 | 2.036332067 |
| 7 | 60,244,257 | 5.304810807 | 2.016174724 |
| 8 | 60,244,257 | 5.203980991 | 1.983169917 |
| 9 | 60,244,257 | 5.335059729 | 1.987386988 |
| 10 | 60,244,257 | 5.490322835 | 2.071632344 |
| **Median, each independent series** | **60,244,257** | **5.322679488** | **2.003659634** |

Retained previous **current** historical (independent/unpaired session, NOT rerun): **60,585,717 B, 5.493743735 seconds encoding, 2.194427893 seconds decoding**. Retained historical external Zstd-3 **66,216,569 B / 1.2561 s encode / 0.8551 s decode** and XZ-6 **49,233,340 B / 116.9892 s encode / 2.3079 s decode**; different external CLI I/O measurement convention.

- Fully framed compressed archive **-341,460 B = -0.563598%**.
- Original reduction changes from **71.413550%** to **71.574662%** (+0.161112 percentage points).
- Median encoding: **-3.113801%** versus retained previous accepted.
- Median decoding: **-8.693303%** versus retained previous accepted.
- Measurements are **independent**, rather than controlled paired reference/candidate timing. Improvements are observations on this host, not universal processor guarantees.

**ACCEPT only the entire paired *native BCJ→BWT plus compact Go PGO build*, not unprofiled BCJ alone.** The unprofiled ten-run candidate was **60,244,257 B / 5.557846 s encode / 2.108719 s decode**, which missed the user's earlier encoding gate and remains an independently documented REJECT.

## Full validation and compatibility

- Canonical original twelve-file SHA-256 **12/12 PASS**; **120/120 candidate roundtrips** exactly reproduce original full files with embedded checksums.
- **Uncached `go test -count=1 ./...` PASS**, **`go vet ./...` PASS**, targeted `go test -race -count=1 ./engine -run '^TestLayeredBCJBWTMalformedAndExact$'` PASS, **CGO=0 Windows/amd64 build PASS** with Go1.23.2 (not a native Windows runtime test).
- Verified ordinary Linux `go build .` detects `default.pgo` automatically and produces a root native CLI binary; independent CLI compressed original `ooffice` to **2,692,252 B** and restored exact SHA-256.
- New candidate decoder successfully decompressed an actual previously accepted original `mr` archive; original SHA-256 verified. Existing older binaries reject a new mode19 archive as expected (**forward compatibility not claimed**), but previously accepted archive decoding remains.
- Source-alignment sanity: the previous unprofiled local worktree inherited four noncritical source-layout differences compared with GitHub `main` (BWT local `span` vs constant, mixed-cross output prefix sums vs fixed 256KiB groups, plus equivalent finite-state/BWT MDL code layout). Before the **promotable PGO measurement**, the candidate was reconciled with main's 256KiB source/history decoding and BWT block layout; all ten full original Silesia archive identities remained exact. Permanent integration is a minimal patch **fetched from live GitHub main**, not blindly copying an entire older source tree.
- Profile-guided optimization is an optional Go compiler feature tied to target Go/architecture behavior, not an algorithmic cross-platform guarantee. The strict numerical Pareto evidence is on this recorded Go1.23.2 Linux host, not native Windows.

## Decision and further research

**ACCEPT H1:** additive placement that captures genuine missed executable structure, with built-in Go PGO reducing encoder/decoder costs enough to satisfy all three full-original-Silesia median metrics. This is a real architecture+build result, not a theoretical proposal or a mature-codec wrapper.

Next research should analyze actual per-file and pipeline differences against Zstd/XZ and seek similarly complementary mechanism improvements, rather than repeating opcode-density thresholds, cost-aware parser constants, BWT block-size sweeps, or sparse-defect layouts. Both the empirical new success and old related failures stay in [strategy ledger](../research/strategy-ledger.md) and [placement matrix](../research/layer-placement-matrix.md). Do not invent new numbered round identifiers.