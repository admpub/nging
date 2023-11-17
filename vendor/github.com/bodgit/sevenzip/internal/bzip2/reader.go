package bzip2

import (
	"compress/bzip2"
	"errors"
	"io"
)

type readCloser struct {
	c io.Closer
	r io.Reader
}

func (rc *readCloser) Close() error {
	var err error
	if rc.c != nil {
		err = rc.c.Close()
		rc.c, rc.r = nil, nil
	}

	return err
}

func (rc *readCloser) Read(p []byte) (int, error) {
	if rc.r == nil {
		return 0, errors.New("bzip2: Read after Close")
	}

	return rc.r.Read(p)
}

// NewReader returns a new bzip2 io.ReadCloser.
func NewReader(_ []byte, _ uint64, readers []io.ReadCloser) (io.ReadCloser, error) {
	if len(readers) != 1 {
		return nil, errors.New("bzip2: need exactly one reader")
	}

	return &readCloser{
		c: readers[0],
		r: bzip2.NewReader(readers[0]),
	}, nil
}
