package reactor

import (
	"sync"
	"sync/atomic"
)

type (
	cbMap map[uint32]Callback
	shard struct {
		nonces    []uint32
		callbacks []cbMap
		boundary  int

		slowCallbacks map[int]cbMap
		slowNonces    map[int]uint32

		sync.Mutex
	}
	cbRegistry struct {
		shards      []*shard
		granularity int
		shardCnt    int
	}
)

func newShard(fCap int) *shard {
	buff := make([]cbMap, fCap)
	for i := range buff {
		buff[i] = make(cbMap, 10)
	}

	return &shard{
		callbacks:     buff,
		nonces:        make([]uint32, fCap),
		boundary:      fCap,
		slowCallbacks: make(map[int]cbMap),
		slowNonces:    make(map[int]uint32),
	}
}

func (sh *shard) add(idx int, cb Callback) (n uint32) {
	if idx < sh.boundary {
		n = atomic.AddUint32(&sh.nonces[idx], 1)
		sh.Lock()
		sh.callbacks[idx][n] = cb
		sh.Unlock()
		return n
	}

	//slow path, for big fd values
	sh.Lock()
	sh.slowNonces[idx]++
	n = sh.slowNonces[idx]

	if _, exists := sh.slowCallbacks[idx]; !exists {
		sh.slowCallbacks[idx] = make(cbMap, 10)
	}

	sh.slowCallbacks[idx][n] = cb
	sh.Unlock()
	return n
}

func (sh *shard) pop(idx int, nonce uint32) Callback {
	if idx < sh.boundary {
		sh.Lock()
		cb := sh.callbacks[idx][nonce]
		delete(sh.callbacks[idx], nonce)
		sh.Unlock()
		return cb
	}

	sh.Lock()
	cb := sh.slowCallbacks[idx][nonce]
	delete(sh.slowCallbacks[idx], nonce)
	sh.Unlock()

	return cb
}

func newCbRegistry(shardCount int, granularity int) *cbRegistry {
	shards := make([]*shard, shardCount)
	for i := 0; i < shardCount; i++ {
		shards[i] = newShard((1 << 16) / shardCount)
	}

	return &cbRegistry{
		granularity: granularity,
		shardCnt:    shardCount,
		shards:      shards,
	}
}

func (r *cbRegistry) shardNumAndFlattenIdx(fd int) (int, int) {
	// fd / r.granularity - granule number
	gNum := fd / r.granularity

	// gNum/r.shardCnt - granule number in shard
	// granule number in shard *r.granularity - index of first el in granule
	// index of granule start + fd % r.granularity - index of file descriptor in shard
	return gNum % r.shardCnt, (gNum/r.shardCnt)*r.granularity + (fd % r.granularity)
}

func (r *cbRegistry) add(fd int, cb Callback) uint32 {
	shardNum, idx := r.shardNumAndFlattenIdx(fd)
	return r.shards[shardNum].
		add(idx, cb)
}

func (r *cbRegistry) pop(fd int, nonce uint32) Callback {
	shardNum, idx := r.shardNumAndFlattenIdx(fd)
	return r.shards[shardNum].
		pop(idx, nonce)
}
