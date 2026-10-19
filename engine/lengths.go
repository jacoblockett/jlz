package engine

import (
	"bytes"
	"errors"
)

// splitOperationLengths puts short copy lengths and literal-run lengths in
// different entropy models, while retaining exactly the current varint values.
func splitOperationLengths(ops []byte, commands int, lengths []byte) ([]byte, []byte, error) {
	r := bytes.NewReader(lengths)
	var literal, reference []byte
	if len(ops) != (commands+1)/2 {
		return nil, nil, errors.New("split op count")
	}
	for i := 0; i < commands; i++ {
		op := ops[i/2]
		if i%2 == 0 {
			op &= 15
		} else {
			op >>= 4
		}
		if op > matchRepAdd {
			return nil, nil, errors.New("split invalid opcode")
		}
		value, err := readVar(r)
		if err != nil {
			return nil, nil, err
		}
		if op == matchLiteral {
			literal = appendUvar(literal, value)
		} else {
			reference = appendUvar(reference, value)
		}
	}
	if r.Len() != 0 {
		return nil, nil, errors.New("split trailing length bytes")
	}
	return literal, reference, nil
}

func mergeOperationLengths(ops []byte, commands int, literal, reference []byte) ([]byte, error) {
	a, b := bytes.NewReader(literal), bytes.NewReader(reference)
	if len(ops) != (commands+1)/2 {
		return nil, errors.New("merge op count")
	}
	result := make([]byte, 0, len(literal)+len(reference))
	for i := 0; i < commands; i++ {
		op := ops[i/2]
		if i%2 == 0 {
			op &= 15
		} else {
			op >>= 4
		}
		if op > matchRepAdd {
			return nil, errors.New("merge invalid opcode")
		}
		rd := b
		if op == matchLiteral {
			rd = a
		}
		n, err := readVar(rd)
		if err != nil {
			return nil, err
		}
		result = appendUvar(result, n)
	}
	if a.Len() != 0 || b.Len() != 0 {
		return nil, errors.New("merge trailing length bytes")
	}
	return result, nil
}

func packSplitLengthHeader(st nativeStreams, packedOp, packedDistances, packedTransforms []byte) ([]byte, bool) {
	if st.ncommands < 128 {
		return nil, false
	}
	literal, copyLengths, err := splitOperationLengths(st.op, st.ncommands, st.lengths)
	if err != nil || len(literal) == 0 || len(copyLengths) == 0 {
		return nil, false
	}
	var b bytes.Buffer
	putVar(&b, uint64(st.ncommands))
	b.Write(packedOp)
	b.Write(packFiniteState(literal))
	b.Write(packFiniteState(copyLengths))
	b.Write(packedDistances)
	b.Write(packedTransforms)
	return b.Bytes(), true
}
