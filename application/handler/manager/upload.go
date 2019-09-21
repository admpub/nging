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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"path"
	"path/filepath"
	"strings"

	"github.com/admpub/checksum"
	imageproxy "github.com/admpub/imageproxy"
	"github.com/admpub/log"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/collector/exec"
	"github.com/admpub/nging/application/library/common"
	modelFile "github.com/admpub/nging/application/model/file"
	"github.com/admpub/nging/application/registry/upload"
	"github.com/admpub/nging/application/registry/upload/driver/filesystem"
	"github.com/admpub/nging/application/registry/upload/helper"
	"github.com/admpub/qrcode"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"
	"github.com/webx-top/echo/param"

	uploadClient "github.com/webx-top/client/upload"
	_ "github.com/webx-top/client/upload/driver"
)

// ResponseDataForUpload 根据不同的上传方式响应不同的数据格式
func ResponseDataForUpload(ctx echo.Context, field string, err error, imageURLs []string) (result echo.H, embed bool) {
	return upload.ResponserGet(field)(ctx, field, err, imageURLs)
}

var (
	StorerEngine = filesystem.Name
)

func File(ctx echo.Context) error {
	typ := ctx.Param(`type`)
	file := ctx.Param(`*`)
	file = filepath.Join(helper.UploadDir, typ, file)
	return ctx.File(file)
}

