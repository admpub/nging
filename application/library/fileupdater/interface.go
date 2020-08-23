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

package fileupdater

import (
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
)

type Reler interface {
	RelationFiles(project string, table string, field string, tableID string, content string, seperator ...string) error
	RelationEmbeddedFiles(project string, table string, field string, tableID string, content string) error
	DeleteByTableID(project string, table string, tableID string) error
	FileIDs() []uint64
	MoveFileToOwner(table string, fileIDs []uint64, ownerID string) (replaces map[string]string, err error)
	Context() echo.Context
}

type Listener interface {
	SetTableName(table string) Listener
	Listen()
	ListenByOptions(fieldName string, setters ...OptionSetter)
	Add(fieldName string, setters ...OptionSetter) Listener
	BuildOptions(fieldName string, setters ...OptionSetter) *Options
	Delete(fieldNames ...string) Listener
}

type CallbackFunc func(m factory.Model) (tableID string, content string, property *Property)
