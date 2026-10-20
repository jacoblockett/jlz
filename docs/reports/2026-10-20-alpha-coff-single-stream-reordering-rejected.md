# Reversible Alpha COFF-aware encoding without duplicate full archives: tested, not promoted

**2026-10-20. REJECT both full-corpus integrated H1/H2 production candidates.** Current accepted JLZ on `main` remains **58,767,739 B** over the entire original 12-file Silesia corpus (211,938,580 B), approved historical ten-run medians **4.7873709325 s encode / 1.8721679705 s decode** with root `default.pgo`. No source edits or new stream subtype were committed. All alternatives used original native Go JLZ; no foreign compression backend or runtime dependency.

## Source-first motivation and exact experimental scope

Source identification corrected in [previous COFF report](2026-10-20-alpha-coff-source-structure-screens.md): original `mozilla` is a USTAR containing 118 genuine DEC Alpha ECOFF `.so` objects with native `.text` section spans. These sections total **23,975,888 of 51,220,480 original bytes**, leaving **27,244,592 noncode bytes**. The previously tested source-class split produced a **325,643 B** improvement by independently encoding two entire JLZ subarchives, but whole-corpus compression was ~42.6% slower in one diagnostic, and a new parser/format was costly. An alternative must retain the source partition without performing multiple full match-search/entropy encodings.

### H0: current math metadata and literals conditioned on *actual* .text source positions (component-only)

A disposable native-Go diagnostic ran JLZ's **current native match tokenizer** over all original Mozilla 256KiB input regions, preserving current 128KiB cross-group history decisions. Tokens were parsed to retain precise original source byte positions for literals and exact encoded op/length/distance fields. A verified 118-entry original COFF `.text` span map, not a guessed text/UTF8 heuristic, classifies each relevant token. The actual existing JLZ `packLiteralStream`, `packFiniteState`, and `packVarintStream` were used, including complete per-stream framing and an extra candidate mode byte per block. These are **component-cost projections, not an integrated archive codec** and no new decoder was implemented for class-conditioned token streams.

- Among **112 source blocks with both classes and sufficiently large literal streams**: incumbent packed literal component bytes **5,968,739 B**; code/noncode separately entropy-coded literals **5,864,470 B**. Net **-104,269 B**; **111/112** local block wins; 15,444,569 literal source bytes were encountered, 8,564,162 within real COFF `.text`.
- Across **113 source blocks with both class token metadata**: existing length components **1,590,668 B** vs split **1,583,516 B** (**-7,152**); distance **1,423,997→1,422,370 B (-1,627)**; opcode **540,259→535,129 B (-5,130)**. Net metadata **-13,909 B** after candidate per-stream costs; length/distance/opcode local block wins 65/75/89 respectively.
- The combined stage-level evidence is about **118,178 B**, before shared source class map (**696 B**) and new whole-file/region framing, decoder work, and any effect on JLZ's mode selection. It does not justify a specialized bitstream decoder on its own. No token-column inverse implementation was executed for this diagnostic, so **do not call it a reversible integrated coding scheme**. All source labels originate from bytes in the original corpus, not information magically available to a decoder.

### H1: one reordered native child JLZ archive, preserving inverse section positions

Implemented a new **temporary native-Go archive mode20**, scoped only to source files with a verified USTAR header and at least sixteen valid, disjoint real Alpha ECOFF `.text` sections. Existing source `.text` ranges are copied into one concatenated code group, all remaining original bytes into another, and a delta/length section map **696 B** restores exact original ordering. Unlike the previous expensive experiment, **the original JLZ compressor is invoked exactly once** on the *combined* group-ordered source, with the outer archive carrying the original digest and one complete child native JLZ archive. The native child encoder's recursive COFF detector is disabled; legacy decoders remain available, but old binaries cannot read new temporary mode20. No third-party compressor or new entropy mechanism is employed. Original source restoration is independently implemented, then the regular archive SHA-256 integrity check is applied.

Two original-full-Mozilla component screens and their exact native inverses:
- `code||noncode`: child native JLZ **18,129,543 B**, fully costed estimated frame/map **18,130,283 B** (about **-322,347 B** vs accepted full Mozilla **18,452,630 B**). One child encode diagnostic **2.307 s**.
- `noncode||code` (selected): child native JLZ **18,115,807 B**, additional map and outer frame actual complete **18,116,548 B**, saving **336,082 B** vs accepted Mozilla. One child encode diagnostic **2.116 s**; independent one-process run varied with environment. The noncode-first arrangement is the stronger of the two actual arrangements tested; do not sweep additional class orders.

The one-child format is a real reversible implementation: entire original **51,220,480 B Mozilla source** reconstructed **exactly**, with decoded-child digest and original-archive digest checked. Candidate size of 18,116,548 includes both full archive frames and the 118-entry map.

### Ten completed original-corpus benchmark passes, H1

