package request

import (
	"github.com/docker/docker/api/types/swarm"
	"github.com/webx-top/echo"
)

type SwarmNodeEdit struct {
	swarm.NodeSpec
}

func (a *SwarmNodeEdit) ValueDecoders(echo.Context) echo.BinderValueCustomDecoders {
	return map[string]echo.BinderValueCustomDecoder{
		`Labels`: mapDecoder,
	}
}
