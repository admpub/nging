/*
 * Copyright 1999-2020 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package file

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var osType string
var path string

const WINDOWS = "windows"

func init() {
	osType = runtime.GOOS
	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		path = "\\"
	} else {
		path = "/"
	}
}

func MkdirIfNecessary(createDir string) (err error) {
	s := strings.Split(createDir, path)
	startIndex := 0
	dir := ""
	if s[0] == "" {
		startIndex = 1
	} else {
		dir, _ = os.Getwd() //当前的目录
	}
	for i := startIndex; i < len(s); i++ {
		var d string
		if osType == WINDOWS && filepath.IsAbs(createDir) {
			d = strings.Join(s[startIndex:i+1], path)
		} else {
			d = dir + path + strings.Join(s[startIndex:i+1], path)
		}
		if _, e := os.Stat(d); os.IsNotExist(e) {
			err = os.Mkdir(d, os.ModePerm) //在当前目录下生成md目录
			if err != nil {
				break
			}
		}
	}

	return err
}

func GetCurrentPath() string {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println("can not get current path")
	}
	return dir
}
