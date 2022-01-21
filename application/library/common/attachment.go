/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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

package common

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/errors"
	"github.com/admpub/nging/v4/application/registry/upload/helper"
)

// @deprecated 本文件作废

// IsRightUploadFile 是否是正确的上传文件
func IsRightUploadFile(ctx echo.Context, src string) error {
	return helper.IsRightUploadFile(ctx, src)
}

// RemoveAvatar 删除头像
func RemoveAvatar(ctx echo.Context, typ string, id uint64) error {
	userDir := filepath.Join(helper.UploadDir, typ, fmt.Sprint(id))
	if !com.IsDir(userDir) {
		return nil
	}
	err := os.RemoveAll(userDir)
	if err != nil {
		return err
	}
	return OnRemoveOwnerFile(ctx, typ, id, userDir)
}

// MoveAvatarToUserDir 移动临时文件夹中的头像到用户目录
func MoveAvatarToUserDir(ctx echo.Context, src string, typ string, id uint64) (string, error) {
	return MoveUploadedFileToOwnerDirCommon(ctx, src, typ, id, true)
}

// DirShardingNum 文件夹分组基数
const DirShardingNum = float64(50000)

var (
	reNumberOrWord = regexp.MustCompile(`^[0-9a-zA-Z_-]+$`)
	// ErrIncorrectFileOwnerID .
	ErrIncorrectFileOwnerID = errors.New("Incorrect File Owner ID")
)

// ValidFileOwnerID 验证文件宿主ID
func ValidFileOwnerID(id string) error {
	if reNumberOrWord.MatchString(id) {
		return nil
	}
	return ErrIncorrectFileOwnerID
}

// DirSharding 文件夹分组(暂不使用)
func DirSharding(id uint64) uint64 {
	return IDSharding(id, DirShardingNum)
}

// RemoveUploadedFile 删除被上传的文件
func RemoveUploadedFile(ctx echo.Context, typ string, id interface{}) error {
	idv := fmt.Sprint(id)
	if err := ValidFileOwnerID(idv); err != nil {
		return err
	}
	sdir := filepath.Join(helper.UploadDir, typ, idv)
	if !com.IsDir(sdir) {
		return nil
	}
	err := os.RemoveAll(sdir)
	if err != nil {
		return err
	}
	return OnRemoveOwnerFile(ctx, typ, id, sdir)
}

// OnUpdateOwnerFilePath 当更新文件路径时的通用操作
var OnUpdateOwnerFilePath = func(ctx echo.Context,
	src string, typ string, id interface{},
	newSavePath string, newViewURL string) error {
	return nil
}

// OnRemoveOwnerFile 当删除文件时的通用操作
var OnRemoveOwnerFile = func(ctx echo.Context, typ string, id interface{}, ownerDir string) error {
	return nil
}

// MoveUploadedFileToOwnerDir 移动上传的文件到所有者目录
func MoveUploadedFileToOwnerDir(ctx echo.Context, src string, typ string, id interface{}) (string, error) {
	return MoveUploadedFileToOwnerDirCommon(ctx, src, typ, id, false)
}

// MoveUploadedFileToOwnerDirCommon 移动上传的文件到所有者目录
func MoveUploadedFileToOwnerDirCommon(ctx echo.Context, src string, typ string, id interface{}, isAvatar bool) (string, error) {
	var newPath string
	if err := helper.IsRightUploadFile(ctx, src); err != nil {
		return newPath, err
	}
	name := path.Base(src)
	// unownedFile := filepath.Join(helper.UploadDir, typ, `0`, name)
	// 无主文件
	unownedFile := helper.URLToFile(src)
	if !com.FileExists(unownedFile) {
		return src, nil
	}
	idv := fmt.Sprint(id)
	if err := ValidFileOwnerID(idv); err != nil {
		return newPath, err
	}
	sdir := filepath.Join(helper.UploadDir, typ, idv)
	com.MkdirAll(sdir, os.ModePerm)
	// 迁移目的地
	ownedFile := sdir + echo.FilePathSeparator + name
	if isAvatar {
		ext := path.Ext(src)
		ownedFile = sdir + echo.FilePathSeparator + `avatar` + ext
	}
	err := os.Rename(unownedFile, ownedFile)
	if err != nil {
		return newPath, err
	}
	// 迁移后文件的访问网址
	newPath = helper.UploadURLPath + typ + `/` + idv + `/` + name
	p := strings.LastIndex(unownedFile, `.`)
	if p > 0 {
		filePrefix := unownedFile[0:p] + `_`
		unownedFiles, err := filepath.Glob(filePrefix + `*` + unownedFile[p:])
		if err != nil {
			return newPath, err
		}
		for _, file := range unownedFiles {
			name := filepath.Base(file)
			userFile := sdir + echo.FilePathSeparator + name
			if isAvatar {
				name = strings.TrimPrefix(name, filePrefix)
				userFile = sdir + echo.FilePathSeparator + `avatar_` + name
			}
			err := os.Rename(file, userFile)
			if err != nil {
				return newPath, err
			}
		}
	}
	return newPath, OnUpdateOwnerFilePath(ctx, src, typ, id, ownedFile, newPath)
}

// ModifyAsThumbnailName 将指向临时文件夹的缩略图路径改为新位置上的缩略图路径
// originName 为新位置上的原始图路径
// thumbnailName 为临时位置上的缩略图路径
func ModifyAsThumbnailName(originName, thumbnailName string) string {
	name := path.Base(thumbnailName)
	position := strings.Index(name, `_`)
	var suffix string
	if position > 0 {
		suffix = name[position:]
	}
	if len(suffix) > 0 {
		return originName[0:strings.LastIndex(originName, `.`)] + suffix
	}
	return originName
}

// Replacex 根据map替换
func Replacex(s string, oldAndNew map[string]string) string {
	for oldName, newName := range oldAndNew {
		s = strings.Replace(s, oldName, newName, -1)
	}
	return s
}
