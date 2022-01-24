package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/webx-top/com"
	dbPkg "github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo/param"
)

var (
	reTable = regexp.MustCompile("CREATE TABLE `[^`]+` \\(")
)

// CopyTableStruct 复制表结构到新表
func CopyTableStruct(srcLinkID int, srcDBName string, srcTableName string,
	destLinkID int, destDBName string, destTableName string) error {
	ddl, err := CreateTableSQL(srcLinkID, srcDBName, srcTableName)
	if err != nil {
		return err
	}
	db := factory.NewParam().SetIndex(destLinkID).DB()
	ddl = reTable.ReplaceAllString(ddl, "CREATE TABLE `"+destTableName+"` (")
	if factory.Debug() {
		fmt.Println(ddl)
	}
	_, err = db.Exec(ddl)
	return err
}

// ReplacePrefix UPDATE table_name SET `field`=REPLACE(`field`,'oldPrefix','newPrefix') WHERE `field` LIKE 'oldPrefix%';
func ReplacePrefix(linkID int, tableName string, field string, oldPrefix string, newPrefix string, notPrefix bool, extWheres ...string) (int64, error) {
	var extWhere string
	if len(extWheres) > 0 {
		extWhere = extWheres[0] + ` AND `
	}
	oldPrefix = com.AddSlashes(oldPrefix, '_', '%')
	tableName = strings.ReplaceAll(tableName, "`", "``")
	field = strings.ReplaceAll(field, "`", "``")
	sqlStr := "UPDATE `" + tableName + "` SET `" + field + "`=REPLACE(`" + field + "`,?,?) WHERE " + extWhere + "`" + field + "` LIKE ?"
	db := factory.NewParam().SetIndex(linkID).DB()
	likeValue := oldPrefix + `%`
	if notPrefix {
		likeValue = `%` + likeValue
	}
	result, err := db.Exec(sqlStr, oldPrefix, newPrefix, likeValue)
	if factory.Debug() {
		fmt.Println(dbPkg.BuildSQL(sqlStr, oldPrefix, newPrefix, likeValue))
	}
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	return affected, err
	//value := db.Raw("REPLACE(`"+field+"`, ?, ?)", oldPrefix, newPrefix)
	//return m.UpdateField(nil, field, value, field, db.Like(oldPrefix+`%`))
}

// TableExists 查询表是否存在
func TableExists(linkID int, dbName string, tableName string) (bool, error) {
	ctx := context.Background()
	db := factory.NewParam().SetIndex(linkID).DB()
	stmt, err := db.PrepareContext(ctx, SQLTableExists)
	if factory.Debug() {
		fmt.Println(SQLTableExists, `[`, dbName, tableName, `]`)
	}
	if err != nil {
		return false, err
	}
	recv := sql.NullString{}
	err = stmt.QueryRowContext(ctx, dbName, tableName).Scan(&recv)
	if err != nil {
		return false, err
	}
	return param.String(recv.String).Int64() > 0, err
}

// FieldExists 查询表字段是否存在
func FieldExists(linkID int, tableName string, fieldName string) (bool, error) {
	ctx := context.Background()
	db := factory.NewParam().SetIndex(linkID).DB()
	stmt, err := db.PrepareContext(ctx, SQLFieldExists)
	if factory.Debug() {
		fmt.Println(SQLFieldExists, `[`, tableName, fieldName, `]`)
	}
	if err != nil {
		return false, err
	}
	recv := sql.NullString{}
	err = stmt.QueryRowContext(ctx, tableName, fieldName).Scan(&recv)
	if err != nil {
		return false, err
	}
	return param.String(recv.String).Int64() > 0, err
}

// MoveToTable 移动数据到新表
func MoveToTable(linkID int, dbName string, srcTableName string, destTableName string, condition string, src2DestFieldMapping ...map[string]string) (int64, error) {
	exists, err := TableExists(linkID, dbName, destTableName)
	if err != nil {
		return 0, err
	}
	if !exists {
		err = CopyTableStruct(linkID, dbName, srcTableName,
			linkID, dbName, destTableName)
		if err != nil {
			return 0, err
		}
	}
	sqlStr := "INSERT IGNORE INTO `" + destTableName + "`"
	var srcFields, destFields string
	if len(src2DestFieldMapping) > 0 {
		var sep string
		for srcField, destField := range src2DestFieldMapping[0] {
			srcFields += sep + "`" + srcField + "`"
			destFields += sep + "`" + destField + "`"
			sep = ","
		}
		destFields = `(` + destFields + `)`
	}
	sqlStr += destFields
	if len(srcFields) < 1 {
		srcFields = `*`
	}
	sqlStr += ` SELECT ` + srcFields + ` FROM ` + "`" + srcTableName + "`"
	if len(condition) > 0 {
		sqlStr += ` WHERE ` + condition
	}
	db := factory.NewParam().SetIndex(linkID).DB()
	result, err := db.Exec(sqlStr)
	if factory.Debug() {
		fmt.Println(sqlStr)
	}
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return affected, err
	}
	if affected > 0 {
		sqlStr := "DELETE FROM `" + srcTableName + "`"
		if len(condition) > 0 {
			sqlStr += ` WHERE ` + condition
		}
		db := factory.NewParam().SetIndex(linkID).DB()
		result, err := db.Exec(sqlStr)
		if factory.Debug() {
			fmt.Println(sqlStr)
		}
		if err != nil {
			return affected, err
		}
		affected, err = result.RowsAffected()
	}
	return affected, err
}
