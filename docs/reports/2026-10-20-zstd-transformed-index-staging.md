# Zstd-style mathematical match-index work allocation: exact-only and sparse XOR/ADD insertion/probing

**Date:** 2026-10-20. **Status:** source-researched, multiple concrete integrated **one-pass screens**, plus independently order-balanced whole-corpus paired timings. **No ten-run candidate qualified or production promotion.** Current root PGO, two-command automatic codec, archive format, backward decoding and Go source unchanged. This is one part of the [active Zstd-only architecture matrix](../research/zstd-architecture-coverage.md); the family is NOT claimed universally exhausted.

## Upstream source and novelty relative to older JLZ failures

Studied actual Meta Zstd [`zstd_double_fast.c`](https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_double_fast.c) (Git blob `1a266e7d955f7533640f738311d1711d996f0b27`), [`zstd_fast.c`](https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_fast.c) and [`zstd_compress_internal.h`](https://github.com/facebook/zstd/blob/dev/lib/compress/zstd_compress_internal.h). Its fast parser separates cheap short/long recency hash probes and adapts insertion/advancement rather than continuously maintaining every candidate class. Zstd **does not have JLZ's constant-byte XOR/ADD reference grammar**; this research adapts the principle of work allocation, NOT Zstd source code or Zstd backends.

**Nearest previous rejection:** JLZ's dual short/long hash **exact-only replacement parser** was integrated and rejected on ten complete runs, producing 61,866,415 B versus then-current 60,838,336 B (+1,028,079 B) while nominally improving encode speed. This new trial retains the **accepted JLZ exact-match hash/index, exact/XOR/ADD token grammar, repeat distances, match scorer and lazy parsing** unchanged and manipulates *only* maintenance and availability of the two transformed candidate hash indexes. The prior result does **not** predict the size penalty of holding the accepted exact parser invariant. Current selected decoder workload confirms XOR/ADD match expansions total only **3,410 commands / 82,880 B** out of 53,688,643 B reconstructed through all reference types ([command census](2026-10-20-zstd-sequence-execution-workload-screen.md)); that is a motivation for cheap transformed-work gating, not proof that discarded matches are harmless.

**Mechanism tested:** In the existing `engine/match.go` `tokenizeMatchV2RangeStride`, current source instantiates `newMatchIndexTransforms(input, !exactOnly && likelyTransformedMatches(input))`; when active, each inserted position populates exact/XOR/ADD recent-fingerprint entries, and each candidate search probes all active indexes. Original-Go disposable variants:
1. **Current same-parser exact-only control:** instantiate the unchanged exact index with transformed indexes disabled; all downstream parse choices unchanged, unlike prior dual-hash exact parser replacement.
2. **Sparse transformed insertion stride4:** always insert exact index at *every* eligible position; only insert XOR/ADD at positions congruent to 0 modulo 4, retaining normal transformed candidate queries. This can change matches but retains exact-index parity and live decoder grammar.
3. **Sparse insertion stride8:** same as stride4, period8 (a mechanistic sampling probe, NOT evidence an infinite stride-parameter optimization was exhausted).
4. **Staged loop implementation of sparse insertion8:** precompute `types`=1 or3 before the per-index loop, reducing repeated conditional checks; exact archive bytes unchanged relative to candidate3.
5. **Sparse transformed queries stride8, unchanged full insertion:** probe exact every candidate point and transformed references only at congruent positions.
6. **Both insertion and query sparse stride8:** suppress transformed index insertion and transformed lookups, always retaining exact.
All are supplementary scheduling variants of the accepted native parser, do not run two full parsers, and do not introduce a decoder mode or metadata.

## Entire original 12-file Silesia correctness screens

Used all twelve **unmodified original** Silesia inputs (211,938,580 B), SHA-256 matched committed source manifest, native Go encoder and decoder, original root `default.pgo` retained, source snapshot verified against live main (production Git blob `engine/match.go` `6c9714114a056bf49f8124f81826f0eded68a5d8`). **Every reported size is a fully framed archive summed across all twelve original files**, every one-pass variant encoded and reconstructed all twelve byte-exact. The initial encode/decode wall times below are **one-pass screens** and are not ten-run medians or directly comparable with historical accepted wall times.

| Implementation | Full archive B | Delta B vs accepted | One-pass encode s | One-pass decode s | Decision scope |
|---|---:|---:|---:|---:|---|
| Accepted JLZ, historical reference | **58,767,739** | anchor | **4.7873709325** (ten-run historical) | **1.8721679705** (ten-run historical) | current unchanged |
| **Same parser, exact-only** | **58,789,865** | **+22,126 B (worse)** | 5.417 (single screen) | 2.018 (single screen) | size-for-speed tradeoff, not promoted |
| Sparse transformed insertion / 4 | **58,766,452** | **-1,287 B (better)** | 5.541958 | 2.123919 | negligible size win; no speed evidence |
| Sparse transformed insertion / 8 | **58,764,729** | **-3,010 B (better)** | 5.433392 | 2.046764 | tiny size win; slower in paired tests below |
| Sparse transformed probing / 8, full insertion | **58,770,768** | **+3,029 B (worse)** | 5.458411 | 2.078050 | screened / unpromising |
| Sparse insertion + probing / 8 | **58,773,026** | **+5,287 B (worse)** | 5.888648 | 2.162904 | screened / unpromising |

**Exact-only per-file complete archives** against accepted (all other nine files bit-size unchanged): Mozilla **18,452,630→18,472,075 (+19,445 B worse)**, ooffice **2,692,252→2,699,583 (+7,331 B worse)**, Samba **5,134,262→5,129,612 (-4,650 B better)**. This surprising Samba gain is a signal that transformed match competition can perturb parser choices unfavorably, not that all transformed modes should be removed.

Sparse insertion/8's 3,010-byte improvement is **0.00512% of accepted archive bytes**, or ~0.00142 percentage points of original-source compression reduction. There are no new decoder modes, so any decode timing variation across these trials is almost entirely attributable to different command/literal workload and host noise, not a new faster decoder architecture.

## Controlled same-session whole-original-corpus paired comparisons

To distinguish real local encode throughput from the stale ten-run historical baseline, used disposable test programs that, within a single Go process and with explicit `-pgo=./default.pgo`, encoded each of the **same twelve actual originals** once with incumbent and once with modified same-parser mode; completed paired order permutations (not a new ten-run candidate or accepted reference). All outputs exact; decode timing in the sparse runs also measured independently.

**Exact-only vs current, two completed whole-corpus pairs (third attempted process reached execution timeout, excluded):**

| Completed pair | Current encode | Exact-only encode | Exact-only size | Outcome |
|---|---:|---:|---:|---|
| Pair A | **5.621423678 s** | **5.249979837 s** | 58,789,865 B | ~6.6% faster encode, +22,126 B worse size |
| Pair B | **5.485406318 s** | **5.350677636 s** | 58,789,865 B | ~2.5% faster encode, +22,126 B worse size |

**Sparse insertion/8 vs current, two complete whole-corpus pairs:**

| Completed pair | Current encode | Sparse insert/8 encode | Current decode | Sparse insert/8 decode | Outcome |
|---|---:|---:|---:|---:|---|
| Pair A | **5.295401484 s** | **5.424317543 s** | 1.935500119 s | 1.920931652 s | sparse ~2.4% slower encode, 3,010 B smaller |
| Pair B, hoisted per-position type-loop branch | **5.185729909 s** | **5.362159066 s** | 1.938846621 s | 1.934699224 s | sparse ~3.4% slower encode, 3,010 B smaller |

**Interpretation:** the unconditional acceptance of even a small modulo and cursor divergence into the existing tight matcher loop can exceed the saved transformed-index writes, and could alter hash and cache locality. These two non-randomized, short paired screens are insufficient for precise causal CPU attribution or a universal performance claim. They do establish that the **current concrete sparse insertion loop has not recovered throughput**. Do not report an accepted +3KB Pareto win or invent ten-run results.

The exact-only control **did** avoid transformed hash maintenance/probes and reduced same-session encode time in both complete pairs, albeit with a mere **+0.03765% archived size regression**. No user approval has been obtained for trading a small ratio regression against moderate speed gains, and a broader-than-two-pair statistical benchmark would still be required before proposing it as a new production policy.

## Decision, scope and next Zstd-only work

- **REJECT concrete sparse transformed insertion/4 and insertion/8 as performance gains for now**, despite very small real ratio improvements, because paired whole-corpus timing is *worse* and additional modulo/branch work is substantial relative to saved model writes. The hoisted-loop refinement was specifically tested and did not change the conclusion.
- **REJECT sparse transformed probing/8 and both insertion+probing/8 as direct candidates**, since even their full archive sizes regress and their one-pass encode results do not promise clear speed improvement.
- **DEFER, NOT conclusively reject, exact-only same-parser policy**: both controlled pairs are faster, but archive expands 22,126 bytes overall, with per-file mixed behavior. This is a potentially user-approvable speed/ratio tradeoff if later paired and candidate ten-run evidence confirms it; need not hardcode a small gain just to claim progress.
- **DO NOT claim the Zstd-inspired transformed-index staging family exhausted.** Untested structurally different strategies include conditional transformed indexing based on a cheap, actual *quantitative evidence* gate rather than the current one-hit binary gate, row/tag cache-line packed search for all three transformations, or delayed-index construction triggered only by matched evidence. These would require a genuinely different algorithm/work bound, not another period/threshold sweep.
- No permanent mode, no root PGO modification, all experimental files disposable; no integrated ten-run qualified variant and therefore no new archive compatibility, race/Windows candidate acceptance runs needed. Accepted source uncached tests and vet already passed in this research session.
- **Next Zstd work:** interrogate Zstd's row-based cache-local matchfinder/skip policy source with realistic current-JLZ exact/XOR/ADD candidate frequency and profiler cost; or a *math sequence-superblock* table-sharing architecture that differs from the failed BWT table reuse. Preserve the Zstd coverage matrix as the authoritative gate.

**Fixed historical refs:** original 211,938,580 B; accepted JLZ **58,767,739 B / 4.7873709325s / 1.8721679705s**, Zstd-3 **66,216,569 B / 1.2561 / 0.8551s**, XZ-6 **49,233,340 B / 116.9892 / 2.3079s**. No third-party codec was used as a JLZ backend.
