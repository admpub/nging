package initconfig

import (
	"github.com/webx-top/echo/param"
)

var events = param.NewMap()

func GetEvent() *Event {
	return events.Get(typ, func() interface{} {
		return NewEvent()
	}).(*Event)
}

func OnBefore(typ string, eventFunc EventFunc) {
	e := GetEvent()
	e.AddBefore(eventFunc)
}

func OnAfter(typ string, eventFunc EventFunc) {
	e := GetEvent()
	e.AddAfter(eventFunc)
}
