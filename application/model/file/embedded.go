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
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/fileupdater"
	"github.com/admpub/nging/application/model/base"
	uploadStorer "github.com/admpub/nging/application/registry/upload/driver"
	uploadHelper "github.com/admpub/nging/application/registry/upload/helper"
)

func NewEmbedded(ctx echo.Context, fileMdls ...*File) *Embedded {
	if ctx == nil {
		panic(`ctx is nil`)
	}
	var fileM *File
	if len(fileMdls) > 0 {
		fileM = fileMdls[0]
	} else {
		fileM = NewFile(ctx)
	}
	m := &Embedded{
		FileEmbedded: &dbschema.FileEmbedded{},
		base:         base.New(ctx),
		File:         fileM,
	}
	m.FileEmbedded.SetContext(ctx)
	return m
}

type Embedded struct {
	*dbschema.FileEmbedded
	base    *base.Base
	File    *File
	updater *fileupdater.FileUpdater
}

func (f *Embedded) Updater(table string, field string, tableID string) *fileupdater.FileUpdater {
	if f.updater == nil {
		f.updater = fileupdater.New(f)
	}
	f.updater.Set(table, field, tableID)
	return f.updater
}

func (f *Embedded) FileIDs() []uint64 {
	fileIDs := []uint64{}
	if len(f.FileIds) == 0 {
		return fileIDs
	}
	for _, fileID := range strings.Split(f.FileIds, `,`) {
		fileIDs = append(fileIDs, param.AsUint64(fileID))
	}
	return fileIDs
}

func (f *Embedded) MoveFileToOwner(fileIDs []uint64, ownerID string) (map[string]string, error) {
	replaces := make(map[string]string)
	if len(fileIDs) == 0 {
		return replaces, nil
	}
	_, err := f.File.ListByOffset(nil, nil, 0, -1, db.Cond{`id`: db.In(fileIDs)})
	if err != nil {
		return replaces, err
	}
	replaceFrom := `/0/`
	replaceTo := `/` + ownerID + `/`
	storers := map[string]uploadStorer.Storer{}
	defer func() {
		for _, storer := range storers {
			storer.Close()
		}
	}()
	for _, file := range f.File.Objects() {
		if !strings.Contains(file.SavePath, replaceFrom) {
			continue
		}
		storer, ok := storers[file.StorerName]
		if !ok {
			newStore := uploadStorer.Get(file.StorerName)
			if newStore == nil {
				return replaces, f.base.E(`存储引擎“%s”未被登记`, file.StorerName)
			}
			storer = newStore(f.base.Context, ``)
			storers[file.StorerName] = storer
		}
		var newSavePath, newViewURL, prefix string
		if file.FieldName == `avatar` {
			tmp := strings.SplitN(file.ViewUrl, replaceFrom, 2)
			ext := path.Ext(tmp[1])
			prefix = strings.TrimSuffix(path.Base(tmp[1]), ext)
			newViewURL = tmp[0] + replaceTo + `avatar` + ext
			tmp = strings.SplitN(file.SavePath, replaceFrom, 2)
			newSavePath = tmp[0] + replaceTo + `avatar` + ext
		} else {
			newSavePath = strings.Replace(file.SavePath, replaceFrom, replaceTo, 1)
			newViewURL = strings.Replace(file.ViewUrl, replaceFrom, replaceTo, 1)
		}
		if errMv := storer.Move(file.SavePath, newSavePath); errMv != nil && !os.IsNotExist(errMv) {
			return replaces, errMv
		}
		replaces[file.ViewUrl] = newViewURL
		err = file.SetFields(nil, echo.H{
			`save_path`:  newSavePath,
			`view_url`:   newViewURL,
			`save_name`:  path.Base(newViewURL),
			`used_times`: 1,
		}, db.Cond{`id`: file.Id})
		if err != nil {
			return replaces, err
		}
		thumbM := &dbschema.FileThumb{}
		_, err = thumbM.ListByOffset(nil, nil, 0, -1, db.Cond{`file_id`: file.Id})
		if err != nil {
			return replaces, err
		}
		for _, thumb := range thumbM.Objects() {
			if !strings.Contains(thumb.SavePath, replaceFrom) {
				continue
			}
			var newSavePath, newViewURL string
			if file.FieldName == `avatar` {
				tmp := strings.SplitN(thumb.ViewUrl, replaceFrom, 2)
				suffix := strings.TrimPrefix(path.Base(tmp[1]), prefix)
				newViewURL = tmp[0] + replaceTo + `avatar` + suffix
				tmp = strings.SplitN(file.SavePath, replaceFrom, 2)
				newSavePath = tmp[0] + replaceTo + `avatar` + suffix
			} else {
				newSavePath = strings.Replace(thumb.SavePath, replaceFrom, replaceTo, 1)
				newViewURL = strings.Replace(thumb.ViewUrl, replaceFrom, replaceTo, 1)
			}
			if errMv := storer.Move(thumb.SavePath, newSavePath); errMv != nil && !os.IsNotExist(errMv) {
				return replaces, errMv
			}
			replaces[thumb.ViewUrl] = newViewURL
			err = thumb.SetFields(nil, echo.H{
				`save_path`: newSavePath,
				`view_url`:  newViewURL,
				`save_name`: path.Base(newViewURL),
			}, db.Cond{`id`: thumb.Id})
			if err != nil {
				return replaces, err
			}
		}
	}
	return replaces, err
}

