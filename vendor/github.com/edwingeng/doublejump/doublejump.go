// Package doublejump provides a revamped Google's jump consistent hash.
package doublejump

import (
	"math/rand"

	"github.com/dgryski/go-jump"
)

type looseHolder struct {
	a []interface{}
	m map[interface{}]int
	f []int
}

func (holder *looseHolder) add(obj interface{}) {
	if _, ok := holder.m[obj]; ok {
		return
	}

	if nf := len(holder.f); nf == 0 {
		holder.a = append(holder.a, obj)
		holder.m[obj] = len(holder.a) - 1
	} else {
		idx := holder.f[nf-1]
		holder.f = holder.f[:nf-1]
		holder.a[idx] = obj
		holder.m[obj] = idx
	}
}

func (holder *looseHolder) remove(obj interface{}) {
	if idx, ok := holder.m[obj]; ok {
		holder.f = append(holder.f, idx)
		holder.a[idx] = nil
		delete(holder.m, obj)
	}
}

func (holder *looseHolder) get(key uint64) interface{} {
	na := len(holder.a)
	if na == 0 {
		return nil
	}

	h := jump.Hash(key, na)
	return holder.a[h]
}

func (holder *looseHolder) shrink() {
	if len(holder.f) == 0 {
		return
	}

	var a []interface{}
	for _, obj := range holder.a {
		if obj != nil {
			a = append(a, obj)
			holder.m[obj] = len(a) - 1
		}
	}
	holder.a = a
	holder.f = nil
}

type compactHolder struct {
	a []interface{}
	m map[interface{}]int
}

func (holder *compactHolder) add(obj interface{}) {
	if _, ok := holder.m[obj]; ok {
		return
	}

	holder.a = append(holder.a, obj)
	holder.m[obj] = len(holder.a) - 1
}

func (holder *compactHolder) shrink(a []interface{}) {
	for i, obj := range a {
		holder.a[i] = obj
		holder.m[obj] = i
	}
}

func (holder *compactHolder) remove(obj interface{}) {
	if idx, ok := holder.m[obj]; ok {
		n := len(holder.a)
		holder.a[idx] = holder.a[n-1]
		holder.m[holder.a[idx]] = idx
		holder.a[n-1] = nil
		holder.a = holder.a[:n-1]
		delete(holder.m, obj)
	}
}

func (holder *compactHolder) get(key uint64) interface{} {
	na := len(holder.a)
	if na == 0 {
		return nil
	}

	h := jump.Hash(key*0xc6a4a7935bd1e995, na)
	return holder.a[h]
}

// Hash is a revamped Google's jump consistent hash. It overcomes the shortcoming of
// the original implementation - being unable to remove nodes.
//
// Hash is NOT thread-safe.
type Hash struct {
	loose   looseHolder
	compact compactHolder
}

// NewHash creates a new doublejump hash instance.
func NewHash() *Hash {
	hash := &Hash{}
	hash.loose.m = make(map[interface{}]int)
	hash.compact.m = make(map[interface{}]int)
	return hash
}

// Add adds an object to the hash.
func (h *Hash) Add(obj interface{}) {
	if obj == nil {
		return
	}

	h.loose.add(obj)
	h.compact.add(obj)
}

// Remove removes an object from the hash.
func (h *Hash) Remove(obj interface{}) {
	if obj == nil {
		return
	}

	h.loose.remove(obj)
	h.compact.remove(obj)
}

// Len returns the number of objects in the hash.
func (h *Hash) Len() int {
	return len(h.compact.a)
}

// LooseLen returns the size of the inner loose object holder.
func (h *Hash) LooseLen() int {
	return len(h.loose.a)
}

// Shrink removes all empty slots from the hash.
func (h *Hash) Shrink() {
	h.loose.shrink()
	h.compact.shrink(h.loose.a)
}

// Get returns the existing object for the key, or nil if there is no object in the hash.
func (h *Hash) Get(key uint64) interface{} {
	obj := h.loose.get(key)
	switch obj {
	case nil:
		return h.compact.get(key)
	default:
		return obj
	}
}

// All returns all the objects in this Hash.
func (h *Hash) All() []interface{} {
	n := len(h.compact.a)
	if n == 0 {
		return nil
	}
	all := make([]interface{}, n)
	copy(all, h.compact.a)
	return all
}

// Random returns a random object, or nil if there is no object in the hash.
func (h *Hash) Random() interface{} {
	if n := len(h.compact.a); n > 0 {
		idx := rand.Intn(n)
		return h.compact.a[idx]
	}
	return nil
}
