package handler

import "io"

type consoleReader struct {
	dst io.ReadCloser
	out io.Writer
}

func (w *consoleReader) Read(p []byte) (n int, err error) {
	n, err = w.dst.Read(p)
	if n > 0 {
		w.out.Write(p[:n])
	}
	return
}

func (w *consoleReader) Close() error {
	return w.dst.Close()
}
