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
	"fmt"
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/formfilter"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/config"

	"github.com/nging-plugins/frpmanager/application/dbschema"
	"github.com/nging-plugins/frpmanager/application/library/cmder"
	"github.com/nging-plugins/frpmanager/application/library/frp"
	"github.com/nging-plugins/frpmanager/application/library/utils"
	"github.com/nging-plugins/frpmanager/application/model"
)

func ServerIndex(ctx echo.Context) error {
	groupId := ctx.Formx(`groupId`).Uint()
	m := model.NewFrpServer(ctx)
	cond := db.Compounds{}
	if groupId > 0 {
		cond.AddKV(`group_id`, groupId)
	}
	common.SelectPageCond(ctx, &cond)
	var serverAndGroup []*model.FrpServerAndGroup
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, &serverAndGroup, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()))
	for k, u := range serverAndGroup {
		serverAndGroup[k].Running = config.FromCLI().IsRunning(`frpserver.` + fmt.Sprint(u.Id))
	}

	mg := model.NewFrpGroup(ctx)
	var groupList []*dbschema.NgingFrpGroup
	mg.ListByOffset(&groupList, nil, 0, -1)
	ctx.Set(`listData`, serverAndGroup)
	ctx.Set(`groupList`, groupList)
	ctx.Set(`groupId`, groupId)
	ctx.Set(`isRunning`, config.FromCLI().CmdHasGroup(`frpserver`))
	return ctx.Render(`frp/server_index`, handler.Err(ctx, err))
}

func serverFormFilter(opts ...formfilter.Options) echo.FormDataFilter {
	opts = append(opts,
		formfilter.JoinValues(`Plugins`),
	)
	return formfilter.Build(opts...)
}

func ServerAdd(ctx echo.Context) error {
	cm, err := cmder.GetServer()
	if err != nil {
		return err
	}
	m := model.NewFrpServer(ctx)
	user := handler.User(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingFrpServer, serverFormFilter())
		if err == nil {
			m.Uid = user.Id
			_, err = m.Add()
			if err == nil {
				err = utils.SaveConfigFile(m.NgingFrpServer)
			}
			if err == nil {
				if m.Disabled == `N` {
					err = cm.StartBy(m.NgingFrpServer.Id)
				}
				if err != nil {
					handler.SendOk(ctx, ctx.T(`保存成功。但启动失败: %v`, err.Error()))
				} else {
					handler.SendOk(ctx, ctx.T(`操作成功`))
				}
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
		if len(ctx.Form(`token`)) == 0 {
			defaultToken, _ := common.GenPassword()
			ctx.Request().Form().Set(`token`, defaultToken)
		}
		if len(ctx.Form(`logFile`)) == 0 {
			logRandName := time.Now().Format(`20060102`) + `-` + com.RandomAlphanumeric(8)
			ctx.Request().Form().Set(`logFile`, `./data/logs/frp/server.`+logRandName+`.log`)
		}
	}
	mg := model.NewFrpGroup(ctx)
	_, e := mg.List(nil, nil, 1, -1)
	if err == nil {
		err = e
	}
	ctx.Set(`groupList`, mg.Objects())
	ctx.Set(`pluginList`, frp.ServerPluginSlice())
	ctx.SetFunc(`isChecked`, func(name string) bool {
		return false
	})
	ctx.Set(`activeURL`, `/frp/server_index`)
	return ctx.Render(`frp/server_edit`, err)
}

func ServerEdit(ctx echo.Context) error {
	cm, err := cmder.GetServer()
	if err != nil {
		return err
	}
	id := ctx.Formx(`id`).Uint()
	m := model.NewFrpServer(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.NewError(code.DataNotFound, `数据不存在`)
		}
		return err
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingFrpServer, serverFormFilter(formfilter.Exclude(`created`)))
		if err == nil {
			m.Id = id
			// 强制设置, select2 在没有选中任何值的情况下不会提交此字段，所以需要手动检查和设置
			if len(ctx.FormValues(`plugins`)) == 0 {
				m.Plugins = ``
			}
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				err = utils.SaveConfigFile(m.NgingFrpServer)
			}
			if err == nil {
				var opType string
				if m.Disabled == `N` {
					err = cm.RestartBy(fmt.Sprintf(`%d`, m.Id))
					opType = ctx.T(`启动失败`)
				} else {
					err = cm.StopBy(fmt.Sprintf(`%d`, m.Id))
					opType = ctx.T(`关闭失败`)
				}
				if err != nil {
					handler.SendOk(ctx, ctx.T(`保存成功。但%s: %v`, opType, err.Error()))
				} else {
					handler.SendOk(ctx, ctx.T(`操作成功`))
				}
				return ctx.Redirect(handler.URLFor(`/frp/server_index`))
			}
		}
	} else if ctx.IsAjax() {
		disabled := ctx.Query(`disabled`)
		if len(disabled) > 0 {
			m.Disabled = disabled
			data := ctx.Data()
			err = m.Update(nil, db.Cond{`id`: id})
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			err = utils.SaveConfigFile(m.NgingFrpServer)
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			var opType, status string
			if m.Disabled == `N` {
				err = cm.RestartBy(fmt.Sprintf(`%d`, m.Id))
				opType = ctx.T(`启动失败`)
				status = `started`
			} else {
				err = cm.StopBy(fmt.Sprintf(`%d`, m.Id))
				opType = ctx.T(`关闭失败`)
				status = `stopped`
			}
			if err != nil {
				data.SetData(echo.H{`status`: `failed`})
				data.SetInfo(ctx.T(`状态已经更改成功。但%s: %v`, opType, err.Error()))
			} else {
				data.SetData(echo.H{`status`: status})
				data.SetInfo(ctx.T(`状态已经更改成功`))
			}
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
	ctx.Set(`pluginList`, frp.ServerPluginSlice())
	ctx.Set(`activeURL`, `/frp/server_index`)
	var plugins []string
	if len(m.Plugins) > 0 {
		plugins = param.StringSlice(strings.Split(m.Plugins, `,`)).Filter().String()
	}
	ctx.SetFunc(`isChecked`, func(name string) bool {
		for _, rid := range plugins {
			if rid == name {
				return true
			}
		}
		return false
	})
	return ctx.Render(`frp/server_edit`, err)
}

func ServerDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewFrpServer(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		err = utils.SaveConfigFile(&dbschema.NgingFrpServer{Disabled: `Y`, Id: id})
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
