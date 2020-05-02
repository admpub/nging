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

	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/model"
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
	var serverAndGroup []*model.FrpServerAndGroup
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, &serverAndGroup, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()))
	for k, u := range serverAndGroup {
		serverAndGroup[k].Running= config.DefaultCLIConfig.IsRunning(`frpserver.` + fmt.Sprint(u.Id))
	}

	mg := model.NewFrpGroup(ctx)
	var groupList []*dbschema.NgingFrpGroup
	mg.ListByOffset(&groupList, nil, 0, -1)
	ctx.Set(`listData`, serverAndGroup)
	ctx.Set(`groupList`, groupList)
	ctx.Set(`groupId`, groupId)
	ctx.Set(`isRunning`, config.DefaultCLIConfig.CmdHasGroup(`frpserver`))
	return ctx.Render(`frp/server_index`, handler.Err(ctx, err))
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
			err = ctx.MustBind(m.NgingFrpServer)
		}

		if err == nil {
			_, err = m.Add()
			if err == nil {
				err = config.DefaultCLIConfig.FRPSaveConfigFile(m.NgingFrpServer)
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
				echo.StructToForm(ctx, m.NgingFrpServer, ``, func(topName, fieldName string) string {
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
			err = ctx.MustBind(m.NgingFrpServer, echo.ExcludeFieldName(`created`))
		}

		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				err = config.DefaultCLIConfig.FRPSaveConfigFile(m.NgingFrpServer)
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
			err = config.DefaultCLIConfig.FRPSaveConfigFile(m.NgingFrpServer)
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`状态已经更改成功，请重启服务端令其生效`))
			return ctx.JSON(data)
		}
	}
	if err == nil {
		echo.StructToForm(ctx, m.NgingFrpServer, ``, func(topName, fieldName string) string {
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
		err = config.DefaultCLIConfig.FRPSaveConfigFile(&dbschema.NgingFrpServer{Disabled: `Y`, Id: id})
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/frp/server_index`))
}

func ServerLog(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	if id < 1 {
		return ctx.JSON(ctx.Data().SetError(ctx.E(`id无效`)))
	}
	var err error
	m := model.NewFrpServer(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.E(`不存在id为%d的配置`)
		}
		return ctx.JSON(ctx.Data().SetError(err))
	}
	return common.LogShow(ctx, m.LogFile, echo.H{`title`: m.Name})
}
