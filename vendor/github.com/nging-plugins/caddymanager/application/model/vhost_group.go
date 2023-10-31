package model

import (
	"github.com/nging-plugins/caddymanager/application/dbschema"
	"github.com/webx-top/echo"
)

func NewVhostGroup(ctx echo.Context) *VhostGroup {
	return &VhostGroup{
		NgingVhostGroup: dbschema.NewNgingVhostGroup(ctx),
	}
}

type VhostGroup struct {
	*dbschema.NgingVhostGroup
}
