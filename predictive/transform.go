package predictive

// Reversible XOR/subtraction against a fixed previous-byte lag.
const (
	predNone byte = 0
	predXor  byte = 1
)
func Transform(input []byte, typ byte, lag int) []byte {
	out := make([]byte, len(input))
	if typ == predNone {
		copy(out, input)
		return out
	}
	copy(out[:lag], input[:lag])
	for i := lag; i < len(input); i++ {
		if typ == predXor {
			out[i] = input[i] ^ input[i-lag]
		} else {
			out[i] = input[i] - input[i-lag]
		}
	}
	return out
}
func Reverse(input []byte, typ byte, lag int) []byte {
	out := append([]byte(nil), input...)
	if typ == predNone {
		return out
	}
	for i := lag; i < len(input); i++ {
		if typ == predXor {
			out[i] ^= out[i-lag]
		} else {
			out[i] += out[i-lag]
		}
	}
	return out
}
