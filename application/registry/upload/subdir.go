/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package upload

import (
	"fmt"
	"strings"
)

var subdirs = map[string]*SubdirInfo{
	`user`: &SubdirInfo{
		Allowed:     true,
		Key:         "user",
		Name:        "后台用户个人文件",
		Description: "",
	}, //后台用户个人文件
	`user-avatar`: &SubdirInfo{
		Allowed:     true,
		Key:         "user-avatar",
		Name:        "后台用户头像",
		Description: "",
	}, //后台用户头像
	`config`: &SubdirInfo{
		Allowed:     true,
		Key:         "config",
		Name:        "站点公告图片",
		Description: "",
		checker:     ConfigChecker,
	}, //后台系统设置中的图片
}

func SubdirRegister(subdir string, allow interface{}, nameAndDescription ...string) *SubdirInfo {
	var isAllow bool
	switch v := allow.(type) {
	case bool:
		isAllow = v
	case *SubdirInfo:
		return SubdirRegisterObject(subdir, v)
	case SubdirInfo:
		return SubdirRegisterObject(subdir, &v)
	default:
		panic(fmt.Sprintf(`Unsupported type: %T`, v))
	}
	var name, nameEN, description string
	switch len(nameAndDescription) {
	case 3:
		description = nameAndDescription[2]
		fallthrough
	case 2:
		nameEN = nameAndDescription[1]
		fallthrough
	case 1:
		name = nameAndDescription[0]
	}
	info := &SubdirInfo{
		Allowed:     isAllow,
		Key:         subdir,
		Name:        name,
		NameEN:      nameEN,
		Description: description,
	}

	r := strings.SplitN(info.Key, `-`, 2)
	switch len(r) {
	case 2:
		info.fieldName = r[1]
		fallthrough
	case 1:
		info.tableName = r[0]
	}
	SubdirRegisterObject(subdir, info)
	return info
}

func SubdirRegisterObject(subdir string, info *SubdirInfo) *SubdirInfo {
	subdirs[subdir] = info
	return info
}

func SubdirUnregister(subdirList ...string) {
	for _, subdir := range subdirList {
		_, ok := subdirs[subdir]
		if ok {
			delete(subdirs, subdir)
		}
	}
}

func SubdirAll() map[string]*SubdirInfo {
	return subdirs
}

func SubdirIsAllowed(subdir string) bool {
	info, ok := subdirs[subdir]
	if !ok || info == nil {
		return false
	}
	return info.Allowed
}

func SubdirGet(subdir string) *SubdirInfo {
	info, ok := subdirs[subdir]
	if !ok {
		return nil
	}
	return info
}

// CleanTempFile 清理临时文件
func CleanTempFile(prefix string, deleter func(folderPath string) error) error {
	if !strings.HasSuffix(prefix, `/`) {
		prefix += `/`
	}
	for subdir := range subdirs {
		err := deleter(prefix + subdir + `/0/`)
		if err != nil {
			return err
		}
	}
	return nil
}
