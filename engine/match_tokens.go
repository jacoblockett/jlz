package engine

// MDC5 Stage 2: standalone native sequence coding, no external compressor.
// The v2 match parser is reused without a second search. Metadata and literal
// symbols are encoded separately; legacy-style v2 streams use the in-project packer under JLZ1 framing.

import (
	"bytes"
	"errors"
	"io"
)

const (
	modeMatchNative     byte = 11
	modeMatchNativePred byte = 13
)

type nativeStreams struct {
	ncommands                                    int
	op, lengths, distances, transforms, literals []byte
}

func splitMatchTokens(tokens []byte, originalLen int) (nativeStreams, error) {
	r := bytes.NewReader(tokens)
	var st nativeStreams
	reconstructed := 0
	for r.Len() > 0 {
		op, err := r.ReadByte()
		if err != nil {
			return st, err
		}
		if op > matchRepAdd {
			return st, errors.New("bad native opcode")
		}
		if st.ncommands%2 == 0 {
			st.op = append(st.op, op)
		} else {
			st.op[len(st.op)-1] |= op << 4
		}
		st.ncommands++
		if op == matchLiteral {
			n, err := readVar(r)
			if err != nil {
				return st, err
			}
			if n == 0 || n > uint64(r.Len()) || n > uint64(originalLen-reconstructed) {
				return st, errors.New("bad native literal")
			}
			st.lengths = appendUvar(st.lengths, n)
			start := len(tokens) - r.Len()
			st.literals = append(st.literals, tokens[start:start+int(n)]...)
			if _, err = r.Seek(int64(n), io.SeekCurrent); err != nil {
				return st, err
			}
			reconstructed += int(n)
			continue
		}
		if op == matchCopy || op == matchXor || op == matchAdd {
			d, err := readVar(r)
			if err != nil {
				return st, err
			}
			st.distances = appendUvar(st.distances, d)
		} else if op == matchRepXor || op == matchRepAdd {
			idx, err := r.ReadByte()
			if err != nil {
				return st, err
			}
			st.distances = append(st.distances, idx)
		}
		n, err := readVar(r)
		if err != nil {
			return st, err
		}
		if n == 0 || n > uint64(originalLen-reconstructed) {
			return st, errors.New("bad native length")
		}
		st.lengths = appendUvar(st.lengths, n)
		reconstructed += int(n)
		if op == matchXor || op == matchAdd || op == matchRepXor || op == matchRepAdd {
			mask, err := r.ReadByte()
			if err != nil {
				return st, err
			}
			st.transforms = append(st.transforms, mask)
		}
	}
	if reconstructed != originalLen {
		return st, errors.New("native token coverage mismatch")
	}
	return st, nil
}

func appendUvar(b []byte, v uint64) []byte {
	for v >= 128 {
		b = append(b, byte(v)|0x80)
		v >>= 7
	}
	return append(b, byte(v))
}

func encodeMatchStage2(input []byte, digest [32]byte) ([]byte, string, error) {
	tokens, err := tokenizeMatchV2(input)
	if err != nil || len(tokens) == 0 {
		return nil, "", err
	}
	native, err := encodeMatchNativeTokens(input, digest, tokens)
	if err != nil {
		return nil, "", err
	}
	legacy, err := encodeMatchV2Tokens(input, digest, tokens)
	if err != nil {
		return nil, "", err
	}
	best, label := legacy, "match-v2"
	if len(native) < len(best) {
		best, label = native, "match-native"
	}
	pred, err := encodeMatchNativeVariant(input, digest, tokens, true)
	if err != nil {
		return nil, "", err
	}
	if len(pred) < len(best) {
		best, label = pred, "match-native-pred"
	}
	packed, packedLabel, err := bestPackedMatch(input, digest, tokens)
	if err != nil {
		return nil, "", err
	}
	if len(packed) < len(best) {
		best, label = packed, packedLabel
	}
	return best, label, nil
}

