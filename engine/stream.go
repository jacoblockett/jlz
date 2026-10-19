package engine

// Compact, independently coded native streams. The prior native format stays
// readable. Streams use raw storage, single-byte RLE, or sparse-table canonical
// Huffman according to measured byte cost; no DEFLATE backend is involved.
import (
	"bytes"
	"errors"
	"io"
)

const (
	modeMatchPacked     byte = 14
	modeMatchPackedPred byte = 15
)

func packStream(data []byte) []byte {
	// Each stream begins with raw and payload lengths followed by a codec byte.
	best := append([]byte(nil), data...)
	codec := byte(0)
	if len(data) > 2 {
		all := true
		for _, v := range data[1:] {
			if v != data[0] {
				all = false
				break
			}
		}
		if all {
			best = []byte{data[0]}
			codec = 1
		}
	}
	if len(data) > 16 && codec == 0 {
		lengths, ok := nativeLengths(data)
		if ok {
			codes, _, _, _, _, err := nativeCodes(lengths)
			if err == nil {
				cnt, bits := 0, 0
				for _, l := range lengths {
					if l > 0 {
						cnt++
					}
				}
				for _, v := range data {
					bits += int(lengths[v])
				}
				// The sparse canonical table stores (symbol,length) for nonzero symbols.
				header := varSize(uint64(cnt)) + 2*cnt
				if header+(bits+7)/8 < len(best) {
					var out bytes.Buffer
					putVar(&out, uint64(cnt))
					for v, l := range lengths {
						if l > 0 {
							out.WriteByte(byte(v))
							out.WriteByte(l)
						}
					}
					payload := make([]byte, (bits+7)/8)
					// Huffman codes are at most 15 bits. Keep at most 7
					// pending bits and write complete bytes, preserving
					// exactly the prior MSB-first packed format.
					var pending uint64
					pendingBits := 0
					at := 0
					for _, v := range data {
						l := int(lengths[v])
						pending = (pending << uint(l)) | uint64(codes[v])
						pendingBits += l
						for pendingBits >= 8 {
							pendingBits -= 8
							payload[at] = byte(pending >> uint(pendingBits))
							at++
							pending &= (uint64(1) << uint(pendingBits)) - 1
						}
					}
					if pendingBits > 0 {
						payload[at] = byte(pending << uint(8-pendingBits))
					}
					out.Write(payload)
					if out.Len() < len(best) {
						best = out.Bytes()
						codec = 2
					}
				}
			}
		}
	}
	var result bytes.Buffer
	putVar(&result, uint64(len(data)))
	putVar(&result, uint64(len(best)))
	result.WriteByte(codec)
	result.Write(best)
	return result.Bytes()
}

// packLiteralStream reuses MDC5's established order-1 arithmetic coder when
// it beats the existing raw/RLE/Huffman representation INCLUDING framing.
// Other native metadata streams keep their cheaper packed representation.
func packLiteralStream(data []byte) []byte {
	best := packStream(data)
	// Limit context coding to the native match regions where its measured
	// gains justify the extra CPU. Large whole-file match candidates should
	// not spend seconds encoding literals only to lose to the baseline.
	if len(data) < 384 || len(data) > 256*1024 {
		return best
	}
	coded := arithmeticEncode(data)
	var candidate bytes.Buffer
	putVar(&candidate, uint64(len(data)))
	putVar(&candidate, uint64(len(coded)))
	candidate.WriteByte(3) // order-1 adaptive arithmetic; no learned table transmitted
	candidate.Write(coded)
	if candidate.Len() < len(best) {
		return candidate.Bytes()
	}
	return best
}

func unpackStream(r *bytes.Reader, limit int) ([]byte, error) {
	n, err := readVar(r)
	if err != nil {
		return nil, err
	}
	length, err := readVar(r)
	if err != nil {
		return nil, err
	}
	if n > uint64(limit) || length > uint64(r.Len()) {
		return nil, errors.New("packed stream bounds")
	}
	codec, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	if length > uint64(r.Len()) {
		return nil, io.ErrUnexpectedEOF
	}
	buf := make([]byte, int(length))
	if _, err = io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	switch codec {
	case 0:
		if uint64(len(buf)) != n {
			return nil, errors.New("packed raw stream size")
		}
		return buf, nil
	case 1:
		if len(buf) != 1 || n == 0 {
			return nil, errors.New("packed run invalid")
		}
		return bytes.Repeat(buf, int(n)), nil
	case 2:
		h := bytes.NewReader(buf)
		cnt, err := readVar(h)
		if err != nil {
			return nil, err
		}
		if cnt == 0 || cnt > 256 || uint64(h.Len()) < cnt*2 {
			return nil, errors.New("packed huffman table")
		}
		var lengths [256]byte
		for i := uint64(0); i < cnt; i++ {
			sym, e := h.ReadByte()
			if e != nil {
				return nil, e
			}
			l, e := h.ReadByte()
			if e != nil {
				return nil, e
			}
			if l == 0 || l > 15 || lengths[sym] != 0 {
				return nil, errors.New("packed huffman duplicate/length")
			}
			lengths[sym] = l
		}
		coded := buf[len(buf)-h.Len():]
		return decodeCanonicalHuffman(lengths, coded, int(n))
	case 4:
		return unpackVarintSlots(buf, int(n))
	case 3:
		if n == 0 || len(buf) == 0 {
			return nil, errors.New("packed arithmetic stream invalid")
		}
		return arithmeticDecode(buf, int(n))
	default:
		return nil, errors.New("packed stream codec")
	}
}

