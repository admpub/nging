package roleutils

import (
	"github.com/admpub/nging/v4/application/library/perm"
	"github.com/admpub/nging/v4/application/model"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

var UserRolePermissionType = echo.NewKVData().
	Add(UserRolePermissionTypePage, `页面权限`, echo.KVOptX(
		perm.NewHandle().SetTmpl(`/manager/role_edit_perm_page`).SetTmpl(`/manager/role_edit_perm_page_foot`, `foot`).
			SetGenerator(PermPageGenerator).
			SetChecker(PermPageChecker).
			SetItemLister(PermPageList).
			OnRender(PermPageOnRender),
	)).
	Add(UserRolePermissionTypeCommand, `指令集权限`, echo.KVOptX(
		perm.NewHandle().SetTmpl(`/manager/role_edit_perm_command`).
			SetGenerator(PermPageGenerator).
			SetChecker(PermCommandChecker).
			SetItemLister(PermCommandList).
			OnRender(PermCommandOnRender).
			SetIsValid(PermCommandIsValid),
	)).
	Add(UserRolePermissionTypeBehavior, `行为权限`, echo.KVOptX(
		perm.NewHandle().SetTmpl(`/manager/role_edit_perm_behavior`).
			SetGenerator(PermBehaviorGenerator).
			SetChecker(PermBehaviorChecker).
			SetItemLister(PermBehaviorList).
			OnRender(PermBehaviorOnRender).
			SetIsValid(PermBehaviorIsValid),
	))

const (
	UserRolePermissionTypePage     = `page`
	UserRolePermissionTypeCommand  = `command`
	UserRolePermissionTypeBehavior = `behavior`
)

func UserRolePermissionTypeFireRender(ctx echo.Context) (err error) {
	return perm.HandleFireRender(ctx, UserRolePermissionType)
}

func UserRolePermissionTypeGenerate(ctx echo.Context) (mp map[string]string, err error) {
	return perm.HandleGenerate(ctx, UserRolePermissionType)
}

func UserRolePermissionTypeCheck(ctx echo.Context, current string, typ string, permission string, parsed interface{}) (mp interface{}, err error) {
	return perm.HandleCheck(ctx, UserRolePermissionType, current, typ, permission, parsed)
}

func AddUserRolePermission(ctx echo.Context, roleID uint) (err error) {
	var perms map[string]string
	perms, err = UserRolePermissionTypeGenerate(ctx)
	if err != nil {
		return
	}
	rpM := model.NewUserRolePermission(ctx)
	for typ, perm := range perms {
		rpM.RoleId = roleID
		rpM.Type = typ
		rpM.Permission = perm
		_, err = rpM.Add()
		if err != nil {
			break
		}
	}
	return
}

func EditUserRolePermission(ctx echo.Context, roleID uint) (err error) {
	var perms map[string]string
	perms, err = UserRolePermissionTypeGenerate(ctx)
	if err != nil {
		return
	}
	rpM := model.NewUserRolePermission(ctx)
	_, err = rpM.ListByOffset(nil, nil, 0, -1, db.Cond{`role_id`: roleID})
	if err != nil {
		return
	}
	var deleted []string
	for _, rule := range rpM.Objects() {
		_, ok := perms[rule.Type]
		if !ok {
			deleted = append(deleted, rule.Type)
		}
	}
	if len(deleted) > 0 {
		err = rpM.Delete(nil, db.And(
			db.Cond{`role_id`: roleID},
			db.Cond{`type`: db.In(deleted)},
		))
		if err != nil {
			return
		}
	}
	for typ, perm := range perms {
		rpM.RoleId = roleID
		rpM.Type = typ
		rpM.Permission = perm
		_, err = rpM.Add()
		if err != nil {
			break
		}
	}
	return
}
