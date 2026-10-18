package main

// Graph and lag predictors exploit exact mathematical relationships among
// high-entropy byte blocks. Both are invertible, deterministic transforms.

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"sort"
)

const (
	blockRaw byte = iota
	blockCopy
	blockXorConstant
	blockXorPattern
	blockXorPair
	blockXorResidual
)

type graphOp struct {
	kind    byte
	a, b    int
	mask    []byte
	payload []byte
}
type fingerprint struct {
	sig  uint64
	norm uint64
}

func blockSig(s []byte) uint64 {
	var sig uint64
	for i := 0; i < 8 && i < len(s); i++ {
		sig |= uint64(s[i]) << (8 * i)
	}
	return sig
}
func normHash(s []byte) uint64 {
	// FNV64 over the prefix, invariant under a constant-byte XOR mask.
	var h uint64 = 14695981039346656037
	for i := 1; i < len(s) && i < 33; i++ {
		h ^= uint64(s[i] ^ s[0])
		h *= 1099511628211
	}
	return h
}
func isExact(a, b []byte) bool { return bytes.Equal(a, b) }
func xorConstant(a, b []byte) (byte, bool) {
	if len(a) != len(b) || len(a) == 0 {
		return 0, false
	}
	mask := a[0] ^ b[0]
	for i := 1; i < len(a); i++ {
		if a[i]^b[i] != mask {
			return 0, false
		}
	}
	return mask, true
}
func xorPattern(a, b []byte) ([]byte, bool) {
	if len(a) != len(b) || len(a) < 4 {
		return nil, false
	}
	for period := 2; period <= 16 && period < len(a)/2; period++ {
		// Period 1 is already considered by xorConstant.
		p := make([]byte, period)
		for k := 0; k < period; k++ {
			p[k] = a[k] ^ b[k]
		}
		matches := true
		for k := period; k < len(a); k++ {
			if a[k]^b[k] != p[k%period] {
				matches = false
				break
			}
		}
		if matches {
			return p, true
		}
	}
	return nil, false
}
func xorPairEqual(a, b, c []byte) bool {
	if len(a) != len(b) || len(a) != len(c) {
		return false
	}
	for k := range c {
		if a[k]^b[k] != c[k] {
			return false
		}
	}
	return true
}
func diffScore(a, b []byte) int {
	// Rough savings proxy: how concentrated is the bytewise XOR residual?
	// A high concentration gives the entropy coder a simple target.
	var hist [256]int
	for i := range a {
		hist[a[i]^b[i]]++
	}
	max := 0
	for _, n := range hist {
		if n > max {
			max = n
		}
	}
	return max
}
func selectGraphOps(src []byte, width int) ([]graphOp, string) {
	n := (len(src) + width - 1) / width
	chunks := make([][]byte, n)
	for i := 0; i < n; i++ {
		end := (i + 1) * width
		if end > len(src) {
			end = len(src)
		}
		chunks[i] = src[i*width : end]
	}
	ops := make([]graphOp, 0, n)
	sigs := make(map[uint64][]int, n)
	norms := make(map[uint64][]int, n)
	counts := make(map[byte]int)
	for i, b := range chunks {
		selected := graphOp{kind: blockRaw, payload: b}
		// Exact references, complementary masks, and short periodic XOR masks.
		// Normalized indexes offer O(1) candidate retrieval for constant masks.
		for j := len(norms[normHash(b)]) - 1; j >= 0; j-- {
			k := norms[normHash(b)][j]
			if len(chunks[k]) != len(b) {
				continue
			}
			m, ok := xorConstant(b, chunks[k])
			if !ok {
				continue
			}
			if m == 0 {
				selected = graphOp{kind: blockCopy, a: k}
			} else {
				selected = graphOp{kind: blockXorConstant, a: k, mask: []byte{m}}
			}
			break
		}
		// XOR relations C = A XOR B, using the XOR-linear first-eight-byte key.
		if selected.kind == blockRaw && len(b) >= 8 {
			for j := 0; j < i; j++ {
				want := blockSig(b) ^ blockSig(chunks[j])
				cand := sigs[want]
				found := false
				for _, k := range cand {
					if k >= j || len(chunks[k]) != len(b) || len(chunks[j]) != len(b) {
						continue
					}
					if xorPairEqual(chunks[j], chunks[k], b) {
						selected = graphOp{kind: blockXorPair, a: j, b: k}
						found = true
						break
					}
				}
				if found {
					break
				}
			}
		}
		if selected.kind == blockRaw {
			// Compare with earlier blocks for near-matches, masks, or residuals.
			// Limit the lookback for bounded time on large files.
			bestScore := 0
			bestRef := -1
			from := i - 128
			if from < 0 {
				from = 0
			}
			for j := i - 1; j >= from; j-- {
				if len(chunks[j]) != len(b) {
					continue
				}
				if p, ok := xorPattern(b, chunks[j]); ok {
					selected = graphOp{kind: blockXorPattern, a: j, mask: p}
					break
				}
				score := diffScore(b, chunks[j])
				if score > bestScore {
					bestScore = score
					bestRef = j
				}
			}
			if selected.kind == blockRaw && bestRef >= 0 && bestScore*100 >= len(b)*35 {
				// These heuristics only select the predicted residual. The final
				// actual compressed archive must beat the raw-file fallback.
				residual := make([]byte, len(b))
				for j := range b {
					residual[j] = b[j] ^ chunks[bestRef][j]
				}
				selected = graphOp{kind: blockXorResidual, a: bestRef, payload: residual}
			}
		}
		ops = append(ops, selected)
		counts[selected.kind]++
		sigs[blockSig(b)] = append(sigs[blockSig(b)], i)
		norms[normHash(b)] = append(norms[normHash(b)], i)
	}
	return ops, fmt.Sprintf("chunk=%d; raw=%d, equal=%d, XOR-constant=%d, XOR-periodic=%d, XOR-pair=%d, XOR-residual=%d", width, counts[blockRaw], counts[blockCopy], counts[blockXorConstant], counts[blockXorPattern], counts[blockXorPair], counts[blockXorResidual])
}
func encodeGraph(src []byte, digest [32]byte, width int) ([]byte, string, error) {
	ops, detail := selectGraphOps(src, width)
	var meaningful bool
	for _, op := range ops {
		if op.kind != blockRaw {
			meaningful = true
			break
		}
	}
	if !meaningful {
		return nil, "", nil
	}
	var out bytes.Buffer
	out.WriteString(magic)
	out.WriteByte(modeGraph)
	putVar(&out, uint64(len(src)))
	out.Write(digest[:])
	putVar(&out, uint64(width))
	fw, err := newNativeWriter(&out)
	if err != nil {
		return nil, "", err
	}
	bw := bufio.NewWriterSize(fw, 64*1024)
	for _, op := range ops {
		if err := bw.WriteByte(op.kind); err != nil {
			return nil, "", err
		}
		switch op.kind {
		case blockCopy, blockXorConstant, blockXorPattern, blockXorResidual:
			if err := putVar(bw, uint64(op.a)); err != nil {
				return nil, "", err
			}
		case blockXorPair:
			if err := putVar(bw, uint64(op.a)); err != nil {
				return nil, "", err
			}
			if err := putVar(bw, uint64(op.b)); err != nil {
				return nil, "", err
			}
		}
		switch op.kind {
		case blockXorConstant:
			if err := bw.WriteByte(op.mask[0]); err != nil {
				return nil, "", err
			}
		case blockXorPattern:
			if err := bw.WriteByte(byte(len(op.mask))); err != nil {
				return nil, "", err
			}
			if _, err := bw.Write(op.mask); err != nil {
				return nil, "", err
			}
		case blockRaw, blockXorResidual:
			if _, err := bw.Write(op.payload); err != nil {
				return nil, "", err
			}
		}
	}
	if err := bw.Flush(); err != nil {
		return nil, "", err
	}
	if err := fw.Close(); err != nil {
		return nil, "", err
	}
	return out.Bytes(), detail, nil
}
func bestGraph(src []byte, digest [32]byte) ([]byte, string, error) {
	widths := []int{128, 192, 256, 384, 512, 768, 1024, 1536, 2048, 3072, 4096, 8192}
	var best []byte
	detail := ""
	for _, w := range widths {
		if w > len(src) || (len(src)+w-1)/w > 512 {
			continue
		}
		arc, d, err := encodeGraph(src, digest, w)
		if err != nil {
			return nil, "", err
		}
		if len(arc) > 0 && (best == nil || len(arc) < len(best)) {
			best, detail = arc, d
		}
	}
	return best, detail, nil
}
func decodeGraph(reader *bytes.Reader, n uint64) ([]byte, error) {
	width64, err := readVar(reader)
	if err != nil {
		return nil, err
	}
	if width64 == 0 || width64 > n+1 || width64 > math.MaxInt32 {
		return nil, errors.New("invalid graph width")
	}
	width := int(width64)
	chunkCount := (int(n) + width - 1) / width
	if chunkCount > 1000000 {
		return nil, errors.New("excessive graph chunks")
	}
	out := make([]byte, int(n))
	fl := newNativeReader(reader)
	defer fl.Close()
	br := bufio.NewReader(fl)
	for i := 0; i < chunkCount; i++ {
		offset := i * width
		count := width
		if offset+count > len(out) {
			count = len(out) - offset
		}
		dest := out[offset : offset+count]
		kind, err := br.ReadByte()
		if err != nil {
			return nil, err
		}
		var ref1, ref2 uint64
		if kind == blockCopy || kind == blockXorConstant || kind == blockXorPattern || kind == blockXorResidual || kind == blockXorPair {
			ref1, err = readVar(br)
			if err != nil {
				return nil, err
			}
			if ref1 >= uint64(i) {
				return nil, errors.New("forward graph reference")
			}
		}
		if kind == blockXorPair {
			ref2, err = readVar(br)
			if err != nil {
				return nil, err
			}
			if ref2 >= uint64(i) {
				return nil, errors.New("forward XOR reference")
			}
		}
		getRef := func(r uint64) ([]byte, error) {
			st := int(r) * width
			if st+count > offset {
				return nil, errors.New("graph reference crosses boundary")
			}
			return out[st : st+count], nil
		}
		switch kind {
		case blockRaw:
			if _, err = io.ReadFull(br, dest); err != nil {
				return nil, err
			}
		case blockCopy, blockXorConstant, blockXorPattern, blockXorResidual:
			previous, e := getRef(ref1)
			if e != nil {
				return nil, e
			}
			switch kind {
			case blockCopy:
				copy(dest, previous)
			case blockXorConstant:
				mask, e := br.ReadByte()
				if e != nil {
					return nil, e
				}
				for k := range dest {
					dest[k] = previous[k] ^ mask
				}
			case blockXorPattern:
				period, e := br.ReadByte()
				if e != nil {
					return nil, e
				}
				if period < 2 || period > 16 {
					return nil, errors.New("bad graph mask period")
				}
				mask := make([]byte, int(period))
				if _, e = io.ReadFull(br, mask); e != nil {
					return nil, e
				}
				for k := range dest {
					dest[k] = previous[k] ^ mask[k%len(mask)]
				}
			case blockXorResidual:
				if _, e = io.ReadFull(br, dest); e != nil {
					return nil, e
				}
				for k := range dest {
					dest[k] ^= previous[k]
				}
			}
		case blockXorPair:
			a, e := getRef(ref1)
			if e != nil {
				return nil, e
			}
			b, e := getRef(ref2)
			if e != nil {
				return nil, e
			}
			for k := range dest {
				dest[k] = a[k] ^ b[k]
			}
		default:
			return nil, errors.New("unknown graph operation")
		}
	}
	_, err = br.ReadByte()
	if err != io.EOF {
		if err == nil {
			return nil, errors.New("extra graph bytes")
		}
		return nil, err
	}
	return out, nil
}

