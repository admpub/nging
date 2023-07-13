//go:build linux && !amd64_atomic

package uring

import (
	"sync/atomic"
)

func SmpStoreReleaseUint32(p *uint32, v uint32) {
	atomic.StoreUint32(p, v)
}

func SmpLoadAcquireUint32(p *uint32) uint32 {
	return atomic.LoadUint32(p)
}

func ReadOnceUint32(p *uint32) uint32 {
	return atomic.LoadUint32(p)
}
