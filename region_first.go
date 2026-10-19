package main

// tryEarlyRegional avoids encoding the full file with unrelated candidates
// when three distributed regions give strong independent evidence that the
// selected local/regional match representation is far better than native.
// This is an encoder-only work heuristic. The selected bytes are still the
// original standalone grouped-history format with the same decoder.
func tryEarlyRegional(input []byte, digest [32]byte, rawSize int) ([]byte, bool, error) {
	count := (len(input) + stage3BlockSize - 1) / stage3BlockSize
	if count < 12 {
		return nil, false, nil
	}
	var samples mathRegionSamples
	samples.positions = [3]int{0, count / 2, count - 1}
	regionalBytes, nativeBytes, rawBytes := 0, 0, 0
	for i, idx := range samples.positions {
		at := idx * stage3BlockSize
		block := input[at:min(at+stage3BlockSize, len(input))]
		mode, payload, err := encodeMathRegionBlockFast(block)
		if err != nil {
			return nil, false, err
		}
		samples.results[i] = mathRegionResult{mode: mode, payload: payload}
		native, err := nativePackBytes(block)
		if err != nil {
			return nil, false, err
		}
		// All sampled blocks must support the same clear direction.
		if (len(payload)+32)*10 >= len(native)*9 {
			return nil, false, nil
		}
		regionalBytes += len(payload)
		nativeBytes += len(native)
		rawBytes += len(block)
	}
	// Broader preflight still demands a substantial sampled advantage for the
	// mathematical representation, then lets the measured BWT alternative
	// compete instead of requiring extreme (>50%) input reduction.
	if regionalBytes*5 >= nativeBytes*4 || regionalBytes*5 >= rawBytes*4 {
		return nil, false, nil
	}
	// A sampled alternative block-sorting path competes with the selected
	// mathematical regions. Only run the full alternate after a 5% margin.
	bwtOK, bwtSamples, err := sampleBWTPreference(input, samples)
	if err != nil {
		return nil, false, err
	}
	if bwtOK {
		// A mathematically distinct additive-progression representation can
		// replace suffix sorting when actual distributed sample payloads
		// reliably beat BWT. It writes a complete original-byte JLZ1 archive.
		if additiveWordPreference(input, bwtSamples) {
			additive, e := encodeParallelAdditiveArchive(input, digest)
			if e != nil {
				return nil, false, e
			}
			if len(additive) > 0 && len(additive)*20 < len(input)*19 {
				return additive, true, nil
			}
		}
		bwtArchive, err := encodeParallelBWTRegions(input, digest, rawSize, bwtSamples)
		if err != nil {
			return nil, false, err
		}
		if len(bwtArchive) > 0 && len(bwtArchive)*20 < len(input)*19 {
			return bwtArchive, true, nil
		}
	}
	regional, err := encodeParallelMathRegions(input, digest, rawSize, samples)
	if err != nil {
		return nil, false, err
	}
	if len(regional) == 0 || len(regional)*20 >= len(input)*19 {
		return nil, false, nil
	}
	return regional, true, nil
}