func encodeMatchNativeTokens(input []byte, digest [32]byte, tokens []byte) ([]byte, error) {
	return encodeMatchNativeVariant(input, digest, tokens, false)
}

func encodeMatchNativeVariant(input []byte, digest [32]byte, tokens []byte, matched bool) ([]byte, error) {
	st, err := splitMatchTokens(tokens, len(input))
	if err != nil {
		return nil, err
	}
	if matched {
		st.literals, err = matchedLiteralResiduals(tokens, input)
		if err != nil {
			return nil, err
		}
	}
	codec, table, encoded := encodeNativeLiterals(st.literals)
	var b bytes.Buffer
	b.WriteString(magic)
	if matched {
		b.WriteByte(modeMatchNativePred)
	} else {
		b.WriteByte(modeMatchNative)
	}
	putVar(&b, uint64(len(input)))
	b.Write(digest[:])
	for _, value := range []int{st.ncommands, len(st.literals), len(st.lengths), len(st.distances), len(st.transforms), len(encoded)} {
		putVar(&b, uint64(value))
	}
	b.WriteByte(codec)
	b.Write(st.op)
	b.Write(st.lengths)
	b.Write(st.distances)
	b.Write(st.transforms)
	if codec == 1 {
		b.Write(table[:])
	}
	b.Write(encoded)
	return b.Bytes(), nil
}

func decodeMatchNative(r *bytes.Reader, n uint64) ([]byte, error) {
	return decodeMatchNativeVariant(r, n, false)
}

func decodeMatchNativeVariant(r *bytes.Reader, n uint64, matched bool) ([]byte, error) {
	return decodeMatchNativeVariantHistory(r, n, matched, nil)
}

