package engine

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"testing"
)

func (m *matchIndex) findReference(pos int, recent [3]int) matchChoice {
	n := len(m.input)
	if pos+5 > n {
		return matchChoice{}
	}
	var best matchChoice
	best.repeat = -1
	data := m.input
	// The oracle includes the intentional repeat-offset candidates, while
	// checking that the optimized hash-bucket search preserves its results.
	for _, distance := range recent {
		if distance <= 0 || distance > pos || distance > matchMaxDistance {
			continue
		}
		from := pos - distance
		if !bytes.Equal(data[pos:pos+5], data[from:from+5]) {
			continue
		}
		size := 5
		limit := min(n-pos, matchMaximum)
		for size < limit && data[pos+size] == data[from+size] {
			size++
		}
		c := scoreMatchChoice(matchChoice{length: size, distance: distance}, recent)
		if c.score > best.score {
			best = c
		}
	}
	for typ := byte(0); typ < 3; typ++ {
		key := matchFingerprint(data, pos, typ)
		bucket := &m.heads[typ][matchHash(key)]
		for k := 0; k < 2; k++ {
			value := bucket.latest(k)
			if value == 0 {
				continue
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

func TestShortCircuitMatchSearchEquivalence(t *testing.T) {
	rng := rand.New(rand.NewSource(20261016))
	for caseNum := 0; caseNum < 12; caseNum++ {
		data := make([]byte, 16384)
		switch caseNum % 4 {
		case 0:
			_, _ = rng.Read(data)
		case 1:
			for i := range data {
				data[i] = byte((i / 37) % 29)
			}
		case 2:
			_, _ = rng.Read(data[:1024])
			for i := 1024; i < len(data); i++ {
				data[i] = data[i-1024] ^ byte(caseNum)
			}
		case 3:
			_, _ = rng.Read(data[:512])
			for i := 512; i < len(data); i++ {
				data[i] = data[i-512] + byte(caseNum+1)
			}
		}
		m := newMatchIndex(data)
		var recent [3]int
		for pos := 0; pos < len(data); pos++ {
			m.addBefore(pos)
			a, b := m.findReference(pos, recent), m.find(pos, recent)
			if a != b {
				t.Fatalf("case %d pos %d old=%+v new=%+v", caseNum, pos, a, b)
			}
			if b.score > 0 {
				updateMatchRecent(&recent, b.distance)
			}
		}
	}
}
