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
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/formfilter"

	"github.com/admpub/nging/v4/application/handler"

	"github.com/nging-plugins/ftpmanager/application/dbschema"
	"github.com/nging-plugins/ftpmanager/application/model"
)

func AccountIndex(ctx echo.Context) error {
	groupId := ctx.Formx(`groupId`).Uint()
	m := model.NewFtpUser(ctx)
	cond := db.Cond{}
	if groupId > 0 {
		cond[`group_id`] = groupId
	}
	var userAndGroup []*model.FtpUserAndGroup
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, &userAndGroup, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond))
	mg := model.NewFtpUserGroup(ctx)
	var groupList []*dbschema.NgingFtpUserGroup
	mg.ListByOffset(&groupList, nil, 0, -1)
	ctx.Set(`listData`, userAndGroup)
	ctx.Set(`groupList`, groupList)
	ctx.Set(`groupId`, groupId)
	return ctx.Render(`ftp/account`, handler.Err(ctx, err))
}

func AccountAdd(ctx echo.Context) error {
	var err error
	m := model.NewFtpUser(ctx)
	if ctx.IsPost() {
		if ctx.Form(`confirmPassword`) != ctx.Form(`password`) {
			err = ctx.E(`两次输入密码不匹配，请输入一样的密码，以便确认自己没有输入错误`)
		} else if len(ctx.Form(`password`)) < 6 {
			err = ctx.E(`密码不能少于6个字符`)
		} else {
			err = ctx.MustBind(m.NgingFtpUser)
		}
		if err == nil {
			_, err = m.Add()
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/ftp/account`))
			}
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.NgingFtpUser, ``, func(topName, fieldName string) string {
					if topName == `` && fieldName == `Password` {
						return ``
					}
					return echo.LowerCaseFirstLetter(topName, fieldName)
				})
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}
	mg := model.NewFtpUserGroup(ctx)
	_, groupList, e := mg.ListByActive(1, 1000)
	if err == nil {
		err = e
	}
	ctx.Set(`groupList`, groupList)
	return ctx.Render(`ftp/account_edit`, err)
}

func AccountEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := model.NewFtpUser(ctx)
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
			err = ctx.E(`两次输入密码不匹配，请输入一样的密码，以便确认自己没有输入错误`)
		} else if length > 0 && length < 6 {
			err = ctx.E(`密码不能少于6个字符`)
		} else {
			err = ctx.MustBind(m.NgingFtpUser, formfilter.Build(formfilter.Exclude(`created`)))
		}
		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/ftp/account`))
			}
		}
	} else {
		echo.StructToForm(ctx, m.NgingFtpUser, ``, func(topName, fieldName string) string {
			if topName == `` && fieldName == `Password` {
				return ``
			}
			return echo.LowerCaseFirstLetter(topName, fieldName)
		})
	}

	mg := model.NewFtpUserGroup(ctx)
	_, groupList, e := mg.ListByActive(1, 1000)
	if err == nil {
		err = e
	}
	ctx.Set(`groupList`, groupList)
	ctx.Set(`activeURL`, `/ftp/account`)
	return ctx.Render(`ftp/account_edit`, err)
}

func AccountDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewFtpUser(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/ftp/account`))
}

func GroupIndex(ctx echo.Context) error {
	m := model.NewFtpUserGroup(ctx)
	_, err := handler.PagingWithLister(ctx, m)
	ret := handler.Err(ctx, err)
	ctx.Set(`listData`, m.Objects())
	return ctx.Render(`ftp/group`, ret)
}

func GroupAdd(ctx echo.Context) error {
	var err error
	m := model.NewFtpUserGroup(ctx)
	if ctx.IsPost() {
		name := ctx.Form(`name`)
		if len(name) == 0 {
			err = ctx.E(`用户组名称不能为空`)
		} else if y, e := m.Exists(name); e != nil {
			err = e
		} else if y {
			err = ctx.E(`用户组名称已经存在`)
		} else {
			err = ctx.MustBind(m.NgingFtpUserGroup)
		}
		if err == nil {
			_, err = m.Insert()
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/ftp/group`))
			}
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.NgingFtpUserGroup, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}

	return ctx.Render(`ftp/group_edit`, err)
}

func GroupEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := model.NewFtpUserGroup(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if ctx.IsPost() {
		name := ctx.Form(`name`)
		if len(name) < 1 {
			err = ctx.E(`用户组名称不能为空`)
		} else if y, e := m.ExistsOther(name, id); e != nil {
			err = e
		} else if y {
			err = ctx.E(`用户组名称已经存在`)
		} else {
			err = ctx.MustBind(m.NgingFtpUserGroup, echo.ExcludeFieldName(`created`))
		}

		if err == nil {
			m.Id = id
			err = m.Update(nil, db.Cond{`id`: id})
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/ftp/group`))
			}
		}
	} else if err == nil {
		echo.StructToForm(ctx, m.NgingFtpUserGroup, ``, echo.LowerCaseFirstLetter)
	}

	ctx.Set(`activeURL`, `/ftp/group`)
	return ctx.Render(`ftp/group_edit`, err)
}

func GroupDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewFtpUserGroup(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/ftp/group`))
}
