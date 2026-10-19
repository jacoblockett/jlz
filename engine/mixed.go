package engine

// Stage 3: independently coded blocks and matched-literal prediction.
// Block references are local to their block; no independent dictionaries,
// checksums or parser states are serialized. Global SHA-256 remains mandatory.

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"io"
	"math"
)

const (
	modeMixed            byte = 12
	stage3BlockSize           = 256 * 1024
	stage3Raw            byte = 0
	blockNativeBytes     byte = 1
	blockPredictive      byte = 2
	blockNative          byte = 3
	blockMatchedLiteral  byte = 4
	blockLegacyMatch     byte = 5
	blockOldDictionary   byte = 6
	blockMatchPacked     byte = 7
	blockMatchPackedPred byte = 8
)

// An inexpensive distributed check for changing data distributions. This is
// only a work-limiting heuristic: a false negative retains Stage 2's codec.
// Screen block mode switching with three distributed local predictor profiles.
// Only invest in a full extra compression pass if different regions strongly
// favor different predictors, or if a poorly compressed stream has highly
// divergent byte distributions. This is deliberately conservative.
func stage3SampleUseful(input []byte, bestSize int) bool {
	if len(input) < 2*stage3BlockSize {
		return false
	}
	const sampleSize = 4096
	var hist [3][256]int
	positions := [3]int{0, len(input)/2 - sampleSize/2, len(input) - sampleSize}
	type winner struct {
		kind, lag int
		gain      float64
	}
	var winners [3]winner
	entropy := func(h *[256]int, count int) float64 {
		sum := 0.0
		for _, c := range h {
			if c > 0 {
				p := float64(c) / float64(count)
				sum -= p * math.Log2(p)
			}
		}
		return sum
	}
	for k, at := range positions {
		chunk := input[at : at+sampleSize]
		for _, v := range chunk {
			hist[k][v]++
		}
		base := entropy(&hist[k], sampleSize)
		for _, lag := range []int{1, 2, 4, 8, 16} {
			for kind := 1; kind <= 2; kind++ {
				var h [256]int
				for i := lag; i < sampleSize; i++ {
					a, b := chunk[i], chunk[i-lag]
					if kind == 1 {
						h[a^b]++
					} else {
						h[a-b]++
					}
				}
				gain := base - entropy(&h, sampleSize-lag)
				if gain > winners[k].gain {
					winners[k] = winner{kind, lag, gain}
				}
			}
		}
	}
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 3; j++ {
			a, b := winners[i], winners[j]
			if a.gain >= 0.8 && b.gain >= 0.8 && (a.kind != b.kind || a.lag != b.lag) {
				return true
			}
		}
	}
	if bestSize*10 < len(input)*7 {
		return false
	} // only when baseline is relatively expensive
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 3; j++ {
			distance := 0
			for v := 0; v < 256; v++ {
				d := hist[i][v] - hist[j][v]
				if d < 0 {
					d = -d
				}
				distance += d
			}
			if distance > sampleSize*7/5 {
				return true
			}
		}
	}
	return false
}

func innerPayload(archive []byte, mode byte, n int) ([]byte, error) {
	if len(archive) < 5 || string(archive[:4]) != magic || archive[4] != mode {
		return nil, errors.New("invalid internal archive")
	}
	r := bytes.NewReader(archive[5:])
	length, err := readVar(r)
	if err != nil {
		return nil, err
	}
	if length != uint64(n) {
		return nil, errors.New("internal length mismatch")
	}
	if r.Len() < 32 {
		return nil, io.ErrUnexpectedEOF
	}
	if _, err = r.Seek(32, io.SeekCurrent); err != nil {
		return nil, err
	}
	return archive[len(archive)-r.Len():], nil
}

