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

	imageproxy "git.webx.top/coscms/app/base/lib/image/proxy"
	"github.com/admpub/log"
	"github.com/admpub/nging/application/library/collector/exec"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/registry/upload"
	"github.com/admpub/nging/application/registry/upload/driver/filesystem"
	"github.com/admpub/nging/application/registry/upload/helper"
	"github.com/admpub/qrcode"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"

	uploadClient "github.com/webx-top/client/upload"
	_ "github.com/webx-top/client/upload/driver"
)

// ResponseDataForUpload 根据不同的上传方式响应不同的数据格式
func ResponseDataForUpload(ctx echo.Context, field string, err error, imageURLs []string) (result echo.H, embed bool) {
	return upload.ResponserGet(field)(ctx, field, err, imageURLs)
}

var StorerEngine = filesystem.Name

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
	var err error
	typ := ctx.Param(`type`)
	field := ctx.Query(`field`)
	pipe := ctx.Form(`pipe`)
	var files []string
	if len(typ) == 0 {
		err = ctx.E(`请提供参数“%s”`, ctx.Path())
		datax, embed := ResponseDataForUpload(ctx, field, err, files)
		if !embed {
			return ctx.JSON(datax)
		}
		return err
	}
	if !upload.SubdirIsAllowed(typ) {
		err = ctx.E(`参数“%s”未被登记`, typ)
		datax, embed := ResponseDataForUpload(ctx, field, err, files)
		if !embed {
			return ctx.JSON(datax)
		}
		return err
	}
	//echo.Dump(ctx.Forms())
	newStore := upload.StorerGet(StorerEngine)
	if newStore == nil {
		err := errors.New(ctx.T(`存储引擎“%s”未被登记`, StorerEngine))
		datax, embed := ResponseDataForUpload(ctx, field, err, files)
		if !embed {
			return ctx.JSON(datax)
		}
		return err
	}
	storer := newStore(typ)
	var subdir, name string
	subdir, name, err = upload.CheckerGet(typ)(ctx)
	if err != nil {
		return err
	}
	clientName := ctx.Form(`client`)
	if len(clientName) > 0 {
		result := &uploadClient.Result{}
		result.SetDistFileGenerator(func(filename string) (string, error) {
			return SaveFilename(subdir, name, filename)
		})
		return uploadClient.Upload(ctx, clientName, result, storer)
	}
	files, err = upload.BatchUpload(ctx, `files[]`, func(hd *multipart.FileHeader) (string, error) {
		return SaveFilename(subdir, name, hd.Filename)
	}, storer)
	datax, embed := ResponseDataForUpload(ctx, field, err, files)
	if err != nil {
		if !embed {
			return ctx.JSON(datax)
		}
		return err
	}

	if pipe == `deqr` { //解析二维码
		if len(files) > 0 {
			reader, err := storer.Get(files[0])
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
			raw, err := qrcode.Decode(reader, strings.TrimPrefix(path.Ext(files[0]), `.`))
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
	newStore := upload.StorerGet(StorerEngine)
	if newStore == nil {
		return ctx.E(`存储引擎“%s”未被登记`, StorerEngine)
	}
	typ := ctx.Param(`type`)
	storer := newStore(typ)
	_ = storer //TODO: WIP
	src := ctx.Form(`src`)
	src, _ = com.URLDecode(src)
	if err := common.IsRightUploadFile(ctx, src); err != nil {
		return err
	}
	x := ctx.Formx(`x`).Float64()
	y := ctx.Formx(`y`).Float64()
	w := ctx.Formx(`w`).Float64()
	h := ctx.Formx(`h`).Float64()

	//{"x":528,"y":108,"height":864,"width":864,"rotate":0}
	//fmt.Println(avatard)
	option := &imageproxy.CropOptions{
		X:      x, //裁剪X轴起始位置
		Y:      y, //裁剪Y轴起始位置
		Width:  w, //裁剪宽度
		Height: h, //裁剪高度
	}
	opt := imageproxy.Options{
		CropOptions:    option,
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
	absSrcFile := filepath.Join(echo.Wd(), src)
	dstFile := tplfunc.AddSuffix(src, fmt.Sprintf(`_%v_%v`, opt.Width, opt.Height))
	absFile := filepath.Join(echo.Wd(), dstFile)
	cropped := com.FileExists(absFile)
	name := path.Base(src)
	var onSuccess func()

	//对于头像图片，可以根据原图文件的md5值来判断是否需要重新生成缩略图
	if len(name) > 7 && name[0:7] == `avatar.` {
		md5file := filepath.Join(filepath.Dir(absFile), `avatar.md5`)
		onSuccess = func() {
			originMd5 := com.Md5file(absSrcFile)
			err := ioutil.WriteFile(md5file, []byte(originMd5), 0666)
			if err != nil {
				log.Error(err)
			}
		}
		cropped = cropped && com.FileExists(md5file)
		if cropped {
			b, _ := ioutil.ReadFile(md5file)
			originMd5 := com.Md5file(absSrcFile)
			if string(b) == originMd5 {
				goto END
			}
			cropped = false
			onSuccess = func() { //直接使用上面读到的md5
				err := ioutil.WriteFile(md5file, []byte(originMd5), 0666)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}

END:
	if cropped {
		if ctx.Format() == `json` {
			return ctx.JSON(ctx.Data().SetData(dstFile))
		}
		return ctx.File(absFile)
	}

	b, err := ioutil.ReadFile(absSrcFile)
	if err != nil {
		return err
	}
	thumb, err := Resize(bytes.NewReader(b), opt)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(absFile, thumb, 0666) //r-4;w-2;x-1
	if err != nil {
		return err
	}
	if onSuccess != nil {
		onSuccess()
	}
	if ctx.Format() == `json` {
		return ctx.JSON(ctx.Data().SetData(dstFile))
	}
	return ctx.File(absFile)
}

func Resize(r io.Reader, opt imageproxy.Options) (b []byte, err error) {
	return imageproxy.TransformFromReader(r, opt)
}
