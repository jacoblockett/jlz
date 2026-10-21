# Zstd sequence expansion versus actual JLZ match-copy workload

**Date:** 2026-10-20. **Type:** source-level decoder hotloop audit with instrumented **real accepted JLZ archive**. **Disposition:** SOURCE-AUDITED / WORKLOAD SCREEN; defer new word-wide transformed XOR/ADD decoder experiment on the current corpus due vanishingly low workload. **No candidate code change or 10-run benchmark.** This is not a global rejection of optimized native match expansion.

## Explicit novelty check

Primary Zstd source inspected: [`lib/decompress/zstd_decompress_block.c`](https://github.com/facebook/zstd/blob/dev/lib/decompress/zstd_decompress_block.c), Git blob `8d108c72ece8b889b6b6a5701412686afa7544ff`. Its sequence decoder is optimized for copying literal and match spans, with fast/long paths and carefully managed overlapping matches. Existing JLZ `match/expand.go` already calls Go's optimized `copy` for exact backward references, then grows overlapping exact references by repeated doubling; this part is **not new**. The earlier [native word-wide transformed-match verification study](2026-10-20-transformed-match-wide-check-and-fse-abort.md) wrote/benchmarked **encoder-side** XOR/ADD match-length validation and rejected its concrete implementation. This screen instead considers **decoder-side** word-wide arithmetic on XOR/ADD matches, a genuinely different placement, and asks how much accepted-source work it could ever affect before implementing a candidate.

## Entire original-Silesia accepted-archive census

With disposable, race-safe atomic counters in the *existing* `match.Expand` entry (no change to output or parser), compressed and decompressed **all 12 original user-supplied Silesia files**, **211,938,580 original B**, accepted native current JLZ, **58,767,739 B fully framed**, 12/12 byte- and SHA-exact after decoding. Corpus SHA-256 against committed manifest previously verified this same session. Counted all match commands from actually selected codecs, not just hypothetical math-source blocks or a simulated parser.

| Match-expansion category | Actual count | Actual reconstructed match bytes | Percent of all match reference bytes |
|---|---:|---:|---:|
| Exact match | **3,665,755** | **53,605,763 B** | **99.84563%** |
| Constant-byte XOR match | **572** | **11,176 B** | **0.02082%** |
| Constant-byte ADD match | **2,838** | **71,704 B** | **0.13355%** |
| **All reference expansions** | **3,669,165** | **53,688,643 B** | **100%** |

Only **3,410 / 3,669,165 = 0.09294%** of actual reference-expansion commands, and **82,880 / 53,688,643 = 0.15437%** of reference bytes, are transformed XOR/ADD. Three source files actually contain those operations: `mozilla` 31,423 B (XOR 3,875 + ADD 27,548), `ooffice` 25,508 B (1,425 + 24,083), and `samba` 25,949 B (5,876 + 20,073). The other nine have zero transformed expansions in the current chosen native archive. Exact matches account for 53.6 MB and 3.666 million commands. Existing baseline decoder profiling attributes ~46% of CPU to adaptive arithmetic, so eliminating the relatively rare transformed copy loop is unlikely to rescue JLZ's speed gap by itself.

**No inference of exact decode wall-time delta was measured here.** Transformation loops could have greater cost per byte than exact copy; the 0.15% byte fraction is a *workload opportunity indicator*, not a mathematical theorem bounding all CPU savings. Instrumentation itself uses atomics and can slow decompression. Do not compare this instrumented run's wall times with retained JLZ medians or claim a codec speed improvement.

## Decision and further Zstd work

Do not build a word-wide transformed XOR/ADD decoder solely for the current Silesia performance objective: its live workload is too small to promise a material gain, and the previous **encoder** wordwide check already failed a separate practical screening. This is a **workload-screen deferral**, not an implemented failed candidate and not a closure of other fast-copy methods. No ten-run result, no code promotion or new decoder subtype.

The **exact** match executor remains a potentially interesting Zstd-inspired target: 3,665,755 match commands were issued, typical short copies; compare actual call/disptach/buffer churn against the already bulk-optimized Go implementation first, accounting for past failed direct L/M/D column decoding before proposing any new fused executor. For higher-priority ratio opportunities, continue the Zstd row-based multi-candidate matching and shared-table superblock subtask list in [the Zstd architecture coverage ledger](../research/zstd-architecture-coverage.md).

**Validation:** instrumented test `TestZstdMathExpansionFrequency` completed all twelve entire original inputs; archives summed to accepted 58,767,739 B and all restored sources exact. Source path `match/expand.go` modified in disposable `jlz_zstd_expand` only; production source and `default.pgo` unchanged. Current accepted-source uncached Go suite and vet previously passed this session; no integrated candidate, corruption test, race-enabled full corpus or Windows runtime attempted for this profiling-only instrument.
