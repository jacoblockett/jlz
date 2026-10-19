package jlz

// The v2 matcher indexes short byte fingerprints in fixed-size tables.
// It complements the original research matchers without changing their modes.
// References are distances rather than absolute source positions, can overlap,
// and share a three-entry repeat-distance cache. No index is stored in an archive.

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"io"
)

const (
	modeMatchV2  byte = 10
	matchLiteral byte = iota
	matchCopy
	matchXor
	matchAdd
	matchRep0
	matchRep1
	matchRep2
	matchRepXor
	matchRepAdd
)

const (
	matchHashBits           = 16
	matchMaxDistance        = 8 << 20
	matchMaximum            = 32 << 10
	matchMinimumExact       = 5
	matchMinimumTransformed = 12
)

type matchBucket struct {
	positions [4]uint32 // Stored as (position + 1), 0 means empty.
	next      uint32
}

type matchIndex struct {
	heads         [3][]matchBucket
	input         []byte
	indexed       int
	useTransforms bool
}

func newMatchIndex(data []byte) *matchIndex {
	return newMatchIndexTransforms(data, true)
}

func newMatchIndexTransforms(data []byte, transforms bool) *matchIndex {
	m := &matchIndex{input: data, useTransforms: transforms}
	m.heads[0] = make([]matchBucket, 1<<matchHashBits)
	if transforms {
		m.heads[1] = make([]matchBucket, 1<<matchHashBits)
		m.heads[2] = make([]matchBucket, 1<<matchHashBits)
	}
	return m
}

func matchFingerprint(data []byte, at int, typ byte) uint32 {
	if typ == 0 {
		return binary.LittleEndian.Uint32(data[at:]) ^ uint32(data[at+4])*0x9e3779b1
	}
	a := data[at]
	b, c, d, e := data[at+1], data[at+2], data[at+3], data[at+4]
	if typ == 1 {
		return uint32(b^a) | uint32(c^a)<<8 | uint32(d^a)<<16 | uint32(e^a)<<24
	}
	return uint32(b-a) | uint32(c-a)<<8 | uint32(d-a)<<16 | uint32(e-a)<<24
}

func matchHash(key uint32) uint32 {
	key ^= key >> 15
	key *= 0x2c1b3c6d
	key ^= key >> 12
	return key >> (32 - matchHashBits)
}

func (m *matchIndex) addBefore(pos int) {
	n := len(m.input)
	for m.indexed < pos {
		at := m.indexed
		m.indexed++
		if at+5 > n {
			continue
		}
		for typ := byte(0); typ < 3; typ++ {
			if !m.useTransforms && typ != 0 {
				break
			}
			bucket := &m.heads[typ][matchHash(matchFingerprint(m.input, at, typ))]
			bucket.positions[bucket.next&3] = uint32(at + 1)
			bucket.next++
		}
	}
}

type matchChoice struct {
	length, distance int
	typ, mask        byte
	repeat           int
	score            int
}

func matchCost(c matchChoice, recent [3]int) (int, int) {
	for i, distance := range recent {
		if distance > 0 && c.distance == distance {
			if c.typ == 0 {
				return 1 + varSize(uint64(c.length)), i
			}
			return 3 + varSize(uint64(c.length)), i
		}
	}
	cost := 1 + varSize(uint64(c.length)) + varSize(uint64(c.distance))
	if c.typ != 0 {
		cost++
	}
	return cost, -1
}

// Approximate the opportunity cost of leaving bytes as literals. The former
// 3-bit literal / 1-bit local-copy scores greatly underselected short exact
// matches in common real-file content. The calibrated 4 / 5 scores choose
// more profitable local references, with no on-disk format change. These are
// encoder heuristics, not mathematical guarantees or user-selectable levels.
const (
	matchLiteralValueBits   = 5
	matchLocalCopyValueBits = 6
	matchLocalCopyWindow    = 32 << 10
)

func scoreMatchChoice(c matchChoice, recent [3]int) matchChoice {
	// The scorer is on the matcher hot path. Compute the encoded byte cost
	// directly without another function call or loop over recent offsets.
	c.repeat = -1
	if c.distance == recent[0] && recent[0] != 0 {
		c.repeat = 0
	} else if c.distance == recent[1] && recent[1] != 0 {
		c.repeat = 1
	} else if c.distance == recent[2] && recent[2] != 0 {
		c.repeat = 2
	}
	cost := 1 + varSize(uint64(c.length))
	if c.repeat < 0 {
		cost += varSize(uint64(c.distance))
	} else if c.typ != 0 {
		cost += 2
	}
	if c.repeat < 0 && c.typ != 0 {
		cost++
	}
	bitsPerByte := matchLiteralValueBits
	if c.typ == 0 && c.distance <= matchLocalCopyWindow {
		bitsPerByte = matchLocalCopyValueBits
	}
	c.score = c.length*bitsPerByte - cost*8
	return c
}

