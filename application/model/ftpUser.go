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
package model

import (
	"errors"

	"github.com/admpub/nging/application/dbschema"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

type FtpUserAndGroup struct {
	*dbschema.FtpUser
	Group *dbschema.FtpUserGroup
}

var DefaultSalt = ``

func NewFtpUser(ctx echo.Context) *FtpUser {
	return &FtpUser{
		FtpUser: &dbschema.FtpUser{},
		Base:    &Base{Context: ctx},
	}
}

type FtpUser struct {
	*dbschema.FtpUser
	*Base
}

func (f *FtpUser) Exists(username string) (bool, error) {
	n, e := f.Param().SetArgs(db.Cond{`username`: username}).Count()
	return n > 0, e
}

func (f *FtpUser) CheckPasswd(username string, password string) (bool, error) {
	n, e := f.Param().SetArgs(db.Cond{`username`: username, `password`: com.MakePassword(password, DefaultSalt)}).Count()
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
	if f.FtpUser.GroupId > 0 {
		m := NewFtpUserGroup(f.Base.Context)
		err = m.Get(nil, db.Cond{`id`: f.FtpUser.GroupId})
		if err != nil {
			return
		}
		if m.FtpUserGroup.Banned == `Y` {
			err = ErrBannedFtpUser
			return
		}
		basePath = m.FtpUserGroup.Directory
	}
	if f.FtpUser.Banned == `Y` {
		err = ErrBannedFtpUser
		return
	}
	if len(f.FtpUser.Directory) > 0 {
		basePath = f.FtpUser.Directory
		return
	}
	if len(basePath) < 1 {
		err = ErrNoneFtpDirectory
	}
	return
}
