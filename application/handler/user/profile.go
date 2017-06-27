/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package user

import (
	"errors"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func init() {
	handler.RegisterToGroup(`/user`, func(g *echo.Group) {
		g.Route("GET,POST", `/edit`, Edit)
	})
}

func Edit(ctx echo.Context) error {
	var err error
	user := handler.User(ctx)
	if user == nil {
		return errors.New(ctx.T(`登录信息获取失败，请重新登录`))
	}
	m := model.NewUser(ctx)
	err = m.Get(nil, `id`, user.Id)
	if err != nil {
		return err
	}
	needCheckU2F := m.NeedCheckU2F(user.Id)
	if ctx.IsPost() {
		email := ctx.Form(`email`)
		modifyPass := ctx.Form(`modifyPass`) == `1`

		//新密码
		newPass := ctx.Form(`newPass`)
		confirmPass := ctx.Form(`confirmPass`)

		//旧密码
		passwd := ctx.Form(`pass`)

		if len(email) == 0 {
			err = errors.New(ctx.T(`Email不能为空`))
		} else if modifyPass && len(newPass) < 8 {
			err = errors.New(ctx.T(`新密码不能少于8个字符`))
		} else if modifyPass && newPass != confirmPass {
			err = errors.New(ctx.T(`新密码与确认新密码不一致`))
		} else if m.User.Password != com.MakePassword(passwd, m.User.Salt) {
			err = errors.New(ctx.T(`旧密码输入不正确`))
		} else if needCheckU2F {
			//两步验证码
			err = GAuthVerify(ctx, `u2fCode`)
		}
		if err == nil && !ctx.Validate(`email`, email, `email`).Ok() {
			err = errors.New(ctx.T(`Email地址格式不正确`))
		}
		if err == nil {
			set := map[string]interface{}{
				`email`: email,
			}
			if modifyPass {
				set[`password`] = com.MakePassword(newPass, m.User.Salt)
			}
			err = m.Param().SetSend(set).SetArgs(`id`, user.Id).Update()
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`修改成功`))
			m.Get(nil, `id`, user.Id)
			m.SetSession()
			return ctx.Redirect(`/user/edit`)
		}
	}
	ctx.Set(`needCheckU2F`, needCheckU2F)
	return ctx.Render(`user/edit`, handler.Err(ctx, err))
}
