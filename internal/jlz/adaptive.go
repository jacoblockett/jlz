package jlz

// MDC4 research mode: approximate matching at arbitrary byte offsets.
// It finds earlier substrings with common 8-byte anchors, tests different
// candidate lengths, stores only sparse XOR exceptions, then uses a
// shortest-path dynamic program to choose a sequence of literal/match tokens.
// All decoding depends solely on previously reconstructed source bytes.

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const modeAdaptive byte = 6
const (
	adaptiveLiteral  byte = 0
	adaptiveCopy     byte = 1
	adaptivePatch    byte = 2
	adaptiveXor      byte = 3
	adaptiveAdd      byte = 4
	adaptiveXorPatch byte = 5
	adaptiveAddPatch byte = 6
)

type edit struct {
	position int
	xor      byte
}
type adaptiveChoice struct {
	from, length, cost int
	transform, mask    byte
	edits              []edit
}
type adaptiveStats struct {
	copies, patches, literals, edits, xorRefs, addRefs int
	commands                                           int
}

// Invariant anchor keys detect bytewise XOR and modular-add relations even
// when the transformed blocks share no literal 8-byte substring.
func adaptiveKey(data []byte, p int, transform byte) uint64 {
	if transform == 0 {
		return binary.LittleEndian.Uint64(data[p:])
	}
	var key uint64
	for k := 1; k < 8; k++ {
		v := data[p+k]
		if transform == 1 {
			v ^= data[p]
		} else {
			v -= data[p]
		}
		key |= uint64(v) << uint(8*(k-1))
	}
	return key
}

// Proposed segments are not tied to global chunks. References can begin at
// any byte in the earlier input, including positions that differ in alignment.
func adaptiveCandidates(data []byte) [][]adaptiveChoice {
	n := len(data)
	choices := make([][]adaptiveChoice, n)
	if n < 128 {
		return choices
	}
	indexes := [3]map[uint64][]int{}
	for j := range indexes {
		indexes[j] = make(map[uint64][]int, n/8)
	}
	positions := []int{0, 8, 16, 32, 64}
	lengths := []int{64, 128, 256, 512, 1024, 2048, 4096}
	for pos := 0; pos < n; pos++ {
		// Keep a lookbehind, so candidate references finish before the
		// corresponding destination starts. This also bounds self-overlap.
		prior := pos - 64
		if prior >= 0 && prior+8 <= n {
			for transform := byte(0); transform < 3; transform++ {
				h := adaptiveKey(data, prior, transform)
				existing := indexes[transform][h]
				if len(existing) < 8 {
					indexes[transform][h] = append(existing, prior)
				} else {
					copy(existing[1:], existing[2:])
					existing[7] = prior
				}
			}
		}
		if pos%4 != 0 || pos+64 > n {
			continue
		}
		// Find exact matching anchors a little *inside* a prospective
		// segment, then test whether their surrounding bytes mostly match.
		seen := make(map[[3]int]bool, 60)
		best := make(map[int]adaptiveChoice, 6)
		for _, offset := range positions {
			if pos+offset+8 > n {
				continue
			}
			for transform := byte(0); transform < 3; transform++ {
				key := adaptiveKey(data, pos+offset, transform)
				for _, where := range indexes[transform][key] {
					from := where - offset
					if from < 0 || from+64 > pos {
						continue
					}
					var mask byte
					if transform == 1 {
						mask = data[pos+offset] ^ data[where]
					}
					if transform == 2 {
						mask = data[pos+offset] - data[where]
					}
					if transform != 0 && mask == 0 {
						continue
					}
					signature := [3]int{from, int(transform), int(mask)}
					if seen[signature] {
						continue
					}
					seen[signature] = true
					limit := 4096
					if n-pos < limit {
						limit = n - pos
					}
					if pos-from < limit {
						limit = pos - from
					}
					if limit < 64 {
						continue
					}
					predict := func(j int) byte {
						v := data[from+j]
						if transform == 1 {
							return v ^ mask
						}
						if transform == 2 {
							return v + mask
						}
						return v
					}
					wrong := 0
					for j := 0; j < 64; j++ {
						if data[pos+j] != predict(j) {
							wrong++
						}
					}
					if wrong > 16 {
						continue
					}
					edits := make([]edit, 0, 64)
					editBytes, prev, nextCheck := 0, 0, 0
					for j := 0; j < limit; j++ {
						v := predict(j)
						if data[pos+j] != v {
							gap := j - prev + 1
							edits = append(edits, edit{j, data[pos+j] ^ v})
							editBytes += varSize(uint64(gap)) + 1
							prev = j + 1
						}
						if nextCheck >= len(lengths) || j+1 != lengths[nextCheck] {
							continue
						}
						width := lengths[nextCheck]
						nextCheck++
						if len(edits) > width/5 {
							break
						}
						cost := 1 + varSize(uint64(from)) + varSize(uint64(width))
						if transform != 0 {
							cost++
						}
						if len(edits) > 0 {
							cost += varSize(uint64(len(edits))) + editBytes
						}
						if cost+4 >= width {
							continue
						}
						old, exists := best[width]
						if !exists || cost < old.cost {
							saved := append([]edit(nil), edits...)
							best[width] = adaptiveChoice{from: from, length: width, cost: cost, transform: transform, mask: mask, edits: saved}
						}
					}
				}
			}
		}
		for _, l := range lengths {
			if v, ok := best[l]; ok {
				choices[pos] = append(choices[pos], v)
			}
		}
	}
	return choices
}

