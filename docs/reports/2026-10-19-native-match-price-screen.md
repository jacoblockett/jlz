# Native mathematical match-price calibration: reversible component diagnostic

**Status: SCREENED; no integrated codec candidate.** Date 2026-10-19 local project session. This followed the [source-level audit of Zstd, XZ and Brotli](2026-10-19-industry-architecture-audit.md), whose high-quality parsers estimate actual costs of mutually exclusive match/insert decisions.

## Evidence and distinction from historical trials

The **current** JLZ Go matcher in `engine/match.go` uses `scoreMatchChoice`: each candidate length is valued at a *fixed* five bits per replaced byte, except a short exact local copy at six bits, minus an approximate count of varint-like command bytes. Selection uses a one-byte lazy lookahead. The final encoder, however, sends independently native-packed opcode, command lengths, distance extra-bit slots, transforms and separately context-arithmetic-coded literals. There is no guarantee that the fixed five/six-bit score accurately reflects marginal output price.

This inquiry differs mechanically from:
- **Rejected exact-only dual-anchor parser**: replaced original XOR/ADD candidate generation and enlarged archives. The proposed research **keeps** transformed matches.
- **Rejected larger-history search**: expanded candidate search at heavy cost. Proposed research changes **selection among the matches already available**.
- **Rejected original+alternate full re-encoding**: consumed substantial duplicate CPU. The relevant idea is a **bounded cost estimate within one parsing pass**, not two complete entropy encodes.
- **Rejected match-length-conditioned distance stream splitting**: added separately framed entropy models and expanded bytes. This changes *which tokens to select*, **not** how current tokens are entropy-framed.

## Actual executed reversible standalone diagnostic

Ran an original disposable local Go unit test against the verified current `engine/match.go` content (Git blob `ff344090beddf2578a27af1354d51d90cac24b71`), using full original Silesia source files `mozilla`, `ooffice`, `samba` (the three math-selected whole-file paths). For each complete original file, selected its first, middle, and final natural 256 KiB source region, for **nine full region diagnostics**. No non-Silesia benchmark corpus or external compressor.

For each, actually ran the native JLZ match tokenizer, parsed its exact token columns, encoded that region's native legacy match archive with embedded SHA-256 and decoded it through JLZ's decoder; **9/9 exact byte-for-byte verified**. Measured Shannon 0th-order entropy of the **actual literals left by this tokenizer**, together with actual costs of existing native packed opcode, length, distance and literal streams. The H0 measure is a *diagnostic proxy*, **NOT the final context-adaptive arithmetic coder's exact marginal symbol costs**, and selected source blocks can take a different final winner after archive preflight.

Literal H0 in bits per output literal byte:

| Complete original Silesia source | Start block | Middle block | Final block |
|---|---:|---:|---:|
| mozilla | 7.946 | 6.852 | 6.847 |
| ooffice | 6.755 | 6.765 | 6.653 |
| samba | 7.895 | 4.474 | 5.686 |

**Seven of nine selected literal streams have H0 greater than 6 bits/byte**, one lies between 5 and 6, and one below 5. JLZ's fixed scorer values a replaced source byte at 5–6 bits without conditioning on where it lands in an entropy-coded literal stream. This is a **measured mismatch in model assumptions**, not evidence that changing the score would improve compressed bytes. H0 of raw input is also not its post-parse encoded cost, and transformed-references' mask/slot cache state complicates accurate marginal pricing.

Selected actual numbers from reversible native components:
- mozilla first block: 110,118 literal bytes / 152,026 referenced source bytes; literal H0 **7.946**, actual packed literal stream **107,968 B**, encoded op **1,651 B**, encoded lengths **5,595 B**, encoded distance stream **4,352 B**.
- ooffice middle block: 144,482 literal bytes / 117,662 referenced bytes; literal H0 **6.765**, actual packed literal stream **105,127 B**, op **3,561 B**, lengths **14,003 B**, distance **15,984 B**.
- samba middle block: 17,188 literal bytes / 244,956 referenced bytes; literal H0 **4.474**, actual packed literals **8,739 B**, op **4,926 B**, lengths **10,088 B**, distance **10,079 B**.

All units above are **one-pass selected-stream component outputs**, not whole JLZ archived size or an encoder performance benchmark. The test executable's duration is not comparable to full-corpus JLZ timing. Other fallback modes can still win the complete region; this test **does not prove** each nine source blocks is a final selected mathematical archive block. No candidate version, user-selectable mode, or output format was created.

## Next falsifiable mechanism and anti-spiral stop gate

Investigate one native **bounded-cost, repeat-cache-aware parser decision** that reuses existing exact/XOR/ADD candidates and the current op/length/distance/literal schema. The cost model should account for expected byte costs of actual native encoding streams (including repeat distances and literal-run framing). Before full integrated candidate work, measure how many distinct candidate decisions it would change and whether those changes produce favorable **actual packed output** on originally math-selected whole Silesia source blocks. It must **not** run a second full matching/entropy pass or create new coder modes; if computing the model dominates compression time, stop. A constant tweak from five to six/seven bits is a threshold sweep, **NOT** the novel hypothesis, and is prohibited by the anti-spiral rule.

A candidate is only promitable if **fully integrated, fully framed, byte-exact** and a ten-complete-original-Silesia-run median is nonregressing in all **three** accepted metrics: archive bytes, encoding time and decoding time. Do not infer a candidate win from H0 mismatch alone. The current accepted engine remains untouched.

## Provenance

Disposable actual test was `/mnt/data/jlz_cost_audit/engine/zz_price_audit_test.go`, run with `go test ./engine -run '^TestPriceAuditSelectedMathSources$' -count=1 -v`. The separate test file was not committed to production, as it is a one-off diagnostic. Source filenames containing historical numbered reports remain archival links only; no new round identifier was assigned.
