package factory

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

func (f FieldValidator) ValidField(table string, field string) bool {
	if tb, ok := f[table]; ok {
		_, ok = tb[field]
		return ok
	}
	return false
}

func (f FieldValidator) ValidTable(table string) bool {
	_, ok := f[table]
	return ok
}

var Fields FieldValidator = map[string]map[string]*FieldInfo{}

func ValidField(table string, field string) bool {
	return Fields.ValidField(table, field)
}

func ValidTable(table string) bool {
	return Fields.ValidTable(table)
}
