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
	"bytes"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/middleware/tplfunc"
	"github.com/webx-top/echo/param"

	"github.com/admpub/checksum"
	"github.com/admpub/errors"
	imageproxy "github.com/admpub/imageproxy"
	"github.com/admpub/log"
	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/common"
	modelFile "github.com/admpub/nging/v4/application/model/file"
	uploadChecker "github.com/admpub/nging/v4/application/registry/upload/checker"
	"github.com/admpub/nging/v4/application/registry/upload/convert"
	"github.com/admpub/nging/v4/application/registry/upload/helper"
	uploadPrepare "github.com/admpub/nging/v4/application/registry/upload/prepare"
	"github.com/admpub/nging/v4/application/registry/upload/thumb"
)

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
	cropSize := ctx.Form(`size`, thumb.DefaultSize.String())
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
	srcURL := ctx.Form(`src`)
	srcURL, err = com.URLDecode(srcURL)
	if err != nil {
		return err
	}
	fileM := modelFile.NewFile(ctx)
	err = fileM.GetByViewURL(srcURL)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.NewError(code.DataNotFound, ctx.T(`文件数据不存在`))
		}
		return err
	}
	if fileM.Type != `image` {
		return ctx.NewError(code.InvalidParameter, ctx.T(`只支持裁剪图片文件`))
	}
	subdir := fileM.Subdir
	if len(subdir) == 0 {
		subdir = helper.ParseSubdir(srcURL)
	}
	var unlimitResize bool
	unlimitResizeToken := ctx.Form(`token`)
	if len(unlimitResizeToken) > 0 {
		unlimitResize = unlimitResizeToken == uploadChecker.Token(`file`, srcURL, `width`, thumbWidth, `height`, thumbHeight)
	}
	storerInfo := StorerEngine()
	prepareData, err := uploadPrepare.Prepare(ctx, subdir, ``, storerInfo)
	if err != nil {
		return err
	}
	defer prepareData.Close()
	var thumbSize *thumb.Size
	if !unlimitResize { // 是否检查尺寸
		// 获取缩略图尺寸
		thumbSize = thumb.Registry.Get(subdir).Get(thumbWidth, thumbHeight)
		if thumbSize == nil {
			return ctx.E(`不支持裁剪图片尺寸: %vx%v`, thumbWidth, thumbHeight)
		}
	} else {
		thumbSize = &thumb.Size{
			Width:  thumbWidth,
			Height: thumbHeight,
		}
	}
	ctx.Internal().Set(`storerID`, fileM.StorerId)
	storer, err := prepareData.Storer(ctx)
	if err != nil {
		return err
	}
	srcURL = storer.URLToPath(srcURL)
	if err = common.IsRightUploadFile(ctx, srcURL); err != nil {
		return errors.WithMessage(err, srcURL)
	}
	thumbM := modelFile.NewThumb(ctx)
	var editable bool
	if ownerType == `user` && ownerID == 1 { //管理员可编辑
		editable = true
	} else if fileM.OwnerType == ownerType &&
		fileM.OwnerId == ownerID { //上传者可编辑
		editable = true
	} else { //其它验证方式
		//editable = true //TODO: 验证
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
		CropX:          x,                        //裁剪X轴起始位置
		CropY:          y,                        //裁剪Y轴起始位置
		CropWidth:      w,                        //裁剪宽度
		CropHeight:     h,                        //裁剪高度
		Width:          thumb.DefaultSize.Width,  //缩略图宽度
		Height:         thumb.DefaultSize.Height, //缩略图高度
		Fit:            false,
		Rotate:         0,
		FlipVertical:   false,
		FlipHorizontal: false,
		Quality:        thumb.DefaultSize.Quality,
		Signature:      "",
		ScaleUp:        true,
	}
	if thumbSize != nil {
		opt.Width = thumbSize.Width
		opt.Height = thumbSize.Height
		if thumbSize.Quality > 0 {
			if thumbSize.Quality > 100 {
				opt.Quality = 100
			} else {
				opt.Quality = thumbSize.Quality
			}
		}
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
			b, err := io.ReadAll(md5reader)
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
			return ctx.JSON(ctx.Data().SetInfo(`skipped`).SetData(storer.BaseURL() + thumbURL))
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
		Options:          opt,
		File:             fileM.NgingFile,
		SrcReader:        reader,
		Storer:           storer,
		DestFile:         storer.URLToFile(thumbURL),
		FileMD5:          fileMd5,
		WatermarkOptions: GetWatermarkOptions(),
	}
	//panic(cropOpt.DestFile)
	err = thumbM.Crop(cropOpt)
	if err != nil {
		return err
	}
	otherFormatExtensions := convert.Extensions()
	for _, extension := range otherFormatExtensions {
		if err := storer.Delete(thumbURL + extension); err != nil && !storer.ErrIsNotExist(err) {
			return err
		}
	}
	if ctx.Format() == `json` {
		return ctx.JSON(ctx.Data().SetInfo(`cropped`).SetData(storer.BaseURL() + thumbURL))
	}
	return storer.SendFile(ctx, thumbURL)
}