func adaptiveCommands(data []byte) ([]byte, adaptiveStats) {
	return adaptiveCommandsWithCost(data, adaptiveCandidates(data), 1000, 1000)
}
func adaptiveCommandsWithCost(data []byte, candidates [][]adaptiveChoice, literalWeight, commandWeight int) ([]byte, adaptiveStats) {
	n := len(data)
	// Solve an approximate minimum-byte description. A literal costs 1
	// byte; selecting any dependency pays its exact token cost. The final
	// encoder coalesces adjacent literals to reduce header overhead.
	dp := make([]int, n+1)
	choice := make([]int, n)
	for p := n - 1; p >= 0; p-- {
		dp[p] = literalWeight + dp[p+1]
		choice[p] = -1
		for i, c := range candidates[p] {
			if commandWeight*c.cost+dp[p+c.length] < dp[p] {
				dp[p] = commandWeight*c.cost + dp[p+c.length]
				choice[p] = i
			}
		}
	}
	var encoded bytes.Buffer
	var st adaptiveStats
	for p := 0; p < n; {
		if choice[p] < 0 {
			start := p
			p++
			for p < n && choice[p] < 0 {
				p++
			}
			encoded.WriteByte(adaptiveLiteral)
			putVar(&encoded, uint64(p-start))
			encoded.Write(data[start:p])
			st.literals++
			continue
		}
		c := candidates[p][choice[p]]
		switch c.transform {
		case 0:
			if len(c.edits) == 0 {
				encoded.WriteByte(adaptiveCopy)
				st.copies++
			} else {
				encoded.WriteByte(adaptivePatch)
				st.patches++
			}
		case 1:
			if len(c.edits) == 0 {
				encoded.WriteByte(adaptiveXor)
			} else {
				encoded.WriteByte(adaptiveXorPatch)
				st.patches++
			}
			st.xorRefs++
		case 2:
			if len(c.edits) == 0 {
				encoded.WriteByte(adaptiveAdd)
			} else {
				encoded.WriteByte(adaptiveAddPatch)
				st.patches++
			}
			st.addRefs++
		}
		putVar(&encoded, uint64(c.from))
		putVar(&encoded, uint64(c.length))
		if c.transform != 0 {
			encoded.WriteByte(c.mask)
		}
		if len(c.edits) > 0 {
			putVar(&encoded, uint64(len(c.edits)))
			prev := 0
			for _, e := range c.edits {
				putVar(&encoded, uint64(e.position-prev+1))
				encoded.WriteByte(e.xor)
				prev = e.position + 1
			}
			st.edits += len(c.edits)
		}
		p += c.length
	}
	st.commands = encoded.Len()
	return encoded.Bytes(), st
}

func encodeAdaptive(input []byte, digest [32]byte) ([]byte, string, error) {
	candidates := adaptiveCandidates(input)
	profiles := [][2]int{{1000, 1000}, {650, 1000}, {1000, 650}, {600, 700}}
	var best []byte
	bestInfo := ""
	// Model several different literal/reference costs. Their real native stream
	// sizes can change the ranking, so choose after compression.
	for _, profile := range profiles {
		tokens, st := adaptiveCommandsWithCost(input, candidates, profile[0], profile[1])
		for _, codec := range []byte{0, 1, 2} {
			payload := tokens
			if codec == 1 {
				var err error
				payload, err = nativeBytes(tokens)
				if err != nil {
					return nil, "", err
				}
			} else if codec == 2 {
				payload = arithmeticEncode(tokens)
			}
			var archive bytes.Buffer
			archive.WriteString(magic)
			archive.WriteByte(modeAdaptive)
			putVar(&archive, uint64(len(input)))
			archive.Write(digest[:])
			archive.WriteByte(codec)
			if codec != 0 {
				putVar(&archive, uint64(len(payload)))
				if codec == 2 {
					putVar(&archive, uint64(len(tokens)))
				}
			}
			archive.Write(payload)
			if best == nil || archive.Len() < len(best) {
				best = archive.Bytes()
				codecName := []string{"raw", "native", "arithmetic"}[codec]
				bestInfo = fmt.Sprintf("adaptive: %d copies, %d XOR, %d ADD, %d patches, %d edited bytes, %d literals; %s tokens; weights %d/%d", st.copies, st.xorRefs, st.addRefs, st.patches, st.edits, st.literals, codecName, profile[0], profile[1])
			}
		}
	}
	return best, bestInfo, nil
}

