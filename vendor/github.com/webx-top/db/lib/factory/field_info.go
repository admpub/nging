package factory

import (
	"sort"
)

type ModelInstancer func(connID int) Model
type ModelInstancers map[string]ModelInstancer

type FieldInfo struct {
	//以下为数据库中的信息
	Name          string   `json:"name" xml:"name" bson:"name"`                            //字段名
	DataType      string   `json:"dataType" xml:"dataType" bson:"dataType"`                //数据库中数据类型
	Unsigned      bool     `json:"unsigned" xml:"unsigned" bson:"unsigned"`                //是否是无符号类型
	PrimaryKey    bool     `json:"primaryKey" xml:"primaryKey" bson:"primaryKey"`          //是否是主键
	AutoIncrement bool     `json:"autoIncrement" xml:"autoIncrement" bson:"autoIncrement"` //是否是自增字段
	Min           int      `json:"min" xml:"min" bson:"min"`                               //最小值
	Max           int      `json:"max" xml:"max" bson:"max"`                               //最大值
	Precision     int      `json:"precision" xml:"precision" bson:"precision"`             //小数精度(小数保留位数)
	MaxSize       int      `json:"maxSize" xml:"maxSize" bson:"maxSize"`                   //最大尺寸
	Options       []string `json:"options" xml:"options" bson:"options"`                   //选项值
	DefaultValue  string   `json:"defaultValue" xml:"defaultValue" bson:"defaultValue"`    //默认值
	Comment       string   `json:"comment" xml:"comment" bson:"comment"`                   //备注

	//以下为Golang中的信息
	GoType string `json:"goType" xml:"goType" bson:"goType"` //Golang数据类型
	GoName string `json:"goName" xml:"goName" bson:"goName"` //Golang字段名
}

type FieldValidator map[string]map[string]*FieldInfo

func (f FieldValidator) ExistField(table string, field string) bool {
	if tb, ok := f[table]; ok {
		_, ok = tb[field]
		return ok
	}
	return false
}

func (f FieldValidator) ExistTable(table string) bool {
	_, ok := f[table]
	return ok
}

func (f FieldValidator) FieldList(table string, excludeField ...string) []string {
	fields := []string{}
	if tb, ok := f[table]; ok {
		for field := range tb {
			var exists bool
			for _, ex := range excludeField {
				if field == ex {
					exists = true
					break
				}
			}
			if exists {
				continue
			}
			fields = append(fields, field)
		}
	}
	return fields
}

func (f FieldValidator) SortedFieldList(table string, excludeField ...string) []string {
	fields := f.FieldList(table, excludeField...)
	sort.Strings(fields)
	return fields
}

func (f FieldValidator) SortedFieldLists(table string, excludeField ...string) []interface{} {
	fields := f.SortedFieldList(table, excludeField...)
	returns := make([]interface{}, len(fields))
	for i, v := range fields {
		returns[i] = v
	}
	return returns
}

func (f FieldValidator) FieldLists(table string, excludeField ...string) []interface{} {
	fields := []interface{}{}
	if tb, ok := f[table]; ok {
		for field := range tb {
			var exists bool
			for _, ex := range excludeField {
				if field == ex {
					exists = true
					break
				}
			}
			if exists {
				continue
			}
			fields = append(fields, field)
		}
	}
	return fields
}

var (
	// Fields {table:{field:FieldInfo}}
	Fields FieldValidator = map[string]map[string]*FieldInfo{}
	// Models {StructName:ModelInstancer}
	Models = ModelInstancers{}
	// TableNamers {table:NewName}
	TableNamers = map[string]func(obj interface{}) string{}
	// DefaultTableNamer default
	DefaultTableNamer = func(table string) func(obj interface{}) string {
		return func(obj interface{}) string {
			return table
		}
	}
)

func TableNamerRegister(namers map[string]func(obj interface{}) string) {
	for table, namer := range namers {
		TableNamers[table] = namer
	}
}

func TableNamerGet(table string) func(obj interface{}) string {
	if namer, ok := TableNamers[table]; ok {
		return namer
	}
	return DefaultTableNamer(table)
}

func FieldRegister(tables map[string]map[string]*FieldInfo) {
	for table, info := range tables {
		Fields[table] = info
	}
}

func ModelRegister(instancers map[string]ModelInstancer) {
	for structName, instancer := range instancers {
		Models[structName] = instancer
	}
}

func ExistField(table string, field string) bool {
	return Fields.ExistField(table, field)
}

func ExistTable(table string) bool {
	return Fields.ExistTable(table)
}

func NewModel(structName string, connID int) Model {
	return Models[structName](connID)
}
