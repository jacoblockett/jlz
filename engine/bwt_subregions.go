package engine

import (
	"bytes"
	"errors"
	"github.com/jacoblockett/jlz/bwt"
	"math"
)

const blockBWTSubregions byte = 16

// regimeCost is a zero-order minimum-description-length estimate for a native
// independently framed entropy segment; no source-specific trained tables.
// The predicted table costs conservatively 3 bytes per distinct symbol.
func regimeCost(symbols []byte) float64 {
	if len(symbols) == 0 {
		return 1e100
	}
	var counts [256]int
	for _, b := range symbols {
		counts[b]++
	}
	total := float64(len(symbols))
	distinct := 0
	entropyBits := 0.0
	for _, n := range counts {
		if n == 0 {
			continue
		}
		distinct++
		entropyBits += float64(n) * math.Log2(total/float64(n))
	}
	return entropyBits/8 + float64(distinct*3+8)
}

func encodeBWTSubregions(input []byte) ([]byte, byte, error) {
	transformed, primary, err := bwtNativeTransform(input)
	if err != nil {
		return nil, 0, err
	}
	syms := bwt.BijectiveZeroEncode(transformed)
	var old bytes.Buffer
	putVar(&old, uint64(primary))
	old.Write(packFiniteState(syms))
	best := old.Bytes()
	if len(syms) < 2048 {
		return best, blockBWTBijective, nil
	}
	// Estimate the total cost of a single model, two neighborhood models,
	// and four neighborhood models before considering any further compression.
	estimateBase := regimeCost(syms)
	choice := 0
	predicted := estimateBase
	for _, parts := range []int{2, 4} {
		cost := float64(1) // partition count
		for j := 0; j < parts; j++ {
			a := len(syms) * j / parts
			z := len(syms) * (j + 1) / parts
			cost += regimeCost(syms[a:z])
		}
		if cost+128 < predicted {
			choice = parts
			predicted = cost
		}
	}
	if choice == 0 {
		return best, blockBWTBijective, nil
	}
	var candidate bytes.Buffer
	putVar(&candidate, uint64(primary))
	candidate.WriteByte(byte(choice))
	for j := 0; j < choice; j++ {
		a := len(syms) * j / choice
		z := len(syms) * (j + 1) / choice
		candidate.Write(packFiniteState(syms[a:z]))
		if candidate.Len() >= len(best) {
			return best, blockBWTBijective, nil
		}
	}
	if candidate.Len() < len(best) {
		return candidate.Bytes(), blockBWTSubregions, nil
	}
	return best, blockBWTBijective, nil
}

func decodeBWTSubregions(payload []byte, expected int) ([]byte, error) {
	r := bytes.NewReader(payload)
	primary, err := readVar(r)
	if err != nil {
		return nil, err
	}
	if expected <= 0 || primary == 0 || primary > uint64(expected) {
		return nil, errors.New("segmented BWT primary")
	}
	count, e := r.ReadByte()
	if e != nil {
		return nil, e
	}
	if count != 2 && count != 4 {
		return nil, errors.New("segmented BWT count")
	}
	maxSymbols := expected*2 + 20
	syms := make([]byte, 0, expected)
	for j := 0; j < int(count); j++ {
		part, e := unpackStream(r, maxSymbols-len(syms))
		if e != nil {
			return nil, e
		}
		if len(part) == 0 || len(part) > maxSymbols-len(syms) {
			return nil, errors.New("bad segment length")
		}
		syms = append(syms, part...)
	}
	if r.Len() != 0 {
		return nil, errors.New("segmented BWT trailing")
	}
	transformed, e := bwt.BijectiveZeroDecode(syms, expected)
	if e != nil {
		return nil, e
	}
	return inverseBWTNative(transformed, int(primary))
}