func decodeMatchNativeVariantHistory(r *bytes.Reader, n uint64, matched bool, history []byte) ([]byte, error) {
	var counts [6]int
	for i := range counts {
		v, err := readVar(r)
		if err != nil {
			return nil, err
		}
		if v > n*8+1024 {
			return nil, errors.New("oversized native stream")
		}
		counts[i] = int(v)
	}
	if counts[0] == 0 && n != 0 {
		return nil, errors.New("missing native commands")
	}
	if uint64(counts[0]) > n || uint64(counts[1]) > n {
		return nil, errors.New("invalid native counts")
	}
	codec, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	opSize := (counts[0] + 1) / 2
	expected := int64(opSize) + int64(counts[2]) + int64(counts[3]) + int64(counts[4]) + int64(counts[5])
	if codec == 1 {
		expected += 128
	} else if codec != 0 {
		return nil, errors.New("unknown native literal codec")
	}
	if expected != int64(r.Len()) {
		return nil, errors.New("native payload length mismatch")
	}
	take := func(sz int) ([]byte, error) { a := make([]byte, sz); _, e := io.ReadFull(r, a); return a, e }
	ops, err := take(opSize)
	if err != nil {
		return nil, err
	}
	sizes, err := take(counts[2])
	if err != nil {
		return nil, err
	}
	distances, err := take(counts[3])
	if err != nil {
		return nil, err
	}
	transforms, err := take(counts[4])
	if err != nil {
		return nil, err
	}
	var table [128]byte
	if codec == 1 {
		if _, err = io.ReadFull(r, table[:]); err != nil {
			return nil, err
		}
	}
	payload, err := take(counts[5])
	if err != nil {
		return nil, err
	}
	var literals []byte
	if codec == 0 {
		if len(payload) != counts[1] {
			return nil, errors.New("wrong native raw literal size")
		}
		literals = payload
	} else {
		literals, err = decodeNativeLiterals(table, payload, counts[1])
		if err != nil {
			return nil, err
		}
	}
	lr, dr := bytes.NewReader(sizes), bytes.NewReader(distances)
	out := make([]byte, len(history)+int(n))
	copy(out, history)
	li, ti, at := 0, 0, len(history)
	var recent [3]int
	lastKind, lastMask := byte(0), byte(0)
	for i := 0; i < counts[0]; i++ {
		op := ops[i/2] & 15
		if i&1 != 0 {
			op = ops[i/2] >> 4
		}
		if op > matchRepAdd {
			return nil, errors.New("invalid native match command")
		}
		if op == matchLiteral {
			count, e := readVar(lr)
			if e != nil {
				return nil, e
			}
			if count == 0 || count > uint64(len(out)-at) || count > uint64(len(literals)-li) {
				return nil, errors.New("invalid native literal run")
			}
			if !matched {
				copy(out[at:at+int(count)], literals[li:li+int(count)])
			} else {
				for j := 0; j < int(count); j++ {
					v := literals[li+j]
					distance := recent[0]
					if distance > 0 && distance <= at+j {
						predicted := out[at+j-distance]
						if lastKind == 1 {
							predicted ^= lastMask
						} else if lastKind == 2 {
							predicted += lastMask
						}
						v ^= predicted
					}
					out[at+j] = v
				}
			}
			li += int(count)
			at += int(count)
			continue
		}
		var distance int
		if op >= matchRep0 && op <= matchRep2 {
			distance = recent[op-matchRep0]
		} else if op == matchRepXor || op == matchRepAdd {
			ri, e := dr.ReadByte()
			if e != nil {
				return nil, e
			}
			if ri >= 3 {
				return nil, errors.New("bad native repeat index")
			}
			distance = recent[ri]
		} else {
			d, e := readVar(dr)
			if e != nil {
				return nil, e
			}
			if d > uint64(at) {
				return nil, errors.New("bad native distance")
			}
			distance = int(d)
		}
		length, e := readVar(lr)
		if e != nil {
			return nil, e
		}
		if distance <= 0 || distance > at || distance > matchMaxDistance || length == 0 || length > matchMaximum || length > uint64(len(out)-at) {
			return nil, errors.New("bad native reference")
		}
		kind := byte(0)
		if op == matchXor || op == matchRepXor {
			kind = 1
		}
		if op == matchAdd || op == matchRepAdd {
			kind = 2
		}
		var mask byte
		if kind != 0 {
			if ti >= len(transforms) {
				return nil, io.ErrUnexpectedEOF
			}
			mask = transforms[ti]
			ti++
			if mask == 0 {
				return nil, errors.New("redundant native transform")
			}
		}
		expandMatch(out, at, distance, int(length), kind, mask)
		at += int(length)
		lastKind, lastMask = kind, mask
		updateMatchRecent(&recent, distance)
	}
	if at != len(out) || li != len(literals) || ti != len(transforms) || lr.Len() != 0 || dr.Len() != 0 {
		return nil, errors.New("unused or missing native tokens")
	}
	if counts[0]%2 != 0 && ops[len(ops)-1]>>4 != 0 {
		return nil, errors.New("invalid native opcode padding")
	}
	return out[len(history):], nil
}

type huffNode struct {
	freq        int
	symbol      int
	left, right *huffNode
}

// The previous canonical length builder repeatedly sorted all nodes.
// This stable min-heap uses the same (frequency, lowest symbol) ordering.
func huffLess(a, b *huffNode) bool {
	return a.freq < b.freq || (a.freq == b.freq && a.symbol < b.symbol)
}
func huffDown(h []*huffNode, i int) {
	n := len(h)
	for {
		child := i*2 + 1
		if child >= n {
			return
		}
		if child+1 < n && huffLess(h[child+1], h[child]) {
			child++
		}
		if !huffLess(h[child], h[i]) {
			return
		}
		h[i], h[child] = h[child], h[i]
		i = child
	}
}
func huffHeapify(h []*huffNode) {
	for i := len(h)/2 - 1; i >= 0; i-- {
		huffDown(h, i)
	}
}
func huffPop(h *[]*huffNode) *huffNode {
	a := *h
	top := a[0]
	last := a[len(a)-1]
	a = a[:len(a)-1]
	if len(a) > 0 {
		a[0] = last
		huffDown(a, 0)
	}
	*h = a
	return top
}
func huffPush(h *[]*huffNode, node *huffNode) {
	a := append(*h, node)
	i := len(a) - 1
	for i > 0 {
		parent := (i - 1) / 2
		if !huffLess(a[i], a[parent]) {
			break
		}
		a[i], a[parent] = a[parent], a[i]
		i = parent
	}
	*h = a
}

