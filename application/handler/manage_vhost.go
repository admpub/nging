/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"

	"github.com/admpub/caddyui/application/library/config"
	"github.com/admpub/caddyui/application/library/modal"
	"github.com/admpub/caddyui/application/library/notice"
	"github.com/admpub/caddyui/application/model"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func ManageIndex(ctx echo.Context) error {
	m := model.NewVhost(ctx)
	page, size := Paging(ctx)
	cnt, err := m.List(nil, nil, page, size)
	var ret interface{}
	if err == nil {
		flash := ctx.Flash()
		if flash != nil {
			if errMsg, ok := flash.(string); ok {
				ret = errors.New(errMsg)
			} else {
				ret = flash
			}
		}
	} else {
		ret = err
	}
	ctx.SetFunc(`totalRows`, cnt)
	ctx.Set(`listData`, m.Objects())
	return ctx.Render(`manage/index`, ret)
}

func ManageVhostAdd(ctx echo.Context) error {
	var err error
	if ctx.IsPost() {
		m := model.NewVhost(ctx)
		m.Domain = ctx.Form(`domain`)
		m.Disabled = ctx.Form(`disabled`)
		m.Root = ctx.Form(`root`)
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
			err = saveVhostData(ctx, m)
		}
		if err == nil {
			ctx.Session().AddFlash(Ok(ctx.T(`操作成功`)))
			return ctx.Redirect(`/manage`)
		}

		ctx.SetFunc(`Val`, func(name, defaultValue string) string {
			return ctx.Form(name)
		})
	} else {
		ctx.SetFunc(`Val`, func(name, defaultValue string) string {
			return defaultValue
		})
	}
	return ctx.Render(`manage/vhost_edit`, err)
}

func saveVhostData(ctx echo.Context, m *model.Vhost) (err error) {
	var b []byte
	var saveFile string
	SetCaddyfileFunc(ctx)
	b, err = ctx.Fetch(`manage/caddyfile`, nil)
	if err != nil {
		return
	}
	saveFile, err = filepath.Abs(config.DefaultConfig.Sys.VhostsfileDir)
	if err != nil {
		return
	}
	saveFile = filepath.Join(saveFile, fmt.Sprint(m.Id))
	if m.Disabled == `Y` {
		err = os.Remove(saveFile)
		if os.IsNotExist(err) {
			err = nil
		}
	} else {
		err = ioutil.WriteFile(saveFile, b, os.ModePerm)
		if len(ctx.Form(`restart`)) > 0 {
			err = config.DefaultCLIConfig.CaddyRestart()
		}
	}
	return
}

func ManageVhostDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	if id < 1 {
		ctx.Session().AddFlash(ctx.T(`id无效`))
		return ctx.Redirect(`/manage`)
	}
	m := model.NewVhost(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err != nil {
		ctx.Session().AddFlash(err)
	} else {
		var saveFile string
		saveFile, err = filepath.Abs(config.DefaultConfig.Sys.VhostsfileDir)
		if err == nil {
			saveFile = filepath.Join(saveFile, fmt.Sprint(m.Id))
			err = os.Remove(saveFile)
			if os.IsNotExist(err) {
				err = nil
			}
			if err == nil {
				ctx.Session().AddFlash(Ok(ctx.T(`操作成功`)))
			}
		}
	}
	return ctx.Redirect(`/manage`)
}

func ManageVhostEdit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	if id < 1 {
		ctx.Session().AddFlash(ctx.T(`id无效`))
		return ctx.Redirect(`/manage`)
	}

	var err error
	m := model.NewVhost(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		ctx.Session().AddFlash(err.Error())
		return ctx.Redirect(`/manage`)
	}
	if ctx.IsPost() {
		m.Domain = ctx.Form(`domain`)
		m.Disabled = ctx.Form(`disabled`)
		m.Root = ctx.Form(`root`)
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
			err = saveVhostData(ctx, m)
		}
		if err == nil {
			ctx.Session().AddFlash(Ok(ctx.T(`操作成功`)))
			return ctx.Redirect(`/manage`)
		}
	} else {
		var formData url.Values
		if e := json.Unmarshal([]byte(m.Setting), &formData); e == nil {
			for key, values := range formData {
				for _, v := range values {
					ctx.Request().Form().Add(key, v)
				}
			}
		}
	}
	ctx.SetFunc(`Val`, func(name, defaultValue string) string {
		return ctx.Form(name)
	})
	ctx.Set(`activeURL`, `/manage`)
	return ctx.Render(`manage/vhost_edit`, err)
}

func ManageClearCache(ctx echo.Context) error {
	if err := modal.Clear(); err != nil {
		return err
	}
	notice.Clear()
	return ctx.String(ctx.T(`已经清理完毕`))
}

func ManageVhostFile(ctx echo.Context) error {
	var err error
	vhostId := ctx.Formx(`id`).Uint()
	file := ctx.Form(`file`)
	do := ctx.Form(`do`)
	m := model.NewVhost(ctx)
	err = m.Get(nil, db.Cond{`id`: vhostId})
	absFile := m.Root
	if err == nil && len(m.Root) > 0 {
		var exit bool
		switch do {
		case `edit`:
		case `delete`:
		default:
			if len(file) > 0 {
				file = filepath.Clean(file)
				absFile = filepath.Join(m.Root, file)
			}
			err, exit = fileList(ctx, absFile)
		}
		if exit {
			return err
		}
	}
	ctx.Set(`data`, m)
	if file == `.` {
		file = ``
	}
	ctx.Set(`file`, file)
	ctx.Set(`absFile`, absFile)
	ctx.Set(`activeURL`, `/manage`)
	return ctx.Render(`manage/file`, err)
}

func fileList(ctx echo.Context, absFile string) (err error, exit bool) {
	fs := http.Dir(filepath.Dir(absFile))
	var d http.File
	fileName := filepath.Base(absFile)
	d, err = fs.Open(fileName)
	if err == nil {
		defer d.Close()
		var fi os.FileInfo
		fi, err = d.Stat()
		if err == nil {
			if !fi.IsDir() {
				return ctx.Attachment(d, fileName), true
			}
			var dirs []os.FileInfo
			dirs, err = d.Readdir(-1)
			sort.Sort(byFileType(dirs))
			ctx.Set(`dirs`, dirs)
		}
	}
	return
}

type byFileType []os.FileInfo

func (s byFileType) Len() int { return len(s) }
func (s byFileType) Less(i, j int) bool {
	if s[i].IsDir() {
		if !s[j].IsDir() {
			return true
		}
	} else if s[j].IsDir() {
		if !s[i].IsDir() {
			return false
		}
	}
	return s[i].Name() < s[j].Name()
}
func (s byFileType) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
