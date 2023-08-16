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
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/handler"
	"github.com/nging-plugins/caddymanager/application/dbschema"
	"github.com/nging-plugins/caddymanager/application/library/cmder"
	"github.com/nging-plugins/caddymanager/application/model"
)

func VhostIndex(ctx echo.Context) error {
	m := model.NewVhost(ctx)
	groupID := ctx.Formx(`groupId`).Uint()
	cond := db.Compounds{}
	if groupID > 0 {
		cond.AddKV(`group_id`, groupID)
	}
	q := ctx.Formx(`q`).String()
	if len(q) > 0 {
		cond.AddKV(`name`, db.Like(`%`+q+`%`))
	}
	var rowAndGroup []*model.VhostAndGroup
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, &rowAndGroup, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()))
	mg := dbschema.NewNgingVhostGroup(ctx)
	var groupList []*dbschema.NgingVhostGroup
	mg.ListByOffset(&groupList, nil, 0, -1)
	ctx.Set(`listData`, rowAndGroup)
	ctx.Set(`groupList`, groupList)
	ctx.Set(`groupId`, groupID)
	currentHost := ctx.Host()
	ctx.SetFunc(`generateHostURL`, func(hosts string) []template.HTML {
		return generateHostURL(currentHost, hosts)
	})
	return ctx.Render(`caddy/vhost`, handler.Err(ctx, err))
}

var reSplitRegexp = regexp.MustCompile(`[\s]+`)

func hasEnvVar(v string) bool {
	for _, r := range v {
		if r == '$' || r == '%' {
			return true
		}
	}
	return false
}

func generateHostURL(currentHost string, hosts string) []template.HTML {
	hosts = strings.TrimSpace(hosts)
	hostsSlice := reSplitRegexp.Split(hosts, -1)
	urls := make([]template.HTML, 0, len(hostsSlice))
	for _, v := range hostsSlice {
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		if hasEnvVar(v) {
			urls = append(urls, template.HTML(v))
		} else {
			switch {
			case v[0] == ':':
				urls = append(urls, template.HTML(`<a href="http://`+currentHost+v+`" target="_blank">`+v+`</a>`))
			case strings.HasPrefix(v, `0.0.0.0:`):
				urls = append(urls, template.HTML(`<a href="http://`+currentHost+strings.TrimPrefix(v, `0.0.0.0`)+`" target="_blank">`+v+`</a>`))
			case !strings.Contains(v, `//`):
				urls = append(urls, template.HTML(`<a href="http://`+v+`" target="_blank">`+v+`</a>`))
			default:
				urls = append(urls, template.HTML(`<a href="`+strings.ReplaceAll(v, `*`, `test`)+`" target="_blank">`+v+`</a>`))
			}
		}
	}
	return urls
}

func Vhostbuild(ctx echo.Context) error {
	saveDir, err := getSaveDir()
	if err == nil {
		err = filepath.Walk(saveDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, `.conf`) {
				return nil
			}
			log.Info(`Delete the Caddy configuration file: `, path)
			return os.Remove(path)
		})
	}
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
	}
	m := model.NewVhost(ctx)
	n := 100
	cnt, err := m.ListByOffset(nil, nil, 0, n, `disabled`, `N`)
	if err != nil {
		return err
	}
	for i, j := 0, cnt(); int64(i) < j; i += n {
		if i > 0 {
			_, err = m.ListByOffset(nil, nil, i, n, `disabled`, `N`)
			if err != nil {
				handler.SendFail(ctx, err.Error())
				return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
			}
		}
		for _, m := range m.Objects() {
			var formData url.Values
			err := json.Unmarshal([]byte(m.Setting), &formData)
			if err == nil {
				file := filepath.Join(saveDir, fmt.Sprint(m.Id)+`.conf`)
				err = saveVhostConf(ctx, file, formData)
			}
			if err != nil {
				handler.SendFail(ctx, err.Error())
				return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
			}
		}
	}
	err = cmder.GetCaddyCmd().ReloadServer()
	if err != nil {
		ctx.Logger().Error(err)
	}
	handler.SendOk(ctx, ctx.T(`操作成功`))
	return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
}

