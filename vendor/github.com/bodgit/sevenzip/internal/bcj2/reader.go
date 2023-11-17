package bcj2

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/bodgit/sevenzip/internal/util"
	"github.com/hashicorp/go-multierror"
)

const (
	numMoveBits               = 5
	numbitModelTotalBits      = 11
	bitModelTotal        uint = 1 << numbitModelTotalBits
	numTopBits                = 24
	topValue             uint = 1 << numTopBits
)

func isJcc(b0, b1 byte) bool {
	return b0 == 0x0f && (b1&0xf0) == 0x80
}

func isJ(b0, b1 byte) bool {
	return (b1&0xfe) == 0xe8 || isJcc(b0, b1)
}

func index(b0, b1 byte) int {
	switch b1 {
	case 0xe8:
		return int(b0)
	case 0xe9:
		return 256
	default:
		return 257
	}
}

type readCloser struct {
	main util.ReadCloser
	call io.ReadCloser
	jump io.ReadCloser

	rd     util.ReadCloser
	nrange uint
	code   uint

	sd [256 + 2]uint

	previous byte
	written  uint64

	buf *bytes.Buffer
}

// NewReader returns a new BCJ2 io.ReadCloser.
func NewReader(_ []byte, _ uint64, readers []io.ReadCloser) (io.ReadCloser, error) {
	if len(readers) != 4 {
		return nil, errors.New("bcj2: need exactly four readers")
	}

	rc := &readCloser{
		main:   util.ByteReadCloser(readers[0]),
		call:   readers[1],
		jump:   readers[2],
		rd:     util.ByteReadCloser(readers[3]),
		nrange: 0xffffffff,
		buf:    new(bytes.Buffer),
	}
	rc.buf.Grow(1 << 16)

	b := make([]byte, 5)
	if _, err := io.ReadFull(rc.rd, b); err != nil {
		return nil, err
	}

	for _, x := range b {
		rc.code = (rc.code << 8) | uint(x)
	}

	for i := range rc.sd {
		rc.sd[i] = bitModelTotal >> 1
	}

	return rc, nil
}

func (rc *readCloser) Close() error {
	var err *multierror.Error
	if rc.main != nil {
		err = multierror.Append(err, rc.main.Close(), rc.call.Close(), rc.jump.Close(), rc.rd.Close())
	}

	return err.ErrorOrNil()
}

func (rc *readCloser) Read(p []byte) (int, error) {
	if rc.main == nil {
		return 0, errors.New("bcj2: Read after Close")
	}

	if err := rc.read(); err != nil && !errors.Is(err, io.EOF) {
		return 0, err
	}

	return rc.buf.Read(p)
}

func (rc *readCloser) update() error {
	if rc.nrange < topValue {
		b, err := rc.rd.ReadByte()
		if err != nil {
			return err
		}

		rc.code = (rc.code << 8) | uint(b)
		rc.nrange <<= 8
	}

	return nil
}

func (rc *readCloser) decode(i int) (bool, error) {
	newBound := (rc.nrange >> numbitModelTotalBits) * rc.sd[i]

	if rc.code < newBound {
		rc.nrange = newBound
		rc.sd[i] += (bitModelTotal - rc.sd[i]) >> numMoveBits

		if err := rc.update(); err != nil {
			return false, err
		}

		return false, nil
	}

	rc.nrange -= newBound
	rc.code -= newBound
	rc.sd[i] -= rc.sd[i] >> numMoveBits

	if err := rc.update(); err != nil {
		return false, err
	}

	return true, nil
}

func (rc *readCloser) read() error {
	var (
		b   byte
		err error
	)

	for {
		if b, err = rc.main.ReadByte(); err != nil {
			return err
		}

		rc.written++
		_ = rc.buf.WriteByte(b)

		if isJ(rc.previous, b) {
			break
		}

		rc.previous = b

		if rc.buf.Len() == rc.buf.Cap() {
			return nil
		}
	}

	bit, err := rc.decode(index(rc.previous, b))
	if err != nil {
		return err
	}

	if bit {
		var r io.Reader
		if b == 0xe8 {
			r = rc.call
		} else {
			r = rc.jump
		}

		var dest uint32
		if err = binary.Read(r, binary.BigEndian, &dest); err != nil {
			return err
		}

		dest -= uint32(rc.written + 4)
		_ = binary.Write(rc.buf, binary.LittleEndian, dest)

		rc.previous = byte(dest >> 24)
		rc.written += 4
	} else {
		rc.previous = b
	}

	return nil
}