func nativeLengths(data []byte) (lengths [256]byte, ok bool) {
	if len(data) == 0 {
		return lengths, false
	}
	var frequency [256]int
	for _, b := range data {
		frequency[b]++
	}
	// This project's original Huffman builder refused any frequency
	// distribution yielding a depth over 15. On skewed natural-language
	// bytes that silently disabled compression entirely. Preserve the
	// original optimal lengths whenever valid; otherwise rebalance the
	// same in-project tree using a lower bound on rare-symbol weights.
	// This is an independently implemented, length-limited approximation:
	// all symbols remain exactly reconstructible, and the on-disk format
	// and decoder remain unchanged.
	for floor := 0; ; {
		nodes := make([]*huffNode, 0, 256)
		for i, c := range frequency {
			if c > 0 {
				if c < floor {
					c = floor
				}
				nodes = append(nodes, &huffNode{freq: c, symbol: i})
			}
		}
		if len(nodes) == 1 {
			lengths[nodes[0].symbol] = 1
			return lengths, true
		}
		huffHeapify(nodes)
		for len(nodes) > 1 {
			a := huffPop(&nodes)
			b := huffPop(&nodes)
			huffPush(&nodes, &huffNode{freq: a.freq + b.freq, symbol: min(a.symbol, b.symbol), left: a, right: b})
		}
		var provisional [256]byte
		valid := true
		var visit func(*huffNode, int)
		visit = func(node *huffNode, depth int) {
			if node.left == nil && node.right == nil {
				if depth > 15 {
					valid = false
					return
				}
				provisional[node.symbol] = byte(depth)
				return
			}
			visit(node.left, depth+1)
			visit(node.right, depth+1)
		}
		visit(nodes[0], 0)
		if valid {
			return provisional, true
		}
		if floor == 0 {
			floor = max(1, len(data)/1024)
		} else {
			floor *= 2
		}
		// Uniform frequencies always fit within 8 bits for 256 symbols.
		// Clamping avoids integer overflow on very large inputs.
		if floor > len(data) {
			floor = len(data)
		}
	}
}

func nativeCodes(lengths [256]byte) (codes [256]uint16, first [16]uint32, starts [16]int, counts [16]int, ordered []byte, err error) {
	for _, l := range lengths {
		if l > 15 {
			return codes, first, starts, counts, nil, errors.New("huffman depth out of range")
		}
		if l > 0 {
			counts[l]++
		}
	}
	code := uint32(0)
	for l := 1; l < 16; l++ {
		code = (code + uint32(counts[l-1])) << 1
		if code+uint32(counts[l]) > uint32(1<<l) {
			return codes, first, starts, counts, nil, errors.New("oversubscribed huffman tree")
		}
		first[l] = code
		starts[l] = len(ordered)
		for s, lval := range lengths {
			if int(lval) == l {
				ordered = append(ordered, byte(s))
			}
		}
		for j := 0; j < counts[l]; j++ {
			codes[ordered[starts[l]+j]] = uint16(code + uint32(j))
		}
	}
	if len(ordered) == 0 {
		return codes, first, starts, counts, nil, errors.New("empty huffman code")
	}
	return codes, first, starts, counts, ordered, nil
}

