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

package file

import (
	"bytes"
	"io"

	imageproxy "github.com/admpub/imageproxy"
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/registry/upload/driver"
	"github.com/webx-top/image"
)

func ImageOptions(width, height float64, cropOptionsSetters ...func(*imageproxy.Options)) *imageproxy.Options {
	cropOptions := &imageproxy.Options{
		//CropX:          x,   //裁剪X轴起始位置
		//CropY:          y,   //裁剪Y轴起始位置
		//CropWidth:      width,  //裁剪宽度
		//CropHeight:     height, //裁剪高度
		Width:     width,  //缩略图宽度
		Height:    height, //缩略图高度
		Quality:   75,
		ScaleUp:   true,
		SmartCrop: true,
	}
	for _, set := range cropOptionsSetters {
		set(cropOptions)
	}
	return cropOptions
}

// CropOptions 图片裁剪选项
type CropOptions struct {
	Options          *imageproxy.Options     //裁剪方式设置
	File             *dbschema.NgingFile     //原图信息
	SrcReader        io.Reader               //原图reader
	Storer           driver.Storer           //存储器
	DestFile         string                  //保存文件路径
	FileMD5          string                  //原图MD5
	WatermarkOptions *image.WatermarkOptions //水印图片文件
	thumbData        *bytes.Reader
}

type CropOptionsSetter func(options *CropOptions)

func CropOptOptions(opt *imageproxy.Options, setters ...func(*imageproxy.Options)) CropOptionsSetter {
	return func(options *CropOptions) {
		for _, set := range setters {
			set(opt)
		}
		options.Options = opt
	}
}

func CropOptFileSchema(fileSchema *dbschema.NgingFile) CropOptionsSetter {
	return func(options *CropOptions) {
		options.File = fileSchema
	}
}

func CropOptSrcReader(reader io.Reader) CropOptionsSetter {
	return func(options *CropOptions) {
		options.SrcReader = reader
	}
}

func CropOptStorer(storer driver.Storer) CropOptionsSetter {
	return func(options *CropOptions) {
		options.Storer = storer
	}
}

func CropOptDestFile(destFile string) CropOptionsSetter {
	return func(options *CropOptions) {
		options.DestFile = destFile
	}
}

func CropOptFileMD5(fileMD5 string) CropOptionsSetter {
	return func(options *CropOptions) {
		options.FileMD5 = fileMD5
	}
}

func CropOptWatermarkOptions(watermarkOptions *image.WatermarkOptions) CropOptionsSetter {
	return func(options *CropOptions) {
		options.WatermarkOptions = watermarkOptions
	}
}

func (c *CropOptions) ApplySetter(setters ...CropOptionsSetter) *CropOptions {
	for _, set := range setters {
		set(c)
	}
	return c
}

func (c *CropOptions) ThumbData() *bytes.Reader {
	return c.thumbData
}

func (c *CropOptions) SetThumbData(thumbData *bytes.Reader) *CropOptions {
	c.thumbData = thumbData
	return c
}
