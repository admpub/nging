package mysql

import (
	"strings"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/webx-top/echo"
)

func newHTMLDocExportor(dbName string) DocExportor {
	return &mysqlExportHTMLDoc{
		dbName: dbName,
	}
}

type mysqlExportHTMLDoc struct {
	dbName string
}

func (a *mysqlExportHTMLDoc) Open(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEEventStream)
	encodedName := echo.URLEncode(a.dbName+`_doc.html`, true)
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename="+encodedName+"; filename*=utf-8''"+encodedName)
	c.Response().Write([]byte(`<doctype html><html><head><title>Table Documention</title><meta name="viewport" content="width=device-width, initial-scale=1"><link rel="stylesheet" href="` + common.BackendURL(c) + `/public/assets/backend/js/bootstrap/dist/css/bootstrap.min.css?t=20220807115920" /></head><body><div class="container"><div class="row"><div class="col-md-12">`))
	return nil
}

func (a *mysqlExportHTMLDoc) Write(c echo.Context, table *TableStatus, fields []*Field) error {
	c.Response().Write([]byte(`<h2 id="` + strings.ReplaceAll(table.Name.String, `"`, ``) + `">` + table.Name.String + `</h2>`))
	c.Response().Write([]byte(`<em>` + table.Comment.String + `</em>`))
	c.Response().Write([]byte(`<table class="table table-bordered table-hover table-condensed">`))
	c.Response().Write([]byte(`<thead><tr><th>` + c.T(`字段名`) + `</th><th>` + c.T(`数据类型`) + `</th><th>` + c.T(`默认值`) + `</th><th>` + c.T(`必要字段`) + `</th><th>` + c.T(`说明`) + `</th></tr></thead>`))
	c.Response().Write([]byte(`<tbody>`))
	for _, v := range fields {
		dataType := v.Full_type
		if v.Null {
			dataType += ` <span>NULL</span>`
		}
		if v.AutoIncrement.Valid {
			dataType += ` <em>` + c.T("自动增量") + `</em>`
		}
		if len(v.On_update) > 0 {
			dataType += ` ON UPDATE <b>` + v.On_update + `</b>`
		}
		var defaultValue string
		if v.Default.Valid {
			defaultValue = ` [<b>` + v.Default.String + `</b>]`
		}
		required := `<b>` + c.T(`是`) + `</b>`
		if !v.IsRequired() {
			required = c.T(`否`)
		}
		c.Response().Write([]byte(`<tr><td>` + v.Field + `</td><td>` + dataType + `</td><td>` + defaultValue + `</td><td>` + required + `</td><td>` + v.Comment + `</td></tr>`))
	}
	c.Response().Write([]byte(`</tbody>`))
	c.Response().Write([]byte(`</table>`))
	return nil
}

func (a *mysqlExportHTMLDoc) Close(c echo.Context) error {
	c.Response().Write([]byte(`</div></div></div></body></html>`))
	return nil
}
