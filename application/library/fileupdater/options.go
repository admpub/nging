package fileupdater

type Options struct {
	TableName  string   // 数据表名称
	FieldName  string   // 数据表字段名
	SameFields []string // 数据表类似字段名
	Embedded   bool     // 是否为嵌入图片
	Seperator  string   // 文件字段中多个文件路径之间的分隔符，空字符串代表为单个文件
	Callback   CallbackFunc
}

type OptionSetter func(o *Options)

func TableName(tableName string) OptionSetter {
	return func(o *Options) {
		o.TableName = tableName
	}
}

func FieldName(fieldName string) OptionSetter {
	return func(o *Options) {
		o.FieldName = fieldName
	}
}

func SameFields(sameFields ...string) OptionSetter {
	return func(o *Options) {
		o.SameFields = sameFields
	}
}

func Embedded(embedded bool) OptionSetter {
	return func(o *Options) {
		o.Embedded = embedded
	}
}

func Seperator(seperator string) OptionSetter {
	return func(o *Options) {
		o.Seperator = seperator
	}
}

func Callback(callbackFunc CallbackFunc) OptionSetter {
	return func(o *Options) {
		o.Callback = callbackFunc
	}
}
