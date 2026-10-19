// JLZ is an experimental native-only, lossless mathematical-dependency compressor.
// It discovers distant chunk relationships, XOR/subtract lag dependencies,
// and overlapping bit-pattern candidates with hierarchical grammar rules.
//
// Usage:
//
//	go run ./cmd/jlz compress input.bin [output.jlz]
//	go run ./cmd/jlz decompress input.bin.jlz [restored.bin]
//
// Go 1.23+, standard library only.
package engine

import (
 "github.com/jacoblockett/jlz/archive"
	"bufio"
	"bytes"
	"crypto/sha256"
		"errors"
	"fmt"
	"io"
	"math"
	"time"
)

const (
	magic                = "JLZ1"
	modeRaw         byte = 0
	modeGrammar     byte = 1
	modeRelation    byte = 2
	modeGraph       byte = 3
	modeNativeBytes byte = 4
	modeHybrid      byte = 5
	opConcat        byte = 0
	opNot           byte = 1
	opXor           byte = 2
	defaultLimit         = 512 * 1024 * 1024 // Maximum input and decoded output size; not an encoding level
)

type rule struct {
	op   byte
	a, b uint32
}
type pair struct{ a, b uint32 }
type stats struct{ seen, usable, last int }
type model struct {
	length uint64
	value  uint64
	packed bool
}
type summary struct {
	original, archive  int
	rules, symbols     int
	grammarBytes       int
	notRules, xorRules int
	raw                bool
	method             string
	detail             string
	sha                string
	elapsed            time.Duration
}

func show(s summary, path string) {
	kind := s.method
	if kind == "" {
		kind = "raw fallback"
	}
	pct := 0.0
	if s.original > 0 {
		pct = 100 * (1 - float64(s.archive)/float64(s.original))
	}
	fmt.Printf("%d -> %d bytes (%.2f%% savings), %s, %s\n", s.original, s.archive, pct, kind, s.elapsed.Round(time.Millisecond))
	if s.grammarBytes > 0 {
		fmt.Printf("Grammar candidate: %d bytes; %d rules (%d NOT, %d XOR); %d remaining tokens\n", s.grammarBytes, s.rules, s.notRules, s.xorRules, s.symbols)
	}
	if s.detail != "" {
		fmt.Println("Detail:", s.detail)
	}
	fmt.Printf("SHA-256: %s\nSaved: %s\n", s.sha, path)
}

func putVar(w io.Writer, n uint64) error { return archive.PutVar(w, n) }
func readVar(r io.ByteReader) (uint64,error) { return archive.ReadVar(r) }
func varSize(n uint64) int { return archive.VarSize(n) }

// compress is JLZ's sole compression entry point. All mode selection is automatic.
func compress(input []byte) ([]byte, summary, error) {
	return compressFast(input)
}

// learnGrammar constructs a Re-Pair-style grammar. We count overlapping pair
// candidates, but only replace nonoverlapping instances in the final stream.
func learnGrammar(tokens []uint32, maxRules, minUses int) ([]uint32, []rule) {
	rules := make([]rule, 0, maxRules)
	for len(rules) < maxRules && len(tokens) >= minUses*2 {
		counts := make(map[pair]stats, 32)
		for i := 0; i+1 < len(tokens); i++ {
			key := pair{tokens[i], tokens[i+1]}
			st, ok := counts[key]
			st.seen++
			if !ok || i > st.last+1 {
				st.usable++
				st.last = i
			}
			counts[key] = st
		}
		var best pair
		bestCount := 0
		for p, st := range counts {
			// Tie breaks are deterministic regardless of map iteration order.
			if st.usable > bestCount || (st.usable == bestCount && (p.a < best.a || p.a == best.a && p.b < best.b)) {
				best, bestCount = p, st.usable
			}
		}
		if bestCount < minUses {
			break
		}
		id := uint32(2 + len(rules))
		out := make([]uint32, 0, len(tokens)-bestCount)
		for i := 0; i < len(tokens); {
			if i+1 < len(tokens) && tokens[i] == best.a && tokens[i+1] == best.b {
				out = append(out, id)
				i += 2
			} else {
				out = append(out, tokens[i])
				i++
			}
		}
		rules = append(rules, rule{op: opConcat, a: best.a, b: best.b})
		tokens = out
	}
	return tokens, rules
}

