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
	"html/template"
	"net/url"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/nging-plugins/caddymanager/application/dbschema"
	"github.com/nging-plugins/caddymanager/application/library/engine"
	"github.com/nging-plugins/caddymanager/application/model"
)

func VhostIndex(ctx echo.Context) error {
	m := model.NewVhost(ctx)
	groupID := ctx.Formx(`groupId`).Uint()
	serverIdent := ctx.Formx(`serverIdent`).String()
	engineType := ctx.Formx(`engine`).String()
	cond := db.Compounds{}
	if groupID > 0 {
		cond.AddKV(`a.group_id`, groupID)
	}
	if len(engineType) > 0 {
		if engineType == `default` {
			serverIdent = engineType
		} else {
			cond.AddKV(`b.engine`, engineType)
		}
	}
	if len(serverIdent) > 0 {
		cond.AddKV(`a.server_ident`, serverIdent)
	}
	q := ctx.Formx(`q`).String()
	if len(q) > 0 {
		cond.AddKV(`a.name`, db.Like(`%`+q+`%`))
	}
	var rowAndGroup []*model.VhostAndGroup
	p := m.NewParam().SetCols(`a.*`, `b.name AS serverName`, `b.engine AS serverEngine`).SetAlias(`a`).SetMW(func(r db.Result) db.Result {
		return r.OrderBy(`-a.id`)
	}).SetRecv(&rowAndGroup)
	p.AddJoin(`LEFT`, dbschema.NewNgingVhostServer(ctx).Short_(), `b`, `b.ident=a.server_ident`)
	p.AddArgs(cond.And())
	_, err := handler.PagingWithList(ctx, p)
	mg := dbschema.NewNgingVhostGroup(ctx)
	var groupList []*dbschema.NgingVhostGroup
	mg.ListByOffset(&groupList, nil, 0, -1)
	ctx.Set(`listData`, rowAndGroup)
	ctx.Set(`groupList`, groupList)
	ctx.Set(`engineList`, engine.Engines.Slice())
	ctx.Set(`groupId`, groupID)
	currentHost := ctx.Host()
	ctx.SetFunc(`generateHostURL`, func(hosts string) []template.HTML {
		return generateHostURL(currentHost, hosts)
	})
	ctx.SetFunc(`engineName`, engine.Engines.Get)
	return ctx.Render(`caddy/vhost`, handler.Err(ctx, err))
}

func Vhostbuild(ctx echo.Context) error {
	groupID := ctx.Formx(`groupId`).Uint()
	serverIdent := ctx.Formx(`serverIdent`).String()
	engineType := ctx.Formx(`engine`).String()
	err := vhostbuild(ctx, groupID, serverIdent, engineType)
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
	}
	handler.SendOk(ctx, ctx.T(`操作成功`))
	return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
}

