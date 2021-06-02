package common

import (
	"regexp"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	mysqlUtil "github.com/webx-top/db/lib/factory/mysql"
	"github.com/webx-top/db/mysql"
	"github.com/webx-top/echo"
)

func SQLLineParser(exec func(string) error) func(string) error {
	var sqlStr string
	return func(line string) error {
		if strings.HasPrefix(line, `--`) {
			return nil
		}
		line = strings.TrimRight(line, "\r ")
		if strings.HasPrefix(line, `/*`) && strings.HasSuffix(line, `*/;`) {
			return nil
		}
		sqlStr += line
		if strings.HasSuffix(line, `;`) {
			defer func() {
				sqlStr = ``
			}()
			//println(sqlStr)
			if sqlStr == `;` {
				return nil
			}
			return exec(sqlStr)
		}
		sqlStr += "\n"
		return nil
	}
}

func ParseSQL(sqlFile string, isFile bool, installer func(string) error) (err error) {
	installFunction := SQLLineParser(installer)
	if isFile {
		return com.SeekFileLines(sqlFile, installFunction)
	}
	sqlContent := sqlFile
	for _, line := range strings.Split(sqlContent, "\n") {
		err = installFunction(line)
		if err != nil {
			return err
		}
	}
	return err
}

// ReplacePrefix 替换前缀数据
func ReplacePrefix(m factory.Model, field string, oldPrefix string, newPrefix string) error {
	oldPrefix = com.AddSlashes(oldPrefix, '_', '%')
	value := db.Raw("REPLACE(`"+field+"`, ?, ?)", oldPrefix, newPrefix)
	return m.SetField(nil, field, value, field, db.Like(oldPrefix+`%`))
}

var (
	sqlCharsetRegexp     = regexp.MustCompile(`(?i) (CHARACTER SET |CHARSET=)utf8mb4 `)
	sqlCollateRegexp     = regexp.MustCompile(`(?i) (COLLATE[= ])utf8mb4_general_ci`)
	sqlCreateTableRegexp = regexp.MustCompile(`(?i)^CREATE TABLE `)
	mysqlNetworkRegexp   = regexp.MustCompile(`^[/]{2,}`)
)

// ReplaceCharset 替换DDL语句中的字符集
func ReplaceCharset(sqlStr string, charset string, checkCreateDDL ...bool) string {
	if charset == `utf8mb4` {
		return sqlStr
	}
	if len(checkCreateDDL) > 0 && checkCreateDDL[0] {
		if !sqlCreateTableRegexp.MatchString(sqlStr) {
			return sqlStr
		}
	}
	sqlStr = sqlCharsetRegexp.ReplaceAllString(sqlStr, ` ${1}`+charset+` `)
	sqlStr = sqlCollateRegexp.ReplaceAllString(sqlStr, ` ${1}`+charset+`_general_ci`)
	return sqlStr
}

func ParseMysqlConnectionURL(settings *mysql.ConnectionURL) {
	if strings.HasPrefix(settings.Host, `unix:`) {
		settings.Socket = strings.TrimPrefix(settings.Host, `unix:`)
		settings.Socket = mysqlNetworkRegexp.ReplaceAllString(settings.Socket, `/`)
		settings.Host = ``
	}
}

func SelectPageCond(ctx echo.Context, cond *db.Compounds, pkAndLabelFields ...string) {
	pk := `id`
	lb := `name`
	switch len(pkAndLabelFields) {
	case 2:
		if len(pkAndLabelFields[1]) > 0 {
			lb = pkAndLabelFields[1]
		}
		fallthrough
	case 1:
		if len(pkAndLabelFields[0]) > 0 {
			pk = pkAndLabelFields[0]
		}
	}
	searchValue := ctx.Formx(`searchValue`).String()
	if len(searchValue) > 0 {
		cond.AddKV(pk, searchValue)
	} else {
		keywords := ctx.FormValues(`q_word[]`)
		andOr := ctx.Form(`andOr`)
		var sep string
		if len(andOr) == 0 || andOr == `AND` {
			sep = ` `
		} else {
			sep = `||`
		}
		q := strings.Join(keywords, `"`+sep+`"`)
		if len(q) == 0 {
			q = ctx.Formx(`q`).String()
		} else {
			q = `"` + q + `"`
		}
		if len(q) > 0 {
			cond.From(mysqlUtil.SearchField(lb, q, pk))
		}
	}
}
