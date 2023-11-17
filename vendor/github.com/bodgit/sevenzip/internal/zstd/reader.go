package zstd

import (
	"errors"
	"io"
	"runtime"
	"sync"

	"github.com/klauspost/compress/zstd"
)

//nolint:gochecknoglobals
var zstdReaderPool sync.Pool

type readCloser struct {
	c io.Closer
	r *zstd.Decoder
}

func (rc *readCloser) Close() (err error) {
	if rc.c != nil {
		zstdReaderPool.Put(rc.r)
		err = rc.c.Close()
		rc.c, rc.r = nil, nil
	}

	return
}

func (rc *readCloser) Read(p []byte) (int, error) {
	if rc.r == nil {
		return 0, errors.New("zstd: Read after Close")
	}

	return rc.r.Read(p)
}

// NewReader returns a new Zstandard io.ReadCloser.
func NewReader(_ []byte, _ uint64, readers []io.ReadCloser) (io.ReadCloser, error) {
	if len(readers) != 1 {
		return nil, errors.New("zstd: need exactly one reader")
	}

	var err error

	r, ok := zstdReaderPool.Get().(*zstd.Decoder)
	if ok {
		if err = r.Reset(readers[0]); err != nil {
			return nil, err
		}
	} else {
		if r, err = zstd.NewReader(readers[0]); err != nil {
			return nil, err
		}
		runtime.SetFinalizer(r, (*zstd.Decoder).Close)
	}

	return &readCloser{
		c: readers[0],
		r: r,
	}, nil
}
