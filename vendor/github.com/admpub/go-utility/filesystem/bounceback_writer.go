package rifs

import (
	"io"
	"os"

	"github.com/dsoprea/go-logging"
)

// BouncebackWriter wraps a WriteSeeker, keeps track of our position, and
// seeks back to it before writing. This allows an underlying ReadWriteSeeker
// with an unstable position can still be used for a prolonged series of writes.
type BouncebackWriter struct {
	ws              io.WriteSeeker
	currentPosition int64

	statsWrites int
	statsSeeks  int
}

// NewBouncebackWriter returns a new `BouncebackWriter` struct.
func NewBouncebackWriter(ws io.WriteSeeker) (bw *BouncebackWriter, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	initialPosition, err := ws.Seek(0, os.SEEK_CUR)
	log.PanicIf(err)

	bw = &BouncebackWriter{
		ws:              ws,
		currentPosition: initialPosition,
	}

	return bw, nil
}

// Position returns the position that we're supposed to be at.
func (bw *BouncebackWriter) Position() int64 {
	return bw.currentPosition
}

// StatsWrites returns the number of write operations.
func (bw *BouncebackWriter) StatsWrites() int {
	return bw.statsWrites
}

// StatsSeeks returns the number of seek operations.
func (bw *BouncebackWriter) StatsSeeks() int {
	return bw.statsSeeks
}

// Seek puts us at a specific position in the internal writer for the next
// write/seek.
func (bw *BouncebackWriter) Seek(offset int64, whence int) (newPosition int64, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	newPosition, err = bw.ws.Seek(offset, whence)
	log.PanicIf(err)

	// Update our internal tracking.
	bw.currentPosition = newPosition

	return newPosition, nil
}

// Write performs a write against the internal `WriteSeeker` starting at the
// position that we're supposed to be at.
func (bw *BouncebackWriter) Write(p []byte) (n int, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	bw.statsWrites++

	// Make sure we're where we're supposed to be.

	realCurrentPosition, err := bw.ws.Seek(0, os.SEEK_CUR)
	log.PanicIf(err)

	if realCurrentPosition != bw.currentPosition {
		bw.statsSeeks++

		_, err = bw.ws.Seek(bw.currentPosition, os.SEEK_SET)
		log.PanicIf(err)
	}

	// Do write.

	n, err = bw.ws.Write(p)
	log.PanicIf(err)

	// Update our internal tracking.
	bw.currentPosition += int64(n)

	return n, nil
}
