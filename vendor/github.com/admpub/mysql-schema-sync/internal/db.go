package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

type DBOperator interface {
	DBEngine() string
	GetTableNames() ([]string, error)
	GetTableSchema(name string) (schema string, err error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string) (sql.Result, error)
	Begin() (*sql.Tx, error)
	Close() error
}

type Comparer interface {
	AlterData(sc *SchemaSync, tableName string) (*TableAlterData, error)
}

var _ DBOperator = new(MyDb)

// MyDb db struct
type MyDb struct {
	*sql.DB
	dbType string
	engine string
}

// NewMyDb parse dsn
func NewMyDb(dsn string, dbType string) *MyDb {
	if len(dsn) == 0 {
		log.Fatal(dbType + " dns is empty")
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Sprintf("connect to db [%s] failed, %v", dsn, err))
	}
	return &MyDb{
		DB:     db,
		dbType: dbType,
		engine: `mysql`,
	}
}

func (db *MyDb) DBEngine() string {
	return db.engine
}

// GetTableNames table names
func (mydb *MyDb) GetTableNames() ([]string, error) {
	rs, err := mydb.Query("show table status")
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
		if valObj["Engine"] != nil {
			tables = append(tables, valObj["Name"].(string))
		}
	}
	return tables, nil
}

// GetTableSchema table schema
func (mydb *MyDb) GetTableSchema(name string) (schema string, err error) {
	var rs *sql.Rows
	rs, err = mydb.Query(fmt.Sprintf("show create table `%s`", name))
	if err != nil {
		var mySQLErr *mysql.MySQLError
		if errors.As(err, &mySQLErr) && mySQLErr.Number == 1146 { // Error 1146 (42S02): Table 'table_name' doesn't exist
			err = nil
		}
		return
	}
	defer rs.Close()
	for rs.Next() {
		var vname string
		if err = rs.Scan(&vname, &schema); err != nil {
			err = fmt.Errorf("get table %s 's schema failed, %s", name, err)
			return
		}
	}
	return
}

// Query execute sql query
func (mydb *MyDb) Query(query string, args ...interface{}) (*sql.Rows, error) {
	log.Println("[SQL]", "["+mydb.dbType+"]", query, args)
	return mydb.DB.Query(query, args...)
}

// Exec execute sql query
func (mydb *MyDb) Exec(query string) (sql.Result, error) {
	log.Println("[SQL]", "["+mydb.dbType+"]", query)
	return Exec(mydb.DB, query)
}
