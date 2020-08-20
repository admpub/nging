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

package subdir

import (
	"fmt"
	"sort"
	"strings"

	"github.com/admpub/nging/application/library/fileupdater"

	"github.com/admpub/nging/application/registry/upload/checker"
	"github.com/admpub/nging/application/registry/upload/table"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

type ThumbSize struct {
	AutoCrop bool
	Width    float64
	Height   float64
}

func (t ThumbSize) String() string {
	return fmt.Sprintf("%vx%v", t.Width, t.Height)
}

type FieldInfo struct {
	Key     string
	Name    string
	Thumb   []ThumbSize
	checker checker.Checker
}

func (info *FieldInfo) String() string {
	if len(info.Name) > 0 {
		return info.Name
	}
	return info.Key
}

func (info *FieldInfo) IsEmpty() bool {
	if info == EmptyFieldInfo {
		return true
	}
	return len(info.Key) == 0 && len(info.Name) == 0 && len(info.Thumb) == 0 && info.checker == nil
}

func (info *FieldInfo) AddChecker(checkerFn checker.Checker) *FieldInfo {
	if checkerFn == nil {
		return info
	}
	if info.checker == nil {
		info.checker = checkerFn
		return info
	}
	oldChecker := info.checker
	info.checker = func(ctx echo.Context, tbl table.TableInfoStorer) (subdir string, name string, err error) {
		subdir, name, err = oldChecker(ctx, tbl)
		if err != nil {
			return
		}
		subdir2, name2, err2 := checkerFn(ctx, tbl)
		if err2 != nil {
			err = err2
		}
		if len(subdir2) > 0 {
			subdir = subdir2
		}
		if len(name2) > 0 {
			name = name2
		}
		return
	}
	return info
}

var FileListenerGenerator func() fileupdater.Listener

// SubdirInfo 子目录信息
type SubdirInfo struct {
	Allowed     bool   // 是否允许上传
	Key         string // 子目录文件夹名
	Name        string // 子目录中文名称
	NameEN      string // 子目录英文名称
	Description string // 子目录简述

	tableName    string
	fieldInfos   map[string]*FieldInfo
	checker      checker.Checker
	fileListener fileupdater.Listener
}

func NewSubdirInfo(key, name string, checkers ...checker.Checker) *SubdirInfo {
	var (
		checkerFn    checker.Checker
		fileListener fileupdater.Listener
	)
	if len(checkers) > 0 {
		checkerFn = checkers[0]
	}
	if FileListenerGenerator != nil {
		fileListener = FileListenerGenerator()
	}
	return &SubdirInfo{
		Allowed:      true,
		Key:          key,
		Name:         name,
		checker:      checkerFn,
		fileListener: fileListener,
	}
}

var EmptyFieldInfo = &FieldInfo{}

func (i *SubdirInfo) GetField(field string) *FieldInfo {
	if len(field) == 0 {
		return &FieldInfo{Key: i.Key, Name: i.Name}
	}
	if i.fieldInfos == nil {
		return EmptyFieldInfo
	}
	if info, ok := i.fieldInfos[field]; ok {
		return info
	}
	return EmptyFieldInfo
}

func (i *SubdirInfo) SetFileListener(listener fileupdater.Listener) *SubdirInfo {
	i.fileListener = listener
	return i
}

func (i *SubdirInfo) FileListener() fileupdater.Listener {
	if i.fileListener == nil && FileListenerGenerator != nil {
		i.fileListener = FileListenerGenerator()
		i.fileListener.SetTableName(i.tableName)
	}
	return i.fileListener
}

func (i *SubdirInfo) CopyFrom(other *SubdirInfo) *SubdirInfo {
	i.Allowed = other.Allowed
	if len(other.Key) > 0 {
		i.Key = other.Key
	}
	if len(other.Name) > 0 {
		i.Name = other.Name
	}
	if len(other.NameEN) > 0 {
		i.NameEN = other.NameEN
	}
	if len(other.Description) > 0 {
		i.Description = other.Description
	}
	if len(other.tableName) > 0 {
		i.SetTableName(other.tableName)
	}
	if other.checker != nil {
		if i.checker != nil {
			oldChecker := i.checker
			i.checker = func(ctx echo.Context, tbl table.TableInfoStorer) (subdir string, name string, err error) {
				subdir, name, err = oldChecker(ctx, tbl)
				if err != nil {
					return
				}
				subdir2, name2, err2 := other.checker(ctx, tbl)
				if err2 != nil {
					err = err2
				}
				if len(subdir2) > 0 {
					subdir = subdir2
				}
				if len(name2) > 0 {
					name = name2
				}
				return
			}
		} else {
			i.checker = other.checker
		}
	}
	if len(other.fieldInfos) > 0 {
		if i.fieldInfos == nil {
			i.fieldInfos = make(map[string]*FieldInfo)
		}
		for k, in := range other.fieldInfos {
			i.fieldInfos[k] = in
		}
	}
	return i
}

func (i *SubdirInfo) String() string {
	if len(i.Name) > 0 {
		return i.Name
	}
	return i.Key
}

func (i *SubdirInfo) TableName() string {
	if len(i.tableName) == 0 {
		r := strings.SplitN(i.Key, `.`, 2)
		switch len(r) {
		case 2:
			i.SetFieldName(r[1])
			fallthrough
		case 1:
			i.SetTableName(r[0])
		}
	}
	return i.tableName
}

func (i *SubdirInfo) ValidFieldName(fieldName string) bool {
	if i.fieldInfos == nil {
		return false
	}
	_, ok := i.fieldInfos[fieldName]
	return ok
}

func (i *SubdirInfo) ThumbSize(fieldNames ...string) []ThumbSize {
	if i.fieldInfos == nil {
		return nil
	}
	if len(fieldNames) > 0 {
		info, ok := i.fieldInfos[fieldNames[0]]
		if ok {
			return info.Thumb
		}
	}
	for _, info := range i.fieldInfos {
		for _, thumbSize := range info.Thumb {
			if thumbSize.Width > 0 && thumbSize.Height > 0 {
				return info.Thumb
			}
		}
	}
	return nil
}

func (i *SubdirInfo) AutoCropThumbSize(fieldName string) []ThumbSize {
	r := []ThumbSize{}
	sizes := i.ThumbSize(fieldName)
	for _, size := range sizes {
		if size.AutoCrop {
			r = append(r, size)
		}
	}
	return r
}

func (i *SubdirInfo) FieldNames() []string {
	fieldNames := make([]string, len(i.fieldInfos))
	if i.fieldInfos == nil {
		return fieldNames
	}
	var index int
	for fieldName := range i.fieldInfos {
		fieldNames[index] = fieldName
		index++
	}
	sort.Strings(fieldNames)
	return fieldNames
}

func (i *SubdirInfo) FieldInfos() []echo.KV {
	fieldNames := i.FieldNames()
	r := make([]echo.KV, len(fieldNames))
	for index, fieldName := range fieldNames {
		info := i.fieldInfos[fieldName]
		r[index] = echo.KV{K: fieldName, V: info.Name}
	}
	return r
}

func (i *SubdirInfo) SetTableName(tableName string) *SubdirInfo {
	i.tableName = tableName
	table2dir[i.tableName] = i.Key
	return i
}

func (i *SubdirInfo) SetTable(tableName string, fieldNames ...string) *SubdirInfo {
	i.SetTableName(tableName)
	if len(fieldNames) > 0 {
		i.SetFieldName(fieldNames...)
	}
	return i
}

func (i *SubdirInfo) SetFieldName(fieldNames ...string) *SubdirInfo {
	i.fieldInfos = map[string]*FieldInfo{}
	for _, fieldName := range fieldNames {
		i.AddFieldName(fieldName)
	}
	return i
}

var AutoCropPrefix = `[auto]`

// parseFieldInfo fieldName:fieldDescription:thumbWidth x thumbHeight,thumbWidth x thumbHeight
// example parseFieldInfo(`user:用户:200x200,300x600`)
func (i *SubdirInfo) parseFieldInfo(field string) (fieldName string, fieldText string, thumbSize []ThumbSize) {
	r := strings.SplitN(field, ":", 3)
	switch len(r) {
	case 3:
		r[2] = strings.TrimSpace(r[2])
		for _, size := range strings.Split(r[2], `,`) {
			size := strings.TrimSpace(size)
			if len(size) == 0 {
				continue
			}
			autoCrop := strings.HasPrefix(size, AutoCropPrefix)
			if autoCrop {
				size = strings.TrimPrefix(size, AutoCropPrefix)
			}
			sz := strings.SplitN(size, `x`, 2)
			ts := ThumbSize{AutoCrop: autoCrop}
			switch len(sz) {
			case 2:
				sz[1] = strings.TrimSpace(sz[1])
				sz[0] = strings.TrimSpace(sz[0])
				ts.Height = param.AsFloat64(sz[1])
				ts.Width = param.AsFloat64(sz[0])
			case 1:
				sz[0] = strings.TrimSpace(sz[0])
				ts.Width = param.AsFloat64(sz[0])
				ts.Height = ts.Width
			}
			if ts.Width > 0 && ts.Height > 0 {
				if autoCrop {

				}
				thumbSize = append(thumbSize, ts)
			}
		}
		fallthrough
	case 2:
		r[1] = strings.TrimSpace(r[1])
		fieldText = r[1]
		fallthrough
	case 1:
		r[0] = strings.TrimSpace(r[0])
		fieldName = r[0]
	}
	return
}

func (i *SubdirInfo) AddFieldName(fieldName string, checkers ...checker.Checker) *SubdirInfo {
	var (
		fieldText string
		checkerFn checker.Checker
		thumbSize []ThumbSize
	)
	if len(checkers) > 0 {
		checkerFn = checkers[0]
	}
	fieldName, fieldText, thumbSize = i.parseFieldInfo(fieldName)
	info, ok := i.fieldInfos[fieldName]
	if !ok {
		i.fieldInfos[fieldName] = &FieldInfo{
			Key:     fieldName,
			Name:    fieldText,
			Thumb:   thumbSize,
			checker: checkerFn,
		}
	} else {
		if len(fieldText) > 0 {
			info.Name = fieldText
		}
		info.AddChecker(checkerFn)
		if len(thumbSize) > 0 {
			info.Thumb = append(info.Thumb, thumbSize...)
		}
	}
	return i
}

func (i *SubdirInfo) GetNameEN() string {
	if len(i.NameEN) > 0 {
		return i.NameEN
	}
	return i.Name
}

func (i *SubdirInfo) SetChecker(checkerFn checker.Checker, fieldNames ...string) *SubdirInfo {
	if len(fieldNames) > 0 {
		if i.fieldInfos == nil {
			i.fieldInfos = make(map[string]*FieldInfo)
		}
		info, ok := i.fieldInfos[fieldNames[0]]
		if !ok {
			i.AddFieldName(fieldNames[0])
			info, _ = i.fieldInfos[fieldNames[0]]
			//panic(`not found: ` + i.Key + `.` + fieldNames[0])
		}
		info.AddChecker(checkerFn)
		return i
	}
	i.checker = checkerFn
	return i
}

func (i *SubdirInfo) Checker() checker.Checker {
	return i.checker
}

func (i *SubdirInfo) MustChecker() checker.Checker {
	return func(ctx echo.Context, tab table.TableInfoStorer) (subdir string, name string, err error) {
		tab.SetTableName(i.TableName())
		if !i.ValidFieldName(tab.FieldName()) {
			err = table.ErrInvalidFieldName
		}
		//echo.Dump(echo.H{`field`: tab.FieldName(), `fields`: i.fieldInfos})
		var checkerFn checker.Checker
		info, ok := i.fieldInfos[tab.FieldName()]
		if ok {
			checkerFn = info.checker
		}
		if checkerFn == nil {
			if i.checker != nil {
				checkerFn = i.checker
			} else {
				checkerFn = checker.Default
			}
		}
		subdir, name, err = checkerFn(ctx, tab)
		return
	}
}

func (i *SubdirInfo) SetAllowed(allowed bool) *SubdirInfo {
	i.Allowed = allowed
	return i
}

func (i *SubdirInfo) SetName(name string) *SubdirInfo {
	i.Name = name
	return i
}

func (i *SubdirInfo) SetNameEN(nameEN string) *SubdirInfo {
	i.NameEN = nameEN
	return i
}

func (i *SubdirInfo) SetDescription(description string) *SubdirInfo {
	i.Description = description
	return i
}
