package factory

import (
	"sort"
	"strings"

	"github.com/admpub/errors"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/param"
)

// ModelInstancer 模型实例化
type ModelInstancer func(connID int) Model

// ModelInstancers 模型实例化
type ModelInstancers map[string]ModelInstancer

// FieldInfo 字段信息
type FieldInfo struct {
	//以下为数据库中的信息
	Name          string   `json:"name" xml:"name" bson:"name"`                            //字段名
	DataType      string   `json:"dataType" xml:"dataType" bson:"dataType"`                //数据库中数据类型
	Unsigned      bool     `json:"unsigned" xml:"unsigned" bson:"unsigned"`                //是否是无符号类型
	PrimaryKey    bool     `json:"primaryKey" xml:"primaryKey" bson:"primaryKey"`          //是否是主键
	AutoIncrement bool     `json:"autoIncrement" xml:"autoIncrement" bson:"autoIncrement"` //是否是自增字段
	Min           float64  `json:"min" xml:"min" bson:"min"`                               //最小值
	Max           float64  `json:"max" xml:"max" bson:"max"`                               //最大值
	Precision     int      `json:"precision" xml:"precision" bson:"precision"`             //小数精度(小数保留位数)
	MaxSize       int      `json:"maxSize" xml:"maxSize" bson:"maxSize"`                   //最大尺寸
	Options       []string `json:"options" xml:"options" bson:"options"`                   //选项值
	DefaultValue  string   `json:"defaultValue" xml:"defaultValue" bson:"defaultValue"`    //默认值
	Comment       string   `json:"comment" xml:"comment" bson:"comment"`                   //备注

	//以下为Golang中的信息
	GoType string `json:"goType" xml:"goType" bson:"goType"` //Golang数据类型
	GoName string `json:"goName" xml:"goName" bson:"goName"` //Golang字段名
}

// Validate 验证值是否符合数据库要求
func (f *FieldInfo) Validate(value interface{}) error {
	switch {
	case f.DataType == `enum`:
		if len(f.Options) == 0 {
			return nil
		}
		r := param.AsString(value)
		if !com.InSlice(r, f.Options) {
			return errors.Errorf(`The value "%v" does not exist in [%s]`, value, strings.Join(f.Options, `,`))
		}
	case f.DataType == `set`:
		if len(f.Options) == 0 {
			return nil
		}
		var values []string
		switch v := value.(type) {
		case []string:
			values = v
		default:
			values = param.Split(value, `,`)
		}
		for _, r := range values {
			if !com.InSlice(r, f.Options) {
				return errors.Errorf(`The value "%v" does not exist in [%s]`, value, strings.Join(f.Options, `,`))
			}
		}
	case strings.HasPrefix(f.DataType, `char`):
		r := param.AsString(value)
		if len(r) > f.MaxSize {
			return errors.Errorf(`Content "%v" cannot exceed %d characters`, value, f.MaxSize)
		}
	default:
		//case f.DataType == `decimal`:
		//case f.DataType == `double`:
		//case f.DataType == `float`:
		//case strings.HasPrefix(f.DataType, `int`):
		if f.Max <= 0 {
			if f.MaxSize > 0 {
				r := param.AsString(value)
				if len(r) > f.MaxSize {
					return errors.Errorf(`Content "%v" cannot exceed %d characters`, value, f.MaxSize)
				}
			}
			return nil
		}
		r := param.AsFloat64(value)
		if r < f.Min {
			return errors.Errorf(`The value "%v" cannot be less than %v`, value, f.Min)
		}
		if r > f.Max {
			return errors.Errorf(`The value "%v" cannot be greater than %v`, value, f.Max)
		}
	}
	return nil
}

// FieldValidator 字段验证器
type FieldValidator map[string]map[string]*FieldInfo

// ExistField 字段是否存在(表名不带前缀)
func (f FieldValidator) ExistField(table string, field string) bool {
	if tb, ok := f[table]; ok {
		_, ok = tb[field]
		return ok
	}
	return false
}

// Find 获取字段信息(表名不带前缀)
func (f FieldValidator) Find(table string, field string) (*FieldInfo, bool) {
	if tb, ok := f[table]; ok {
		fi, ok := tb[field]
		return fi, ok
	}
	return nil, false
}

