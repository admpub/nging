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
	"github.com/admpub/nging/v3/application/dbschema"
	"github.com/admpub/nging/v3/application/library/fileupdater"
	modelFile "github.com/admpub/nging/v3/application/model/file"
)

var DBI = func() *factory.DBI {
	return dbschema.DBI
}

// New 实例化监听器具
func New(cb fileupdater.CallbackFunc, embedded bool, seperators ...string) *FileRelation {
	var seperator string
	if len(seperators) > 0 {
		seperator = seperators[0]
	}
	return &FileRelation{
		Options: &fileupdater.Options{
			Embedded:  embedded,
			Seperator: seperator,
			Callback:  cb,
		},
	}
}

// NewWithOptions 实例化监听器具
func NewWithOptions(options *fileupdater.Options) *FileRelation {
	return &FileRelation{
		Options: options,
	}
}

type Callback = fileupdater.CallbackFunc

// FileRelation 文件关联数据监听
// FileRelation.SetTable(`table`,`field`).ListenDefault()
type FileRelation struct {
	*fileupdater.Options
	dbi *factory.DBI
}

func (f *FileRelation) SetSeperator(seperator string) *FileRelation {
	f.Seperator = seperator
	return f
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
	return f.Listen(factory.EventCreated, factory.EventUpdated, factory.EventDeleting)
}

func (f *FileRelation) attachUpdateEvent(event string) func(m factory.Model, editColumns ...string) error {
	seperator := f.Seperator
	embedded := f.Embedded
	callback := f.Callback
	return func(m factory.Model, editColumns ...string) error {
		if len(editColumns) > 0 && !com.InSlice(f.FieldName, editColumns) {
			return nil
		}
		fileM := modelFile.NewEmbedded(m.Context())
		tableID, content, property := callback(m)
		var updater func(string, string) error
		if property != nil {
			if property.Exit() {
				return nil
			}
			if property.Seperator().Valid {
				seperator = property.Seperator().String
			}
			if property.Embedded().Valid {
				embedded = property.Embedded().Bool
			}
			updater = property.Updater()
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
	callback := f.Callback
	return func(m factory.Model, _ ...string) error {
		fileM := modelFile.NewEmbedded(m.Context())
		tableID, content, property := callback(m)
		var updater func(string, string) error
		if property != nil {
			if property.Exit() {
				return nil
			}
			if property.Seperator().Valid {
				seperator = property.Seperator().String
			}
			if property.Embedded().Valid {
				embedded = property.Embedded().Bool
			}
			updater = property.Updater()
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
	callback := f.Callback
	return func(m factory.Model, _ ...string) error {
		fileM := modelFile.NewEmbedded(m.Context())
		tableID, content, property := callback(m)
		if property != nil {
			if property.Exit() {
				return nil
			}
			if property.Seperator().Valid {
				seperator = property.Seperator().String
			}
			if property.Embedded().Valid {
				embedded = property.Embedded().Bool
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
		case factory.EventUpdating, factory.EventUpdated:
			f.On(event, f.attachUpdateEvent(event))
		case factory.EventDeleting, factory.EventDeleted:
			f.On(event, f.attachDeleteEvent(event))
		default:
			f.On(event, f.attachEvent(event))
		}
	}
	return f
}

func (f *FileRelation) On(event string, h factory.EventHandler) *FileRelation {
	f.DBI().On(event, h, f.TableName)
	log.Debug(color.MagentaString(`listener.`+event+`:`), f.TableName+`.`+f.FieldName)
	RecordUpdaterInfo(``, f.TableName, f.FieldName, f.Seperator, f.Embedded, f.SameFields...)
	return f
}

func (f *FileRelation) OnRead(event string, h factory.EventReadHandler) *FileRelation {
	f.DBI().OnRead(event, h, f.TableName)
	log.Debug(color.MagentaString(`listener.`+event+`:`), f.TableName+`.`+f.FieldName)
	RecordUpdaterInfo(``, f.TableName, f.FieldName, f.Seperator, f.Embedded, f.SameFields...)
	return f
}
