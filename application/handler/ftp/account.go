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
package ftp

import (
	"errors"
	"strings"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func init() {
	handler.RegisterToGroup(`/ftp`, func(g *echo.Group) {
		g.Route(`GET`, `/account`, AccountIndex)
		g.Route(`GET,POST`, `/account_add`, AccountAdd)
		g.Route(`GET,POST`, `/account_edit`, AccountEdit)
		g.Route(`GET,POST`, `/account_delete`, AccountDelete)

		g.Route(`GET`, `/group`, GroupIndex)
		g.Route(`GET,POST`, `/group_add`, GroupAdd)
		g.Route(`GET,POST`, `/group_edit`, GroupEdit)
		g.Route(`GET,POST`, `/group_delete`, GroupDelete)
	})
}

func AccountIndex(ctx echo.Context) error {
	m := model.NewFtpUser(ctx)
	page, size, totalRows, p := handler.PagingWithPagination(ctx)
	cnt, err := m.List(nil, nil, page, size)
	ret := handler.Err(ctx, err)
	if totalRows <= 0 {
		totalRows = int(cnt())
		p.SetRows(totalRows)
	}
	ctx.Set(`pagination`, p)
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
		has := false
		for _, gid := range gIds {
			if gid == u.GroupId {
				has = true
				break
			}
		}
		if !has {
			gIds = append(gIds, u.GroupId)
		}
	}

	mg := model.NewFtpUserGroup(ctx)
	var groupList []*dbschema.FtpUserGroup
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
	ctx.Set(`listData`, userAndGroup)
	return ctx.Render(`ftp/account`, ret)
}

func AccountAdd(ctx echo.Context) error {
	var err error
	if ctx.IsPost() {
		m := model.NewFtpUser(ctx)
		username := ctx.Form(`username`)
		if ctx.Form(`confirmPassword`) != ctx.Form(`password`) {
			err = errors.New(ctx.T(`两次输入的密码之间不匹配，请输入一样的密码，以确认自己没有输入错误`))
		} else if len(ctx.Form(`password`)) < 6 {
			err = errors.New(ctx.T(`密码不能少于6个字符`))
		} else if len(username) == 0 {
			err = errors.New(ctx.T(`账户名不能为空`))
		} else if y, e := m.Exists(username); e != nil {
			err = e
		} else if y {
			err = errors.New(ctx.T(`账户名已经存在`))
		} else {
			err = ctx.MustBind(m.FtpUser)
		}

		if err == nil {
			m.Password = com.MakePassword(m.Password, model.DefaultSalt)
			_, err = m.Add()
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(`/ftp/account`)
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
			err = errors.New(ctx.T(`两次输入的密码之间不匹配，请输入一样的密码，以确认自己没有输入错误`))
		} else if length > 0 && length < 6 {
			err = errors.New(ctx.T(`密码不能少于6个字符`))
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
				return ctx.Redirect(`/ftp/account`)
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

	return ctx.Redirect(`/ftp/account`)
}

func GroupIndex(ctx echo.Context) error {
	m := model.NewFtpUserGroup(ctx)
	page, size, totalRows, p := handler.PagingWithPagination(ctx)
	cnt, err := m.List(nil, nil, page, size)
	ret := handler.Err(ctx, err)
	if totalRows <= 0 {
		totalRows = int(cnt())
		p.SetRows(totalRows)
	}
	ctx.Set(`pagination`, p)
	ctx.Set(`listData`, m.Objects())
	return ctx.Render(`ftp/group`, ret)
}

func GroupAdd(ctx echo.Context) error {
	var err error
	if ctx.IsPost() {
		m := model.NewFtpUserGroup(ctx)
		name := ctx.Form(`name`)
		if len(name) == 0 {
			err = errors.New(ctx.T(`用户组名称不能为空`))
		} else if y, e := m.Exists(name); e != nil {
			err = e
		} else if y {
			err = errors.New(ctx.T(`用户组名称已经存在`))
		} else {
			err = ctx.MustBind(m.FtpUserGroup)
		}
		if err == nil {
			_, err = m.Add()
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(`/ftp/group`)
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
		if len(name) < 6 {
			err = errors.New(ctx.T(`用户组名称不能为空`))
		} else if y, e := m.ExistsOther(name, id); e != nil {
			err = e
		} else if y {
			err = errors.New(ctx.T(`用户组名称已经存在`))
		} else {
			err = ctx.MustBind(m.FtpUserGroup, func(k string, v []string) (string, []string) {
				switch k {
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
				return ctx.Redirect(`/ftp/group`)
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

	return ctx.Redirect(`/ftp/group`)
}
