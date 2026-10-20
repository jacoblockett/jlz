# Current JLZ architecture (no numbered releases)

`current` is the native Go JLZ engine on `main`. `candidate` means an isolated, unapproved research implementation; none is currently integrated or benchmarked. Current state and required metrics: [STATE.md](STATE.md). Concrete failed strategies and prior outcomes: [strategy ledger](research/strategy-ledger.md). Historical numbered report filenames are only evidence links, not names of the production compressor.

## Package structure

```text
main.go                    two-command CLI: compress/decompress
engine/                    automatic selection, stream codecs, JLZ1 construction/decoding
archive/                   native wire framing and varint primitives
additive/                  reversible signed modular residual transforms
bwt/                       original-byte BWT, move-to-front/bijective zero coding, inverse
entropy/                   native adaptive arithmetic coding
lane/                      independent lane transforms and statistical estimates
match/                     native exact and XOR/ADD match finding/expansion
predictive/                native reversible lag predictions
workers/                   quota-aware bounded parallelism
engine/testdata/fixtures/  deterministic byte-exact correctness fixtures
docs/research/             strategy ledger, hypotheses, negative results and research queue
docs/reports/              dated experimental methods and historically benchmarked data
```

No third-party mature compression codec as a backend or fallback. The full JLZ engine remains one automatically selected, exact-byte lossless compressor with exactly two commands and no user levels.

## Current actual production path

- Encoder chooses among native mathematical match/predictive streams, additive/word transforms, native entropy/raw fallbacks, and context-preflight-selected BWT blocks using original-byte exact fully framed costs where candidate comparisons occur.
- **Source regions remain 256 KiB** for both mathematical matching and BWT. Four mathematical reference regions form an independently decodable history group; the original mathematical lookback is **128 KiB**, with strict decoder limits. Larger BWT regions and enlarged conventional math history were tested and rejected for slowing down both encoding/decoding.
- Existing mathematical match tokens support exact copying, XOR/ADD transformed relations, repeat distances, and separately packed opcode, length, distance, transform and literal streams.
- BWT uses a native byte transform whose forward suffix order is currently obtained by ordinary Go `index/suffixarray` and its serialized varint index representation. The selected BWT intermediate uses bijective move-to-front zero-run coding and native 12-bit finite-state entropy when shorter than legacy packing.
- For winning BWT blocks, inexpensive minimum-description-length estimates select a single optional 2-/4-part local finite-state representation, **internal BWT mode 16**. The old BWT format and prior raw/RLE/Huffman/math modes continue decoding unchanged. Older binaries that predate mode16 cannot decode new mode16 payloads. No new forward archive mode has been added by subsequent rejected experiments.
- Local compression/decoding tests, corruption checks, Windows crossbuild, corpus hash matching, exact original bytes, and historical decoder compatibility remain mandatory for accepted changes.

## Things that are deliberately NOT in production

- Dual-anchor exact-only fast math matcher that grew archive bytes.
- BWT 512 KiB/1 MiB/2 MiB source-block variants that improved size but slowed decode/encode; no new extended BWT frame accepted.
- Three-model conditional BWT entropy, dynamic BWT change-point searches, conditional rank remaps, fixed math metadata segmentation and incremental-frequency post-BWT models with disproportionate size/time costs.
- More expensive expanded conventional math cross-history, ad-hoc sparse far-file absolute-offset references and packed LF-link inversion without an accepted tri-metric advantage.
- A generic radix-prefix-doubling suffix constructor. It was byte-exact but 2–7× slower on 21 original BWT source-block diagnostic samples than the current standard-library implementation. The generic standard-library serialized suffix array remains in production until a mathematically equivalent direct-order method demonstrates a net all-metric win.

For measured provenance and limits, use the [strategy ledger](research/strategy-ledger.md) and [dated experiment reports](reports/), rather than a chronological list of numbered rounds.

## Local build and strict acceptance

```sh
go test -count=1 ./...
go vet ./...
go build -o jlz .
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o jlz.exe .
```

GitHub Actions are disabled. Only permanent validated source edits to `main`, disposable local research otherwise. No candidate may worsen the complete original 12-file Silesia **archive size, ten-run median encoding time, or ten-run median decoding time**, unless the user explicitly approves that precise tradeoff.


## Accepted layered BCJ→BWT plus Go PGO (2026-10-20)

**Current** uses reversible per-block E8/E9 32-bit little-endian relative operand normalization before existing native 256KiB BWT/MTF/finite-state coding on automatically selected math regions. Mixed-cross mode 19 holds existing BWT subtype 15 or 16. The decoder reconstructs via existing BWT inverse then branch inverse; old archive modes remain backward-decodable, while old binaries reject newly written mode 19. A root checked-in **default.pgo** profile (16,804 compressed bytes; SHA-256 3b35cda74832360a252b63c1805e025aa1f4ed7a4fe4124d9a55929a6e86289f) is automatically read by ordinary Go 1.23+ compilation. Without the measured profile the same codec was historically REJECTED for encoding regression; **do not claim the approved timings for an unprofiled build**.

Current **12-file original Silesia** accepted complete archive: **60,244,257 B**; candidate-only ten-run medians **5.322679488 s encode / 2.003659634 s decode**, 120 exact reconstructions. Previous accepted historical (not rerun): 60,585,717 B / 5.493743735 s / 2.194427893 s. Net -341,460 B, nominal -3.1138% encode, -8.6933% decode (independently timed sessions). [Actual ten-run data, method and backward-compatibility checks](reports/2026-10-20-native-bcj-bwt-pgo-accepted.md). No pending candidate and no numbered release label.
