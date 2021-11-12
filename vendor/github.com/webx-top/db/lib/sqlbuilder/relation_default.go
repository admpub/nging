package sqlbuilder

import (
	"fmt"
	"reflect"

	"github.com/admpub/errors"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/reflectx"
)

var (
	ErrUnableDetermineTableName                       = errors.New(`Unable to determine table name`)
	ErrModelCannotNil                                 = errors.New("model argument cannot be nil pointer passed")
	GetTableName                TableNameFunc         = DefaultGetTableName
	StructToTableName           StructToTableNameFunc = DefaultStructToTableName
	GetDBConn                   DBConnFunc            = DefaultDBConn
	GetSQLBuilder               SQLBuilderFunc        = DefaultGetSQLBuilder
	DBConnTagName               string                = `dbconn`
)

func DefaultDBConn(name string) db.Database {
	return nil
}

func DefaultGetSQLBuilder(fieldInfo *reflectx.FieldInfo, defaults ...SQLBuilder) SQLBuilder {
	if len(DBConnTagName) > 0 {
		if dbConnName, ok := fieldInfo.Options[DBConnTagName]; ok && len(dbConnName) > 0 {
			if database := GetDBConn(dbConnName); database != nil {
				if sqlBuilder, ok := database.(SQLBuilder); ok {
					return sqlBuilder
				}
			}
		}
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return nil
}

func DefaultGetTableName(fieldInfo *reflectx.FieldInfo, data interface{}) (table string, err error) {
	var ok bool
	table, ok = fieldInfo.Options[`table`] // table=table1
	if !ok || len(table) == 0 {
		table, err = StructToTableName(data)
	}
	return
}

func DefaultStructToTableName(data interface{}, retry ...bool) (string, error) {
	switch m := data.(type) {
	case Name_:
		return m.Name_(), nil
	case db.TableName:
		return m.TableName(), nil
	default:
		if len(retry) > 0 && retry[0] {
			return ``, ErrUnableDetermineTableName
		}
	}
	value := reflect.ValueOf(data)
	if value.IsNil() {
		return ``, errors.WithMessagef(ErrModelCannotNil, `%T`, data)
	}
	tp := reflect.Indirect(value).Type()
	if tp.Kind() == reflect.Interface {
		tp = reflect.Indirect(value).Elem().Type()
	}

	if tp.Kind() != reflect.Slice {
		return ``, fmt.Errorf("model argument must slice, but get %T", data)
	}

	tpEl := tp.Elem()
	//Compatible with []*Struct or []Struct
	if tpEl.Kind() == reflect.Ptr {
		tpEl = tpEl.Elem()
	}
	name, err := DefaultStructToTableName(reflect.New(tpEl).Interface(), true)
	if err == ErrUnableDetermineTableName {
		name = com.SnakeCase(tpEl.Name())
		err = nil
	}
	return name, err
}
