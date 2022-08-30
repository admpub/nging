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
	"os"
	"path/filepath"
	"strings"

	uploadClient "github.com/webx-top/client/upload"
	uploadDropzone "github.com/webx-top/client/upload/driver/dropzone"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/library/filemanager"
	"github.com/admpub/nging/v4/application/library/notice"
	"github.com/admpub/nging/v4/application/library/respond"
	uploadChunk "github.com/admpub/nging/v4/application/registry/upload/chunk"

	"github.com/nging-plugins/caddymanager/application/model"
)

func VhostFile(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	filePath := ctx.Form(`path`)
	do := ctx.Form(`do`)
	m := model.NewVhost(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	mgr := filemanager.New(m.Root, config.FromFile().Sys.EditableFileMaxBytes(), ctx)
	absPath := m.Root
	user := handler.User(ctx)
	if err == nil && len(m.Root) > 0 {

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
			next := ctx.Referer()
			if len(next) == 0 {
				next = ctx.Request().URL().Path() + fmt.Sprintf(`?id=%d&path=%s`, id, com.URLEncode(filepath.Dir(filePath)))
			}
			return ctx.Redirect(next)
		case `upload`:
			var cu *uploadClient.ChunkUpload
			var opts []uploadClient.ChunkInfoOpter
			if user != nil {
				_cu := uploadChunk.ChunkUploader()
				_cu.UID = fmt.Sprintf(`user/%d`, user.Id)
				cu = &_cu
				opts = append(opts, uploadClient.OptChunkInfoMapping(uploadDropzone.MappingChunkInfo))
			}
			err = mgr.Upload(absPath, cu, opts...)
			if err != nil {
				user := handler.User(ctx)
				if user != nil {
					notice.OpenMessage(user.Username, `upload`)
					notice.Send(user.Username, notice.NewMessageWithValue(`upload`, ctx.T(`文件上传出错`), err.Error()))
				}
			}
			return respond.Dropzone(ctx, err, nil)
		default:
			var dirs []os.FileInfo
			var exit bool
			err, exit, dirs = mgr.List(absPath)
			if exit {
				return err
			}
			ctx.Set(`dirs`, dirs)
		}
	}
	ctx.Set(`data`, m)
	if filePath == `.` {
		filePath = ``
	}
	pathSlice := strings.Split(strings.Trim(filePath, echo.FilePathSeparator), echo.FilePathSeparator)
	pathLinks := make(echo.KVList, len(pathSlice))
	encodedSep := filemanager.EncodedSepa
	urlPrefix := ctx.Request().URL().Path() + fmt.Sprintf(`?id=%d&path=`, id) + encodedSep
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
	ctx.SetFunc(`Playable`, func(fileName string) string {
		mime, _ := Playable(fileName)
		return mime
	})
	ctx.Set(`activeURL`, `/caddy/vhost`)
	return ctx.Render(`caddy/file`, err)
}

func Editable(fileName string) (string, bool) {
	return config.FromFile().Sys.Editable(fileName)
}

func Playable(fileName string) (string, bool) {
	return config.FromFile().Sys.Playable(fileName)
}
