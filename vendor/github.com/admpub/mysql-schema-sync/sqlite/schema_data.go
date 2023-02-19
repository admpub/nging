package sqlite

import (
	"database/sql"
	"log"
	"regexp"
)

var (
	sqlTableName = regexp.MustCompile("CREATE TABLE [^`\"]*[\"`]([^`\"]+)[\"`] \\(")
)

// SchemaData db schema data
type SchemaData struct {
	Data   string
	dbType string
	engine string
}

// NewSchemaData object
func NewSchemaData(schema string, dbType string) *SchemaData {
	return &SchemaData{
		Data:   schema,
		dbType: dbType,
		engine: `sqlite`,
	}
}

func (m *SchemaData) DBEngine() string {
	return m.engine
}

// GetTableNames table names
func (m *SchemaData) GetTableNames() ([]string, error) {
	matches := sqlTableName.FindAllStringSubmatch(m.Data, -1)
	var tables []string
	for _, match := range matches {
		tables = append(tables, match[1])
	}
	return tables, nil
}

// GetTableSchema table schema
func (m *SchemaData) GetTableSchema(name string) (schema string, err error) {
	var schemaStruct *regexp.Regexp
	schemaStruct, err = regexp.Compile("(?sm)CREATE TABLE [^`\"]*[\"`]" + name + "[\"`] \\((.+?)\\)[;]?(?:[\\r]?\\n|$)")
	if err != nil {
		return
	}
	matches := schemaStruct.FindStringSubmatch(m.Data)
	//log.Printf("%#v\n", matches)
	if matches != nil {
		schema = matches[0]
	}
	if len(schema) > 0 {
		schema = FormatSchema(schema)
	}
	var schemaIndex *regexp.Regexp
	schemaIndex, err = regexp.Compile("(?sm)CREATE (?:UNIQUE )?INDEX (?:IF NOT EXISTS )?[`\"][^`\"]*[`\"] ON [`\"]" + name + "[`\"]([^\\r\\n]*)[\\r]?\\n")
	if err != nil {
		return
	}
	matches2 := schemaIndex.FindAllStringSubmatch(m.Data, -1)
	//log.Printf("%#v\n", matches2)
	for _, matches := range matches2 {
		schema += matches[0]
	}
	return
}

// Query execute sql query
func (m *SchemaData) Query(query string, args ...interface{}) (*sql.Rows, error) {
	log.Println("[SQL]", "["+m.dbType+"]", query, args)
	return nil, nil
}

// Exec execute sql query
func (m *SchemaData) Exec(query string) (sql.Result, error) {
	log.Println("[SQL]", "["+m.dbType+"]", query)
	return nil, nil
}

func (m *SchemaData) Begin() (*sql.Tx, error) {
	return nil, nil
}

func (m *SchemaData) Close() error {
	return nil
}
