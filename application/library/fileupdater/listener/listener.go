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

type Callback func(m factory.Model) (tableID string, content string)

// FileRelation 文件关联数据监听
// FileRelation.SetTable(`table`,`field`).ListenDefault()
type FileRelation struct {
	TableName string   // 数据表名称
	FieldName string   // 数据表字段名
	Embedded  bool     // 是否为嵌入图片
	Seperator string   // 文件字段中多个文件路径之间的分隔符，空字符串代表为单个文件
	callback  Callback //根据模型获取表行ID和内容
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

func (f *FileRelation) SetEmedded(embedded bool) *FileRelation {
	f.Embedded = embedded
	return f
}

func (f *FileRelation) ListenDefault() {
	f.Listen(`created`, `updating`, `deleted`)
}

func (f *FileRelation) Listen(events ...string) {
	for _, event := range events {
		switch event {
		case `updating`, `updated`:
			DBI().On(event, func(m factory.Model, editColumns ...string) error {
				if len(editColumns) > 0 && !com.InSlice(f.FieldName, editColumns) {
					return nil
				}
				fileM := modelFile.NewEmbedded(m.Context())
				tableID, content := f.callback(m)
				return fileM.Updater(f.TableName, f.FieldName, tableID).SetSeperator(f.Seperator).Add(content, f.Embedded)
			}, f.TableName)
		default:
			DBI().On(event, func(m factory.Model, _ ...string) error {
				fileM := modelFile.NewEmbedded(m.Context())
				tableID, content := f.callback(m)
				return fileM.Updater(f.TableName, f.FieldName, tableID).SetSeperator(f.Seperator).Add(content, f.Embedded)
			}, f.TableName)
		}
	}
}
