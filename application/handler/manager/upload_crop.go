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
	"os"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"
	"github.com/webx-top/echo/param"

	"github.com/admpub/checksum"
	imageproxy "github.com/admpub/imageproxy"
	"github.com/admpub/log"
	"github.com/admpub/errors"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/middleware"
	modelFile "github.com/admpub/nging/application/model/file"
	"github.com/admpub/nging/application/registry/upload"
	"github.com/admpub/nging/application/registry/upload/convert"
)

// cropPermCheckers 裁剪权限检查
var cropPermCheckers = []func(ctx echo.Context, f *modelFile.File) error{
	func(ctx echo.Context, f *modelFile.File) error {
		if f.FieldName() == `avatar` && f.OwnerType == `user` {
			err := middleware.CheckAnyPerm(ctx, `manager/user_add`, `manager/user_edit`)
			if err != nil {
				return err
			}
			return nil
		}
		return common.ErrNext
	},
}

var watermarkFile = filepath.Join(echo.Wd(), `public/assets/backend/images/nging-gear.png`)

func SetWatermark(markfile string) {
	watermarkFile = markfile
}

func WatermarkFile() string {
	return watermarkFile
}

// CropPermCheckerAdd 添加裁剪权限检查逻辑
func CropPermCheckerAdd(checker func(ctx echo.Context, f *modelFile.File) error) {
	cropPermCheckers = append(cropPermCheckers, checker)
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
	return CropByOwner(ctx, ownerType, ownerID)
}

// CropByOwner 图片裁剪
func CropByOwner(ctx echo.Context, ownerType string, ownerID uint64) error {
	var err error
	uploadType := ctx.Param(`type`)
	if len(uploadType) == 0 {
		uploadType = ctx.Form(`type`)
	}
	storerInfo := StorerEngine()
	prepareData, err := upload.Prepare(ctx, uploadType, ``, storerInfo)
	if err != nil {
		return err
	}
	defer prepareData.Close()
	subdirInfo := upload.SubdirGet(prepareData.TableName)
	if subdirInfo == nil {
		return ctx.E(`“%s”未被登记`, uploadType)
	}

	// 获取缩略图尺寸
	thumbSizes := subdirInfo.ThumbSize(prepareData.FieldName)
	cropSize := ctx.Form(`size`, `200x200`)
	var thumbWidth, thumbHeight float64
	cropSizeArr := strings.SplitN(cropSize, `x`, 2)
	switch len(cropSizeArr) {
	case 2:
		cropSizeArr[1] = strings.TrimSpace(cropSizeArr[1])
		thumbHeight = param.AsFloat64(cropSizeArr[1])
		cropSizeArr[0] = strings.TrimSpace(cropSizeArr[0])
		thumbWidth = param.AsFloat64(cropSizeArr[0])
	case 1:
		cropSizeArr[0] = strings.TrimSpace(cropSizeArr[0])
		thumbWidth = param.AsFloat64(cropSizeArr[0])
		thumbHeight = thumbWidth
	}
	var thumbSize *upload.ThumbSize
	if len(thumbSizes) > 0 {
		for _, ts := range thumbSizes {
			if ts.Width == thumbWidth && ts.Height == thumbHeight {
				thumbSize = &ts
				break
			}
		}
		if thumbSize == nil {
			return ctx.E(`“%s”不支持裁剪图片`, uploadType)
		}
	}

	storer, err := prepareData.Storer(ctx)
	if err != nil {
		return err
	}
	srcURL := ctx.Form(`src`)
	srcURL, err = com.URLDecode(srcURL)
	if err != nil {
		return err
	}
	srcURL = `/` + storer.URLToPath(srcURL)
	if err = common.IsRightUploadFile(ctx, srcURL); err != nil {
		return errors.WithMessage(err, srcURL)
	}
	thumbM := modelFile.NewThumb(ctx)
	fileM := modelFile.NewFile(ctx)
	err = fileM.GetByViewURL(storerInfo, srcURL)
	if err != nil {
		return err
	}
	var editable bool
	if ownerType == `user` && ownerID == 1 { //管理员可编辑
		editable = true
	} else if fileM.OwnerType == ownerType &&
		fileM.OwnerId == ownerID { //上传者可编辑
		editable = true
	} else { //其它验证方式
		for _, check := range cropPermCheckers {
			err := check(ctx, fileM)
			if err == nil { //验证到权限
				editable = true
				break
			}
			if err != common.ErrNext {
				return err
			}
		}
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
	opt := &imageproxy.Options{
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
	if thumbSize != nil {
		opt.Width = thumbSize.Width
		opt.Height = thumbSize.Height
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

	var fileMd5 string
	if onSuccess != nil {
		fileMd5 = onSuccess()
	}
	cropOpt := &modelFile.CropOptions{
		Options:       opt,
		File:          fileM.NgingFile,
		SrcReader:     reader,
		Storer:        storer,
		DestFile:      storer.URLToFile(thumbURL),
		FileMD5:       fileMd5,
		WatermarkFile: watermarkFile,
	}
	err = thumbM.Crop(cropOpt)
	if err != nil {
		return err
	}
	otherFormatExtensions := convert.Extensions()
	for _, extension := range otherFormatExtensions {
		if err := storer.Delete(thumbURL + extension); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	if ctx.Format() == `json` {
		return ctx.JSON(ctx.Data().SetInfo(`cropped`).SetData(thumbURL))
	}
	return storer.SendFile(ctx, thumbURL)
}
