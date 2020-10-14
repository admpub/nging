package goloader

import "io"

type readAtSeeker struct {
	io.ReadSeeker
}

func (r *readAtSeeker) ReadAt(p []byte, offset int64) (n int, err error) {
	n = -1
	_, err = r.Seek(offset, io.SeekStart)
	if err == nil {
		n, err = r.Read(p)
	}
	return
}

func (r *readAtSeeker) ReadAtWithSize(p *[]byte, size, offset int64) (n int, err error) {
	n = -1
	_, err = r.Seek(offset, io.SeekStart)
	if err == nil {
		b := make([]byte, size)
		n, err = r.Read(b)
		*p = append(*p, b...)
	}
	return
}
