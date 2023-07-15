package user

import (
	"strings"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/codec"
	"github.com/admpub/nging/v5/application/model"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func EditPassword(ctx echo.Context) error {
	var err error
	user := handler.User(ctx)
	if user == nil {
		return ctx.NewError(code.Unauthenticated, `登录信息获取失败，请重新登录`)
	}
	m := model.NewUser(ctx)
	err = m.Get(nil, `id`, user.Id)
	if err != nil {
		return err
	}
	needCheckU2F, err := m.NeedCheckU2F(model.AuthTypePassword, user.Id, 2)
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		//新密码
		newPass := strings.TrimSpace(ctx.Form(`newPass`))
		confirmPass := strings.TrimSpace(ctx.Form(`confirmPass`))

		//旧密码
		passwd := strings.TrimSpace(ctx.Form(`pass`))

		passwd, err = codec.DefaultSM2DecryptHex(passwd)
		if err != nil {
			return ctx.NewError(code.InvalidParameter, `旧密码解密失败: %v`, err).SetZone(`pass`)
		}
		newPass, err = codec.DefaultSM2DecryptHex(newPass)
		if err != nil {
			return ctx.NewError(code.InvalidParameter, `新密码解密失败: %v`, err).SetZone(`newPass`)
		}
		confirmPass, err = codec.DefaultSM2DecryptHex(confirmPass)
		if err != nil {
			return ctx.NewError(code.InvalidParameter, `您输入的确认密码解密失败: %v`, err).SetZone(`confirmPass`)
		}

		if len(newPass) < 8 {
			err = ctx.E(`新密码不能少于8个字符`)
		} else if newPass != confirmPass {
			err = ctx.E(`新密码与确认新密码不一致`)
		} else if m.NgingUser.Password != com.MakePassword(passwd, m.NgingUser.Salt) {
			err = ctx.E(`旧密码输入不正确`)
		} else if needCheckU2F {
			//两步验证码
			err = GAuthVerify(ctx, `u2fCode`)
		}
		if err == nil {
			set := echo.H{
				`password`: com.MakePassword(newPass, m.NgingUser.Salt),
			}
			err = m.UpdateFields(nil, set, `id`, user.Id)
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`修改成功`))
			return ctx.Redirect(handler.URLFor(`/user/password`))
		}
	}
	ctx.Set(`needCheckU2F`, needCheckU2F)
	ctx.Set(`activeSafeItem`, `password`)
	ctx.Set(`safeItems`, model.SafeItems.Slice())
	return ctx.Render(`user/password`, handler.Err(ctx, err))
}
