package sqlite

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/admpub/mysql-schema-sync/internal"
)

var (
	_ internal.DBOperator = &MyDb{}
	_ internal.DBOperator = &SchemaData{}
)

// MyDb db struct
type MyDb struct {
	*sql.DB
	dbType string
	engine string
}

// New parse dsn
func New(dsn string, dbType string) *MyDb {
	if len(dsn) == 0 {
		log.Fatal(dbType + " dns is empty")
	}
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		panic(fmt.Sprintf("connect to db [%s] failed, %v", dsn, err))
	}
	return &MyDb{
		DB:     db,
		dbType: dbType,
		engine: `sqlite`,
	}
}

func (mydb *MyDb) DBEngine() string {
	return mydb.engine
}

// GetTableNames table names
func (mydb *MyDb) GetTableNames() ([]string, error) {
	rs, err := mydb.Query("SELECT tbl_name FROM sqlite_master WHERE type='table'")
	if err != nil {
		return nil, fmt.Errorf("show tables failed: %s", err.Error())
	}
	defer rs.Close()
	tables := []string{}
	columns, err := rs.Columns()
	if err != nil {
		return nil, err
	}
	for rs.Next() {
		var values = make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		if err := rs.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("show tables failed when scan, %s", err.Error())
		}
		var valObj = make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			valObj[col] = v
		}
		tables = append(tables, valObj["tbl_name"].(string))
	}
	return tables, nil
}

// GetTableSchema table schema
func (mydb *MyDb) GetTableSchema(name string) (schema string, err error) {
	var rs *sql.Rows
	rs, err = mydb.Query(fmt.Sprintf("SELECT sql FROM sqlite_master WHERE type='table' and tbl_name='%s'", name))
	if err != nil {
		return
	}
	for rs.Next() {
		if err = rs.Scan(&schema); err != nil {
			rs.Close()
			err = fmt.Errorf("get table %s 's schema failed,%s", name, err)
			return
		}
	}
	rs.Close()
	if len(schema) == 0 {
		return
	}
	schema = FormatSchema(schema)
	schema += ";\n"
	rs, err = mydb.Query(fmt.Sprintf("SELECT sql FROM sqlite_master WHERE type='index' and tbl_name='%s'", name))
	if err != nil {
		return
	}
	defer rs.Close()
	for rs.Next() {
		var schemaIndex sql.NullString
		if err = rs.Scan(&schemaIndex); err != nil {
			err = fmt.Errorf("get table %s 's schema failed,%s", name, err)
			return
		}
		if len(schemaIndex.String) > 0 {
			indexString, _ := parseIndex(schemaIndex.String)
			schema += indexString + ";\n"
		}
	}
	return
}

// Query execute sql query
func (mydb *MyDb) Query(query string, args ...interface{}) (*sql.Rows, error) {
	log.Println("[SQL]", "["+mydb.dbType+"]", query, args)
	return mydb.DB.Query(query)
}

// Exec execute sql query
func (mydb *MyDb) Exec(query string) (sql.Result, error) {
	log.Println("[SQL]", "["+mydb.dbType+"]", query)
	return internal.Exec(mydb.DB, query)
}
