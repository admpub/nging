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

	uploadClient "github.com/webx-top/client/upload"
	_ "github.com/webx-top/client/upload/driver"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"
	"github.com/webx-top/echo/param"

	"github.com/admpub/checksum"
	imageproxy "github.com/admpub/imageproxy"
	"github.com/admpub/log"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/collector/exec"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/middleware"
	modelFile "github.com/admpub/nging/application/model/file"
	"github.com/admpub/nging/application/registry/upload"
	"github.com/admpub/nging/application/registry/upload/driver/filesystem"
	"github.com/admpub/nging/application/registry/upload/helper"
	"github.com/admpub/qrcode"
)

// 文件上传保存路径规则：
// 表名称/表行ID/文件名

// ResponseDataForUpload 根据不同的上传方式响应不同的数据格式
func ResponseDataForUpload(ctx echo.Context, field string, err error, imageURLs []string) (result echo.H, embed bool) {
	return upload.ResponserGet(field)(ctx, field, err, imageURLs)
}

var (
	StorerEngine   = filesystem.Name
	DefaultChecker = func(r *uploadClient.Result) error {
		return nil
	}
)

func File(ctx echo.Context) error {
	uploadType := ctx.Param(`type`)
	typ, _, _ := getTableInfo(uploadType)
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
	tableName, fieldName, defaults := getTableInfo(uploadType)
	if !upload.SubdirIsAllowed(uploadType, defaults...) {
		err = ctx.E(`参数“%s”未被登记`, uploadType)
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
	fileM.TableId = ``
	fileM.SetFieldName(fieldName)
	fileM.SetTableName(tableName)
	fileM.OwnerId = ownerID
	fileM.OwnerType = ownerType
	fileType := ctx.Form(`filetype`)
	fileM.Type = fileType

	storer := newStore(ctx, tableName) // 使用表名称作为文件夹名
	defer storer.Close()
	var subdir, name string
	subdir, name, err = upload.CheckerGet(uploadType, defaults...)(ctx, fileM)
	if err != nil {
		return err
	}
	dbsaver := upload.DBSaverGet(uploadType, defaults...)
	checker := func(r *uploadClient.Result) error {
		extension := path.Ext(r.FileName)
		if len(r.FileType) > 0 {
			if !uploadClient.CheckTypeExtension(fileType, extension) {
				return ctx.E(`不支持将扩展名为“%v”的文件作为“%v”类型的文件来进行上传`, extension, fileType)
			}
		} else {
			r.FileType = uploadClient.FileType(uploadClient.DetectType(extension))
		}
		return DefaultChecker(r) //fileM.FnGetByMd5()
	}

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
		err = dbsaver(fileM, result, reader)
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
			return dbsaver(fileM, result, file)
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

func getTableInfo(uploadType string) (tableName string, fieldName string, defaults []string) {
	return upload.GetTableInfo(uploadType)
}

// Crop 图片裁剪
func Crop(ctx echo.Context) error {
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
	return CropByOwner(ctx, ownerType, ownerID, func(f *modelFile.File) error {
		if f.FieldName() == `avatar` && f.OwnerType == `user` {
			err := middleware.CheckAnyPerm(ctx, `manager/user_add`, `manager/user_edit`)
			if err != nil {
				return err
			}
			return nil
		}
		return common.ErrUserNoPerm
	})
}

// CropByOwner 图片裁剪
func CropByOwner(ctx echo.Context, ownerType string, ownerID uint64, permChecker func(*modelFile.File) error) error {
	var err error
	newStore := upload.StorerGet(StorerEngine)
	if newStore == nil {
		return ctx.E(`存储引擎“%s”未被登记`, StorerEngine)
	}
	uploadType := ctx.Param(`type`)
	typ, _, _ := getTableInfo(uploadType)
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
	var editable bool
	if ownerType == `user` && ownerID == 1 { //管理员可编辑
		editable = true
	} else if fileM.OwnerType == ownerType &&
		fileM.OwnerId == ownerID { //上传者可编辑
		editable = true
	} else if err = permChecker(fileM); err != nil { //其它验证方式
		return err
	} else {
		editable = true
	}
	if !editable {
		return common.ErrUserNoPerm
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
