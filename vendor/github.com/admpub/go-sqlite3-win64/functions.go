// +build windows
// Copyright (C) 2016 Samuel Melrose <sam@infitialis.com>.
//
// Based on work by Yasuhiro Matsumoto <mattn.jp@gmail.com>
// https://github.com/mattn/go-sqlite3
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package sqlite3

import (
	"database/sql/driver"
	"os"
	"path/filepath"
	"regexp"
	"unsafe"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func basePath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}

	return dir + string(os.PathSeparator)
}

func BytePtrToString(p *byte) string {
	var (
		sizeTest byte
		finalStr = make([]byte, 0)
	)
	for {
		if *p == byte(0) {
			break
		}

		finalStr = append(finalStr, *p)
		p = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + unsafe.Sizeof(sizeTest)))
	}
	return string(finalStr[0:])
}

var tableNameRegexp = regexp.MustCompile(`(?is)^[\s]*INSERT[\s]+INTO[\s]+([^(\s]+)`)

func LastInsertID(c *SQLiteConn, table string, isTableName bool) (int64, error) {
	if len(table) == 0 {
		return 0, nil
	}
	if !isTableName {
		matches := tableNameRegexp.FindStringSubmatch(table)
		if matches == nil {
			return 0, nil
		}
		table = matches[1]
	}

	rows, err := c.Query("SELECT last_insert_rowid() FROM `"+table+"`", nil)
	if err != nil {
		return 0, err
	}
	v := make([]driver.Value, 1)
	err = rows.Next(v)
	if err != nil {
		return 0, err
	}
	rowid, _ := v[0].(int64)
	return rowid, nil
}
