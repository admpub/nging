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

	"github.com/admpub/color"
	"github.com/admpub/log"
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

// FileRelation 文件关联数据监听
// FileRelation.SetTable(`table`,`field`).ListenDefault()
type FileRelation struct {
	TableName string   // 数据表名称
	FieldName string   // 数据表字段名
	SameFields []string   // 数据表类似字段名
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

func (f *FileRelation) SetTable(table string, field string, samesFields ...string) *FileRelation {
	f.TableName = table
	f.FieldName = field
	f.SameFields = samesFields
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
	return f.Listen(`created`, `updated`, `deleting`)
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
		var updater func(string, string) error
		if property != nil {
			if property.exit {
				return nil
			}
			if property.seperator.Valid {
				seperator = property.seperator.String
			}
			if property.embedded.Valid {
				embedded = property.embedded.Bool
			}
			updater = property.updater
		}
		//println(event+`=========`, f.TableName, f.FieldName, tableID, content)
		if updater == nil {
			return fileM.Updater(f.TableName, f.FieldName, tableID).SetSeperator(seperator).Handle(event, &content, embedded)
		}
		rawContent := content
		err := fileM.Updater(f.TableName, f.FieldName, tableID).SetSeperator(seperator).Handle(event, &content, embedded)
		if err != nil {
			return err
		}
		if rawContent != content {
			err = updater(f.FieldName, content)
		}
		return err
	}
}

func (f *FileRelation) attachEvent(event string) func(m factory.Model, editColumns ...string) error {
	seperator := f.Seperator
	embedded := f.Embedded
	callback := f.callback
	return func(m factory.Model, _ ...string) error {
		fileM := modelFile.NewEmbedded(m.Context())
		tableID, content, property := callback(m)
		var updater func(string, string) error
		if property != nil {
			if property.exit {
				return nil
			}
			if property.seperator.Valid {
				seperator = property.seperator.String
			}
			if property.embedded.Valid {
				embedded = property.embedded.Bool
			}
			updater = property.updater
		}
		if updater == nil {
			return fileM.Updater(f.TableName, f.FieldName, tableID).SetSeperator(seperator).Handle(event, &content, embedded)
		}
		rawContent := content
		err := fileM.Updater(f.TableName, f.FieldName, tableID).SetSeperator(seperator).Handle(event, &content, embedded)
		if err != nil {
			return err
		}
		if rawContent != content {
			err = updater(f.FieldName, content)
		}
		return err
	}
}

func (f *FileRelation) attachDeleteEvent(event string) func(m factory.Model, editColumns ...string) error {
	seperator := f.Seperator
	embedded := f.Embedded
	callback := f.callback
	return func(m factory.Model, _ ...string) error {
		fileM := modelFile.NewEmbedded(m.Context())
		tableID, content, property := callback(m)
		if property != nil {
			if property.exit {
				return nil
			}
			if property.seperator.Valid {
				seperator = property.seperator.String
			}
			if property.embedded.Valid {
				embedded = property.embedded.Bool
			}
		}
		return fileM.Updater(f.TableName, f.FieldName, tableID).SetSeperator(seperator).Handle(event, &content, embedded)
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
	for _, event := range events {
		switch event {
		case `updating`, `updated`:
			f.On(event, f.attachUpdateEvent(event))
		case `deleting`, `deleted`:
			f.On(event, f.attachDeleteEvent(event))
		default:
			f.On(event, f.attachEvent(event))
		}
	}
	return f
}

func (f *FileRelation) On(event string, h factory.EventHandler) *FileRelation {
	f.DBI().On(event, h, f.TableName)
	log.Info(color.MagentaString(`listener.`+event+`:`), f.TableName+`.`+f.FieldName)
	RecordUpdaterInfo(``, f.TableName, f.FieldName, f.Seperator, f.Embedded, f.SameFields...)
	return f
}

func (f *FileRelation) OnRead(event string, h factory.EventReadHandler) *FileRelation {
	f.DBI().OnRead(event, h, f.TableName)
	log.Info(color.MagentaString(`listener.`+event+`:`), f.TableName+`.`+f.FieldName)
	RecordUpdaterInfo(``, f.TableName, f.FieldName, f.Seperator, f.Embedded, f.SameFields...)
	return f
}
