package engine

// Logarithmic distance slots, with the low bits stored verbatim. This is
// native JLZ entropy coding, not a separate compressor or optional mode.
import (
	"bytes"
	"errors"
	"math/bits"
)

// packVarintStream retains the ordinary representation unless the complete
// slot-based representation, including framing, is smaller.
func packVarintStream(data []byte) []byte {
	plain := packStream(data)
	if len(data) < 128 {
		return plain
	}
	r := bytes.NewReader(data)
	slots := make([]byte, 0, len(data))
	extras := make([]byte, 0, len(data))
	var pending uint64
	pendingBits := 0
	for r.Len() > 0 {
		v, e := readVar(r)
		if e != nil || v == 0 || v > 0xffffffff {
			return plain
		}
		width := bits.Len64(v)
		slots = append(slots, byte(width))
		k := width - 1
		pending = (pending << uint(k)) | (v & ((uint64(1) << uint(k)) - 1))
		pendingBits += k
		for pendingBits >= 8 {
			pendingBits -= 8
			extras = append(extras, byte(pending>>uint(pendingBits)))
			pending &= (uint64(1) << uint(pendingBits)) - 1
		}
	}
	// Bitstream is exact-length, with at most seven zero padding bits.
	packedSlots := packStream(slots)
	if pendingBits > 0 {
		extras = append(extras, byte(pending<<uint(8-pendingBits)))
	}
	extra := extras
	payloadLen := len(packedSlots) + len(extra)
	// Framing cost of alternative representation.
	total := varSize(uint64(len(data))) + varSize(uint64(payloadLen)) + 1 + payloadLen
	if total >= len(plain) {
		return plain
	}
	out := make([]byte, 0, total)
	out = appendUvar(out, uint64(len(data)))
	out = appendUvar(out, uint64(payloadLen))
	out = append(out, 4)
	out = append(out, packedSlots...)
	out = append(out, extra...)
	return out
}

func unpackVarintSlots(payload []byte, orig int) ([]byte, error) {
	rd := bytes.NewReader(payload)
	// The nested slot array is encoded only in the preexisting 0/1/2 formats.
	start := len(payload)
	n, e := readVar(rd)
	if e != nil {
		return nil, e
	}
	sz, e := readVar(rd)
	if e != nil {
		return nil, e
	}
	codec, e := rd.ReadByte()
	if e != nil {
		return nil, e
	}
	if codec > 2 || n > uint64(orig) || sz > uint64(rd.Len()) {
		return nil, errors.New("invalid nested distance slot stream")
	}
	// Bound the nested decoder to the declared number of slot symbols.
	rd = bytes.NewReader(payload)
	slots, e := unpackStream(rd, orig)
	if e != nil {
		return nil, e
	}
	packedLen := start - rd.Len()
	extra := payload[packedLen:]
	needed := 0
	for _, width := range slots {
		if width == 0 || width > 32 {
			return nil, errors.New("invalid distance slot")
		}
		needed += int(width) - 1
	}
	if (needed+7)/8 != len(extra) {
		return nil, errors.New("distance low-bit length mismatch")
	}
	out := make([]byte, 0, orig)
	var pending uint64
	pendingBits := 0
	at := 0
	for _, width := range slots {
		k := int(width) - 1
		for pendingBits < k {
			pending = (pending << 8) | uint64(extra[at])
			at++
			pendingBits += 8
		}
		pendingBits -= k
		low := (pending >> uint(pendingBits)) & ((uint64(1) << uint(k)) - 1)
		pending &= (uint64(1) << uint(pendingBits)) - 1
		value := (uint64(1) << uint(k)) | low
		out = appendUvar(out, value)
		if len(out) > orig {
			return nil, errors.New("distance slots exceed raw length")
		}
	}
	if len(out) != orig {
		return nil, errors.New("distance slots raw length mismatch")
	}
	if needed%8 != 0 && len(extra) > 0 && extra[len(extra)-1]&byte((1<<(8-needed%8))-1) != 0 {
		return nil, errors.New("distance slots nonzero padding")
	}
	return out, nil
}