// Validate 验证值是否符合数据库要求(表名不带前缀)
func (f FieldValidator) Validate(table string, field string, value interface{}) error {
	tb, ok := f[table]
	if !ok {
		return errors.WithMessage(ErrNotFoundTable, table)
	}
	fieldInfo, ok := tb[field]
	if !ok {
		return errors.WithMessage(ErrNotFoundField, field)
	}
	return fieldInfo.Validate(value)
}

// BatchValidate 批量验证值是否符合数据库要求(表名不带前缀)
func (f FieldValidator) BatchValidate(table string, row map[string]interface{}) error {
	tb, ok := f[table]
	if !ok {
		return errors.WithMessage(ErrNotFoundTable, table)
	}
	for field, value := range row {
		fieldInfo, ok := tb[field]
		if !ok {
			return errors.WithMessage(ErrNotFoundField, field)
		}
		err := fieldInfo.Validate(value)
		if err != nil {
			return errors.WithMessage(err, field)
		}
	}
	return nil
}

// ExistTable 表是否存在(表名不带前缀)
func (f FieldValidator) ExistTable(table string) bool {
	_, ok := f[table]
	return ok
}

// FieldList 获取表的字段列表(表名不带前缀)
func (f FieldValidator) FieldList(table string, excludeField ...string) []string {
	var fields []string
	if tb, ok := f[table]; ok {
		for field := range tb {
			if !com.InSlice(field, excludeField) {
				fields = append(fields, field)
			}
		}
	}
	return fields
}

// SortedFieldList 获取表的字段列表，并排序(表名不带前缀)
func (f FieldValidator) SortedFieldList(table string, excludeField ...string) []string {
	fields := f.FieldList(table, excludeField...)
	sort.Strings(fields)
	return fields
}

// SortedFieldLists 获取表的字段列表，并排序，返回[]interface{}(表名不带前缀)
func (f FieldValidator) SortedFieldLists(table string, excludeField ...string) []interface{} {
	fields := f.SortedFieldList(table, excludeField...)
	returns := make([]interface{}, len(fields))
	for i, v := range fields {
		returns[i] = v
	}
	return returns
}

// FieldLists 获取表的字段列表，返回[]interface{}(表名不带前缀)
func (f FieldValidator) FieldLists(table string, excludeField ...string) []interface{} {
	var fields []interface{}
	if tb, ok := f[table]; ok {
		for field := range tb {
			if !com.InSlice(field, excludeField) {
				fields = append(fields, field)
			}
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

// TableNamerRegister 注册表名称生成器(表名不带前缀)
func TableNamerRegister(namers map[string]func(obj interface{}) string) {
	for table, namer := range namers {
		TableNamers[table] = namer
	}
}

// TableNamerGet 获取表名称生成器(表名不带前缀)
func TableNamerGet(table string) func(obj interface{}) string {
	if namer, ok := TableNamers[table]; ok {
		return namer
	}
	return DefaultTableNamer(table)
}

// FieldRegister 注册字段信息(表名不带前缀)
func FieldRegister(tables map[string]map[string]*FieldInfo) {
	for table, info := range tables {
		Fields[table] = info
	}
}

// FieldFind 获取字段信息(表名不带前缀)
func FieldFind(table string, field string) (*FieldInfo, bool) {
	return Fields.Find(table, field)
}

// ModelRegister 模型构造函数登记
func ModelRegister(instancers map[string]ModelInstancer) {
	for structName, instancer := range instancers {
		Models[structName] = instancer
	}
}

// ExistField 字段是否存在(表名不带前缀)
func ExistField(table string, field string) bool {
	return Fields.ExistField(table, field)
}

// ExistTable 表是否存在(表名不带前缀)
func ExistTable(table string) bool {
	return Fields.ExistTable(table)
}

// NewModel 模型实例化
func NewModel(structName string, connID int) Model {
	return Models[structName](connID)
}

// Validate 验证值是否符合数据库要求
func Validate(table string, field string, value interface{}) error {
	return Fields.Validate(table, field, value)
}

// BatchValidate 批量验证值是否符合数据库要求
func BatchValidate(table string, row map[string]interface{}) error {
	return Fields.BatchValidate(table, row)
}
