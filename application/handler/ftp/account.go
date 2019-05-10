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
package ftp

import (
	"strings"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func AccountIndex(ctx echo.Context) error {
	groupId := ctx.Formx(`groupId`).Uint()
	m := model.NewFtpUser(ctx)
	cond := db.Cond{}
	if groupId > 0 {
		cond[`group_id`] = groupId
	}
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond))
	ret := handler.Err(ctx, err)
	users := m.Objects()
	gIds := []uint{}
	userAndGroup := make([]*model.FtpUserAndGroup, len(users))
	for k, u := range users {
		userAndGroup[k] = &model.FtpUserAndGroup{
			FtpUser: u,
		}
		if u.GroupId < 1 {
			continue
		}
		if !com.InUintSlice(u.GroupId, gIds) {
			gIds = append(gIds, u.GroupId)
		}
	}

	mg := model.NewFtpUserGroup(ctx)
	var groupList []*dbschema.FtpUserGroup
	if len(gIds) > 0 {
		_, err = mg.List(&groupList, nil, 1, 1000, db.Cond{`id IN`: gIds})
		if err != nil {
			if ret == nil {
				ret = err
			}
		} else {
			for k, v := range userAndGroup {
				for _, g := range groupList {
					if g.Id == v.GroupId {
						userAndGroup[k].Group = g
						break
					}
				}
			}
		}
	}
	ctx.Set(`listData`, userAndGroup)
	mg.ListByOffset(&groupList, nil, 0, -1)
	ctx.Set(`groupList`, groupList)
	ctx.Set(`groupId`, groupId)
	return ctx.Render(`ftp/account`, ret)
}

func AccountAdd(ctx echo.Context) error {
	var err error
	m := model.NewFtpUser(ctx)
	if ctx.IsPost() {
		username := ctx.Form(`username`)
		if ctx.Form(`confirmPassword`) != ctx.Form(`password`) {
			err = ctx.E(`两次输入的密码之间不匹配，请输入一样的密码，以确认自己没有输入错误`)
		} else if len(ctx.Form(`password`)) < 6 {
			err = ctx.E(`密码不能少于6个字符`)
		} else if len(username) == 0 {
			err = ctx.E(`账户名不能为空`)
		} else if y, e := m.Exists(username); e != nil {
			err = e
		} else if y {
			err = ctx.E(`账户名已经存在`)
		} else {
			err = ctx.MustBind(m.FtpUser)
		}

		if err == nil {
			m.Password = com.MakePassword(m.Password, model.DefaultSalt)
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
				echo.StructToForm(ctx, m.FtpUser, ``, func(topName, fieldName string) string {
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
	if ctx.IsPost() {
		password := ctx.Form(`password`)
		length := len(password)
		if ctx.Form(`confirmPassword`) != password {
			err = ctx.E(`两次输入的密码之间不匹配，请输入一样的密码，以确认自己没有输入错误`)
		} else if length > 0 && length < 6 {
			err = ctx.E(`密码不能少于6个字符`)
		} else {
			err = ctx.MustBind(m.FtpUser, func(k string, v []string) (string, []string) {
				switch strings.ToLower(k) {
				case `password`:
					if len(v) < 1 || v[0] == `` {
						//忽略密码为空的情况
						return ``, v
					}
					v[0] = com.MakePassword(v[0], model.DefaultSalt)
				case `created`, `username`: //禁止修改创建时间和用户名
					return ``, v
				}
				return k, v
			})
		}

		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/ftp/account`))
			}
		}
	} else if err == nil {
		echo.StructToForm(ctx, m.FtpUser, ``, func(topName, fieldName string) string {
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
			err = ctx.MustBind(m.FtpUserGroup)
		}
		if err == nil {
			_, err = m.Add()
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
				echo.StructToForm(ctx, m.FtpUserGroup, ``, echo.LowerCaseFirstLetter)
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
			err = ctx.MustBind(m.FtpUserGroup, func(k string, v []string) (string, []string) {
				switch strings.ToLower(k) {
				case `created`: //禁止修改创建时间
					return ``, v
				}
				return k, v
			})
		}

		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/ftp/group`))
			}
		}
	} else if err == nil {
		echo.StructToForm(ctx, m.FtpUserGroup, ``, echo.LowerCaseFirstLetter)
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
