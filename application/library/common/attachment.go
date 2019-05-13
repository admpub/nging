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

package common

import (
	stdErr "errors"
	"fmt"
	"math"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/admpub/nging/application/registry/upload/helper"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

// AllowedUploadFileExtensions 被允许上传的文件的扩展名
var AllowedUploadFileExtensions = []string{
	`.jpeg`, `.jpg`, `.gif`, `.png`,
}

// IsRightUploadFile 是否是正确的上传文件
func IsRightUploadFile(ctx echo.Context, src string) error {
	src = path.Clean(src)
	ext := strings.ToLower(path.Ext(src))
	var ok bool
	var invalidExt string
	for _, ex := range AllowedUploadFileExtensions {
		if ext == ex {
			ok = true
			invalidExt = ext
			break
		}
	}
	if !ok {
		return stdErr.New(ctx.T(`不支持的文件扩展名`) + `: ` + invalidExt)
	}
	if !strings.Contains(src, helper.UploadDir) {
		return stdErr.New(ctx.T(`路径不合法`))
	}
	return nil
}

// RemoveAvatar 删除头像
func RemoveAvatar(typ string, id uint64) error {
	userDir := filepath.Join(echo.Wd(), strings.TrimPrefix(helper.UploadDir, `/`)+typ, fmt.Sprint(id))
	if !com.IsDir(userDir) {
		return nil
	}
	return os.RemoveAll(userDir)
}

// MoveAvatarToUserDir 移动临时文件夹中的头像到用户目录
func MoveAvatarToUserDir(ctx echo.Context, src string, typ string, id uint64) (string, error) {
	if strings.Contains(src, `://`) {
		return src, nil
	}
	var newPath string
	if err := IsRightUploadFile(ctx, src); err != nil {
		return newPath, err
	}
	name := path.Base(src)
	updir := strings.TrimPrefix(helper.UploadDir, `/`)
	guestFile := filepath.Join(echo.Wd(), updir+typ+`/0`, name)
	if !com.FileExists(guestFile) {
		return src, nil
	}
	userDir := filepath.Join(echo.Wd(), updir+typ, fmt.Sprint(id))
	os.MkdirAll(userDir, os.ModePerm)
	ext := path.Ext(src)
	userFile := userDir + echo.FilePathSeparator + `avatar` + ext
	err := os.Rename(guestFile, userFile)
	if err != nil {
		return newPath, err
	}
	newPath = helper.UploadDir + typ + `/` + fmt.Sprint(id) + `/` + name
	p := strings.LastIndex(guestFile, `.`)
	if p > 0 {
		filePrefix := guestFile[0:p] + `_`
		guestFiles, err := filepath.Glob(filePrefix + `*` + guestFile[p:])
		if err != nil {
			return newPath, err
		}
		for _, file := range guestFiles {
			name := filepath.Base(file)
			name = strings.TrimPrefix(name, filePrefix)
			userFile := userDir + echo.FilePathSeparator + `avatar_` + name
			err := os.Rename(file, userFile)
			if err != nil {
				return newPath, err
			}
		}
	}
	return newPath, err
}

// DirSharding 文件夹分组(暂不使用)
func DirSharding(id uint64) uint64 {
	return uint64(math.Ceil(float64(id) / float64(50000)))
}

// RemoveUploadedFile 删除被上传的文件
func RemoveUploadedFile(typ string, id uint64) error {
	sdir := filepath.Join(echo.Wd(), strings.TrimPrefix(helper.UploadDir, `/`)+typ, fmt.Sprint(id))
	if !com.IsDir(sdir) {
		return nil
	}
	return os.RemoveAll(sdir)
}

// MoveUploadedFileToOwnerDir 移动上传的文件到所有者目录
func MoveUploadedFileToOwnerDir(ctx echo.Context, src string, typ string, id uint64) (string, error) {
	var newPath string
	if err := IsRightUploadFile(ctx, src); err != nil {
		return newPath, err
	}
	name := path.Base(src)
	updir := strings.TrimPrefix(helper.UploadDir, `/`)
	unownedFile := filepath.Join(echo.Wd(), updir+typ+`/0`, name)
	if !com.FileExists(unownedFile) {
		return src, nil
	}
	sdir := filepath.Join(echo.Wd(), updir+typ, fmt.Sprint(id))
	os.MkdirAll(sdir, os.ModePerm)
	ownedFile := sdir + echo.FilePathSeparator + name
	err := os.Rename(unownedFile, ownedFile)
	if err != nil {
		return newPath, err
	}
	newPath = helper.UploadDir + typ + `/` + fmt.Sprint(id) + `/` + name
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
			err := os.Rename(file, userFile)
			if err != nil {
				return newPath, err
			}
		}
	}
	return newPath, err
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

var (
	temporaryFileRegexp  = regexp.MustCompile(helper.UploadDir + `[\w-]+/0/[\w]+\.[a-zA-Z]+`)
	persistentFileRegexp = regexp.MustCompile(helper.UploadDir + `[\w-]+/([^0]|[0-9]{2,})/[\w]+\.[a-zA-Z]+`)
)

// ParseTemporaryFileName 从文本中解析出临时文件名称
func ParseTemporaryFileName(s string) []string {
	files := temporaryFileRegexp.FindAllString(s, -1)
	return files
}

// Replacex 根据map替换
func Replacex(s string, oldAndNew map[string]string) string {
	for oldName, newName := range oldAndNew {
		s = strings.Replace(s, oldName, newName, -1)
	}
	return s
}

// MoveEmbedTemporaryFiles 转移被嵌入到文本内容中临时文件
func MoveEmbedTemporaryFiles(ctx echo.Context, content string, typ string, id uint64) (int, string, error) {
	files := ParseTemporaryFileName(content)
	oldAndNew := map[string]string{}
	for _, fileN := range files {
		if _, ok := oldAndNew[fileN]; ok {
			continue
		}
		newPath, err := MoveUploadedFileToOwnerDir(ctx, fileN, typ, id)
		if err != nil {
			return 0, content, err
		}
		oldAndNew[fileN] = newPath
	}
	return len(oldAndNew), Replacex(content, oldAndNew), nil
}
