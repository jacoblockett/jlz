# Round 58 — Variable BWT change points and length-conditioned distance coding (REJECT H1)

**Date:** 2026-10-19. **Production baseline remains approved Round57, unchanged.** This was a native-Go original experimental research round using exclusively the full original 12-file Silesia corpus (211,938,580 bytes). No external compression backend, source-code port, GitHub Actions, additional production wire format, or main-branch production modification. The user’s non-negotiable rule against sacrificing compression ratio for speed remains in force, but a large speed penalty also requires proportionate compression improvement.

## H1A: entropy-optimal variable neighborhoods after BWT

Round57’s mode16 independently frames 2 or 4 *equal-length* entropy segments after BWT/bijective MTF. **Key exact-format observation:** The existing mode16 decoder reads each segment’s decoded length from its native packed-stream frame and never requires equal lengths; therefore variable cut points need **no new decoder or explicit split-position metadata**. Fully reversible experimental Go code uses the same existing mode16 byte format, same BWT/MTF inverse, same finite-state entropy implementation, and shorter-of-current-versus-candidate exact-byte selection.

Implemented a **16-observation-bin minimum-description-length dynamic program**. For each candidate segment, estimate native zero-order bits + three bytes per distinct symbol + eight bytes framing, then choose four segment endpoints minimizing total modeled bytes. Only **one** chosen variable-cut partition is actually entropy-encoded, compared against Round57’s already encoded incumbent; all fallback archive bytes remain valid Round57 output. Cuts are not transmitted because packed stream headers supply exact lengths. Complete 12-file Silesia candidate archive bytes **60,377,031**, **208,686 B below R57**, with six complete files smaller and six unchanged, no regressions.

**Independent 451-block component screen:** examined *all 451 BWT-selected source blocks* in the seven BWT-selected original Silesia files; original existing Round57 native packed symbol-stream aggregate **23,266,779 B** (component-only). Evaluated exact native packing of alternate variable 2-/4-part DP partitions, including explicit boundary lengths for a conservative component-cost estimate, and verified the complete packed/unpacked byte stream for each candidate. Best selectable variable scheme estimated **23,020,889 B**, an optimistic component saving **245,890 B** across 447 blocks. This is NOT the final JLZ archive improvement. The integrated mode16 encoding does not need the separate boundary-length fields, but source selection and preflight sample reuse constrain actual full-archive savings.

### Ten COMPLETE original 12-file Silesia candidate-only runs

Each run compressed and decoded all twelve complete files with checksum verification; all 120 original files reconstructed exact byte-for-byte, every run emitted the identical archive sum **60,377,031 B**. The first two-pass executable call completed two runs; the second batch completed a third run before its command wrapper timed out. Seven additional completed one-pass invocations finished the ten-run set. The incomplete fourth attempt is not counted.

| Complete candidate run | Archive B | Encode s | Decode s |
|---:|---:|---:|---:|
| 1 | 60,377,031 | 6.206331935 | 2.328088723 |
| 2 | 60,377,031 | 6.386983795 | 2.457327695 |
| 3 | 60,377,031 | 6.502634300 | 2.381989633 |
| 4 | 60,377,031 | 5.961584796 | 2.249146090 |
| 5 | 60,377,031 | 5.665102552 | 2.269998219 |
| 6 | 60,377,031 | 5.919499759 | 2.298386556 |
| 7 | 60,377,031 | 6.121885337 | 2.316072808 |
| 8 | 60,377,031 | 5.809469356 | 2.165951548 |
| 9 | 60,377,031 | 5.904505123 | 2.207995076 |
| 10 | 60,377,031 | 6.149127121 | 2.320535411 |
| **Median** | **60,377,031** | **6.041735067** | **2.307229682** |

Retained approved **Round57** ten-run historical medians **60,585,717 B / 5.493743735 s encode / 2.194427893 s decode**, NOT re-run. **Real archive improvement -208,686 B (-0.344448% of compressed bytes)**, nominal **+9.975% encode / +5.140% decode**, unpaired across independent sessions and not precise causal attribution. Six files shrink (dickens -16,150; mr -74,708; nci -3,123; reymont -23,412; webster -89,243; xml -2,050 B); mozilla, ooffice, osdb, samba, sao, x-ray unchanged. **No per-file compression regression**.

Two structurally cheaper boundary searches were also actually integrated and tested on all 12 complete files, **one whole-corpus screen each (NOT ten-run)**:

