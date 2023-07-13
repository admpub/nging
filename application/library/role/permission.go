package role

import (
	"github.com/webx-top/echo"
)

type ICheckByType interface {
	CheckByType(ctx echo.Context, typ string, permPath string) interface{}
}

type PermissionConfiger interface {
	GetType() string
	GetPermission() string
}

type PermissionsGetter interface {
	GetPermissions() []PermissionConfiger
}
