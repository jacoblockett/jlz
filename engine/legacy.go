package engine

// Retained fast v1 matcher for short files where its coarse match parsing
// outperforms the new short-match parser; original archives remain readable.

import (
	"bytes"
)

type matchEntry struct {
	vals  [4]int
	count int
}

// A fixed-depth hash chain, O(n) expected construction and bounded lookup.
// Three invariant fingerprints permit exact, constant-XOR, and constant-ADD
// byte matches without enumerating masks. References never self-overlap.
func fastAdaptive(input []byte, digest [32]byte) ([]byte, error) {
	n := len(input)
	if n < 128 {
		return nil, nil
	}
	idx := [3]map[uint64]matchEntry{}
	for i := range idx {
		idx[i] = make(map[uint64]matchEntry, n/10)
	}
	insert := func(at int) {
		if at+8 > n {
			return
		}
		for t := range idx {
			key := adaptiveKey(input, at, byte(t))
			e := idx[t][key]
			e.vals[e.count%len(e.vals)] = at
			e.count++
			idx[t][key] = e
		}
	}
	var cmds bytes.Buffer
	indexed := 0
	litStart := 0
	flush := func(end int) {
		if end == litStart {
			return
		}
		cmds.WriteByte(adaptiveLiteral)
		putVar(&cmds, uint64(end-litStart))
		cmds.Write(input[litStart:end])
	}
	hasRefs := false
	for pos := 0; pos < n; {
		// An insertion delay guarantees that candidates finish before pos.
		for indexed+8 <= n && indexed <= pos-64 {
			insert(indexed)
			indexed++
		}
		if pos+64 > n {
			break
		}
		bestLen, bestFrom, bestTransform := 0, 0, byte(0)
		var bestMask byte
		for transform := byte(0); transform < 3; transform++ {
			key := adaptiveKey(input, pos, transform)
			entry, found := idx[transform][key]
			if !found {
				continue
			}
			for k := 0; k < len(entry.vals) && k < entry.count; k++ {
				from := entry.vals[k]
				var mask byte
				if transform == 1 {
					mask = input[pos] ^ input[from]
				}
				if transform == 2 {
					mask = input[pos] - input[from]
				}
				if transform != 0 && mask == 0 {
					continue
				}
				maxLen := n - pos
				if maxLen > 4096 {
					maxLen = 4096
				}
				if pos-from < maxLen {
					maxLen = pos - from
				}
				if maxLen <= bestLen {
					continue
				}
				size := 8 // anchor invariant already checked
				for size < maxLen {
					v := input[from+size]
					switch transform {
					case 1:
						v ^= mask
					case 2:
						v += mask
					}
					if v != input[pos+size] {
						break
					}
					size++
				}
				if size > bestLen {
					bestLen, bestFrom, bestTransform, bestMask = size, from, transform, mask
				}
			}
		}
		// The token's reference costs at most eight bytes for this range.
		if bestLen < 32 {
			pos++
			continue
		}
		flush(pos)
		opcode := adaptiveCopy
		if bestTransform == 1 {
			opcode = adaptiveXor
		} else if bestTransform == 2 {
			opcode = adaptiveAdd
		}
		cmds.WriteByte(opcode)
		putVar(&cmds, uint64(bestFrom))
		putVar(&cmds, uint64(bestLen))
		if bestTransform != 0 {
			cmds.WriteByte(bestMask)
		}
		pos += bestLen
		litStart = pos
		hasRefs = true
	}
	flush(n)
	if !hasRefs {
		return nil, nil
	}
	payload, err := nativeBytes(cmds.Bytes())
	if err != nil {
		return nil, err
	}
	var archive bytes.Buffer
	archive.WriteString(magic)
	archive.WriteByte(modeAdaptive)
	putVar(&archive, uint64(n))
	archive.Write(digest[:])
	archive.WriteByte(1) // in-project native-packed adaptive tokens
	putVar(&archive, uint64(len(payload)))
	archive.Write(payload)
	return archive.Bytes(), nil
}
