# Round 31: Three geometric boundary encoders (rejected)

**Date:** 2026-10-18. **Approved:** Round 29, unchanged. The user’s theoretical catalog drives this experiment. Only original Silesia files were used. All candidate code was isolated locally and is not part of the released JLZ compressor.

## Implemented, not estimated

Three distinct reversible Go encoder/decoder prototypes operated on the eight bitplanes of original bytes interpreted as a 2D Euclidean square grid:

1. **Dyadic region tree:** A square with uniform labels emits token 0 or 1; a mixed square emits token 2 and recursively encodes four quadrants. The decoder recreates each bitplane exactly. Serialized node count and native Huffman/RLE/raw packed node stream are included in payload size.
2. **Horizontal boundary distances:** Encode each row’s initial bit, the positive integer distances between label-changing edges, and row delimiters. Reconstruct by toggling bits at exact crossing positions. Header includes row width and decoded boundary-stream length, and two native-packed streams include all labels and distances.
3. **Vertical edge-incidence boundaries:** XOR each row with the previous row to identify two-dimensional vertical boundaries; encode horizontal crossings and row starts using variant 2, then integrate the vertical differences with reconstructed preceding rows. All width metadata and compressed streams are included.

This is a **bounded finite CAT(0)-inspired cubical-grid representation**, not an implementation, proof, or refutation of the full geometric CAT(0) boundary-reconstruction theory. All three had working exact inverses.

## Actual packed Silesia screen

The first, middle and last input regions for each of the twelve Silesia files total **7,728,404 bytes (36 regions)**. The prototypes were executed on these real source bytes; every reconstructed region matched exactly. Each result includes actual serialized model metadata and JLZ’s own entropy stream framing, but excludes complete JLZ1 file headers and digest. Native BWT and additive competitors were computed on exactly the same sample inputs. Neither Zstd nor XZ was run.

| Packed representation | Total across 36 regions |
|---|---:|
| Quadtree region tokens | 11,373,991 B |
| Horizontal boundary, width 128 | 6,580,943 B |
| Horizontal boundary, width 256 | 6,473,582 B |
| Vertical row-difference boundaries, width 64 | 6,984,395 B |
| Vertical, width 128 | 6,906,682 B |
| Vertical, width 256 | 6,876,873 B |
| Vertical, width 512 | 6,791,454 B |
| **Best geometric choice per region (optimistic oracle)** | **6,404,465 B** |
| Native BWT payload | 2,708,041 B |
| Native rank-one additive payload | 5,966,477 B |
| **Best BWT/additive choice per region (optimistic incumbent)** | **2,676,525 B** |

**All 36 geometry samples lose to the existing per-region BWT/additive minimum.** Even the optimistic geometry choice is **3,727,940 bytes (139.3%) larger** on the same input sample. This does not justify constructing a complete new JLZ archive or an incompatible decoder mode. There are no valid new full-archive candidate bytes or ten-run times.

## Verification and limits

- Exact original-byte roundtrip for all 36 Silesia blocks with three working encoder/decoder implementations; additional short, empty, partial, structured and high-transition inputs were used for correctness only.
- Negative corruption cases for truncated payloads, invalid row widths, impossible event counts and invalid quadtree headers passed.
- Uncached targeted Go tests, Go vet, and targeted race tests passed in the isolated disposable working copy.
- Round 29 remains **62,644,663 B** complete Silesia, accepted encode/decode median **7.516364 / 2.821363 seconds**. Historical unrerun Zstd-3 Silesia **66,216,569 B**, CLI compress/decompress **1.2561 / 0.8551 s**; XZ-6 **49,233,340 B**, CLI **116.9892 / 2.3079 s**. CLI and in-memory JLZ times are not directly comparable.

**Verdict:** Reject the **three finite geometric implementations**, not CAT(0) theory. Retain Round 29 and the active theoretical queue.

**Next distinct theory:** Try a materially different Tietze/grammar-presentation representation with full generator definitions, reconstructible parsing and properly charged metadata, rather than another minor entropy tuning.
