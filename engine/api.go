package engine

import "time"

// MaxInputSize is the CLI's existing safety bound, not a compression level.
const MaxInputSize = defaultLimit

// Report retains the encoder's method and accounting information.
type Report = summary

// Compress uses JLZ's single automatic lossless encoder.
func Compress(input []byte) ([]byte, Report, error) {
	return compress(input)
}

// Decompress restores an archive, checking its embedded SHA-256 digest.
func Decompress(archive []byte, limit int64) ([]byte, error) {
	return decompress(archive, limit)
}

// Print reports a completed compression without exposing encoder internals.
func (s summary) Print(path string, elapsed time.Duration) {
	s.elapsed = elapsed
	show(s, path)
}
