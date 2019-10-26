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

package listener

import (
	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"

	"github.com/admpub/nging/application/dbschema"
	modelFile "github.com/admpub/nging/application/model/file"
)

var DBI = func() *factory.DBI {
	return dbschema.DBI
}

// New 实例化监听器具
func New(cb Callback, embedded bool, seperators ...string) *FileRelation {
	var seperator string
	if len(seperators) > 0 {
		seperator = seperators[0]
	}
	return &FileRelation{Embedded: embedded, Seperator: seperator, callback: cb}
}

type Callback func(m factory.Model) (tableID string, content string, property *Property)

type Property struct {
	Embedded  bool   // 是否为嵌入图片
	Seperator string // 文件字段中多个文件路径之间的分隔符，空字符串代表为单个文件
	Exit      bool
}

// FileRelation 文件关联数据监听
// FileRelation.SetTable(`table`,`field`).ListenDefault()
type FileRelation struct {
	TableName string   // 数据表名称
	FieldName string   // 数据表字段名
	Embedded  bool     // 是否为嵌入图片
	Seperator string   // 文件字段中多个文件路径之间的分隔符，空字符串代表为单个文件
	callback  Callback //根据模型获取表行ID和内容
	dbi       *factory.DBI
}

func (f *FileRelation) SetSeperator(seperator string) *FileRelation {
	f.Seperator = seperator
	return f
}

func (f *FileRelation) Callback() Callback {
	return f.callback
}

func (f *FileRelation) SetTable(table string, field string) *FileRelation {
	f.TableName = table
	f.FieldName = field
	return f
}

func (f *FileRelation) SetDBI(dbi *factory.DBI) *FileRelation {
	f.dbi = dbi
	return f
}

func (f *FileRelation) SetEmbedded(embedded bool) *FileRelation {
	f.Embedded = embedded
	return f
}

func (f *FileRelation) ListenDefault() *FileRelation {
	return f.Listen(`created`, `updated`, `deleted`)
}

func (f *FileRelation) attachUpdateEvent(event string) func(m factory.Model, editColumns ...string) error {
	seperator := f.Seperator
	embedded := f.Embedded
	callback := f.callback
	return func(m factory.Model, editColumns ...string) error {
		if len(editColumns) > 0 && !com.InSlice(f.FieldName, editColumns) {
			return nil
		}
		fileM := modelFile.NewEmbedded(m.Context())
		tableID, content, property := callback(m)
		if property != nil {
			if property.Exit {
				return nil
			}
			seperator = property.Seperator
			embedded = property.Embedded
		}
		//println(event+`=========`, f.TableName, f.FieldName, tableID, content)
		return fileM.Updater(f.TableName, f.FieldName, tableID).SetSeperator(seperator).Handle(event, content, embedded)
	}
}

func (f *FileRelation) attachEvent(event string) func(m factory.Model, editColumns ...string) error {
	seperator := f.Seperator
	embedded := f.Embedded
	callback := f.callback
	return func(m factory.Model, _ ...string) error {
		fileM := modelFile.NewEmbedded(m.Context())
		tableID, content, property := callback(m)
		if property != nil {
			if property.Exit {
				return nil
			}
			seperator = property.Seperator
			embedded = property.Embedded
		}
		return fileM.Updater(f.TableName, f.FieldName, tableID).SetSeperator(seperator).Handle(event, content, embedded)
	}
}

func (f *FileRelation) DBI() *factory.DBI {
	dbi := f.dbi
	if dbi == nil {
		dbi = DBI()
	}
	return dbi
}

func (f *FileRelation) Listen(events ...string) *FileRelation {
	dbi := f.DBI()
	for _, event := range events {
		switch event {
		case `updating`, `updated`:
			dbi.On(event, f.attachUpdateEvent(event), f.TableName)
		default:
			dbi.On(event, f.attachEvent(event), f.TableName)
		}
	}
	return f
}

func (f *FileRelation) On(event string, h factory.EventHandler) *FileRelation {
	f.DBI().On(event, h, f.TableName)
	return f
}
