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
	"github.com/admpub/caddyui/application/dbschema"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

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
	return n > 0, e
}
