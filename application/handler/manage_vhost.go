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
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func ManageIndex(ctx echo.Context) error {
	m := model.NewVhost(ctx)
	page, size := Paging(ctx)
	cnt, err := m.List(nil, nil, page, size)
	ret := Err(ctx, err)
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
	filePath := ctx.Form(`path`)
	do := ctx.Form(`do`)
	m := model.NewVhost(ctx)
	err = m.Get(nil, db.Cond{`id`: vhostId})
	absPath := m.Root
	if err == nil && len(m.Root) > 0 {
		var exit bool

		if len(filePath) > 0 {
			filePath = filepath.Clean(filePath)
			absPath = filepath.Join(m.Root, filePath)
		}

		switch do {
		case `edit`:
			data := ctx.NewData()
			dat, err := fileEdit(ctx, absPath)
			if err != nil {
				data.Info = err.Error()
			} else {
				data.Code = 1
				data.Data = dat
			}
			return ctx.JSON(data)
		case `rename`:
			data := ctx.NewData()
			newName := ctx.Form(`name`)
			if len(newName) > 0 {
				err = os.Rename(absPath, filepath.Join(filepath.Dir(absPath), filepath.Base(newName)))
			} else {
				err = errors.New(ctx.T(`请输入文件名称`))
			}
			if err != nil {
				data.Info = err.Error()
			} else {
				data.Code = 1
			}
			return ctx.JSON(data)
		case `delete`:
			err = fileRemove(absPath)
			if err != nil {
				ctx.Session().AddFlash(err)
			}
			return ctx.Redirect(ctx.Referer())
		case `upload`:
			err = fileUpload(ctx, absPath)
			if err != nil {
				user, _ := ctx.Get(`user`).(string)
				if len(user) > 0 {
					notice.OpenMessage(user, `upload`)
					notice.Send(user, notice.NewMessageWithValue(`upload`, ctx.T(`文件上传出错`), err.Error()))
				}
				return err
			}
			return ctx.String(`OK`)
		default:
			err, exit = fileList(ctx, absPath)
		}
		if exit {
			return err
		}
	}
	ctx.Set(`data`, m)
	if filePath == `.` {
		filePath = ``
	}
	ctx.Set(`path`, filePath)
	ctx.Set(`absPath`, absPath)
	ctx.Set(`activeURL`, `/manage`)
	return ctx.Render(`manage/file`, err)
}

func fileEdit(ctx echo.Context, absPath string) (interface{}, error) {
	fi, err := os.Stat(absPath)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return nil, errors.New(ctx.T(`不能编辑文件夹`))
	}
	if ctx.IsPost() {
		content := ctx.Form(`content`)
		err = ioutil.WriteFile(absPath, []byte(content), fi.Mode())
		return nil, err
	}
	b, err := ioutil.ReadFile(absPath)
	return string(b), err
}

func fileRemove(absPath string) error {
	fi, err := os.Stat(absPath)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return os.RemoveAll(absPath)
	}
	return os.Remove(absPath)
}

func enterPath(absPath string) (d http.File, fi os.FileInfo, err error) {
	fs := http.Dir(filepath.Dir(absPath))
	fileName := filepath.Base(absPath)
	d, err = fs.Open(fileName)
	if err != nil {
		return
	}
	//defer d.Close()
	fi, err = d.Stat()
	return
}

func fileUpload(ctx echo.Context, absPath string) (err error) {
	var (
		d  http.File
		fi os.FileInfo
	)
	d, fi, err = enterPath(absPath)
	if d != nil {
		defer d.Close()
	}
	if err != nil {
		return
	}
	if !fi.IsDir() {
		return errors.New(ctx.T(`路径不正确`))
	}
	pipe := ctx.Form(`pipe`)
	switch pipe {
	case `unzip`:
		fileHdr, err := ctx.SaveUploadedFile(`file`, absPath)
		if err != nil {
			return err
		}
		filePath := filepath.Join(absPath, fileHdr.Filename)
		err = com.Unzip(filePath, absPath)
		if err == nil {
			err = os.Remove(filePath)
			if err != nil {
				err = errors.New(ctx.T(`压缩包已经成功解压，但是删除压缩包失败：`) + err.Error())
			}
		}
		return err
	default:
		_, err = ctx.SaveUploadedFile(`file`, absPath)
	}
	return
}

func fileList(ctx echo.Context, absPath string) (err error, exit bool) {
	var (
		d  http.File
		fi os.FileInfo
	)
	d, fi, err = enterPath(absPath)
	if d != nil {
		defer d.Close()
	}
	if err != nil {
		return
	}
	if !fi.IsDir() {
		fileName := filepath.Base(absPath)
		return ctx.Attachment(d, fileName), true
	}

	var dirs []os.FileInfo
	dirs, err = d.Readdir(-1)
	sort.Sort(byFileType(dirs))
	ctx.Set(`dirs`, dirs)
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
