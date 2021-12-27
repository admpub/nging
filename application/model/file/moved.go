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

package file

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/dbschema"
)

func NewMoved(ctx echo.Context) *Moved {
	m := &Moved{
		NgingFileMoved: dbschema.NewNgingFileMoved(ctx),
	}
	return m
}

type Moved struct {
	*dbschema.NgingFileMoved
}

func (t *Moved) ListByViewURLs(viewURLs []interface{}) ([]*dbschema.NgingFileMoved, error) {
	_, err := t.ListByOffset(nil, func(r db.Result) db.Result {
		return r.Select(`id`, `from`, `to`)
	}, 0, -1, db.Cond{
		`from`: db.In(viewURLs),
	})
	if err != nil {
		return nil, err
	}
	return t.Objects(), err
}

func (t *Moved) Add() (err error) {
	_, err = t.NgingFileMoved.Add()
	return
}

func (t *Moved) Save() (err error) {
	m := dbschema.NewNgingFileMoved(t.Context())
	err = m.Get(nil, db.And(
		db.Cond{`from`: t.From},
	))
	if err != nil {
		if err != db.ErrNoMoreRows {
			return
		}
		_, err = t.NgingFileMoved.Add()
		return
	}
	t.NgingFileMoved = m
	err = t.SetFields(nil, echo.H{
		`to`: t.To,
	}, db.Cond{`id`: m.Id})
	return
}
