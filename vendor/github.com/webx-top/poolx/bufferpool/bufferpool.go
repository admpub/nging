package bufferpool

import (
	"bytes"
	"sync"
)

// Package bufferpool is a simple wrapper around sync.Pool that is specific
// to bytes.Buffer.

var global = New()

// Get returns a bytes.Buffer from a global pool
func Get() *bytes.Buffer {
	return global.Get()
}

// Release puts the given bytes.Buffer instance back in the global pool.
func Release(buf *bytes.Buffer) {
	global.Release(buf)
}

// BufferPool is a sync.Pool for bytes.Buffer objects
type BufferPool struct {
	pool sync.Pool
}

// New creates a new BufferPool instance
func New() *BufferPool {
	var bp BufferPool
	bp.pool.New = allocBuffer
	return &bp
}

func allocBuffer() interface{} {
	return &bytes.Buffer{}
}

// Get returns a bytes.Buffer from the specified pool
func (bp *BufferPool) Get() *bytes.Buffer {
	return bp.pool.Get().(*bytes.Buffer)
}

// Release puts the given bytes.Buffer back in the specified pool after
// resetting it
func (bp *BufferPool) Release(buf *bytes.Buffer) {
	buf.Reset()
	bp.pool.Put(buf)
}
