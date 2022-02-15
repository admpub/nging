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

package cloud

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/library/filemanager"
	"github.com/admpub/nging/v4/application/library/notice"
	"github.com/admpub/nging/v4/application/library/respond"
	"github.com/admpub/nging/v4/application/library/s3manager/s3client"
	"github.com/admpub/nging/v4/application/model"

	uploadChunk "github.com/admpub/nging/v4/application/registry/upload/chunk"
	uploadClient "github.com/webx-top/client/upload"
	uploadDropzone "github.com/webx-top/client/upload/driver/dropzone"
)

func StorageFile(ctx echo.Context) error {
	ctx.Set(`activeURL`, `/cloud/storage`)
	id := ctx.Formx(`id`).Uint()
	m := model.NewCloudStorage(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		return err
	}
	mgr, err := s3client.New(m.NgingCloudStorage, config.DefaultConfig.Sys.EditableFileMaxBytes())
	if err != nil {
		return err
	}

	ppath := ctx.Form(`path`)
	do := ctx.Form(`do`)
	var parentPath string
	if len(ppath) == 0 {
		ppath = `/`
	}
	if ppath != `/` {
		parentPath = strings.TrimSuffix(ppath, `/`)
		parentPath = path.Dir(parentPath)
	}
	if len(parentPath) > 0 && parentPath != `/` {
		parentPath += `/`
	}
	user := handler.User(ctx)
	switch do {
	case `edit`:
		data := ctx.Data()
		if _, ok := config.DefaultConfig.Sys.Editable(ppath); !ok {
			data.SetInfo(ctx.T(`此文件不能在线编辑`), 0)
		} else {
			content := ctx.Form(`content`)
			encoding := ctx.Form(`encoding`)
			dat, err := mgr.Edit(ctx, ppath, content, encoding)
			if err != nil {
				data.SetInfo(err.Error(), 0)
			} else {
				if ctx.IsPost() {
					data.SetInfo(ctx.T(`保存成功`), 1)
				}
				data.SetData(dat, 1)
			}
		}
		return ctx.JSON(data)
	case `mkdir`:
		data := ctx.Data()
		newName := ctx.Form(`name`)
		if len(newName) == 0 {
			data.SetInfo(ctx.T(`请输入文件夹名`), 0)
		} else {
			err = mgr.Mkdir(ctx, ppath, newName)
			if err != nil {
				data.SetError(err)
			}
			if data.GetCode() == 1 {
				data.SetInfo(ctx.T(`创建成功`))
			}
		}
		return ctx.JSON(data)
	case `rename`:
		data := ctx.Data()
		newName := ctx.Form(`name`)
		err = mgr.Rename(ctx, ppath, newName)
		if err != nil {
			data.SetInfo(err.Error(), 0)
		} else {
			data.SetInfo(ctx.T(`重命名成功`), 1)
		}
		return ctx.JSON(data)
	case `search`:
		prefix := ctx.Form(`query`)
		num := ctx.Formx(`size`, `10`).Int()
		if num <= 0 {
			num = 10
		}
		paths := mgr.Search(ctx, ppath, prefix, num)
		data := ctx.Data().SetData(paths)
		return ctx.JSON(data)
	case `delete`:
		err = mgr.Remove(ctx, ppath)
		if err != nil {
			handler.SendFail(ctx, err.Error())
		}
		return ctx.Redirect(ctx.Referer())
	case `signedPutObjectURL`:
		return echo.ErrNotFound
		data := ctx.Data()
		name := ctx.Form(`name`)
		if len(name) == 0 {
			return ctx.JSON(data.SetInfo(ctx.T(`参数name的值无效`), 0).SetZone(`name`))
		}
		objectName := path.Join(ppath, name)
		urlData, err := mgr.PresignedPutObject(ctx, objectName, time.Hour*24)
		if err != nil {
			data.SetError(err)
		} else {
			data.SetURL(urlData.String())
		}
		return ctx.JSON(data)
	case `upload`:
		var cu *uploadClient.ChunkUpload
		var opts []uploadClient.ChunkInfoOpter
		if user != nil {
			_cu := uploadChunk.ChunkUploader()
			_cu.UID = fmt.Sprintf(`user/%d`, user.Id)
			cu = &_cu
			opts = append(opts, uploadClient.OptChunkInfoMapping(uploadDropzone.MappingChunkInfo))
		}
		err = mgr.Upload(ctx, ppath, cu, opts...)
		if err != nil {
			user := handler.User(ctx)
			if user != nil {
				notice.OpenMessage(user.Username, `upload`)
				notice.Send(user.Username, notice.NewMessageWithValue(`upload`, ctx.T(`文件上传出错`), err.Error()))
			}
		}
		return respond.Dropzone(ctx, err, nil)
	case `download`:
		return mgr.Download(ctx, ppath)
	default:
		var dirs []os.FileInfo
		var exit bool
		dirs, exit, err = mgr.List(ctx, ppath)
		if exit {
			return err
		}
		ctx.Set(`dirs`, dirs)
	}
	ctx.Set(`parentPath`, parentPath)
	ctx.Set(`path`, ppath)
	var pathPrefix string
	pathPrefix = ppath
	pathSlice := strings.Split(strings.Trim(pathPrefix, `/`), `/`)
	pathLinks := make(echo.KVList, len(pathSlice))
	encodedSep := filemanager.EncodedSep
	urlPrefix := ctx.Request().URL().Path() + fmt.Sprintf(`?id=%d&path=`, id) + encodedSep
	for k, v := range pathSlice {
		urlPrefix += com.URLEncode(v)
		pathLinks[k] = &echo.KV{K: v, V: urlPrefix + `/`}
		urlPrefix += encodedSep
	}
	if !strings.HasSuffix(pathPrefix, `/`) {
		pathPrefix += `/`
	}
	ctx.Set(`pathLinks`, pathLinks)
	ctx.Set(`pathPrefix`, pathPrefix)
	ctx.SetFunc(`Editable`, func(fileName string) bool {
		_, ok := config.DefaultConfig.Sys.Editable(fileName)
		return ok
	})
	ctx.SetFunc(`Playable`, func(fileName string) string {
		mime, _ := config.DefaultConfig.Sys.Playable(fileName)
		return mime
	})
	ctx.Set(`data`, m.NgingCloudStorage)
	return ctx.Render(`cloud/storage_file`, handler.Err(ctx, err))
}
