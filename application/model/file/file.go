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
	"github.com/admpub/nging/v3/application/dbschema"
	"github.com/admpub/nging/v3/application/model/file/storer"
)

func NewFile(ctx echo.Context) *File {
	m := &File{
		NgingFile: dbschema.NewNgingFile(ctx),
	}
	return m
}

type File struct {
	*dbschema.NgingFile
}

func (f *File) NewFile(m *dbschema.NgingFile) *File {
	r := &File{
		NgingFile: m,
	}
	r.SetContext(f.Context())
	return r
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

func (f *File) FillData(reader io.Reader, forceReset bool, schemas ...*dbschema.NgingFile) error {
	var m *dbschema.NgingFile
	if len(schemas) > 0 {
		m = schemas[0]
	} else {
		m = f.NgingFile
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
	_, err := f.NgingFile.Add()
	return err
}

func (f *File) fireDelete() error {
	files := []string{f.SavePath}
	thumbM := NewThumb(f.Context())
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
	err = f.Context().Fire(f.OwnerType+`-file-deleted`, events.ModeSync, map[string]interface{}{
		`ctx`:     f.Context(),
		`data`:    f.NgingFile,
		`ownerID`: f.OwnerId,
	})
	if err != nil {
		return err
	}
	err = f.Context().Fire(`file-deleted`, events.ModeSync, map[string]interface{}{
		`ctx`:   f.Context(),
		`data`:  f.NgingFile,
		`files`: files,
	})
	return err
}

func (f *File) DeleteByID(id uint64, ownerType string, ownerID uint64) (err error) {
	f.Context().Begin()
	defer func() {
		f.Context().End(err == nil)
	}()
	err = f.Get(nil, db.Cond{`id`: id})
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		return nil
	}
	if f.UsedTimes > 0 && (ownerType != `user` || ownerID != 1) {
		return f.Context().E(`文件正在使用中，不能删除(只有创始人才能强制删除)`)
	}
	err = f.Delete(nil, db.Cond{`id`: id})
	if err != nil {
		return err
	}
	return f.fireDelete()
}

func (f *File) GetBySavePath(storerInfo storer.Info, savePath string) (err error) {
	err = f.Get(nil, db.And(
		db.Cond{`storer_name`: storerInfo.Name},
		db.Cond{`storer_id`: storerInfo.ID},
		db.Cond{`save_path`: savePath},
	))
	return
}

func (f *File) GetByStorerAndURL(storerInfo storer.Info, viewURL string) (err error) {
	err = f.Get(nil, db.And(
		db.Cond{`storer_name`: storerInfo.Name},
		db.Cond{`storer_id`: storerInfo.ID},
		db.Cond{`view_url`: viewURL},
	))
	return
}

func (f *File) GetByViewURL(viewURL string) (err error) {
	err = f.Get(nil, db.And(
		db.Cond{`view_url`: viewURL},
	))
	return
}

func (f *File) FnGetByMd5() func(r *uploadClient.Result) error {
	fileD := dbschema.NewNgingFile(f.Context())
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
		return uploadClient.ErrExistsFile
	}
}

func (f *File) DeleteBySavePath(savePath string) (err error) {
	f.Context().Begin()
	defer func() {
		f.Context().End(err == nil)
	}()
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

func (f *File) RemoveUnusedAvatar(ownerType string, excludeID uint64) error {
	return f.DeleteBy(db.And(
		db.Cond{`subdir`: `avatar`},
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

// CondByOwner 所有者条件
func (f *File) CondByOwner(ownerType string, ownerID uint64) db.Compound {
	return db.And(
		db.Cond{`owner_id`: ownerID},
		db.Cond{`owner_type`: ownerType},
	)
}

func (f *File) DeleteBy(cond db.Compound) error {
	size := 500
	_, err := f.ListByOffset(nil, nil, 0, size, cond)
	if err != nil {
		return err
	}
	totalRows, _ := f.Count(nil, cond)
	var start int64
	for ; start < totalRows; start += int64(size) {
		if start > 0 {
			_, err = f.ListByOffset(nil, nil, 0, size, cond)
			if err != nil {
				return err
			}
		}
		rows := f.Objects()
		for _, fm := range rows {
			f.Context().Begin()
			err = f.Delete(nil, db.Cond{`id`: fm.Id})
			if err != nil {
				f.Context().Rollback()
				return err
			}
			err = f.fireDelete()
			f.Context().End(err == nil)
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

func (f *File) GetAvatar() (*dbschema.NgingFile, error) {
	m := dbschema.NewNgingFile(nil)
	m.CPAFrom(f.NgingFile)
	err := m.Get(nil, db.Cond{`view_url`: f.ViewUrl})
	return m, err
}
