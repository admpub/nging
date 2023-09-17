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

package model

import (
	"encoding/json"
	"errors"
	"sort"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v5/application/library/common"

	"github.com/nging-plugins/ftpmanager/application/dbschema"
	"github.com/nging-plugins/ftpmanager/application/library/fileperm"
)

type FtpUserAndGroup struct {
	*dbschema.NgingFtpUser
	Group *dbschema.NgingFtpUserGroup `db:"-,relation=id:group_id|gtZero"`
}

func NewFtpUser(ctx echo.Context) *FtpUser {
	return &FtpUser{
		NgingFtpUser: dbschema.NewNgingFtpUser(ctx),
	}
}

type FtpUser struct {
	*dbschema.NgingFtpUser
	rootPath           string
	groupDefaultModify bool
	groupPermission    fileperm.User
	userPermission     fileperm.User
}

func (f *FtpUser) Exists(username string, excludeIDs ...uint) (bool, error) {
	cond := db.NewCompounds()
	cond.AddKV("username", username)
	if len(excludeIDs) > 0 {
		cond.AddKV("id", db.NotEq(excludeIDs[0]))
	}
	return f.NgingFtpUser.Exists(nil, cond.And())
}

func (f *FtpUser) ExistsAvailable() (bool, error) {
	cond := db.NewCompounds()
	cond.AddKV("banned", "N")
	return f.NgingFtpUser.Exists(nil, cond.And())
}

func (f *FtpUser) CheckPasswd(username string, password string) (bool, error) {
	salt := common.CookieConfig().BlockKey
	err := f.NgingFtpUser.Get(nil, db.Cond{`username`: username})
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = echo.NewError(`Account does not exist`, code.UserNotFound)
		}
		return false, err
	}
	if f.Password != com.MakePassword(password, salt) {
		return false, echo.NewError(`Incorrect password`, code.Unauthenticated)
	}

	// 获取根路径
	f.rootPath, err = f.RootPath(username)
	if err != nil {
		return false, err
	}

	// 获取组权限
	var groupPermission fileperm.Rules
	if f.NgingFtpUser.GroupId > 0 {
		groupPermM := NewFtpPermission(f.Context())
		err = groupPermM.GetByTarget(`group`, f.NgingFtpUser.GroupId)
		if err != nil {
			if err != db.ErrNoMoreRows {
				return false, err
			}
			err = nil
		} else if len(groupPermM.NgingFtpPermission.Permission) > 0 {
			jsonBytes := []byte(groupPermM.NgingFtpPermission.Permission)
			err = json.Unmarshal(jsonBytes, &groupPermission)
			if err != nil {
				return false, common.JSONBytesParseError(err, jsonBytes)
			}
		}
	}
	err = groupPermission.Init()
	if err != nil {
		return false, err
	}
	sort.Sort(groupPermission)
	f.groupPermission = fileperm.User{
		Writeable: f.groupDefaultModify,
		RootDir:   f.rootPath,
		Rules:     groupPermission,
	}

	// 获取用户权限
	var userPermission fileperm.Rules
	userPermM := NewFtpPermission(f.Context())
	err = userPermM.GetByTarget(`user`, f.NgingFtpUser.Id)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return false, err
		}
		err = nil
	} else if len(userPermM.NgingFtpPermission.Permission) > 0 {
		jsonBytes := []byte(userPermM.NgingFtpPermission.Permission)
		err = json.Unmarshal(jsonBytes, &userPermission)
		if err != nil {
			return false, common.JSONBytesParseError(err, jsonBytes)
		}
	}
	err = userPermission.Init()
	sort.Sort(userPermission)
	f.userPermission = fileperm.User{
		Writeable: f.NgingFtpUser.Modify == `Y`,
		RootDir:   f.rootPath,
		Rules:     userPermission,
	}

	return true, err
}

var (
	ErrNoneDirectory      = errors.New(`No accessible directory`)
	ErrBannedUser         = errors.New(`The current account has been disabled`)
	ErrIPAddressIsBlocked = errors.New(`IP is blocked`)
)

func (f *FtpUser) GetRootPathOnce(username string) (string, error) {
	if len(f.rootPath) > 0 {
		return f.rootPath, nil
	}
	rootPath, err := f.RootPath(username)
	f.rootPath = rootPath
	return rootPath, err
}

func (f *FtpUser) Allowed(path string, modification bool) bool {
	return f.userPermission.Allowed(path, modification) && f.groupPermission.Allowed(path, modification)
}

func (f *FtpUser) RootPath(username string) (basePath string, err error) {
	if f.Id <= 0 {
		err = f.Get(nil, db.Cond{`username`: username})
		if err != nil {
			return
		}
	}
	if f.NgingFtpUser.GroupId > 0 {
		m := NewFtpUserGroup(f.Context())
		err = m.Get(nil, db.Cond{`id`: f.NgingFtpUser.GroupId})
		if err != nil {
			return
		}
		if m.NgingFtpUserGroup.Banned == `Y` {
			err = ErrBannedUser
			return
		}
		basePath = m.NgingFtpUserGroup.Directory
		f.groupDefaultModify = m.NgingFtpUserGroup.Modify == `Y`
	} else {
		f.groupDefaultModify = f.Modify == `Y`
	}
	if f.NgingFtpUser.Banned == `Y` {
		err = ErrBannedUser
		return
	}
	if len(f.NgingFtpUser.Directory) > 0 {
		basePath = f.NgingFtpUser.Directory
		return
	}
	if len(basePath) < 1 {
		err = ErrNoneDirectory
	}
	return
}

func (f *FtpUser) check() error {
	if len(f.Username) == 0 {
		return f.Context().NewError(code.InvalidParameter, `用户名不能为空`).SetZone(`username`)
	}
	var exists bool
	var err error
	if f.Id > 0 {
		exists, err = f.Exists(f.Username, f.Id)
	} else {
		exists, err = f.Exists(f.Username)
	}
	if err != nil {
		return err
	}
	if exists {
		return f.Context().NewError(code.DataAlreadyExists, `用户名已经存在`).SetZone(`username`)
	}
	return err
}

func (f *FtpUser) Add() (pk interface{}, err error) {
	if err := f.check(); err != nil {
		return nil, err
	}
	salt := common.CookieConfig().BlockKey
	f.Password = com.MakePassword(f.Password, salt)
	return f.NgingFtpUser.Insert()
}

func (f *FtpUser) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := f.check(); err != nil {
		return err
	}
	old := dbschema.NewNgingFtpUser(f.Context())
	err := old.Get(func(r db.Result) db.Result {
		return r.Select(`password`)
	}, `id`, f.Id)
	if err != nil {
		return err
	}
	if len(f.Password) == 0 {
		f.Password = old.Password
	} else {
		salt := common.CookieConfig().BlockKey
		f.Password = com.MakePassword(f.Password, salt)
	}
	return f.NgingFtpUser.Update(mw, args...)
}
