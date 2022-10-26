package perm

import (
	"github.com/admpub/nging/v5/application/library/role"
)

func New() *role.RolePermission {
	return role.NewRolePermission()
}

type RolePermission = role.RolePermission
