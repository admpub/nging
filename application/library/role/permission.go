package role

import (
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/nging/v4/application/library/perm"
	"github.com/admpub/nging/v4/application/registry/navigate"
	"github.com/webx-top/echo"
)

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

func NewCommonPermission(d *echo.KVData, c navigate.Checker) *CommonPermission {
	return &CommonPermission{
		DefinedType: d,
		filter:      navigate.NewFilter(c),
		Combined:    map[string]string{},
		parsed:      map[string]interface{}{},
	}
}

type CommonPermission struct {
	DefinedType *echo.KVData
	Combined    map[string]string
	parsed      map[string]interface{}
	filter      *navigate.Filter
}

type PermissionConfiger interface {
	GetType() string
	GetPermission() string
}

type PermissionsGetter interface {
	GetPermissions() []PermissionConfiger
}

func (r *CommonPermission) Init(roleList []PermissionsGetter) *CommonPermission {
	checkeds := map[string]map[string]struct{}{}
	seperators := map[string]string{}
	for _, role := range roleList {
		for _, rolePerm := range role.GetPermissions() {
			if _, ok := checkeds[rolePerm.GetType()]; !ok {
				checkeds[rolePerm.GetType()] = map[string]struct{}{}
			}
			if _, ok := checkeds[rolePerm.GetType()][`*`]; ok {
				continue
			}
			for _, pa := range strings.Split(rolePerm.GetPermission(), `,`) {
				if _, ok := checkeds[rolePerm.GetType()][pa]; !ok {
					checkeds[rolePerm.GetType()][pa] = struct{}{}
					r.Combined[rolePerm.GetType()] += seperators[rolePerm.GetType()] + pa
					seperators[rolePerm.GetType()] = `,`
				}
			}
		}
	}
	return r
}

func (r *CommonPermission) onceParse(ctx echo.Context, typ string) bool {
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

func (r *CommonPermission) CheckByType(ctx echo.Context, typ string, permPath string) interface{} {
	if !r.onceParse(ctx, typ) {
		return nil
	}
	rs, err := perm.HandleCheck(ctx, r.DefinedType, permPath, typ, r.Combined[typ], r.parsed[typ])
	if err != nil {
		log.Error(err)
	}
	return rs
}

func (r *CommonPermission) Check(ctx echo.Context, permPath string) bool {
	rs := r.CheckByType(ctx, RolePermissionTypePage, permPath)
	if rs == nil {
		return false
	}
	if bl, ok := rs.(bool); ok {
		return bl
	}
	return false
}

func (r *CommonPermission) CheckCmd(ctx echo.Context, permPath string) bool {
	rs := r.CheckByType(ctx, RolePermissionTypeCommand, permPath)
	if rs == nil {
		return false
	}
	if bl, ok := rs.(bool); ok {
		return bl
	}
	return false
}

func (r *CommonPermission) CheckBehavior(ctx echo.Context, permPath string) *perm.CheckedBehavior {
	rs := r.CheckByType(ctx, RolePermissionTypeBehavior, permPath)
	if rs == nil {
		return &perm.CheckedBehavior{}
	}
	if bv, ok := rs.(*perm.CheckedBehavior); ok {
		return bv
	}
	return &perm.CheckedBehavior{}
}

//FilterNavigate 过滤导航菜单，只显示有权限的菜单
func (r *CommonPermission) FilterNavigate(ctx echo.Context, navList *navigate.List) navigate.List {
	return r.filter.FilterNavigate(ctx, navList)
}