// DeleteByTableID 删除嵌入文件
func (f *Embedded) DeleteByTableID(project string, table string, tableID string) error {
	_, err := f.ListByOffset(nil, nil, 0, -1, db.And(
		db.Cond{`table_id`: tableID},
		db.Cond{`table_name`: table},
		db.Cond{`project`: project},
	))
	if err != nil {
		return err
	}
	var ids []uint64
	for _, row := range f.Objects() {
		err = f.File.SetField(nil, `used_times`, db.Raw(`used_times-1`), db.And(
			db.Cond{`used_times`: db.Gt(0)},
			db.Cond{`id`: db.In(strings.Split(row.FileIds, `,`))},
		))
		if err != nil {
			return err
		}
		ids = append(ids, row.Id)
	}
	if len(ids) > 0 {
		err = f.Delete(nil, db.Cond{`id`: db.In(ids)})
	}
	return err
}

func (f *Embedded) UpdateByFileID(project string, table string, field string, tableID string, fileID uint64) error {
	err := f.File.UpdateUnrelation(project, table, field, tableID, fileID)
	if err != nil {
		return err
	}
	err = f.File.Incr(fileID)
	if err != nil {
		return err
	}
	m := &dbschema.FileEmbedded{}
	err = m.Get(nil, db.And(
		db.Cond{`table_id`: tableID},
		db.Cond{`table_name`: table},
		db.Cond{`field_name`: field},
	))
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		m.Reset()
		m.FieldName = field
		m.TableName = table
		m.Project = project
		m.TableId = tableID
		m.FileIds = fmt.Sprint(fileID)
		_, err = m.Add()
	}
	return err
}

func (f *Embedded) UpdateEmbedded(embedded bool, project string, table string, field string, tableID string, fileIds ...interface{}) (err error) {
	f.base.Begin()
	defer func() {
		f.base.End(err == nil)
	}()
	f.Use(f.base.Tx())
	f.File.Use(f.Trans())

	err = f.File.UpdateUnrelation(project, table, field, tableID, fileIds...)
	if err != nil {
		return err
	}

	m := &dbschema.FileEmbedded{}
	err = m.Use(f.Trans()).Get(nil, db.And(
		db.Cond{`table_id`: tableID},
		db.Cond{`table_name`: table},
		db.Cond{`field_name`: field},
	))
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		if len(fileIds) < 1 {
			return nil
		}
		// 不存在时，添加
		m.Reset()
		m.FieldName = field
		m.TableName = table
		m.Project = project
		m.TableId = tableID
		if embedded {
			m.Embedded = `Y`
		} else {
			m.Embedded = `N`
		}
		m.FileIds = ""
		err = f.File.Incr(fileIds...)
		if err != nil {
			return err
		}
		for _, v := range fileIds {
			m.FileIds += fmt.Sprintf("%v,", v)
		}
		m.FileIds = strings.TrimSuffix(m.FileIds, ",")
		f.FileIds = m.FileIds // 供FileIDs()使用
		_, err = m.Add()
		return err
	}
	isEmpty := len(fileIds) < 1
	if isEmpty { // 删除关联记录
		return f.DeleteByInstance(m)
	}
	var fidsString string
	fidList := make([]string, len(fileIds))
	for k, v := range fileIds {
		s := fmt.Sprint(v)
		fidList[k] = s
		fidsString += s + `,`
	}
	fidsString = strings.TrimSuffix(fidsString, ",")
	if m.FileIds == fidsString {
		return nil
	}
	ids := strings.Split(m.FileIds, ",")
	//新增引用
	err = f.AddFileByIds(fidList, ids...)
	if err != nil {
		return err
	}
	//已删除引用
	err = f.DeleteFileByIds(ids, fidList...)
	if err != nil {
		return err
	}
	m.FileIds = fidsString
	f.FileIds = m.FileIds // 供FileIDs()使用
	err = f.SetField(nil, `file_ids`, m.FileIds, db.Cond{`id`: m.Id})
	return err
}

// RelationEmbeddedFiles 关联嵌入的文件
// @param project 项目名称
// @param table 表名称
// @param field 被嵌入的字段名
// @param tableID 表中行主键ID
// @param v 内容
// @return
// @author AdamShen <swh@admpub.com>
func (f *Embedded) RelationEmbeddedFiles(project string, table string, field string, tableID string, v string) error {
	var (
		files []interface{}
		fids  []interface{} //旧文件ID
	)
	uploadHelper.EmbeddedRes(v, func(file string, fid int64) {
		var exists bool
		if fid > 0 {
			exists = com.InSliceIface(fid, fids)
		} else {
			exists = com.InSliceIface(file, files)
		}
		if exists {
			return
		}
		files = append(files, file)
		if fid > 0 {
			fids = append(fids, fid)
		}
	})
	if len(fids) < 1 && len(files) > 0 {
		fids = f.File.GetIDByViewURLs(files)
	}
	err := f.UpdateEmbedded(true, project, table, field, tableID, fids...)
	return err
}

func (f *Embedded) RelationFiles(project string, table string, field string, tableID string, v string, seperator ...string) error {
	var (
		files []interface{}
		fids  []interface{} //旧文件ID
	)
	//println(`RelationFiles:`, v)
	uploadHelper.RelatedRes(v, func(file string, fid int64) {
		var exists bool
		if fid > 0 {
			exists = com.InSliceIface(fid, fids)
		} else {
			exists = com.InSliceIface(file, files)
		}
		if exists {
			return
		}
		files = append(files, file)
		if fid > 0 {
			fids = append(fids, fid)
		}
	}, seperator...)
	if len(fids) < 1 && len(files) > 0 {
		fids = f.File.GetIDByViewURLs(files)
	}
	err := f.UpdateEmbedded(false, project, table, field, tableID, fids...)
	return err
}
