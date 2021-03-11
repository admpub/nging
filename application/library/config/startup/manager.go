package startup

import (
	"fmt"
	"sort"

	"github.com/webx-top/echo/param"
)

var events = param.NewMap()

// MustGetEvent 获取事件
func MustGetEvent(typ string) *Event {
	e, ok := events.GetOk(typ)
	if !ok {
		e = NewEvent()
		events.Set(typ, e)
	}
	return e.(*Event)
}

// TypeList 事件类型列表
func TypeList() []string {
	var list []string
	events.Range(func(key, val interface{}) bool {
		list = append(list, fmt.Sprint(key))
		return true
	})
	sort.Strings(list)
	return list
}

// OnBefore 监听前置事件
func OnBefore(typ string, eventFunc EventFunc) {
	e := MustGetEvent(typ)
	e.AddBefore(eventFunc)
}

// OnAfter 监听后置事件
func OnAfter(typ string, eventFunc EventFunc) {
	e := MustGetEvent(typ)
	e.AddAfter(eventFunc)
}

// FireBefore 触发前置事件
func FireBefore(typ string) {
	e, ok := events.GetOk(typ)
	if !ok {
		return
	}
	e.(*Event).RunBefore()
}

// FireAfter 触发后置事件
func FireAfter(typ string) {
	e, ok := events.GetOk(typ)
	if !ok {
		return
	}
	e.(*Event).RunAfter()
}
