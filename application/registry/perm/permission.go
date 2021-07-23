package perm

import (
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/nging/v3/application/dbschema"
	"github.com/admpub/nging/v3/application/library/perm"
	"github.com/admpub/nging/v3/application/registry/navigate"
)

func New() *RolePermission {
	return &RolePermission{}
}

type RolePermission struct {
	Actions       string
	Commands      string
	Behaviors     string
	Roles         []*dbschema.NgingUserRole
	permActions   *perm.Map
	permCommonds  *perm.Map
	permBehaviors perm.BehaviorPerms
	filter        *navigate.Filter
}

func (r *RolePermission) Init(roleList []*dbschema.NgingUserRole) *RolePermission {
	if r.filter == nil {
		r.filter = navigate.NewFilter(r)
	}
	r.Roles = roleList
	cmdChecked := map[string]struct{}{}
	actChecked := map[string]struct{}{}
	behChecked := map[string]struct{}{}
	var cmdSep, actSep, behSep string
	r.Actions = ``
	r.Commands = ``
	r.Behaviors = ``
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
		if len(role.PermBehavior) > 0 {
			if _, ok := behChecked[`*`]; !ok {
				for _, pa := range strings.Split(role.PermBehavior, `,`) {
					if _, ok := behChecked[pa]; !ok {
						behChecked[pa] = struct{}{}
						r.Behaviors += behSep + pa
						behSep = `,`
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

func (r *RolePermission) CheckBehavior(permPath string) *perm.CheckedBehavior {
	if r.permBehaviors == nil {
		var err error
		r.permBehaviors, err = perm.ParseBehavior(r.Behaviors, Behaviors)
		if err != nil {
			log.Error(err)
		}
	}

	return r.permBehaviors.CheckBehavior(permPath)
}

//FilterNavigate 过滤导航菜单，只显示有权限的菜单
func (r *RolePermission) FilterNavigate(navList *navigate.List) navigate.List {
	return r.filter.FilterNavigate(navList)
}
