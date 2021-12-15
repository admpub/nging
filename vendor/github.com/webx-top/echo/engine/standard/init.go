package standard

import (
	"github.com/webx-top/echo/engine"
)

const Name = `standard`

func init() {
	engine.Register(Name, func(c *engine.Config) engine.Engine {
		return NewWithConfig(c)
	})
}
