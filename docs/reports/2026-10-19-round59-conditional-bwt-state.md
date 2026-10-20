# Round 59 — Decoder-synchronized BWT conditional entropy (REJECT H1 for production, positive compression evidence)

**Date:** 2026-10-19. **Production baseline remains Round57**, unmodified. Research-only native Go prototypes were built and executed locally against exclusively the **12 complete original Silesia files**, total **211,938,580 B**, with **12/12 canonical input SHA-256 PASS**. No external compression backend, imported Zstd/FSE implementation, branches, GitHub Actions, CLI flags or production format changes were made.

## Hypothesis

Round57's cheap minimum-description-length (MDL) BWT model partitioning saves archive bytes with negligible extra CPU, while Round58's dynamic split-point search saved more but cost ~10% encoding speed for 0.34% archive savings. Investigate **one-pass decoder-synchronized statistical context** over BWT's *existing bijective zero-run / move-to-front rank bytes*. A context can be selected using the previously **already decoded** symbol. Therefore no per-byte role metadata, segment boundaries, or cross-block model state is required. This is a genuinely different model from Round48's static role maps and Round58's match-length-conditioned metadata **stream fragmentation**.

Built original native three-context interleaved byte-rANS with 12-bit normalized tables, using existing Round53 finite-state mathematical state equations. Context rule: preceding rank-code byte 0–1, 2–5, or >=6; initial preceding byte 0. Each block serializes **three normalized independent sparse tables**, a decoded symbol count, and one interleaved rANS state stream; decoder chooses the exact table from the previously restored rank code. A new **experimental BWT block mode17** wraps the whole block's existing BWT primary index and those entropy bytes. Decoder validates table frequencies, initial/final state, symbol count and exact BWT reconstruction. Other native old BWT, mathematical match and fallback paths remain intact. Tested only as an alternate BWT model and only selected when its **actual, completely serialized payload is smaller**.

## Stage 1: All selected BWT-stream component evaluation (451 blocks, actual original corpus files)

