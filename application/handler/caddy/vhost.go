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
package caddy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/filemanager"
	"github.com/admpub/nging/application/library/notice"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"
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
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()))
	ret := handler.Err(ctx, err)
	rows := m.Objects()
	gIds := []uint{}
	rowAndGroup := make([]*model.VhostAndGroup, len(rows))
	for k, u := range rows {
		rowAndGroup[k] = &model.VhostAndGroup{
			Vhost: u,
		}
		if u.GroupId < 1 {
			continue
		}
		if !com.InUintSlice(u.GroupId, gIds) {
			gIds = append(gIds, u.GroupId)
		}
	}

	mg := &dbschema.VhostGroup{}
	var groupList []*dbschema.VhostGroup
	if len(gIds) > 0 {
		_, err = mg.List(&groupList, nil, 1, 1000, db.Cond{`id IN`: gIds})
		if err != nil {
			if ret == nil {
				ret = err
			}
		} else {
			for k, v := range rowAndGroup {
				for _, g := range groupList {
					if g.Id == v.GroupId {
						rowAndGroup[k].Group = g
						break
					}
				}
			}
		}
	}
	ctx.Set(`listData`, rowAndGroup)
	mg.ListByOffset(&groupList, nil, 0, -1)
	ctx.Set(`groupList`, groupList)
	ctx.Set(`groupId`, groupID)
	return ctx.Render(`caddy/vhost`, ret)
}

func Vhostbuild(ctx echo.Context) error {
	saveFile, err := getSaveDir()
	if err == nil {
		err = filepath.Walk(saveFile, func(path string, info os.FileInfo, err error) error {
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
				file := filepath.Join(saveFile, fmt.Sprint(m.Id)+`.conf`)
				err = saveVhostConf(ctx, file, formData)
			}
			if err != nil {
				handler.SendFail(ctx, err.Error())
				return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
			}
		}
	}
	err = config.DefaultCLIConfig.CaddyReload()
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
			_, err = m.Add()
			if err != nil {
				break
			}
			fallthrough
		case 0 == 1:
			err = saveVhostData(ctx, m.Vhost, ctx.Forms(), false)
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
	g := &dbschema.VhostGroup{}
	g.ListByOffset(nil, nil, 0, -1)
	ctx.Set(`groupList`, g.Objects())
	return ctx.Render(`caddy/vhost_edit`, err)
}

func getSaveDir() (saveFile string, err error) {
	saveFile, err = filepath.Abs(config.DefaultConfig.Sys.VhostsfileDir)
	if err != nil {
		return
	}
	if fi, er := os.Stat(saveFile); er != nil || !fi.IsDir() {
		err = os.MkdirAll(saveFile, 0666)
		if err != nil {
			return
		}
	}
	return
}

func saveVhostConf(ctx echo.Context, saveFile string, values url.Values) error {
	SetCaddyfileFunc(ctx, values)
	ctx.Set(`values`, values)
	b, err := ctx.Fetch(`caddy/caddyfile`, nil)
	if err != nil {
		return err
	}
	b = com.CleanSpaceLine(b)
	log.Info(`Generate a Caddy configuration file: `, saveFile)
	err = ioutil.WriteFile(saveFile, b, os.ModePerm)
	return err
}

