package background

import (
	"sync"

	"github.com/webx-top/echo/param"
)

// Backgrounds 后台任务集合
var Backgrounds = sync.Map{}

// All 所有任务
func All() map[string]map[string]*Background {
	r := map[string]map[string]*Background{}
	Backgrounds.Range(func(key, val interface{}) bool {
		r[param.AsString(key)] = val.(*Group).Map()
		return true
	})
	return r
}

// ListBy 获取某个操作的所有任务
func ListBy(op string) *Group {
	old, exists := Backgrounds.Load(op)
	if !exists {
		return nil
	}
	exec := old.(*Group)
	return exec
}

// Cancel 取消执行
func Cancel(op string, cacheKeys ...string) {
	if len(cacheKeys) == 0 {
		return
	}
	exec := ListBy(op)
	if exec == nil {
		return
	}
	exec.Cancel(cacheKeys...)
	Backgrounds.Store(op, exec)
}
