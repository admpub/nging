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
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/formfilter"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/config"

	"github.com/nging-plugins/frpmanager/pkg/dbschema"
	"github.com/nging-plugins/frpmanager/pkg/library/cmder"
	"github.com/nging-plugins/frpmanager/pkg/library/utils"
	"github.com/nging-plugins/frpmanager/pkg/model"
)

func ClientIndex(ctx echo.Context) error {
	groupId := ctx.Formx(`groupId`).Uint()
	m := model.NewFrpClient(ctx)
	cond := db.Compounds{}
	if groupId > 0 {
		cond.AddKV(`group_id`, groupId)
	}
	q := ctx.Formx(`q`).String()
	if len(q) > 0 {
		cond.AddKV(`name`, db.Like(`%`+q+`%`))
	}
	var clientAndGroup []*model.FrpClientAndGroup
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, &clientAndGroup, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()))
	for k, u := range clientAndGroup {
		clientAndGroup[k].Running = config.FromCLI().IsRunning(`frpclient.` + fmt.Sprint(u.Id))
	}

	mg := model.NewFrpGroup(ctx)
	var groupList []*dbschema.NgingFrpGroup
	mg.ListByOffset(&groupList, nil, 0, -1)
	ctx.Set(`listData`, clientAndGroup)
	ctx.Set(`groupList`, groupList)
	ctx.Set(`groupId`, groupId)
	ctx.Set(`isRunning`, config.FromCLI().CmdHasGroup(`frpclient`))
	return ctx.Render(`frp/client_index`, handler.Err(ctx, err))
}

func clientFormFilter(opts ...formfilter.Options) echo.FormDataFilter {
	opts = append(opts, formfilter.Exclude(`extra`))
	return formfilter.Build(opts...)
}

func ClientAdd(ctx echo.Context) error {
	m := model.NewFrpClient(ctx)
	user := handler.User(ctx)
	cm, err := cmder.GetClient()
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingFrpClient, clientFormFilter())
		if err == nil {
			m.NgingFrpClient.Extra, err = form2ExtraStr(ctx)
		}
		if err == nil {
			m.Uid = user.Id
			_, err = m.Add()
			if err == nil {
				err = utils.SaveConfigFile(m.NgingFrpClient)
			}
			if err == nil {
				if m.Disabled == `N` {
					err = cm.StartBy(m.NgingFrpClient.Id)
				}
				if err != nil {
					handler.SendOk(ctx, ctx.T(`保存成功。但启动失败: %v`, err.Error()))
				} else {
					handler.SendOk(ctx, ctx.T(`操作成功`))
				}
				return ctx.Redirect(handler.URLFor(`/frp/client_index`))
			}
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.NgingFrpClient, ``, func(topName, fieldName string) string {
					return echo.LowerCaseFirstLetter(topName, fieldName)
				})
				err = copyClientExtra2Form(ctx, m.NgingFrpClient)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
		if len(ctx.Form(`logFile`)) == 0 {
			logRandName := time.Now().Format(`20060102`) + `-` + com.RandomAlphanumeric(8)
			ctx.Request().Form().Set(`logFile`, `./data/logs/frp/client.`+logRandName+`.log`)
		}
	}
	mg := model.NewFrpGroup(ctx)
	_, e := mg.List(nil, nil, 1, -1)
	if err == nil {
		err = e
	}
	ctx.Set(`groupList`, mg.Objects())
	ctx.Set(`activeURL`, `/frp/client_index`)
	setAddonFunc(ctx)
	return ctx.Render(`frp/client_edit`, err)
}

func ClientEdit(ctx echo.Context) error {
	cm, err := cmder.GetClient()
	if err != nil {
		return err
	}
	id := ctx.Formx(`id`).Uint()
	m := model.NewFrpClient(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.NewError(code.DataNotFound, `数据不存在`)
		}
		return err
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingFrpClient, clientFormFilter(formfilter.Exclude(`created`)))
		if err == nil {
			m.NgingFrpClient.Extra, err = form2ExtraStr(ctx)
		}
		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				err = utils.SaveConfigFile(m.NgingFrpClient)
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
				return ctx.Redirect(handler.URLFor(`/frp/client_index`))
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
			err = utils.SaveConfigFile(m.NgingFrpClient)
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			var opType, status string
			if m.Disabled == `N` {
				err = cm.RestartBy(fmt.Sprintf(`%d`, m.NgingFrpClient.Id))
				opType = ctx.T(`启动失败`)
				status = `started`
			} else {
				err = cm.StopBy(fmt.Sprintf(`%d`, m.NgingFrpClient.Id))
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
		echo.StructToForm(ctx, m.NgingFrpClient, ``, func(topName, fieldName string) string {
			return echo.LowerCaseFirstLetter(topName, fieldName)
		})
		err = copyClientExtra2Form(ctx, m.NgingFrpClient)
	}

	mg := model.NewFrpGroup(ctx)
	_, e := mg.List(nil, nil, 1, -1)
	if err == nil {
		err = e
	}
	ctx.Set(`groupList`, mg.Objects())
	ctx.Set(`activeURL`, `/frp/client_index`)
	setAddonFunc(ctx)
	return ctx.Render(`frp/client_edit`, err)
}

func form2ExtraStr(ctx echo.Context) (extraStr string, err error) {
	extra := url.Values{}
	for key, val := range ctx.Forms() {
		if !strings.HasPrefix(key, `extra[`) {
			continue
		}
		extra[key] = val
	}
	var b []byte
	b, err = json.Marshal(extra)
	if err != nil {
		return
	}
	extraStr = string(b)
	return
}

func copyClientExtra2Form(ctx echo.Context, cfg *dbschema.NgingFrpClient) (err error) {
	if len(cfg.Extra) == 0 {
		return nil
	}
	extra := url.Values{}
	f := ctx.Request().Form()
	err = json.Unmarshal([]byte(cfg.Extra), &extra)
	if err != nil {
		return
	}
	mapx := echo.NewMapx(extra)
	mapx = mapx.Get(`extra`)
	if mapx == nil {
		return
	}
	sections := []*Section{}
	for section := range mapx.Map {
		if cfg.Type == section {
			continue
		}
		s := &Section{
			Section: section,
			Addon:   regexNumEnd.ReplaceAllString(section, ``),
		}
		sections = append(sections, s)
	}
	f.Merge(extra)
	ctx.Set(`sections`, sections)
	return
}

func ClientDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewFrpClient(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		err = utils.SaveConfigFile(&dbschema.NgingFrpClient{Disabled: `Y`, Id: id})
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/frp/client_index`))
}

func ClientLog(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	if id < 1 {
		return ctx.JSON(ctx.Data().SetError(ctx.E(`id无效`)))
	}
	var err error
	m := model.NewFrpClient(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.E(`不存在id为%d的配置`)
		}
		return ctx.JSON(ctx.Data().SetError(err))
	}
	return common.LogShow(ctx, m.LogFile, echo.H{`title`: m.Name})
}
