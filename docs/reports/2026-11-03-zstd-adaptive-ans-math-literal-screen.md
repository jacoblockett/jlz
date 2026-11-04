# Zstd-inspired decoder-synchronous adaptive ANS (2026-11-03)

**Status: BENCHMARKED INTEGRATED CANDIDATE, NOT PROMOTED.** Current accepted JLZ engine, root default.pgo and archive format on main are unchanged. This work is a new native-Go entropy-coder architecture screen, not a switch to a preexisting codec backend.

## Source evidence, collision audit and priority

The active Zstd-only architecture queue was checked in docs/research/zstd-architecture-coverage.md and the historical strategy ledger, placement matrix and reports. Primary reference: Meta Zstd's documented reverse-encoded finite-state entropy and independent sequence tables (https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md), and current JLZ entropy/arithmetic.go, engine/stream.go, engine/finite_state.go. Earlier native binary arithmetic/range variations, direct EMA probability changes, static multi-context rANS, predecessor-symbol adaptive Fenwick, and indirect state ICM are already measured negatives or expensive wins. Unlike static rANS, this experiment **keeps precisely the existing decoder-synchronous bitModel probability and update policy**; unlike E3 re-normalization experiments, it changes the representation to a reverse ANS state stream.

Other established compressors were cross-checked against older already-executed JLZ tests rather than needlessly retried: XZ/LZMA and libdeflate costed parsing (current match scorer and earlier slow candidates), Brotli's context models and clustered histograms (already sampled/expensive in math literals), libbsc's BWT/QLFC and LZP (past selected-stream trials rejected), and ZPAQ ICM/ISSE (BWT size benefit, excessive decode/encode work). No other external compressor source was imported or used as a JLZ backend.

## Original native algorithm

At each original mathematical literal bit, compute the exact current JLZ adaptive order-0/order-1 combined binary probability p0 in 1..65535 from the previous decoded literal and the current symbol prefix. During an initial forward encoder pass, store p0 (16 bits per input bit), update the same model as existing arithmetic E3, then rANS-encode these bits **in reverse** using total frequency 65536 and normalized L=2^23. The decoder reads the new state stream and generates the same p0 values from previously restored bits in ordinary forward order. It demands exact final state L and complete payload consumption, and retains existing outer SHA256 verification. The new stream codec subtype is 6. Existing subtype 3 and older modes remain readable; older JLZ decoder binaries cannot read new subtype6 output. Raw/RLE/Huffman fallback remains unchanged. This is a **supplement to existing math literal placement**, not a parser or BWT replacement.

Encoder temporary frequency-buffer requirement is 16 bytes per input literal byte, up to ~4 MiB for a 256KiB stream, additional to native model and output. No aggregate peak-heap measurement was conducted. A deliberate smaller 24-bit-state variation was also tried and rejected because it expanded the 18 sampled input streams by 6,359 B versus existing E3, as opposed to +55 B for the four-byte state.

## Component test (natural original input; not whole-codec benchmarks)

Current V2 tokenizer mathematical literal streams from ten actual original Mozilla 256KiB source blocks and eight Samba 256KiB blocks. This scope is eligible tokenizer data, not independently verified to be each block's final selected archived literal mode. All 18 streams were separately encoded with accepted E3 and new adaptive ANS, decoded and byte-equal verified:

| Source samples | Literal input B | E3 body B | Adaptive ANS body B | E3 encode ms | ANS encode ms | E3 decode ms | ANS decode ms |
|---|---:|---:|---:|---:|---:|---:|---:|
| Mozilla, 10 blocks | 888,549 | 775,744 | 775,774 | 233.431 | 133.305 | 232.886 | 177.191 |
| Samba, 8 blocks | 554,315 | 473,281 | 473,306 | 130.669 | 90.395 | 141.151 | 91.916 |

One original-Go component pass, not repeated timing medians. 18/18 exact and new state payload finalization checked. Sample compressed size **+55 B** on 1,249,025 coded E3 bytes (about +0.0044%); speed improved materially.

## Actual integrated complete original twelve-file Silesia

