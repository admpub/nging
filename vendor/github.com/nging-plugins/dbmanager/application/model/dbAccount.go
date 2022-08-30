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

	"github.com/nging-plugins/dbmanager/application/dbschema"
	"github.com/nging-plugins/dbmanager/application/library/dbmanager/driver/mysql"
)

func NewDbAccount(ctx echo.Context) *DbAccount {
	return &DbAccount{
		NgingDbAccount: dbschema.NewNgingDbAccount(ctx),
	}
}

type DbAccount struct {
	*dbschema.NgingDbAccount
}

func (a *DbAccount) SetOptions() error {
	options := echo.H{}
	charset := a.Context().Formx(`charset`).String()
	if len(charset) > 0 {
		if !com.InSlice(charset, mysql.Charsets) {
			return a.Context().E(`字符集charset值无效`)
		}
		options.Set(`charset`, charset)
	}
	b, _ := com.JSONEncode(options)
	a.Options = com.Bytes2str(b)
	return nil
}

func (a *DbAccount) setDefaultValue() {
	if len(a.NgingDbAccount.Engine) == 0 {
		a.NgingDbAccount.Engine = `mysql`
	}
	switch a.NgingDbAccount.Engine {
	case `mysql`:
		if len(a.NgingDbAccount.Host) == 0 {
			a.NgingDbAccount.Host = `127.0.0.1:3306`
		}
		if len(a.NgingDbAccount.User) == 0 {
			a.NgingDbAccount.User = `root`
		}
	case `redis`:
		if len(a.NgingDbAccount.Host) == 0 {
			a.NgingDbAccount.Host = `127.0.0.1:6379`
		}
		a.NgingDbAccount.Name = ``
		/*
			if len(a.NgingDbAccount.Name) == 0 {
				a.NgingDbAccount.Name = `0`
			}
		*/
	}
}

func (a *DbAccount) Add() (interface{}, error) {
	if len(a.NgingDbAccount.Title) == 0 {
		return nil, errors.New(a.Context().T(`请输入标题`))
	}
	num, err := a.Count(nil, db.And(db.Cond{`uid`: a.Uid}, db.Cond{`title`: a.Title}))
	if err != nil {
		return nil, err
	}
	if num > 0 {
		return nil, errors.New(a.Context().T(`标题已存在，请设置为一个从未使用过的标题`))
	}
	a.setDefaultValue()
	return a.NgingDbAccount.Insert()
}

func (a *DbAccount) Edit(id uint, mw func(db.Result) db.Result, args ...interface{}) error {
	if len(a.NgingDbAccount.Title) == 0 {
		return errors.New(a.Context().T(`请输入标题`))
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
		return errors.New(a.Context().T(`标题已存在，请设置为一个从未使用过的标题`))
	}
	a.setDefaultValue()
	return a.NgingDbAccount.Update(mw, args...)
}
