# Round 34: selective match-length grammar integration (rejected)
Date: 2026-10-18/19. Current approved baseline remains Round 29. This report records an actual integrated JLZ candidate, not a source-region proxy.

## Hypothesis and implementation
A bounded 16-rule Re-Pair/Tietze-inspired straight-line grammar might improve already selected match-length streams with little CPU overhead. Implemented new native packed stream codec 5; 16 hierarchical non-overlapping pair substitutions over lengths only; four-byte rules; independently native-packed low/high coordinates; backwards-only generator references and bounded expansion; malformed input rejection. Added a first-8192-byte byte-pair concentration gate requiring >=384 repeats of a pair before grammar search. Added only as an optional stream representation in ordinary, split-length and cross-region match writers, with measured full-stream selection and old decoders retained.

## Performance and source evidence
Only the original 12-file Silesia corpus, 211,938,580 B, tested; all twelve SHA-256 hashes verified against the canonical Round 23 manifest. The accepted JLZ and external baselines below are **historical reused** measurements, NOT rerun. Candidate 10 complete-run medians obtained on GitHub Actions ubuntu-24.04, Go 1.23.12; all runs perform full in-memory compression and byte-for-byte decompression of each source, including archive framing/digest. Since baseline sessions were separate, CPU comparisons are indicative, not paired.

| Model | Complete archive | Reduction | Encode s | Decode s |
|---|---:|---:|---:|---:|
| Zstd-3 historical CLI | 66,216,569 B | 68.76% | 1.2561 | 0.8551 |
| XZ-6 historical CLI | 49,233,340 B | 76.77% | 116.9892 | 2.3079 |
| JLZ Round 29 accepted | 62,644,663 B | 70.44% | 7.516364 | 2.821363 |
| Round 34 candidate, ten runs | 62,417,637 B | 70.55% | 16.231375 | 5.415786 |

| Candidate vs accepted archive-size delta | Encode time delta | Decode time delta |
|---:|---:|---:|
| -0.362403% (227,026 B saved) | +115.9472% | +91.9564% |

All ten complete runs produced precisely 62,417,637 B and restored all 12 originals. Per-file candidate archive sizes: dickens 3,471,769; mozilla 18,461,596; mr 2,898,511; nci 2,456,164; ooffice 3,022,590; osdb 3,459,257; reymont 1,617,920; samba 5,040,328; sao 6,261,713; webster 10,883,616; xml 577,186; x-ray 4,266,987. Compared with the accepted complete total, savings are 227,026 B; preserve per-file baseline comparison as a separate audit when necessary.

The first synthetic-match metadata *proxy* scan used a different simple LZ-style tokenizer on the same 12 original Silesia files; it saved 147,290 B in its own packed length streams at a gate threshold of 384, but this is **not** JLZ archive improvement or a selected-stream measurement. The complete production-tokenized candidate above is authoritative.

## Validation, decision, and next direction
Go source tests including stream inverse, truncation and malformed-input checks, vet, targeted race checks, six-platform packaged cross-compilation and normal Test rolling prerelease all passed after correcting two test-harness errors (copied OSDB SHA-256 typo and malformed-rule byte-offset test). Candidate-only 10-run complete corpus CI succeeded. The encoder's extra pair passes, temporary dense counters and nested grammar packing inflate CPU far beyond its modest size gain.

**REJECT Round 34 and retain Round 29 source and decoder.** Restore all modified production Go files on main, remove candidate codec/tests and research-only CI workflow, retain this report and update the active backlog. This rejects the particular bounded implementation and gate, not Tietze's mathematical family. Next investigate better structural coding in existing selected streams with a significantly cheaper selector/inverse, particularly other still-open theory×stream pairings (low-rank bitplane relationships, shared group models, or algebraic dependencies). Zstd/XZ timings historically include CLI I/O and must not be compared as matched CPU results.
