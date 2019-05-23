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

package frp

import (
	"fmt"
	"strings"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func ServerIndex(ctx echo.Context) error {
	groupId := ctx.Formx(`groupId`).Uint()
	m := model.NewFrpServer(ctx)
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
	servers := m.Objects()
	gIds := []uint{}
	serverAndGroup := make([]*model.FrpServerAndGroup, len(servers))
	for k, u := range servers {
		serverAndGroup[k] = &model.FrpServerAndGroup{
			FrpServer: u,
			Running:   config.DefaultCLIConfig.IsRunning(`frpserver.` + fmt.Sprint(u.Id)),
		}
		if u.GroupId < 1 {
			continue
		}
		if !com.InUintSlice(u.GroupId, gIds) {
			gIds = append(gIds, u.GroupId)
		}
	}

	mg := model.NewFrpGroup(ctx)
	var groupList []*dbschema.FrpGroup
	if len(gIds) > 0 {
		_, err = mg.List(&groupList, nil, 1, 1000, db.Cond{`id IN`: gIds})
		if err != nil {
			if ret == nil {
				ret = err
			}
		} else {
			for k, v := range serverAndGroup {
				for _, g := range groupList {
					if g.Id == v.GroupId {
						serverAndGroup[k].Group = g
						break
					}
				}
			}
		}
	}
	ctx.Set(`listData`, serverAndGroup)
	mg.ListByOffset(&groupList, nil, 0, -1)
	ctx.Set(`groupList`, groupList)
	ctx.Set(`groupId`, groupId)
	ctx.Set(`isRunning`, config.DefaultCLIConfig.CmdHasGroup(`frpserver`))
	return ctx.Render(`frp/server_index`, ret)
}

func ServerAdd(ctx echo.Context) error {
	var err error
	m := model.NewFrpServer(ctx)
	if ctx.IsPost() {
		name := ctx.Form(`name`)
		if len(name) == 0 {
			err = ctx.E(`名称不能为空`)
		} else if y, e := m.Exists(name); e != nil {
			err = e
		} else if y {
			err = ctx.E(`名称已经存在`)
		} else {
			err = ctx.MustBind(m.FrpServer)
		}

		if err == nil {
			_, err = m.Add()
			if err == nil {
				err = config.DefaultCLIConfig.FRPSaveConfigFile(m.FrpServer)
			}
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/frp/server_index`))
			}
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.FrpServer, ``, func(topName, fieldName string) string {
					return echo.LowerCaseFirstLetter(topName, fieldName)
				})
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}
	mg := model.NewFrpGroup(ctx)
	_, e := mg.List(nil, nil, 1, -1)
	if err == nil {
		err = e
	}
	ctx.Set(`groupList`, mg.Objects())
	ctx.Set(`activeURL`, `/frp/server_index`)
	return ctx.Render(`frp/server_edit`, err)
}

func ServerEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := model.NewFrpServer(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if ctx.IsPost() {
		name := ctx.Form(`name`)
		if len(name) == 0 {
			err = ctx.E(`名称不能为空`)
		} else if y, e := m.Exists(name, id); e != nil {
			err = e
		} else if y {
			err = ctx.E(`名称已经存在`)
		} else {
			err = ctx.MustBind(m.FrpServer, func(k string, v []string) (string, []string) {
				switch strings.ToLower(k) {
				case `created`: //禁止修改创建时间和用户名
					return ``, v
				}
				return k, v
			})
		}

		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				err = config.DefaultCLIConfig.FRPSaveConfigFile(m.FrpServer)
			}
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/frp/server_index`))
			}
		}
	} else if ctx.IsAjax() {
		disabled := ctx.Query(`disabled`)
		if len(disabled) > 0 {
			m.Disabled = disabled
			data := ctx.Data()
			err = m.Edit(nil, db.Cond{`id`: id})
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			err = config.DefaultCLIConfig.FRPSaveConfigFile(m.FrpServer)
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`状态已经更改成功，请重启服务端令其生效`))
			return ctx.JSON(data)
		}
	}
	if err == nil {
		echo.StructToForm(ctx, m.FrpServer, ``, func(topName, fieldName string) string {
			return echo.LowerCaseFirstLetter(topName, fieldName)
		})
	}

	mg := model.NewFrpGroup(ctx)
	_, e := mg.List(nil, nil, 1, -1)
	if err == nil {
		err = e
	}
	ctx.Set(`groupList`, mg.Objects())
	ctx.Set(`activeURL`, `/frp/server_index`)
	return ctx.Render(`frp/server_edit`, err)
}

func ServerDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewFrpServer(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		err = config.DefaultCLIConfig.FRPSaveConfigFile(&dbschema.FrpServer{Disabled: `Y`, Id: id})
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/frp/server_index`))
}
