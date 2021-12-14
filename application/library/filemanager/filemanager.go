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

package filemanager

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/admpub/nging/v4/application/library/charset"
	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var (
	EncodedSep   = com.URLEncode(`/`)
	EncodedSlash = com.URLEncode(`\`)
	EncodedSepa  = com.URLEncode(echo.FilePathSeparator)
)

func New(root string, editableMaxSize int64, ctx echo.Context) *fileManager {
	return &fileManager{
		Context:         ctx,
		Root:            root,
		EditableMaxSize: editableMaxSize,
	}
}

type fileManager struct {
	echo.Context
	Root            string
	EditableMaxSize int64
}

func (f *fileManager) RealPath(filePath string) string {
	absPath := f.Root
	if len(filePath) > 0 {
		filePath = filepath.Clean(filePath)
		absPath = filepath.Join(f.Root, filePath)
	}
	return absPath
}

func (f *fileManager) Edit(absPath string, content string, encoding string) (interface{}, error) {
	fi, err := os.Stat(absPath)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return nil, errors.New(f.T(`不能编辑文件夹`))
	}
	if f.EditableMaxSize > 0 && fi.Size() > f.EditableMaxSize {
		return nil, errors.New(f.T(`很抱歉，不支持编辑超过%v的文件`, com.FormatByte(f.EditableMaxSize)))
	}
	encoding = strings.ToLower(encoding)
	isUTF8 := encoding == `` || encoding == `utf-8`
	if f.IsPost() {
		b := []byte(content)
		if !isUTF8 {
			b, err = charset.Convert(`utf-8`, encoding, b)
			if err != nil {
				return ``, err
			}
		}
		err = ioutil.WriteFile(absPath, b, fi.Mode())
		return nil, err
	}
	b, err := ioutil.ReadFile(absPath)
	if err == nil && !isUTF8 {
		b, err = charset.Convert(encoding, `utf-8`, b)
	}
	return string(b), err
}

func (f *fileManager) Remove(absPath string) error {
	fi, err := os.Stat(absPath)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return os.RemoveAll(absPath)
	}
	return os.Remove(absPath)
}

func (f *fileManager) Mkdir(absPath string, mode os.FileMode) error {
	return com.MkdirAll(absPath, mode)
}

func (f *fileManager) Rename(absPath string, newName string) (err error) {
	if len(newName) > 0 {
		err = os.Rename(absPath, filepath.Join(filepath.Dir(absPath), filepath.Base(newName)))
	} else {
		err = errors.New(f.T(`请输入有效的文件名称`))
	}
	return
}

func (f *fileManager) enterPath(absPath string) (d http.File, fi os.FileInfo, err error) {
	absPath = strings.TrimRight(absPath, `/`)
	absPath = strings.TrimRight(absPath, `\`)
	fs := http.Dir(filepath.Dir(absPath))
	fileName := filepath.Base(absPath)
	d, err = fs.Open(fileName)
	if err != nil {
		return
	}
	//defer d.Close()
	fi, err = d.Stat()
	return
}

func (f *fileManager) Upload(absPath string,
	chunkUpload *uploadClient.ChunkUpload,
	chunkOpts ...uploadClient.ChunkInfoOpter) (err error) {
	var (
		d  http.File
		fi os.FileInfo
	)
	d, fi, err = f.enterPath(absPath)
	if d != nil {
		defer d.Close()
	}
	if err != nil {
		return
	}
	if !fi.IsDir() {
		return errors.New(f.T(`路径不正确: %s`, absPath))
	}
	var filePath string
	var chunked bool // 是否支持分片
	if chunkUpload != nil {
		_, err := chunkUpload.Upload(f.Request().StdRequest(), chunkOpts...)
		if err != nil {
			if !errors.Is(err, uploadClient.ErrChunkUnsupported) {
				if errors.Is(err, uploadClient.ErrChunkUploadCompleted) ||
					errors.Is(err, uploadClient.ErrFileUploadCompleted) {
					return nil
				}
				return err
			}
		} else {
			if !chunkUpload.Merged() {
				return nil
			}
			chunked = true
			filePath = chunkUpload.GetSavePath()
		}
	}
	if !chunked {
		fileHdr, err := f.SaveUploadedFile(`file`, absPath)
		if err != nil {
			return err
		}
		filePath = filepath.Join(absPath, fileHdr.Filename)
	}
	pipe := f.Form(`pipe`)
	switch pipe {
	case `unzip`:
		err = com.Unzip(filePath, absPath)
		if err == nil {
			err = os.Remove(filePath)
			if err != nil {
				err = errors.New(f.T(`压缩包已经成功解压，但是删除压缩包失败：`) + err.Error())
			}
		}
		return err
	default:
		if chunked {
			newfile := filepath.Join(absPath, filepath.Base(filePath))
			err = os.Rename(filePath, newfile)
			if err != nil {
				return fmt.Errorf(`move %s to %s: %w`, filePath, newfile, err)
			}
		}
	}
	return
}

func (f *fileManager) List(absPath string, sortBy ...string) (err error, exit bool, dirs []os.FileInfo) {
	var (
		d  http.File
		fi os.FileInfo
	)
	d, fi, err = f.enterPath(absPath)
	if d != nil {
		defer d.Close()
	}
	if err != nil {
		return
	}
	if !fi.IsDir() {
		fileName := filepath.Base(absPath)
		inline := f.Formx(`inline`).Bool()
		return f.Attachment(d, fileName, inline), true, nil
	}

	dirs, err = d.Readdir(-1)
	if len(sortBy) > 0 {
		switch sortBy[0] {
		case `time`:
			sort.Sort(SortByModTime(dirs))
		case `-time`:
			sort.Sort(SortByModTimeDesc(dirs))
		case `name`:
		case `-name`:
			sort.Sort(SortByNameDesc(dirs))
		case `type`:
			fallthrough
		default:
			sort.Sort(SortByFileType(dirs))
		}
	} else {
		sort.Sort(SortByFileType(dirs))
	}
	if f.Format() == "json" {
		dirList, fileList := f.ListTransfer(dirs)
		data := f.Data()
		data.SetData(echo.H{
			`dirList`:  dirList,
			`fileList`: fileList,
		})
		return f.JSON(data), true, nil
	}
	return
}

func (f *fileManager) ListTransfer(dirs []os.FileInfo) (dirList []echo.H, fileList []echo.H) {
	dirList = []echo.H{}
	fileList = []echo.H{}
	for _, d := range dirs {
		item := echo.H{
			`name`:  d.Name(),
			`size`:  d.Size(),
			`mode`:  d.Mode().String(),
			`mtime`: d.ModTime().Format(`2006-01-02 15:04:05`),
			//`sys`:   d.Sys(),
		}
		if d.IsDir() {
			dirList = append(dirList, item)
			continue
		}
		fileList = append(fileList, item)
	}
	return
}
