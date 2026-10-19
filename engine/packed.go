package engine

import (
	"bytes"
	"math"
)

// Share parsed tokens and packed metadata across two final match encodings.
// Equal-sized candidates keep the plain mode; archive bytes are unchanged.
func packedRegionVariantsShared(input []byte, tokens []byte) (byte, []byte, error) {
	st, err := splitMatchTokens(tokens, len(input))
	if err != nil {
		return 0, nil, err
	}
	var head bytes.Buffer
	putVar(&head, uint64(st.ncommands))
	// Pack unchanged token columns once, independently of which length model wins.
	packedOp := packFiniteState(st.op)
	packedLengths := packFiniteState(st.lengths)
	packedDistances := packVarintStream(st.distances)
	packedTransforms := packStream(st.transforms)
	head.Write(packedOp)
	head.Write(packedLengths)
	head.Write(packedDistances)
	head.Write(packedTransforms)
	packLiterals := func(lit []byte) []byte {
		if len(lit)*4 > len(input)*3 {
			return packStream(lit)
		}
		return packLiteralStream(lit)
	}
	p0 := packLiterals(st.literals)
	selected := make([]byte, 0, head.Len()+len(p0))
	selected = append(selected, head.Bytes()...)
	selected = append(selected, p0...)
	mode := blockMatchPacked
	// A second way to serialize the same token sequence, with separate
	// Huffman models for literal and reference lengths. Choose only if the
	// complete stream representation is actually smaller.
	if splitHead, ok := packSplitLengthHeader(st, packedOp, packedDistances, packedTransforms); ok && len(splitHead) < head.Len() {
		head.Reset()
		head.Write(splitHead)
		selected = append(selected[:0], head.Bytes()...)
		selected = append(selected, p0...)
		mode = blockMatchPackedSplit
	}
	residual, err := matchedLiteralResiduals(tokens, input)
	if err != nil {
		return 0, nil, err
	}
	// The second context-coded literal stream is usually much worse, but
	// coding it fully costs substantial CPU. Test a short prefix of BOTH
	// streams before investing in a full alternative. A conservative 256-byte
	// margin retains contenders that may improve after the sample.
	// This is an encoder selection heuristic; the exact decoder and wire
	// format remain unchanged.
	if !considerMatchedLiteralResidual(st.literals, residual) {
		return mode, selected, nil
	}
	p1 := packLiterals(residual)
	if len(p1) < len(p0) {
		selected = append(selected[:0], head.Bytes()...)
		selected = append(selected, p1...)
		if mode == blockMatchPackedSplit {
			mode = blockMatchPackedSplitPred
		} else {
			mode = blockMatchPackedPred
		}
	}
	return mode, selected, nil
}

// A bounded preflight for the second full arithmetic pass. The selected
// archive remains a complete, lossless JLZ1 representation. This is a
// speed-versus-search-depth heuristic: on highly nonstationary inputs a
// late improvement may be missed if the first 8 KiB is unrepresentative.
func considerMatchedLiteralResidual(literals, residual []byte) bool {
	if len(literals) != len(residual) {
		return false
	}
	// A necessary-looking (but heuristic) signal for profitable masked-literal
	// coding is a noticeably more concentrated symbol distribution. Computing
	// two 256-bin histograms is far cheaper than a second arithmetic pass.
	// Tiny streams retain the old exact comparison path. This conservative
	// histogram heuristic can miss rare context-only improvements; selected
	// candidates are still verified against complete encoded byte sizes.
	if len(residual) >= 8192 {
		var plain, predicted [256]int
		for i, b := range literals {
			plain[b]++
			predicted[residual[i]]++
		}
		entropy := func(freq *[256]int) float64 {
			var h float64
			inv := 1 / float64(len(literals))
			for _, n := range freq {
				if n > 0 {
					p := float64(n) * inv
					h -= p * math.Log2(p)
				}
			}
			return h
		}
		if entropy(&plain)-entropy(&predicted) < 0.08 {
			return false
		}
	}
	if len(residual) < 48*1024 {
		return true
	}
	const preview = 8192
	plain := arithmeticEncode(literals[:preview])
	pred := arithmeticEncode(residual[:preview])
	return len(pred) <= len(plain)+256
}
