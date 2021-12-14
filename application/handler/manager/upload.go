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

package manager

import (
	"fmt"
	"io"
	"path/filepath"

	uploadClient "github.com/webx-top/client/upload"
	_ "github.com/webx-top/client/upload/driver"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/middleware/tplfunc"
	"github.com/webx-top/echo/param"

	"github.com/admpub/log"
	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/handler/manager/file"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/config"
	modelFile "github.com/admpub/nging/v4/application/model/file"
	"github.com/admpub/nging/v4/application/model/file/storer"

	"github.com/admpub/nging/v4/application/registry/upload"
	uploadChunk "github.com/admpub/nging/v4/application/registry/upload/chunk"
	_ "github.com/admpub/nging/v4/application/registry/upload/client"
	uploadPipe "github.com/admpub/nging/v4/application/registry/upload/pipe"
	uploadPrepare "github.com/admpub/nging/v4/application/registry/upload/prepare"
	"github.com/admpub/nging/v4/application/registry/upload/thumb"
)

var (
	File                = file.File
	GetWatermarkOptions = storer.GetWatermarkOptions
	CropOptions         = modelFile.ImageOptions
)

// 文件上传保存路径规则：
// 子文件夹/表行ID/文件名

func StorerEngine() storer.Info {
	return storer.Get()
}

// SaveFilename SaveFilename(`0/`,``,`img.jpg`)
func SaveFilename(subdir, name, postFilename string) (string, error) {
	ext := filepath.Ext(postFilename)
	fname := name
	if len(fname) == 0 {
		var err error
		fname, err = common.UniqueID()
		if err != nil {
			return ``, err
		}
	}
	fname += ext
	return subdir + fname, nil
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
func UploadByOwner(ctx echo.Context, ownerType string, ownerID uint64) error {
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
	var err error
	result := &uploadClient.Result{}
	client := uploadClient.Get(clientName)
	client.Init(ctx, result)
	cu := uploadChunk.ChunkUploader()
	cu.UID = fmt.Sprintf(`%s/%d`, ownerType, ownerID)
	client.SetChunkUpload(&cu)
	client.SetUploadMaxSize(int64(config.DefaultConfig.GetMaxRequestBodySize()))
	subdir := ctx.Form(`subdir`, `default`)
	if !upload.AllowedSubdir(subdir) {
		err = ctx.E(`参数subdir的值无效: %s`, subdir)
		return client.SetError(err).Response()
	}
	fileType := ctx.Form(`filetype`)
	storerInfo := StorerEngine()
	prepareData, err := uploadPrepare.Prepare(ctx, subdir, fileType, storerInfo)
	if err != nil {
		return client.SetError(err).Response()
	}
	storer, err := prepareData.Storer(ctx)
	if err != nil {
		return client.SetError(err).Response()
	}
	defer prepareData.Close()
	fileM := modelFile.NewFile(ctx)
	fileM.StorerName = storerInfo.Name
	fileM.StorerId = storerInfo.ID
	fileM.OwnerId = ownerID
	fileM.OwnerType = ownerType
	fileM.Type = fileType
	fileM.Subdir = subdir

	subdir, name, err := prepareData.Checkin(ctx)
	if err != nil {
		return client.SetError(err).Response()
	}
	result.SetFileNameGenerator(func(filename string) (string, error) {
		return SaveFilename(subdir, name, filename)
	})

	callback := func(result *uploadClient.Result, originalReader io.Reader, _ io.Reader) error {
		fileM.Id = 0
		fileM.SetByUploadResult(result)
		if err := ctx.Begin(); err != nil {
			return err
		}
		fileM.Use(common.Tx(ctx))
		err := prepareData.DBSaver(fileM, result, originalReader)
		if err != nil {
			ctx.Rollback()
			return err
		}
		if result.FileType.String() != `image` {
			ctx.Commit()
			return nil
		}
		thumbSizes := thumb.Registry.Get(subdir).AutoCrop()
		if len(thumbSizes) > 0 {
			thumbM := modelFile.NewThumb(ctx)
			thumbM.CPAFrom(fileM.NgingFile)
			for _, thumbSize := range thumbSizes {
				thumbM.Reset()
				if seek, ok := originalReader.(io.Seeker); ok {
					seek.Seek(0, 0)
				}
				thumbURL := tplfunc.AddSuffix(result.FileURL, fmt.Sprintf(`_%v_%v`, thumbSize.Width, thumbSize.Height))
				cropOpt := &modelFile.CropOptions{
					Options:          CropOptions(thumbSize.Width, thumbSize.Height),
					File:             fileM.NgingFile,
					SrcReader:        originalReader,
					Storer:           storer,
					DestFile:         storer.URLToFile(thumbURL),
					FileMD5:          ``,
					WatermarkOptions: GetWatermarkOptions(),
				}
				err = thumbM.Crop(cropOpt)
				if err != nil {
					ctx.Rollback()
					return err
				}
			}
		}
		ctx.Commit()
		return nil
	}

	optionsSetters := []uploadClient.OptionsSetter{
		uploadClient.OptClientName(clientName),
		uploadClient.OptResult(result),
		uploadClient.OptStorer(storer),
		uploadClient.OptWatermarkOptions(GetWatermarkOptions()),
		uploadClient.OptChecker(prepareData.Checker),
		uploadClient.OptCallback(callback),
	}
	if clientName == `default` {
		client.BatchUpload(optionsSetters...)
	} else {
		client.Upload(optionsSetters...)
	}
	if client.GetError() != nil {
		log.Error(client.GetError())
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
		if results == nil {
			results = uploadClient.Results{result}
		}
		err = pipeFunc(ctx, storer, results, recv)
		if err != nil {
			return client.SetError(err).Response()
		}
	}
	return client.Response()
}
