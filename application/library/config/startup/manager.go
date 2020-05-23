package startup

import (
	"fmt"
	"sort"
	"github.com/webx-top/echo/param"
)

var events = param.NewMap()

func MustGetEvent(typ string) *Event {
	e, ok := events.GetOk(typ)
	if !ok {
		e = NewEvent()
		events.Set(typ, e)
	}
	return e.(*Event)
}

func TypeList() []string {
	var list []string
	events.Range(func(key, val interface{}) bool {
		list = append(list, fmt.Sprint(key))
		return true
	})
	sort.Strings(list)
	return list
}

func OnBefore(typ string, eventFunc EventFunc) {
	e := MustGetEvent(typ)
	e.AddBefore(eventFunc)
}

func OnAfter(typ string, eventFunc EventFunc) {
	e := MustGetEvent(typ)
	e.AddAfter(eventFunc)
}

func FireBefore(typ string) {
	e, ok := events.GetOk(typ)
	if !ok {
		return
	}
	e.(*Event).RunBefore()
}

func FireAfter(typ string) {
	e, ok := events.GetOk(typ)
	if !ok {
		return
	}
	e.(*Event).RunAfter()
}