func (m *matchIndex) find(pos int, recent [3]int) matchChoice {
	n := len(m.input)
	if pos+5 > n {
		return matchChoice{}
	}
	var best matchChoice
	best.repeat = -1
	data := m.input
	// Probe the three recent distances directly, even if their positions were
	// evicted from the short fingerprint buckets. Repeated offsets are cheap
	// to encode and should not depend on hash-bucket retention.
	for _, distance := range recent {
		if distance <= 0 || distance > pos || distance > matchMaxDistance {
			continue
		}
		from := pos - distance
		if !bytes.Equal(data[pos:pos+5], data[from:from+5]) {
			continue
		}
		limit := min(n-pos, matchMaximum)
		size := 5
		for size+8 <= limit && binary.LittleEndian.Uint64(data[pos+size:]) == binary.LittleEndian.Uint64(data[from+size:]) {
			size += 8
		}
		for size < limit && data[pos+size] == data[from+size] {
			size++
		}
		c := scoreMatchChoice(matchChoice{length: size, distance: distance}, recent)
		if c.score > best.score {
			best = c
		}
	}
	for typ := byte(0); typ < 3; typ++ {
		if !m.useTransforms && typ != 0 {
			break
		}
		key := matchFingerprint(data, pos, typ)
		bucket := &m.heads[typ][matchHash(key)]
		for k := 0; k < 4; k++ {
			slot := (bucket.next - 1 - uint32(k)) & 3
			value := bucket.positions[slot]
			if value == 0 {
				break // earlier slots cannot be filled if the newest is empty
			}
			from := int(value) - 1
			distance := pos - from
			if distance <= 0 || distance > matchMaxDistance {
				continue
			}
			var mask byte
			if typ == 1 {
				mask = data[pos] ^ data[from]
			}
			if typ == 2 {
				mask = data[pos] - data[from]
			}
			if typ != 0 && mask == 0 {
				continue
			}
			// Hash bucket collisions are expected. Check every anchor byte.
			anchor := 4
			if typ != 0 {
				anchor = 5
			}
			good := true
			for j := 0; j < anchor; j++ {
				predicted := data[from+j]
				if typ == 1 {
					predicted ^= mask
				} else if typ == 2 {
					predicted += mask
				}
				if data[pos+j] != predicted {
					good = false
					break
				}
			}
			if !good {
				continue
			}
			limit := n - pos
			if limit > matchMaximum {
				limit = matchMaximum
			}
			if typ == 0 {
				size := anchor
				// Common case: compare 8 bytes at a time before the tail.
				for size+8 <= limit && binary.LittleEndian.Uint64(data[pos+size:]) == binary.LittleEndian.Uint64(data[from+size:]) {
					size += 8
				}
				for size < limit && data[pos+size] == data[from+size] {
					size++
				}
				c := matchChoice{length: size, distance: distance, typ: typ}
				if size < matchMinimumExact {
					continue
				}
				c = scoreMatchChoice(c, recent)
				if c.score > best.score {
					best = c
				}
			} else {
				limit = min(limit, 4096)
				size := anchor
				for size < limit {
					predicted := data[from+size]
					if typ == 1 {
						predicted ^= mask
					} else {
						predicted += mask
					}
					if data[pos+size] != predicted {
						break
					}
					size++
				}
				if size < matchMinimumTransformed {
					continue
				}
				c := matchChoice{length: size, distance: distance, typ: typ, mask: mask}
				c = scoreMatchChoice(c, recent)
				if c.score > best.score {
					best = c
				}
			}
			// Further searching cannot materially improve a strong match.
			if best.length >= 8192 {
				return best
			}
		}
	}
	return best
}

func updateMatchRecent(recent *[3]int, distance int) {
	old := *recent
	recent[0] = distance
	dst := 1
	for _, v := range old {
		if v != 0 && v != distance && dst < 3 {
			recent[dst] = v
			dst++
		}
	}
	for dst < 3 {
		recent[dst] = 0
		dst++
	}
}

