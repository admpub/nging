package ttlmap

import (
	"container/heap"
	"sync"
)

type store struct {
	sync.RWMutex
	kv           map[string]*pqitem
	pq           pqueue
	onWillExpire func(key string, item Item)
	onWillEvict  func(key string, item Item)
}

func newStore(opts *Options) *store {
	return &store{
		kv:           make(map[string]*pqitem, opts.InitialCapacity),
		pq:           make(pqueue, 0, opts.InitialCapacity),
		onWillExpire: opts.OnWillExpire,
		onWillEvict:  opts.OnWillEvict,
	}
}

func (s *store) set(pqi *pqitem) {
	s.kv[pqi.key] = pqi
	heap.Push(&s.pq, pqi)
}

func (s *store) delete(pqi *pqitem) {
	delete(s.kv, pqi.key)
	heap.Remove(&s.pq, pqi.index)
}

func (s *store) fix(pqi *pqitem) {
	heap.Fix(&s.pq, pqi.index)
}

func (s *store) tryExpire(pqi *pqitem) bool {
	if pqi.item.Expired() {
		if s.onWillExpire != nil {
			s.onWillExpire(pqi.key, *pqi.item)
		}
		s.evict(pqi)
		return true
	}
	return false
}

func (s *store) evict(pqi *pqitem) {
	if s.onWillEvict != nil {
		s.onWillEvict(pqi.key, *pqi.item)
	}
	s.delete(pqi)
}

func (s *store) evictExpired() {
	for pqi := s.pq.peek(); pqi != nil; pqi = s.pq.peek() {
		if !s.tryExpire(pqi) {
			return
		}
	}
}

func (s *store) drain() {
	for _, pqi := range s.pq {
		if s.onWillEvict != nil {
			s.onWillEvict(pqi.key, *pqi.item)
		}
	}
	s.kv = nil
	s.pq = nil
}
