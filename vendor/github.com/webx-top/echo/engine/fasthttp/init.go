package fasthttp

import (
	"github.com/webx-top/echo/engine"
)

const Name = `fasthttp`

func init() {
	engine.Register(Name, func(c *engine.Config) engine.Engine {
		return NewWithConfig(c)
	})
}
