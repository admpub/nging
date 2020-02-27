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

	"github.com/admpub/color"
	"github.com/admpub/log"
)

var subdirs = map[string]*SubdirInfo{
	`user`: (&SubdirInfo{
		Allowed:     true,
		Key:         "user",
		Name:        "后台用户",
		Description: "",
	}).SetTableName("user").SetFieldName(`:个人文件`, `avatar:头像`), //后台用户文件
	`config`: (&SubdirInfo{
		Allowed:     true,
		Key:         "config",
		Name:        "站点公告图片",
		Description: "",
		checker:     ConfigChecker,
	}).SetTableName("config"), //后台系统设置中的图片
}

func SubdirRegister(subdir interface{}, nameAndDescription ...string) *SubdirInfo {
	var key string
	switch v := subdir.(type) {
	case string:
		key = v
	case *SubdirInfo:
		return SubdirRegisterObject(v)
	case SubdirInfo:
		return SubdirRegisterObject(&v)
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
		Allowed:     true,
		Key:         key,
		Name:        name,
		NameEN:      nameEN,
		Description: description,
	}

	r := strings.SplitN(info.Key, `-`, 2)
	switch len(r) {
	case 2:
		info.SetFieldName(r[1])
		fallthrough
	case 1:
		info.tableName = r[0]
	}
	SubdirRegisterObject(info)
	return info
}

func SubdirRegisterObject(info *SubdirInfo) *SubdirInfo {
	in, ok := subdirs[info.Key]
	if ok {
		return in.CopyFrom(info)
	}
	subdirs[info.Key] = info
	log.Info(color.MagentaString(`subdir.register:`), info.Key)
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

func SubdirIsAllowed(subdir string, defaults ...string) bool {
	info, ok := subdirs[subdir]
	if !ok || info == nil {
		if len(defaults) > 0 {
			return SubdirIsAllowed(defaults[0])
		}
		return false
	}
	return info.Allowed
}

func SubdirGet(subdir string) *SubdirInfo {
	info, ok := subdirs[subdir]
	if !ok {
		return SubdirRegisterObject(NewSubdirInfo(subdir, ``))
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
