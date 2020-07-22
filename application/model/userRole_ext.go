package model

import (
	"strings"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/perm"
	"github.com/admpub/nging/application/registry/navigate"
)

func NewPermission() *RolePermission {
	return &RolePermission{}
}

type RolePermission struct {
	Actions      string
	Commands     string
	Roles        []*dbschema.NgingUserRole
	permActions  *perm.Map
	permCommonds *perm.Map
}

func (r *RolePermission) Init(roleList []*dbschema.NgingUserRole) *RolePermission {
	r.Roles = roleList
	cmdChecked := map[string]struct{}{}
	actChecked := map[string]struct{}{}
	var cmdSep, actSep string
	r.Actions = ``
	r.Commands = ``
	for _, role := range roleList {
		if len(role.PermAction) > 0 {
			if _, ok := actChecked[`*`]; !ok {
				for _, pa := range strings.Split(role.PermAction, `,`) {
					if _, ok := actChecked[pa]; !ok {
						actChecked[pa] = struct{}{}
						r.Actions += actSep + pa
						actSep = `,`
					}
				}
			}
		}
		if len(role.PermCmd) > 0 {
			if _, ok := cmdChecked[`*`]; !ok {
				for _, pa := range strings.Split(role.PermCmd, `,`) {
					if _, ok := cmdChecked[pa]; !ok {
						cmdChecked[pa] = struct{}{}
						r.Commands += cmdSep + pa
						cmdSep = `,`
					}
				}
			}
		}
	}
	return r
}

func (r *RolePermission) Check(permPath string) bool {
	permPath = strings.TrimPrefix(permPath, `/`)
	if len(r.Actions) == 0 {
		return perm.NavTreeCached().Check(permPath, nil)
	}
	navTree := perm.NavTreeCached()
	if r.permActions == nil {
		r.permActions = perm.NewMap()
		r.permActions.Parse(r.Actions, navTree)
	}
	return r.permActions.Check(permPath, navTree)
}

func (r *RolePermission) CheckCmd(permPath string) bool {
	if r.permCommonds == nil {
		r.permCommonds = perm.NewMap().ParseCmd(r.Commands)
	}

	return r.permCommonds.CheckCmd(permPath)
}

//FilterNavigate 过滤导航菜单，只显示有权限的菜单
func (r *RolePermission) FilterNavigate(navList *navigate.List) navigate.List {
	var result navigate.List
	if navList == nil {
		return result
	}
	for _, nav := range *navList {
		if !nav.Unlimited && !r.Check(nav.Action) {
			continue
		}
		navCopy := *nav
		navCopy.Children = &navigate.List{}
		for _, child := range *nav.Children {
			var perm string
			if len(child.Action) > 0 {
				perm = nav.Action + `/` + child.Action
			} else {
				perm = nav.Action
			}
			if !nav.Unlimited && !r.Check(perm) {
				continue
			}
			*navCopy.Children = append(*navCopy.Children, child)
		}
		result = append(result, &navCopy)
	}
	return result
}
