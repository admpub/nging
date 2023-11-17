package delta

import (
	"errors"
	"io"
)

const (
	stateSize = 256
)

type readCloser struct {
	rc    io.ReadCloser
	state [stateSize]byte
	delta int
}

func (rc *readCloser) Close() (err error) {
	if rc.rc != nil {
		err = rc.rc.Close()
		rc.rc = nil
	}

	return
}

func (rc *readCloser) Read(p []byte) (int, error) {
	if rc.rc == nil {
		return 0, errors.New("delta: Read after Close")
	}

	n, err := rc.rc.Read(p)
	if err != nil {
		return n, err
	}

	var (
		buffer [stateSize]byte
		j      int
	)

	copy(buffer[:], rc.state[:rc.delta])

	for i := 0; i < n; {
		for j = 0; j < rc.delta && i < n; i++ {
			p[i] = buffer[j] + p[i]
			buffer[j] = p[i]
			j++
		}
	}

	if j == rc.delta {
		j = 0
	}

	copy(rc.state[:], buffer[j:rc.delta])
	copy(rc.state[rc.delta-j:], buffer[:j])

	return n, nil
}

// NewReader returns a new Delta io.ReadCloser.
func NewReader(p []byte, _ uint64, readers []io.ReadCloser) (io.ReadCloser, error) {
	if len(readers) != 1 {
		return nil, errors.New("delta: need exactly one reader")
	}

	if len(p) != 1 {
		return nil, errors.New("delta: not enough properties")
	}

	return &readCloser{
		rc:    readers[0],
		delta: int(p[0] + 1),
	}, nil
}
