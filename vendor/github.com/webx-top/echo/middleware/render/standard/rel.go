package standard

import (
	htmlTpl "html/template"
	"sync"
)

type tplInfo struct {
	Template *htmlTpl.Template
	Blocks   map[string]struct{}
}

func NewTplInfo(t *htmlTpl.Template) *tplInfo {
	return &tplInfo{
		Template: t,
		Blocks:   map[string]struct{}{},
	}
}

func NewRel(cachedKey string) *CcRel {
	return &CcRel{
		Rel: map[string]uint8{cachedKey: 0},
		Tpl: [2]*tplInfo{NewTplInfo(nil), NewTplInfo(nil)},
	}
}

type CcRel struct {
	Rel map[string]uint8
	Tpl [2]*tplInfo //0是独立模板；1是子模板
	l   sync.RWMutex
}

func (c *CcRel) GetOk(cacheKey string) (uint8, bool) {
	c.l.RLock()
	r, y := c.Rel[cacheKey]
	c.l.RUnlock()
	return r, y
}

func (c *CcRel) Set(cacheKey string, r uint8) {
	c.l.Lock()
	c.Rel[cacheKey] = r
	c.l.Unlock()
}

func (c *CcRel) Remove(cacheKey string) {
	c.l.Lock()
	delete(c.Rel, cacheKey)
	c.l.Unlock()
}

func (c *CcRel) Reset() {
	c.l.Lock()
	c.Rel = map[string]uint8{}
	c.Tpl = [2]*tplInfo{NewTplInfo(nil), NewTplInfo(nil)}
	c.l.Unlock()
}