// stage3MatchPotential is only a CPU prefilter, not an archive-size model.
// Short byte matches, XOR and additive equivalences all contribute evidence;
// a good native byte result does not disqualify an independently better match.
// We cap probes with a sparse hash table; exact sizes decide after selection.
func stage3MatchPotential(data []byte) bool {
	const step = 8
	if len(data) < 1024 {
		return true
	}
	// Probe 8-byte exact repeats and 8-byte XOR/additive-difference repeats.
	const buckets = 1 << 15
	type key struct {
		hash uint32
		kind byte
	}
	seen := make(map[key]int, min(len(data)/step, buckets))
	hits := 0
	for i := 0; i+16 <= len(data); i += step {
		for kind := byte(0); kind < 3; kind++ {
			k := key{hash: matchFingerprint(data, i, kind), kind: kind}
			if prev, ok := seen[k]; ok && i-prev >= 16 {
				// Verify beyond the short fingerprint to reject hash collisions.
				m := byte(0)
				if kind == 1 {
					m = data[i] ^ data[prev]
				}
				if kind == 2 {
					m = data[i] - data[prev]
				}
				valid := true
				for j := 0; j < 16; j++ {
					v := data[prev+j]
					if kind == 1 {
						v ^= m
					} else if kind == 2 {
						v += m
					}
					if data[i+j] != v {
						valid = false
						break
					}
				}
				if valid {
					hits++
					if hits >= 3 {
						return true
					}
				}
			}
			seen[k] = i
		}
	}
	return false
}

func encodeStage3Block(block []byte) (byte, []byte, error) {
	return encodeStage3BlockWithDict(block, nil)
}

func encodeStage3BlockWithDict(block []byte, history []byte) (byte, []byte, error) {
	bestMode, best := stage3Raw, block
	packed, err := nativePackBytes(block)
	if err != nil {
		return 0, nil, err
	}
	if len(packed) < len(best) {
		bestMode, best = blockNativeBytes, packed
	}
	// No cross-block external-compressor dictionary path in JLZ.
	if len(block) >= 128 {
		pred, err := fastPredictive(block, sha256.Sum256(block))
		if err != nil {
			return 0, nil, err
		}
		if pred != nil {
			payload, err := innerPayload(pred, modePredictive, len(block))
			if err != nil {
				return 0, nil, err
			}
			if len(payload) < len(best) {
				bestMode, best = blockPredictive, payload
			}
		}
	}
	// A native match-token block is a first-class candidate. Previously,
	// only native match streams were considered, although the legacy stream
	// repeatedly won during the match-cost audit. Use the very same tokenization
	// for all match encoders; only the final representations compete.
	// Bounded match-potential screening is a fixed part of JLZ.
	if len(block) >= 64 && stage3MatchPotential(block) {
		tokens, err := tokenizeMatchV2(block)
		if err != nil {
			return 0, nil, err
		}
		if len(tokens) > 0 {
			digest := sha256.Sum256(block)
			legacy, err := encodeMatchV2Tokens(block, digest, tokens)
			if err != nil {
				return 0, nil, err
			}
			payload, err := innerPayload(legacy, modeMatchV2, len(block))
			if err != nil {
				return 0, nil, err
			}
			if len(payload) < len(best) {
				bestMode, best = blockLegacyMatch, payload
			}
			for _, kind := range []struct {
				matched bool
				mode    byte
			}{{false, blockNative}, {true, blockMatchedLiteral}} {
				arc, err := encodeMatchNativeVariant(block, digest, tokens, kind.matched)
				if err != nil {
					return 0, nil, err
				}
				nativeMode := modeMatchNative
				if kind.matched {
					nativeMode = modeMatchNativePred
				}
				payload, err := innerPayload(arc, nativeMode, len(block))
				if err != nil {
					return 0, nil, err
				}
				if len(payload) < len(best) {
					bestMode, best = kind.mode, payload
				}
			}
			// Native stream packing shares the same parser and competes on encoded size.
			for _, v := range []struct {
				predicted  bool
				kind       byte
				standalone byte
			}{
				{false, blockMatchPacked, modeMatchPacked},
				{true, blockMatchPackedPred, modeMatchPackedPred},
			} {
				arc, e := encodeMatchPackedVariant(block, digest, tokens, v.predicted)
				if e != nil {
					return 0, nil, e
				}
				payload, e := innerPayload(arc, v.standalone, len(block))
				if e != nil {
					return 0, nil, e
				}
				if len(payload) < len(best) {
					bestMode, best = v.kind, payload
				}
			}
		}
	}
	return bestMode, best, nil
}