func VhostAdd(ctx echo.Context) error {
	var err error
	m := model.NewVhost(ctx)
	if ctx.IsPost() {
		receiveFormData(ctx, m.NgingVhost)
		var b []byte
		b, err = json.Marshal(ctx.Forms())
		switch {
		case err == nil:
			ctx.Begin()
			m.Setting = string(b)
			_, err = m.Add()
			if err != nil {
				ctx.Rollback()
				break
			}
			fallthrough
		case 0 == 1:
			err = saveVhostData(ctx, m.NgingVhost, ctx.Forms(),
				m.Disabled == common.BoolN && len(ctx.Form(`restart`)) > 0,
				m.Disabled == common.BoolN && len(ctx.Form(`removeCachedCert`)) > 0,
			)
			ctx.End(err == nil)
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, db.Cond{`id`: id})
			if err == nil {
				var formData url.Values
				if e := json.Unmarshal([]byte(m.Setting), &formData); e == nil {
					for key, values := range formData {
						for k, v := range values {
							if k == 0 {
								ctx.Request().Form().Set(key, v)
								continue
							}
							ctx.Request().Form().Add(key, v)
						}
					}
				}
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}
	setVhostForm(ctx)
	ctx.Set(`isAdd`, true)
	ctx.Set(`title`, ctx.T(`添加网站`))
	return ctx.Render(`caddy/vhost_edit`, err)
}

func setVhostForm(ctx echo.Context) {
	ctx.SetFunc(`Val`, func(name, defaultValue string) string {
		return ctx.Form(name, defaultValue)
	})
	g := model.NewVhostGroup(ctx)
	g.ListByOffset(nil, nil, 0, -1)
	ctx.Set(`groupList`, g.Objects())
	svr := model.NewVhostServer(ctx)
	svr.ListByOffset(nil, nil, 0, -1, db.Cond{`disabled`: common.BoolN})
	ctx.Set(`serverList`, svr.Objects())
}

func VhostDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	if id < 1 {
		handler.SendFail(ctx, ctx.T(`id无效`))
		return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
	}
	m := model.NewVhost(ctx)
	err := m.Get(func(r db.Result) db.Result {
		return r.Select(`server_ident`)
	}, `id`, id)
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
	}
	err = m.Delete(nil, db.Cond{`id`: id})
	if err != nil {
		handler.SendFail(ctx, err.Error())
	} else {
		err = DeleteCaddyfileByID(ctx, m.ServerIdent, id)
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
		}
	}
	return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
}

func VhostEdit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	if id < 1 {
		handler.SendFail(ctx, ctx.T(`id无效`))
		return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
	}

	var err error
	m := model.NewVhost(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
	}
	if ctx.IsPost() {
		old := *m.NgingVhost
		receiveFormData(ctx, m.NgingVhost)
		var b []byte
		b, err = json.Marshal(ctx.Forms())
		switch {
		case err == nil:
			ctx.Begin()
			m.Setting = string(b)
			err = m.Edit(nil, db.Cond{`id`: id})
			if err != nil {
				ctx.Rollback()
				break
			}
			fallthrough
		case 0 == 1:
			removeCachedCert := ctx.Form(`removeCachedCert`)
			if len(removeCachedCert) > 0 && removeCachedCert == `1` {
				m.RemoveCachedCert()
			}
			if old.ServerIdent != m.ServerIdent {
				DeleteCaddyfileByID(ctx, old.ServerIdent, m.Id)
			}
			err = saveVhostData(ctx, m.NgingVhost, ctx.Forms(),
				len(ctx.Form(`restart`)) > 0,
				m.Disabled == common.BoolN && removeCachedCert == `1`,
			)
			ctx.End(err == nil)
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
		}
	} else if ctx.IsAjax() {
		data := ctx.Data()
		disabled := ctx.Query(`disabled`)
		if len(disabled) > 0 {
			if !common.IsBoolFlag(disabled) {
				return ctx.NewError(code.InvalidParameter, ``).SetZone(`disabled`)
			}
			m.Disabled = disabled
			err = m.UpdateField(nil, `disabled`, disabled, db.Cond{`id`: id})
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			if m.Disabled == `Y` {
				err = DeleteCaddyfileByID(ctx, m.ServerIdent, id)
			} else {
				var formData url.Values
				err = json.Unmarshal([]byte(m.Setting), &formData)
				if err == nil {
					err = saveVhostData(ctx, m.NgingVhost, formData, true, false)
				}
			}
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`操作成功`))
		}
		return ctx.JSON(data)
	} else {
		var formData url.Values
		if e := json.Unmarshal([]byte(m.Setting), &formData); e == nil {
			for key, values := range formData {
				for k, v := range values {
					if k == 0 {
						ctx.Request().Form().Set(key, v)
						continue
					}
					ctx.Request().Form().Add(key, v)
				}
			}
		}
		echo.StructToForm(ctx, m.NgingVhost, ``, echo.LowerCaseFirstLetter)
	}
	setVhostForm(ctx)
	ctx.Set(`activeURL`, `/caddy/vhost`)
	ctx.Set(`isAdd`, false)
	ctx.Set(`title`, ctx.T(`修改网站`))
	return ctx.Render(`caddy/vhost_edit`, err)
}
