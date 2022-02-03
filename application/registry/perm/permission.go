package perm

import (
	"github.com/admpub/nging/v4/application/library/roleutils"
)

func New() *roleutils.RolePermission {
	return roleutils.NewRolePermission()
}

type RolePermission = roleutils.RolePermission