- **8-bin DP with four segments:** **60,454,475 B**, one-run **5.658036 s encode / 2.239969 s decode**. Gains **131,242 B** vs R57, but a smaller potential ratio gain and encoding still directionally slower.
- **16-bin greedy divisive MDL** (choose the largest split gain three times; cache only visited segment costs): **60,384,224 B**, one-run **6.159224 s encode / 2.376452 s decode**. Saves 201,493 B vs R57 but does not solve its performance problem.

Both schemes are fully reversible and use no new decoder. The full dynamic-programming candidate passed `go test -count=1 ./...`, `go vet ./...`, targeted `go test -race -count=1 ./engine -run 'Test(BWTSubregions|BWTPreference|FiniteState)'`, and Windows/amd64 CGO-disabled cross-compilation. Sampled/complete source bytes were all drawn exclusively from the canonical full Silesia originals, and the complete ten-run test compared every decoded file against its input. The unchanged R57 local production snapshot also passed uncached Go package tests and vet. No native Windows runtime test is claimed.

## H1B: conditioning mathematical distance bytes on already decoded match lengths

Implemented a separate fully reversible **decoder-synchronous cross-field entropy partition** on actual selected mathematical archives, not a generic byte-stream shuffling heuristic. Read complete correct R57 archives from the three math-selected full Silesia files (mozilla, ooffice, samba); extracted **299 selected block distance streams** along with their actual packed operation and length streams. Where a split-length block was used, remerged exact original operation lengths before testing.

For every command, use its *already decoded associated match length* to choose one of 2 (length ≤16 or >16) or 3 (≤16, 17..128, >128) metadata contexts. Route each original exact encoded distance varint or repeat-transform byte to its context and pack that context separately with existing native finite-state or legacy packed coding. All new parts were decoded and **reinterleaved according to the original op/length stream**, with exact original distance bytes compared and fully charged per-context framing. No per-symbol role labels or learned model side-information are assumed.

| Selected math distance component | Complete packed bytes across 299 streams | Delta from original |
|---|---:|---:|
| R57 original distance streams, native `packVarintStream` | **4,146,842** | Reference |
| Two match-length contexts | 4,397,180 | **+250,338 B** |
| Three match-length contexts | 4,398,608 | **+251,766 B** |
| Best per-stream oracle among old/two/three | **4,146,746** | **-96 B**, only **1 of 299** streams wins |

**REJECT H1B** as an actual selected-stream compression path. The finite-state context fragmentation and per-context model framing swamp any conditional information, even though the context selection is reversible and essentially free to the decoder. No full JLZ integrated archive was created for this clearly unfavorable component hypothesis, and no promotional ten-run results are implied.

## Required complete-corpus comparison: strongest full JLZ candidate H1A

| Model | Complete archive bytes | Reduction of 211,938,580 B | Compression time | Decompression time |
|---|---:|---:|---:|---:|
| Zstd-3, retained historical external CLI | 66,216,569 | 68.76% | 1.2561 s | 0.8551 s |
| XZ-6, retained historical external CLI | 49,233,340 | 76.77% | 116.9892 s | 2.3079 s |
| **Approved JLZ Round57, retained historical ten-run median** | **60,585,717** | **71.41355%** | **5.493744 s** | **2.194428 s** |
| **Candidate JLZ Round58 DP, new ten-run median** | **60,377,031** | **71.51201%** | **6.041735 s** | **2.307230 s** |

| Candidate archive-size delta | Candidate encode-time delta | Candidate decode-time delta |
|---:|---:|---:|
| **-0.34445%** | **+9.97483% nominal** | **+5.14037% nominal** |

Historical Zstd/XZ external CLI times and JLZ in-memory times are not methodologically equivalent. Accepted JLZ and external benchmark references were not rerun.

## Final disposition and genuinely open next steps

**REJECT H1 as a production promotion.** Variable BWT cut points have real exact-byte compression value and require **no new wire subtype**, but the tested dynamic program’s 0.344% compressed-size saving currently incurs roughly 10% nominal encoding and 5% decoding wall-time cost; simpler greedy search did not correct the trade. The eight-bin version preserves less savings and is also directionally slower. Under YAGNI and practical user expectations this is not competitive enough despite meeting the strict *no ratio regression* prerequisite. **Keep Round57 as approved source; all variable-point and distance-context experiment code stays disposable.**

Theoretically **OPEN**: fundamentally cheaper change-point *discovery* from existing BWT suffix/context statistics or one-pass model transitions; a useful adaptive coding approach need not transmit explicit splits because mode16 is already compatible. **CLOSED for these concrete implementations:** 16-bin DP, 16-bin greedy, 8-bin DP, and simple two-/three-way match-length-routed distance varint streams. Future experiments should prioritize mechanisms with a material net archive advantage and no disproportionate CPU or decoder-state overhead. Do not simply sweep more bin counts or classifier thresholds.