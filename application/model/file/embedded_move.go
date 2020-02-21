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
	"os"
	"path"
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/dbschema"
	uploadStorer "github.com/admpub/nging/application/registry/upload/driver"
)

func (f *Embedded) getSeperator(content string) string {
	for _, v := range content {
		if v == '/' {
			return "/"
		}
		if v == '\\' {
			return `\`
		}
	}
	return "/"
}

func (f *Embedded) renameFile(
	tableName string, fieldName string,
	savePath string, viewURL string,
	replaceFrom string, replaceTo string,
	savePathSep string, prefixes ...string) (newSavePath string, newViewURL string, prefix string) {
	viewURLSep := `/`
	viewURLTmp := strings.SplitN(viewURL, viewURLSep+replaceFrom+viewURLSep, 2)
	if len(viewURLTmp) != 2 {
		newViewURL = viewURL
		newSavePath = savePath
		return
	}
	pos := strings.LastIndex(viewURLTmp[0], viewURLSep)
	if viewURLTmp[0][pos+1:] != tableName {
		viewURLTmp[0] = viewURLTmp[0][0:pos] + viewURLSep + tableName
	}
	savePathTmp := strings.SplitN(savePath, savePathSep+replaceFrom+savePathSep, 2)
	pos = strings.LastIndex(savePathTmp[0], savePathSep)
	if savePathTmp[0][pos+1:] != tableName {
		savePathTmp[0] = savePathTmp[0][0:pos] + savePathSep + tableName
	}
	if fieldName == `avatar` {
		ext := path.Ext(viewURLTmp[1])
		suffix := ext
		if len(prefixes) > 0 { // 缩略图 (表file_thumb中的数据)
			suffix = strings.TrimPrefix(path.Base(viewURLTmp[1]), prefixes[0])
		} else { // 表file中的数据
			prefix = strings.TrimSuffix(path.Base(viewURLTmp[1]), ext)
		}
		newViewURL = viewURLTmp[0] + viewURLSep + replaceTo + viewURLSep + `avatar` + suffix
		newSavePath = savePathTmp[0] + savePathSep + replaceTo + savePathSep + `avatar` + suffix
	} else {
		newViewURL = viewURLTmp[0] + viewURLSep + replaceTo + viewURLSep + viewURLTmp[1]
		newSavePath = savePathTmp[0] + savePathSep + replaceTo + savePathSep + savePathTmp[1]
	}
	return
}

func (f *Embedded) MoveFileToOwner(tableName string, fileIDs []uint64, ownerID string) (map[string]string, error) {
	replaces := make(map[string]string)
	if len(fileIDs) == 0 {
		return replaces, nil
	}
	_, err := f.File.ListByOffset(nil, nil, 0, -1, db.Cond{`id`: db.In(fileIDs)})
	if err != nil {
		return replaces, err
	}
	replaceFrom := `0`
	replaceTo := ownerID
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
		savePathSep := f.getSeperator(file.SavePath)
		newSavePath, newViewURL, prefix := f.renameFile(tableName, file.FieldName, file.SavePath, file.ViewUrl, replaceFrom, replaceTo, savePathSep)
		if newSavePath != file.SavePath {
			if errMv := storer.Move(file.SavePath, newSavePath); errMv != nil && !os.IsNotExist(errMv) {
				return replaces, errMv
			}
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
		thumbM := &dbschema.NgingFileThumb{}
		_, err = thumbM.ListByOffset(nil, nil, 0, -1, db.Cond{`file_id`: file.Id})
		if err != nil {
			return replaces, err
		}
		for _, thumb := range thumbM.Objects() {
			if !strings.Contains(thumb.SavePath, replaceFrom) {
				continue
			}
			newSavePath, newViewURL, _ := f.renameFile(tableName, file.FieldName, thumb.SavePath, thumb.ViewUrl, replaceFrom, replaceTo, savePathSep, prefix)
			if newSavePath != thumb.SavePath {
				if errMv := storer.Move(thumb.SavePath, newSavePath); errMv != nil && !os.IsNotExist(errMv) {
					return replaces, errMv
				}
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
