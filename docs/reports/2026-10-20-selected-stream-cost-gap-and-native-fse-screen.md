# Selected-stream coding-cost census, CPU profile, and plain native-FSE math-literal feasibility screen

**Date:** 2026-10-20. **Status:** PROFILED/SCREENED; native-FSE supplementary literal placement **REJECTED**. **No new codec, integrated candidate, ten-run result, production Go edits, decoder change or PGO change.**

## Purpose and novelty gate

After a series of rejected ZPAQ ICM/ISSE and PPM models, first locate **actual currently selected archive coding costs** rather than merely screen another substitute on unselected source blocks. The 2026-10-20 [bytewise context entropy screens](2026-10-20-bytewise-context-entropy-screens.md) already measured ~47.56% decoder CPU in the binary order-1 arithmetic decoder; this CPU hotspot is **not a new discovery**. The added evidence here is a complete current-selected-mode census of all 12 originals, actual packed literal entropy allocations and a new **single-model, current native finite-state coding placement on *selected* match-literal streams**. This is not a new entropy algorithm: JLZ already uses single-state 12-bit rANS subtype5 on BWT and opcode/length streams. Earlier mathematically conditioned static 8-model byte/rANS and 8-prefix-bit rANS placements were tested and rejected; unlike them, this test reuses **one current exact native model/decoder**, with no novel context tables and no new subtype.

The specific question: is it possible to bypass expensive existing adaptive per-bit arithmetic and still retain most of the **real** current selected literal savings, at a negligible cost of introducing native-FSE selection into `packLiteralStream`? Do not confuse the measured stream bytes or decompression CPU fractions with an integrated candidate benchmark.

## Provenance and exact accepted archive

Tested source from user-supplied original `jlz-6a388ba91682.zip`, verified earlier as matching accepted production at `6a388ba91682bb6fd7f0240e2e59406c9a6ce73d`. Subsequent live `main` updates through the immediately preceding report head `734634f85c32f30290520701b75e7da1eaed2c85` changed research documentation only, not relevant engine source. All **12 complete original Silesia source files, 211,938,580 B**, matched the canonical committed [SHA-256 manifest](2026-10-18-round23-silesia-sha256.txt) individually. Root `default.pgo` SHA-256 **3b35cda74832360a252b63c1805e025aa1f4ed7a4fe4124d9a55929a6e86289f**, unchanged; Go 1.23.2 linux/amd64. All diagnostic source existed in an **uncommitted disposable working copy only**. No external compressor was used as a backend, and Zstd/XZ were **not rerun**.

An entire accepted JLZ encode/decode pass of all original files reproduced exact **58,767,739 B**, reconstructing **12/12 SHA-identical original bytes**. Under profiler overhead the cumulative observed encoding wall time was about **5.24s**, decoding about **2.05s**. **These are deliberately not a new ten-run baseline**, not speed comparisons or qualification results.

## Current accepted representation: ALL 12 files and ALL 815 independent source blocks

Parsed the **actual complete framed accepted archives**, skipping their existing headers/varints and counting real selected block modes and nested packed entropy frames. There is no guessing based on a potential-mode oracle. All per-mode stored totals below exclude individual outer block mode/varints and top-level archive headers, which total **6,141 B** across all 12 archives.

| Accepted region representation | Blocks | Stored payload bytes | Share of full JLZ archive |
| --- | ---: | ---: | ---: |
| Native BWT selected modes15/16, plus BCJ→BWT mode19 | 471 | **25,657,006** | **43.66%** |
| Native mathematical matched/reference modes7/9/10/11 | 279 | **23,661,860** | **40.26%** |
| Native field-aligned byte lanes mode13 | 32 | **5,176,017** | **8.81%** |
| Native additive-word mode14 | 33 | **4,266,715** | **7.26%** |
| Outer block framing and file headers | | **6,141** | **0.01%** |
| **Total actual JLZ archive** | **815** | **58,767,739** | **100%** |

