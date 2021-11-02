package sync

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	sqlCreateTableRegexp = regexp.MustCompile(`(?i)^CREATE TABLE `)
)

// ReplaceCharset 替换DDL语句中的字符集
func ReplaceCharset(sqlStr string, srcCharset string, destCharset string, checkCreateDDL ...bool) string {
	if destCharset == srcCharset {
		return sqlStr
	}
	if len(checkCreateDDL) > 0 && checkCreateDDL[0] {
		if !sqlCreateTableRegexp.MatchString(sqlStr) {
			return sqlStr
		}
	}
	srcCharsetQuoted := regexp.QuoteMeta(srcCharset)
	sqlCharsetRegexp := regexp.MustCompile(`(?i) (CHARACTER SET |CHARSET=)` + srcCharsetQuoted + ` `)
	sqlCollateRegexp := regexp.MustCompile(`(?i) (COLLATE[= ])` + srcCharsetQuoted + `([\w]+)`)
	sqlStr = sqlCharsetRegexp.ReplaceAllString(sqlStr, ` ${1}`+destCharset+` `)
	sqlStr = sqlCollateRegexp.ReplaceAllString(sqlStr, ` ${1}`+destCharset+`${2}`)
	return sqlStr
}

func GetSQLFileContent(filePattern string) (string, error) {
	var content string
	sourceSQLFiles, err := filepath.Glob(filePattern)
	if err != nil {
		return content, err
	}
	for _, sourceSQLFile := range sourceSQLFiles {
		if strings.HasSuffix(sourceSQLFile, `.sql`) {
			b, err := ioutil.ReadFile(sourceSQLFile)
			if err != nil {
				return content, err
			}
			content += "\n" + string(b)
		}
	}
	return content, nil
}

var HTMLDocTemplate = `
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Document</title>
  <link rel="stylesheet" href="https://b.webx.top/public/assets/backend/js/bootstrap/dist/css/bootstrap.min.css?t=20211101213650" />
</head>
<body>
  %s
</body>
</html>`