// encodeMatchV2 returns nil if no profitable references are discovered.
// Every candidate is verified; no source-side index or model is serialized.
func tokenizeMatchV2(input []byte) ([]byte, error) { return tokenizeMatchV2Range(input, 0) }

// tokenizeMatchV2Range emits commands only for input[start:], while the
// preceding bytes are referenceable history already available to the decoder.
func tokenizeMatchV2Range(input []byte, start int) ([]byte, error) {
	return tokenizeMatchV2RangeStride(input, start, 1, false)
}

func tokenizeMatchV2RangeStride(input []byte, start, stride int, exactOnly bool) ([]byte, error) {
	if len(input)-start < 64 {
		return nil, nil
	}
	idx := newMatchIndexTransforms(input, !exactOnly && likelyTransformedMatches(input))
	var commands bytes.Buffer
	var recent [3]int
	literalStart := start
	refs := 0
	flush := func(end int) {
		if end > literalStart {
			commands.WriteByte(matchLiteral)
			putVar(&commands, uint64(end-literalStart))
			commands.Write(input[literalStart:end])
		}
	}
	cachedPos := -1
	var cachedChoice matchChoice
	for pos := start; pos < len(input); {
		idx.addBefore(pos)
		if stride > 1 && (pos-start)%stride != 0 {
			pos++
			continue
		}
		var choice matchChoice
		if cachedPos == pos {
			choice = cachedChoice
			cachedPos = -1
		} else {
			choice = idx.find(pos, recent)
		}
		if choice.score <= 0 {
			pos++
			continue
		}
		// Single-byte lazy parsing. Only take a one-byte delay when its
		// score improves enough to compensate for the new literal byte.
		if stride == 1 && pos+1 < len(input) {
			idx.addBefore(pos + 1)
			next := idx.find(pos+1, recent)
			if next.score > choice.score+2*matchLiteralValueBits {
				cachedPos = pos + 1
				cachedChoice = next
				pos++
				continue
			}
		}
		// Extend a verified reference backwards through pending literals.
		// This is safe for overlapping copies because source bytes precede
		// the output position, and the underlying input is immutable.
		for pos > literalStart && pos > choice.distance {
			from := pos - choice.distance - 1
			actual := input[pos-1]
			predicted := input[from]
			if choice.typ == 1 {
				predicted ^= choice.mask
			} else if choice.typ == 2 {
				predicted += choice.mask
			}
			if actual != predicted || choice.length >= matchMaximum {
				break
			}
			pos--
			choice.length++
		}
		flush(pos)
		opcode := matchCopy
		if choice.typ == 1 {
			opcode = matchXor
		} else if choice.typ == 2 {
			opcode = matchAdd
		}
		if choice.repeat >= 0 {
			if choice.typ == 0 {
				opcode = matchRep0 + byte(choice.repeat)
			}
			if choice.typ == 1 {
				opcode = matchRepXor
			}
			if choice.typ == 2 {
				opcode = matchRepAdd
			}
		}
		commands.WriteByte(opcode)
		if choice.repeat < 0 {
			putVar(&commands, uint64(choice.distance))
		}
		if choice.typ != 0 && choice.repeat >= 0 {
			commands.WriteByte(byte(choice.repeat))
		}
		putVar(&commands, uint64(choice.length))
		if choice.typ != 0 {
			commands.WriteByte(choice.mask)
		}
		updateMatchRecent(&recent, choice.distance)
		pos += choice.length
		literalStart = pos
		refs++
	}
	flush(len(input))
	if refs == 0 {
		return nil, nil
	}
	return commands.Bytes(), nil
}

// Legacy match archive. Stage 2 still evaluates it as a fallback,
// and archives previously written with this mode remain decodable.
func encodeMatchV2(input []byte, digest [32]byte) ([]byte, error) {
	tokens, err := tokenizeMatchV2(input)
	if err != nil || len(tokens) == 0 {
		return nil, err
	}
	return encodeMatchV2Tokens(input, digest, tokens)
}

func encodeMatchV2Tokens(input []byte, digest [32]byte, tokens []byte) ([]byte, error) {
	payload, err := fastNativeBytes(tokens)
	if err != nil {
		return nil, err
	}
	codec := byte(1)
	if len(tokens) <= len(payload) {
		payload = tokens
		codec = 0
	}
	var archive bytes.Buffer
	archive.WriteString(magic)
	archive.WriteByte(modeMatchV2)
	putVar(&archive, uint64(len(input)))
	archive.Write(digest[:])
	archive.WriteByte(codec)
	putVar(&archive, uint64(len(payload)))
	archive.Write(payload)
	return archive.Bytes(), nil
}

