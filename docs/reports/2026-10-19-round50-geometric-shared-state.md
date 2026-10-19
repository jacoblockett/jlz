# Round 50 — geometric data layouts and shared mathematical metadata state (REJECT H1)

**Date:** 2026-10-19. **Production decision:** keep approved Round 44 unchanged, no experimental native stream modes promoted. This research was executed with disposable Go code based on the provided approved `b41ce26` source snapshot. Corpus is **only the original 12 complete Silesia files, 211,938,580 B**, canonical SHA-256 manifest **12/12 PASS**. No external compressor as backend and no hosted CI.

## H1: geometries of varying sizes improve selected native streams or compressor throughput

Implemented reversible byte layouts with explicit dest-to-source bijective maps and exact inverses:

- 2D tile transposes over 256/1,024/4,096/16,384-byte tiles and widths 4/16/64/256, as applicable;
- 2D Morton Z-order squares of side 8/16/32/64;
- 3D axis permutations and Morton space-filling traversal of cubes side 4/8/16.

**25 distinct layouts, all 676 archived opcode/length metadata streams from all 338 actually math-selected R44 blocks:** the original native packed size was **4,506,404 B**. Every single permuted native-packed stream had **the same byte length** as the unpermuted stream, so one paid mode byte per stream made each layout **4,507,080 B**, **0/676 wins, zero selective savings**. All permutations were checked to be exact bijections. This was expected from Huffman frequency invariance: permuting symbols does not change the zero-order frequencies, code lengths, or payload bit count. It can only add work to this entropy model; pure layout cannot accelerate zero-order Huffman by supplying new information.

**29 distinct reversible spatial residual tests on all 676 mathematical and 412 BWT-selected symbol streams:** tested modular subtraction at lags 1,2,4,8,16,32,64,128,256,512,1024; XOR at lags 1..256; and subtraction along nine representative 2D/3D rearranged scans. Every alternate fully packs its transformed bytes through JLZ's native coder and is charged a mode byte. For **math metadata**, actual incumbent **4,506,404 B**, optimistic best-per-stream among all candidates **4,465,636 B**, i.e. only **40,768 B selective savings from 94/676 streams**, before selector/preflight computation and complete-archive framing. **Every individual fixed mode grew the full collection**: even first-lag XOR was **5,090,885 B (+584,481)**. For **BWT** all 412 actual selected intermediate streams had **0/412 wins** even under per-stream oracle; incumbent and best remain **20,533,358 B**; all fixed spatial modes expand. All selected winning residuals underwent actual native pack → native unpack → inverse prediction exact reconstruction, not entropy estimates.

**Rearrangement *before* BWT** was independently tested on 18 original Silesia source regions (three distributed blocks for each of the six BWT-selected files), measuring actual full reversible `encodeBWTBijective` / `decodeBWTBijective` paths plus inverse geometric layout. Baseline **763,048 B**, encode **192.80 ms**, decode **81.08 ms** over those 18 regions. Seven layouts (2D transpose at 256/1024/4096, 2D Morton 16/32, 3D axis-8 and 3D Morton-8) each lost on **18/18**. Best layout *by size* was Morton 2D/16 with **1,603,016 B** (**+839,968 B**); fastest candidate encode was Morton 2D/16 at **259.61 ms** (+34.65% nominal in these same-session screens). Sample decode times are all >109 ms versus 81.08 ms baseline. These are **18-block screens only**, not whole-file/corpus timing medians. No geometric candidate warranted integration.

## H2: cross-channel/cross-block shared coding and decoder-visible predictive state

**Joint Huffman models across math metadata channels:** With all **338 actual selected opcode+length stream pairs** decoded from the original R44 archives, concatenating opcode then length and using one native model (plus original opcode byte count and a mode marker) yielded **5,522,728 B** versus independent actual/native **4,506,404 B**, **+1,016,324 B**, **0/338 pair wins**. Independent predictive distributions are valuable; naive joint histograms cannot encode cross-channel information without conditioning.

**Group-level shared Huffman model**, concatenating same-channel metadata across consecutive 2-, 4- and 8-block groups within each file, including per-part byte lengths and group marker:

