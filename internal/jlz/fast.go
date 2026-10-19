package jlz

// JLZ single encoder. Internal representations compete automatically on
// complete archive size; there are no levels or user-facing algorithm switches.

import (
	"bytes"

	"crypto/sha256"
	"encoding/hex"
	"math"
)

func compressFast(input []byte) ([]byte, summary, error) {
	digest := sha256.Sum256(input)
	s := summary{original: len(input), sha: hex.EncodeToString(digest[:]), raw: true, method: "raw"}
	// Delay allocating/copying the entire raw input until we actually know
	// it won. On full files this otherwise duplicates tens of MiB for nothing.
	rawSize := len(magic) + 1 + varSize(uint64(len(input))) + len(digest) + len(input)
	makeRaw := func() []byte {
		var raw bytes.Buffer
		raw.Grow(rawSize)
		raw.WriteString(magic)
		raw.WriteByte(modeRaw)
		putVar(&raw, uint64(len(input)))
		raw.Write(digest[:])
		raw.Write(input)
		return raw.Bytes()
	}
	var best []byte
	bestSize := rawSize
	if len(input) == 0 {
		best = makeRaw()
		s.archive = len(best)
		return best, s, nil
	}
	// First evaluate strong regional candidates before allocating and coding a
	// redundant whole-file baseline. The conservative preflight compares the
	// same independent regions using both native and mathematical coding.
	// Any ambiguous case follows the original exhaustive selection path.
	if len(input) > 2*1024*1024 {
		regional, ok, err := tryEarlyRegional(input, digest, rawSize)
		if err != nil {
			return nil, s, err
		}
		if ok {
			s.archive = len(regional)
			s.raw = false
			s.method = "cross-region-candidate"
			return regional, s, nil
		}
	}
	var base []byte
	var err error
	parallel := len(input) >= parallelNativeMinSize
	if parallel {
		base, err = encodeParallelNative(input, digest, 0)
	} else {
		base, err = encodeFastNative(input, digest)
	}
	if err != nil {
		return nil, s, err
	}
	if len(base) < bestSize {
		best = base
		bestSize = len(base)
		if parallel {
			s.method = "parallel-native"
		} else {
			s.method = "native-bytes"
		}
		s.raw = false
	}
	// Reuse the cheap structural probe for whole-file matching and mixed-block
	// screening. A positive result authorizes further evaluation, never selects
	// a lossy or unverified archive.
	structuralChecked := false
	structuralEvidence := false
	structuralUseful := func() bool {
		if !structuralChecked {
			structuralEvidence = stage3StructuralSampleUseful(input)
			structuralChecked = true
		}
		return structuralEvidence
	}
	// A bounded hash index tests exact and mathematically transformed copies.
	// Large-file searches are gated cheaply before the expensive tokenization.
	// Whole-file mathematical parsing has superlinear worst-case costs; it is
	// never launched automatically above 1 MiB in the native-only research build.
	// This limit is part of the single JLZ compression algorithm.
	if len(input) >= 64 && len(input) <= 1024*1024 {
		// For smaller files, check whether the complete matcher is economical.
		// Larger eligible inputs also require cheap structural evidence.
		eligible := false
		if len(input) <= 384*1024 || structuralUseful() {
			var err error
			eligible, err = matchSampleUseful(input)
			if err != nil {
				return nil, s, err
			}
		}
		if eligible {
			arc, method, err := encodeMatchStage2(input, digest)
			if err != nil {
				return nil, s, err
			}
			if len(arc) > 0 && len(arc) < bestSize {
				best, bestSize, s.method, s.raw = arc, len(arc), method, false
			}
		}
		// Keep the proven small-input matcher as a cheap alternative rather
		// than deleting a working mathematical representation. This also
		// protects archive-size behavior while the newer parser matures.
		if len(input) >= 128 && len(input) <= 384*1024 {
			arc, err := fastAdaptive(input, digest)
			if err != nil {
				return nil, s, err
			}
			if len(arc) > 0 && len(arc) < bestSize {
				best, bestSize, s.method, s.raw = arc, len(arc), "adaptive-fast", false
			}
		}
	}
	// Sampled histograms exclude unpromising transforms before full encoding.
	if len(input) >= 128 {
		arc, err := fastPredictive(input, digest)
		if err != nil {
			return nil, s, err
		}
		if len(arc) > 0 && len(arc) < bestSize {
			best, bestSize, s.method, s.raw = arc, len(arc), "predictive-fast", false
		}
	}
	// Native multi-mode regional search is bounded by file size in normal
	// operation. There is no user-selectable exhaustive mode.
	if len(input) >= 2*stage3BlockSize && len(input) <= 2*1024*1024 {
		if stage3SampleUseful(input, bestSize) || structuralUseful() {
			arc, err := encodeMixedBlocks(input, digest, bestSize)
			if err != nil {
				return nil, s, err
			}
			if len(arc) > 0 && len(arc) < bestSize {
				best, bestSize, s.method, s.raw = arc, len(arc), "mixed-blocks", false
			}
		}
	}
	// Large files: perform fixed-work, independently decodable 256 KiB
	// mathematical region searches only when representative aligned blocks
	// demonstrate a meaningful advantage over the whole-file candidate.
	// This is part of the single automatic JLZ encoder, never a user flag.
	if len(input) > 2*1024*1024 {
		useRegions, samples := regionalMathUseful(input, bestSize)
		if useRegions {
			arc, err := encodeParallelMathRegions(input, digest, bestSize, samples)
			if err != nil {
				return nil, s, err
			}
			if len(arc) > 0 && len(arc) < bestSize {
				best, bestSize, s.method, s.raw = arc, len(arc), "cross-region-candidate", false
			}
		}
	}
	// A different structural representation can exploit interleaved fields
	// missed by whole-file native Huffman and mathematical references.
	// Actual full-archive bytes always decide whether to use it.
	if len(input) > 2*1024*1024 {
		if stride := sampledLaneStride(input, bestSize); stride != 0 {
			candidate, e := encodeLaneArchive(input, digest, bestSize, stride)
			if e != nil {
				return nil, s, e
			}
			if len(candidate) != 0 && len(candidate) < bestSize {
				best, bestSize, s.method, s.raw = candidate, len(candidate), "lane-regions", false
			}
		}
	}
	if best == nil {
		best = makeRaw()
	}
	s.archive = len(best)
	return best, s, nil
}

