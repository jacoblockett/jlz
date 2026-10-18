# Round 28: Activate theory backlog; bounded prime-index and information-synergy screens

**Date:** 2026-10-18. **Decision:** two narrowly scoped theory-derived prototypes REJECTED at screening; **no full-theory rejection** and **no new JLZ archive candidate or source promotion**. The approved standalone native compressor remains Round 25. Research scope is the full twelve-file original Silesia corpus only; no external-file benchmarks and no existing reference codec or approved JLZ benchmark was rerun.

## Research direction correction

The user's theoretical ideas catalog is **an active source of experimental candidates, not a veto list of historically failed ideas**. When incremental codec tuning stalls, pick an idea from that catalog and actually prototype it. `docs/research/attempted-ideas.md` remains useful as an account of particular historical prototypes and warnings about synthetic-only wins or external-codec dependencies. Neither the title nor past negative tests should be used to dismiss a theory in general. Track concrete follow-up attempts and future variants in [theory backlog](../research/theory-backlog.md).

The historical inventory explicitly says several ideas were *not implemented*, including prime-index/irreducible-prime strategies and higher-order partial information decomposition. Two **limited implementations** were screened in disposable copies on the current Silesia bytes. These are not claimed to implement the entire mathematical theories.

## Screen A: prime-indexed byte partition

Hypothesis: treating prime-numbered within-region positions as an independent arithmetic-index stream might expose structure ignored by ordinary contiguous byte coding.

Method: For each 256 KiB original input region, enumerate index primality with a simple sieve, partition original bytes into prime-indexed and other-positioned streams, independently encode both using **JLZ's existing native raw/RLE/Huffman packer**, and reconstruct by deterministic index scattering. No primality bitmap is stored; both lane lengths are known from the region length. Both decoder functions were implemented in the disposable test and proven byte-exact on all 36 sampled regions (first, middle, last for each of 12 files).

Sampled native packer output was **5,140,826 B**, versus split native packed streams **5,133,446 B**, a component saving of **7,380 B (0.144%)** before any new JLZ block framing. `x-ray` accounted for **6,309 B** of the saving; other files were mixed. This prompted a full-corpus **component estimate**, not a new JLZ archive.

The full 211,938,580-byte Silesia corpus was processed through the prime split and the two native packers region by region. A notional archive includes 40 B file overhead plus 8 B of region metadata (approximate, not canonical JLZ framing). The estimated prime-only compressed size was **133,180,437 B** versus the approved Round 25 **63,104,117 B** archive size. **All twelve prime-only estimated files were larger than their Round 25 archives**, even before any extra actual decoder-mode overhead. Therefore no actual complete-prime JLZ1 archive was built and no ten-run candidate promotion benchmark was warranted. The 133 MB number is a *component-derived archive estimate*, not a verified alternative JLZ complete archive.

This result rejects only the **prime-position byte segregation + ordinary native Huffman** proposal. It does not test irreducible prime factor bases, prime-index recurrence models or general prime-based descriptions; those remain eligible catalog ideas.

## Screen B: higher-order information-synergy inspiration

Hypothesis: two individually weak historical contexts might jointly predict a target bit through XOR, exposing synergy that order-one modeling misses. This is a narrow **PID-inspired parity proxy**, not a formal partial information decomposition or proof of its absence.

Selected 11 fixed pairs of byte lags (`1,2`, `1,4`, `2,4`, `1,8`, `4,8`, `1,16`, `7,11`, `11,17`, `31,127`, `127,509`, `509,1021`). For each Silesia file, measured the eight bitplanes' conditional Shannon entropy for the two independent lagged bits and for the XOR of both. Picked a lag pair on first/middle 64 KiB windows and assessed the same pair on the **held-out last 64 KiB window**. These entropy numbers are **diagnostic, not a coder or archived byte gain**.

Eleven of twelve held-out files favored the best single-lag model over the XOR pair. `x-ray` was the exception: lag `7,11` had a **+0.8927 bits per source byte** difference in favor of parity versus the better of the two binary single-lag conditional models on held-out data. This was not sufficient evidence to improve native compression, so an actual deterministic and byte-exact original-byte lagged-XOR residual was implemented in the disposable copy.

For each file, used the selected lag pair and compared **three 256 KiB blocks** encoded as (1) native Huffman packer on original bytes, (2) native packer on the reversible XOR residual, and (3) existing native **BWT/MTF/zero-run/Huffman** (all native JLZ components). Every residual restored its exact original bytes. For `x-ray`: original native packed **634,928 B**, XOR-residual packed **778,895 B**, and existing BWT packed **429,274 B**, across three sampled blocks. **The residual was larger than BWT for all twelve sampled files**, with a substantial margin. No new archive mode, repeated full corpus encoding, or candidate timing benchmark was justified.

This result rejects only the **fixed two-lag XOR residual packed with the current native packer** as an improvement on the sampled blocks. It does not reject a formal high-order PID model, data-dependent multivariable contexts, learned graph relationships or other catalog theories.

## Benchmark policy and references

The approved complete Silesia JLZ Round 25 remains **63,104,117 B**, with its previously measured median encode **7.549920 s** and decode **2.839905 s**; these were **not rerun**. Historical Silesia fixed reference sizes remain **Zstd-3 66,216,569 B** (historical CLI encode 1.256 s) and **XZ-6 49,233,340 B** (historical CLI encode 116.989 s); no repeat runs were performed and their CLI times cannot be compared directly with JLZ in-memory medians. This screening did not create a promotable full-archive candidate, so **there are no new candidate archive sizes or ten-run medians to put in the four-row promotional comparison**.

All work was on original Silesia bytes in disposable local Go 1.23.2 and numpy diagnostic scripts. Three custom Go test screens passed: `TestPrimeIndexResearch`, `TestPrimeFullFileScreen`, `TestParityTheoryResearch`; deterministic inverses were checked for every encoded sampled block. Go tests/vet of a new source candidate, full compatibility or race audits were not applicable, because no candidate source edits were proposed or retained. No compiled binaries or copyrighted corpora were committed.

## Verdict and forward work

**REJECT the two *specific* low-cost screens, not their parent theories. Do not promote anything.** Prime-position segregation is far above the existing complete archive size, and simple parity residuals lose against the native BWT coder despite one promising idealized information signal.

**Next active theoretical investigation:** Polynomial Freiman–Ruzsa-inspired *structured additive sets/progressions* or a higher-order multivariable information decomposition with a genuinely new predictive coding structure. Start with bounded, held-out Silesia screens and graduate to a full native, byte-exact JLZ candidate only where there is credible potential to beat the current archives. Do not return to incremental BWT/lane coder tuning by default. Update the catalog statuses after each attempt.
