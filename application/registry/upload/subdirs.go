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

import "strings"

// SubdirInfo 子目录信息
type SubdirInfo struct {
	Allowed     bool
	Key         string
	Name        string
	NameEN      string
	Description string
	Checker     Checker `json:"checker" xml:"-"`
}

func (i *SubdirInfo) String() string {
	if len(i.Name) > 0 {
		return i.Name
	}
	return i.Key
}

func (i *SubdirInfo) GetNameEN() string {
	if len(i.NameEN) > 0 {
		return i.NameEN
	}
	return i.Name
}

func (i *SubdirInfo) SetChecker(checker Checker) *SubdirInfo {
	i.Checker = checker
	return i
}

func (i *SubdirInfo) SetAllowed(allowed bool) *SubdirInfo {
	i.Allowed = allowed
	return i
}

func (i *SubdirInfo) SetName(name string) *SubdirInfo {
	i.Name = name
	return i
}

func (i *SubdirInfo) SetNameEN(nameEN string) *SubdirInfo {
	i.NameEN = nameEN
	return i
}

func (i *SubdirInfo) SetDescription(description string) *SubdirInfo {
	i.Description = description
	return i
}

var subdirs = map[string]*SubdirInfo{
	`user`: &SubdirInfo{
		Allowed:     true,
		Key:         "user",
		Name:        "后台用户个人文件",
		Description: "",
	}, //后台用户个人文件
	`customer`: &SubdirInfo{
		Allowed:     true,
		Key:         "customer",
		Name:        "前台客户个人文件",
		Description: "",
	}, //前台客户个人文件
	`comment`: &SubdirInfo{
		Allowed:     true,
		Key:         "comment",
		Name:        "评论附件",
		Description: "",
	}, //评论附件
	`customer-avatar`: &SubdirInfo{
		Allowed:     true,
		Key:         "customer-avatar",
		Name:        "前台客户头像",
		Description: "",
	}, //前台客户头像
	`user-avatar`: &SubdirInfo{
		Allowed:     true,
		Key:         "user-avatar",
		Name:        "后台用户头像",
		Description: "",
		Checker:     UserAvatarChecker,
	}, //后台用户头像
	`pay-product`: &SubdirInfo{
		Allowed:     true,
		Key:         "pay-product",
		Name:        "产品图片",
		Description: "",
	}, //产品图片
	`product-version`: &SubdirInfo{
		Allowed:     true,
		Key:         "product-version",
		Name:        "产品版本文件",
		Description: "",
	}, //产品版本文件
	`site-announcement`: &SubdirInfo{
		Allowed:     true,
		Key:         "site-announcement",
		Name:        "站点公告图片",
		Description: "",
	}, //站点公告图片
	`news`: &SubdirInfo{
		Allowed:     true,
		Key:         "news",
		Name:        "新闻图片",
		Description: "",
	}, //新闻图片
}

func SubdirRegister(subdir string, allow bool, nameAndDescription ...string) *SubdirInfo {
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
		Key:         subdir,
		Name:        name,
		NameEN:      nameEN,
		Description: description,
	}
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
