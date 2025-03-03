/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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

package manager

import (
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/nerrors"
	"github.com/coscms/webcore/library/notice"
	"github.com/coscms/webcore/library/nsql"
	"github.com/coscms/webcore/model"
)

var userLinks = []func(c *dbschema.NgingUser) string{}

func UserLink(c *dbschema.NgingUser) string {
	var t string
	for _, cl := range userLinks {
		t += cl(c)
	}
	return t
}

func AddUserLink(fn func(c *dbschema.NgingUser) string) {
	userLinks = append(userLinks, fn)
}

func User(ctx echo.Context) error {
	user := backend.User(ctx)
	username := ctx.Formx(`username`).String()
	online := ctx.Form(`online`)
	cond := db.Compounds{}
	if len(username) > 0 {
		cond.AddKV(`username`, db.Like(username+`%`))
	}
	if len(online) > 0 {
		if online == `Y` {
			cond.Add(db.Or(
				db.Cond{`id`: user.Id},
				db.Cond{`online`: online},
			))
		} else {
			cond.Add(db.Cond{`online`: `N`})
		}
	}
	nsql.SelectPageCond(ctx, &cond, `id`, `username`)
	m := model.NewUser(ctx)
	_, err := common.PagingWithLister(ctx, common.NewLister(m, nil, func(r db.Result) db.Result {
		return r.Select(factory.DBIGet().OmitSelect(m, `password`, `salt`, `safe_pwd`)...).OrderBy(`-id`)
	}, cond.And()))
	ret := common.Err(ctx, err)
	rows := m.Objects()
	var offlineUserIDs []uint
	for index, row := range rows {
		if row.Online == `Y` {
			if row.Id != user.Id && !notice.IsOnline(row.Username) {
				row.Online = `N`
				rows[index] = row
				offlineUserIDs = append(offlineUserIDs, row.Id)
			}
		} else {
			if row.Id == user.Id {
				row.Online = `Y`
				rows[index] = row
			}
		}
	}
	if len(offlineUserIDs) > 0 {
		m.NgingUser.UpdateField(nil, `online`, `N`, `id`, db.In(offlineUserIDs))
	}
	ctx.Set(`listData`, rows)
	ctx.SetFunc(`userLink`, UserLink)
	return ctx.Render(`/manager/user`, ret)
}

func UserAdd(ctx echo.Context) error {
	var err error
	m := model.NewUser(ctx)
	if ctx.IsPost() {
		m.Username = strings.TrimSpace(ctx.Form(`username`))
		m.Email = strings.TrimSpace(ctx.Form(`email`))
		m.Mobile = strings.TrimSpace(ctx.Form(`mobile`))
		m.Password = strings.TrimSpace(ctx.Form(`password`))
		confirmPwd := strings.TrimSpace(ctx.Form(`confirmPwd`))
		m.Password, err = backend.DecryptPassword(ctx, m.Username, m.Password)
		if err != nil {
			return ctx.NewError(code.InvalidParameter, `密码解密失败: %v`, err).SetZone(`password`)
		}
		confirmPwd, err = backend.DecryptPassword(ctx, m.Username, confirmPwd)
		if err != nil {
			return ctx.NewError(code.InvalidParameter, `您输入的确认密码解密失败: %v`, err).SetZone(`confirmPwd`)
		}
		if m.Password != confirmPwd {
			err = ctx.E(`密码与确认密码不一致`)
			goto END
		}
		m.Avatar = strings.TrimSpace(ctx.Form(`avatar`))
		m.Gender = strings.TrimSpace(ctx.Form(`gender`))
		m.RoleIds = strings.Join(ctx.FormValues(`roleIds`), `,`)
		err = m.Add()
		if err == nil {
			common.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(backend.URLFor(`/manager/user`))
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				setFormData(ctx, m)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}

END:
	ctx.Set(`activeURL`, `/manager/user`)
	roleM := model.NewUserRole(ctx)
	roleM.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Select(`id`, `name`, `description`)
	}, 0, -1, db.And(db.Cond{`parent_id`: 0}))
	ctx.Set(`roleList`, roleM.Objects())
	ctx.SetFunc(`isChecked`, func(roleId uint) bool {
		return false
	})
	return ctx.Render(`/manager/user_edit`, common.Err(ctx, err))
}

