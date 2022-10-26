package model

import (
	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/code"
)

func (u *User) Register(user, pass, email, roleIds string) error {
	ctx := u.Context()
	if len(user) == 0 {
		return ctx.NewError(code.InvalidParameter, `用户名不能为空`).SetZone(`username`)
	}
	if len(email) == 0 {
		return ctx.NewError(code.InvalidParameter, `Email不能为空`).SetZone(`email`)
	}
	if len(pass) < 8 {
		return ctx.NewError(code.InvalidParameter, `密码不能少于8个字符`).SetZone(`password`)
	}
	if !com.IsUsername(user) {
		return ctx.NewError(code.InvalidParameter, `用户名不能包含特殊字符(只能由字母、数字、下划线和汉字组成)`).SetZone(`username`)
	}
	if !ctx.Validate(`email`, email, `email`).Ok() {
		return ctx.NewError(code.InvalidParameter, `Email地址格式不正确`).SetZone(`email`)
	}
	exists, err := u.Exists(user)
	if err != nil {
		return err
	}
	if exists {
		return ctx.NewError(code.InvalidParameter, `用户名已经存在`).SetZone(`username`)
	}
	userSchema := dbschema.NewNgingUser(ctx)
	userSchema.Username = user
	userSchema.Email = email
	userSchema.Salt = com.Salt()
	userSchema.Password = com.MakePassword(pass, userSchema.Salt)
	userSchema.Disabled = `N`
	userSchema.RoleIds = roleIds
	_, err = userSchema.EventOFF().Insert()
	u.NgingUser = userSchema
	return err
}
