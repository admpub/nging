package com

import (
	"sync"
)

type SafeMap struct {
	lock *sync.RWMutex
	bm   map[interface{}]interface{}
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		lock: new(sync.RWMutex),
		bm:   make(map[interface{}]interface{}),
	}
}

//Get from maps return the k's value
func (m *SafeMap) Get(k interface{}) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if val, ok := m.bm[k]; ok {
		return val
	}
	return nil
}

// Set maps the given key and value. Returns false
// if the key is already in the map and changes nothing.
func (m *SafeMap) Set(k interface{}, v interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if val, ok := m.bm[k]; !ok {
		m.bm[k] = v
	} else if val != v {
		m.bm[k] = v
	} else {
		return false
	}
	return true
}

// Check returns true if k is exist in the map.
func (m *SafeMap) Check(k interface{}) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if _, ok := m.bm[k]; !ok {
		return false
	}
	return true
}

func (m *SafeMap) Delete(k interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.bm, k)
}

func (m *SafeMap) Items() map[interface{}]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.bm
}

func NewOrderlySafeMap() *OrderlySafeMap {
	return &OrderlySafeMap{
		SafeMap: NewSafeMap(),
		keys:    make([]interface{}, 0),
	}
}

type OrderlySafeMap struct {
	*SafeMap
	keys   []interface{} // map keys
	values []interface{} // map values
}

func (m *OrderlySafeMap) Set(k interface{}, v interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if val, ok := m.bm[k]; !ok {
		m.bm[k] = v
		m.keys = append(m.keys, k)
	} else if val != v {
		m.bm[k] = v
	} else {
		return false
	}
	return true
}

func (m *OrderlySafeMap) Delete(k interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.bm, k)
	endIndex := len(m.keys) - 1
	for index, mapKey := range m.keys {
		if mapKey != k {
			continue
		}
		if index == endIndex {
			m.keys = m.keys[0:index]
			break
		}
		if index == 0 {
			m.keys = m.keys[1:]
			break
		}
		m.keys = append(m.keys[0:index], m.keys[index+1:]...)
		break
	}
}

func (m *OrderlySafeMap) Keys() []interface{} {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.keys
}

func (m *OrderlySafeMap) Values(force ...bool) []interface{} {
	m.lock.Lock()
	defer m.lock.Unlock()
	if (len(force) == 0 || !force[0]) && m.values != nil {
		return m.values
	}
	m.values = []interface{}{}
	for _, mapKey := range m.keys {
		m.values = append(m.values, m.bm[mapKey])
	}
	return m.values
}

func (m *OrderlySafeMap) VisitAll(callback func(int, interface{}, interface{})) {
	m.lock.Lock()
	for index, mapKey := range m.keys {
		callback(index, mapKey, m.bm[mapKey])
	}
	m.lock.Unlock()
}
