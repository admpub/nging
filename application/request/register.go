package request

import (
	"github.com/admpub/nging/v5/application/library/codec"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

type Register struct {
	InvitationCode       string `validate:"required,min=16,max=32"`
	Username             string `validate:"required,username"`
	Email                string `validate:"required,email"`
	Password             string `validate:"required,min=8,max=64"`
	ConfirmationPassword string `validate:"required,eqfield=Password"`
}

func (r *Register) BeforeValidate(ctx echo.Context) error {
	if len(r.Password) == 0 {
		return ctx.NewError(code.InvalidParameter, `请输入密码`).SetZone(`password`)
	}
	if len(r.ConfirmationPassword) == 0 {
		return ctx.NewError(code.InvalidParameter, `请输入确认密码`).SetZone(`confirmationPassword`)
	}
	passwd, err := codec.DefaultSM2DecryptHex(r.Password)
	if err != nil {
		err = ctx.NewError(code.InvalidParameter, `密码解密失败: %v`, err).SetZone(`password`)
	} else {
		r.Password = passwd
	}
	cpasswd, err := codec.DefaultSM2DecryptHex(r.ConfirmationPassword)
	if err != nil {
		err = ctx.NewError(code.InvalidParameter, `密码解密失败: %v`, err).SetZone(`confirmationPassword`)
	} else {
		r.ConfirmationPassword = cpasswd
	}
	return err
}
