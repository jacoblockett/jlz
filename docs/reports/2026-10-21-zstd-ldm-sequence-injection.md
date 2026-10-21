# Zstd LDM external sequence injection and source-index work elimination

**Date:** 2026-10-21. **Baseline:** original complete 12-file Silesia 211,938,580 B, current JLZ 58,767,739 B (72.271335% reduction) / 4.7873709325 s encode /1.8721679705 s decode (retained accepted ten-run medians). **Result:** four genuinely distinct original-Go implementation placements integrated into disposable full JLZ encode/decode. None demonstrated a convincing compression-time improvement, and each traded away some of the previously screened bounded S2 candidate's compression advantage. **NO PROMOTION**. Production source, root PGO and decoder unchanged.

## Novelty gate and mechanism

The primary [Meta Zstd source-side LDM sequence producer](https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_ldm.c), especially `ZSTD_ldm_generateSequences_internal` / `ZSTD_ldm_blockCompress`, and [multi-threaded LDM job feed](https://github.com/facebook/zstd/blob/dev/lib/compress/zstdmt_compress.c) provide an upstream example of **precomputed, source-verified long-distance sequences consumed by independently scheduled block parsers**. This is an architectural analogy, not a port or external codec.

Nearest failed/screened JLZ work: [historic R63 512KiB ordinary history](2026-10-19-round62-63-math-sequences-far-history.md), [G1 sparse group matching](2026-10-21-zstd-group-local-sparse-ldm.md) and [bounded once-per-group hint producer](2026-10-21-zstd-ldm-sequence-list-scheduling.md). G1 re-indexed far history per block, S2 amortized the producer per four-block group but still **added** the regular near exact/XOR/ADD parser and one-byte lazy comparison. This experiment **directly substitutes** prevalidated long references at known parse positions, avoids inserting some already consumed exact-match bytes, or changes the **producer data structure** rather than trivially adjusting anchor strides or restoring previously rejected near-index depth.

All prototypes originate from the verified accepted native source snapshot plus disposable bounded S2, retain exactly the ordinary 128KiB two-slot exact/XOR/ADD match index, three repeat-distance entries, native syntax/metadata coder, BWT/predictive alternatives, four-block independence and unchanged per-file digests. Only source precomputed hints with distance >128KiB can change match selection. No Zstd backend, mature compressor imports, native assembly, `unsafe`, new user-visible options or code committed.

## Four distinct source-integrated designs

1. **FORCE:** When a group-produced 64B-verified far anchor is encountered at the current 16B sample position, directly emit its reverified longest exact reference rather than also executing ordinary exact/XOR/ADD candidate comparisons and one-byte lazy parse. All non-hint positions retain full original parser.
2. **FORCE+INSERT-SKIP:** Additionally skip filling the near exact/XOR/ADD hash index over the already consumed external exact-match interior, retaining the last 32 bytes for later near matches. Distinct work-elimination technique from changing bucket depth; may lose future nearby reference opportunities.
3. **DIRECT-INDEX:** Instead of a per-group `map[uint64]uint32` with `clear`, use a single file-local 2^17 entry `uint64` key + `uint32` absolute position direct-mapped table for the sparse producer. Absolute previous-source positions act as group generations without clearing; verify source content and group bounds on use. Collisions can lose opportunities but cannot inject false references. S2 downstream parsing remains unchanged.
4. **MAXIMAL-PRODUCER:** At a 64B-verified far match, extend up to native 32KiB maximum and skip subsequent 16B source anchor scans throughout that coverage, emitting one source-side candidate as a non-overlapping sequence. Preserve a single >=128B match as a valid block gate as well as the previous two-hit gate. Uses FORCE consumer.

Each was an **actually coded JLZ integration**, not a mere component test. The standalone map/hash idea was not confused with the historically rejected ordinary match-bucket tag/depth and prior matcher-index `sync.Pool` experiments.

## Original 12-file complete archive results

All original corpus file hashes matched committed `docs/reports/2026-10-18-round23-silesia-sha256.txt` (12/12 PASS). Each implemented variant ran one or more complete encode and decode passes, with existing JLZ framing, file checksum verification and output byte-by-byte comparison to all 12 original sources. Whole-archive sizes were deterministic across repeated complete runs. These are **diagnostic full-corpus screens, not qualifying ten-run candidate performance medians.**

| Actual original JLZ integration | Compression reduction | Δ vs accepted JLZ | Archive B | Size Δ relative to bounded S2 |
|---|---:|---:|---:|---:|
| Accepted JLZ anchor | 72.271335% | 0 pp | 58,767,739 | n/a |
| Previous bounded S2 sequence hints | 72.312068% | +0.040733 pp (better) | 58,681,409 | 0 |
| FORCE | **72.311549%** | **+0.040214 pp (better)** | 58,682,510 | +1,101 B (worse) |
| FORCE+INSERT-SKIP | **72.310844%** | **+0.039510 pp (better)** | 58,684,003 | +2,594 B (worse) |
| DIRECT-INDEX | **72.311631%** | **+0.040296 pp (better)** | 58,682,336 | +927 B (worse) |
| MAXIMAL-PRODUCER | **72.308847%** | **+0.037512 pp (better)** | 58,688,236 | +6,827 B (worse) |

Per-file sizes (B) verify **only Mozilla and Samba differ from accepted JLZ**; ten other files' compressed sizes and all twelve sources remain the same. No per-file regression relative to accepted JLZ; every implementation is worse than earlier S2 for at least one of the two affected files:

| File | Accepted | Previous S2 | FORCE | FORCE+INSERT-SKIP | DIRECT-INDEX | MAXIMAL |
|---|---:|---:|---:|---:|---:|---:|
| Mozilla | 18,452,630 | 18,370,668 | 18,371,113 | 18,372,210 | 18,371,482 | 18,374,758 |
| Samba | 5,134,262 | 5,129,894 | 5,130,550 | 5,130,946 | 5,130,007 | 5,132,631 |

## Timing confidence, runtime and validation

Original Go host Go1.23.2 linux/amd64, original committed `default.pgo`, `GOMAXPROCS=4`, in-memory complete JLZ archive encode/decode. Interleaved diagnostic **S2→DIRECT→DIRECT→S2** compression times [5.712643271, 5.765796198, 5.804706669, 5.721076028] seconds: S2 median 5.716859650 s, DIRECT median 5.785251434 s, **DIRECT +1.20% slower (worse)**. Independent reversed **MAXIMAL→S2→S2→MAXIMAL** compression times [5.897281798, 6.217470471, 5.441601674, 5.714195715]: S2 median 5.829536073 s, MAXIMAL median 5.805738757 s, **MAXIMAL -0.41% nominally faster (better), unstable and much smaller than within-control drift**. FORCE and FORCE+SKIP completed several all-12 source-exact passes with mixed one-pass times (~5.33–6.36 s); one attempted six-run sequential comparison timed out after four completed runs. No missing/timed-out run is counted as a completed pass, and no claim of 10-run qualification is made for any variant. Unchanged accepted JLZ ten-run reference and external Zstd-3 66,216,569 B/1.2561/0.8551 s and XZ-6 49,233,340 B/116.9892/2.3079 s were not rerun. Local timing drift, altered codec workload and historical reference methodology prohibit an unpaired candidate-vs-accepted speed claim.

`go test -count=1 ./...` **PASS** and `go vet ./...` **PASS** on FORCE, FORCE+SKIP, DIRECT and MAXIMAL disposable sources. Targeted `go test -race -count=1 ./engine -run 'TestExistingArchiveCompatibility|TestParallelNativeIntegrity|TestCorruptionDetected|TestMatchV2OverlappingRepeatAndTransforms'` **PASS** for FORCE and DIRECT. `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ./...` **PASS** for FORCE and DIRECT (build only, not Windows native execution). Race and cross-build not separately rerun for FORCE+SKIP or MAXIMAL because no production promotion was contemplated. Existing archived old-JLZ1 decoder compatibility and malformed-input checks passed in unchanged unit coverage. **Forward decode using old JLZ binaries for new >128KiB cross-block reference distances is not established and should be considered incompatible.** No candidate archival format was promoted.

## Defensible conclusion and open mechanisms

**REJECT these four concrete external-sequence injection/probe/producer implementations as production upgrades.** The maximum compression improvement even before CPU and decoder complexity was only **+0.040733 pp** in previous S2; all work-elimination variants slightly reduced that gain. Bypassing candidate work at sparse verified hits saves too little of the total parser work, and sparse producer substitution risks missed near match candidates. Do not waste work tuning 16B anchor masks, 2^17 table size, last-32B insertion tail or adjacent jump thresholds as a "new architecture".

**Do not declare Zstd LDM exhausted.** Potential genuinely different placements still include (a) selection over *whole precomputed command intervals* with native token-cost comparisons, where producer data **replaces** an entire parse interval rather than only a single candidate lookup; (b) jointly optimizing group mode-selection preflight and source sequence discovery so group source is scanned at most once; and (c) a bounded cost-aware raw-sequence schedule that can cheaply prove distant references cannot beat near references. Any future source implementation needs an actual work/cost bound and explicit compatibility accounting. Continue through remaining **OPEN** Zstd families instead of repeating these local implementations.
