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
	"errors"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/model/base"
)

type FtpUserAndGroup struct {
	*dbschema.NgingFtpUser
	Group *dbschema.NgingFtpUserGroup
}

var DefaultSalt = ``

func NewFtpUser(ctx echo.Context) *FtpUser {
	return &FtpUser{
		NgingFtpUser: &dbschema.NgingFtpUser{},
		Base:         base.New(ctx),
	}
}

type FtpUser struct {
	*dbschema.NgingFtpUser
	*base.Base
}

func (f *FtpUser) Exists(username string) (bool, error) {
	n, e := f.Param(nil, db.Cond{`username`: username}).Count()
	return n > 0, e
}

func (f *FtpUser) CheckPasswd(username string, password string) (bool, error) {
	n, e := f.Param(nil, db.Cond{
		`username`: username,
		`password`: com.MakePassword(password, DefaultSalt),
	}).Count()
	y := n > 0
	if y {
		_, e = f.RootPath(username)
		if e != nil {
			y = false
		}
	}
	return y, e
}

var (
	ErrNoneFtpDirectory = errors.New(`No accessible directory`)
	ErrBannedFtpUser    = errors.New(`The current account has been disabled`)
)

func (f *FtpUser) RootPath(username string) (basePath string, err error) {
	err = f.Get(nil, db.Cond{`username`: username})
	if err != nil {
		return
	}
	if f.NgingFtpUser.GroupId > 0 {
		m := NewFtpUserGroup(f.Base.Context)
		err = m.Get(nil, db.Cond{`id`: f.NgingFtpUser.GroupId})
		if err != nil {
			return
		}
		if m.NgingFtpUserGroup.Banned == `Y` {
			err = ErrBannedFtpUser
			return
		}
		basePath = m.NgingFtpUserGroup.Directory
	}
	if f.NgingFtpUser.Banned == `Y` {
		err = ErrBannedFtpUser
		return
	}
	if len(f.NgingFtpUser.Directory) > 0 {
		basePath = f.NgingFtpUser.Directory
		return
	}
	if len(basePath) < 1 {
		err = ErrNoneFtpDirectory
	}
	return
}
