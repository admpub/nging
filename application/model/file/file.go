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
   along with f program.  If not, see <https://www.gnu.org/licenses/>.
*/

package file

import (
	"io"
	"mime"
	"strings"

	"github.com/admpub/events"
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/model/base"
	"github.com/coscms/go-imgparse/imgparse"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func NewFile(ctx echo.Context) *File {
	return &File{
		File: &dbschema.File{},
		base: base.New(ctx),
	}
}

type File struct {
	*dbschema.File
	base *base.Base
}

func (f *File) Add(reader io.Reader) error {
	if len(f.Mime) == 0 {
		f.Mime = mime.TypeByExtension(f.Ext)
		if len(f.Mime) == 0 {
			f.Mime = echo.MIMEOctetStream
		}
	}
	f.Md5 = com.Md5(f.ViewUrl)
	if f.Type == `image` {
		typ := strings.TrimPrefix(f.Ext, `.`)
		if typ == `jpg` {
			typ = `jpeg`
		}
		width, height, err := imgparse.ParseRes(reader, typ)
		if err != nil {
			return err
		}
		f.Width = uint(width)
		f.Height = uint(height)
		f.Dpi = 0
	}
	_, err := f.File.Add()
	return err
}

func (f *File) fireDelete() error {
	files := []string{f.SavePath}
	thumbM := NewThumb(f.base.Context)
	cnt, err := thumbM.ListByOffset(nil, nil, 0, -1, db.Cond{`file_id`: f.Id})
	if err != nil {
		return err
	}
	thumbNum := cnt()
	if thumbNum > 0 {
		err = thumbM.Delete(nil, db.Cond{`file_id`: f.Id})
		if err != nil {
			return err
		}
		for _, thumb := range thumbM.Objects() {
			files = append(files, thumb.SavePath)
		}
	}
	err = f.base.Fire(f.OwnerType+`-file-deleted`, events.ModeSync, map[string]interface{}{
		`ctx`:     f.base.Context,
		`data`:    f.File,
		`ownerID`: f.OwnerId,
	})
	if err != nil {
		return err
	}
	err = f.base.Fire(`file-deleted`, events.ModeSync, map[string]interface{}{
		`ctx`:   f.base.Context,
		`data`:  f.File,
		`files`: files,
	})
	return err
}

func (f *File) DeleteByID(id uint64) (err error) {
	err = f.Get(nil, db.Cond{`id`: id})
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		return nil
	}
	err = f.Delete(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	return f.fireDelete()
}

func (f *File) DeleteBySavePath(savePath string) (err error) {
	err = f.Get(nil, db.Cond{`save_path`: savePath})
	if err != nil {
		if err != db.ErrNoMoreRows {
			return
		}
		return nil
	}
	err = f.Delete(nil, db.Cond{`id`: f.Id})
	if err != nil {
		return
	}
	return f.fireDelete()
}

func (f *File) UpdateAvatar(project string, ownerType string, ownerID uint64) error {
	err := f.SetFields(nil, echo.H{
		`project`:    project,
		`table_id`:   ownerID,
		`table_name`: ownerType,
		`field_name`: `avatar`,
		`used_times`: 1,
	}, db.Cond{`id`: f.Id})
	if err != nil {
		return err
	}
	err = f.RemoveUnusedAvatar(ownerType, f.Id)
	return err
}

func (f *File) RemoveUnusedAvatar(ownerType string, excludeID uint64) error {
	return f.DeleteBy(db.And(
		db.Cond{`table_name`: ownerType},
		db.Cond{`field_name`: `avatar`},
		db.Cond{`table_id`: 0},
		db.Cond{`id`: db.NotEq(excludeID)},
	))
}

func (f *File) DeleteBy(cond db.Compound) error {
	_, err := f.ListByOffset(nil, nil, 0, -1, cond)
	if err != nil {
		return err
	}
	for _, fm := range f.Objects() {
		err = f.Delete(nil, db.Cond{`id`: fm.Id})
		if err != nil {
			return err
		}
		err = f.fireDelete()
		if err != nil {
			return err
		}
	}
	return err
}

func (f *File) RemoveAvatar(ownerType string, ownerID int64) error {
	return f.DeleteBy(db.And(
		db.Cond{`table_name`: ownerType},
		db.Cond{`field_name`: `avatar`},
		db.Cond{`table_id`: ownerID},
	))
}
