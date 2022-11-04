package request

import (
	"github.com/admpub/nging/v5/application/library/codec"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

type Login struct {
	User string `validate:"required,username"`
	Pass string `validate:"required,min=8,max=64"`
	Code string `validate:"required"`
}

func (r *Login) BeforeValidate(ctx echo.Context) error {
	if len(r.Pass) == 0 {
		return ctx.NewError(code.InvalidParameter, `请输入密码`).SetZone(`password`)
	}
	passwd, err := codec.DefaultSM2DecryptHex(r.Pass)
	if err != nil {
		err = ctx.NewError(code.InvalidParameter, `密码解密失败: %v`, err).SetZone(`password`)
	} else {
		r.Pass = passwd
	}
	return err
}
