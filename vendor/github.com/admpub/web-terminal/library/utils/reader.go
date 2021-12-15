package utils

import "io"

type ConsoleReader struct {
	dst io.ReadCloser
	out io.Writer
}

func (w *ConsoleReader) Read(p []byte) (n int, err error) {
	n, err = w.dst.Read(p)
	if n > 0 {
		w.out.Write(p[:n])
	}
	return
}

func (w *ConsoleReader) Close() error {
	return w.dst.Close()
}