func decodeMatchV2(r *bytes.Reader, n uint64, limit int64) ([]byte, error) {
	codec, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	size, err := readVar(r)
	if err != nil {
		return nil, err
	}
	if size != uint64(r.Len()) || size > uint64(limit)*8+1024*1024 {
		return nil, errors.New("bad v2 payload length")
	}
	var tokens []byte
	if codec == 0 {
		tokens = make([]byte, int(size))
		if _, err = io.ReadFull(r, tokens); err != nil {
			return nil, err
		}
	} else if codec == 1 {
		dec := newNativeReader(r)
		bound := int64(n)*4 + 1024*1024
		tokens, err = io.ReadAll(io.LimitReader(dec, bound+1))
		closeErr := dec.Close()
		if err != nil {
			return nil, err
		}
		if closeErr != nil {
			return nil, closeErr
		}
		if int64(len(tokens)) > bound {
			return nil, errors.New("v2 token stream exceeds limit")
		}
	} else {
		return nil, errors.New("unknown v2 token codec")
	}
	stream := bytes.NewReader(tokens)
	out := make([]byte, int(n))
	at := 0
	var recent [3]int
	for at < len(out) {
		op, e := stream.ReadByte()
		if e != nil {
			return nil, e
		}
		if op == matchLiteral {
			length, e := readVar(stream)
			if e != nil {
				return nil, e
			}
			if length == 0 || length > uint64(len(out)-at) {
				return nil, errors.New("bad v2 literal length")
			}
			if _, e = io.ReadFull(stream, out[at:at+int(length)]); e != nil {
				return nil, e
			}
			at += int(length)
			continue
		}
		if op > matchRepAdd {
			return nil, errors.New("unknown v2 match opcode")
		}
		distance := 0
		typ := byte(0)
		if op >= matchRep0 && op <= matchRep2 {
			distance = recent[op-matchRep0]
		} else if op == matchRepXor || op == matchRepAdd {
			idx, e := stream.ReadByte()
			if e != nil {
				return nil, e
			}
			if idx >= 3 {
				return nil, errors.New("invalid v2 repeat index")
			}
			distance = recent[idx]
		} else {
			d, e := readVar(stream)
			if e != nil {
				return nil, e
			}
			if d > uint64(at) {
				return nil, errors.New("invalid v2 match distance")
			}
			distance = int(d)
		}
		if op == matchXor || op == matchRepXor {
			typ = 1
		}
		if op == matchAdd || op == matchRepAdd {
			typ = 2
		}
		length, e := readVar(stream)
		if e != nil {
			return nil, e
		}
		if distance <= 0 || distance > at || distance > matchMaxDistance || length == 0 || length > uint64(len(out)-at) || length > matchMaximum {
			return nil, errors.New("invalid v2 backward reference")
		}
		var mask byte
		if typ != 0 {
			mask, e = stream.ReadByte()
			if e != nil {
				return nil, e
			}
			if mask == 0 {
				return nil, errors.New("redundant v2 transformation")
			}
		}
		expandMatch(out, at, distance, int(length), typ, mask)
		updateMatchRecent(&recent, distance)
		at += int(length)
	}
	if stream.Len() != 0 {
		return nil, errors.New("trailing v2 tokens")
	}
	return out, nil
}

func matchSampleUseful(input []byte) (bool, error) {
	if len(input) <= 384*1024 {
		return true, nil
	}
	// A short, fixed work budget to avoid a full matching pass on data where
	// the current native byte packer is already stronger. The sample is not archived.
	size := min(len(input), 24*1024)
	sample := input[:size]
	candidate, _, err := encodeMatchStage2(sample, sha256.Sum256(sample))
	if err != nil {
		return false, err
	}
	if len(candidate) == 0 {
		return false, nil
	}
	baseline, err := encodeFastNative(sample, sha256.Sum256(sample))
	if err != nil {
		return false, err
	}
	// Screen using the best of ALL representations the whole-file encoder
	// would actually choose. Previously this gate considered legacy tokens
	// alone, silently rejecting potentially winning native encodings.
	// This remains a bounded sample heuristic, not a full-file proof.
	return len(candidate) > 0 && len(candidate)+len(baseline)/50 < len(baseline), nil
}
