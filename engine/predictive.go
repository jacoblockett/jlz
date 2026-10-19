package engine

// MDC5: reversible byte-predictor transforms followed by an adaptive,
// order-0/order-1 context arithmetic coder. The coder is standalone and
// has no third-party dependencies. Only an actual size win is retained.

import (
 "github.com/jacoblockett/jlz/predictive"
 "github.com/jacoblockett/jlz/entropy"
	"bytes"
	"errors"
	"fmt"
	"io"
)

const modePredictive byte = 7
const (
	predNone        byte = 0
	predXor         byte = 1
	predSub         byte = 2
	codecNative     byte = 0
	codecArithmetic byte = 1
)

// Arithmetic streams use the shared entropy package.
func arithmeticEncode(data []byte) []byte { return entropy.EncodeArithmetic(data) }
func arithmeticDecode(payload []byte, length int) ([]byte,error) { return entropy.DecodeArithmetic(payload, length) }

func transform(input []byte, typ byte, lag int) []byte { return predictive.Transform(input,typ,lag) }
func reverseTransform(input []byte, typ byte, lag int) []byte { return predictive.Reverse(input,typ,lag) }

func nativeBytes(b []byte) ([]byte, error)          { return nativePackBytes(b) }
func unnativeBytes(b []byte, n int) ([]byte, error) { return nativeUnpackBytes(b, n) }
func encodePredictive(input []byte, digest [32]byte) ([]byte, string, error) {
	bestSize := int(^uint(0) >> 1)
	var best []byte
	bestName := ""
	types := []struct {
		typeID byte
		lag    int
	}{
		{predNone, 0}, {predXor, 1}, {predSub, 1},
		{predXor, 2}, {predSub, 2}, {predXor, 4}, {predSub, 4},
		{predXor, 8}, {predSub, 8},
	}
	for _, t := range types {
		if t.lag >= len(input) {
			continue
		}
		data := transform(input, t.typeID, t.lag)
		// Test both the simple entropy baseline and a context arithmetic coder.
		for _, codec := range []byte{codecArithmetic, codecNative} {
			var encoded []byte
			var e error
			if codec == codecArithmetic {
				encoded = arithmeticEncode(data)
			} else {
				encoded, e = nativeBytes(data)
				if e != nil {
					return nil, "", e
				}
			}
			var b bytes.Buffer
			b.WriteString(magic)
			b.WriteByte(modePredictive)
			putVar(&b, uint64(len(input)))
			b.Write(digest[:])
			b.WriteByte(t.typeID)
			putVar(&b, uint64(t.lag))
			b.WriteByte(codec)
			b.Write(encoded)
			if b.Len() < bestSize {
				best = b.Bytes()
				bestSize = b.Len()
				bestName = fmt.Sprintf("transform=%d lag=%d codec=%d", t.typeID, t.lag, codec)
			}
		}
	}
	return best, bestName, nil
}
func decodePredictive(r *bytes.Reader, n uint64) ([]byte, error) {
	typ, e := r.ReadByte()
	if e != nil {
		return nil, e
	}
	lag, e := readVar(r)
	if e != nil {
		return nil, e
	}
	codec, e := r.ReadByte()
	if e != nil {
		return nil, e
	}
	if typ > predSub || (typ != predNone && (lag == 0 || lag >= n)) || (typ == predNone && lag != 0) {
		return nil, errors.New("bad predictive transform")
	}
	data := make([]byte, r.Len())
	if _, e = io.ReadFull(r, data); e != nil {
		return nil, e
	}
	var decoded []byte
	switch codec {
	case codecArithmetic:
		decoded, e = arithmeticDecode(data, int(n))
	case codecNative:
		decoded, e = unnativeBytes(data, int(n))
	default:
		return nil, errors.New("unknown predictive codec")
	}
	if e != nil {
		return nil, e
	}
	return reverseTransform(decoded, typ, int(lag)), nil
}
