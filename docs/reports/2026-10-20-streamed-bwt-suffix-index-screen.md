# Direct suffix-index stream projection into native BWT, 2026-10-20

**Disposition: REJECT tested native-Go writer variants H1/H2, no production code promoted.** The currently accepted unchanged complete original-Silesia result is **58,767,739 bytes / 4.7873709325 s encode / 1.8721679705 s decode**, historical medians of ten full passes (not rerun as a reference). Zstd-3 fixed historical 66,216,569 B / 1.2561 / 0.8551 s and XZ-6 49,233,340 B / 116.9892 / 2.3079 s. This study is a **compression-time architecture** test with no prediction of better encoded bytes.

## Source-level research, novelty gate and placement

Read [libbsc actual BWT implementation](https://github.com/IlyaGrebnov/libbsc/blob/master/libbsc/bwt/bwt.cpp), which builds BWT directly using libsais, and the local Go standard library `index/suffixarray/suffixarray.go` `Index.Write` source. No external compressor, libsais dependency, or copied implementation was used in JLZ. The current `bwt.Transform` creates `suffixarray.New(data)`, serializes the entire suffix index including redundant source bytes and chunked varints to a `bytes.Buffer`, then reparses the serialization to emit the BWT bytes and primary index. Serialization is an implementation-dependent indirect exposure of the standard library's otherwise private suffix order.

The closest earlier rejected mechanism was a separately designed **radix-prefix-doubling suffix constructor** on 21 original Silesia samples. It was 2–7× slower; this experiment is **not** another SA-building algorithm or size/cutoff sweep. It preserves the existing standard-library O(N) suffix-array builder and uses a custom original Go `io.Writer` to project its emitted varint chunks directly into the final BWT result. Validates fixed header, original source bytes, record lengths, all suffix offsets, exact record count and primary index. BWT forward result and bitstream are identical to incumbent; no extra encoder pass, archive flags, format, or decoder path.

- **H1:** one 10-byte header, one source write, then independent 16KiB-ish suffix-array chunk writes consumed in place; avoids redundant serialized intermediate bytes.
- **H2:** same writer plus in-place deletion of the sentinel in the BWT output buffer, removing another output allocation/copy. These are distinct implementation forms of the *same* pipeline work-elimination theory, not distinct compression models.

## Verified original-source component screens

Loaded *original, manifest-matching* Silesia samples at beginning/middle/end of each of Mozilla, osdb, webster, mr, nci, sao, xml and samba (24 independently selected 256KiB natural blocks). Every H1/H2 generated exactly the incumbent BWT byte string and primary index. The control implemented the current suffix-array serialization + decode parser and compared outputs. The following are component-only `go test -bench` observations with 24 blocks per iteration, not complete JLZ timings, two measurements per baseline/H1 and three per H2; host effects apply:

| Component path | Time per 24-block iteration | Allocated bytes per iteration | Allocations |
| --- | --- | --- | --- |
| Current serialization/read | 269,855,364 / 264,838,564 ns | 85,111,456 / 85,111,525 B | 248 / 249 |
| H1 direct projection | 250,135,439 / 249,276,215 ns | 39,694,218 / 39,694,112 B | 153 / 151 |
| H2 in-place compaction | 267,529,469 / 288,256,736 / 258,223,657 ns | 33,402,748 / 33,402,696 / 33,402,648 B | 128 / 127 / 127 |

H1 component median estimated ~6.6% faster than baseline and ~53% fewer allocated bytes. H2 substantially reduced allocations but was not reliably faster on the component screen. This improvement occurs only where BWT is actually built; the full engine routes other input through mathematical and record-lane paths.

## Full original twelve-file Silesia integrated experiments

Compiled disposable native Go test executables with explicit `-pgo=./default.pgo` (root PGO SHA-256 unchanged). Every pass compressed and decompressed **all 12 unmodified originals**, totaling 211,938,580 B, including the complete archive framing; checked byte equality and source SHA-256 against `docs/reports/2026-10-18-round23-silesia-sha256.txt`. The ten passes for each variant were organized as five batches of two **separate candidate-only test processes** due to execution-shell timeouts in an initially longer run. Each of the ten completed passes per variant restored all 12 SHA-exact. Size was **58,767,739 B on every pass**, no per-file regressions: the serialized archives are unchanged.

| Model | Complete bytes | Median encode s (ten complete runs) | Median decode s (ten complete runs) | Encode vs approved historical | Decode vs approved historical |
| --- | ---: | ---: | ---: | ---: | ---: |
| Accepted current (historical) | 58,767,739 | 4.7873709325 | 1.8721679705 | anchor | anchor |
| H1 stream projection | 58,767,739 | 4.8466381300 | 1.7574856245 | +1.238% (worse) | -6.126% (better, unpaired/no decoder changes) |
| H2 in-place projection | 58,767,739 | 4.8480981235 | 1.7871587915 | +1.268% (worse) | -4.541% (better, unpaired/no decoder changes) |

H1 encode seconds: 4.832245755, 4.822993146, 5.009524061, 4.933785321, 4.938146034, 4.693877917, 4.938268129, 4.814759246, 4.861030505, 4.747618239.
H1 decode seconds: 1.758531404, 1.796227436, 1.797705488, 1.887713479, 1.754603321, 1.754001269, 1.745113958, 1.756439845, 1.812044663, 1.752934841.

H2 encode seconds: 5.089621284, 4.851894008, 5.011621915, 4.813558404, 4.988981731, 4.801402483, 4.812620220, 4.718061896, 4.980717342, 4.844302239.
H2 decode seconds: 1.904170536, 1.780306208, 1.734349230, 1.786129359, 1.788188224, 1.820273579, 1.717051330, 1.794336563, 1.969385639, 1.768035323.

All observed decoder improvements are **unpaired historical session differences**, and the decoder has *zero source changes* in either variant. They are not evidence of a decoder algorithm gain. H1/H2 also had nominal +1.24/+1.27% encoding time vs the latest accepted historical reference. The historical timing baseline was not rerun as instructed; host/task noise means these shifts should not be treated as proof of slowdown either. Crucially, neither variant established a validated JLZ encoder time reduction despite component gains.

## Checks and disposition

The full native Go uncached test suite (`go test -count=1 ./...`), `go vet ./...`, targeted `go test -race -count=1 ./bwt`, and Windows/amd64 CGO-disabled cross-build **PASS** on H1. H1/H2 identical existing decoder means no archive backward-compatibility changes; existing compatibility and corrupt-input suite passed on H1. No separate corpus-wide legacy archive decoding or native Windows executable tests were performed, and no long-running full race-corpus test was attempted. Samples and ten-run candidate archives were temporary only.

**REJECT both concrete forms:** while removing an unnecessary serialized-memory buffer reduces component allocations and can shorten BWT transform time, it adds implementation-specific `io.Writer` record dispatch, additional validation logic and dependence on Go `Index.Write` call boundaries without establishing a meaningful whole-codec encode win or archive size gain. Accepted original source, root PGO, archive format and decoder remain unchanged. This does **not** reject a future truly direct/native suffix-array-to-BWT algorithm; avoid spending further time on adjacent varint micro-optimizations. Instead pursue source-grounded larger gains in mathematical source redundancy where JLZ still trails best compressors.
