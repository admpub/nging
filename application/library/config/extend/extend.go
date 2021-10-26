package extend

import (
	"github.com/admpub/nging/v3/application/library/hook"
	"github.com/webx-top/echo"
)

type Initer func() interface{}

var Hook = hook.New()

type Reloader interface {
	Reload() error
}

type SetDefaults interface {
	SetDefaults()
}

var extendIniters = map[string]Initer{}

func Register(name string, initer Initer) {
	extendIniters[name] = initer
}

func Range(f func(string, interface{})) {
	for name, initer := range extendIniters {
		f(name, initer())
	}
}

func Get(name string) Initer {
	initer, _ := extendIniters[name]
	return initer
}

func Unregister(name string) {
	if initer, ok := extendIniters[name]; ok {
		if err := Hook.Fire(`unregister`, echo.H{`name`: name, `initer`: initer}); err != nil {
			panic(err)
		}
		delete(extendIniters, name)
	}
}