func decodeAdaptive(r *bytes.Reader, n uint64, limit int64) ([]byte, error) {
	flag, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	var stream []byte
	switch flag {
	case 0:
		stream = make([]byte, r.Len())
		if _, err = io.ReadFull(r, stream); err != nil {
			return nil, err
		}
	case 1:
		sz, e := readVar(r)
		if e != nil {
			return nil, e
		}
		if sz > uint64(r.Len()) || sz > uint64(limit)*8+1024*1024 {
			return nil, errors.New("bad adaptive payload length")
		}
		if sz != uint64(r.Len()) {
			return nil, errors.New("trailing bytes in adaptive archive")
		}
		fl := newNativeReader(r)
		defer fl.Close()
		// Bounded by 8 times declared output + 1MiB token overhead.
		bound := int64(n)*8 + 1024*1024
		stream, err = io.ReadAll(io.LimitReader(fl, bound+1))
		if err != nil {
			return nil, err
		}
		if int64(len(stream)) > bound {
			return nil, errors.New("adaptive token stream too large")
		}
	case 2:
		sz, e := readVar(r)
		if e != nil {
			return nil, e
		}
		if sz > uint64(r.Len()) || sz > uint64(limit)*8+1024*1024 {
			return nil, errors.New("bad adaptive arithmetic payload size")
		}
		decodedLength, e := readVar(r)
		if e != nil {
			return nil, e
		}
		if decodedLength > n*8+1024*1024 {
			return nil, errors.New("adaptive token length too large")
		}
		if sz != uint64(r.Len()) {
			return nil, errors.New("trailing adaptive arithmetic bytes")
		}
		compressed := make([]byte, r.Len())
		if _, e = io.ReadFull(r, compressed); e != nil {
			return nil, e
		}
		stream, e = arithmeticDecode(compressed, int(decodedLength))
		if e != nil {
			return nil, e
		}
	default:
		return nil, errors.New("unknown adaptive encoding")
	}
	tokens := bytes.NewReader(stream)
	output := make([]byte, int(n))
	at := 0
	for at < len(output) {
		opcode, e := tokens.ReadByte()
		if e != nil {
			return nil, e
		}
		if opcode > adaptiveAddPatch {
			return nil, errors.New("invalid adaptive opcode")
		}
		if opcode == adaptiveLiteral {
			length, e := readVar(tokens)
			if e != nil {
				return nil, e
			}
			if length == 0 || length > uint64(len(output)-at) {
				return nil, errors.New("bad adaptive literal length")
			}
			if _, e = io.ReadFull(tokens, output[at:at+int(length)]); e != nil {
				return nil, e
			}
			at += int(length)
			continue
		}
		pos, e := readVar(tokens)
		if e != nil {
			return nil, e
		}
		length, e := readVar(tokens)
		if e != nil {
			return nil, e
		}
		if length == 0 || length > uint64(len(output)-at) || pos > uint64(at) || length > uint64(at)-pos {
			return nil, errors.New("invalid adaptive backward reference")
		}
		copy(output[at:at+int(length)], output[int(pos):int(pos)+int(length)])
		if opcode >= adaptiveXor {
			mask, e := tokens.ReadByte()
			if e != nil {
				return nil, e
			}
			if opcode == adaptiveXor || opcode == adaptiveXorPatch {
				for j := 0; j < int(length); j++ {
					output[at+j] ^= mask
				}
			} else {
				for j := 0; j < int(length); j++ {
					output[at+j] += mask
				}
			}
		}
		if opcode == adaptivePatch || opcode == adaptiveXorPatch || opcode == adaptiveAddPatch {
			count, e := readVar(tokens)
			if e != nil {
				return nil, e
			}
			if count == 0 || count > length {
				return nil, errors.New("bad adaptive patch count")
			}
			editAt := -1
			for j := uint64(0); j < count; j++ {
				delta, e := readVar(tokens)
				if e != nil {
					return nil, e
				}
				if delta == 0 || delta > length {
					return nil, errors.New("invalid adaptive patch gap")
				}
				if editAt+int(delta) >= int(length) {
					return nil, errors.New("adaptive patch outside segment")
				}
				editAt += int(delta)
				mask, e := tokens.ReadByte()
				if e != nil {
					return nil, e
				}
				if mask == 0 {
					return nil, errors.New("redundant adaptive patch")
				}
				output[at+editAt] ^= mask
			}
		}
		at += int(length)
	}
	if tokens.Len() != 0 {
		return nil, errors.New("trailing adaptive token data")
	}
	return output, nil
}
