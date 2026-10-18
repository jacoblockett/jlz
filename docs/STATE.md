# State and next steps

**Last updated:** 2026-10-18. **Status:** Repository migration to canonical source and JLS-compatible delivery workflow. **Approved source:** JLZ Round 15. **No compression experiment is authorized as a replacement beyond Round 15.**

## Accepted performance

Silesia 12 files, 211,938,580 original bytes. Round 15: **66,797,408 archive bytes**, **6.942 s median encode time** (ten paired in-memory corpus runs). Round 12: same bytes, 7.331 s paired median in that comparison. Historical external reference bytes: Zstd-3 66,216,569; XZ-6 49,233,340. Their historical CLI times, 1.256 s and 116.989 s respectively, are not controlled comparisons against JLZ timings.

## Architecture as shipped

Single CLI encoder/decoder in Go, native entropy coding, bounded mathematical match and exact-copy processing, grouped regions with independent decode groups, adaptive worker selection based on usable CPUs, logarithmic distance slots, matching heuristics, and automatically chosen representations. Reconstruct original bytes exactly. `JLZ1` container with decoder compatibility rules documented in source tests. Build remains pure Go, no dependencies beyond stdlib.

## Immediate research direction

Seek a unified match/token/entropy architecture that exposes larger real-file redundancy at acceptable CPU cost, not deeper versions of the same costly matcher. Any new idea begins as a separate branch/prototype with screening on full representative files, then ten paired whole-corpus runs if promising. Examine possible token-stream entropy modeling improvements without duplicating full coding passes.

## Stop signs / pitfalls

- DEFLATE contamination invalidated historical MDC claims. See `docs/research/attempted-ideas.md`.
- Round 11 deep search saved 0.88% bytes but was 11.01% slower, rejected.
- Round 13 sparse history saved 0.63% but was 5.75% slower, rejected.
- Round 14 byte-pair grammar saved 0.12% but was about 10% slower in screening and had a unit-test mismatch, rejected.
- Round 15 skipped redundant coding passes and achieved 5.30% speed improvement with byte-identical archives. This is the current accepted result.

## Handoff checklist

Before ending work, amend the current state, log each experiment in a dated immutable report and update the research index. Keep benchmark methods and raw run records under `docs/reports/` (text/CSV), never compiled binaries. In the user-facing summary, clearly distinguish screening from completed ten-run results. No research corpus redistribution without rights to do so.
