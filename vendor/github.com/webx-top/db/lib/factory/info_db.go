package factory

import (
	"fmt"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
)

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
		TableNamers: map[string]func(Model) string{},
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

func (d *DBI) NewModel(structName string, connID ...int) Model {
	mi, ok := d.Models[structName]
	if !ok {
		return nil
	}
	var _connID int
	if len(connID) > 0 {
		_connID = connID[0]
	}
	return mi.Make(_connID)
}

func (d *DBI) ModelInfo(structName string) *ModelInstancer {
	mi, _ := d.Models[structName]
	return mi
}

type Short interface {
	Short_() string
}

func NoPrefixTableName(v interface{}) string {
	var noPrefixTableName string
	switch a := v.(type) {
	case string:
		noPrefixTableName = a
	case Short:
		noPrefixTableName = a.Short_()
	default:
		panic(fmt.Sprintf(`Unsupported type: %T`, v))
	}
	return noPrefixTableName
}

func (d *DBI) OmitSelect(v interface{}, excludeColumns ...string) []interface{} {
	noPrefixTableName := NoPrefixTableName(v)
	columns := d.TableColumns(noPrefixTableName)
	results := make([]interface{}, 0, len(columns)-len(excludeColumns))
	for _, column := range columns {
		if com.InSlice(column, excludeColumns) {
			continue
		}
		results = append(results, column)
	}
	return results
}

func (d *DBI) OmitColumns(v interface{}, excludeColumns ...string) []string {
	noPrefixTableName := NoPrefixTableName(v)
	columns := d.TableColumns(noPrefixTableName)
	results := make([]string, 0, len(columns)-len(excludeColumns))
	for _, column := range columns {
		if com.InSlice(column, excludeColumns) {
			continue
		}
		results = append(results, column)
	}
	return results
}

func (d *DBI) TableColumns(tableName string) []string {
	cols, ok := d.Columns[tableName]
	if ok {
		return cols
	}
	return nil
}

func (d *DBI) HasEvent(event string, model Model) bool {
	return d.Events.Exists(event, model)
}

func (d *DBI) Fire(event string, model Model, mw func(db.Result) db.Result, args ...interface{}) error {
	return d.Events.Call(event, model, nil, mw, args...)
}

func (d *DBI) FireUpdate(event string, model Model, editColumns []string, mw func(db.Result) db.Result, args ...interface{}) error {
	return d.Events.Call(event, model, editColumns, mw, args...)
}

func (d *DBI) FireReading(model Model, param *Param, rangers ...Ranger) error {
	return d.Events.CallRead(EventReading, model, param, rangers...)
}

func (d *DBI) FireReaded(model Model, param *Param, rangers ...Ranger) error {
	return d.Events.CallRead(EventReaded, model, param, rangers...)
}

func (d *DBI) ParseEventNames(event string) []string {
	switch event {
	case `w+`:
		return AllAfterWriteEvents
	case `+w`:
		return AllBeforeWriteEvents
	case `r+`:
		return AllAfterReadEvents
	case `+r`:
		return AllBeforeReadEvents
	default:
		return strings.Split(event, ",")
	}
}

// - 注册写(CUD)事件

func (d *DBI) On(event string, h EventHandler, tableName ...string) *DBI {
	var table string
	if len(tableName) > 0 {
		table = tableName[0]
	} else {
		set := strings.SplitN(event, ":", 2)
		switch len(set) {
		case 2:
			event = set[1]
			fallthrough
		case 1:
			table = set[0]
		}
	}
	if !d.Fields.ExistTable(table) {
		panic(`Table does not exist: ` + table)
	}
	for _, evt := range d.ParseEventNames(event) {
		d.Events.On(evt, h, table)
	}
	return d
}

func (d *DBI) OnAsync(event string, h EventHandler, tableName ...string) *DBI {
	var table string
	if len(tableName) > 0 {
		table = tableName[0]
	} else {
		set := strings.SplitN(event, ":", 2)
		switch len(set) {
		case 2:
			event = set[1]
			fallthrough
		case 1:
			table = set[0]
		}
	}
	if !d.Fields.ExistTable(table) {
		panic(`Table does not exist: ` + table)
	}
	for _, evt := range d.ParseEventNames(event) {
		d.Events.On(evt, h, table, true)
	}
	return d
}

// - 注册读(R)事件

func (d *DBI) OnRead(event string, h EventReadHandler, tableName ...string) *DBI {
	var table string
	if len(tableName) > 0 {
		table = tableName[0]
	} else {
		set := strings.SplitN(event, ":", 2)
		switch len(set) {
		case 2:
			event = set[1]
			fallthrough
		case 1:
			table = set[0]
		}
	}
	if !d.Fields.ExistTable(table) {
		panic(`Table does not exist: ` + table)
	}
	for _, evt := range d.ParseEventNames(event) {
		d.Events.OnRead(evt, h, table)
	}
	return d
}

func (d *DBI) OnReadAsync(event string, h EventReadHandler, tableName ...string) *DBI {
	var table string
	if len(tableName) > 0 {
		table = tableName[0]
	} else {
		set := strings.SplitN(event, ":", 2)
		switch len(set) {
		case 2:
			event = set[1]
			fallthrough
		case 1:
			table = set[0]
		}
	}
	if !d.Fields.ExistTable(table) {
		panic(`Table does not exist: ` + table)
	}
	for _, evt := range d.ParseEventNames(event) {
		d.Events.OnRead(evt, h, table, true)
	}
	return d
}

// FieldsRegister 注册字段信息(表名不带前缀)
func (d *DBI) FieldsRegister(tables map[string]map[string]*FieldInfo) {
	for table, info := range tables {
		d.Fields[table] = info
	}
}

// ColumnsRegister 注册模型构造函数(map的不带前缀表名)
func (d *DBI) ColumnsRegister(columns map[string][]string) {
	for table, cols := range columns {
		d.Columns[table] = cols
	}
}

// ModelsRegister 注册模型构造函数(map的键为结构体名)
func (d *DBI) ModelsRegister(instancers map[string]*ModelInstancer) {
	d.Models.Register(instancers)
}

// TableNamersRegister 自定义表名称生成函数
func (d *DBI) TableNamersRegister(namers map[string]func(Model) string) {
	d.TableNamers.Register(namers)
}
