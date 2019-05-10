package sync

import (
	"database/sql"
	"log"
	"regexp"

	"github.com/admpub/mysql-schema-sync/internal"
)

var (
	sqlTableName                     = regexp.MustCompile("CREATE TABLE [^`]*`([^`]+)` \\(")
	_            internal.DBOperator = &MySchemaData{}
)

// MySchemaData db schema data
type MySchemaData struct {
	Data   string
	dbType string
}

// NewMySchemaData object
func NewMySchemaData(schema string, dbType string) *MySchemaData {
	return &MySchemaData{
		Data:   schema,
		dbType: dbType,
	}
}

// GetTableNames table names
func (m *MySchemaData) GetTableNames() []string {
	matches := sqlTableName.FindAllStringSubmatch(m.Data, -1)
	var tables []string
	if matches != nil {
		for _, match := range matches {
			tables = append(tables, match[1])
		}
	}
	return tables
}

// GetTableSchema table schema
func (m *MySchemaData) GetTableSchema(name string) (schema string) {
	schemaStruct, err := regexp.Compile("(?sm)CREATE TABLE [^`]*`" + name + "` \\((.+?)\\) ENGINE\\=[^\\r\\n]*;[\\r]?\\n")
	if err != nil {
		log.Println(err)
	}
	matches := schemaStruct.FindStringSubmatch(m.Data)
	//log.Printf("%#v\n", matches)
	if matches != nil {
		schema = matches[0]
	}
	return
}

// Query execute sql query
func (m *MySchemaData) Query(query string, args ...interface{}) (*sql.Rows, error) {
	log.Println("[SQL]", "["+m.dbType+"]", query, args)
	return nil, nil
}

func (m *MySchemaData) Begin() (*sql.Tx, error) {
	return nil, nil
}
