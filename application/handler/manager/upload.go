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
	"bytes"
	"io"
	"mime/multipart"
	"path"
	"path/filepath"
	"strings"
	"time"

	uploadClient "github.com/webx-top/client/upload"
	_ "github.com/webx-top/client/upload/driver"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/common"
	modelFile "github.com/admpub/nging/application/model/file"
	"github.com/admpub/nging/application/registry/upload"
	"github.com/admpub/nging/application/registry/upload/convert"
	"github.com/admpub/nging/application/registry/upload/driver/local"
	"github.com/admpub/nging/application/registry/upload/helper"
	"github.com/admpub/qrcode"
	"github.com/admpub/nging/application/model/file/storer"
)

// 文件上传保存路径规则：
// 表名称/表行ID/文件名

// ResponseDataForUpload 根据不同的上传方式响应不同的数据格式
func ResponseDataForUpload(ctx echo.Context, field string, err error, imageURLs []string) (result echo.H, embed bool) {
	return upload.ResponserGet(field)(ctx, field, err, imageURLs)
}

func StorerEngine() storer.Info {
	return storer.Get()
}

func File(ctx echo.Context) error {
	uploadType := ctx.Param(`type`)
	typ, _, _ := getTableInfo(uploadType)
	file := ctx.Param(`*`)
	file = filepath.Join(helper.UploadDir, typ, file)
	var (
		convertFunc convert.Convert
		ok bool
		originalFile string
	)
	extension := ctx.Query(`ex`)
	if len(extension) > 0 {
		extension = `.` + extension
		convertFunc, ok = convert.GetConverter(extension)
		if !ok {
			return ctx.File(file)
		}
		originalFile = file
	} else {
		originalExtension := filepath.Ext(file)
		extension = strings.ToLower(originalExtension)
		convertFunc, ok = convert.GetConverter(extension)
		if !ok {
			return ctx.File(file)
		}
		originalFile = strings.TrimSuffix(file, originalExtension)
		index := strings.LastIndex(originalFile, `.`)
		// 单扩展名或相同扩展名的情况下不转换格式
		if index < 0 || strings.ToLower(originalFile[index:]) == extension {
			return ctx.File(originalFile)
		}
	}
	supported := strings.Contains(ctx.Header(echo.HeaderAccept), "image/" + strings.TrimPrefix(extension, `.`))
	if !supported {
		return ctx.File(originalFile)
	}
	if err := ctx.File(file); err != echo.ErrNotFound {
		return err
	}
	storerName := local.Name
	return ctx.ServeCallbackContent(func(_ echo.Context) (io.Reader, error) {
		newStore := upload.StorerGet(storerName)
		if newStore == nil {
			return nil, ctx.E(`存储引擎“%s”未被登记`, storerName)
		}
		storer, err := newStore(ctx, typ)
		if err != nil {
			return nil, err
		}
		f, err := storer.Get(`/` + originalFile)
		if err != nil {
			return nil, echo.ErrNotFound
		}
		defer f.Close()
		buf, err := convertFunc(f, 70)
		if err != nil {
			return nil, err
		}
		b := buf.Bytes()
		saveFile := storer.URLToFile(`/` + file)
		_, _, err = storer.Put(saveFile, buf, int64(len(b)))
		return bytes.NewBuffer(b), err
	}, path.Base(file), time.Now())
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
	uploadType := ctx.Param(`type`)
	field := ctx.Query(`field`) // 上传表单file输入框名称
	pipe := ctx.Form(`pipe`)
	var (
		err      error
		fileURLs []string
	)
	if len(uploadType) == 0 {
		err = ctx.E(`请提供参数“%s”`, ctx.Path())
		datax, embed := ResponseDataForUpload(ctx, field, err, fileURLs)
		if !embed {
			return ctx.JSON(datax)
		}
		return err
	}
	fileType := ctx.Form(`filetype`)
	storerInfo := StorerEngine()
	prepareData, err := upload.Prepare(ctx, uploadType, fileType, storerInfo)
	if err != nil {
		datax, embed := ResponseDataForUpload(ctx, field, err, fileURLs)
		if !embed {
			return ctx.JSON(datax)
		}
	}
	storer, err := prepareData.Storer(ctx)
	if err != nil {
		datax, embed := ResponseDataForUpload(ctx, field, err, fileURLs)
		if !embed {
			return ctx.JSON(datax)
		}
	}
	defer prepareData.Close()
	fileM := modelFile.NewFile(ctx)
	fileM.StorerName = storerInfo.Name
	fileM.StorerId = storerInfo.ID
	fileM.TableId = ``
	fileM.SetFieldName(prepareData.FieldName)
	fileM.SetTableName(prepareData.TableName)
	fileM.OwnerId = ownerID
	fileM.OwnerType = ownerType
	fileM.Type = fileType

	subdir, name, err := prepareData.Checkin(ctx, fileM)
	if err != nil {
		datax, embed := ResponseDataForUpload(ctx, field, err, fileURLs)
		if !embed {
			return ctx.JSON(datax)
		}
		return err
	}

	clientName := ctx.Form(`client`)
	if len(clientName) > 0 {
		result := &uploadClient.Result{}
		result.SetFileNameGenerator(func(filename string) (string, error) {
			return SaveFilename(subdir, name, filename)
		})

		client := uploadClient.Upload(ctx, clientName, result, storer, watermarkFile, prepareData.Checker)
		if client.GetError() != nil {
			if client.GetError() == upload.ErrExistsFile {
				client.SetError(nil)
			}
			return client.Response()
		}

		fileM.SetByUploadResult(result)

		var reader io.ReadCloser
		reader, err = storer.Get(result.SavePath)
		if reader != nil {
			defer reader.Close()
		}
		if err != nil {
			return client.SetError(err).Response()
		}
		err = prepareData.DBSaver(fileM, result, reader)
		return client.SetError(err).Response()
	}
	var results uploadClient.Results
	results, err = upload.BatchUpload(
		ctx,
		`files[]`,
		func(r *uploadClient.Result) (string, error) {
			if err := prepareData.Checker(r); err != nil {
				return ``, err
			}
			return SaveFilename(subdir, name, r.FileName)
		},
		storer,
		func(result *uploadClient.Result, file multipart.File) error {
			fileM.Id = 0
			fileM.SetByUploadResult(result)
			return prepareData.DBSaver(fileM, result, file)
		},
		watermarkFile,
	)
	datax, embed := ResponseDataForUpload(ctx, field, err, results.FileURLs())
	if err != nil {
		if !embed {
			return ctx.JSON(datax)
		}
		return err
	}

	if pipe == `deqr` { //解析二维码
		if len(results) > 0 {
			reader, err := storer.Get(results[0].SavePath)
			if reader != nil {
				defer reader.Close()
			}
			if err != nil {
				if !embed {
					datax[`raw`] = err.Error()
					return ctx.JSON(datax)
				}
				return err
			}
			raw, err := qrcode.Decode(reader, strings.TrimPrefix(path.Ext(results[0].SavePath), `.`))
			if err != nil {
				raw = err.Error()
			}
			datax[`raw`] = raw
		}
	}
	if !embed {
		return ctx.JSON(datax)
	}
	data := ctx.Data()
	data.SetData(datax)
	return ctx.JSON(data)
}

func getTableInfo(uploadType string) (tableName string, fieldName string, defaults []string) {
	return upload.GetTableInfo(uploadType)
}
