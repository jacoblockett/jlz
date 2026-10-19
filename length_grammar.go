package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// Length-grammar codec 5 is an optional packed match-length substream only.
// G(256+i) expands to two earlier symbols, making all rules acyclic.
const lengthGrammarCodec byte = 5
const lengthGrammarRules = 16

type lengthPair struct{ a, b uint16 }

func lengthVarSize(x uint64) int {
	var buf [10]byte
	return binary.PutUvarint(buf[:], x)
}
func writeLengthVar(b *bytes.Buffer, x uint64) {
	var buf [10]byte
	n := binary.PutUvarint(buf[:], x)
	b.Write(buf[:n])
}

// This bounded pair sample avoids 16 grammar passes for unpromising streams.
// Fast screening may miss candidates; complete serialized bytes always select.
func lengthGrammarGate(data []byte) bool {
	if len(data) < 2048 || len(data) > 256*1024 {
		return false
	}
	n := min(len(data), 8192)
	var counts [65536]uint16
	best := 0
	for i := 0; i+1 < n; i += 2 {
		k := int(data[i])<<8 | int(data[i+1])
		counts[k]++
		if int(counts[k]) > best {
			best = int(counts[k])
		}
	}
	return best >= 384
}

func packLengthStream(data []byte) []byte {
	incumbent := packStream(data)
	if !lengthGrammarGate(data) {
		return incumbent
	}
	proposal := packLengthGrammar(data, len(incumbent))
	if len(proposal) != 0 && len(proposal) < len(incumbent) {
		return proposal
	}
	return incumbent
}

// An exact 16-rule hierarchical Re-Pair variant, checking packed size every
// fourth rule (and at the final short checkpoint). Definitions are charged.
func packLengthGrammar(data []byte, incumbent int) []byte {
	if len(data) < 2 || len(data) > 256*1024 {
		return nil
	}
	tokens := make([]uint16, len(data))
	for i, c := range data {
		tokens[i] = uint16(c)
	}
	rules := make([]lengthPair, 0, lengthGrammarRules)
	bestSize := incumbent
	var best []byte
	const alphabet = 256 + lengthGrammarRules
	const bins = alphabet * alphabet
	counts := make([]uint32, bins)
	last := make([]int32, bins)
	seen := make([]uint8, bins)
	for len(rules) < lengthGrammarRules {
		stamp := uint8(len(rules) + 1)
		winner, bestUses := 0, uint32(0)
		for i := 0; i+1 < len(tokens); i++ {
			k := int(tokens[i])*alphabet + int(tokens[i+1])
			if seen[k] != stamp {
				seen[k] = stamp
				counts[k] = 0
				last[k] = -2
			}
			if int32(i) < last[k]+2 {
				continue
			}
			last[k] = int32(i)
			counts[k]++
			if counts[k] > bestUses || (counts[k] == bestUses && k < winner) {
				bestUses, winner = counts[k], k
			}
		}
		if bestUses < 24 {
			break
		}
		a, b := uint16(winner/alphabet), uint16(winner%alphabet)
		replaced := make([]uint16, 0, len(tokens)-int(bestUses))
		id := uint16(256 + len(rules))
		for i := 0; i < len(tokens); {
			if i+1 < len(tokens) && tokens[i] == a && tokens[i+1] == b {
				replaced = append(replaced, id)
				i += 2
			} else {
				replaced = append(replaced, tokens[i])
				i++
			}
		}
		rules = append(rules, lengthPair{a, b})
		tokens = replaced
		if len(rules)%4 != 0 && bestUses >= 48 {
			continue
		}
		lo, hi := make([]byte, len(tokens)), make([]byte, len(tokens))
		for i, v := range tokens {
			lo[i], hi[i] = byte(v), byte(v>>8)
		}
		var body bytes.Buffer
		body.WriteByte(byte(len(rules)))
		writeLengthVar(&body, uint64(len(tokens)))
		for _, rule := range rules {
			var pair [4]byte
			binary.LittleEndian.PutUint16(pair[:2], rule.a)
			binary.LittleEndian.PutUint16(pair[2:], rule.b)
			body.Write(pair[:])
		}
		body.Write(packStream(lo))
		body.Write(packStream(hi))
		need := body.Len() + lengthVarSize(uint64(len(data))) + lengthVarSize(uint64(body.Len())) + 1
		if need >= bestSize {
			continue
		}
		var out bytes.Buffer
		writeLengthVar(&out, uint64(len(data)))
		writeLengthVar(&out, uint64(body.Len()))
		out.WriteByte(lengthGrammarCodec)
		out.Write(body.Bytes())
		bestSize, best = out.Len(), append([]byte(nil), out.Bytes()...)
	}
	return best
}

