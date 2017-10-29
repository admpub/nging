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
	"github.com/admpub/nging/application/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func NewDbAccount(ctx echo.Context) *DbAccount {
	return &DbAccount{
		DbAccount: &dbschema.DbAccount{},
		Base:      &Base{Context: ctx},
	}
}

type DbAccount struct {
	*dbschema.DbAccount
	*Base
}

func (a *DbAccount) setDefaultValue() {
	if len(a.DbAccount.Host) == 0 {
		a.DbAccount.Host = `127.0.0.1:3306`
	}
	if len(a.DbAccount.User) == 0 {
		a.DbAccount.User = `root`
	}
	if len(a.DbAccount.Engine) == 0 {
		a.DbAccount.Host = `mysql`
	}
}

func (a *DbAccount) Add() (interface{}, error) {
	a.setDefaultValue()
	return a.DbAccount.Add()
}

func (a *DbAccount) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	a.setDefaultValue()
	return a.DbAccount.Edit(mw, args...)
}
