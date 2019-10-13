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

package download

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	mw "github.com/webx-top/echo/middleware"

	"github.com/admpub/godownloader/service"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/handler/caddy"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/dbmanager/driver/mysql"
	"github.com/admpub/nging/application/library/filemanager"
	"github.com/admpub/nging/application/library/notice"
)

var downloadDir = func() string {
	if len(config.DefaultConfig.Download.SavePath) == 0 {
		return service.GetDownloadPath()
	}
	return config.DefaultConfig.Download.SavePath
}

func init() {
	server := &service.DServ{}
	server.SetTmpl(`download/index`)
	server.SetSavePath(downloadDir)
	mysql.SQLTempDir = downloadDir //将SQL文件缓存到下载目录里面方便管理
	handler.RegisterToGroup(`/download`, func(g echo.RouteRegister) {
		server.Register(g, true)
		g.Route(`GET,POST`, `/file`, File, mw.CORS())
	})
}

func File(ctx echo.Context) error {
	var err error
	filePath := ctx.Form(`path`)
	do := ctx.Form(`do`)
	root := downloadDir()
	mgr := filemanager.New(root, config.DefaultConfig.Sys.EditableFileMaxBytes, ctx)
	absPath := root

	if len(filePath) > 0 {
		filePath = filepath.Clean(filePath)
		absPath = filepath.Join(root, filePath)
	}

	switch do {
	case `edit`:
		data := ctx.Data()
		if _, ok := caddy.Editable(absPath); !ok {
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
		paths := ctx.FormValues(`path`)
		for _, filePath := range paths {
			filePath = strings.TrimSpace(filePath)
			if len(filePath) == 0 {
				continue
			}
			filePath = filepath.Clean(filePath)
			absPath = filepath.Join(root, filePath)
			err = mgr.Remove(absPath)
			if err != nil {
				handler.SendFail(ctx, err.Error())
				return ctx.Redirect(ctx.Referer())
			}
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
		var exit bool
		err, exit, dirs = mgr.List(absPath)
		if exit {
			return err
		}
		ctx.Set(`dirs`, dirs)
	}
	if filePath == `.` {
		filePath = ``
	}
	pathSlice := strings.Split(strings.Trim(filePath, echo.FilePathSeparator), echo.FilePathSeparator)
	pathLinks := make(echo.KVList, len(pathSlice))
	encodedSep := filemanager.EncodedSepa
	urlPrefix := handler.URLFor(`/download/file?path=` + encodedSep)
	for k, v := range pathSlice {
		urlPrefix += com.URLEncode(v)
		pathLinks[k] = &echo.KV{K: v, V: urlPrefix}
		urlPrefix += encodedSep
	}
	ctx.Set(`pathLinks`, pathLinks)
	ctx.Set(`rootPath`, strings.TrimSuffix(root, echo.FilePathSeparator))
	ctx.Set(`path`, filePath)
	ctx.Set(`absPath`, absPath)
	ctx.SetFunc(`Editable`, func(fileName string) bool {
		_, ok := caddy.Editable(fileName)
		return ok
	})
	ctx.SetFunc(`Playable`, func(fileName string) string {
		mime, _ := caddy.Playable(fileName)
		return mime
	})
	return ctx.Render(`download/file`, err)
}
