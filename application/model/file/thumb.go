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
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/admpub/checksum"
	"github.com/admpub/errors"
	imageproxy "github.com/admpub/imageproxy"
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/model/base"
)

func NewThumb(ctx echo.Context) *Thumb {
	m := &Thumb{
		FileThumb: &dbschema.FileThumb{},
		base:      base.New(ctx),
	}
	m.FileThumb.SetContext(ctx)
	return m
}

type Thumb struct {
	*dbschema.FileThumb
	base *base.Base
}

func (t *Thumb) SetByFile(file *dbschema.File) *Thumb {
	t.FileId = file.Id
	t.Dpi = file.Dpi
	return t
}

func (t *Thumb) Save() (err error) {
	m := &dbschema.FileThumb{}
	err = m.Get(nil, db.And(
		db.Cond{`save_path`: t.SavePath},
	))
	if err != nil {
		if err != db.ErrNoMoreRows {
			return
		}
		_, err = t.FileThumb.Add()
		return
	}
	t.FileThumb = m
	err = t.SetFields(nil, echo.H{
		`view_url`: t.ViewUrl,
		`size`:     t.Size,
		`width`:    t.Width,
		`height`:   t.Height,
		`dpi`:      t.Dpi,
	}, db.Cond{`id`: m.Id})
	return
}

// Crop 裁剪图片
func (t *Thumb) Crop(opt *CropOptions) error {
	b, err := ioutil.ReadAll(opt.SrcReader)
	if err != nil {
		return err
	}
	thumb, err := imageproxy.Transform(b, *opt.Options)
	if err != nil {
		return err
	}
	if len(opt.WatermarkFile) > 0 {
		var extension string
		if pos := strings.LastIndex(opt.DestFile, `.`); pos > -1 {
			extension = opt.DestFile[pos:]
		}
		b, err = watermark.Bytes(thumb, extension, opt.WatermarkFile)
		if err != nil {
			return err
		}
	} else {
		b = thumb
	}
	byteReader := bytes.NewReader(b)
	t.SavePath, t.ViewUrl, err = opt.Storer.Put(opt.DestFile, byteReader, byteReader.Size()) //r-4;w-2;x-1
	if err != nil {
		return errors.WithMessage(err, `Put`)
	}
	t.Size = uint64(len(b))
	t.Width = param.AsUint(opt.Options.Width)
	t.Height = param.AsUint(opt.Options.Height)
	t.SaveName = path.Base(t.SavePath)
	t.UsedTimes = 0
	if len(opt.FileMD5) == 0 {
		opt.FileMD5, err = checksum.MD5sumReader(opt.SrcReader)
		if err != nil {
			return err
		}
	}
	t.Md5 = opt.FileMD5 //原图Md5
	return t.SetByFile(opt.File).Save()
}