**Selected-mode distribution by source:** `mozilla` 192 math-selected blocks (modes7/9/10/11) plus four lane blocks; `samba` 83 math-selected blocks; `ooffice` four math-selected blocks and 20 BCJ→BWT; `nci` 128 BWT blocks, `webster` 159 BWT, `osdb` 39 BWT, `sao` 28 28-byte-lane blocks, `x-ray` 33 native additive blocks. The balance of Dickens/MR/Reymont/XML is BWT-selected. Counts and bytes are original archive reality, unlike prior screens of arbitrary possible representations.

**Actual nested framed entropy-stream payload census:** BWT current FSE codec5 **652 independent streams**, **25,534,853 B** entropy bodies (including model tables), only seven BWT canonical Huffman codec2 streams totaling 116,169 B. Mathematical *literal* stream totals across current selected math regions: **279 streams, 18,998,973 input literal bytes → 15,483,417 completely framed stored bytes**, a saving of **3,515,556 B** against raw storage. Of those, **262 adaptive arithmetic codec3** streams contain **15,193,831 literal bytes → 11,678,194 framed archive bytes**, saving **3,515,637 B** themselves. Other math literals are raw/Huffman and contribute the tiny difference. Thus the arithmetic-coded selected math literals alone occupy **19.87%** of the entire accepted archive.

Other native math stream coding sizes (codec/frame bodies, excluding stream header): distance fields **3,832,758 B** split between raw, Huffman and logarithmic slot coding; opcode FSE **1,052,366 B**; length streams use existing FSE and Huffman. Their presence confirms that rediscovering split op/length/distance columns is not a novel architecture.

## CPU profile: identification of throughput costs, not speed claims

Ran Go runtime CPU profiling for an accepted all-twelve corpus full encode pass and a complete exact decode pass, separately. CPU samples cover **all parallel goroutines**, so CPU-seconds exceed elapsed wall time and cumulative call-tree percentages overlap. Sampled encoder CPU total **17.50s**, decoder CPU total **6.78s**. Relevant *cumulative* sampled CPU:

| Function stage | Encoder cumulative CPU | Decoder cumulative CPU | Interpretation |
| --- | ---: | ---: | --- |
| Native BWT.Transform, including Go SA-IS | **6.03s, 34.46%** | n/a | Main forward BWT source transform still costly; suffixarray serialization tweaks have ALREADY been tested/rejected in ten-run results |
| Native mathematical token discovery `tokenizeMatchV2RangeStride` | **4.32s, 24.69%** | n/a | Conventional parser/anchor/pricing variants ALREADY broadly investigated; no implicit novelty |
| Adaptive order-1 arithmetic `entropy.EncodeArithmetic` / `DecodeArithmetic` | **3.61s, 20.63%** | **3.09s, 45.58%** | Expensive eight serial prediction/update steps per decoded byte |
| BWT.Inverse | n/a | **1.55s, 22.86%** | Two-hop LF was previously accepted; 4-hop previously tested and rejected |
| BWT bijective-zero decode | n/a | **0.75s, 11.06%** | Different from BWT inverse; no universal alternative is implied |
| Native finite-state decoder | n/a | **0.42s, 6.19%** | Already comparatively cheap |

The `mozilla` file alone needed **1.840s of 5.235s elapsed encoding** (~35%) and **0.827s of 2.053s elapsed decoding** (~40%) in this one profiled pass; these percentages include profiling overhead and concurrent scheduling. No per-file ten-run speeds are claimed.

## Distinguish insufficient modeling from coder overhead

For each actually selected adaptive-arithmetic literal stream, reversed its existing codec3 payload *using JLZ's own lossless decoder*, computed the byte zero-order Shannon coding lower bound (`H0 = sum(count[s] * log2(N/count[s])) / 8`) on the exact original stream, and compared with its actual complete framed output. The H0 figure is **only a zero-order ideal**, not a universal compression limit and not a promised achievable new coder size.

