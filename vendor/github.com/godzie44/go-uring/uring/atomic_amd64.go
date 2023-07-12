//go:build linux && amd64 && amd64_atomic

package uring

func SmpStoreReleaseUint32(p *uint32, v uint32) {
	*p = v
}

func SmpLoadAcquireUint32(p *uint32) uint32 {
	return *p
}

func ReadOnceUint32(p *uint32) uint32 {
	return *p
}
