package mysql

import (
	"encoding/csv"

	"github.com/webx-top/echo"
)

func newCSVDocExportor(dbName string) DocExportor {
	return &mysqlExportCSVDoc{
		dbName: dbName,
	}
}

type mysqlExportCSVDoc struct {
	dbName string
	writer *csv.Writer
}

func (a *mysqlExportCSVDoc) Open(c echo.Context) error {
	a.writer = csv.NewWriter(c.Response())
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEEventStream)
	encodedName := echo.URLEncode(a.dbName+`_doc.csv`, true)
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename="+encodedName+"; filename*=utf-8''"+encodedName)
	c.Response().Write([]byte("\xEF\xBB\xBF")) // BOM
	return a.writer.Write([]string{`Table Documention`})
}

func (a *mysqlExportCSVDoc) Write(c echo.Context, table *TableStatus, fields []*Field) error {
	title := table.Name.String
	if len(table.Comment.String) > 0 {
		title += ` (` + table.Comment.String + `)`
	}
	err := a.writer.Write([]string{title})
	if err != nil {
		return err
	}
	err = a.writer.Write([]string{c.T(`字段名`), c.T(`数据类型`), c.T(`说明`)})
	if err != nil {
		return err
	}
	for _, v := range fields {
		dataType := v.Full_type
		if v.Null {
			dataType += ` NULL`
		}
		if v.AutoIncrement.Valid {
			dataType += ` ` + c.T("自动增量")
		}
		if v.Default.Valid {
			dataType += ` [` + v.Default.String + `]`
		}
		if len(v.On_update) > 0 {
			dataType += ` ON UPDATE ` + v.On_update
		}
		err = a.writer.Write([]string{v.Field, dataType, v.Comment})
		if err != nil {
			return err
		}
	}

	return a.writer.Write([]string{``})
}

func (a *mysqlExportCSVDoc) Close(c echo.Context) error {
	a.writer.Flush()
	return nil
}
