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
package download

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/admpub/godownloader/service"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/handler/caddy"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/filemanager"
	"github.com/admpub/nging/application/library/notice"
	"github.com/webx-top/echo"
)

func init() {
	server := &service.DServ{}
	server.SetTmpl(`download/index`)
	server.SetSavePath(func() string {
		if len(config.DefaultConfig.Download.SavePath) == 0 {
			return service.GetDownloadPath()
		}
		return config.DefaultConfig.Download.SavePath
	})
	handler.RegisterToGroup(`/download`, func(g *echo.Group) {
		server.Register(g, true)
		g.Route(`GET,POST`, `/file`, File)
	})
}

func File(ctx echo.Context) error {
	var err error
	filePath := ctx.Form(`path`)
	do := ctx.Form(`do`)
	var root string
	if len(config.DefaultConfig.Download.SavePath) == 0 {
		root = service.GetDownloadPath()
	} else {
		root = config.DefaultConfig.Download.SavePath
	}
	mgr := filemanager.New(root, config.DefaultConfig.Sys.EditableFileMaxBytes, ctx)
	absPath := root
	var exit bool

	if len(filePath) > 0 {
		filePath = filepath.Clean(filePath)
		absPath = filepath.Join(root, filePath)
	}

	switch do {
	case `edit`:
		data := ctx.NewData()
		if _, ok := caddy.Editable(absPath); !ok {
			data.Info = errors.New(ctx.T(`此文件不能在线编辑`))
		} else {
			content := ctx.Form(`content`)
			encoding := ctx.Form(`encoding`)
			dat, err := mgr.Edit(absPath, content, encoding)
			if err != nil {
				data.Info = err.Error()
			} else {
				data.Code = 1
				data.Data = dat
			}
		}
		return ctx.JSON(data)
	case `rename`:
		data := ctx.NewData()
		newName := ctx.Form(`name`)
		err = mgr.Rename(absPath, newName)
		if err != nil {
			data.Info = err.Error()
		} else {
			data.Code = 1
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
			return err
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
	if filePath == `.` {
		filePath = ``
	}
	ctx.Set(`path`, filePath)
	ctx.Set(`absPath`, absPath)
	ctx.SetFunc(`Editable`, func(fileName string) bool {
		_, ok := caddy.Editable(fileName)
		return ok
	})
	return ctx.Render(`download/file`, err)
}
