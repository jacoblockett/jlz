package main

// Native-only framed byte stream. This is our project's own raw/RLE/canonical
// Huffman byte packer (packStream and unpackStream in native_packed.go), NOT
// DEFLATE, Zstd, XZ or a library compression implementation.
// The stream uses its own length framing, so callers can read safely to EOF.

import (
	"bytes"
	"errors"
	"io"
)

// A deliberately conservative bound on intermediate decoded token streams.
// Ordinary full files are subject to the fixed safety limit.
const nativeIntermediateLimit = 384 * 1024 * 1024

type nativeStreamWriter struct {
	target io.Writer
	buf    bytes.Buffer
	closed bool
}

func newNativeWriter(dst io.Writer) (*nativeStreamWriter, error) {
	return &nativeStreamWriter{target: dst}, nil
}

func (w *nativeStreamWriter) Write(p []byte) (int, error) {
	if w.closed {
		return 0, errors.New("writer already closed")
	}
	return w.buf.Write(p)
}
func (w *nativeStreamWriter) Close() error {
	if w.closed {
		return errors.New("writer already closed")
	}
	w.closed = true
	_, err := w.target.Write(packStream(w.buf.Bytes()))
	return err
}

type nativeStreamReader struct {
	*bytes.Reader
	err error
}

func (r *nativeStreamReader) Read(p []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	return r.Reader.Read(p)
}
func (r *nativeStreamReader) Close() error { return nil }
func newNativeReader(src io.Reader) io.ReadCloser {
	b, e := io.ReadAll(io.LimitReader(src, nativeIntermediateLimit+1024*1024+1))
	if e != nil {
		return &nativeStreamReader{Reader: bytes.NewReader(nil), err: e}
	}
	if len(b) > nativeIntermediateLimit+1024*1024 {
		return &nativeStreamReader{Reader: bytes.NewReader(nil), err: errors.New("native stream too large")}
	}
	rd := bytes.NewReader(b)
	out, e := unpackStream(rd, nativeIntermediateLimit)
	if e == nil && rd.Len() != 0 {
		e = errors.New("native stream trailing bytes")
	}
	return &nativeStreamReader{Reader: bytes.NewReader(out), err: e}
}
func nativePackBytes(data []byte) ([]byte, error) { return packStream(data), nil }
func nativeUnpackBytes(src []byte, limit int) ([]byte, error) {
	r := bytes.NewReader(src)
	out, err := unpackStream(r, limit)
	if err == nil && r.Len() != 0 {
		err = errors.New("native trailing bytes")
	}
	return out, err
}
