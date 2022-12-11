package background

import "sync"

func NewGroup() *Group {
	return &Group{
		m: map[string]*Background{},
	}
}

// Group 执行信息
type Group struct {
	m  map[string]*Background
	mu sync.RWMutex
}

// Cancel 取消某个任务
func (e *Group) Cancel(cacheKeys ...string) {
	if len(cacheKeys) == 0 {
		return
	}
	e.mu.Lock()
	e.cancel(cacheKeys...)
	e.mu.Unlock()
}

func (e *Group) cancel(cacheKeys ...string) {
	for _, cacheKey := range cacheKeys {
		if bgExec, ok := (*e).m[cacheKey]; ok {
			bgExec.cancel()
			delete((*e).m, cacheKey)
		}
	}
}

// Exists 任务是否存在
func (e *Group) Exists(cacheKey string) bool {
	e.mu.RLock()
	_, ok := (*e).m[cacheKey]
	e.mu.RUnlock()
	return ok
}

// Map 任务列表
func (e *Group) Map() map[string]*Background {
	e.mu.RLock()
	r := (*e).m
	e.mu.RUnlock()
	return r
}

// Add 新增任务
func (e *Group) Add(op string, cacheKey string, bgExec *Background) {
	e.mu.Lock()

	e.cancel(cacheKey) // 避免被覆盖后旧任务失去控制，先取消已存在的任务

	(*e).m[cacheKey] = bgExec

	e.mu.Unlock()

	Backgrounds.Store(op, e)
}