On *all 451 actual BWT source blocks* from seven BWT-selected Silesia files (dickens, mr, nci, osdb, reymont, webster, xml), independently generated native BWT+bijective symbol streams and tried seven deterministic decoder-reconstructible predecessor rank-code context rules. All model tables, encoded state and stream-length framing were charged and **all roundtrips exact**. The standalone one-table entropy streams measured **23,416,721 B**. The best simple fixed context, three ranges of *previous rank-code* [0–1,2–5,6+], cost **23,125,864 B**, saving **290,857 B relative to this component-only single-table baseline**, with 440/451 individual block wins. Other tested predecessor rules (previous-zero/nonzero binary; preceding exact 0/1/other; previous two bytes' zero state, and a 4-context threshold rule) differed in size, with the 4-context rule inapplicable to blocks where a context never appeared. This 290,857 B is **not additive to Round57's already accepted 2-/4-neighborhood savings**.

## Stage 2: Real complete JLZ integration and improvement in selection

Reconstructed the exact approved Round57 Go encoder from the original complete repository source plus changes, and **first confirmed its entire original twelve-file Silesia aggregate was precisely the historical 60,585,717 B**, with 12/12 exact decode. Added research-only mode17 and ran these full-archive configurations:

| Complete 12-file prototype, one entire screen | Archive bytes | Encode time | Decode time |
|---|---:|---:|---:|
| Original R57 reconstructed, correctness-only screen | 60,585,717 | 6.192541 s | 2.428193 s |
| Round57 existing 2-/4-regime winner **plus optional full context encode** | **60,391,212** | 6.488696 s | 2.522852 s |
| **One-choice MDL candidate**: estimate context and region costs, encode only the predicted better alternative (still compare with old fully packed stream) | **60,400,234** | 5.840482 s | 2.242980 s |

The initial full-size fallback generated both optional alternatives; one-choice MDL predicted cost from actual context symbol counts plus conservative sparse-table headers and encoded **at most one** of segmented versus contextual alternatives per eligible BWT block. This gave up only **9,022 B** of the extra compression win while avoiding substantial duplicate work. Neither screen result is a ten-run median. Both candidates restored all twelve original files byte-exactly, and the optional addition never replaced existing mathematical match or BWT preprocessing.

## Ten-run COMPLETE Silesia test: one-choice MDL candidate

Compiled a single local candidate test binary and ran **ten** complete original 12-file Silesia encode/decode passes; each pass verified every reconstructed input byte and embedded digest, **120/120 exact file reconstructions**, and every pass emitted **60,400,234 B**. A long wrapper timed out after three completed runs; the remaining seven runs were completed in short independent batches using the same compiled binary; no partial run counted.

| Full candidate run | Archive B | Encode s | Decode s |
|---:|---:|---:|---:|
| 1 | 60,400,234 | 6.144602090 | 2.309428131 |
| 2 | 60,400,234 | 6.149488874 | 2.526825133 |
| 3 | 60,400,234 | 6.947602985 | 2.566481318 |
| 4 | 60,400,234 | 6.734741214 | 2.482780120 |
| 5 | 60,400,234 | 6.040984885 | 2.349595762 |
| 6 | 60,400,234 | 5.848225577 | 2.347884461 |
| 7 | 60,400,234 | 5.824531018 | 2.233203311 |
| 8 | 60,400,234 | 5.978526333 | 2.319172596 |
| 9 | 60,400,234 | 6.091595842 | 2.322116349 |
| 10 | 60,400,234 | 6.365574080 | 2.337119089 |
| **Median** | **60,400,234** | **6.118098966** | **2.342501775** |

Retained approved Round57 independent ten-run historical medians: **60,585,717 B / 5.493743735 s encode / 2.194427893 s decode**, **not rerun for speed**. The candidate saves exactly **185,483 B (-0.3061497% of prior compressed bytes)** with nominal **+11.36484% encode / +6.74772% decode** runtime. Time differences are **unpaired** and should not be treated as exact causal estimates. The size saving is real but clearly disproportionate to runtime and extra decoder format/state complexity.

### Required comparison: complete original Silesia

| Model | Full compressed bytes | Reduction vs 211,938,580 B original | Compress s | Decompress s |
|---|---:|---:|---:|---:|
| Zstd-3, historical external CLI | 66,216,569 | 68.76% | 1.2561 | 0.8551 |
| XZ-6, historical external CLI | 49,233,340 | 76.77% | 116.9892 | 2.3079 |
| **Current approved JLZ R57, retained historical 10-run median** | **60,585,717** | **71.41355%** | **5.493744** | **2.194428** |
| **R59 one-choice candidate, NEW 10-run median** | **60,400,234** | **71.50107%** | **6.118099** | **2.342502** |

| Compressed size delta R59 vs R57 | Encode time delta | Decode time delta |
|---:|---:|---:|
| **-0.30615%** (-185,483 B) | **+11.36484% nominal** | **+6.74772% nominal** |

External CLI Zstd/XZ timing conventions differ from JLZ in-memory measurements; all earlier codec timings were reused unchanged.

### Per-file full archived sizes

| Full source | Approved R57 B | New candidate R59 B | Delta B |
|---|---:|---:|---:|
| dickens | 3,155,856 | 3,155,856 | 0 |
| mozilla | 18,609,139 | 18,609,139 | 0 |
| mr | 2,729,261 | 2,729,261 | 0 |
| nci | 2,163,199 | 2,163,199 | 0 |
| ooffice | 3,033,712 | 3,033,712 | 0 |
| osdb | 3,425,086 | 3,333,043 | **-92,043** |
| reymont | 1,455,234 | 1,434,877 | **-20,357** |
| samba | 5,134,262 | 5,134,262 | 0 |
| sao | 6,261,713 | 6,261,713 | 0 |
| webster | 9,835,416 | 9,764,819 | **-70,597** |
| xml | 515,852 | 513,366 | **-2,486** |
| x-ray | 4,266,987 | 4,266,987 | 0 |
| **TOTAL** | **60,585,717** | **60,400,234** | **-185,483** |

Four complete files shrink, eight unchanged, worst size regression **zero**. Note that the R57 baseline was locally reconstructed from approved algorithmic source for correctness confirmation; official reference *timings* above remain historical independent ten-run values.

## Validation

- 12/12 canonical original Silesia SHA-256 pass; 10 complete full-corpus passes, 120/120 exact reconstructions and embedded SHA-256 checks.
- \`go test -count=1 ./...\` PASS; \`go vet ./...\` PASS; targeted \`go test -race -count=1 ./engine -run '^TestR59ConditionalExactAndCorrupt$'\` PASS; Windows/amd64 CGO-disabled cross-build PASS (compile only).
- Focused tests verify exact conditional rank state restoration, all original BWT bytes, and reject truncated conditional payloads for three synthetic fixtures; malformed normalized tables and length bounds are checked through the existing native finite-state reader. No supported production archive compatibility change occurred because **mode17 was never committed**.
- Experimental source and test files remain local and disposable, **not committed** to \`main\`. GitHub contains this research report and updated canonical backlog/state only.

## Verdict and next substantive hypothesis

**REJECT H1 as a production integration** under practical speed/complexity requirements: 0.306% archived-size gain for 11.36% nominal encode and 6.75% nominal decode overhead, plus an extra internal mode17 decoder, is not proportionate. Crucially **not a ratio regression** and **not proof that conditional BWT entropy is unhelpful**. It contains demonstrable conditional information beyond simple global static frequencies and beyond the currently selected Round57 local entropy regimes, concentrated in osdb/webster/reymont. Leave the broader *decoder-synchronized but cheaper* context grammar OPEN. Close these precise three-context interleaved byte-rANS and one-choice MDL implementations and do not parameter-sweep adjacent class thresholds.

For next investigation seek **a different encoding mechanism** capable of exploiting context without three models and per-symbol rANS table selection, or investigate another complementary substream where the same decoder-synchronous context is cheaper; maintain the user's non-negotiable no-compression-regression prerequisite and YAGNI.