func SaveFilename(subdir, name, postFilename string) (string, error) {
	ext := filepath.Ext(postFilename)
	fname := name
	if len(fname) == 0 {
		var err error
		fname, err = exec.UniqueID()
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
	return UploadByOwner(ctx, ownerType, ownerID)
}

// UploadByOwner 上传文件
func UploadByOwner(ctx echo.Context, ownerType string, ownerID uint64) error {
	typ := ctx.Param(`type`)
	field := ctx.Query(`field`)
	pipe := ctx.Form(`pipe`)
	var (
		err      error
		fileURLs []string
	)
	if len(typ) == 0 {
		err = ctx.E(`请提供参数“%s”`, ctx.Path())
		datax, embed := ResponseDataForUpload(ctx, field, err, fileURLs)
		if !embed {
			return ctx.JSON(datax)
		}
		return err
	}
	if !upload.SubdirIsAllowed(typ) {
		err = ctx.E(`参数“%s”未被登记`, typ)
		datax, embed := ResponseDataForUpload(ctx, field, err, fileURLs)
		if !embed {
			return ctx.JSON(datax)
		}
		return err
	}
	//echo.Dump(ctx.Forms())
	newStore := upload.StorerGet(StorerEngine)
	if newStore == nil {
		err := errors.New(ctx.T(`存储引擎“%s”未被登记`, StorerEngine))
		datax, embed := ResponseDataForUpload(ctx, field, err, fileURLs)
		if !embed {
			return ctx.JSON(datax)
		}
		return err
	}

	fileM := modelFile.NewFile(ctx)
	fileM.StorerName = StorerEngine
	fileM.TableId = 0
	fileM.TableName = typ
	fileM.FieldName = ``
	fileM.Type = ctx.Form(`filetype`, `image`)

	storer := newStore(ctx, typ)
	defer storer.Close()
	var subdir, name string
	subdir, name, err = upload.CheckerGet(typ)(ctx, fileM)
	if err != nil {
		return err
	}
	checker := fileM.FnGetByMd5()

	clientName := ctx.Form(`client`)
	if len(clientName) > 0 {
		result := &uploadClient.Result{}
		result.SetFileNameGenerator(func(filename string) (string, error) {
			return SaveFilename(subdir, name, filename)
		})

		client := uploadClient.Upload(ctx, clientName, result, storer, checker)
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
		err = fileM.Add(reader)
		return client.SetError(err).Response()
	}
	var results uploadClient.Results
	results, err = upload.BatchUpload(
		ctx,
		`files[]`,
		func(r *uploadClient.Result) (string, error) {
			if err := checker(r); err != nil {
				return ``, err
			}
			return SaveFilename(subdir, name, r.FileName)
		},
		storer,
		func(result *uploadClient.Result, file multipart.File) error {
			fileM.Id = 0
			fileM.SetByUploadResult(result)
			/*
				reader, err := storer.Get(result.SavePath)
				if err != nil {
					return err
				}
				err = fileM.Add(reader)
			*/
			err = fileM.Add(file)
			return err
		},
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

func Crop(ctx echo.Context) error {
	var err error
	newStore := upload.StorerGet(StorerEngine)
	if newStore == nil {
		return ctx.E(`存储引擎“%s”未被登记`, StorerEngine)
	}
	typ := ctx.Param(`type`)
	storer := newStore(ctx, typ)
	defer storer.Close()
	srcURL := ctx.Form(`src`)
	srcURL, err = com.URLDecode(srcURL)
	if err != nil {
		return err
	}
	if err = common.IsRightUploadFile(ctx, srcURL); err != nil {
		return err
	}
	thumbM := modelFile.NewThumb(ctx)
	fileM := modelFile.NewFile(ctx)
	err = fileM.GetByViewURL(StorerEngine, srcURL)
	if err != nil {
		return err
	}

	x := ctx.Formx(`x`).Float64()
	y := ctx.Formx(`y`).Float64()
	w := ctx.Formx(`w`).Float64()
	h := ctx.Formx(`h`).Float64()

	//{"x":528,"y":108,"height":864,"width":864,"rotate":0}
	//fmt.Println(avatard)
	opt := imageproxy.Options{
		CropX:          x,   //裁剪X轴起始位置
		CropY:          y,   //裁剪Y轴起始位置
		CropWidth:      w,   //裁剪宽度
		CropHeight:     h,   //裁剪高度
		Width:          200, //缩略图宽度
		Height:         200, //缩略图高度
		Fit:            false,
		Rotate:         0,
		FlipVertical:   false,
		FlipHorizontal: false,
		Quality:        100,
		Signature:      "",
		ScaleUp:        true,
	}
	thumbURL := tplfunc.AddSuffix(srcURL, fmt.Sprintf(`_%v_%v`, opt.Width, opt.Height))
	var cropped bool
	cropped, err = storer.Exists(thumbURL)
	if err != nil {
		return err
	}
	name := path.Base(srcURL)
	var onSuccess func() string

	//对于头像图片，可以根据原图文件的md5值来判断是否需要重新生成缩略图
	if len(name) > 7 && name[0:7] == `avatar.` {
		md5file := path.Join(path.Dir(srcURL), `avatar.md5`)
		putFile := storer.URLToFile(md5file)
		onSuccess = func() string {
			reader, err := storer.Get(srcURL)
			if reader != nil {
				defer reader.Close()
			}
			if err != nil {
				log.Error(err)
				return ``
			}
			originMd5, err := checksum.MD5sumReader(reader)
			if err != nil {
				log.Error(err)
				return ``
			}
			size := len(originMd5)
			_, _, err = storer.Put(putFile, bytes.NewBufferString(originMd5), int64(size))
			if err != nil {
				log.Error(err)
			}
			return originMd5
		}

		if cropped {
			cropped, err = storer.Exists(md5file)
			if err != nil {
				return err
			}
		}

		if cropped {
			md5reader, err := storer.Get(md5file)
			if md5reader != nil {
				defer md5reader.Close()
			}
			if err != nil {
				return err
			}
			b, err := ioutil.ReadAll(md5reader)
			if err != nil {
				return err
			}
			reader, err := storer.Get(srcURL)
			if reader != nil {
				defer reader.Close()
			}
			originMd5, err := checksum.MD5sumReader(reader)
			if err != nil {
				return err
			}
			if string(b) == originMd5 {
				goto END
			}
			cropped = false
			onSuccess = func() string { //直接使用上面读到的md5
				size := len(originMd5)
				_, _, err = storer.Put(putFile, bytes.NewBufferString(originMd5), int64(size))
				if err != nil {
					log.Error(err)
				}
				return originMd5
			}
		}
	}

END:
	if cropped {
		if ctx.Format() == `json` {
			return ctx.JSON(ctx.Data().SetInfo(`skipped`).SetData(thumbURL))
		}
		return storer.SendFile(ctx, thumbURL)
	}

	var reader io.ReadCloser
	reader, err = storer.Get(srcURL)
	if reader != nil {
		defer reader.Close()
	}
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	thumb, err := imageproxy.Transform(b, opt)
	if err != nil {
		return err
	}
	byteReader := bytes.NewReader(thumb)
	thumbM.SavePath, thumbM.ViewUrl, err = storer.Put(storer.URLToFile(thumbURL), byteReader, byteReader.Size()) //r-4;w-2;x-1
	if err != nil {
		return err
	}
	var fileMd5 string
	if onSuccess != nil {
		fileMd5 = onSuccess()
	} else {
		fileMd5, err = checksum.MD5sumReader(reader)
		if err != nil {
			return err
		}
	}
	size := len(thumb)
	thumbM.Size = uint64(size)
	thumbM.Width = param.AsUint(opt.Width)
	thumbM.Height = param.AsUint(opt.Height)
	thumbM.SaveName = path.Base(thumbM.SavePath)
	thumbM.UsedTimes = 0
	thumbM.Md5 = fileMd5
	err = thumbM.SetByFile(fileM.File).Save()
	if err != nil {
		return err
	}

	if ctx.Format() == `json` {
		return ctx.JSON(ctx.Data().SetInfo(`cropped`).SetData(thumbURL))
	}
	return storer.SendFile(ctx, thumbURL)
}
