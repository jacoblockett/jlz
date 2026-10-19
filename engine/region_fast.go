package engine

import (
	"crypto/sha256"
	"math"
)

// The regional pass is only entered after an actual block benchmark. Reuse
// existing native match tokens and independently coded streams, skipping
// older serialized alternatives that carry redundant metadata for this path.
func encodeMathRegionBlockFast(block []byte) (byte, []byte, error) {
	return encodeMathRegionBlockHistory(block, nil)
}

// huffmanCannotBeat certifies when no raw/RLE/Huffman native stream could
// beat an already encoded regional candidate. Huffman payload length is at
// least the zero-order Shannon entropy, before adding its sparse codebook.
// Prediction and match candidates are still encoded in full and compared
// using actual complete sizes. The margin covers float rounding.
func huffmanCannotBeat(block []byte, selectedSize int) bool {
	if selectedSize >= len(block) || len(block) < 512 {
		return false
	}
	var counts [256]int
	for _, b := range block {
		counts[b]++
	}
	distinct := 0
	lowerBits := 0.0
	for _, c := range counts {
		if c > 0 {
			distinct++
			lowerBits += float64(c) * math.Log2(float64(len(block))/float64(c))
		}
	}
	if distinct <= 1 {
		return false
	} // RLE could win.
	// Deliberately exclude framing bytes for a conservative lower bound.
	lowerBytes := int(lowerBits/8) + distinct*2 + varSize(uint64(distinct))
	return selectedSize+32 < lowerBytes
}

func encodeMathRegionBlockHistory(block []byte, history []byte) (byte, []byte, error) {
	// Determine match and predictor results first, then avoid a redundant
	// native Huffman pass only when an entropy bound proves it cannot win.
	var predPayload, matchPayload []byte
	var predMode, matchMode byte
	if len(block) >= 128 {
		pred, err := fastPredictive(block, sha256.Sum256(block))
		if err != nil {
			return 0, nil, err
		}
		if pred != nil {
			predPayload, err = innerPayload(pred, modePredictive, len(block))
			if err != nil {
				return 0, nil, err
			}
			predMode = blockPredictive
		}
	}
	if len(block) >= 64 && (len(history) > 0 || stage3MatchPotential(block)) {
		var tokens []byte
		var err error
		if len(history) > 0 {
			whole := make([]byte, len(history)+len(block))
			copy(whole, history)
			copy(whole[len(history):], block)
			tokens, err = tokenizeMatchV2Range(whole, len(history))
		} else {
			tokens, err = tokenizeMatchV2(block)
		}
		if err != nil {
			return 0, nil, err
		}
		if len(tokens) > 0 {
			if len(history) > 0 {
				matchPayload, err = packCrossRegionTokens(block, tokens)
				matchMode = blockCrossPacked
			} else {
				matchMode, matchPayload, err = packedRegionVariantsShared(block, tokens)
			}
			if err != nil {
				return 0, nil, err
			}
		}
	}
	bestCandidateSize := len(block)
	if predPayload != nil {
		bestCandidateSize = min(bestCandidateSize, len(predPayload))
	}
	if matchPayload != nil {
		bestCandidateSize = min(bestCandidateSize, len(matchPayload))
	}
	mode, best := stage3Raw, block
	if !huffmanCannotBeat(block, bestCandidateSize) {
		packed, err := nativePackBytes(block)
		if err != nil {
			return 0, nil, err
		}
		if len(packed) < len(best) {
			mode, best = blockNativeBytes, packed
		}
	}
	// Preserve the original tie-breaking order native -> prediction -> match.
	if predPayload != nil && len(predPayload) < len(best) {
		mode, best = predMode, predPayload
	}
	if matchPayload != nil && len(matchPayload) < len(best) {
		mode, best = matchMode, matchPayload
	}
	return mode, best, nil
}
