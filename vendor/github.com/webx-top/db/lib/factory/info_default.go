package factory

// DefaultDBKey 默认数据库标识
const DefaultDBKey = `default`

var (
	// DefaultTableNamer default
	DefaultTableNamer = func(table string) func(Model) string {
		return func(_ Model) string {
			return table
		}
	}
	// DefaultDBI default
	DefaultDBI = NewDBI()
)

// TableNamerRegister 注册表名称生成器(表名不带前缀)
func TableNamerRegister(namers map[string]func(Model) string) {
	DBIGet().TableNamers.Register(namers)
}

// TableNamerGet 获取表名称生成器(表名不带前缀)
func TableNamerGet(table string) func(Model) string {
	return DBIGet().TableNamers.Get(table)
}

// FieldRegister 注册字段信息(表名不带前缀)
func FieldRegister(tables map[string]map[string]*FieldInfo) {
	DBIGet().Fields.Register(tables)
}

// FieldFind 获取字段信息(表名不带前缀)
func FieldFind(table string, field string) (*FieldInfo, bool) {
	return DBIGet().Fields.Find(table, field)
}

// ModelRegister 模型构造函数登记
func ModelRegister(instancers map[string]*ModelInstancer) {
	DBIGet().Models.Register(instancers)
}

// ExistField 字段是否存在(表名不带前缀)
func ExistField(table string, field string) bool {
	return DBIGet().Fields.ExistField(table, field)
}

// ExistTable 表是否存在(表名不带前缀)
func ExistTable(table string) bool {
	return DBIGet().Fields.ExistTable(table)
}

// NewModel 模型实例化
func NewModel(structName string, connID int) Model {
	return DBIGet().Models[structName].Make(connID)
}

// Validate 验证值是否符合数据库要求
func Validate(table string, field string, value interface{}) error {
	return DBIGet().Fields.Validate(table, field, value)
}

// BatchValidate 批量验证值是否符合数据库要求
func BatchValidate(table string, row map[string]interface{}) error {
	return DBIGet().Fields.BatchValidate(table, row)
}
