/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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

package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/library/filemanager"
	"github.com/admpub/nging/v4/application/library/notice"
	"github.com/admpub/nging/v4/application/library/respond"
	uploadLibrary "github.com/admpub/nging/v4/application/library/upload"
	"github.com/admpub/nging/v4/application/registry/upload/chunk"

	uploadClient "github.com/webx-top/client/upload"
	uploadDropzone "github.com/webx-top/client/upload/driver/dropzone"
)

// UploadedFile 本地附件文件管理
func UploadedFile(ctx echo.Context) error {
	ctx.Set(`activeURL`, `/manager/uploaded/file`)
	return Uploaded(ctx, `file`)
}

func UploadedChunk(ctx echo.Context) error {
	ctx.Set(`activeURL`, `/manager/uploaded/chunk`)
	return Uploaded(ctx, `chunk`)
}

func UploadedMerged(ctx echo.Context) error {
	ctx.Set(`activeURL`, `/manager/uploaded/merged`)
	return Uploaded(ctx, `merged`)
}

// Uploaded 本地附件文件管理
func Uploaded(ctx echo.Context, uploadType string) error {
	var (
		root      string
		canUpload bool
		canEdit   bool
		canDelete bool
	)
	switch uploadType {
	case `chunk`:
		root = chunk.ChunkTempDir
		canDelete = true
	case `merged`:
		root = chunk.MergeSaveDir
		canDelete = true
	case `file`:
		canUpload = true
		canEdit = true
		canDelete = true
		root = uploadLibrary.UploadDir
	default:
		return echo.ErrNotFound
	}
	var err error
	id := ctx.Formx(`id`).Uint()
	filePath := ctx.Form(`path`)
	do := ctx.Form(`do`)
	mgr := filemanager.New(root, config.DefaultConfig.Sys.EditableFileMaxBytes(), ctx)
	absPath := root
	if err == nil && len(root) > 0 {

		if len(filePath) > 0 {
			filePath = filepath.Clean(filePath)
			absPath = filepath.Join(root, filePath)
		}

		user := handler.User(ctx)
		switch do {
		case `edit`:
			if !canEdit {
				return echo.ErrNotFound
			}
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
			if !canEdit {
				return echo.ErrNotFound
			}
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
			if !canEdit {
				return echo.ErrNotFound
			}
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
			if !canDelete {
				return echo.ErrNotFound
			}
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
			if !canUpload {
				return echo.ErrNotFound
			}
			var cu *uploadClient.ChunkUpload
			var opts []uploadClient.ChunkInfoOpter
			if user != nil {
				_cu := chunk.ChunkUploader()
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
	if filePath == `.` {
		filePath = ``
	}
	gallery := ctx.Formx(`gallery`).Bool()
	var urlParam string
	if gallery {
		urlParam += `&gallery=1`
	}
	pathSlice := strings.Split(strings.Trim(filePath, echo.FilePathSeparator), echo.FilePathSeparator)
	pathLinks := make(echo.KVList, len(pathSlice))
	encodedSep := filemanager.EncodedSepa
	urlPrefix := fmt.Sprintf(`/manager/uploaded/`+uploadType+`?id=%d`+urlParam+`&path=`, id) + encodedSep
	for k, v := range pathSlice {
		urlPrefix += com.URLEncode(v)
		pathLinks[k] = &echo.KV{K: v, V: urlPrefix}
		urlPrefix += encodedSep
	}
	ctx.Set(`pathLinks`, pathLinks)
	ctx.Set(`rootPath`, strings.TrimSuffix(root, echo.FilePathSeparator))
	ctx.Set(`path`, filePath)
	ctx.Set(`absPath`, absPath)

	ctx.Set(`uploadType`, uploadType)
	ctx.Set(`canUpload`, canUpload)
	ctx.Set(`canEdit`, canEdit)
	ctx.Set(`canDelete`, canDelete)

	ctx.SetFunc(`Editable`, func(fileName string) bool {
		if !canEdit {
			return false
		}
		_, ok := Editable(fileName)
		return ok
	})
	ctx.SetFunc(`Playable`, func(fileName string) string {
		mime, _ := Playable(fileName)
		return mime
	})
	ctx.SetFunc(`URLPrefix`, func() string {
		return `/manager/uploaded/` + uploadType
	})
	if gallery {
		return ctx.Render(`manager/uploaded_photo`, err)
	}
	return ctx.Render(`manager/uploaded_file`, err)
}

func Editable(fileName string) (string, bool) {
	return config.DefaultConfig.Sys.Editable(fileName)
}

func Playable(fileName string) (string, bool) {
	return config.DefaultConfig.Sys.Playable(fileName)
}
