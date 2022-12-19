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

package sftpmanager

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/admpub/nging/v5/application/library/charset"
	"github.com/admpub/nging/v5/application/library/filemanager"

	"github.com/pkg/sftp"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	uploadClient "github.com/webx-top/client/upload"
)

func New(connector Connector, config *Config, editableMaxSize int) *SftpManager {
	if connector == nil {
		connector = DefaultConnector
	}
	return &SftpManager{
		connector:       connector,
		config:          config,
		EditableMaxSize: editableMaxSize,
	}
}

type Connector func(*Config) (*sftp.Client, error)

type SftpManager struct {
	client          *sftp.Client
	config          *Config
	connector       Connector
	connerror       error
	EditableMaxSize int
}

func (s *SftpManager) Connect() error {
	s.client, s.connerror = s.connector(s.config)
	return s.connerror
}

func (s *SftpManager) Client() *sftp.Client {
	if s.client == nil {
		s.Connect()
	}
	return s.client
}

func (s *SftpManager) ConnError() error {
	return s.connerror
}

func (s *SftpManager) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

func (s *SftpManager) Edit(ctx echo.Context, ppath string, content string, encoding string) (interface{}, error) {
	c := s.Client()
	if c == nil {
		return nil, s.ConnError()
	}
	f, err := c.Open(ppath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return nil, ctx.E(`不能编辑文件夹`)
	}
	if s.EditableMaxSize > 0 && fi.Size() > int64(s.EditableMaxSize) {
		return nil, ctx.E(`很抱歉，不支持编辑超过%v的文件`, com.FormatBytes(s.EditableMaxSize))
	}
	encoding = strings.ToLower(encoding)
	isUTF8 := len(encoding) == 0 || encoding == `utf-8`
	if ctx.IsPost() {
		b := []byte(content)
		if !isUTF8 {
			b, err = charset.Convert(`utf-8`, encoding, b)
			if err != nil {
				return nil, err
			}
		}
		f.Close()
		r := bytes.NewReader(b)
		f, err = c.OpenFile(ppath, os.O_CREATE|os.O_RDWR|os.O_TRUNC)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(f, r)
		if err != nil {
			return nil, ctx.E(ppath + `:` + err.Error())
		}
		return nil, err
	}

	dat, err := io.ReadAll(f)
	if err == nil && !isUTF8 {
		dat, err = charset.Convert(encoding, `utf-8`, dat)
	}
	return string(dat), err
}

func (s *SftpManager) Mkdir(ctx echo.Context, ppath, newName string) error {
	c := s.Client()
	if c == nil {
		return s.ConnError()
	}
	dirPath := path.Join(ppath, newName)
	f, err := c.Open(dirPath)
	if err == nil {
		finfo, err := f.Stat()
		if err != nil {
			return err
		}
		if finfo.IsDir() {
			return ctx.E(`已经存在相同名称的文件夹`)
		}
		return ctx.E(`已经存在相同名称的文件`)
	}
	if !os.IsNotExist(err) {
		return err
	}
	err = c.Mkdir(dirPath)
	return err
}

func (s *SftpManager) Rename(ctx echo.Context, ppath, newName string) error {
	if !strings.HasPrefix(newName, `/`) {
		newName = path.Join(path.Dir(ppath), newName)
	}
	c := s.Client()
	if c == nil {
		return s.ConnError()
	}
	_, err := c.Stat(newName)
	if err == nil {
		return ctx.E(`重命名失败，文件“%s”已经存在`, newName)
	}
	return c.Rename(ppath, newName)
}

func (s *SftpManager) Chown(ctx echo.Context, ppath string, uid, gid int) error {
	c := s.Client()
	if c == nil {
		return s.ConnError()
	}
	return c.Chown(ppath, uid, gid)
}

func (s *SftpManager) Chmod(ctx echo.Context, ppath string, mode os.FileMode) error {
	c := s.Client()
	if c == nil {
		return s.ConnError()
	}
	return c.Chmod(ppath, mode)
}

