package deflate

import (
	"errors"
	"io"
	"sync"

	"github.com/bodgit/sevenzip/internal/util"
	"github.com/klauspost/compress/flate"
)

//nolint:gochecknoglobals
var flateReaderPool sync.Pool

type readCloser struct {
	c  io.Closer
	fr io.ReadCloser
}

func (rc *readCloser) Close() error {
	var err error

	if rc.c != nil {
		if err = rc.fr.Close(); err != nil {
			return err
		}

		flateReaderPool.Put(rc.fr)
		err = rc.c.Close()
		rc.c, rc.fr = nil, nil
	}

	return err
}

func (rc *readCloser) Read(p []byte) (int, error) {
	if rc.fr == nil {
		return 0, errors.New("deflate: Read after Close")
	}

	return rc.fr.Read(p)
}

// NewReader returns a new DEFLATE io.ReadCloser.
func NewReader(_ []byte, _ uint64, readers []io.ReadCloser) (io.ReadCloser, error) {
	if len(readers) != 1 {
		return nil, errors.New("deflate: need exactly one reader")
	}

	fr, ok := flateReaderPool.Get().(io.ReadCloser)
	if ok {
		frf, ok := fr.(flate.Resetter)
		if ok {
			if err := frf.Reset(util.ByteReadCloser(readers[0]), nil); err != nil {
				return nil, err
			}
		}
	} else {
		fr = flate.NewReader(util.ByteReadCloser(readers[0]))
	}

	return &readCloser{
		c:  readers[0],
		fr: fr,
	}, nil
}