func encodeMatchPackedVariant(input []byte, digest [32]byte, tokens []byte, matched bool) ([]byte, error) {
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
	var out bytes.Buffer
	out.WriteString(magic)
	if matched {
		out.WriteByte(modeMatchPackedPred)
	} else {
		out.WriteByte(modeMatchPacked)
	}
	putVar(&out, uint64(len(input)))
	out.Write(digest[:])
	putVar(&out, uint64(st.ncommands))
	for i, stream := range [][]byte{st.op, st.lengths, st.distances, st.transforms} {
		if i == 2 {
			out.Write(packVarintStream(stream))
		} else {
			out.Write(packStream(stream))
		}
	}
	// If almost all original bytes are literals, the special match representation
	// is highly unlikely to beat the selected whole-file baseline. Avoid spending CPU on an
	// additional context coder for such a weak match candidate.
	if len(st.literals)*4 > len(input)*3 {
		out.Write(packStream(st.literals))
	} else {
		out.Write(packLiteralStream(st.literals))
	}
	return out.Bytes(), nil
}

func decodeMatchPacked(r *bytes.Reader, n uint64, matched bool) ([]byte, error) {
	return decodeMatchPackedHistory(r, n, matched, nil)
}

func decodeMatchPackedHistory(r *bytes.Reader, n uint64, matched bool, history []byte) ([]byte, error) {
	return decodePackedWithLengthMode(r, n, matched, history, false)
}

// The split representation shares the same token grammar, but encodes literal
// and reference lengths with independent statistical models. Rejoin their
// varints before invoking the existing, exhaustively tested token decoder.
func decodePackedWithLengthMode(r *bytes.Reader, n uint64, matched bool, history []byte, split bool) ([]byte, error) {
	cmds, err := readVar(r)
	if err != nil {
		return nil, err
	}
	if cmds > n || cmds == 0 && n != 0 {
		return nil, errors.New("packed command count")
	}
	var streams [5][]byte
	// All raw streams are bounded individually by a small multiple of output size.
	maxCount := int(n)*8 + 1024
	for i := range streams {
		streams[i], err = unpackStream(r, maxCount)
		if err != nil {
			return nil, err
		}
	}
	if split {
		// Six packed streams: opcode, literal lengths, copy lengths, distance,
		// transforms, and literal bytes. Retain the previous five-stream mode.
		extra := streams[4]
		streams[4], err = unpackStream(r, maxCount)
		if err != nil {
			return nil, err
		}
		merged, e := mergeOperationLengths(streams[0], int(cmds), streams[1], streams[2])
		if e != nil {
			return nil, e
		}
		streams[1], streams[2], streams[3] = merged, streams[3], extra
	}
	if len(streams[0]) != (int(cmds)+1)/2 {
		return nil, errors.New("packed op length")
	}
	if r.Len() != 0 {
		return nil, errors.New("packed trailing bytes")
	}
	// Reuse the extensively tested native-token decoder. The reconstructed
	// ephemeral payload is not persisted or additionally compressed.
	var old bytes.Buffer
	for _, v := range []int{int(cmds), len(streams[4]), len(streams[1]), len(streams[2]), len(streams[3]), len(streams[4])} {
		putVar(&old, uint64(v))
	}
	old.WriteByte(0)
	for _, idx := range []int{0, 1, 2, 3, 4} {
		old.Write(streams[idx])
	}
	decoded, err := decodeMatchNativeVariantHistory(bytes.NewReader(old.Bytes()), n, matched, history)
	return decoded, err
}

// Local selection always uses complete final archive sizes.
func bestPackedMatch(input []byte, digest [32]byte, tokens []byte) ([]byte, string, error) {
	a, err := encodeMatchPackedVariant(input, digest, tokens, false)
	if err != nil {
		return nil, "", err
	}
	b, err := encodeMatchPackedVariant(input, digest, tokens, true)
	if err != nil {
		return nil, "", err
	}
	if len(b) < len(a) {
		return b, "match-packed-pred", nil
	}
	return a, "match-packed", nil
}
