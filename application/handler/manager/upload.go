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
	"github.com/admpub/log"
	uploadClient "github.com/webx-top/client/upload"
	_ "github.com/webx-top/client/upload/driver"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/handler/manager/file"
	modelFile "github.com/admpub/nging/v4/application/model/file"
	"github.com/admpub/nging/v4/application/model/file/storer"

	_ "github.com/admpub/nging/v4/application/registry/upload/client"
	uploadPipe "github.com/admpub/nging/v4/application/registry/upload/pipe"
	uploadPrepare "github.com/admpub/nging/v4/application/registry/upload/prepare"
)

var (
	File                = file.File
	GetWatermarkOptions = storer.GetWatermarkOptions

	// SaveFilename SaveFilename(`0/`,``,`img.jpg`)
	SaveFilename = storer.SaveFilename

	CropOptions = modelFile.ImageOptions
)

// 文件上传保存路径规则：
// 子文件夹/表行ID/文件名

func StorerEngine() storer.Info {
	return storer.Get()
}

// Upload 上传文件
func Upload(ctx echo.Context) error {
	ownerType := `user`
	user := handler.User(ctx)
	var ownerID uint64
	if user != nil {
		ownerID = uint64(user.Id)
	}
	if ownerID < 1 {
		ctx.Data().SetError(ctx.E(`请先登录`))
		return ctx.Redirect(handler.URLFor(`/login`))
	}
	return UploadByOwner(ctx, ownerType, ownerID)
}

// UploadByOwner 上传文件
func UploadByOwner(ctx echo.Context, ownerType string, ownerID uint64, readBeforeHooks ...uploadClient.ReadBeforeHook) error {
	pipe := ctx.Form(`pipe`)
	if len(pipe) > 0 && pipe[0] == '_' {
		pipeFunc := uploadPipe.Get(pipe)
		if pipeFunc == nil {
			return ctx.NewError(code.InvalidParameter, ctx.T(`无效的pipe值`))
		}
		data := echo.H{}
		err := pipeFunc(ctx, nil, nil, data)
		if err != nil {
			return err
		}
		return ctx.JSON(ctx.Data().SetData(data))
	}
	clientName := ctx.Form(`client`, `default`)
	fileType := ctx.Form(`filetype`)
	var err error
	client := uploadPrepare.NewClient(ctx, ownerType, ownerID, clientName, fileType)
	client.SetUploadMaxSize(-1)
	client.AddReadBeforeHook(readBeforeHooks...)
	subdir := ctx.Form(`subdir`, `default`)
	prepareData, err := uploadPrepare.Prepare(ctx, subdir, fileType)
	if err != nil {
		return client.SetError(err).Response()
	}
	defer prepareData.Close()
	fileM := prepareData.MakeModel(ctx, ownerType, ownerID)
	_, err = prepareData.SetMultiple(clientName == `default`).Save(fileM, clientName, client)
	if err != nil {
		log.Error(err.Error())
		return client.Response()
	}
	if len(pipe) > 0 {
		var recv map[string]interface{}
		switch rd := client.GetRespData().(type) {
		case param.Store:
			recv = rd
		case map[string]interface{}:
			recv = rd
		case echo.Data:
			switch dd := rd.GetData().(type) {
			case param.Store:
				recv = dd
			case map[string]interface{}:
				recv = dd
			}
		}
		if recv == nil {
			return client.Response()
		}
		pipeFunc := uploadPipe.Get(pipe)
		if pipeFunc == nil {
			return client.SetError(ctx.NewError(code.InvalidParameter, ctx.T(`无效的pipe值`))).Response()
		}
		results := client.GetBatchUploadResults()
		if results == nil || len(results) == 0 {
			results = uploadClient.Results{client.GetUploadResult()}
		}
		storer, _ := prepareData.Storer(ctx)
		err = pipeFunc(ctx, storer, results, recv)
		if err != nil {
			return client.SetError(err).Response()
		}
	}
	return client.Response()
}
