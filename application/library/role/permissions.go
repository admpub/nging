package role

import (
	"github.com/admpub/nging/v4/application/dbschema"
)

type UserRoleWithPermissions struct {
	*dbschema.NgingUserRole
	Permissions []*dbschema.NgingUserRolePermission `db:"-,relation=role_id:id"`
}