// Nested token-coordinate streams must remain old native codecs 0/1/2.
// In particular they cannot recursively contain this grammar codec.
func unpackGrammarLane(rd *bytes.Reader, count int) ([]byte, error) {
	at, _ := rd.Seek(0, io.SeekCurrent)
	_, err := binary.ReadUvarint(rd)
	if err != nil {
		return nil, err
	}
	_, err = binary.ReadUvarint(rd)
	if err != nil {
		return nil, err
	}
	typ, err := rd.ReadByte()
	if err != nil {
		return nil, err
	}
	if typ > 2 {
		return nil, errors.New("nested grammar stream codec")
	}
	if _, err = rd.Seek(at, io.SeekStart); err != nil {
		return nil, err
	}
	return unpackStream(rd, count)
}

func unpackLengthGrammar(payload []byte, n int) ([]byte, error) {
	if len(payload) < 2 || n < 0 {
		return nil, errors.New("invalid grammar payload")
	}
	rd := bytes.NewReader(payload)
	rc, err := rd.ReadByte()
	if err != nil {
		return nil, err
	}
	if rc == 0 || rc > lengthGrammarRules {
		return nil, errors.New("invalid grammar rules")
	}
	tokCount, err := binary.ReadUvarint(rd)
	if err != nil || tokCount == 0 || tokCount > uint64(n) {
		return nil, errors.New("invalid grammar token count")
	}
	if rd.Len() < int(rc)*4 {
		return nil, io.ErrUnexpectedEOF
	}
	rules := make([]lengthPair, int(rc))
	var expansion [256 + lengthGrammarRules]uint32
	for i := 0; i < 256; i++ {
		expansion[i] = 1
	}
	for i := range rules {
		var pair [4]byte
		if _, err = io.ReadFull(rd, pair[:]); err != nil {
			return nil, err
		}
		a, b := binary.LittleEndian.Uint16(pair[:2]), binary.LittleEndian.Uint16(pair[2:])
		if int(a) >= 256+i || int(b) >= 256+i {
			return nil, errors.New("forward grammar reference")
		}
		rules[i] = lengthPair{a, b}
		size := uint64(expansion[a]) + uint64(expansion[b])
		if size > uint64(n) {
			return nil, errors.New("grammar expansion exceeds bound")
		}
		expansion[256+i] = uint32(size)
	}
	lo, err := unpackGrammarLane(rd, int(tokCount))
	if err != nil {
		return nil, err
	}
	hi, err := unpackGrammarLane(rd, int(tokCount))
	if err != nil {
		return nil, err
	}
	if rd.Len() != 0 || len(lo) != int(tokCount) || len(hi) != int(tokCount) {
		return nil, errors.New("grammar coordinate framing")
	}
	out := make([]byte, 0, n)
	var stack [64]uint16
	for i := range lo {
		v := uint16(lo[i]) | uint16(hi[i])<<8
		if int(v) >= 256+len(rules) {
			return nil, errors.New("unknown grammar symbol")
		}
		if len(out)+int(expansion[v]) > n {
			return nil, errors.New("grammar output overrun")
		}
		stack[0] = v
		top := 1
		for top > 0 {
			top--
			id := stack[top]
			if id < 256 {
				out = append(out, byte(id))
				continue
			}
			pair := rules[int(id)-256]
			if top+2 > len(stack) {
				return nil, errors.New("grammar nesting too deep")
			}
			stack[top], stack[top+1] = pair.b, pair.a
			top += 2
		}
	}
	if len(out) != n {
		return nil, errors.New("grammar reconstructed length")
	}
	return out, nil
}