type candidateTransform struct {
	typ     byte
	lag     int
	benefit float64
}

func fastPredictive(input []byte, digest [32]byte) ([]byte, error) {
	n := len(input)
	sample := n
	if sample > 32768 {
		sample = 32768
	}
	stride := n / sample
	if stride < 1 {
		stride = 1
	}
	entropy := func(t byte, lag int) float64 {
		var hist [256]int
		count := 0
		for p := lag; p < n && count < sample; p += stride {
			v := input[p]
			switch t {
			case predXor:
				v ^= input[p-lag]
			case predSub:
				v -= input[p-lag]
			}
			hist[v]++
			count++
		}
		var h float64
		for _, c := range hist {
			if c != 0 {
				prob := float64(c) / float64(count)
				h -= prob * math.Log2(prob)
			}
		}
		return h
	}
	base := entropy(predNone, 0)
	best := candidateTransform{benefit: 0}
	for _, lag := range []int{1, 2, 4, 8, 16} {
		if lag >= n {
			continue
		}
		for _, typ := range []byte{predXor, predSub} {
			gain := base - entropy(typ, lag)
			if gain > best.benefit {
				best = candidateTransform{typ, lag, gain}
			}
		}
	}
	// Empirically neutral transformations don't warrant another full
	// entropy encoding attempt; save both CPU and allocation pressure.
	if best.benefit < 0.45 {
		return nil, nil
	}
	transformed := transform(input, best.typ, best.lag)
	payload, err := fastPredictiveNative(transformed)
	if err != nil {
		return nil, err
	}
	var archive bytes.Buffer
	archive.WriteString(magic)
	archive.WriteByte(modePredictive)
	putVar(&archive, uint64(n))
	archive.Write(digest[:])
	archive.WriteByte(best.typ)
	putVar(&archive, uint64(best.lag))
	archive.WriteByte(codecNative)
	archive.Write(payload)
	return archive.Bytes(), nil
}

// Uses only the in-project byte packer. No effort-level sampler remains:
// the raw/RLE/Huffman stream format has no external backend or effort levels.
func fastNativeBytes(data []byte) ([]byte, error)      { return nativePackBytes(data) }
func fastPredictiveNative(data []byte) ([]byte, error) { return nativePackBytes(data) }

func encodeFastNative(input []byte, digest [32]byte) ([]byte, error) {
	return encodeNativeBytes(input, digest)
}