func maskFor(n uint64) uint64 {
	if n == 64 {
		return ^uint64(0)
	}
	return (uint64(1) << n) - 1
}
func ruleModel(r rule, prev []model) model {
	a := prev[r.a]
	switch r.op {
	case opConcat:
		b := prev[r.b]
		total := a.length + b.length
		if total <= 64 && a.packed && b.packed {
			return model{total, (a.value << b.length) | b.value, true}
		}
		return model{length: total}
	case opNot:
		return model{a.length, (^a.value) & maskFor(a.length), true}
	case opXor:
		return model{a.length, a.value ^ prev[r.b].value, true}
	}
	panic("bad rule")
}

// For patterns up to 64 bits, choose mathematical encodings that need no more
// operand bytes than their concatenation definition. All references point back.
func rewrite(rules []rule) {
	m := make([]model, 2, 2+len(rules))
	m[0] = model{1, 0, true}
	m[1] = model{1, 1, true}
	for i, r := range rules {
		v := ruleModel(r, m)
		candidate := r
		cost := varSize(uint64(r.a)) + varSize(uint64(r.b))
		if v.packed && v.length >= 2 {
			for j := 2; j < len(m); j++ {
				if !m[j].packed || m[j].length != v.length {
					continue
				}
				if v.value == (m[j].value ^ maskFor(v.length)) {
					c := varSize(uint64(j))
					if c < cost {
						candidate = rule{op: opNot, a: uint32(j)}
						cost = c
					}
				}
			}
			if v.length >= 4 {
				for j := 2; j < len(m); j++ {
					if !m[j].packed || m[j].length != v.length {
						continue
					}
					for k := j + 1; k < len(m); k++ {
						if !m[k].packed || m[k].length != v.length {
							continue
						}
						if m[j].value^m[k].value != v.value {
							continue
						}
						c := varSize(uint64(j)) + varSize(uint64(k))
						// Allow cost-neutral XOR to expose dependencies.
						if c < cost || (c == cost && candidate.op == opConcat) {
							candidate = rule{op: opXor, a: uint32(j), b: uint32(k)}
							cost = c
						}
					}
				}
			}
		}
		rules[i] = candidate
		// The revised definition reconstructs exactly the same bits.
		m = append(m, v)
	}
}

