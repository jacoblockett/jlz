package engine

import (
	"bytes"
	"encoding/binary"
)

// Two independent exact-byte hash anchors, one short and one long. Index
// entries are pos+1, so zero means empty. This is an original, optional
// exact-match parser using existing JLZ match tokens, never a codec backend.
// Current callers pass at most one 256-KiB block plus 128-KiB local history.
func tokenizeDualAnchor(input []byte, start int) ([]byte, error) {
	n := len(input)
	if n-start < 64 {
		return nil, nil
	}
	const hbits = 16
	short := make([]uint32, 1<<hbits)
	long := make([]uint32, 1<<hbits)
	hash4 := func(at int) uint32 {
		k := binary.LittleEndian.Uint32(input[at:])
		return (k * 2654435761) >> (32 - hbits)
	}
	hash8 := func(at int) uint32 {
		k := binary.LittleEndian.Uint64(input[at:])
		return uint32((k * 0x9e3779b185ebca87) >> (64 - hbits))
	}
	// Retain up to 128KiB of already decoded history without an external dictionary.
	for at := max(0, start-crossHistorySize); at+8 <= start; at++ {
		short[hash4(at)] = uint32(at + 1)
		long[hash8(at)] = uint32(at + 1)
	}
	var out bytes.Buffer
	var recent [3]int
	literalStart := start
	refs := 0
	lastHit := start
	for at := start; at < n; {
		if at+8 > n {
			at++
			continue
		}
		h4, h8 := hash4(at), hash8(at)
		a, b := short[h4], long[h8]
		short[h4], long[h8] = uint32(at+1), uint32(at+1)
		var choice matchChoice
		for _, v := range []uint32{b, a} {
			if v == 0 {
				continue
			}
			from := int(v) - 1
			dist := at - from
			if dist <= 0 || dist > matchMaxDistance || from+8 > n {
				continue
			}
			// A 4-byte hash collision is never trusted as an exact match.
			if !bytes.Equal(input[at:at+5], input[from:from+5]) {
				continue
			}
			lim := min(n-at, matchMaximum)
			size := 5
			for size+8 <= lim && binary.LittleEndian.Uint64(input[at+size:]) == binary.LittleEndian.Uint64(input[from+size:]) {
				size += 8
			}
			for size < lim && input[at+size] == input[from+size] {
				size++
			}
			c := scoreMatchChoice(matchChoice{length: size, distance: dist}, recent)
			if c.score > choice.score {
				choice = c
			}
		}
		if choice.score <= 0 {
			step := 1 + (at-lastHit)/4096
			if step > 16 {
				step = 16
			}
			at += step
			continue
		}
		// Backward extension absorbs pending literals while preserving exact output.
		for at > literalStart && at > choice.distance && choice.length < matchMaximum && input[at-1] == input[at-choice.distance-1] {
			at--
			choice.length++
		}
		if at > literalStart {
			out.WriteByte(matchLiteral)
			putVar(&out, uint64(at-literalStart))
			out.Write(input[literalStart:at])
		}
		opcode := matchCopy
		if choice.repeat >= 0 {
			opcode = matchRep0 + byte(choice.repeat)
		}
		out.WriteByte(opcode)
		if choice.repeat < 0 {
			putVar(&out, uint64(choice.distance))
		}
		putVar(&out, uint64(choice.length))
		updateMatchRecent(&recent, choice.distance)
		at += choice.length
		literalStart = at
		lastHit = at
		refs++
	}
	if literalStart < n {
		out.WriteByte(matchLiteral)
		putVar(&out, uint64(n-literalStart))
		out.Write(input[literalStart:])
	}
	if refs == 0 {
		return nil, nil
	}
	return out.Bytes(), nil
}

// Select the existing algebraic matcher when its decoder-independent probe
// detects a verified transformed relation; otherwise use a cheaper exact
// parser. Both produce the same existing JLZ packed-stream token syntax.
func tokenizeRegionHybrid(input []byte, start int) ([]byte, error) {
	if likelyTransformedMatches(input) {
		return tokenizeMatchV2Range(input, start)
	}
	return tokenizeDualAnchor(input, start)
}
