package factory

import "github.com/webx-top/db"

var (
	databases = map[string]*DBI{
		DefaultDBKey: DefaultDBI,
	}
)

func DBIRegister(dbi *DBI, keys ...string) {
	key := DefaultDBKey
	if len(keys) > 0 {
		key = keys[0]
	}
	if _, y := databases[key]; y {
		panic(`DBI key already exists, please do not duplicate registrations`)
	}
	databases[key] = dbi
}

func DBIGet(keys ...string) *DBI {
	if len(keys) > 0 {
		return databases[keys[0]]
	}
	return databases[DefaultDBKey]
}

func DBIExists(key string) bool {
	_, ok := databases[key]
	return ok
}

func NewDBI() *DBI {
	return &DBI{
		Fields:      map[string]map[string]*FieldInfo{},
		Columns:     map[string][]string{},
		Models:      ModelInstancers{},
		TableNamers: map[string]func(obj interface{}) string{},
		Events:      NewEvents(),
	}
}

// DBI 数据库信息
type DBI struct {
	// Fields {table:{field:FieldInfo}}
	Fields FieldValidator
	// Columns {table:[field1,field2]}
	Columns map[string][]string
	// Models {StructName:ModelInstancer}
	Models ModelInstancers
	// TableNamers {table:NewName}
	TableNamers TableNamers
	Events      Events
}

func (d *DBI) TableName(structName string) string {
	m, ok := d.Models[structName]
	if ok {
		return m.Short
	}
	return ``
}

func (d *DBI) TableComment(structName string) string {
	m, ok := d.Models[structName]
	if ok {
		return m.Comment
	}
	return ``
}

func (d *DBI) TableColumns(tableName string) []string {
	cols, ok := d.Columns[tableName]
	if ok {
		return cols
	}
	return nil
}

func (d *DBI) EventFire(event string, model Model, mw func(db.Result) db.Result, args ...interface{}) error {
	return d.Events.Call(event, model, mw, args...)
}

func (d *DBI) EventOn(event string, h EventHandler, table string) {
	d.Events.Add(event, h, table)
}
