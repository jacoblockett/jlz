# Zstd native exact-match short-sequence decoder screen

**Date:** 2026-10-20 (local America/New_York). **Status:** source-workload instrumented; two original Go decoder alternatives implemented in disposable copies, full original Silesia reconstructed exactly, no proven throughput improvement. No candidate promoted or archive mode changed.

## Primary inspiration and novelty collision

Meta Zstd's `zstd_decompress_block.c` has specialized sequence-copy paths, notably short/long and overlapping match execution. JLZ already supports the same exact-copy semantics via Go's native `copy` with exponential doubling for an overlapping reference. The earlier JLZ **direct-column decoding** bypass was actually benchmarked and rejected, and earlier **encoder-side transformed-match word verification** is not a new decoder tactic. A prior accepted-corpus profile counted 3,665,755 exact copies and only 3,410 XOR/ADD copies (82,880 bytes). This trial first adds the previously missing **exact length and overlap** census and then codes two cheap, portable, safe alternatives specifically for the hot exact-copy regime.

## Exact accepted original workload

On the accepted original 12-member Silesia source (211,938,580 input B), actual native encoded output **58,767,739 B**, all 12 decoded byte-exact; temporary atomic counters inserted only into disposable `match.Expand`, not into production. The groups sum exactly to the earlier accepted-corpus census of 3,665,755 exact expansions / 53,605,763 referenced bytes:

| Exact-reference class | Calls | Reconstructed bytes |
|---|---:|---:|
| Nonoverlapping, length <=8 | 1,782,811 | 11,886,344 |
| Nonoverlapping, 9..16 | 1,201,101 | 14,759,285 |
| Nonoverlapping, 17..32 | 467,332 | 11,074,193 |
| Nonoverlapping, >32 | 171,885 | 13,570,307 |
| **Overlapping**, length <=16 | 26,173 | 244,530 |
| **Overlapping**, length >16 | 16,453 | 2,071,104 |

**81.4% of exact match commands** are nonoverlapping copies of at most 16 bytes; **only 42,626** commands (1.16%) overlap. This identifies short *nonoverlapping* exact copies as the only credible target without replaying old direct-column reconstruction work.

## Fully integrated disposable Go implementations

1. **Portable two-word 8..16B copy**: for nonoverlapping lengths 8..16, use native `encoding/binary` little-endian 64-bit load/store at the start and another load/store anchored at the tail, preserving exactly `length` bytes without writing beyond the destination or re-reading any future output. All other cases use existing `match.Expand`; no unsafe/assembly/external dependency.
2. **Portable scalar short copy**: replace native `copy` for nonoverlapping lengths <=8 with a directly indexed byte loop; all other paths unmodified.

Each implementation decoded **all 12 complete archives and verified bytes.Equal on every file** across 3 complete original-corpus runs, with all framing/digests and unchanged archives. All 12 outputs totaled 58,767,739 B per run. Encode source was unchanged, but every benchmark included full encode and decode. In-memory timings are subject to machine/GC variance and are only **three runs**, never ten-run qualification:

| Full original Silesia (3 runs, within session) | Decode s (run 1 / 2 / 3) | Decode **median** s | Delta vs local 3-run control |
|---|---|---:|---|
| Accepted existing native `copy` control | 1.999890 / 1.964327 / 1.909954 | **1.964327** | 0 |
| 8..16 B two-word replacement | 1.979587 / 1.977987 / 1.930280 | **1.977987** | +0.695% **(worse)** |
| <=8 B scalar loop | 1.915672 / 1.947815 / 2.042158 | **1.947815** | -0.841% nominal **(better)**, unstable |

Go compiler/host wall-clock noise overlaps every observed difference; the tiny scalar median shift is not persuasive relative to its 0.126-second run spread. For both disposable versions, full `go test -count=1 ./...` and `go vet ./...` passed. No ten-run improvement, race or Windows candidate validation claimed because no qualifying production change was proposed.

## Decision

**Reject the 8..16 B two-word implementation** as an improvement; **defer the scalar <=8 B change** because there is no reliable speed evidence and added complexity for a code path Go's optimized `copy` already handles. This does not exhaust the broader Zstd decoder family: exact-copy stream-buffer churn, fused bound checks and more substantive group-level sequence execution remain possible, but must be differentiated from the already-rejected direct-column bypass and justified by CPU profiling. Avoid merely sweeping adjacent word-count cutoffs.

Historical accepted full original corpus **58,767,739 B /4.7873709325 s encoder /1.8721679705 s decoder**, Zstd-3 **66,216,569 B /1.2561 /0.8551 s**, XZ-6 **49,233,340 B /116.9892 /2.3079 s** retained, not rebenchmarked as anchors. Local three-run times are diagnostic, not replacement accepted references. No format, production source, PGO, independent decoding groups or backward compatibility changes.
