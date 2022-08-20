// Copyright 2021 Hank Shen. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"io"
)

// WriteCloserTarget writes filtered log messages to a io.WriteCloser.
type WriteCloserTarget struct {
	io.WriteCloser
	*Filter
	close chan bool
}

// NewWriteCloserTarget creates a WriteCloserTarget.
func NewWriteCloserTarget(w io.WriteCloser) *WriteCloserTarget {
	return &WriteCloserTarget{
		WriteCloser: w,
		Filter:      &Filter{MaxLevel: LevelDebug},
		close:       make(chan bool),
	}
}

// Open nothing.
func (t *WriteCloserTarget) Open(errWriter io.Writer) (err error) {
	return nil
}

// Process writes a log message using Writer.
func (t *WriteCloserTarget) Process(e *Entry) {
	if e == nil {
		t.close <- true
		return
	}
	if !t.Allow(e) {
		return
	}
	t.Write([]byte(e.String() + "\n"))
}

// Close closes the file target.
func (t *WriteCloserTarget) Close() {
	<-t.close
	t.WriteCloser.Close()
}
