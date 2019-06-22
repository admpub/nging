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
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
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
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()))
	ret := handler.Err(ctx, err)
	clients := m.Objects()
	gIds := []uint{}
	clientAndGroup := make([]*model.FrpClientAndGroup, len(clients))
	for k, u := range clients {
		clientAndGroup[k] = &model.FrpClientAndGroup{
			FrpClient: u,
			Running:   config.DefaultCLIConfig.IsRunning(`frpclient.` + fmt.Sprint(u.Id)),
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
		_, err = mg.List(&groupList, nil, 1, -1, db.Cond{`id IN`: gIds})
		if err != nil {
			if ret == nil {
				ret = err
			}
		} else {
			for k, v := range clientAndGroup {
				for _, g := range groupList {
					if g.Id == v.GroupId {
						clientAndGroup[k].Group = g
						break
					}
				}
			}
		}
	}
	ctx.Set(`listData`, clientAndGroup)
	mg.ListByOffset(&groupList, nil, 0, -1)
	ctx.Set(`groupList`, groupList)
	ctx.Set(`groupId`, groupId)
	ctx.Set(`isRunning`, config.DefaultCLIConfig.CmdHasGroup(`frpclient`))
	return ctx.Render(`frp/client_index`, ret)
}

func ClientAdd(ctx echo.Context) error {
	var err error
	m := model.NewFrpClient(ctx)
	if ctx.IsPost() {
		name := ctx.Form(`name`)
		if len(name) == 0 {
			err = ctx.E(`名称不能为空`)
		} else if y, e := m.Exists(name); e != nil {
			err = e
		} else if y {
			err = ctx.E(`名称已经存在`)
		} else {
			err = ctx.MustBind(m.FrpClient)
		}
		if err == nil {
			m.FrpClient.Extra, err = form2ExtraStr(ctx)
		}
		if err == nil {
			_, err = m.Add()
			if err == nil {
				err = config.DefaultCLIConfig.FRPSaveConfigFile(m.FrpClient)
			}
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/frp/client_index`))
			}
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.FrpClient, ``, func(topName, fieldName string) string {
					return echo.LowerCaseFirstLetter(topName, fieldName)
				})
				err = copyClientExtra2Form(ctx, m.FrpClient)
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
	ctx.Set(`activeURL`, `/frp/client_index`)
	setAddonFunc(ctx)
	return ctx.Render(`frp/client_edit`, err)
}

func ClientEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := model.NewFrpClient(ctx)
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
			err = ctx.MustBind(m.FrpClient, echo.ExcludeFieldName(`created`))
		}
		if err == nil {
			m.FrpClient.Extra, err = form2ExtraStr(ctx)
		}
		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				err = config.DefaultCLIConfig.FRPSaveConfigFile(m.FrpClient)
			}
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
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
			err = config.DefaultCLIConfig.FRPSaveConfigFile(m.FrpClient)
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`状态已经更改成功，请重启客户端令其生效`))
			return ctx.JSON(data)
		}
	}
	if err == nil {
		echo.StructToForm(ctx, m.FrpClient, ``, func(topName, fieldName string) string {
			return echo.LowerCaseFirstLetter(topName, fieldName)
		})
		err = copyClientExtra2Form(ctx, m.FrpClient)
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

func copyClientExtra2Form(ctx echo.Context, cfg *dbschema.FrpClient) (err error) {
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
		s := &Section{
			Section: section,
			Addon:   regexNumEnd.ReplaceAllString(section, ``),
		}
		if cfg.Type == s.Addon {
			continue
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
		err = config.DefaultCLIConfig.FRPSaveConfigFile(&dbschema.FrpClient{Disabled: `Y`, Id: id})
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
