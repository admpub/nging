package notice

import "io"

func ToReadCloser(r io.Reader) io.ReadCloser {
	if rc, ok := r.(io.ReadCloser); ok {
		return rc
	}
	return io.NopCloser(r)
}

func ToWriteCloser(w io.Writer) io.WriteCloser {
	if wc, ok := w.(io.WriteCloser); ok {
		return wc
	}
	return nopWriteCloser{w}
}

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }

type proxyReader struct {
	io.ReadCloser
	prog *Progress
}

func (x *proxyReader) Read(p []byte) (int, error) {
	n, err := x.ReadCloser.Read(p)
	x.prog.Done(int64(n))
	if err == io.EOF {
		x.prog.SetComplete()
	}
	return n, err
}

type proxyWriter struct {
	io.WriteCloser
	prog *Progress
}

func (x *proxyWriter) Write(p []byte) (int, error) {
	n, err := x.WriteCloser.Write(p)
	x.prog.Done(int64(n))
	if err == io.EOF {
		x.prog.SetComplete()
	}
	return n, err
}

type proxyWriterTo struct {
	io.ReadCloser
	wt   io.WriterTo
	prog *Progress
}

func (x *proxyWriterTo) WriteTo(w io.Writer) (int64, error) {
	n, err := x.wt.WriteTo(w)
	x.prog.Done(n)
	if err == io.EOF {
		x.prog.SetComplete()
	}
	return n, err
}
