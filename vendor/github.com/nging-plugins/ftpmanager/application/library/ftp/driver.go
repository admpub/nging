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
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/nging-plugins/ftpmanager/application/model"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	ftpserver "goftp.io/server/v2"
)

type FileDriver struct {
	ftpserver.Perm
}

func getUserModel(ftpCtx *ftpserver.Context) *model.FtpUser {
	return ftpCtx.Sess.Data[`userModel`].(*model.FtpUser)
}

func (f *FileDriver) realPath(ftpCtx *ftpserver.Context, fpath string, pathType PathType, operate Operate) (string, error) {
	userModel := getUserModel(ftpCtx)
	var allowed bool
	switch operate {
	case OperateCreate:
		allowed = userModel.Allowed(path.Dir(fpath), true)
	case OperateRead:
		allowed = userModel.Allowed(fpath, false)
	case OperateModify:
		allowed = userModel.Allowed(fpath, true)
	}
	if !allowed {
		return fpath, echo.NewError(`permission denied`, code.NonPrivileged)
	}

	rootPath, ok := ftpCtx.Sess.Data[`rootPath`].(string)
	if !ok {
		user := ftpCtx.Sess.LoginUser()
		var err error
		rootPath, err = userModel.GetRootPathOnce(user)
		if err != nil {
			return ``, err
		}
		ftpCtx.Sess.Data[`rootPath`] = rootPath
	}
	return filepath.Join(rootPath, fpath), nil
}

func (f *FileDriver) ChangeDir(ftpCtx *ftpserver.Context, path string) error {
	rPath, err := f.realPath(ftpCtx, path, PathTypeDir, OperateRead)
	if err != nil {
		return err
	}
	fi, err := os.Lstat(rPath)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return nil
	}
	return ErrNotDirectory
}

func (f *FileDriver) Stat(ftpCtx *ftpserver.Context, path string) (os.FileInfo, error) {
	basepath, err := f.realPath(ftpCtx, path, PathTypeBoth, OperateRead)
	if err != nil {
		return nil, err
	}
	rPath, err := filepath.Abs(basepath)
	if err != nil {
		return nil, err
	}
	return os.Lstat(rPath)
}

func (f *FileDriver) ListDir(ftpCtx *ftpserver.Context, path string, callback func(os.FileInfo) error) error {
	basepath, err := f.realPath(ftpCtx, path, PathTypeDir, OperateRead)
	if err != nil {
		return err
	}
	userModel := getUserModel(ftpCtx)
	err = filepath.Walk(basepath, func(f string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rPath, _ := filepath.Rel(basepath, f)
		if !userModel.Allowed(rPath, false) {
			return err
		}
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

	return err
}

func (f *FileDriver) DeleteDir(ftpCtx *ftpserver.Context, path string) error {
	rPath, err := f.realPath(ftpCtx, path, PathTypeDir, OperateModify)
	if err != nil {
		return err
	}
	fi, err := os.Lstat(rPath)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return os.Remove(rPath)
	}
	return ErrNotDirectory
}

func (f *FileDriver) DeleteFile(ftpCtx *ftpserver.Context, path string) error {
	rPath, err := f.realPath(ftpCtx, path, PathTypeFile, OperateModify)
	if err != nil {
		return err
	}
	fi, err := os.Lstat(rPath)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return os.Remove(rPath)
	}
	return ErrNotFile
}

func (f *FileDriver) Rename(ftpCtx *ftpserver.Context, fromPath string, toPath string) error {
	oldPath, err := f.realPath(ftpCtx, fromPath, PathTypeBoth, OperateModify)
	if err != nil {
		return err
	}
	fi, err := os.Lstat(oldPath)
	if err != nil {
		return err
	}
	var pt PathType
	if fi.IsDir() {
		pt = PathTypeDir
	} else {
		pt = PathTypeFile
	}
	newPath, err := f.realPath(ftpCtx, toPath, pt, OperateCreate)
	if err != nil {
		return err
	}
	return com.Rename(oldPath, newPath)
}

func (f *FileDriver) MakeDir(ftpCtx *ftpserver.Context, path string) error {
	rPath, err := f.realPath(ftpCtx, path, PathTypeDir, OperateCreate)
	if err != nil {
		return err
	}
	return os.Mkdir(rPath, os.ModePerm)
}

func (f *FileDriver) GetFile(ftpCtx *ftpserver.Context, path string, offset int64) (int64, io.ReadCloser, error) {
	rPath, err := f.realPath(ftpCtx, path, PathTypeFile, OperateRead)
	if err != nil {
		return 0, nil, err
	}
	fp, err := os.Open(rPath)
	if err != nil {
		return 0, nil, err
	}

	info, err := fp.Stat()
	if err != nil {
		return 0, nil, err
	}

	fp.Seek(offset, io.SeekStart)

	return info.Size(), fp, nil
}

func (f *FileDriver) PutFile(ftpCtx *ftpserver.Context, destPath string, data io.Reader, offset int64) (int64, error) {
	rPath, err := f.realPath(ftpCtx, destPath, PathTypeFile, OperateCreate)
	if err != nil {
		return 0, err
	}
	var isExist bool
	fi, err := os.Lstat(rPath)
	if err == nil {
		isExist = true
		if fi.IsDir() {
			return 0, ErrDirectoryAlreadyExists
		}
	} else {
		if os.IsNotExist(err) {
			isExist = false
		} else {
			return 0, fmt.Errorf(`%w: %v`, ErrPutFile, err)
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

	_, err = of.Seek(offset, io.SeekStart)
	if err != nil {
		return 0, err
	}

	bytes, err := io.Copy(of, data)
	if err != nil {
		return 0, err
	}

	return bytes, nil
}
