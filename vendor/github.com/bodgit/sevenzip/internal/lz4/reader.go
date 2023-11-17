package lz4

import (
	"errors"
	"io"
	"sync"

	lz4 "github.com/pierrec/lz4/v4"
)

//nolint:gochecknoglobals
var lz4ReaderPool sync.Pool

type readCloser struct {
	c io.Closer
	r *lz4.Reader
}

func (rc *readCloser) Close() (err error) {
	if rc.c != nil {
		lz4ReaderPool.Put(rc.r)
		err = rc.c.Close()
		rc.c, rc.r = nil, nil
	}

	return
}

func (rc *readCloser) Read(p []byte) (int, error) {
	if rc.r == nil {
		return 0, errors.New("lz4: Read after Close")
	}

	return rc.r.Read(p)
}

// NewReader returns a new LZ4 io.ReadCloser.
func NewReader(_ []byte, _ uint64, readers []io.ReadCloser) (io.ReadCloser, error) {
	if len(readers) != 1 {
		return nil, errors.New("lz4: need exactly one reader")
	}

	r, ok := lz4ReaderPool.Get().(*lz4.Reader)
	if ok {
		r.Reset(readers[0])
	} else {
		r = lz4.NewReader(readers[0])
	}

	return &readCloser{
		c: readers[0],
		r: r,
	}, nil
}
