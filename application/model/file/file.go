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
	"path/filepath"
	"strings"
	"time"

	"github.com/coscms/go-imgparse/imgparse"
	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/events"
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/model/base"
	"github.com/admpub/nging/application/registry/upload/table"
)

func NewFile(ctx echo.Context) *File {
	m := &File{
		File: &dbschema.File{},
		base: base.New(ctx),
	}
	m.File.SetContext(ctx)
	return m
}

var Enable bool

type File struct {
	*dbschema.File
	base *base.Base
}

func (f *File) NewFile(m *dbschema.File) *File {
	return &File{
		File: m,
		base: f.base,
	}
}

func (f *File) SetTableID(tableID string) table.TableInfoStorer {
	if len(tableID) == 0 {
		tableID = `0`
	}
	f.File.TableId = tableID
	return f
}

func (f *File) SetTableName(table string) table.TableInfoStorer {
	f.File.TableName = table
	return f
}

func (f *File) SetFieldName(field string) table.TableInfoStorer {
	f.File.FieldName = field
	return f
}

func (f *File) SetByUploadResult(result *uploadClient.Result) *File {
	f.Name = result.FileName
	f.SavePath = result.SavePath
	f.SaveName = filepath.Base(f.SavePath)
	f.Ext = filepath.Ext(f.SavePath)
	f.ViewUrl = result.FileURL
	f.Type = result.FileType.String()
	f.Size = uint64(result.FileSize)
	f.Md5 = result.Md5
	return f
}

func (f *File) FillData(reader io.Reader, forceReset bool, schemas ...*dbschema.File) error {
	var m *dbschema.File
	if len(schemas) > 0 {
		m = schemas[0]
	} else {
		m = f.File
	}
	if forceReset || len(m.Mime) == 0 {
		m.Mime = mime.TypeByExtension(m.Ext)
		if len(f.Mime) == 0 {
			f.Mime = echo.MIMEOctetStream
		}
	}
	if m.Type == `image` {
		typ := strings.TrimPrefix(m.Ext, `.`)
		if typ == `jpg` {
			typ = `jpeg`
		}
		width, height, err := imgparse.ParseRes(reader, typ)
		if err != nil {
			return err
		}
		m.Width = uint(width)
		m.Height = uint(height)
		m.Dpi = 0
	}
	return nil
}

func (f *File) Add(reader io.Reader) error {
	if err := f.FillData(reader, false); err != nil {
		return err
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
		thumbM.Use(f.Trans())
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

func (f *File) DeleteByID(id uint64, ownerType string, ownerID uint64) (err error) {
	f.base.Begin()
	defer func() {
		f.base.End(err == nil)
	}()
	f.Use(f.base.Tx())
	err = f.Get(nil, db.Cond{`id`: id})
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		return nil
	}
	if f.UsedTimes > 0 && (ownerType != `user` || ownerID != 1) {
		return f.base.E(`文件正在使用中，不能删除(只有创始人才能强制删除)`)
	}
	err = f.Delete(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	return f.fireDelete()
}

func (f *File) GetBySavePath(storerName string, savePath string) (err error) {
	err = f.Get(nil, db.And(
		db.Cond{`storer_name`: storerName},
		db.Cond{`save_path`: savePath},
	))
	return
}

func (f *File) GetByViewURL(storerName string, viewURL string) (err error) {
	err = f.Get(nil, db.And(
		db.Cond{`storer_name`: storerName},
		db.Cond{`view_url`: viewURL},
	))
	return
}

func (f *File) FnGetByMd5() func(r *uploadClient.Result) error {
	fileD := &dbschema.File{}
	return func(r *uploadClient.Result) error {
		fileD.Reset()
		err := fileD.Get(nil, db.Cond{`md5`: r.Md5})
		if err != nil {
			if err == db.ErrNoMoreRows {
				return nil
			}
			return err
		}
		r.SavePath = fileD.SavePath
		r.FileURL = fileD.ViewUrl
		return table.ErrExistsFile
	}
}

func (f *File) DeleteBySavePath(savePath string) (err error) {
	f.base.Begin()
	defer func() {
		f.base.End(err == nil)
	}()
	f.Use(f.base.Tx())
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
	f.base.Begin()
	f.Use(f.base.Tx())
	err := f.SetFields(nil, echo.H{
		`table_id`:   ownerID,
		`table_name`: ownerType,
		`field_name`: `avatar`,
		`project`:    project,
		`used_times`: 1,
	}, db.Cond{`id`: f.Id})
	defer func() {
		f.base.End(err == nil)
	}()
	if err != nil {
		return err
	}
	err = f.RemoveUnusedAvatar(ownerType, f.Id)
	return err
}

func (f *File) RemoveUnusedAvatar(ownerType string, excludeID uint64) error {
	return f.DeleteBy(db.And(
		db.Cond{`table_id`: 0},
		db.Cond{`table_name`: ownerType},
		db.Cond{`field_name`: `avatar`},
		db.Cond{`id`: db.NotEq(excludeID)},
	))
}

func (f *File) RemoveUnused(ago int64, ownerType string, ownerID uint64) error {
	cond := db.NewCompounds()
	cond.Add(
		db.Cond{`table_id`: 0},
		db.Cond{`used_times`: 0},
	)
	if len(ownerType) > 0 {
		cond.AddKV(`owner_id`, ownerID)
		cond.AddKV(`owner_type`, ownerType)
	}
	cond.AddKV(`created`, db.Lt(time.Now().Unix()-ago))
	return f.DeleteBy(cond.And())
}

func (f *File) CondByOwner(ownerType string, ownerID uint64) db.Compound {
	return db.And(
		db.Cond{`owner_id`: ownerID},
		db.Cond{`owner_type`: ownerType},
	)
}

func (f *File) DeleteBy(cond db.Compound) error {
	size := 500
	cnt, err := f.ListByOffset(nil, nil, 0, size, cond)
	if err != nil {
		return err
	}
	totalRows := cnt()
	var start int64
	for ; start < totalRows; start += int64(size) {
		if start > 0 {
			cnt, err = f.ListByOffset(nil, nil, 0, size, cond)
			if err != nil {
				return err
			}
		}
		rows := f.Objects()
		for _, fm := range rows {
			f.base.Begin()
			f.Use(f.base.Tx())
			err = f.Delete(nil, db.Cond{`id`: fm.Id})
			if err != nil {
				f.base.Rollback()
				return err
			}
			err = f.fireDelete()
			f.base.End(err == nil)
			if err != nil {
				return err
			}
		}
		if len(rows) < size {
			break
		}
	}
	return err
}

func (f *File) RemoveAvatar(ownerType string, ownerID int64) error {
	return f.DeleteBy(db.And(
		db.Cond{`table_id`: ownerID},
		db.Cond{`table_name`: ownerType},
		db.Cond{`field_name`: `avatar`},
	))
}
