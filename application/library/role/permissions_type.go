package role

import (
	"github.com/admpub/nging/v4/application/library/perm"
	"github.com/webx-top/echo"
)

var UserRolePermissionType = echo.NewKVData().
	Add(RolePermissionTypePage, `页面权限`, echo.KVOptX(
		perm.NewHandle().SetTmpl(`/manager/role_edit_perm_page`).SetTmpl(`/manager/role_edit_perm_page_foot`, `foot`).
			SetGenerator(PermPageGenerator).
			SetParser(PermPageParser).
			SetChecker(PermPageChecker).
			SetItemLister(PermPageList).
			OnRender(PermPageOnRender),
	)).
	Add(RolePermissionTypeCommand, `指令集权限`, echo.KVOptX(
		perm.NewHandle().SetTmpl(`/manager/role_edit_perm_command`).
			SetGenerator(PermCommandGenerator).
			SetParser(PermCommandParser).
			SetChecker(PermCommandChecker).
			SetItemLister(PermCommandList).
			OnRender(PermCommandOnRender).
			SetIsValid(PermCommandIsValid),
	)).
	Add(RolePermissionTypeBehavior, `行为权限`, echo.KVOptX(
		perm.NewHandle().SetTmpl(`/manager/role_edit_perm_behavior`).
			SetGenerator(PermBehaviorGenerator).
			SetParser(PermBehaviorParser).
			SetChecker(PermBehaviorChecker).
			SetItemLister(PermBehaviorList).
			OnRender(PermBehaviorOnRender).
			SetIsValid(PermBehaviorIsValid),
	))

const (
	RolePermissionTypePage     = `page`
	RolePermissionTypeCommand  = `command`
	RolePermissionTypeBehavior = `behavior`
)