func VhostAdd(ctx echo.Context) error {
	var err error
	m := model.NewVhost(ctx)
	if ctx.IsPost() {
		m.Domain = ctx.Form(`domain`)
		m.Disabled = ctx.Form(`disabled`)
		m.Root = ctx.Form(`root`)
		m.Name = ctx.Form(`name`)
		m.GroupId = ctx.Formx(`groupId`).Uint()
		var b []byte
		b, err = json.Marshal(ctx.Forms())
		switch {
		case err == nil:
			m.Setting = string(b)
			_, err = m.Insert()
			if err != nil {
				break
			}
			fallthrough
		case 0 == 1:
			err = saveVhostData(ctx, m.NgingVhost, ctx.Forms(), false)
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
	ctx.SetFunc(`Val`, func(name, defaultValue string) string {
		return ctx.Form(name, defaultValue)
	})
	g := dbschema.NewNgingVhostGroup(ctx)
	g.ListByOffset(nil, nil, 0, -1)
	ctx.Set(`groupList`, g.Objects())
	ctx.Set(`isAdd`, true)
	ctx.Set(`title`, ctx.T(`添加网站`))
	return ctx.Render(`caddy/vhost_edit`, err)
}

func getSaveDir() (saveDir string, err error) {
	cfg := cmder.GetCaddyConfig()
	saveDir = cfg.GetVhostConfigDirAbsPath()
	err = com.MkdirAll(saveDir, os.ModePerm)
	return
}

func saveVhostConf(ctx echo.Context, saveFile string, values url.Values) error {
	ctx.Set(`values`, NewFormValues(values))
	b, err := ctx.Fetch(`caddy/caddyfile`, nil)
	if err != nil {
		return err
	}
	b = com.CleanSpaceLine(b)
	log.Info(`Generate a Caddy configuration file: `, saveFile)
	err = os.WriteFile(saveFile, b, os.ModePerm)
	//jsonb, _ := caddyfile.ToJSON(b)
	//err = os.WriteFile(saveFile+`.json`, jsonb, os.ModePerm)
	return err
}

func saveVhostData(ctx echo.Context, m *dbschema.NgingVhost, values url.Values, restart bool) (err error) {
	var saveDir string
	saveDir, err = getSaveDir()
	if err != nil {
		return
	}
	saveFile := filepath.Join(saveDir, fmt.Sprint(m.Id)+`.conf`)
	if m.Disabled == `Y` {
		err = os.Remove(saveFile)
		if os.IsNotExist(err) {
			err = nil
		}
	} else {
		err = saveVhostConf(ctx, saveFile, values)
	}
	if err == nil && restart {
		err = cmder.GetCaddyCmd().ReloadServer()
	}
	return
}

func VhostDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	if id < 1 {
		handler.SendFail(ctx, ctx.T(`id无效`))
		return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
	}
	m := model.NewVhost(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err != nil {
		handler.SendFail(ctx, err.Error())
	} else {
		err = DeleteCaddyfileByID(id)
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
		}
	}
	return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
}

func DeleteCaddyfileByID(id uint) error {
	cfg := cmder.GetCaddyConfig()
	saveDir := cfg.GetVhostConfigDirAbsPath()
	saveFile := filepath.Join(saveDir, fmt.Sprint(id)+`.conf`)
	err := os.Remove(saveFile)
	if err == nil {
		err = cmder.GetCaddyCmd().ReloadServer()
	} else if os.IsNotExist(err) {
		err = nil
	}
	return err
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
		m.Domain = ctx.Form(`domain`)
		m.Disabled = ctx.Form(`disabled`)
		m.Root = ctx.Form(`root`)
		m.Name = ctx.Form(`name`)
		m.GroupId = ctx.Formx(`groupId`).Uint()
		var b []byte
		b, err = json.Marshal(ctx.Forms())
		switch {
		case err == nil:
			m.Setting = string(b)
			err = m.Update(nil, db.Cond{`id`: id})
			if err != nil {
				break
			}
			fallthrough
		case 0 == 1:
			removeCachedCert := ctx.Form(`removeCachedCert`)
			if len(removeCachedCert) > 0 && removeCachedCert == `1` {
				m.RemoveCachedCert()
			}
			err = saveVhostData(ctx, m.NgingVhost, ctx.Forms(), len(ctx.Form(`restart`)) > 0)
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
		}
	} else if ctx.IsAjax() {
		data := ctx.Data()
		disabled := ctx.Query(`disabled`)
		if len(disabled) > 0 {
			m.Disabled = disabled
			err = m.UpdateField(nil, `disabled`, disabled, db.Cond{`id`: id})
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			if m.Disabled == `Y` {
				err = DeleteCaddyfileByID(id)
			} else {
				var formData url.Values
				err = json.Unmarshal([]byte(m.Setting), &formData)
				if err == nil {
					err = saveVhostData(ctx, m.NgingVhost, formData, true)
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
	ctx.SetFunc(`Val`, func(name, defaultValue string) string {
		return ctx.Form(name, defaultValue)
	})
	ctx.Set(`activeURL`, `/caddy/vhost`)
	g := dbschema.NewNgingVhostGroup(ctx)
	g.ListByOffset(nil, nil, 0, -1)
	ctx.Set(`groupList`, g.Objects())
	ctx.Set(`isAdd`, false)
	ctx.Set(`title`, ctx.T(`修改网站`))
	return ctx.Render(`caddy/vhost_edit`, err)
}
