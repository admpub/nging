/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package index

import (
	"fmt"

	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/codec"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/library/license"
	"github.com/admpub/nging/v4/application/middleware"
	"github.com/admpub/nging/v4/application/model"
	"github.com/admpub/nging/v4/application/registry/dashboard"
	stdCode "github.com/webx-top/echo/code"
)

func Index(ctx echo.Context) error {
	user := handler.User(ctx)
	if user == nil {
		return ctx.Redirect(handler.URLFor(`/login`))
	}
	var err error
	ctx.Set(`cards`, dashboard.CardAll(ctx).Build(ctx))
	blocks := dashboard.BlockAll(ctx)
	if err = blocks.Ready(ctx); err != nil {
		return err
	}
	ctx.Set(`blocks`, blocks)
	ctx.Set(`license`, license.License())
	ctx.Set(`showExpirationTime`, config.DefaultConfig.Sys.ShowExpirationTime)
	productURL := license.ProductDetailURL()
	ctx.Set(`productURL`, productURL)
	return ctx.Render(`index`, handler.Err(ctx, err))
}

func Login(ctx echo.Context) error {
	//检查是否已安装
	if !config.IsInstalled() {
		return ctx.Redirect(handler.URLFor(`/setup`))
	}

	next := ctx.Form(`next`)
	if len(next) == 0 {
		next = handler.URLFor(`/index`)
	}

	user := handler.User(ctx)
	if user != nil {
		return ctx.Redirect(next)
	}
	var err error
	if ctx.IsPost() {
		if data := common.VerifyCaptcha(ctx, `backend`, `code`); data.GetCode() == stdCode.CaptchaError {
			err = common.ErrCaptcha
		} else if data.GetCode() != stdCode.Success {
			err = fmt.Errorf("%v", data.GetInfo())
		} else {
			err = middleware.Auth(ctx, true)
			if err == nil {
				return ctx.Redirect(next)
			}
		}
	}
	return ctx.Render(`login`, handler.Err(ctx, err))
}

func Register(ctx echo.Context) error {
	var err error
	if ctx.IsPost() {
		c := model.NewCode(ctx)
		m := model.NewUser(ctx)
		code := ctx.Form(`invitationCode`)
		user := ctx.Form(`username`)
		email := ctx.Form(`email`)
		passwd := ctx.Form(`password`)
		repwd := ctx.Form(`confirmationPassword`)
		passwd, err = codec.DefaultSM2DecryptHex(passwd)
		if err != nil {
			err = ctx.NewError(stdCode.InvalidParameter, ctx.T(`密码解密失败: %v`, err)).SetZone(`password`)
			goto END
		}
		repwd, err = codec.DefaultSM2DecryptHex(repwd)
		if err != nil {
			err = ctx.NewError(stdCode.InvalidParameter, ctx.T(`您输入的确认密码解密失败: %v`, err)).SetZone(`confirmationPassword`)
			goto END
		}
		if len(code) == 0 {
			err = ctx.E(`邀请码不能为空`)
			goto END
		}
		if repwd != passwd {
			err = ctx.E(`密码与确认密码不一致`)
			goto END
		}
		err = c.VerfyInvitationCode(code)
		if err != nil {
			goto END
		}
		err = m.Register(user, passwd, email, c.Invitation.RoleIds)
		if err != nil {
			goto END
		}
		c.UseInvitationCode(c.Invitation, m.NgingUser.Id)
		m.SetSession()
		m.SetField(nil, `session_id`, ctx.Session().ID(), `id`, m.NgingUser.Id)
		next := ctx.Query(`next`)
		if len(next) == 0 {
			next = handler.URLFor(`/index`)
		}
		return ctx.Redirect(next)
	}

END:
	return ctx.Render(`register`, handler.Err(ctx, err))
}

func Logout(ctx echo.Context) error {
	ctx.Session().Delete(`user`)
	return ctx.Redirect(handler.URLFor(`/login`))
}

var (
	DonationAccountTypes       = echo.NewKVData()
	DefaultDonationAccountType = `alipay`
)

func init() {
	DonationAccountTypes.AddItem(&echo.KV{
		K: `alipay`, V: `支付宝付款`, H: echo.H{`qrimg`: `alipay.jpg`},
	})
	DonationAccountTypes.AddItem(&echo.KV{
		K: `wechat`, V: `微信支付`, H: echo.H{`qrimg`: `wechat.png`},
	})
	DonationAccountTypes.AddItem(&echo.KV{
		K: `btcoin`, V: `比特币支付`, H: echo.H{`qrimg`: `btcoin.jpeg`},
	})
}

func Donation(ctx echo.Context) error {
	if len(DonationAccountTypes.Slice()) == 0 {
		return echo.ErrNotFound
	}
	typ := ctx.Param(`type`, DefaultDonationAccountType)
	item := DonationAccountTypes.GetItem(typ)
	if item == nil {
		typ = DefaultDonationAccountType
		item = DonationAccountTypes.GetItem(typ)
		if item == nil {
			return echo.ErrNotFound
		}
	}
	ctx.Set(`data`, item)
	ctx.Set(`list`, DonationAccountTypes.Slice())
	return ctx.Render(`index/donation/donation`, nil)
}
