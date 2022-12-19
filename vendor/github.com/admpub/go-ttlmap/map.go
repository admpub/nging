// Package ttlmap provides a map-like interface with string keys and expirable
// items. Keys are currently limited to strings.
package ttlmap

import "errors"

// Errors returned Map operations.
var (
	ErrNotExist = errors.New("key does not exist")
	ErrExist    = errors.New("key already exists")
	ErrDrained  = errors.New("map was drained")
)

var zeroItem Item

// Map is the equivalent of a map[string]interface{} but with expirable Items.
type Map struct {
	store  *store
	keeper *keeper
}

// New creates a new Map with given options.
func New(opts *Options) *Map {
	if opts == nil {
		opts = &Options{}
	}
	store := newStore(opts)
	m := &Map{
		store:  store,
		keeper: newKeeper(store),
	}
	go m.keeper.run()
	return m
}

// Len returns the number of elements in the map.
func (m *Map) Len() int {
	m.store.RLock()
	n := len(m.store.kv)
	m.store.RUnlock()
	return n
}

// Get returns the item in the map with the given key.
// ErrNotExist will be returned if the key does not exist.
// ErrDrained will be returned if the map is already drained.
func (m *Map) Get(key string) (Item, error) {
	m.store.RLock()
	if m.keeper.drained {
		m.store.RUnlock()
		return zeroItem, ErrDrained
	}
	if pqi := m.store.kv[key]; pqi != nil {
		item := *pqi.item
		m.store.RUnlock()
		return item, nil
	}
	m.store.RUnlock()
	return zeroItem, ErrNotExist
}

// Set assigns an item with the specified key in the map.
// ErrExist or ErrNotExist may be returned depending on opts.KeyExist.
// ErrDrained will be returned if the map is already drained.
func (m *Map) Set(key string, item Item, opts *SetOptions) error {
	m.store.Lock()
	if m.keeper.drained {
		m.store.Unlock()
		return ErrDrained
	}
	err := m.set(key, &item, opts)
	m.store.Unlock()
	return err
}

// Update updates an item with the specified key in the map and returns it.
// ErrNotExist will be returned if the key does not exist.
// ErrDrained will be returned if the map is already drained.
func (m *Map) Update(key string, item Item, opts *UpdateOptions) (Item, error) {
	m.store.Lock()
	if m.keeper.drained {
		m.store.Unlock()
		return zeroItem, ErrDrained
	}
	if pqi := m.store.kv[key]; pqi != nil {
		m.update(pqi, &item, opts)
		item = *pqi.item
		m.store.Unlock()
		return item, nil
	}
	m.store.Unlock()
	return zeroItem, ErrNotExist
}

// Delete deletes the item with the specified key from the map.
// ErrNotExist will be returned if the key does not exist.
// ErrDrained will be returned if the map is already drained.
func (m *Map) Delete(key string) (Item, error) {
	m.store.Lock()
	if m.keeper.drained {
		m.store.Unlock()
		return zeroItem, ErrDrained
	}
	if pqi := m.store.kv[key]; pqi != nil {
		m.delete(pqi)
		item := *pqi.item
		m.store.Unlock()
		return item, nil
	}
	m.store.Unlock()
	return zeroItem, ErrNotExist
}

// Draining returns the channel that is closed when the map starts draining.
func (m *Map) Draining() <-chan struct{} {
	return m.keeper.drainingChan
}

// Drain evicts all remaining elements from the map and terminates the usage of
// this map.
func (m *Map) Drain() {
	m.keeper.signalDrain()
	<-m.keeper.doneChan
}

func (m *Map) set(key string, item *Item, opts *SetOptions) error {
	if pqi := m.store.kv[key]; pqi != nil {
		if opts.keyExist() == KeyExistNotYet {
			return ErrExist
		}
		m.expireOrEvict(pqi)
	} else if opts.keyExist() == KeyExistAlready {
		return ErrNotExist
	}
	pqi := &pqitem{
		key:   key,
		item:  item,
		index: -1,
	}
	m.store.set(pqi)
	if pqi.index == 0 {
		m.keeper.signalUpdate()
	}
	return nil
}

func (m *Map) update(pqi *pqitem, item *Item, opts *UpdateOptions) {
	if opts != nil {
		if opts.KeepValue {
			item.value = pqi.item.value
		}
		if opts.KeepExpiration {
			item.expiration = pqi.item.expiration
			item.expires = pqi.item.expires
		}
	}
	pqi.item = item
	m.store.fix(pqi)
	if pqi.index == 0 {
		m.keeper.signalUpdate()
	}
}

func (m *Map) expireOrEvict(pqi *pqitem) {
	if pqi.index == 0 {
		m.keeper.signalUpdate()
	}
	if !m.store.tryExpire(pqi) {
		m.store.evict(pqi)
	}
}

func (m *Map) delete(pqi *pqitem) {
	if pqi.index == 0 {
		m.keeper.signalUpdate()
	}
	m.store.delete(pqi)
}
