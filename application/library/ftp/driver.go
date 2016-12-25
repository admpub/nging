/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package ftp

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	ftpserver "github.com/admpub/ftpserver"
	"github.com/admpub/nging/application/model"
)

type FileDriver struct {
	conn *ftpserver.Conn
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

func (driver *FileDriver) realPath(path string) (string, error) {
	user := driver.conn.LoginUser()
	rootPath, err := driver.user.RootPath(user)
	if err != nil {
		return ``, err
	}
	paths := strings.Split(path, "/")
	return filepath.Join(append([]string{rootPath}, paths...)...), nil
}

func (driver *FileDriver) Init(conn *ftpserver.Conn) {
	driver.conn = conn
}

func (driver *FileDriver) ChangeDir(path string) error {
	rPath, err := driver.realPath(path)
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

func (driver *FileDriver) Stat(path string) (ftpserver.FileInfo, error) {
	basepath, err := driver.realPath(path)
	if err != nil {
		return nil, err
	}
	rPath, err := filepath.Abs(basepath)
	if err != nil {
		return nil, err
	}
	f, err := os.Lstat(rPath)
	if err != nil {
		return nil, err
	}
	mode, err := driver.Perm.GetMode(path)
	if err != nil {
		return nil, err
	}
	if mode == 0 {
		mode = f.Mode()
	}
	if f.IsDir() {
		mode |= os.ModeDir
	}
	owner, err := driver.Perm.GetOwner(path)
	if err != nil {
		return nil, err
	}
	group, err := driver.Perm.GetGroup(path)
	if err != nil {
		return nil, err
	}
	return &FileInfo{f, mode, owner, group}, nil
}

func (driver *FileDriver) ListDir(path string, callback func(ftpserver.FileInfo) error) error {
	basepath, err := driver.realPath(path)
	if err != nil {
		return err
	}
	filepath.Walk(basepath, func(f string, info os.FileInfo, err error) error {
		rPath, _ := filepath.Rel(basepath, f)
		if rPath == info.Name() {
			mode, err := driver.Perm.GetMode(rPath)
			if err != nil {
				return err
			}
			if mode == 0 {
				mode = info.Mode()
			}
			if info.IsDir() {
				mode |= os.ModeDir
			}
			owner, err := driver.Perm.GetOwner(rPath)
			if err != nil {
				return err
			}
			group, err := driver.Perm.GetGroup(rPath)
			if err != nil {
				return err
			}
			err = callback(&FileInfo{info, mode, owner, group})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return nil
}

func (driver *FileDriver) DeleteDir(path string) error {
	rPath, err := driver.realPath(path)
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

func (driver *FileDriver) DeleteFile(path string) error {
	rPath, err := driver.realPath(path)
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

func (driver *FileDriver) Rename(fromPath string, toPath string) error {
	oldPath, err := driver.realPath(fromPath)
	if err != nil {
		return err
	}
	newPath, err := driver.realPath(toPath)
	if err != nil {
		return err
	}
	return os.Rename(oldPath, newPath)
}

func (driver *FileDriver) MakeDir(path string) error {
	rPath, err := driver.realPath(path)
	if err != nil {
		return err
	}
	return os.Mkdir(rPath, os.ModePerm)
}

func (driver *FileDriver) GetFile(path string, offset int64) (int64, io.ReadCloser, error) {
	rPath, err := driver.realPath(path)
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

func (driver *FileDriver) PutFile(destPath string, data io.Reader, appendData bool) (int64, error) {
	rPath, err := driver.realPath(destPath)
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

	if appendData && !isExist {
		appendData = false
	}

	if !appendData {
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

	_, err = of.Seek(0, os.SEEK_END)
	if err != nil {
		return 0, err
	}

	bytes, err := io.Copy(of, data)
	if err != nil {
		return 0, err
	}

	return bytes, nil
}

type FileDriverFactory struct {
	ftpserver.Perm
}

func (factory *FileDriverFactory) NewDriver() (ftpserver.Driver, error) {
	return &FileDriver{nil, model.NewFtpUser(nil), factory.Perm}, nil
}
