# Zstd cache-local match row: historical four-slot novelty collision and measured alternatives

**Date:** 2026-10-20. **Disposition:** NO CANDIDATE. **Correction:** Historical `docs/reports/2026-10-19-round42-match-depth.md` explicitly tested four candidate offsets and ACCEPTED reducing them to two for ~4.15% faster historical ten-run encoding at +0.395% compressed size. Simply restoring a four-position hash index is **NOT an untested Zstd architectural innovation**. Its screen here is evidence to stop repeating depth variants, not a new research direction.

## Original-Go local experiments

Used fully verified accepted Go snapshot HEAD `8566c8a8b28e3161ed941a98ec6ed4fc69906e47`, committed `default.pgo`, original 12-file Silesia 211,938,580 B. All listed 12-file one-pass archive screens included full framing and restored every member byte-exact. Tested 4 contiguous per-hash source offsets, then 2×uint64 shifting rather than four stores; asymmetric 4 exact / 2 XOR-ADD candidate retention; and early skip of third/fourth exact if a profitable/long match already exists. Never changed decoder, frame, or external dependencies.

| Independent single-pass full-corpus screen | JLZ compressed B | Delta vs current 58,767,739 B | Encode time s (NOT median) |
|---|---:|---:|---:|
| Accepted 2-deep source baseline | 58,767,739 | 0 | 5.934812 |
| 4 slots exact, XOR, ADD in 32-bit arrays | 58,578,358 | -189,381 **(better)** | 6.993296 |
| Packed two 64-bit shifts for all four slots | 58,578,358 | -189,381 **(better)** | 6.166090 |
| 4 exact / 2 XOR-ADD, shared 16-B slot layout | 58,584,808 | -182,931 **(better)** | 5.568148; separate pass 6.127610 |
| Four exact only when current best match shorter than 16 B | 58,618,765 | -148,974 **(better)** | 5.898323; separate pass 5.750145 |
| Probe extra exact only after first two fail to find any profitable match | 58,708,852 | -58,887 **(better)** | 5.563617 |
| 4 exact slots, forcibly disable transformed search | 58,605,867 | -161,872 **(better)** | 5.890907 |
| Conditional extra exact and forcibly disable transformed search | 58,730,161 | -37,578 **(better)** | 5.347561 |

**Paired, separate-process timing controls** for the last conditional four-exact *with transformed search preserved* were baseline/candidate **5.306491/5.710522 s (+7.6% slower, worse)** and **5.605190/5.635479 s (+0.54% slower, worse)**. Another 4-exact-only pass **6.127610s** followed a **5.942464s** baseline (+3.1% slower, worse), although its initial pass was faster. Do not infer causality from noisy one-pass times. Archive-size gains correspond to more exact matches, as already known from Round 42. Transformed-exact-only combinations look fast but **failed five pre-existing unit/selection tests**, including XOR/ADD fixture selection, default mixed selection and regional math selection, when checked by `go test -count=1 ./...` on the disposable prototype. Do not weaken those tests or promote an implementation that loses those capabilities.

**Decision:** Close simple four-offset reinstatement, different storage words, and nearby cutoff tweaks as a re-run of historical R42. No ten-run candidate qualified; no source promotion. A truly cache-aware Zstd batched *fingerprint comparison* or pipeline-prefetch mechanism would require measured hash/anchor miss costs and preserved mathematical XOR/ADD support, not more raw depth. The accepted codec remains **58,767,739 B / 4.7873709325 s / 1.8721679705 s** (historical ten-run; decoder unchanged). Accepted unchanged tests, vet, targeted race, Windows cross-build passed. No mature backend or external compression code used.
