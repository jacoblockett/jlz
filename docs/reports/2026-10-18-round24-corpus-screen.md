# Round 24: heterogeneous BWT generalization screen (no codec promotion)

**Date:** 2026-10-18. **Approved compressor remains Round 23.** All permanent changes to `main` are CI caching and documentation, not the JLZ encoder or decoder.

## Hypothesis

The automatically chosen BWT/MTF/zero-run/Huffman representation accepted in Round 21 and broadened in Round 23 might generalize beyond the 12 Silesia files. Test the current selector and a **forced native BWT-per-region diagnostic** against it on a varied set of original bytes, preserving the complete `JLZ1` archive format and byte-exact restoration. This is not a recommendation to always use BWT, and no external codec is used.

## Data and methodology

Eight locally installed real-world inputs (total **70,159,124 source bytes**) with source types spanning TAR-packed Go/Python/C headers, font TTC, Java JAR, SQLite database, ELF shared library and audio SoundFont. Source/header TARs contain unmodified installed-file bytes with ordinary TAR headers; the TAR containers were created for this screen, so their layout is an experimental choice. System assets are used strictly locally, never redistributed or checked into Git. Note that this is a small convenience sample, not a statistical representation of all file types. Source hashes are supplied below for reproducibility when the same installed bytes can be obtained.

From an isolated local checkout of the exact Round 23 `region_first.go` and `bwt_native.go` Git blob SHA identities, run `compressFast` and decode the resulting archive, then construct a complete all-BWT `modeMixedCross` archive with the existing `encodeBWTNative`, `encodeParallelBWTRegions` and standard `decompress` functions. **All 8/8 normal archives and 8/8 forced-BWT archives roundtrip byte-exactly.** Both sizes count complete framing and SHA-256, not payload-only estimates. Go 1.23.2 on linux/amd64. No encoding source changes. One run per file; times are **screening-only**, not ten-run medians, and the alternatives were not interleaved as paired 10-run benchmarks.

## Complete archive-size results

| Real-world input | Original bytes | Round 23 automatic JLZ | Forced native BWT | BWT minus automatic |
|---|---:|---:|---:|---:|
| Java `batik-all-1.18.jar` | 4,328,384 | 3,882,346 | 3,936,425 | +54,079 |
| C/C++ headers TAR | 8,796,160 | 942,848 | 1,033,378 | +90,530 |
| Noto Sans CJK TTC | 19,484,784 | 15,027,086 | 15,746,076 | +718,990 |
| Go standard-library source TAR | 8,570,880 | 1,396,591 | 1,442,792 | +46,201 |
| PHP module ELF | 5,980,008 | 2,267,179 | 2,380,513 | +113,334 |
| PROJ `proj.db` SQLite | 9,359,360 | 1,563,065 | 1,664,011 | +100,946 |
| Python standard-library source TAR | 7,669,760 | 1,470,970 | 1,521,957 | +50,987 |
| MIDI SoundFont | 5,969,788 | 5,766,534 | **5,670,356** | **-96,178** |
| **Total** | **70,159,124** | **32,316,619** | **33,395,508** | **+1,078,889 (3.34% larger)** |

The Round 23 automatic selector chose mathematical or native representations for all eight files; none of these output archives contained BWT block mode 12. The one smaller forced-BWT file, the SoundFont, took **0.3984 s** to encode via forced BWT versus **0.1040 s** for the automatic codec in two independent one-run screens on the same host, roughly **3.8 times** as long. The delta saves 1.67% size only on this file. Timings are illustrative and too few for a promoted CPU-cost claim.

These eight inputs have no independently measured Zstd/XZ columns. Historical Silesia reference archive sizes, **Zstd-3 66,216,569 B and XZ-6 49,233,340 B**, concern an entirely different corpus and are not comparison denominators for the table above.

### Source SHA-256 provenance

```text
01ed478753fa9c1aa11abdb6757608968fb5779d0843922659b6831d1c37a9f8  c-headers.tar
e5bc198326bd5b49a82b6491448774076140f7f31036f320b16de3ed4b7421d5  go-source.tar
6101eb3503baaadb5924fb7bc83ffdfd168b1d35508d2c1f8e1f415207beff9f  python-source.tar
297c1b91e86addaf35297263c9b3b36cb970b632667702fde6cee090d21b2711  batik-jar
b76b0433203017ca80401b2ee0dd69350349871c4b19d504c34dbdd80541690a  font-ttc
7da07d207d2ba2c28b55e11f5754a1a50523e35a5f29b624b22588fe162db70f  php-elf
8bbb2ea3a88c44383f51e205755c32bff09b8aff7937eca001a1cd519714be60  proj-sqlite
c5378b62028c920cb11e4803327983fee2f2cdff5dc89c708e39da417e51c854  soundfont
```

## Validation and compatibility

The reconstructed checkout successfully ran `go test -count=1 ./...`, `go vet ./...`, and a complete Round 23 12-file Silesia test that exactly reproduced the approved **63,658,799 B** archive total on supplied source files. No new compatibility modes exist. The Silesia historical Round 21 total discrepancy (63,770,766 B historically vs 64,505,708 B for a recreated Round 21 on current source bytes) remains unresolved: the historical input-byte hashes were not retained, so its cause cannot be identified from available evidence. It must not be hidden or attributed to any particular cause without new evidence.

A disposable test used to inspect chosen archive block modes found the automatic output uniformly bypassed BWT on these particular eight samples; it was not kept as production test code.

## Verdict and next investigation

**REJECT unconditional BWT expansion. Retain Round 23 as the approved compression baseline.** The existing multi-representation selector is materially better than forced BWT for this heterogeneous convenience sample, and the SoundFont exception has unfavorable one-run CPU economics. This result does not rule out better gated, per-file or per-block selection under an appropriately low-cost diagnostic.

**Next priority:** move beyond gate-threshold tuning. Investigate a substantially different low-overhead structural representation or a faster native suffix sort only if its cost/benefit can be established. First compare existing entropy/match accounting and the historical `docs/research/attempted-ideas.md` inventory to avoid repeating dead ends. Seek broader licensed inputs and keep complete original-byte archives in all comparisons.
