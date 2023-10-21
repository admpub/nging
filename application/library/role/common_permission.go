package role

import (
	"encoding/json"
	"strings"

	"github.com/admpub/copier"
	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/perm"
	"github.com/admpub/nging/v5/application/registry/navigate"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"
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
	ctx := defaults.NewMockContext()
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
			if len(permissionStr) == 0 {
				continue
			}
			if strings.HasPrefix(permissionStr, `{`) && strings.HasSuffix(permissionStr, `}`) {
				r.combineJSON(ctx, checkeds, rolePerm.GetType(), permissionStr)
				continue
			}
			for _, permVal := range strings.Split(permissionStr, `,`) {
				if _, ok := checkeds[rolePerm.GetType()][permVal]; !ok {
					checkeds[rolePerm.GetType()][permVal] = struct{}{}
					r.Combined[rolePerm.GetType()] += seperators[rolePerm.GetType()] + permVal
					seperators[rolePerm.GetType()] = `,`
				}
			}
		}
	}
	return r
}

func (r *CommonPermission) combineBehaviorType(ctx echo.Context, checkeds map[string]map[string]interface{}, permType string, permRule string) {
	item := r.DefinedType.GetItem(permType)
	if item == nil {
		return
	}
	parsed, err := item.X.(*perm.Handle).Parse(ctx, permRule)
	//echo.Dump(echo.H{`rule`: permRule, `parsed`: parsed})
	if err != nil {
		log.Errorf(`failed to parse permission(%s): %v`, permRule, err)
		return
	}
	r.parsed[permType] = parsed
	perms, ok := parsed.(perm.BehaviorPerms)
	if !ok {
		return
	}
	for permKey, permVal := range perms {
		last, ok := checkeds[permType][permKey]
		if !ok {
			checkeds[permType][permKey] = permVal.Value
		} else if combine, ok := last.(Combiner); ok {
			lastRecv := combine.Combine(permVal.Value)
			checkeds[permType][permKey] = lastRecv
		} else {
			err = copier.Copy(last, permVal.Value)
			if err != nil {
				log.Errorf(`failed to copy %#v to %#v: %v`, permVal.Value, last, err)
				continue
			}
			checkeds[permType][permKey] = last
		}
	}
	b, _ := json.Marshal(checkeds[permType])
	r.Combined[permType] = string(b)
}

func (r *CommonPermission) combineJSON(ctx echo.Context, checkeds map[string]map[string]interface{}, permType string, permRule string) {
	if permType == RolePermissionTypeBehavior {
		r.combineBehaviorType(ctx, checkeds, permType, permRule)
		return
	}
	recv := echo.H{}
	jsonBytes := com.Str2bytes(permRule)
	if err := json.Unmarshal(jsonBytes, &recv); err != nil {
		log.Error(common.JSONBytesParseError(err, jsonBytes).Error())
		return
	}
	for permKey, permVal := range recv {
		last, ok := checkeds[permType][permKey]
		if !ok {
			checkeds[permType][permKey] = permVal
		} else if lastRecv, ok := last.(echo.H); ok {
			for k, v := range param.AsStore(permVal) {
				lastRecv[k] = v
			}
			checkeds[permType][permKey] = lastRecv
		}
	}
	b, _ := json.Marshal(checkeds[permType])
	r.Combined[permType] = string(b)
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
