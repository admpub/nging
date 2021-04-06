package writer

import (
	"bytes"
	"io"
	"os"
)

type Shadow struct {
	buf    *bytes.Buffer
	closed bool
}

type errWriter struct {
	*Shadow
	w io.Writer
}

type outWriter struct {
	*Shadow
	w io.Writer
}

func NewShadow() *Shadow {
	return &Shadow{
		buf: bytes.NewBuffer(nil),
	}
}

func NewOut(s *Shadow) io.Writer {
	return &outWriter{Shadow: s, w: os.Stdout}
}

func (w *outWriter) Write(p []byte) (n int, err error) {
	if !w.closed && w.buf != nil {
		w.buf.Write(p[:])
	}
	return w.w.Write(p)
}

func (w *errWriter) Write(p []byte) (n int, err error) {
	if !w.closed && w.buf != nil {
		w.buf.Write(p[:])
	}
	return w.w.Write(p)
}

func (w *Shadow) String() string {
	r := w.buf.String()
	w.closed = true
	return r
}

func NewErr(s *Shadow) io.Writer {
	return &errWriter{Shadow: s, w: os.Stderr}
}
