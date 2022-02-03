package roleutils

import (
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/nging/v4/application/library/perm"
	"github.com/admpub/nging/v4/application/registry/navigate"
	"github.com/webx-top/echo"
)

func NewRolePermission() *RolePermission {
	return &RolePermission{
		DefinedType: UserRolePermissionType,
		Combined:    map[string]string{},
		parsed:      map[string]interface{}{},
	}
}

type RolePermission struct {
	DefinedType *echo.KVData
	Combined    map[string]string
	Roles       []*UserRoleWithPermissions
	parsed      map[string]interface{}
	filter      *navigate.Filter
}

func (r *RolePermission) Init(roleList []*UserRoleWithPermissions) *RolePermission {
	if r.filter == nil {
		r.filter = navigate.NewFilter(r)
	}
	r.Roles = roleList
	checkeds := map[string]map[string]struct{}{}
	seperators := map[string]string{}
	for _, role := range roleList {
		for _, rolePerm := range role.Permissions {
			if _, ok := checkeds[rolePerm.Type]; !ok {
				checkeds[rolePerm.Type] = map[string]struct{}{}
			}
			if _, ok := checkeds[rolePerm.Type][`*`]; ok {
				continue
			}
			for _, pa := range strings.Split(rolePerm.Permission, `,`) {
				if _, ok := checkeds[rolePerm.Type][pa]; !ok {
					checkeds[rolePerm.Type][pa] = struct{}{}
					r.Combined[rolePerm.Type] += seperators[rolePerm.Type] + pa
					seperators[rolePerm.Type] = `,`
				}
			}
		}
	}
	return r
}

func (r *RolePermission) onceParse(ctx echo.Context, typ string) bool {
	if _, ok := r.parsed[typ]; !ok {
		item := r.DefinedType.GetItem(typ)
		if item == nil {
			return false
		}
		var err error
		r.parsed[typ], err = item.X.(*perm.Handle).Parse(ctx, r.Combined[typ])
		//echo.Dump(echo.H{`combined`: r.Combined, `parsed`: r.parsed[typ]})
		if err != nil {
			log.Error(err)
			return false
		}
	}
	return true
}

func (r *RolePermission) CheckByType(ctx echo.Context, typ string, permPath string) interface{} {
	if !r.onceParse(ctx, typ) {
		return nil
	}
	rs, err := perm.HandleCheck(ctx, r.DefinedType, permPath, typ, r.Combined[typ], r.parsed[typ])
	if err != nil {
		log.Error(err)
	}
	return rs
}

func (r *RolePermission) Check(ctx echo.Context, permPath string) bool {
	rs := r.CheckByType(ctx, UserRolePermissionTypePage, permPath)
	if rs == nil {
		return false
	}
	if bl, ok := rs.(bool); ok {
		return bl
	}
	return false
}

func (r *RolePermission) CheckCmd(ctx echo.Context, permPath string) bool {
	rs := r.CheckByType(ctx, UserRolePermissionTypeCommand, permPath)
	if rs == nil {
		return false
	}
	if bl, ok := rs.(bool); ok {
		return bl
	}
	return false
}

func (r *RolePermission) CheckBehavior(ctx echo.Context, permPath string) *perm.CheckedBehavior {
	rs := r.CheckByType(ctx, UserRolePermissionTypeBehavior, permPath)
	if rs == nil {
		return &perm.CheckedBehavior{}
	}
	if bv, ok := rs.(*perm.CheckedBehavior); ok {
		return bv
	}
	return &perm.CheckedBehavior{}
}

//FilterNavigate 过滤导航菜单，只显示有权限的菜单
func (r *RolePermission) FilterNavigate(ctx echo.Context, navList *navigate.List) navigate.List {
	return r.filter.FilterNavigate(ctx, navList)
}
