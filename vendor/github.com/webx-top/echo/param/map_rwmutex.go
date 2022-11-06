package param

import (
	"encoding/xml"
	"html/template"
	"sync"
	"time"
)

func NewSafeStore(data ...Store) *SafeStore {
	var store Store
	if len(data) > 0 {
		store = data[0]
	}
	if store == nil {
		store = Store{}
	}
	return &SafeStore{store: store}
}

type SafeStore struct {
	store     Store
	storeLock sync.RWMutex
}

// Get retrieves data
func (s *SafeStore) Get(key string, defaults ...interface{}) interface{} {
	s.storeLock.RLock()
	v := s.store.Get(key, defaults...)
	s.storeLock.RUnlock()
	return v
}

// Set saves data
func (s *SafeStore) Set(key string, val interface{}) {
	s.storeLock.Lock()
	s.store.Set(key, val)
	s.storeLock.Unlock()
}

// Delete saves data
func (c *SafeStore) Delete(keys ...string) {
	c.storeLock.Lock()
	c.store.Delete(keys...)
	c.storeLock.Unlock()
}

func (s *SafeStore) Stored() Store {
	s.storeLock.RLock()
	copied := s.store.Clone()
	s.storeLock.RUnlock()
	return copied
}

func (s *SafeStore) Has(key string) bool {
	s.storeLock.RLock()
	has := s.store.Has(key)
	s.storeLock.RUnlock()
	return has
}

func (s *SafeStore) String(key string, defaults ...interface{}) string {
	s.storeLock.RLock()
	str := s.store.String(key, defaults...)
	s.storeLock.RUnlock()
	return str
}

func (s *SafeStore) Split(key string, sep string, limit ...int) StringSlice {
	s.storeLock.RLock()
	val := s.store.Split(key, sep, limit...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) Trim(key string, defaults ...interface{}) String {
	s.storeLock.RLock()
	val := s.store.Trim(key, defaults...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) HTML(key string, defaults ...interface{}) template.HTML {
	s.storeLock.RLock()
	val := s.store.HTML(key, defaults...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) HTMLAttr(key string, defaults ...interface{}) template.HTMLAttr {
	s.storeLock.RLock()
	val := s.store.HTMLAttr(key, defaults...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) JS(key string, defaults ...interface{}) template.JS {
	s.storeLock.RLock()
	val := s.store.JS(key, defaults...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) CSS(key string, defaults ...interface{}) template.CSS {
	s.storeLock.RLock()
	val := s.store.CSS(key, defaults...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) Bool(key string, defaults ...interface{}) bool {
	s.storeLock.RLock()
	val := s.store.Bool(key, defaults...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) Float64(key string, defaults ...interface{}) float64 {
	s.storeLock.RLock()
	val := s.store.Float64(key, defaults...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) Float32(key string, defaults ...interface{}) float32 {
	s.storeLock.RLock()
	val := s.store.Float32(key, defaults...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) Int8(key string, defaults ...interface{}) int8 {
	s.storeLock.RLock()
	val := s.store.Int8(key, defaults...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) Int16(key string, defaults ...interface{}) int16 {
	s.storeLock.RLock()
	val := s.store.Int16(key, defaults...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) Int(key string, defaults ...interface{}) int {
	s.storeLock.RLock()
	val := s.store.Int(key, defaults...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) Int32(key string, defaults ...interface{}) int32 {
	s.storeLock.RLock()
	val := s.store.Int32(key, defaults...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) Int64(key string, defaults ...interface{}) int64 {
	s.storeLock.RLock()
	val := s.store.Int64(key, defaults...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) Decr(key string, n int64, defaults ...interface{}) int64 {
	s.storeLock.Lock()
	val := s.store.Decr(key, n, defaults...)
	s.storeLock.Unlock()
	return val
}

func (s *SafeStore) Incr(key string, n int64, defaults ...interface{}) int64 {
	s.storeLock.Lock()
	val := s.store.Incr(key, n, defaults...)
	s.storeLock.Unlock()
	return val
}

func (s *SafeStore) Uint8(key string, defaults ...interface{}) uint8 {
	s.storeLock.RLock()
	val := s.store.Uint8(key, defaults...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) Uint16(key string, defaults ...interface{}) uint16 {
	s.storeLock.RLock()
	val := s.store.Uint16(key, defaults...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) Uint(key string, defaults ...interface{}) uint {
	s.storeLock.RLock()
	val := s.store.Uint(key, defaults...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) Uint32(key string, defaults ...interface{}) uint32 {
	s.storeLock.RLock()
	val := s.store.Uint32(key, defaults...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) Uint64(key string, defaults ...interface{}) uint64 {
	s.storeLock.RLock()
	val := s.store.Uint64(key, defaults...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) Timestamp(key string, defaults ...interface{}) time.Time {
	s.storeLock.RLock()
	val := s.store.Timestamp(key, defaults...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) Duration(key string, defaults ...time.Duration) time.Duration {
	s.storeLock.RLock()
	val := s.store.Duration(key, defaults...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) DateTime(key string, layouts ...string) time.Time {
	s.storeLock.RLock()
	val := s.store.DateTime(key, layouts...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) Children(keys ...interface{}) Store {
	s.storeLock.RLock()
	val := s.store.Children(keys...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) GetStore(key string, defaults ...interface{}) Store {
	s.storeLock.RLock()
	val := s.store.GetStore(key, defaults...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) GetStoreByKeys(keys ...string) Store {
	s.storeLock.RLock()
	val := s.store.GetStoreByKeys(keys...)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) Select(selectKeys ...string) Store {
	s.storeLock.RLock()
	val := s.store.Select(selectKeys...)
	s.storeLock.RUnlock()
	return val
}

// MarshalXML allows type Store to be used with xml.Marshal
func (s *SafeStore) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	s.storeLock.RLock()
	err := s.store.MarshalXML(e, start)
	s.storeLock.RUnlock()
	return err
}

func (s *SafeStore) DeepMerge(source Store) {
	s.storeLock.Lock()
	s.store.DeepMerge(source)
	s.storeLock.Unlock()
}

func (s *SafeStore) Clone() *SafeStore {
	return NewSafeStore(s.CloneStore())
}

func (s *SafeStore) CloneStore() Store {
	s.storeLock.RLock()
	r := s.store.Clone()
	s.storeLock.RUnlock()
	return r
}

func (s *SafeStore) Transform(transfers map[string]Transfer) Store {
	s.storeLock.RLock()
	val := s.store.Transform(transfers)
	s.storeLock.RUnlock()
	return val
}

func (s *SafeStore) SetMKey(key string, value interface{}) *SafeStore {
	s.storeLock.Lock()
	s.store.SetMKey(key, value)
	s.storeLock.Unlock()
	return s
}

func (s *SafeStore) SetMKeys(keys []string, value interface{}) *SafeStore {
	s.storeLock.Lock()
	s.store.SetMKeys(keys, value)
	s.storeLock.Unlock()
	return s
}

func (s *SafeStore) Clear() *SafeStore {
	s.storeLock.Lock()
	s.store = Store{}
	s.storeLock.Unlock()
	return s
}

func (s *SafeStore) Range(fn func(key string, val interface{}) error) (err error) {
	s.storeLock.Lock()
	for key, val := range s.store {
		if err = fn(key, val); err != nil {
			break
		}
	}
	s.storeLock.Unlock()
	return
}
