package internal

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"reflect"
	"regexp"
	"strings"
)

// Version version
const Version = "0.4"

// AppURL site
const AppURL = "https://github.com/admpub/mysql-schema-sync/"

const timeFormatStd string = "2006-01-02 15:04:05"

// loadJsonFile load json
func loadJSONFile(jsonPath string, val interface{}) error {
	bs, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return err
	}
	return ParseJSON(string(bs), val)
}

// ParseJSON parse json
func ParseJSON(content string, val interface{}) error {
	lines := strings.Split(content, "\n")
	var bf bytes.Buffer
	for _, line := range lines {
		lineNew := strings.TrimSpace(line)
		if (len(lineNew) > 0 && lineNew[0] == '#') || (len(lineNew) > 1 && lineNew[0:2] == "//") {
			continue
		}
		bf.WriteString(lineNew)
	}
	return json.Unmarshal(bf.Bytes(), &val)
}

func inStringSlice(str string, strSli []string) bool {
	for _, v := range strSli {
		if str == v {
			return true
		}
	}
	return false
}

func simpleMatch(patternStr string, str string, msg ...string) bool {
	str = strings.TrimSpace(str)
	patternStr = strings.TrimSpace(patternStr)
	if patternStr == str {
		log.Println("simple_match:suc,equal", msg, "patternStr:", patternStr, "str:", str)
		return true
	}
	pattern := "^" + strings.Replace(patternStr, "*", `.*`, -1) + "$"
	match, err := regexp.MatchString(pattern, str)
	if err != nil {
		log.Println("simple_match:error", msg, "patternStr:", patternStr, "pattern:", pattern, "str:", str, "err:", err)
	}
	if match {
		//log.Println("simple_match:suc", msg, "patternStr:", patternStr, "pattern:", pattern, "str:", str)
	}
	return match
}

func htmlPre(str string) string {
	return "<pre>" + html.EscapeString(str) + "</pre>"
}

func dsnSort(dsn string) string {
	i := strings.Index(dsn, "@")
	if i < 1 {
		return dsn
	}
	return dsn[i+1:]
}

func maxMapKeyLen(data interface{}, ext int) string {
	l := 0
	for _, k := range reflect.ValueOf(data).MapKeys() {
		if k.Len() > l {
			l = k.Len()
		}
	}
	return fmt.Sprintf("%d", l+ext)
}

func Exec(mydb *sql.DB, query string) (sql.Result, error) {
	var sqlStr string
	var ret sql.Result
	tx, err := mydb.Begin()
	if err != nil {
		return nil, err
	}
	execute := func(line string) (rErr error) {
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
			if sqlStr == `;` {
				return nil
			}
			ret, err = tx.Exec(sqlStr)
			log.Println("exec_one:[", sqlStr, "]", err)
			return err
		}
		sqlStr += "\n"
		return nil
	}
	for _, line := range strings.Split(query, "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		err = execute(line)
		if err != nil {
			break
		}
	}
	if err == nil {
		sqlStr = strings.TrimSpace(sqlStr)
		if len(sqlStr) > 0 {
			if sqlStr != `;` {
				ret, err = tx.Exec(sqlStr)
				log.Println("exec_one:[", sqlStr, "]", err)
			}
			sqlStr = ``
		}
	}
	if err == nil {
		err = tx.Commit()
	} else {
		tx.Rollback()
	}
	return ret, err
}