func UserEdit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewUser(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		common.SendFail(ctx, err.Error())
		return ctx.Redirect(backend.URLFor(`/manager/user`))
	}
	if ctx.IsPost() {
		modifyPwd := ctx.Form(`modifyPwd`) == `1`
		password := strings.TrimSpace(ctx.Form(`password`))
		confirmPwd := strings.TrimSpace(ctx.Form(`confirmPwd`))
		if modifyPwd {
			password, err = backend.DecryptPassword(ctx, m.Username, password)
			if err != nil {
				err = ctx.NewError(code.InvalidParameter, `新密码解密失败: %v`, err).SetZone(`newPass`)
				goto END
			}
			confirmPwd, err = backend.DecryptPassword(ctx, m.Username, confirmPwd)
			if err != nil {
				err = ctx.NewError(code.InvalidParameter, `您输入的确认密码解密失败: %v`, err).SetZone(`confirmPwd`)
				goto END
			}
			if password != confirmPwd {
				err = ctx.E(`密码与确认密码不一致`)
				goto END
			}
			m.Password = password
		}
		m.Username = strings.TrimSpace(ctx.Form(`username`))
		m.Email = strings.TrimSpace(ctx.Form(`email`))
		m.Mobile = strings.TrimSpace(ctx.Form(`mobile`))
		m.Avatar = strings.TrimSpace(ctx.Form(`avatar`))
		m.Gender = strings.TrimSpace(ctx.Form(`gender`))
		m.RoleIds = strings.Join(ctx.FormValues(`roleIds`), `,`)
		m.Disabled = strings.TrimSpace(ctx.Form(`disabled`))
		if err == nil {
			m.Id = id
			set := map[string]interface{}{
				`email`:    m.Email,
				`mobile`:   m.Mobile,
				`username`: m.Username,
				`role_ids`: m.RoleIds,
				`avatar`:   m.Avatar,
				`gender`:   m.Gender,
				`disabled`: m.Disabled,
			}
			err = m.UpdateField(id, set)
		}
		if err == nil {
			common.SendOk(ctx, ctx.T(`修改成功`))
			return ctx.Redirect(backend.URLFor(`/manager/user`))
		}
	}

END:
	setFormData(ctx, m)
	ctx.Set(`activeURL`, `/manager/user`)
	roleM := model.NewUserRole(ctx)
	roleM.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Select(`id`, `name`, `description`)
	}, 0, -1, db.And(db.Cond{`parent_id`: 0}))
	ctx.Set(`roleList`, roleM.Objects())
	return ctx.Render(`/manager/user_edit`, common.Err(ctx, err))
}

func setFormData(ctx echo.Context, m *model.User) {
	m.Password = ``
	m.Salt = ``
	echo.StructToForm(ctx, m.NgingUser, ``, echo.LowerCaseFirstLetter)
	var roleIds []uint
	if len(m.RoleIds) > 0 {
		roleIds = param.StringSlice(strings.Split(m.RoleIds, `,`)).Uint()
	}
	ctx.SetFunc(`isChecked`, func(roleId uint) bool {
		for _, rid := range roleIds {
			if rid == roleId {
				return true
			}
		}
		return false
	})
}

func UserDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint64()
	m := model.NewUser(ctx)
	if id == 1 {
		common.SendFail(ctx, ctx.T(`创始人不可删除`))
		return ctx.Redirect(backend.URLFor(`/manager/user`))
	}
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/manager/user`))
}

// UserKick 踢🦶用户下线
func UserKick(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	user := backend.User(ctx)
	if user == nil {
		return nerrors.ErrUserNotLoggedIn
	}
	if id == user.Id {
		return ctx.E(`不能踢自己`)
	}
	m := model.NewUser(ctx)
	err := m.Get(func(r db.Result) db.Result {
		return r.Select(`session_id`)
	}, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	if len(m.SessionId) == 0 {
		common.SendFail(ctx, ctx.T(`此用户没有 session id 记录`))
	} else {
		err = ctx.Session().RemoveID(m.SessionId)
		if err == nil {
			m.NgingUser.UpdateField(nil, `session_id`, ``, `id`, id)
			common.SendOk(ctx, ctx.T(`操作成功`))
		} else {
			common.SendFail(ctx, err.Error())
		}
	}

	return ctx.Redirect(backend.URLFor(`/manager/user`))
}
