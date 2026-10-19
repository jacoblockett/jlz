package engine

import "github.com/jacoblockett/jlz/match"

// Keep the existing verified token decoders; only the byte-copy primitive
// is factored into the match package.
func expandMatch(out []byte, at, distance, length int, kind, mask byte) {
 match.Expand(out, at, distance, length, kind, mask)
}
