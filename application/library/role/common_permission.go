package role

import (
	"encoding/json"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/perm"
	"github.com/admpub/nging/v5/application/registry/navigate"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

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

func (r *CommonPermission) Init(roleList []PermissionsGetter) *CommonPermission {
	checkeds := map[string]map[string]interface{}{}
	seperators := map[string]string{}
	for _, role := range roleList {
		for _, rolePerm := range role.GetPermissions() {
			if _, ok := checkeds[rolePerm.GetType()]; !ok {
				checkeds[rolePerm.GetType()] = map[string]interface{}{}
			}
			if _, ok := checkeds[rolePerm.GetType()][`*`]; ok {
				continue
			}
			permissionStr := rolePerm.GetPermission()
			permissionStr = strings.TrimSpace(permissionStr)
			if strings.HasPrefix(permissionStr, `{`) && strings.HasSuffix(permissionStr, `}`) {
				recv := echo.H{}
				jsonBytes := com.Str2bytes(permissionStr)
				if err := json.Unmarshal(jsonBytes, &recv); err != nil {
					log.Error(common.JSONBytesParseError(err, jsonBytes).Error())
					continue
				}
				for pa, val := range recv {
					last, ok := checkeds[rolePerm.GetType()][pa]
					if !ok {
						checkeds[rolePerm.GetType()][pa] = val
						r.Combined[rolePerm.GetType()] = permissionStr
					} else if lastRecv, ok := last.(echo.H); ok {
						for k, v := range recv {
							lastRecv[k] = v
						}
						checkeds[rolePerm.GetType()][pa] = lastRecv
						b, _ := json.Marshal(lastRecv)
						r.Combined[rolePerm.GetType()] = string(b)
					}
				}
				continue
			}
			for _, pa := range strings.Split(permissionStr, `,`) {
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

func (r *CommonPermission) Get(ctx echo.Context, typ string) interface{} {
	if !r.onceParse(ctx, typ) {
		return nil
	}
	return r.parsed[typ]
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

// FilterNavigate 过滤导航菜单，只显示有权限的菜单
func (r *CommonPermission) FilterNavigate(ctx echo.Context, navList *navigate.List) navigate.List {
	return r.filter.FilterNavigate(ctx, navList)
}
