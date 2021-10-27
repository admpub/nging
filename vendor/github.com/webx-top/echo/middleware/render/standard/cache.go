package standard

import "sync"

func NewCache() *CacheData {
	return &CacheData{
		m: map[string]*CcRel{},
	}
}

type CacheData struct {
	m map[string]*CcRel
	l sync.RWMutex
}

func (c *CacheData) GetOk(cacheKey string) (*CcRel, bool) {
	c.l.RLock()
	defer c.l.RUnlock()
	r, y := c.m[cacheKey]
	return r, y
}

func (c *CacheData) Set(cacheKey string, r *CcRel) {
	c.l.Lock()
	c.m[cacheKey] = r
	c.l.Unlock()
}

func (c *CacheData) Remove(cacheKey string) {
	c.l.Lock()
	delete(c.m, cacheKey)
	c.l.Unlock()
}

func (c *CacheData) Reset() {
	c.l.Lock()
	c.m = map[string]*CcRel{}
	c.l.Unlock()
}
