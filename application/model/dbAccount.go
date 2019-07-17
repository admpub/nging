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

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/model/base"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func NewDbAccount(ctx echo.Context) *DbAccount {
	return &DbAccount{
		DbAccount: &dbschema.DbAccount{},
		Base:      base.New(ctx),
	}
}

type DbAccount struct {
	*dbschema.DbAccount
	*base.Base
}

func (a *DbAccount) SetOptions() {
	options := echo.H{}
	charset := a.Formx(`charset`).String()
	if len(charset) > 0 {
		options.Set(`charset`, charset)
	}
	b, _ := com.JSONEncode(options)
	a.Options = com.Bytes2str(b)
}

func (a *DbAccount) setDefaultValue() {
	if len(a.DbAccount.Engine) == 0 {
		a.DbAccount.Engine = `mysql`
	}
	switch a.DbAccount.Engine {
	case `mysql`:
		if len(a.DbAccount.Host) == 0 {
			a.DbAccount.Host = `127.0.0.1:3306`
		}
		if len(a.DbAccount.User) == 0 {
			a.DbAccount.User = `root`
		}
	case `redis`:
		if len(a.DbAccount.Host) == 0 {
			a.DbAccount.Host = `127.0.0.1:6379`
		}
		if len(a.DbAccount.Name) == 0 {
			a.DbAccount.Name = `0`
		}
	}
}

func (a *DbAccount) Add() (interface{}, error) {
	if len(a.DbAccount.Title) == 0 {
		return nil, errors.New(a.T(`请输入标题`))
	}
	num, err := a.Count(nil, db.And(db.Cond{`uid`: a.Uid}, db.Cond{`title`: a.Title}))
	if err != nil {
		return nil, err
	}
	if num > 0 {
		return nil, errors.New(a.T(`标题已存在，请设置为一个从未使用过的标题`))
	}
	a.setDefaultValue()
	return a.DbAccount.Add()
}

func (a *DbAccount) Edit(id uint, mw func(db.Result) db.Result, args ...interface{}) error {
	if len(a.DbAccount.Title) == 0 {
		return errors.New(a.T(`请输入标题`))
	}
	num, err := a.Count(nil, db.And(
		db.Cond{`uid`: a.Uid},
		db.Cond{`title`: a.Title},
		db.Cond{`id`: db.NotEq(id)},
	))
	if err != nil {
		return err
	}
	if num > 0 {
		return errors.New(a.T(`标题已存在，请设置为一个从未使用过的标题`))
	}
	a.setDefaultValue()
	return a.DbAccount.Edit(mw, args...)
}
