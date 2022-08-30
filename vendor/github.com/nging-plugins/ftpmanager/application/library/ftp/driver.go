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

package ftp

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nging-plugins/ftpmanager/application/model"
	"github.com/webx-top/com"
	ftpserver "goftp.io/server/v2"
)

type FileDriver struct {
	user *model.FtpUser
	ftpserver.Perm
}

type FileInfo struct {
	os.FileInfo

	mode  os.FileMode
	owner string
	group string
}

func (f *FileInfo) Mode() os.FileMode {
	return f.mode
}

func (f *FileInfo) Owner() string {
	return f.owner
}

func (f *FileInfo) Group() string {
	return f.group
}

func (driver *FileDriver) realPath(ftpCtx *ftpserver.Context, path string) (string, error) {
	user := ftpCtx.Sess.LoginUser()
	rootPath, err := driver.user.RootPath(user)
	if err != nil {
		return ``, err
	}
	paths := strings.Split(path, "/")
	return filepath.Join(append([]string{rootPath}, paths...)...), nil
}

func (driver *FileDriver) ChangeDir(ftpCtx *ftpserver.Context, path string) error {
	rPath, err := driver.realPath(ftpCtx, path)
	if err != nil {
		return err
	}
	f, err := os.Lstat(rPath)
	if err != nil {
		return err
	}
	if f.IsDir() {
		return nil
	}
	return errors.New("Not a directory")
}

func (driver *FileDriver) Stat(ftpCtx *ftpserver.Context, path string) (os.FileInfo, error) {
	basepath, err := driver.realPath(ftpCtx, path)
	if err != nil {
		return nil, err
	}
	rPath, err := filepath.Abs(basepath)
	if err != nil {
		return nil, err
	}
	return os.Lstat(rPath)
}

func (driver *FileDriver) ListDir(ftpCtx *ftpserver.Context, path string, callback func(os.FileInfo) error) error {
	basepath, err := driver.realPath(ftpCtx, path)
	if err != nil {
		return err
	}
	filepath.Walk(basepath, func(f string, info os.FileInfo, err error) error {
		rPath, _ := filepath.Rel(basepath, f)
		if rPath == info.Name() {
			err = callback(info)
			if err != nil {
				return err
			}
			if info.IsDir() {
				return filepath.SkipDir
			}
		}
		return nil
	})

	return nil
}

func (driver *FileDriver) DeleteDir(ftpCtx *ftpserver.Context, path string) error {
	rPath, err := driver.realPath(ftpCtx, path)
	if err != nil {
		return err
	}
	f, err := os.Lstat(rPath)
	if err != nil {
		return err
	}
	if f.IsDir() {
		return os.Remove(rPath)
	}
	return errors.New("Not a directory")
}

func (driver *FileDriver) DeleteFile(ftpCtx *ftpserver.Context, path string) error {
	rPath, err := driver.realPath(ftpCtx, path)
	if err != nil {
		return err
	}
	f, err := os.Lstat(rPath)
	if err != nil {
		return err
	}
	if !f.IsDir() {
		return os.Remove(rPath)
	}
	return errors.New("Not a file")
}

func (driver *FileDriver) Rename(ftpCtx *ftpserver.Context, fromPath string, toPath string) error {
	oldPath, err := driver.realPath(ftpCtx, fromPath)
	if err != nil {
		return err
	}
	newPath, err := driver.realPath(ftpCtx, toPath)
	if err != nil {
		return err
	}
	return com.Rename(oldPath, newPath)
}

func (driver *FileDriver) MakeDir(ftpCtx *ftpserver.Context, path string) error {
	rPath, err := driver.realPath(ftpCtx, path)
	if err != nil {
		return err
	}
	return os.Mkdir(rPath, os.ModePerm)
}

func (driver *FileDriver) GetFile(ftpCtx *ftpserver.Context, path string, offset int64) (int64, io.ReadCloser, error) {
	rPath, err := driver.realPath(ftpCtx, path)
	if err != nil {
		return 0, nil, err
	}
	f, err := os.Open(rPath)
	if err != nil {
		return 0, nil, err
	}

	info, err := f.Stat()
	if err != nil {
		return 0, nil, err
	}

	f.Seek(offset, os.SEEK_SET)

	return info.Size(), f, nil
}

func (driver *FileDriver) PutFile(ftpCtx *ftpserver.Context, destPath string, data io.Reader, offset int64) (int64, error) {
	rPath, err := driver.realPath(ftpCtx, destPath)
	if err != nil {
		return 0, err
	}
	var isExist bool
	f, err := os.Lstat(rPath)
	if err == nil {
		isExist = true
		if f.IsDir() {
			return 0, errors.New("A dir has the same name")
		}
	} else {
		if os.IsNotExist(err) {
			isExist = false
		} else {
			return 0, errors.New(fmt.Sprintln("Put File error:", err))
		}
	}

	if offset > -1 && !isExist {
		offset = -1
	}

	if offset == -1 {
		if isExist {
			err = os.Remove(rPath)
			if err != nil {
				return 0, err
			}
		}
		f, err := os.Create(rPath)
		if err != nil {
			return 0, err
		}
		defer f.Close()
		bytes, err := io.Copy(f, data)
		if err != nil {
			return 0, err
		}
		return bytes, nil
	}

	of, err := os.OpenFile(rPath, os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		return 0, err
	}
	defer of.Close()

	info, err := of.Stat()
	if err != nil {
		return 0, err
	}
	if offset > info.Size() {
		return 0, fmt.Errorf("offset %d is beyond file size %d", offset, info.Size())
	}

	_, err = of.Seek(offset, os.SEEK_END)
	if err != nil {
		return 0, err
	}

	bytes, err := io.Copy(of, data)
	if err != nil {
		return 0, err
	}

	return bytes, nil
}
