package role

import (
	"testing"

	"github.com/admpub/nging/v5/application/library/perm"
	"github.com/stretchr/testify/assert"
)

var _ PermissionsGetter = &testPermissions{}

type testPermissions struct {
	permissions []PermissionConfiger
}

func (a *testPermissions) GetPermissions() []PermissionConfiger {
	return a.permissions
}

type testPermission struct {
	Type       string
	Permission string
}

func (a *testPermission) GetType() string {
	return a.Type
}

func (a *testPermission) GetPermission() string {
	return a.Permission
}

type article struct {
	MaxPerDay        int64 `json:"maxPerDay"`
	MaxPending       int64 `json:"maxPending"`
	MaxPendingPerDay int64 `json:"maxPendingPerDay"`
}

type articleCfgWithCombine article

func (a *articleCfgWithCombine) Combine(source interface{}) interface{} {
	s := source.(*articleCfgWithCombine)
	if s.MaxPerDay > a.MaxPerDay {
		a.MaxPerDay = s.MaxPerDay
	}
	if s.MaxPending > a.MaxPending {
		a.MaxPending = s.MaxPending
	}
	if s.MaxPendingPerDay > a.MaxPendingPerDay {
		a.MaxPendingPerDay = s.MaxPendingPerDay
	}
	return a
}

func init() {
	Behaviors.Register(`article`, `文章投稿设置`,
		perm.BehaviorOptValue(&article{}),
		perm.BehaviorOptValueInitor(func() interface{} {
			return &article{}
		}),
		perm.BehaviorOptValueType(`json`),
	)
	Behaviors.Register(`article2`, `文章投稿设置`,
		perm.BehaviorOptValue(&articleCfgWithCombine{}),
		perm.BehaviorOptValueInitor(func() interface{} {
			return &articleCfgWithCombine{}
		}),
		perm.BehaviorOptValueType(`json`),
	)
}

func TestInit(t *testing.T) {
	a := NewCommonPermission(UserRolePermissionType, nil)
	expected := `{"article":{"maxPerDay":0,"maxPending":100,"maxPendingPerDay":10}}`
	ps := []PermissionsGetter{
		&testPermissions{
			permissions: []PermissionConfiger{
				&testPermission{Type: RolePermissionTypeBehavior, Permission: expected},
			},
		},
	}
	a.Init(ps)
	assert.Equal(t, expected, string(a.Combined[RolePermissionTypeBehavior]))

	expected = `{"article":{"maxPerDay":10,"maxPending":200,"maxPendingPerDay":20}}`
	a.Init([]PermissionsGetter{&testPermissions{
		permissions: []PermissionConfiger{
			&testPermission{Type: RolePermissionTypeBehavior, Permission: expected},
		},
	}})
	assert.Equal(t, expected, string(a.Combined[RolePermissionTypeBehavior]))

	expected = `{"article":{"maxPerDay":30,"maxPending":150,"maxPendingPerDay":15}}`
	a.Init([]PermissionsGetter{&testPermissions{
		permissions: []PermissionConfiger{
			&testPermission{Type: RolePermissionTypeBehavior, Permission: expected},
		},
	}})
	assert.Equal(t, expected, string(a.Combined[RolePermissionTypeBehavior]))

}

// ------- with combiner -------

func TestWithCombiner(t *testing.T) {
	a := NewCommonPermission(UserRolePermissionType, nil)
	expected := `{"article2":{"maxPerDay":30,"maxPending":150,"maxPendingPerDay":20}}`
	a.Init([]PermissionsGetter{&testPermissions{
		permissions: []PermissionConfiger{
			&testPermission{Type: RolePermissionTypeBehavior, Permission: `{"article2":{"maxPerDay":30,"maxPending":150,"maxPendingPerDay":15}}`},
		},
	}, &testPermissions{
		permissions: []PermissionConfiger{
			&testPermission{Type: RolePermissionTypeBehavior, Permission: `{"article2":{"maxPerDay":10,"maxPending":100,"maxPendingPerDay":10}}`},
		},
	}, &testPermissions{
		permissions: []PermissionConfiger{
			&testPermission{Type: RolePermissionTypeBehavior, Permission: `{"article2":{"maxPerDay":10,"maxPending":100,"maxPendingPerDay":20}}`},
		},
	}})
	//echo.Dump(a.Combined)
	assert.Equal(t, expected, string(a.Combined[RolePermissionTypeBehavior]))
}
