package perm

import (
	"github.com/admpub/nging/v4/application/library/role"
)

func New() *role.RolePermission {
	return role.NewRolePermission()
}

type RolePermission = role.RolePermission
