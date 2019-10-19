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

	tableName string
	fieldName string
	checker   Checker
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
			i.fieldName = r[1]
			fallthrough
		case 1:
			i.tableName = r[0]
		}
	}
	return i.tableName
}

func (i *SubdirInfo) FieldName() string {
	return i.fieldName
}

func (i *SubdirInfo) SetTableName(tableName string) *SubdirInfo {
	i.tableName = tableName
	return i
}

func (i *SubdirInfo) SetTable(tableName string, fieldName string) *SubdirInfo {
	i.tableName = tableName
	i.fieldName = fieldName
	return i
}

func (i *SubdirInfo) SetFieldName(fieldName string) *SubdirInfo {
	i.fieldName = fieldName
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
		tab.SetTableName(i.TableName())
		tab.SetFieldName(i.FieldName())
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
