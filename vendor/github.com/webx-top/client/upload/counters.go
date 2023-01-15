package upload

import "sync"

func NewCounters() *Counters {
	return &Counters{
		m: map[string]int{},
	}
}

type Counters struct {
	m  map[string]int
	mu sync.RWMutex
}

func (c *Counters) Add(key string, i int) {
	c.mu.Lock()
	c.m[key] += i
	c.mu.Unlock()
}

func (c *Counters) GetCount(key string) int {
	c.mu.RLock()
	i := c.m[key]
	c.mu.RUnlock()
	return i
}

func (c *Counters) Map() map[string]int {
	r := map[string]int{}
	c.Range(func(key string, i int) error {
		r[key] = i
		return nil
	})
	return r
}

func (c *Counters) Range(f func(key string, i int) error) (err error) {
	c.mu.RLock()
	for k, v := range c.m {
		err = f(k, v)
		if err != nil {
			break
		}
	}
	c.mu.RUnlock()
	return
}
