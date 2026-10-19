# Round 35: Speed-first research and adaptive-match trial (CI blocked)

Date: 2026-10-18/19. **No completed candidate benchmark and no production promotion.** Approved standalone compressor remains JLZ Round 29 on `main`.

## New user objective
Prioritize practical compression throughput approaching or exceeding Zstd-3, while retaining competitive Silesia archive size. A modest archive-size sacrifice is acceptable when it yields substantial speed. Historical Zstd-3 archive size is 66,216,569 B on the complete original Silesia corpus, against accepted JLZ Round 29's 62,644,663 B. That provides a 3,571,906 B (~5.70%) size gap before JLZ matches the historical Zstd archive. Accepted reference times: JLZ Round 29 in-memory ten-run medians 7.516364 s encode / 2.821363 s decode; historical Zstd CLI 1.2561 / 0.8551 s and XZ CLI 116.9892 / 2.3079 s. These benchmark methodologies and machines differ, so they are NOT a controlled speed comparison or a measured JLZ speedup target. Keep the single automatic two-command native Go codec.

## External open-source inspiration and original adaptation
- Google Snappy `snappy.cc`: adaptive skip of match positions after long runs without successful matches. Adopt the search-budget principle, not source code. https://github.com/google/snappy/blob/main/snappy.cc
- Facebook Zstandard `zstd.1.md`: fast/greedy/lazy/binary-tree strategies and target-length search budgets; use progressive search and good-enough match early termination selectively. https://github.com/facebook/zstd/blob/dev/programs/zstd.1.md
- LZ4 API: acceleration trades ratio for throughput through fewer searches, a natural heuristic input to JLZ's own native matcher. https://github.com/lz4/lz4

Source assessment: JLZ can revisit repeated match-index searching and lazy parse costs; BWT construction through serialized Go suffix arrays and expensive multiple representations are other potentially important cost centers. No CPU profile has yet quantified their relative shares on the complete corpus; avoid ranking these as measured bottlenecks. Existing early region selection and sample gates already help; do not blindly replace them.

## Actual prototype
Updated `match_v2.go` to count consecutive failed match probes, skip up to eight candidate search positions after long unsuccessful runs while indexing all intervening source bytes, and omit single-byte lazy lookahead for matches of length >=64. Existing token and archive formats, checksum and decoding stayed unchanged. Built a complete candidate-only Silesia test intended to measure ten independent full-corpus encode/decode totals, 12 source hashes, per-file selected modes and CPU profile, plus a dedicated ephemeral-branch workflow.

**CI could not execute the tests.** On 2026-10-19, initial GitHub Actions run IDs 39190107746 (candidate) and 39190107831 (Test Trigger), including retries, failed BEFORE runner assignment: 0 executed steps, runner_id 0 and unavailable logs. An additional Ubuntu 26.04 trial at run IDs 39190319911 and 39190319847 reproduced the same failure. This is a provisioning/billing-or-entitlement concern, not evidence that candidate code or tests failed. No compile, vet, full-corpus correctness, candidate archive sizes, per-file regressions, or ten-run performance values were obtained. Do not infer the precise GitHub billing cause without account-run annotations. All original 12 Silesia SHA-256 inputs were independently verified from the user-provided local archive (211,938,580 B), but this does NOT replace candidate execution.

## Historical references, not new candidate measurements
| Model | Archive bytes | Reduction from 211,938,580 B | Encode | Decode |
|---|---:|---:|---:|---:|
| Zstd-3, historical CLI | 66,216,569 | 68.76% | 1.2561 s | 0.8551 s |
| XZ-6, historical CLI | 49,233,340 | 76.77% | 116.9892 s | 2.3079 s |
| Approved JLZ Round 29, historical in-memory 10 runs | 62,644,663 | 70.44% | 7.516364 s | 2.821363 s |
| Round 35 match-skip candidate | N/A (CI unavailable) | N/A | N/A | N/A |

| Candidate compression delta | Encode-time delta | Decode-time delta |
|---:|---:|---:|
| N/A | N/A | N/A |

## Decision and next action
**DO NOT ACCEPT the unmeasured Round 35 candidate.** Source prototype, benchmark test, and trial workflow removed from `main`; the approved Round 29 source is restored. No theoretical family or speed hypothesis is rejected. First restore GitHub Actions runner eligibility (check Actions spending limits/payment health, Actions settings and run annotations), then rerun the adaptive-skip/good-match hypothesis with 10 complete Silesia passes and per-file/mode CPU profiles. Consider independent self-contained faster original matchfinding, simpler coding passes and competitive Huffman/BWT cost only after profiling. Record wins and losses with fully charged archived byte sizes; prioritize genuinely large speed gains, not cosmetic optimizations.