| Actually selected arithmetic literals | Input bytes | Existing complete framed bytes | H0 ideal byte count | Actual arithmetic saving beyond H0 |
| --- | ---: | ---: | ---: | ---: |
| `mozilla` 179 streams | 12,390,914 | 9,736,327 | 10,777,288 | **1,040,961 B** |
| `ooffice` four streams | 286,468 | 218,358 | 244,042 | **25,684 B** |
| `samba` 79 streams | 2,516,449 | 1,723,509 | 1,896,502 | **172,993 B** |
| **All 262 arithmetic streams** | **15,193,831** | **11,678,194** | **12,917,832** | **1,239,638 B** |

This is decisive evidence of **useful selected-codec contextual conditioning**: just replacing these streams with independently modeled **order-zero** coders forfeits contextual gain even assuming a perfect zero-order code with zero model headers. Conversely, because these context gains cost serial per-bit model estimation and updating, the next viable architecture must retain the predictive information while making reconstruction much cheaper.

## Distinct supplementary placement: use existing single-model FSE on selected math literals

Screened the already-present **`packFiniteState`** (native single-state 12-bit rANS/Huffman/raw automatic choice) on the **raw losslessly decoded contents of all 279 actually chosen mathematical literal streams**, comparing its complete candidate frame to the existing complete framed literal archive bytes. Every proposed winning stream was round-tripped through **unchanged `unpackStream` decoder subtype5** and checked byte-for-byte. Thus a true optional integration would **not** require a new subtype or break older archive decoding, but it still needs another full candidate-model/stream encoding to select winners.

| Source | Actually selected math literal streams | Favorable FSE alternatives | Oracle fully framed bytes saved |
| --- | ---: | ---: | ---: |
| Mozilla | 192 | **0** | **0** |
| Ooffice | 4 | **0** | **0** |
| Samba | 83 | **4** | **1,327 B** |
| **Total** | **279** | **4** | **1,327 B** |

The complete evaluation spent **0.2222 seconds additional wall time** in the native `packFiniteState` call on these streams in this diagnostic pass. This is **component-only** work, not an observed full-archive encode regression, but the favorable oracle is just **1,327/58,767,739 = 0.00226% of compressed bytes**, before any conditional model selection code. Rejected alternative frames on the 275 losing streams would total **1,354,952 B** larger if falsely forced. **REJECT immediate integrated implementation** of optional ordinary FSE on selected math literal streams: extra encoder work is not justified; do not iterate table size, per-stream cutoff, Huffman threshold, normalization, or state count.

This result also **independently supports** previous negative conditional-rANS models, without treating them as the same test: those added costly per-context frequency models whereas this uses one already available native decoder subtype.

## Next research gate and limits

The remaining material bottleneck is **serial decoder-synchronous prediction on a relatively small number of presently selected mathematical literal streams**, especially Mozilla and Samba; another crude static entropy replacement is not the answer. Earlier attempts on *this same domain* that MUST be respected: exact reciprocal LUT, binary EMA probability redesign, previous-byte high-bit conditioned rANS, binary-prefix static rANS, bytewise adaptive Fenwick, byte-renormalized arithmetic and full ZPAQ-style ICM (ICM won primarily on BWT, not match literals). A successor requires an actually distinct **predictive information representation or amortization architecture**, with explicit data-source constraints, table/framing costs and decoder parallelism, and a quick stage-level size-plus-work feasibility proof *before* any production implementation. **Do not reflexively revisit these earlier algorithms or adjust their thresholds.** The BWT ratio opportunities are material in output bytes but already heavily researched, including 886KB ICM savings with near-double full decoder time. More upstream source comparison may be needed before the next concrete implementation.

**Validation:** `go test -count=1 ./...` PASS on disposable accepted-code workspace with temporary diagnostic tests; `go vet ./...` PASS. Real chosen all-file archive bytes reconstructed, every source file matched committed 12-file SHA-256 manifest. No integrated candidate, no format changes, and therefore no need for new ten-run timings, backward archive validation, fuzz/corrupt mode tests, targeted race tests or Windows runtime test. Historical refs retained without reruns: accepted **58,767,739 B / 4.7873709325s encode / 1.8721679705s decode**, Zstd-3 **66,216,569 B / 1.2561/0.8551s**, XZ-6 **49,233,340 B / 116.9892/2.3079s**.
