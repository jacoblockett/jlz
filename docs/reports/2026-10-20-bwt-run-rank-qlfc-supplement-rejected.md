# BWT run-rank coding inspired by libbsc QLFC

**Date:** 2026-10-20. **Disposition:** REJECT both full JLZ integrated candidates; accepted production source unchanged. All experiments used original Go code in a disposable working copy, never a libbsc backend.

## Source-level investigation and distinct mechanism

Inspected actual libbsc sources at https://github.com/IlyaGrebnov/libbsc/blob/master/libbsc/coder/qlfc/qlfc.cpp and https://github.com/IlyaGrebnov/libbsc/blob/master/libbsc/coder/qlfc/qlfc_model.cpp. Unlike JLZ's BWT→MTF→bijective-zero-run→independent native Huffman/FSE, libbsc's QLFC models rank and run statistics separately. Tested a deliberately smaller ORIGINAL JLZ-native supplementary representation: encode each BWT symbol run as an MTF rank and a varint run length, then independently encode the rank/length byte streams with existing packFiniteState (which compares JLZ-native Huffman/FSE). Charge primary index, number of runs, both complete stream frames and outer block framing. The inverse reconstructs every BWT byte and uses the current native inverse BWT. No third-party source was ported.

## Exact original-source component screens

Input was 12 unmodified original Silesia members, 211,938,580 B, verified against the committed SHA-256 manifest. On three 256KiB first/mid/end natural blocks from each original, incumbent BWT component payloads totaled **3,218,996 B**, run-rank **3,258,392 B** (+39,396 B); 1/3 sampled mozilla, 3/3 osdb, 2/3 sao won in this isolated comparison, but sao/mozilla are not necessarily BWT-selected.

To avoid the mode-selection error, separately compressed ALL TWELVE complete originals with unchanged production source, extracted all **451 actually selected BWT blocks** from their real archives, recomputed reversible rank/run columns and compared fully framed per-block alternatives. Only **41/451** could improve, with an optimistic **90,001 B gross** saving, concentrated in all 39 osdb blocks (**89,665 B**) plus mr (**60 B**) and xml (**276 B**). This was a component-level oracle, NOT an end-to-end JLZ result. A full-block MDL gate (rank entropy + run-length entropy + modeled table cost) with 512 B margin predicted all 39 osdb wins and no actual loss. The gate excluded two tiny cross-file wins.

## Fully integrated disposable candidates

Integrated native block subtype 20 in the existing mode17 regional archive, as an optional BWT supplement, with bounds/trailing verification and valid inverse. In the faster second implementation, reuse already generated bijective MTF symbols to recover rank/run columns without a redundant BWT/MTF transform, then preflight an 8KiB symbol prefix to skip unpromising full-block model work. This is content-based and automatic; no file names or per-file rules.

**Retained accepted JLZ reference:** 58,767,739 B / 4.7873709325 s encode / 1.8721679705 s decode, historical ten full runs, not rerun as a ten-run reference. **Fixed Zstd-3:** 66,216,569 B / 1.2561 / 0.8551 s. **Fixed XZ-6:** 49,233,340 B / 116.9892 / 2.3079 s. External CLI timings are not measured identically to JLZ's in-memory timing.

| Model | Full original Silesia archive bytes | 10-run median encode s | 10-run median decode s | Size vs accepted | Encode vs accepted | Decode vs accepted |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Accepted current, historical | 58,767,739 | 4.7873709325 | 1.8721679705 | anchor | anchor | anchor |
| Integrated H1, raw rank/run | 58,683,911 | 5.5528247760 | 2.0205714915 | -83,828 (-0.14264%, better) | +15.989% (worse) | +7.927% (worse) |
| Integrated H2, prefix-gated rank/run | 58,683,911 | 5.5393677540 | 2.0352384510 | -83,828 (-0.14264%, better) | +15.708% (worse) | +8.710% (worse) |

Each H1/H2 run encoded and decoded ALL TWELVE complete original files, checked exact bytes and SHA-256 for **120/120** full files per candidate, confirmed 58,683,911 archive bytes on all ten repetitions. **Only osdb changed**: 3,425,086→3,341,258 B (-83,828). Absolute source reduction rose from 72.27133493% to 72.31088790% (+0.03955297 percentage points). Timing comparisons use historical unpaired sessions and are affected by host variation; they do not establish a Pareto speed improvement. The saved size does not justify a new wire mode, encoder path or source complexity with the measured speed regressions.

H1 encode seconds: 5.8299651350, 5.5398684610, 5.8497277120, 5.4741186730, 6.2211369750, 5.5657810910, 5.1786022670, 5.6593031960, 5.3150979550, 5.3910921500.
H1 decode seconds: 2.0406762740, 2.1057576760, 1.9813904760, 2.0269649960, 2.0703612880, 2.0340993090, 1.8784695340, 2.0141779870, 1.9664715020, 1.9171290830.

H2 encode seconds: 6.5589694490, 6.0317331140, 5.9781390660, 6.0063765970, 5.6684308080, 5.3650415670, 5.4103047000, 5.3028363660, 5.3089678170, 5.3131216100.
H2 decode seconds: 2.2456990550, 2.2181825970, 2.1456627560, 2.0937068190, 2.0571488800, 1.9351778470, 1.9226119750, 2.0133280220, 1.8969534170, 1.9866739540.

## Tests, caveats and disposition

Go 1.23.2 linux/amd64, root default.pgo SHA-256 3b35cda74832360a252b63c1805e025aa1f4ed7a4fe4124d9a55929a6e86289f (unchanged). New rank/run targeted synthetic exact and malformed-input tests PASS; all 36 source-sampled natural BWT blocks correctly roundtrip; ten full runs for BOTH variants PASS; uncached go test -count=1 ./..., go vet ./..., targeted run-rank race test, and Windows/amd64 cross-build PASS. Original old archive byte-decoding was not separately timed; prior modes remained in the candidate decoder and the compatibility suite passed. New block subtype would require a new decoder if accepted, while older archives remain structurally supported.

**REJECT H1/H2 concrete implementations**, not the broader run-aware conditional coding theory. Both experimental source variants and their tests were removed; all original production source files and root PGO were restored exactly from the validated ZIP and uncached Go suite passed again. Accepted baseline and decoder format remain unchanged. Future investigations should seek material source-model gains in mozilla's mathematical columns or webster's BWT text that do not require an additional expensive entropy pass on almost every BWT block.