func (s *SftpManager) Search(ppath string, prefix string, num int) []string {
	var paths []string
	c := s.Client()
	if c == nil {
		return []string{}
	}
	dirs, _ := c.ReadDir(ppath)
	for _, d := range dirs {
		if len(paths) >= num {
			break
		}
		name := d.Name()
		if strings.HasPrefix(name, prefix) {
			paths = append(paths, name)
			continue
		}
	}
	return paths
}

func (s *SftpManager) Remove(ppath string) error {
	c := s.Client()
	if c == nil {
		return s.ConnError()
	}
	return c.Remove(ppath)
}

func (s *SftpManager) Upload(ctx echo.Context, ppath string,
	chunkUpload *uploadClient.ChunkUpload,
	chunkOpts ...uploadClient.ChunkInfoOpter) error {
	c := s.Client()
	if c == nil {
		return s.ConnError()
	}
	d, err := c.Open(ppath)
	if err != nil {
		return err
	}
	defer d.Close()
	fi, err := d.Stat()
	if err == nil && !fi.IsDir() {
		return ctx.E(`路径不正确`)
	}
	var fileSrc io.Reader
	var filename string
	var chunked bool // 是否支持分片
	if chunkUpload != nil {
		_, err := chunkUpload.Upload(ctx.Request().StdRequest(), chunkOpts...)
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
			_fp, err := os.Open(chunkUpload.GetSavePath())
			if err != nil {
				return err
			}
			fileSrc = _fp
			defer func() {
				_fp.Close()
				os.Remove(chunkUpload.GetSavePath())
			}()
			chunked = true
			filename = filepath.Base(chunkUpload.GetSavePath())
		}
	}
	if !chunked {
		_fileSrc, _fileHdr, err := ctx.Request().FormFile(`file`)
		if err != nil {
			return err
		}
		fileSrc = _fileSrc
		defer _fileSrc.Close()

		// Destination
		filename = _fileHdr.Filename

	}
	fileDst, err := c.Create(path.Join(ppath, filename))
	if err != nil {
		return err
	}
	defer fileDst.Close()

	_, err = io.Copy(fileDst, fileSrc)
	return err
}

func (s *SftpManager) List(ctx echo.Context, ppath string, sortBy ...string) (err error, exit bool, dirs []os.FileInfo) {
	c := s.Client()
	if c == nil {
		return s.ConnError(), false, nil
	}
	d, err := c.Open(ppath)
	if err != nil {
		return err, false, nil
	}
	defer d.Close()
	fi, err := d.Stat()
	if !fi.IsDir() {
		fileName := path.Base(ppath)
		inline := ctx.Formx(`inline`).Bool()
		return ctx.Attachment(d, fileName, fi.ModTime(), inline), true, nil
	}

	dirs, err = c.ReadDir(ppath)
	if len(sortBy) > 0 {
		switch sortBy[0] {
		case `time`:
			sort.Sort(filemanager.SortByModTime(dirs))
		case `-time`:
			sort.Sort(filemanager.SortByModTimeDesc(dirs))
		case `name`:
		case `-name`:
			sort.Sort(filemanager.SortByNameDesc(dirs))
		case `type`:
			fallthrough
		default:
			sort.Sort(filemanager.SortByFileType(dirs))
		}
	} else {
		sort.Sort(filemanager.SortByFileType(dirs))
	}
	if ctx.Format() == "json" {
		dirList, fileList := s.ListTransfer(dirs)
		data := ctx.Data()
		data.SetData(echo.H{
			`dirList`:  dirList,
			`fileList`: fileList,
		})
		return ctx.JSON(data), true, nil
	}
	return
}

func (s *SftpManager) ListTransfer(dirs []os.FileInfo) (dirList []echo.H, fileList []echo.H) {
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
