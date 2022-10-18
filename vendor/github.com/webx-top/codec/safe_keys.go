package codec

import "sync"

func NewSafeKeys() *SafeKeys {
	return &SafeKeys{
		keys: map[string][]byte{},
	}
}

type SafeKeys struct {
	keys map[string][]byte
	mu   sync.RWMutex
}

func (s *SafeKeys) Set(rawKey string, fixedKey []byte) {
	s.mu.Lock()
	s.keys[rawKey] = fixedKey
	s.mu.Unlock()
}

func (s *SafeKeys) Get(rawKey string) (fixedKey []byte, ok bool) {
	s.mu.RLock()
	fixedKey, ok = s.keys[rawKey]
	s.mu.RUnlock()
	return
}