Original user ZIP 12/12 full source SHA256 matched committed manifest docs/reports/2026-10-18-round23-silesia-sha256.txt; **211,938,580 B** exactly. The disposable integrated encoder used adaptive ANS instead of running the old adaptive arithmetic alternative; old E3 decoder remains. Candidate archive full twelve-original total **58,768,442 B**, versus retained accepted **58,767,739 B**, **+703 B (+0.001196% archived size worse; original reduction 72.27100323% versus accepted 72.27133493%, -0.00033170 percentage points worse)**. Mozilla **+461 B**, Samba **+230 B**, Ooffice **+12 B**; other nine files **unchanged compressed size**. Ten complete candidate-only runs: **120/120 full original-byte and SHA256 restorations pass**, identical total every pass.

| Candidate full run | Complete original archive B | Encode seconds | Decode seconds |
|---|---:|---:|---:|
| 1 | 58,768,442 | 5.297234437 | 1.804216716 |
| 2 | 58,768,442 | 4.965405000 | 1.743472093 |
| 3 | 58,768,442 | 5.158765814 | 1.720638837 |
| 4 | 58,768,442 | 5.007475641 | 1.642134173 |
| 5 | 58,768,442 | 5.031469809 | 1.685532555 |
| 6 | 58,768,442 | 5.094434183 | 1.809418789 |
| 7 | 58,768,442 | 4.989903624 | 1.671815769 |
| 8 | 58,768,442 | 5.066580154 | 1.693058323 |
| 9 | 58,768,442 | 5.077004619 | 1.724598551 |
| 10 | 58,768,442 | 5.308170003 | 1.694327995 |
| **Candidate median** | **58,768,442** | **5.0717923865** | **1.7074834160** |

The historical accepted JLZ complete ten-run median 4.7873709325 /1.8721679705s and external Zstd-3 66,216,569 B /1.2561/0.8551s, XZ-6 49,233,340 B /116.9892/2.3079s were **not rebenchmarked as references**. Cross-host historical-vs-candidate speed differences are not a valid causal speed comparison.

## Controlled same-host full-original A/B diagnostic, not a new accepted reference

One additional four-pass **BASELINE/ANS/ANS/BASELINE** sequence and one four-pass **ANS/BASELINE/BASELINE/ANS** sequence, with the same executable and all twelve full originals in every pass, produced four current-E3 controls and four ANS candidate trials. All outputs exact, baseline archives 58,767,739 B and ANS archives 58,768,442 B. Medians across both alternating sets:

| Same-host control | Encode seconds | Decode seconds |
|---|---:|---:|
| Original E3, four passes | 6.4325047080 | 2.0451690705 |
| New adaptive ANS, four passes | 5.2870131900 | 1.7215687025 |
| **ANS relative** | **-17.80786% (better)** | **-15.82267% (better)** |

The fastest baseline control was still slower than the slowest ANS trial for both encode and decode, although host drift and small sample size limit generalization to other platforms. These control timings are *diagnostics*, not a replacement for the previously retained accepted reference.

## Validation, decision and next work

Full native-Go tests **PASS** after updating the two previously E3-hardcoded test expectations to accept the new optional codec6, plus explicit old codec3 decode fixture. go vet ./... **PASS**; targeted adaptive-ANS corruption/truncation, roundtrip, and race tests **PASS**; CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ./... **PASS** (cross-build only, no Windows runtime). Corrupted and truncated adaptive ANS bodies cannot produce the original bytes without detection; new state termination and outer archive digest validated. Any new archive using subtype6 is **forward-incompatible** with earlier JLZ binaries. No additional production code, binary, PGO or branch committed.

**BENCH / DO NOT PROMOTE** until user explicitly chooses the small space/format/memory trade. This is currently the strongest measured source-audited *speed* opportunity, not an absolute compression-size improvement. A material size-oriented separate priority is cheaper BWT conditional prediction, where historical ZPAQ-ICM-inspired JLZ integration saved **886,032 B** but was far too slow. The existing Zstd-inspired cross-history length split saved **285,798 B** but worsened paired CPU; do not relabel either as untried. Avoid repeating prior parameter sweeps. If continuing adaptive ANS, first reduce the 16x temporary probability buffer or amortize terminal-state overhead through a genuinely different model/framing mechanism. User-distributable original-Go candidate patch and raw benchmark logs are retained outside the repository; no rejected experiment source is committed.
