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

package handler

import (
	"time"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/formfilter"

	"github.com/admpub/nging/v4/application/handler"

	"github.com/nging-plugins/frpmanager/application/model"
)

func AccountIndex(ctx echo.Context) error {
	m := model.NewFrpUser(ctx)
	cond := db.NewCompounds()
	list := []*model.FrpUserAndServer{}
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, &list, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()))
	ctx.Set(`listData`, list)
	return ctx.Render(`frp/account`, handler.Err(ctx, err))
}

func accountFormFilter(opts ...formfilter.Options) echo.FormDataFilter {
	opts = append(opts,
		formfilter.StartDateToTimestamp(`Start`),
		formfilter.EndDateToTimestamp(`End`),
	)
	return formfilter.Build(opts...)
}

func AccountAdd(ctx echo.Context) error {
	var err error
	m := model.NewFrpUser(ctx)
	user := handler.User(ctx)
	if ctx.IsPost() {
		if ctx.Form(`confirmPassword`) != ctx.Form(`password`) {
			err = ctx.E(`两次输入的密码之间不匹配，请输入一样的密码`)
		} else if len(ctx.Form(`password`)) < 6 {
			err = ctx.E(`密码不能少于6个字符`)
		} else {
			err = ctx.MustBind(m.NgingFrpUser, accountFormFilter())
		}
		if err == nil {
			m.Uid = user.Id
			_, err = m.Add()
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/frp/account`))
			}
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.NgingFrpUser, ``, func(topName, fieldName string) string {
					if topName == `` && fieldName == `Password` {
						return ``
					}
					return echo.LowerCaseFirstLetter(topName, fieldName)
				})
				ctx.Request().Form().Set(`id`, `0`)
				var startDate, endDate string
				if m.Start > 0 {
					startDate = time.Unix(int64(m.Start), 0).Format(`2006-01-02`)
				}
				ctx.Request().Form().Set(`start`, startDate)
				if m.End > 0 {
					endDate = time.Unix(int64(m.End), 0).Format(`2006-01-02`)
				}
				ctx.Request().Form().Set(`end`, endDate)
			}
		}
	}
	ctx.Set(`activeURL`, `/frp/account`)
	return ctx.Render(`frp/account_edit`, err)
}

func AccountEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint64()
	m := model.NewFrpUser(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.NewError(code.DataNotFound, `数据不存在`)
		}
		return err
	}
	if ctx.IsPost() {
		password := ctx.Form(`password`)
		length := len(password)
		if ctx.Form(`confirmPassword`) != password {
			err = ctx.E(`两次输入的密码之间不匹配，请输入一样的密码`)
		} else if length > 0 && length < 6 {
			err = ctx.E(`密码不能少于6个字符`)
		} else {
			err = ctx.MustBind(m.NgingFrpUser, accountFormFilter(formfilter.Exclude(`Created`)))
		}

		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/frp/account`))
			}
		}
	} else {
		echo.StructToForm(ctx, m.NgingFrpUser, ``, func(topName, fieldName string) string {
			if topName == `` && fieldName == `Password` {
				return ``
			}
			return echo.LowerCaseFirstLetter(topName, fieldName)
		})
		var startDate, endDate string
		if m.Start > 0 {
			startDate = time.Unix(int64(m.Start), 0).Format(`2006-01-02`)
		}
		ctx.Request().Form().Set(`start`, startDate)
		if m.End > 0 {
			endDate = time.Unix(int64(m.End), 0).Format(`2006-01-02`)
		}
		ctx.Request().Form().Set(`end`, endDate)
	}

	ctx.Set(`activeURL`, `/frp/account`)
	return ctx.Render(`frp/account_edit`, err)
}

func AccountDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewFrpUser(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/frp/account`))
}