func encodeNativeLiterals(data []byte) (byte, [128]byte, []byte) {
	var table [128]byte
	if len(data) < 384 {
		return 0, table, data
	}
	lengths, ok := nativeLengths(data)
	if !ok {
		return 0, table, data
	}
	codes, _, _, _, _, err := nativeCodes(lengths)
	if err != nil {
		return 0, table, data
	}
	bits := 0
	for _, b := range data {
		bits += int(lengths[b])
	}
	if (bits+7)/8+128 >= len(data) {
		return 0, table, data
	}
	for i, l := range lengths {
		if i&1 == 0 {
			table[i/2] = l
		} else {
			table[i/2] |= l << 4
		}
	}
	dst := make([]byte, (bits+7)/8)
	pos := 0
	for _, b := range data {
		l := int(lengths[b])
		c := codes[b]
		for k := l - 1; k >= 0; k-- {
			if (c>>k)&1 != 0 {
				dst[pos/8] |= 1 << uint(7-pos%8)
			}
			pos++
		}
	}
	return 1, table, dst
}

func decodeNativeLiterals(table [128]byte, payload []byte, n int) ([]byte, error) {
	var lengths [256]byte
	for i := 0; i < 256; i++ {
		if i&1 == 0 {
			lengths[i] = table[i/2] & 15
		} else {
			lengths[i] = table[i/2] >> 4
		}
	}
	return decodeCanonicalHuffman(lengths, payload, n)
}

// A literal following a reference is represented by its XOR difference from
// that reference's predicted byte. Encoder and decoder can calculate the
// reference and transform from already decoded match commands; no new
// dictionary, offset, or per-literal signaling is transmitted.
func matchedLiteralResiduals(tokens, input []byte) ([]byte, error) {
	r := bytes.NewReader(tokens)
	var recent [3]int
	var out []byte
	at := 0
	lastKind, lastMask := byte(0), byte(0)
	for r.Len() > 0 {
		op, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if op == matchLiteral {
			size, err := readVar(r)
			if err != nil {
				return nil, err
			}
			if size == 0 || size > uint64(r.Len()) || size > uint64(len(input)-at) {
				return nil, errors.New("invalid matched literal")
			}
			for j := 0; j < int(size); j++ {
				b, err := r.ReadByte()
				if err != nil {
					return nil, err
				}
				v := b
				d := recent[0]
				if d > 0 && d <= at+j {
					pred := input[at+j-d]
					if lastKind == 1 {
						pred ^= lastMask
					} else if lastKind == 2 {
						pred += lastMask
					}
					v ^= pred
				}
				out = append(out, v)
			}
			at += int(size)
			continue
		}
		if op > matchRepAdd {
			return nil, errors.New("invalid matched opcode")
		}
		var distance int
		if op >= matchRep0 && op <= matchRep2 {
			distance = recent[op-matchRep0]
		} else if op == matchRepXor || op == matchRepAdd {
			index, err := r.ReadByte()
			if err != nil {
				return nil, err
			}
			if index >= 3 {
				return nil, errors.New("invalid repeat index")
			}
			distance = recent[index]
		} else {
			value, err := readVar(r)
			if err != nil {
				return nil, err
			}
			if value > uint64(at) {
				return nil, errors.New("invalid match distance")
			}
			distance = int(value)
		}
		size, err := readVar(r)
		if err != nil {
			return nil, err
		}
		if size == 0 || size > uint64(len(input)-at) || distance <= 0 || distance > at {
			return nil, errors.New("invalid matched reference")
		}
		lastKind, lastMask = 0, 0
		if op == matchXor || op == matchRepXor {
			lastKind = 1
		} else if op == matchAdd || op == matchRepAdd {
			lastKind = 2
		}
		if lastKind != 0 {
			lastMask, err = r.ReadByte()
			if err != nil {
				return nil, err
			}
		}
		at += int(size)
		updateMatchRecent(&recent, distance)
	}
	if at != len(input) {
		return nil, errors.New("incomplete matched stream")
	}
	return out, nil
}
