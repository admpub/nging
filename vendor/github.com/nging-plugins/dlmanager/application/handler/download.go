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

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/filemanager"
	"github.com/admpub/nging/v5/application/library/notice"
	"github.com/admpub/nging/v5/application/library/respond"

	uploadChunk "github.com/admpub/nging/v5/application/registry/upload/chunk"
	uploadClient "github.com/webx-top/client/upload"
	uploadDropzone "github.com/webx-top/client/upload/driver/dropzone"
)

func File(ctx echo.Context) error {
	var err error
	filePath := ctx.Form(`path`)
	do := ctx.Form(`do`)
	root := downloadDir()
	mgr := filemanager.New(root, config.FromFile().Sys.EditableFileMaxBytes(), ctx)
	absPath := root

	if len(filePath) > 0 {
		filePath = filepath.Clean(filePath)
		absPath = filepath.Join(root, filePath)
	}

	user := handler.User(ctx)
	switch do {
	case `edit`:
		data := ctx.Data()
		if _, ok := config.FromFile().Sys.Editable(absPath); !ok {
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
		next := ctx.Referer()
		if len(next) == 0 {
			next = ctx.Request().URL().Path() + fmt.Sprintf(`?path=%s`, com.URLEncode(filepath.Dir(filePath)))
		}
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
				return ctx.Redirect(next)
			}
		}
		return ctx.Redirect(next)
	case `upload`:
		var cu *uploadClient.ChunkUpload
		var opts []uploadClient.ChunkInfoOpter
		if user != nil {
			cu = uploadChunk.NewUploader(fmt.Sprintf(`user/%d`, user.Id))
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
		_, ok := config.FromFile().Sys.Editable(fileName)
		return ok
	})
	ctx.SetFunc(`Playable`, func(fileName string) string {
		mime, _ := config.FromFile().Sys.Playable(fileName)
		return mime
	})
	return ctx.Render(`download/file`, err)
}
