package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo/param"
)

const (
	// SQLTableComment 查询表注释的SQL
	SQLTableComment = `SELECT TABLE_NAME,TABLE_COMMENT FROM information_schema.TABLES WHERE table_schema=? AND TABLE_NAME=?`
	// SQLColumnComment 查询列注释的SQL
	SQLColumnComment = "SELECT COLUMN_NAME as `field`, column_comment as `description`, DATA_TYPE as `type`, CHARACTER_MAXIMUM_LENGTH as `max_length`, CHARACTER_OCTET_LENGTH as `octet_length`, NUMERIC_PRECISION as `precision` FROM INFORMATION_SCHEMA.COLUMNS WHERE table_schema=? AND table_name=?"
	// SQLShowCreate 查询建表语句的SQL
	SQLShowCreate = "SHOW CREATE TABLE "
	// SQLTableExists 查询表是否存在的SQL
	SQLTableExists = "SELECT COUNT(1) AS count FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA=? AND TABLE_NAME=?"
	// SQLFieldExists 查询字段是否存在的SQL
	SQLFieldExists = "SELECT COUNT(1) AS count FROM information_schema.columns WHERE table_name=? AND column_name=?"
)

// CreateTableSQL 查询建表语句
func CreateTableSQL(linkID int, dbName string, tableName string) (string, error) {
	ctx := context.Background()
	db := factory.NewParam().SetIndex(linkID).DB()
	sqlStr := SQLShowCreate + "`" + dbName + "`.`" + tableName + "`"
	stmt, err := db.PrepareContext(ctx, sqlStr)
	if factory.Debug() {
		fmt.Println(sqlStr)
	}
	if err != nil {
		return ``, fmt.Errorf(`CreateTableSQL: %v`, err)
	}
	recvTableName := sql.NullString{}
	recvCreateTableSQL := sql.NullString{}
	err = stmt.QueryRowContext(ctx).Scan(&recvTableName, &recvCreateTableSQL)
	if err != nil {
		return ``, fmt.Errorf(`CreateTableSQL.Scan: %v`, err)
	}
	return recvCreateTableSQL.String, err
}

// GetTables 获取数据表列表
func GetTables(linkID int) ([]string, error) {
	rows, err := factory.NewParam().SetIndex(linkID).SetCollection(`SHOW TABLES`).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := []string{}
	for rows.Next() {
		var v sql.NullString
		err := rows.Scan(&v)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v.String)
	}
	return ret, nil
}

// TableComment 查询表注释
func TableComment(linkID int, dbName string, tableName string) (string, error) {
	ctx := context.Background()
	db := factory.NewParam().SetIndex(linkID).DB()
	stmt, err := db.PrepareContext(ctx, SQLTableComment)
	if factory.Debug() {
		fmt.Println(SQLTableComment, `[`, dbName, tableName, `]`)
	}
	if err != nil {
		return ``, fmt.Errorf(`TableComment: %v`, err)
	}
	recvTableName := sql.NullString{}
	recvTableComment := sql.NullString{}
	err = stmt.QueryRowContext(ctx, dbName, tableName).Scan(&recvTableName, &recvTableComment)
	if err != nil {
		return ``, fmt.Errorf(`TableComment.Scan: %v`, err)
	}
	return recvTableComment.String, err
}

// ColumnComment 查询表中某些列的注释
func ColumnComment(linkID int, dbName string, tableName string, fieldNames ...string) (map[string]param.StringMap, []string, error) {
	ctx := context.Background()
	db := factory.NewParam().SetIndex(linkID).DB()
	sqlStr := SQLColumnComment
	if len(fieldNames) > 0 {
		if len(fieldNames) == 1 {
			sqlStr += `AND COLUMN_NAME = '` + com.AddSlashes(fieldNames[0]) + `'`
		} else {
			for key, val := range fieldNames {
				fieldNames[key] = com.AddSlashes(val)
			}
			sqlStr += `AND COLUMN_NAME IN ('` + strings.Join(fieldNames, `','`) + `')`
		}
	}
	stmt, err := db.PrepareContext(ctx, sqlStr)
	if factory.Debug() {
		fmt.Println(sqlStr, `[`, dbName, tableName, `]`)
	}
	if err != nil {
		return nil, nil, err
	}
	results := map[string]param.StringMap{}
	rows, err := stmt.QueryContext(ctx, dbName, tableName)
	if err != nil {
		return nil, nil, err
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	indexes := make(map[string]int, len(cols))
	for idx, col := range cols {
		indexes[col] = idx
	}
	fields := []string{}
	for rows.Next() {
		recv := make([]interface{}, len(cols))
		for idx := range cols {
			recv[idx] = interface{}(&sql.NullString{})
		}
		err := rows.Scan(recv...)
		if err != nil {
			return results, fields, err
		}
		result := param.StringMap{}
		for col, idx := range indexes {
			result[col] = param.String(recv[idx].(*sql.NullString).String)
		}
		field := result.String(`field`)
		results[field] = result
		fields = append(fields, field)
	}
	return results, fields, err
}
