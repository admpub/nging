package utils

import (
	"bytes"
	"io"
)

type MatchWriter struct {
	out      io.Writer
	excepted []byte
	buf      bytes.Buffer
	cb       func()
	matched  bool
}

func (w *MatchWriter) match(p []byte) {
	if len(p) > len(w.excepted) {
		if bytes.Contains(p, w.excepted) {
			w.matched = true
			w.buf.Reset()
			w.cb()
			return
		}
		w.buf.Write(p[:len(w.excepted)-1])

		if bytes.Contains(w.buf.Bytes(), w.excepted) {
			w.matched = true
			w.buf.Reset()
			w.cb()
			return
		}
		w.buf.Reset()
		w.buf.Write(p[len(p)-len(w.excepted):])
		return
	}
	w.buf.Write(p)
	if w.buf.Len() <= len(w.excepted) {
		return
	}

	if bytes.Contains(w.buf.Bytes(), w.excepted) {
		w.matched = true
		w.buf.Reset()
		w.cb()
		return
	}

	reserved := w.buf.Bytes()[w.buf.Len()-len(w.excepted):]
	copy(w.buf.Bytes(), reserved)
	w.buf.Truncate(len(reserved))
}

func (w *MatchWriter) Write(p []byte) (c int, e error) {
	c, e = w.out.Write(p)
	if !w.matched {
		w.match(p)
	}
	return
}