func saveVhostData(ctx echo.Context, m *dbschema.Vhost, values url.Values, restart bool) (err error) {
	var saveFile string
	saveFile, err = getSaveDir()
	if err != nil {
		return
	}
	saveFile = filepath.Join(saveFile, fmt.Sprint(m.Id)+`.conf`)
	if m.Disabled == `Y` {
		err = os.Remove(saveFile)
		if os.IsNotExist(err) {
			err = nil
		}
	} else {
		err = saveVhostConf(ctx, saveFile, values)
	}
	if err == nil && restart {
		err = config.DefaultCLIConfig.CaddyReload()
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
	saveFile, err := filepath.Abs(config.DefaultConfig.Sys.VhostsfileDir)
	if err == nil {
		saveFile = filepath.Join(saveFile, fmt.Sprint(id)+`.conf`)
		err = os.Remove(saveFile)
		if err == nil {
			err = config.DefaultCLIConfig.CaddyReload()
		} else if os.IsNotExist(err) {
			err = nil
		}
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
			err = m.Edit(nil, db.Cond{`id`: id})
			if err != nil {
				break
			}
			fallthrough
		case 0 == 1:
			err = saveVhostData(ctx, m.Vhost, ctx.Forms(), len(ctx.Form(`restart`)) > 0)
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/caddy/vhost`))
		}
	} else if ctx.IsAjax() {
		disabled := ctx.Query(`disabled`)
		if len(disabled) > 0 {
			m.Disabled = disabled
			data := ctx.Data()
			err = m.SetField(nil, `disabled`, disabled, db.Cond{`id`: id})
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
					for name, fun := range tplfunc.TplFuncMap {
						ctx.SetFunc(name, fun)
					}
					err = saveVhostData(ctx, m.Vhost, formData, true)
				}
			}
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`操作成功`))
			return ctx.JSON(data)
		}
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
	}
	ctx.SetFunc(`Val`, func(name, defaultValue string) string {
		return ctx.Form(name, defaultValue)
	})
	ctx.Set(`activeURL`, `/caddy/vhost`)
	g := &dbschema.VhostGroup{}
	g.ListByOffset(nil, nil, 0, -1)
	ctx.Set(`groupList`, g.Objects())
	return ctx.Render(`caddy/vhost_edit`, err)
}

func VhostFile(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	filePath := ctx.Form(`path`)
	do := ctx.Form(`do`)
	m := model.NewVhost(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	mgr := filemanager.New(m.Root, config.DefaultConfig.Sys.EditableFileMaxBytes, ctx)
	absPath := m.Root
	if err == nil && len(m.Root) > 0 {
		var exit bool

		if len(filePath) > 0 {
			filePath = filepath.Clean(filePath)
			absPath = filepath.Join(m.Root, filePath)
		}

		switch do {
		case `edit`:
			data := ctx.Data()
			if _, ok := Editable(absPath); !ok {
				data.SetInfo(ctx.T(`此文件不能在线编辑`), 0)
			} else {
				content := ctx.Form(`content`)
				encoding := ctx.Form(`encoding`)
				dat, err := mgr.Edit(absPath, content, encoding)
				if err != nil {
					data.SetInfo(err.Error(), 0)
				} else {
					data.SetData(dat, 1)
				}
			}
			return ctx.JSON(data)
		case `rename`:
			data := ctx.Data()
			newName := ctx.Form(`name`)
			err = mgr.Rename(absPath, newName)
			if err != nil {
				data.SetInfo(err.Error(), 0)
			} else {
				data.SetCode(1)
			}
			return ctx.JSON(data)
		case `mkdir`:
			data := ctx.Data()
			newName := ctx.Form(`name`)
			err = mgr.Mkdir(filepath.Join(absPath, newName), os.ModePerm)
			if err != nil {
				data.SetInfo(err.Error(), 0)
			} else {
				data.SetCode(1)
			}
			return ctx.JSON(data)
		case `delete`:
			err = mgr.Remove(absPath)
			if err != nil {
				handler.SendFail(ctx, err.Error())
			}
			return ctx.Redirect(ctx.Referer())
		case `upload`:
			err = mgr.Upload(absPath)
			if err != nil {
				user := handler.User(ctx)
				if user != nil {
					notice.OpenMessage(user.Username, `upload`)
					notice.Send(user.Username, notice.NewMessageWithValue(`upload`, ctx.T(`文件上传出错`), err.Error()))
				}
				return ctx.JSON(echo.H{`error`: err.Error()}, 500)
			}
			return ctx.String(`OK`)
		default:
			var dirs []os.FileInfo
			err, exit, dirs = mgr.List(absPath)
			ctx.Set(`dirs`, dirs)
		}
		if exit {
			return err
		}
	}
	ctx.Set(`data`, m)
	if filePath == `.` {
		filePath = ``
	}
	pathSlice := strings.Split(strings.Trim(filePath, echo.FilePathSeparator), echo.FilePathSeparator)
	pathLinks := make(echo.KVList, len(pathSlice))
	encodedSep := filemanager.EncodedSepa
	urlPrefix := fmt.Sprintf(`/caddy/vhost_file?id=%d&path=`, id) + encodedSep
	for k, v := range pathSlice {
		urlPrefix += com.URLEncode(v)
		pathLinks[k] = &echo.KV{K: v, V: urlPrefix}
		urlPrefix += encodedSep
	}
	ctx.Set(`pathLinks`, pathLinks)
	ctx.Set(`rootPath`, strings.TrimSuffix(m.Root, echo.FilePathSeparator))
	ctx.Set(`path`, filePath)
	ctx.Set(`absPath`, absPath)
	ctx.SetFunc(`Editable`, func(fileName string) bool {
		_, ok := Editable(fileName)
		return ok
	})
	ctx.Set(`activeURL`, `/caddy/vhost`)
	return ctx.Render(`caddy/file`, err)
}

func Editable(fileName string) (string, bool) {
	ext := strings.TrimPrefix(filepath.Ext(fileName), `.`)
	ext = strings.ToLower(ext)
	typ, ok := config.DefaultConfig.Sys.EditableFileExtensions[ext]
	return typ, ok
}
