package render

import (
	"github.com/webx-top/echo/logger"
	. "github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/echo/middleware/render/standard"
)

var engines = make(map[string]func(string) Driver)

func New(key string, tmplDir string, args ...logger.Logger) Driver {
	if fn, ok := engines[key]; ok {
		return fn(tmplDir)
	}
	return standard.New(tmplDir, args...)
}

func Reg(key string, val func(string) Driver) {
	engines[key] = val
}

func Del(key string) {
	if _, ok := engines[key]; ok {
		delete(engines, key)
	}
}
