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

package upload

import (
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/registry/upload/table"
)

// SubdirInfo 子目录信息
type SubdirInfo struct {
	Allowed     bool
	Key         string
	Name        string
	NameEN      string
	Description string

	tableName  string
	fieldNames []string
	fieldDescs []string
	checker    Checker
}

func (i *SubdirInfo) String() string {
	if len(i.Name) > 0 {
		return i.Name
	}
	return i.Key
}

func (i *SubdirInfo) TableName() string {
	if len(i.tableName) == 0 {
		r := strings.SplitN(i.Key, `-`, 2)
		switch len(r) {
		case 2:
			i.SetFieldName(r[1])
			fallthrough
		case 1:
			i.tableName = r[0]
		}
	}
	return i.tableName
}

func (i *SubdirInfo) ValidFieldName(fieldName string) bool {
	return com.InSlice(fieldName, i.fieldNames)
}

func (i *SubdirInfo) FieldNames() []string {
	return i.fieldNames
}

func (i *SubdirInfo) FieldInfos() []echo.KV {
	r := make([]echo.KV, len(i.fieldNames))
	for index, fieldName := range i.fieldNames {
		r[index] = echo.KV{K: fieldName, V: i.fieldDescs[index]}
	}
	return r
}

func (i *SubdirInfo) SetTableName(tableName string) *SubdirInfo {
	i.tableName = tableName
	return i
}

func (i *SubdirInfo) SetTable(tableName string, fieldNames ...string) *SubdirInfo {
	i.tableName = tableName
	if len(fieldNames) > 0 {
		i.SetFieldName(fieldNames...)
	}
	return i
}

func (i *SubdirInfo) SetFieldName(fieldNames ...string) *SubdirInfo {
	i.fieldNames = []string{}
	for _, fieldName := range fieldNames {
		i.AddFieldName(fieldName)
	}
	return i
}

func (i *SubdirInfo) parseFieldInfo(field string) (fieldName string, fieldText string) {
	r := strings.SplitN(field, ":", 2)
	switch len(r) {
	case 2:
		fieldText = r[1]
		fallthrough
	case 1:
		fieldName = r[0]
	}
	return
}

func (i *SubdirInfo) AddFieldName(fieldName string) *SubdirInfo {
	if !com.InSlice(fieldName, i.fieldNames) {
		var fieldText string
		fieldName, fieldText = i.parseFieldInfo(fieldName)
		i.fieldNames = append(i.fieldNames, fieldName)
		i.fieldDescs = append(i.fieldDescs, fieldText)
	}
	return i
}

func (i *SubdirInfo) GetNameEN() string {
	if len(i.NameEN) > 0 {
		return i.NameEN
	}
	return i.Name
}

func (i *SubdirInfo) SetChecker(checker Checker) *SubdirInfo {
	i.checker = checker
	return i
}

func (i *SubdirInfo) Checker() Checker {
	return i.checker
}

func (i *SubdirInfo) MustChecker() Checker {
	if i.checker != nil {
		return i.checker
	}
	return func(ctx echo.Context, tab table.TableInfoStorer) (subdir string, name string, err error) {
		subdir, name, err = DefaultChecker(ctx, tab)
		if err != nil {
			return
		}
		tab.SetTableName(i.TableName())
		if len(tab.FieldName()) > 0 && !i.ValidFieldName(tab.FieldName()) {
			err = table.ErrInvalidFieldName
		}
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
