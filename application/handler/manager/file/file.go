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

// Package file 上传文件管理
package file

import (
	"time"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/model/file"
	"github.com/coscms/webcore/registry/upload"
	"github.com/coscms/webcore/registry/upload/checker"
)

func setUploadURL(ctx echo.Context) error {
	subdir := ctx.Form(`subdir`, `default`)
	if !upload.AllowedSubdir(subdir) {
		return ctx.NewError(code.InvalidParameter, `无效的subdir值`).SetZone(`subdir`)
	}
	ctx.Set(`subdir`, subdir)
	ctx.Set(`uploadURL`, checker.BackendUploadURL(subdir))
	return nil
}

func FileList(ctx echo.Context) error {
	return FileListWithOwner(ctx, ``, 0)
}

func FileDelete(ctx echo.Context) (err error) {
	user := backend.User(ctx)
	ownerID := uint64(user.Id)
	return FileDeleteWithOwner(ctx, `user`, ownerID)
}

// FileClean 删除未使用文件
func FileClean(ctx echo.Context) (err error) {
	fileM := file.NewFile(ctx)
	ago := ctx.Form(`ago`)
	var seconds int64 = 86400 * 365
	if len(ago) > 0 {
		t, e := time.ParseDuration(ago)
		if e != nil {
			return e
		}
		seconds = int64(t.Seconds())
	}
	err = fileM.RemoveUnused(seconds, ``, 0)
	if err != nil {
		return err
	}

	return ctx.Redirect(backend.URLFor(`/manager/file/list`))
}
