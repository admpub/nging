package mysql

import (
	"github.com/webx-top/echo"
)

func newMarkdownDocExportor(dbName string) DocExportor {
	return &mysqlExportMarkdownDoc{
		dbName: dbName,
	}
}

type mysqlExportMarkdownDoc struct {
	dbName string
}

func (a *mysqlExportMarkdownDoc) Open(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEEventStream)
	encodedName := echo.URLEncode(a.dbName+`_doc.md`, true)
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename="+encodedName+"; filename*=utf-8''"+encodedName)
	c.Response().Write([]byte(`# Table Documention` + "\n"))
	return nil
}

func (a *mysqlExportMarkdownDoc) Write(c echo.Context, table *TableStatus, fields []*Field) error {
	c.Response().Write([]byte(`## ` + table.Name.String + "\n\n"))
	c.Response().Write([]byte(`> ` + table.Comment.String + "\n\n"))
	c.Response().Write([]byte(`| ` + c.T(`字段名`) + ` | ` + c.T(`数据类型`) + ` | ` + c.T(`默认值`) + ` | ` + c.T(`是否必填`) + ` | ` + c.T(`说明`) + ` |` + "\n"))
	c.Response().Write([]byte(`| :------------ | :------------ |  :------------ |  :------------ | :------------ |` + "\n"))
	for _, v := range fields {
		dataType := v.Full_type
		if v.Null {
			dataType += ` NULL`
		}
		if v.AutoIncrement.Valid {
			dataType += ` *` + c.T("自动增量") + `*`
		}
		if v.Default.Valid {
			if len(v.Default.String) > 0 {
				dataType += ` [**` + v.Default.String + `**]`
			} else {
				dataType += ` []`
			}
		}
		if len(v.On_update) > 0 {
			dataType += ` ON UPDATE **` + v.On_update + `**`
		}
		required := c.T(`是`)
		if v.Null || v.Default.Valid {
			required = c.T(`否`)
		}
		c.Response().Write([]byte(`| ` + v.Field + ` | ` + dataType + ` | ` + v.Default.String + ` | ` + required + ` | ` + v.Comment + ` |` + "\n"))
	}
	c.Response().Write([]byte("\n"))
	return nil
}

func (a *mysqlExportMarkdownDoc) Close(c echo.Context) error {
	c.Response().Write([]byte("\n\n"))
	return nil
}
