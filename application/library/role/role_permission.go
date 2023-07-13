package role

func NewRolePermission() *RolePermission {
	r := &RolePermission{}
	r.CommonPermission = NewCommonPermission(UserRolePermissionType, r)
	return r
}

type RolePermission struct {
	*CommonPermission
	Roles []*UserRoleWithPermissions
}

func (r *RolePermission) Init(roleList []*UserRoleWithPermissions) *RolePermission {
	r.Roles = roleList
	gts := make([]PermissionsGetter, len(r.Roles))
	for k, v := range r.Roles {
		gts[k] = v
	}
	r.CommonPermission.Init(gts)
	return r
}