All **12/12 original Silesia files** compressed, decoded, byte-compared, and SHA-256 checked per pass: **120/120 full-file roundtrips PASS**. Identical complete archive size every run: **58,431,657 B**, versus current **58,767,739 B (-336,082 B, -0.57188179%)**. Only Mozilla's individual archive bytes changed; all eleven other original file sizes and modes unchanged. No established JLZ, Zstd or XZ baseline was rerun as a whole-corpus benchmark.

| H1 candidate run | Complete bytes | Encode s | Decode s |
|---:|---:|---:|---:|
|1|58,431,657|4.915773884|2.143455163|
|2|58,431,657|5.039755770|1.931097586|
|3|58,431,657|4.946704275|2.036679539|
|4|58,431,657|5.007567372|1.808437137|
|5|58,431,657|4.781392546|1.928384915|
|6|58,431,657|4.945281398|1.907450511|
|7|58,431,657|4.916583995|1.870430891|
|8|58,431,657|4.795879339|1.883464640|
|9|58,431,657|4.731377177|1.846127945|
|10|58,431,657|4.845651575|2.004311420|
|**Median**|**58,431,657**|**4.9161789395**|**1.917917713**|

Versus previous approved historical complete-corpus medians **4.7873709325 encode / 1.8721679705 decode**, unpaired nominal difference **+2.690579% encoding time (worse)** and **+2.443677% decoding time (worse)**. These close temporal comparisons may include variable host noise, but **there is no demonstrated performance win**.

### H2: remove avoidable source concatenation and extra allocations

The H1 encoder first built separate **23.98MB code** and **27.24MB noncode** slices, then appended them, making a third **51.22MB temporary**. H2 writes original field segments **directly into exactly one 51.22MB destination**, filling the noncode prefix and code suffix, and produces **bit-identical native child archives and 58,431,657 B final corpus output**. This is a genuine localized elimination of redundant memory copies and allocations. Source-byte exact inverse validated. It did **not** improve measured whole-engine time, illustrating that JLZ's match/entropy traversal rather than the temporary copy is the bottleneck.

Ten **additional complete candidate-only full original Silesia** passes, **120/120 file exact**:

| H2 candidate run | Complete bytes | Encode s | Decode s |
|---:|---:|---:|---:|
|1|58,431,657|4.810485246|1.951292138|
|2|58,431,657|5.148586981|1.890551811|
|3|58,431,657|5.109605762|1.925444047|
|4|58,431,657|5.076311304|2.043325865|
|5|58,431,657|5.461488788|2.120503105|
|6|58,431,657|5.433333595|2.035760897|
|7|58,431,657|5.023246911|1.854128544|
|8|58,431,657|4.820332823|1.940877883|
|9|58,431,657|4.842442690|1.943387925|
|10|58,431,657|4.993420583|1.912117230|
|**Median**|**58,431,657**|**5.0497791075**|**1.942132904**|

Unpaired versus accepted historical: **+5.481258% encode time (worse), +3.737108% decode time (worse)**. Between H1 and H2 the median shifts are not proof that a reduced-copy algorithm is intrinsically slower; session and scheduler variation can dominate this small code change. Both fail the production metric gate.

## Reproducibility, scope, and negative findings

- Go1.23.2 Linux amd64 with exact accepted snapshot source (three key Go source blobs and root `default.pgo` matched latest accepted GitHub `main`), disposable trial under `/mnt/data/jlz_coff_literal_trial`.
- Prototype `go test -count=1 ./...` **PASS**, `go vet ./...` **PASS**, `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build` **PASS**. Tests include exact complete source reversibility and source-map reconstruction. A race-enabled full Mozilla COFF reorder run was started but **did not complete within the tool execution window**; **DO NOT claim race PASS**.
- Exact original Silesia corpus and existing checksums were used; both whole-corpus series completed all 120 original-file reconstructions. Existing archive subtype decoders were retained in experiment; no separate complete historical archive backward-decoding check was completed for this prototype. As the candidate is not promoted, no claim of production backwards-validation needed.
- No new production decoder mode, COFF parser, C++ source or external library was merged. Production `main` retains its approved **58,767,739 B /4.7873709325 s /1.8721679705 s** median metrics; ten-pass external Zstd-3 and XZ-6 historical references remain untouched.

## Decision

**REJECT H1 and H2 as production changes:** a **0.572% whole-corpus byte improvement** does not justify the specialized tar/ECOFF parser, decoder span-map reconstruction, one nested native archive and about **2.4–5.5% historical unpaired timing regression**. These results substantially improve the rejected previous duplicate-two-subarchive design **(+42.6% one-pass encode slowdown)** but still do not reach the accepted JLZ Pareto frontier. The component code/noncode classification savings of ~118 KB and one-native-child layout savings 336KB remain real reusable structural knowledge, not evidence of a standalone compression breakthrough.

**Next distinct experiments:** seek larger byte reductions without source-format-specific container modes; prioritize already-selected math token literals and reference parse economics, or cross-section learning using the current native codec's existing matched references. Do not parameter-sweep ECOFF section sorting or silently trade off 3–5% extra CPU for 0.57% size.