// XOR and subtraction-lag predictors are fully reversible: the decoder
// rebuilds each byte from a previous byte and the stored residual.
const (
	lagXOR byte = iota
	lagSubtract
)

type lagCandidate struct {
	lag   int
	op    byte
	score float64
}

func candidateLags(src []byte) []int {
	base := []int{1, 2, 3, 4, 8, 16, 32, 64, 128, 192, 256, 384, 512, 768, 1024, 1536, 2048, 3072, 4096, 6144, 8192, 12288, 16384}
	counts := make(map[int]int)
	// Adjacent-byte-XOR windows are invariant to a uniform XOR mask.
	// Repeating them across the input yields candidate nonlocal boundaries.
	var last = make(map[uint64]int)
	for i := 0; i+9 <= len(src); i += 4 {
		var sig uint64
		for j := 0; j < 8; j++ {
			sig |= uint64(src[i+j]^src[i+j+1]) << (8 * j)
		}
		if prev, ok := last[sig]; ok {
			lag := i - prev
			if lag >= 16 && lag <= 32768 {
				counts[lag]++
			}
		}
		last[sig] = i
	}
	type kv struct{ lag, count int }
	rank := make([]kv, 0, len(counts))
	for lag, count := range counts {
		rank = append(rank, kv{lag, count})
	}
	sort.Slice(rank, func(i, j int) bool {
		if rank[i].count != rank[j].count {
			return rank[i].count > rank[j].count
		}
		return rank[i].lag < rank[j].lag
	})
	for i := 0; i < len(rank) && i < 10; i++ {
		base = append(base, rank[i].lag)
	}
	seen := make(map[int]bool)
	out := make([]int, 0, len(base))
	for _, n := range base {
		if n > 0 && n < len(src)/2 && !seen[n] {
			seen[n] = true
			out = append(out, n)
		}
	}
	sort.Ints(out)
	return out
}
func lagResidual(src []byte, lag int, op byte) ([]byte, float64) {
	res := make([]byte, len(src))
	var hist [256]int
	for i, b := range src {
		var r byte
		if i < lag {
			r = b
		} else if op == lagXOR {
			r = b ^ src[i-lag]
		} else {
			r = b - src[i-lag]
		}
		res[i] = r
		hist[r]++
	}
	entropy := 0.0
	for _, h := range hist {
		if h > 0 {
			p := float64(h) / float64(len(src))
			entropy -= p * math.Log2(p)
		}
	}
	// Initial lag bytes are untransformed; penalize very large lags.
	return res, entropy * float64(len(src)) / 8
}
func encodeLag(src []byte, digest [32]byte, lag int, op byte, payload []byte) ([]byte, error) {
	var out bytes.Buffer
	out.WriteString(magic)
	out.WriteByte(modeRelation)
	putVar(&out, uint64(len(src)))
	out.Write(digest[:])
	out.WriteByte(op)
	putVar(&out, uint64(lag))
	fw, err := newNativeWriter(&out)
	if err != nil {
		return nil, err
	}
	if _, err = fw.Write(payload); err != nil {
		return nil, err
	}
	if err = fw.Close(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
func bestLag(src []byte, digest [32]byte) ([]byte, string, error) {
	var candidates []lagCandidate
	for _, lag := range candidateLags(src) {
		for op := byte(0); op < 2; op++ {
			_, score := lagResidual(src, lag, op)
			candidates = append(candidates, lagCandidate{lag, op, score})
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].score != candidates[j].score {
			return candidates[i].score < candidates[j].score
		}
		if candidates[i].lag != candidates[j].lag {
			return candidates[i].lag < candidates[j].lag
		}
		return candidates[i].op < candidates[j].op
	})
	var best []byte
	description := ""
	// Entropy estimate chooses candidates, but only the actual encoded archive
	// decides whether a transform was beneficial.
	for i := 0; i < len(candidates) && i < 10; i++ {
		c := candidates[i]
		residual, _ := lagResidual(src, c.lag, c.op)
		arc, err := encodeLag(src, digest, c.lag, c.op, residual)
		if err != nil {
			return nil, "", err
		}
		if len(arc) > 0 && (best == nil || len(arc) < len(best)) {
			best = arc
			name := "XOR"
			if c.op == lagSubtract {
				name = "subtract"
			}
			description = fmt.Sprintf("%s residual with backward lag %d bytes", name, c.lag)
		}
	}
	return best, description, nil
}
func decodeLag(reader *bytes.Reader, n uint64) ([]byte, error) {
	op, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if op > lagSubtract {
		return nil, errors.New("unknown lag operation")
	}
	lag64, err := readVar(reader)
	if err != nil {
		return nil, err
	}
	if lag64 == 0 || lag64 >= n {
		return nil, errors.New("invalid lag")
	}
	lag := int(lag64)
	fl := newNativeReader(reader)
	defer fl.Close()
	out := make([]byte, int(n))
	if _, err = io.ReadFull(fl, out); err != nil {
		return nil, err
	}
	var extra [1]byte
	extraN, e := fl.Read(extra[:])
	if extraN != 0 || e != io.EOF {
		if e == nil {
			return nil, errors.New("extra lag data")
		}
		return nil, e
	}
	for i := lag; i < len(out); i++ {
		if op == lagXOR {
			out[i] ^= out[i-lag]
		} else {
			out[i] += out[i-lag]
		}
	}
	return out, nil
}

// Used by test suites to verify that relation matching uses true XOR,
// not a hash collision or an accidental approximate equality.

// Ordinary entropy coding baseline, with SHA-256 coverage supplied by header.
func encodeNativeBytes(src []byte, digest [32]byte) ([]byte, error) {
	var out bytes.Buffer
	out.WriteString(magic)
	out.WriteByte(modeNativeBytes)
	if err := putVar(&out, uint64(len(src))); err != nil {
		return nil, err
	}
	out.Write(digest[:])
	fw, err := newNativeWriter(&out)
	if err != nil {
		return nil, err
	}
	if _, err = fw.Write(src); err != nil {
		return nil, err
	}
	if err = fw.Close(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
func decodeNativeBytes(r *bytes.Reader, n uint64) ([]byte, error) {
	fr := newNativeReader(r)
	defer fr.Close()
	out := make([]byte, int(n))
	if _, err := io.ReadFull(fr, out); err != nil {
		return nil, err
	}
	var extra [1]byte
	more, err := fr.Read(extra[:])
	if more != 0 || err != io.EOF {
		if err == nil {
			return nil, errors.New("extra native stream bytes")
		}
		return nil, err
	}
	return out, nil
}
