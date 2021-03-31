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
	"io/ioutil"
	"path"
	"strings"

	"github.com/webx-top/client/upload/watermark"
	"github.com/webx-top/echo/param"

	"github.com/admpub/checksum"
	"github.com/admpub/errors"
	imageproxy "github.com/admpub/imageproxy"
)

// Crop 裁剪图片
func (t *Thumb) Crop(opt *CropOptions) error {
	b, err := ioutil.ReadAll(opt.SrcReader)
	if err != nil {
		return errors.WithMessage(err, `Thumb.Crop.ReadAll`)
	}
	thumb, err := imageproxy.Transform(b, *opt.Options)
	if err != nil {
		return errors.WithMessage(err, `Thumb.Crop.Transform`)
	}
	if len(opt.FileMD5) == 0 {
		opt.FileMD5, err = checksum.MD5sumReader(bytes.NewReader(b))
		if err != nil {
			return errors.WithMessage(err, `Thumb.Crop.MD5`)
		}
	}
	if opt.WatermarkOptions != nil && opt.WatermarkOptions.IsEnabled() {
		var extension string
		if pos := strings.LastIndex(opt.DestFile, `.`); pos > -1 {
			extension = opt.DestFile[pos:]
		}
		b, err = watermark.Bytes(thumb, extension, opt.WatermarkOptions)
		if err != nil {
			return errors.WithMessage(err, `Thumb.Crop.Bytes`)
		}
	} else {
		b = thumb
	}
	byteReader := bytes.NewReader(b)
	t.SavePath, t.ViewUrl, err = opt.Storer.Put(opt.DestFile, byteReader, byteReader.Size()) //r-4;w-2;x-1
	if err != nil {
		return errors.WithMessage(err, `Thumb.Crop.Put`)
	}
	opt.SetThumbData(byteReader)
	t.Size = uint64(len(b))
	t.Width = param.AsUint(opt.Options.Width)
	t.Height = param.AsUint(opt.Options.Height)
	t.SaveName = path.Base(t.SavePath)
	t.UsedTimes = 0
	t.Md5 = opt.FileMD5 //原图Md5
	return t.SetByFile(opt.File).Save()
}