// Return nil unless the complete mixed archive is smaller than the current
// best. Fixed blocks, one global checksum, no changes to legacy mode payloads.
func encodeMixedBlocks(input []byte, digest [32]byte, bestSize int) ([]byte, error) {
	var out bytes.Buffer
	out.WriteString(magic)
	out.WriteByte(modeMixed)
	putVar(&out, uint64(len(input)))
	out.Write(digest[:])
	for at := 0; at < len(input); at += stage3BlockSize {
		stop := min(len(input), at+stage3BlockSize)
		chunk := input[at:stop]
		var history []byte
		if at > 0 {
			history = input[max(0, at-32768):at]
		}
		mode, payload, err := encodeStage3BlockWithDict(chunk, history)
		if err != nil {
			return nil, err
		}
		out.WriteByte(mode)
		putVar(&out, uint64(len(chunk)))
		putVar(&out, uint64(len(payload)))
		out.Write(payload)
		if out.Len() >= bestSize {
			return nil, nil
		}
	}
	return out.Bytes(), nil
}

func decodeMixedBlocks(r *bytes.Reader, n uint64) ([]byte, error) {
	if n >= 2*1024*1024 {
		return decodeMixedBlocksParallel(r, n)
	}
	return decodeMixedBlocksSequential(r, n)
}

func decodeMixedBlocksSequential(r *bytes.Reader, n uint64) ([]byte, error) {
	if n > uint64(math.MaxInt) {
		return nil, errors.New("mixed archive too large")
	}
	output := make([]byte, 0, int(n))
	for uint64(len(output)) < n {
		mode, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		size, err := readVar(r)
		if err != nil {
			return nil, err
		}
		length, err := readVar(r)
		if err != nil {
			return nil, err
		}
		if size == 0 || size > stage3BlockSize || size > n-uint64(len(output)) || length > uint64(r.Len()) {
			return nil, errors.New("invalid mixed block size")
		}
		if length > size*8+1024*1024 {
			return nil, errors.New("mixed payload exceeds bound")
		}
		content := make([]byte, int(length))
		if _, err = io.ReadFull(r, content); err != nil {
			return nil, err
		}
		var restored []byte
		cr := bytes.NewReader(content)
		switch mode {
		case stage3Raw:
			if length != size {
				return nil, errors.New("wrong mixed raw length")
			}
			restored = content
		case blockNativeBytes:
			restored, err = nativeUnpackBytes(content, int(size))
			cr.Seek(0, io.SeekEnd)
		case blockOldDictionary:
			return nil, errors.New("external dictionary mode is not part of native-only format")
		case blockPredictive:
			restored, err = decodePredictive(cr, size)
		case blockLegacyMatch:
			restored, err = decodeMatchV2(cr, size, int64(size))
		case blockNative:
			restored, err = decodeMatchNative(cr, size)
		case blockMatchedLiteral:
			restored, err = decodeMatchNativeVariant(cr, size, true)
		case blockMatchPacked, blockMatchPackedPred:
			restored, err = decodeMatchPacked(cr, size, mode == blockMatchPackedPred)
		case blockMatchPackedSplit, blockMatchPackedSplitPred:
			restored, err = decodePackedWithLengthMode(cr, size, mode == blockMatchPackedSplitPred, nil, true)
		default:
			return nil, errors.New("unknown mixed block mode")
		}
		if err != nil {
			return nil, err
		}
		if uint64(len(restored)) != size || (mode != stage3Raw && cr.Len() != 0) {
			return nil, errors.New("mixed block length mismatch or trailing data")
		}
		output = append(output, restored...)
	}
	if r.Len() != 0 {
		return nil, errors.New("trailing mixed archive data")
	}
	return output, nil
}