func decompress(archive []byte, limit int64) ([]byte, error) {
	if len(archive) < 5 || (string(archive[:4]) != magic) {
		return nil, errors.New("not a JLZ archive (JLZ1 format required)")
	}
	mode := archive[4]
	reader := bytes.NewReader(archive[5:])
	n, err := readVar(reader)
	if err != nil {
		return nil, err
	}
	if n > uint64(limit) || n > uint64(math.MaxInt/8) {
		return nil, errors.New("declared output exceeds JLZ safety limit")
	}
	var expected [32]byte
	if _, err := io.ReadFull(reader, expected[:]); err != nil {
		return nil, err
	}
	var output []byte
	if mode == modeRaw {
		if uint64(reader.Len()) != n {
			return nil, errors.New("raw payload has wrong size")
		}
		output = make([]byte, int(n))
		_, err = io.ReadFull(reader, output)
		if err != nil {
			return nil, err
		}
	} else if mode == modeGrammar {
		nr, err := readVar(reader)
		if err != nil {
			return nil, err
		}
		if nr > 4096 {
			return nil, errors.New("excessive rule count")
		}
		rules := make([]rule, 0, int(nr))
		models := []model{{1, 0, true}, {1, 1, true}}
		for i := uint64(0); i < nr; i++ {
			op, err := reader.ReadByte()
			if err != nil {
				return nil, err
			}
			a, err := readVar(reader)
			if err != nil {
				return nil, err
			}
			var b uint64
			if op != opNot {
				b, err = readVar(reader)
				if err != nil {
					return nil, err
				}
			}
			if a >= uint64(len(models)) || (op != opNot && b >= uint64(len(models))) {
				return nil, errors.New("rule has forward/invalid reference")
			}
			if op > opXor {
				return nil, errors.New("unknown rule operation")
			}
			r := rule{op, uint32(a), uint32(b)}
			ma, mb := models[a], models[b]
			if op == opXor && (ma.length != mb.length || !ma.packed || !mb.packed) {
				return nil, errors.New("invalid XOR rule")
			}
			if op == opNot && !ma.packed {
				return nil, errors.New("invalid NOT rule")
			}
			if op == opConcat && (ma.length > n*8 || mb.length > n*8 || ma.length+mb.length > n*8) {
				return nil, errors.New("rule expands beyond output")
			}
			m := ruleModel(r, models)
			if m.length > n*8 {
				return nil, errors.New("rule exceeds declared output")
			}
			rules = append(rules, r)
			models = append(models, m)
		}
		nsym, err := readVar(reader)
		if err != nil {
			return nil, err
		}
		if nsym > n*8 {
			return nil, errors.New("too many symbols")
		}
		output = make([]byte, int(n))
		fl := newNativeReader(reader)
		defer fl.Close()
		br := bufio.NewReader(fl)
		var written uint64
		stack := make([]uint32, 0, 64)
		appendBits := func(v uint64, bits uint64) error {
			if written+bits > n*8 {
				return errors.New("expanded data exceeds output size")
			}
			for j := bits; j > 0; j-- {
				if v&(uint64(1)<<(j-1)) != 0 {
					output[written/8] |= 1 << (7 - (written % 8))
				}
				written++
			}
			return nil
		}
		for i := uint64(0); i < nsym; i++ {
			symbol, err := readVar(br)
			if err != nil {
				return nil, fmt.Errorf("malformed token %d: %w", i, err)
			}
			if symbol >= uint64(len(models)) {
				return nil, errors.New("unknown token")
			}
			stack = append(stack, uint32(symbol))
			for len(stack) > 0 {
				id := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				m := models[id]
				if m.packed {
					if err := appendBits(m.value, m.length); err != nil {
						return nil, err
					}
				} else {
					r := rules[id-2] // nonpacked nodes must be concatenations
					if r.op != opConcat {
						return nil, errors.New("invalid rule structure")
					}
					stack = append(stack, r.b, r.a)
				}
			}
		}
		if written != n*8 {
			return nil, fmt.Errorf("reconstructed %d bits, expected %d", written, n*8)
		}
		// Reject extra decoded data (including trailing token bytes).
		_, err = br.ReadByte()
		if err != io.EOF {
			if err == nil {
				return nil, errors.New("extra data after final token")
			}
			return nil, err
		}
	} else if mode == modeMosaic {
		output, err = decodeMosaic(reader, n)
		if err != nil {
			return nil, err
		}
	} else if mode == modePredictive {
		output, err = decodePredictive(reader, n)
		if err != nil {
			return nil, err
		}
	} else if mode == modeParallelNative {
		output, err = decodeParallelNative(reader, n)
		if err != nil {
			return nil, err
		}
	} else if mode == modeMixed {
		output, err = decodeMixedBlocks(reader, n)
	} else if mode == modeMixedCross {
		output, err = decodeMixedCross(reader, n)
		if err != nil {
			return nil, err
		}
	} else if mode == modeMatchPacked || mode == modeMatchPackedPred {
		output, err = decodeMatchPacked(reader, n, mode == modeMatchPackedPred)
		if err != nil {
			return nil, err
		}
	} else if mode == modeMatchNativePred {
		output, err = decodeMatchNativeVariant(reader, n, true)
		if err != nil {
			return nil, err
		}
	} else if mode == modeMatchNative {
		output, err = decodeMatchNative(reader, n)
		if err != nil {
			return nil, err
		}
	} else if mode == modeMatchV2 {
		output, err = decodeMatchV2(reader, n, limit)
		if err != nil {
			return nil, err
		}
	} else if mode == modeAdaptive {
		output, err = decodeAdaptive(reader, n, limit)
		if err != nil {
			return nil, err
		}
	} else if mode == modeHybrid {
		output, err = decodeHybrid(reader, n, limit)
		if err != nil {
			return nil, err
		}
	} else if mode == modeNativeBytes {
		output, err = decodeNativeBytes(reader, n)
		if err != nil {
			return nil, err
		}
	} else if mode == modeRelation {
		output, err = decodeLag(reader, n)
		if err != nil {
			return nil, err
		}
	} else if mode == modeGraph {
		output, err = decodeGraph(reader, n)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("unknown archive mode")
	}
	if sha256.Sum256(output) != expected {
		return nil, errors.New("SHA-256 mismatch: corrupted archive")
	}
	return output, nil
}
