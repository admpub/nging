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

func newProxyReader(r io.Reader, prog *Progress) io.ReadCloser {
	rc := toReadCloser(r)
	pr := proxyReader{rc, prog}
	if _, ok := r.(io.WriterTo); ok {
		return proxyWriterTo{pr}
	}
	return pr
}

type proxyReader struct {
	io.ReadCloser
	prog *Progress
}

func (x proxyReader) Read(p []byte) (int, error) {
	n, err := x.ReadCloser.Read(p)
	x.prog.Done(int64(n))
	return n, err
}

func newProxyWriter(w io.Writer, prog *Progress) io.WriteCloser {
	wc := toWriteCloser(w)
	pw := proxyWriter{wc, prog *Progress}
	if _, ok := w.(io.ReaderFrom); ok {
		return proxyReaderFrom{pw}
	}
	return pw
}

type proxyWriter struct {
	io.WriteCloser
	prog *Progress
}

func (x proxyWriter) Write(p []byte) (int, error) {
	n, err := x.WriteCloser.Write(p)
	x.prog.Done(int64(n))
	return n, err
}

type proxyWriterTo struct {
	proxyReader
}

func (x proxyWriterTo) WriteTo(w io.Writer) (int64, error) {
	n, err := x.ReadCloser.(io.WriterTo).WriteTo(w)
	x.prog.Done(int64(n))
	return n, err
}


type proxyReaderFrom struct {
	proxyWriter
}

func (x proxyReaderFrom) ReadFrom(r io.Reader) (int64, error) {
	n, err := x.WriteCloser.(io.ReaderFrom).ReadFrom(r)
	x.prog.Done(int64(n))
	return n, err
}

