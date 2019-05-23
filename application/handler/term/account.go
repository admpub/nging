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
package term

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func AccountIndex(ctx echo.Context) error {
	groupId := ctx.Formx(`groupId`).Uint()
	m := model.NewSshUser(ctx)
	cond := db.Compounds{}
	if groupId > 0 {
		cond.AddKV(`group_id`, groupId)
	}
	q := ctx.Formx(`q`).String()
	if len(q) > 0 {
		cond.AddKV(`name`, db.Like(`%`+q+`%`))
	}
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()))
	ret := handler.Err(ctx, err)
	users := m.Objects()
	gIds := []uint{}
	userAndGroup := make([]*model.SshUserAndGroup, len(users))
	for k, u := range users {
		userAndGroup[k] = &model.SshUserAndGroup{
			SshUser: u,
		}
		if u.GroupId < 1 {
			continue
		}
		if !com.InUintSlice(u.GroupId, gIds) {
			gIds = append(gIds, u.GroupId)
		}
	}

	mg := model.NewSshUserGroup(ctx)
	var groupList []*dbschema.SshUserGroup
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
	return ctx.Render(`term/account`, ret)
}

func AccountAdd(ctx echo.Context) error {
	var err error
	m := model.NewSshUser(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.SshUser, func(k string, v []string) (string, []string) {
			switch strings.ToLower(k) {
			case `password`, `passphrase`:
				if len(v) > 0 && len(v[0]) > 0 {
					v[0] = config.DefaultConfig.Encode(v[0])
				}
				return k, v
			}
			return k, v
		})
		if err == nil {
			_, err = m.Add()
			if err == nil {
				if ctx.IsAjax() {
					data := ctx.Data().SetInfo(ctx.T(`SSH账号添加成功`)).SetData(m.SshUser)
					return ctx.JSON(data)
				}
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/term/account`))
			}
		}
		if err != nil && ctx.IsAjax() {
			return ctx.JSON(ctx.Data().SetError(err))
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.SshUser, ``, func(topName, fieldName string) string {
					return echo.LowerCaseFirstLetter(topName, fieldName)
				})
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}
	mg := model.NewSshUserGroup(ctx)
	_, e := mg.List(nil, nil, 1, -1)
	if err == nil {
		err = e
	}
	ctx.Set(`groupList`, mg.Objects())
	return ctx.Render(`term/account_edit`, err)
}

func AccountEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := model.NewSshUser(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if ctx.IsPost() {
		err = ctx.MustBind(m.SshUser, func(k string, v []string) (string, []string) {
			switch strings.ToLower(k) {
			case `created`: //禁止修改创建时间和用户名
				return ``, v
			case `password`, `passphrase`:
				if len(v) > 0 && len(v[0]) > 0 {
					v[0] = config.DefaultConfig.Encode(v[0])
				}
				return k, v
			}
			return k, v
		})

		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/term/account`))
			}
		}
	} else if err == nil {
		if len(m.SshUser.Password) > 0 {
			m.SshUser.Password = config.DefaultConfig.Decode(m.SshUser.Password)
		}
		if len(m.SshUser.Passphrase) > 0 {
			m.SshUser.Passphrase = config.DefaultConfig.Decode(m.SshUser.Passphrase)
		}
		echo.StructToForm(ctx, m.SshUser, ``, func(topName, fieldName string) string {
			return echo.LowerCaseFirstLetter(topName, fieldName)
		})
	}

	mg := model.NewSshUserGroup(ctx)
	_, e := mg.List(nil, nil, 1, -1)
	if err == nil {
		err = e
	}
	ctx.Set(`groupList`, mg.Objects())
	ctx.Set(`activeURL`, `/term/account`)
	return ctx.Render(`term/account_edit`, err)
}

func AccountDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewSshUser(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/term/account`))
}

func GroupIndex(ctx echo.Context) error {
	m := model.NewSshUserGroup(ctx)
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}))
	ret := handler.Err(ctx, err)
	ctx.Set(`listData`, m.Objects())
	return ctx.Render(`term/group`, ret)
}

func GroupAdd(ctx echo.Context) error {
	var err error
	m := model.NewSshUserGroup(ctx)
	if ctx.IsPost() {
		name := ctx.Form(`name`)
		if len(name) == 0 {
			err = ctx.E(`组名称不能为空`)
		} else if y, e := m.Exists(name); e != nil {
			err = e
		} else if y {
			err = ctx.E(`组名称已经存在`)
		} else {
			err = ctx.MustBind(m.SshUserGroup)
		}
		if err == nil {
			_, err = m.Add()
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/term/group`))
			}
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.SshUserGroup, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}

	return ctx.Render(`term/group_edit`, err)
}

func GroupEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := model.NewSshUserGroup(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if ctx.IsPost() {
		name := ctx.Form(`name`)
		if len(name) < 1 {
			err = ctx.E(`组名称不能为空`)
		} else if y, e := m.ExistsOther(name, id); e != nil {
			err = e
		} else if y {
			err = ctx.E(`组名称已经存在`)
		} else {
			err = ctx.MustBind(m.SshUserGroup, func(k string, v []string) (string, []string) {
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
				return ctx.Redirect(handler.URLFor(`/term/group`))
			}
		}
	} else if err == nil {
		echo.StructToForm(ctx, m.SshUserGroup, ``, echo.LowerCaseFirstLetter)
	}

	ctx.Set(`activeURL`, `/term/group`)
	return ctx.Render(`term/group_edit`, err)
}

func GroupDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewSshUserGroup(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/term/group`))
}

func Client(ctx echo.Context) error {
	ctx.Set(`activeURL`, `/term/account`)
	id := ctx.Formx(`id`).Uint()
	var err error
	m := model.NewSshUser(ctx)
	err = m.Get(nil, `id`, id)
	if err == nil {
		/*Test Code
			m.Passphrase = config.DefaultConfig.Decode(m.Passphrase)
			m.Password = config.DefaultConfig.Decode(m.Password)
			return m.ExecMultiCMD(ctx.Response(), "ls .", "export TESTENV=123", "echo $TESTENV")
		//*/
		q := url.Values{}
		q.Add(`id`, fmt.Sprint(m.Id))
		q.Add(`protocol`, m.Protocol)
		q.Add(`hostname`, m.Host)
		q.Add(`name`, m.Name)
		q.Add(`port`, fmt.Sprint(m.Port))
		q.Add(`user`, m.Username)
		q.Add(`password`, m.Password)
		q.Add(`url_prefix`, `term/client`)
		return ctx.Redirect(`/public/assets/backend/js/xterm/index.html?` + q.Encode())
		//return ctx.Redirect(handler.URLFor(`/term/client?`) + q.Encode())
	}
	return ctx.Render(`/term/client`, err)
}
