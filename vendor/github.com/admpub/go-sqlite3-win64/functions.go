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
	"os"
	"path/filepath"
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
		finalStr []byte = make([]byte, 0)
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
