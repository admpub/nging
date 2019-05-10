package rest

import (
	"container/list"
	"sync"
	"time"
)

// ResourceCache, is an LRU-TTL Cache, that caches Responses base on headers
// It uses 3 goroutines -> one for LRU, and the other two for TTL.

// The cache itself.
var resourceCache *resourceTtlLruMap

// ByteSize is a helper for configuring MaxCacheSize
type ByteSize int64

const (
	_ = iota

	// KB = KiloBytes
	KB ByteSize = 1 << (10 * iota)

	// MB = MegaBytes
	MB

	// GB = GigaBytes
	GB
)

// MaxCacheSize is the Maxium Byte Size to be hold by the ResourceCache
// Default is 1 GigaByte
// Type: rest.ByteSize
var MaxCacheSize = 1 * GB

// Current Cache Size.
var cacheSize int64

type lruOperation int

const (
	move lruOperation = iota
	push
	del
	last
)

type lruMsg struct {
	operation lruOperation
	resp      *Response
}

type resourceTtlLruMap struct {
	cache    map[string]*Response
	skipList *skipList    // skiplist for TTL
	lruList  *list.List   // List for LRU
	lruChan  chan *lruMsg // Channel for LRU messages
	ttlChan  chan bool    // Channel for TTL messages
	popChan  chan string
	rwMutex  sync.RWMutex //Read Write Locking Mutex
}

func init() {

	resourceCache = &resourceTtlLruMap{
		cache:    make(map[string]*Response),
		skipList: newSkipList(),
		lruList:  list.New(),
		lruChan:  make(chan *lruMsg, 10000),
		ttlChan:  make(chan bool, 1000),
		popChan:  make(chan string),
		rwMutex:  sync.RWMutex{},
	}

	go resourceCache.lruOperations()
	go resourceCache.ttl()

}

func (rCache *resourceTtlLruMap) lruOperations() {

	for {
		msg := <-rCache.lruChan

		switch msg.operation {
		case move:
			rCache.lruList.MoveToFront(msg.resp.listElement)
		case push:
			msg.resp.listElement = rCache.lruList.PushFront(msg.resp.Request.URL.String())
		case del:
			rCache.lruList.Remove(msg.resp.listElement)
		case last:
			rCache.popChan <- rCache.lruList.Back().Value.(string)
		}

	}

}

func (rCache *resourceTtlLruMap) get(key string) *Response {

	//Read lock only
	rCache.rwMutex.RLock()
	resp := rCache.cache[key]
	rCache.rwMutex.RUnlock()

	//If expired, remove it
	if resp != nil && resp.ttl != nil && resp.ttl.Sub(time.Now()) <= 0 {

		//Full lock
		rCache.rwMutex.Lock()
		defer rCache.rwMutex.Unlock()

		//JIC, get the freshest version
		resp = rCache.cache[key]

		//Check again with the lock
		if resp != nil && resp.ttl != nil && resp.ttl.Sub(time.Now()) <= 0 {
			rCache.remove(key, resp)
			return nil //return. Do not send the move message
		}

	}

	if resp != nil {

		//Buffered msg to LruList
		//Move forward
		rCache.lruChan <- &lruMsg{
			operation: move,
			resp:      resp,
		}

	}

	return resp
}

// Set if key not exist
func (rCache *resourceTtlLruMap) setNX(key string, value *Response) {

	//Full Lock
	rCache.rwMutex.Lock()
	defer rCache.rwMutex.Unlock()

	v := rCache.cache[key]

	if v == nil {

		rCache.cache[key] = value

		//PushFront in LruList
		rCache.lruChan <- &lruMsg{
			operation: push,
			resp:      value,
		}

		//Set ttl if necessary
		if value.ttl != nil {
			value.skipListElement = rCache.skipList.insert(key, *value.ttl)
			rCache.ttlChan <- true
		}

		// Add Response Size to Cache
		// Not necessary to use atomic
		cacheSize += value.size()

		for i := 0; ByteSize(cacheSize) >= MaxCacheSize && i < 10; i++ {

			rCache.lruChan <- &lruMsg{
				last,
				nil,
			}

			k := <-rCache.popChan
			r := rCache.cache[k]

			rCache.remove(k, r)

		}

	}

}

//
func (rCache *resourceTtlLruMap) remove(key string, resp *Response) {

	delete(rCache.cache, key)                    //Delete from map
	rCache.skipList.remove(resp.skipListElement) //Delete from skipList
	rCache.lruChan <- &lruMsg{                   //Delete from LruList
		operation: del,
		resp:      resp,
	}

	// Delete bytes cache
	// Not need for atomic
	cacheSize -= resp.size()
}

func (rCache *resourceTtlLruMap) ttl() {

	// Function to send a message when the timer expires
	backToFuture := func() {
		rCache.ttlChan <- true
	}

	// A timer.
	future := time.AfterFunc(24*time.Hour, backToFuture)

	for {

		<-rCache.ttlChan

		//Full Lock
		rCache.rwMutex.Lock()

		now := time.Now()

		// Traverse the skiplist which is ordered by ttl.
		// We do this by looping at level 0
		for node := rCache.skipList.head.next[0]; node != nil; node = node.next[0] {

			timeLeft := node.ttl.Sub(now)

			// If we still have time, check the timer and break
			if timeLeft > 0 {
				if !future.Reset(timeLeft) {
					future = time.AfterFunc(timeLeft, backToFuture)
				}

				break
			}

			// Remove from cache if time's up
			rCache.remove(node.key, rCache.cache[node.key])
		}

		rCache.rwMutex.Unlock()
	}
}
