package mysql

import (
	"sort"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo/param"
)

// BuildSelectCase : 生成CASE...ELSE...END
// SELECT status,
// CASE status
// WHEN 1 THEN '待查询'
// WHEN 2 THEN '审核中'
// WHEN 3 THEN '审核通过'
// WHEN 4 THEN '审核拒绝'
// ELSE 'Unknown'
// END statusName,
// count(status)
// FROM order GROUP BY status;
func BuildSelectCase(field string, dict map[string]string, defaultValue string, asFields ...string) string {
	var asField string
	if len(asFields)> 0{
		asField = asFields[0]
	}
	values := make([]string, len(dict))
	var index int
	for value := range dict {
		values[index] = value
		index++
	}
	sort.Strings(values)
	sqlString := "CASE `"+field+"`\n"
	for _, value := range values {
		text := dict[value]
		sqlString += "WHEN '"+com.AddSlashes(value)+"' THEN '"+com.AddSlashes(text)+"'\n"
	}
	sqlString += "ELSE '"+com.AddSlashes(defaultValue)+"'\n"
	sqlString += "END"
	if len(asField)>0 {
		sqlString += " `"+asField+"`"
	}
	return sqlString
}

// BuildBatchUpdate 生成批量更新SQL
func BuildBatchUpdate(table string, rows []param.Store, whereFields ...string) string {
	if len(rows) < 1 {
		return ``
	}
	var set string
	firstRow := rows[0]
	fields := make([]string, len(firstRow))
	var index int
	for field := range firstRow {
		fields[index] = field
		index++
	}
	sort.Strings(fields)
	joins := make([]string, len(rows))
	for index, row := range rows {
		join := `SELECT `;
		for _, field := range fields {
			join += "'"+com.AddSlashes(row.String(field))+"' AS `"+field+"`,"
		}
		join = strings.TrimSuffix(join, `,`)
		joins[index] = join
	}
	using := "USING(`" + strings.Join(whereFields, "`, `") + "`)"
	for _, whereField := range whereFields {
		delete(firstRow, whereField)
	}
	for _, field := range fields {
		set += "`a`.`"+field+"`=`b`.`"+field+"`,";
	}
	set = strings.TrimSuffix(set, `,`)
	sqlString := "UPDATE `"+table+"` a"
	sqlString += "\nJOIN ("
	sqlString += "\n " + strings.Join(joins, "\n UNION\n ")
	sqlString += "\n) b"
	sqlString += "\n" + using
	sqlString += "\nSET " + set
	return sqlString
}

// BuildBatchInsert 生成批量插入SQL
func BuildBatchInsert(table string, rows []param.Store, force ...bool) string {
	if len(rows) < 1 {
		return ``
	}
	firstRow := rows[0]
	fields := make([]string, len(firstRow))
	var index int
	for field := range firstRow {
		fields[index] = field
		index++
	}
	sort.Strings(fields)
	setFields := "(`" + strings.Join(fields, "`, `") + "`)"
	setValues := make([]string, len(rows))
	for index, row := range rows {
		setValue := `(`;
		for _, field := range fields {
			setValue += "'"+com.AddSlashes(row.String(field))+"',"
		}
		setValue = strings.TrimSuffix(setValue, `,`)
		setValue += `)`
		setValues[index] = setValue
	}
	sqlString := `INSERT IGNORE INTO`
	if len(force) > 0 && force[0] {
		sqlString = `REPLACE INTO`
	}
	sqlString += " `"+table+"` " + setFields + " VALUES"
	sqlString += "\n" + strings.Join(setValues, ",\n")
	return sqlString
}
