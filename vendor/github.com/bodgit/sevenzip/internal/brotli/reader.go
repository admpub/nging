package brotli

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"sync"

	"github.com/andybalholm/brotli"
	"github.com/bodgit/plumbing"
)

//nolint:gochecknoglobals
var brotliReaderPool sync.Pool

type readCloser struct {
	c io.Closer
	r *brotli.Reader
}

const (
	frameMagic  uint32 = 0x184d2a50
	frameSize   uint32 = 8
	brotliMagic uint16 = 0x5242 // 'B', 'R'
)

// This isn't part of the Brotli format but is prepended by the 7-zip implementation.
type headerFrame struct {
	FrameMagic       uint32
	FrameSize        uint32
	CompressedSize   uint32
	BrotliMagic      uint16
	UncompressedSize uint16 // * 64 KB
}

func (rc *readCloser) Close() (err error) {
	if rc.c != nil {
		brotliReaderPool.Put(rc.r)
		err = rc.c.Close()
		rc.c, rc.r = nil, nil
	}

	return
}

func (rc *readCloser) Read(p []byte) (int, error) {
	if rc.r == nil {
		return 0, errors.New("brotli: Read after Close")
	}

	return rc.r.Read(p)
}

// NewReader returns a new Brotli io.ReadCloser.
func NewReader(_ []byte, _ uint64, readers []io.ReadCloser) (io.ReadCloser, error) {
	if len(readers) != 1 {
		return nil, errors.New("brotli: need exactly one reader")
	}

	hr, b := new(headerFrame), new(bytes.Buffer)
	b.Grow(binary.Size(hr))

	// The 7-Zip Brotli compressor adds a 16 byte frame to the beginning of
	// the data which will confuse a pure Brotli implementation. Read it
	// but keep a copy so we can add it back if it doesn't look right
	if err := binary.Read(io.TeeReader(readers[0], b), binary.LittleEndian, hr); err != nil {
		return nil, err
	}

	var reader io.ReadCloser

	// If the header looks right, continue reading from that point
	// onwards, otherwise prepend it again and hope for the best
	if hr.FrameMagic == frameMagic && hr.FrameSize == frameSize && hr.BrotliMagic == brotliMagic {
		reader = readers[0]
	} else {
		reader = plumbing.MultiReadCloser(io.NopCloser(b), readers[0])
	}

	r, ok := brotliReaderPool.Get().(*brotli.Reader)
	if ok {
		_ = r.Reset(reader)
	} else {
		r = brotli.NewReader(reader)
	}

	return &readCloser{
		c: readers[0],
		r: r,
	}, nil
}
