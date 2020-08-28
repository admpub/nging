package fileupdater

import (
	"github.com/webx-top/com"
	"github.com/webx-top/db"
)

type Options struct {
	TableName  string       // 数据表名称
	FieldName  string       // 数据表字段名
	SameFields []string     // 数据表类似字段名
	Embedded   bool         // 是否为嵌入图片
	Seperator  string       // 文件字段中多个文件路径之间的分隔符，空字符串代表为单个文件
	Callback   CallbackFunc `json:"-" xml:"-"`
	FieldValue FieldValue   `json:"-" xml:"-"`
}

type OptionSetter func(o *Options)

func OptTableName(tableName string) OptionSetter {
	return func(o *Options) {
		o.TableName = tableName
	}
}

func OptFieldName(fieldName string) OptionSetter {
	return func(o *Options) {
		o.FieldName = fieldName
	}
}

func OptSameFields(sameFields ...string) OptionSetter {
	return func(o *Options) {
		o.SameFields = sameFields
	}
}

func OptEmbedded(embedded bool) OptionSetter {
	return func(o *Options) {
		o.Embedded = embedded
	}
}

func OptSeperator(seperator string) OptionSetter {
	return func(o *Options) {
		o.Seperator = seperator
	}
}

func OptCallback(callbackFunc CallbackFunc) OptionSetter {
	return func(o *Options) {
		o.Callback = callbackFunc
	}
}

func OptGenCallback(cond db.Compound, fieldValues ...FieldValue) OptionSetter {
	return OptCallback(GenCallbackWithCond(cond, fieldValues...))
}

func OptFieldValue(field string, value ValueFunc) OptionSetter {
	return func(o *Options) {
		if o.FieldValue == nil {
			o.FieldValue = FieldValueWith(field, value)
		} else {
			o.FieldValue.Set(field, value)
		}
		if !com.InSlice(field, o.SameFields) {
			o.SameFields = append(o.SameFields, field)
		}
	}
}