| Group size (blocks) | Independent native total | Unconditional shared total | Unconditional growth | Best-of-shared-vs-independent oracle saving |
|---|---:|---:|---:|---:|
| 2 | 4,506,404 B | 4,593,908 B | +87,504 B | 8,026 B (125 winning plane-groups) |
| 4 | 4,506,404 B | 4,690,971 B | +184,567 B | 4,904 B (42 groups) |
| 8 | 4,506,404 B | 4,729,594 B | +223,190 B | 4,129 B (15 groups) |

An automatic group selector would incur *additional* identification/framing. These modest oracle savings are not integration-ready and sharing would reduce independent block decoding parallelism.

**Different, decoder-synchronized cross-channel predictive relationship:** opcode stream is already decoded, so use each operation type to select one of nine small last-successor state machines for the next associated length varint's **first exact byte**. A native reversible small-list transposition (LRU 1/2/4 last observed first bytes per opcode) maps learned frequent values toward low byte identities; untouched varint continuation bytes remain verbatim. Decoder derives the operation and updates states only from recovered length bytes; **no model table serialized**. This is not the rejected previous-byte model: conditioning is on the **other existing stream** (operation type), and the length varint structure is preserved. Exactly tested **229 aligned/unsplit length streams**; 109 split or inapplicable records were correctly excluded, never silently forced into an invalid one-length-per-op pairing. On the eligible subset, native packed incumbent **2,800,618 B**:

- One-slot state: **2,771,942 B**, -28,676 B (100 wins).
- Two-slot state: **2,690,909 B**, -109,709 B (117 wins).
- Four-slot state: **2,666,980 B**, -133,638 B (126 wins).
- Best-per-stream oracle across 1/2/4 modes: **2,635,005 B**, **-165,613 B**, 128 wins. This is a *component-only upper bound* even before whole-archive subtype choice.
- Carrying exact four-slot state across consecutive decoded blocks within each four-region independently decodable group instead of resetting per stream cost **2,667,132 B**, *152 B worse* than fresh four-slot states. No measurable group benefit.
- All tested transforms passed actual encoded native inverse and exact source-stream reconstruction; independent random/patterned and multi-byte-varint fixtures pass.

Even the **-165,613 B oracle** is substantially below historical rejected complete conditional-role models that achieved -670,853 B but degraded practical speed. New opcode-state code requires an additional encoding/prediction and decoder path; the observed size opportunity does **not** justify another experimental on-wire mode or ten expensive complete-Silesia promotions.

## Validation and historic comparisons

- Canonical twelve-source SHA-256s **12/12 PASS**.
- Exact reversible layout, spatial predictor, and opcode-conditioned varint-state unit tests PASS. Every selected gain was confirmed with true packed bytes and actual inverse.
- `go test -count=1 ./...` PASS, `go vet ./...` PASS, targeted `go test -race -count=1 ./engine -run 'Test(GeometryExact|GeoPredictExact|OpStateExact)$'` PASS, Windows/amd64 cross-build PASS (not a native Windows runtime execution).
- All source alterations were disposable **research tests only**, added no production decoder or codec and did not modify the accepted on-wire format. No accepted JLZ, Zstd or XZ measurements were rerun.

| Model | Complete Silesia archive | Original reduction | Encode s | Decode s |
|---|---:|---:|---:|---:|
| Zstd-3, historical external CLI | 66,216,569 B | 68.76% | 1.2561 | 0.8551 |
| XZ-6, historical external CLI | 49,233,340 B | 76.77% | 116.9892 | 2.3079 |
| Current approved JLZ Round 44, retained ten-run median | **61,418,211 B** | **71.02075%** | **5.945576** | **2.445270** |
| Candidate JLZ Round 50 | **N/A: no viable full-archive candidate** | N/A | N/A | N/A |

| Candidate-vs-current size % | Candidate-vs-current encode % | Candidate-vs-current decode % |
|---:|---:|---:|
| N/A, component screens only | N/A, no full candidate | N/A, no full candidate |

**REJECT H1** for these specific geometry and simple cross-stream state-sharing implementations. Preserve open broader stateful geometry and cross-channel theories; do **not** infer impossibility from these finite negatives. Key structural conclusion: pure permutation cannot help the current frequency-only code, while spatial differences ruin already sparse BWT MTF symbols. Existing opcode semantics do supply real conditional prediction for length bytes, but too little gain for another slower format. Next pursue **low-overhead, one-pass structural metadata prediction with cheap preflight or directly optimize the selected mathematical matcher/BWT bottlenecks**, rather than more parameterized grid scans.