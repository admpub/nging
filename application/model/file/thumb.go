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

package file

import (
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/model/base"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func NewThumb(ctx echo.Context) *Thumb {
	return &Thumb{
		FileThumb: &dbschema.FileThumb{},
		base:      base.New(ctx),
	}
}

type Thumb struct {
	*dbschema.FileThumb
	base *base.Base
}

func (t *Thumb) SetByFile(file *dbschema.File) *Thumb {
	t.FileId = file.Id
	t.Dpi = file.Dpi
	return t
}

func (t *Thumb) Save() (err error) {
	m := &dbschema.FileThumb{}
	err = m.Get(nil, db.And(
		db.Cond{`save_path`: t.SavePath},
	))
	if err != nil {
		if err != db.ErrNoMoreRows {
			return
		}
		_, err = t.FileThumb.Add()
		return
	}
	t.FileThumb = m
	err = t.SetFields(nil, echo.H{
		`view_url`: t.ViewUrl,
		`size`:     t.Size,
		`width`:    t.Width,
		`height`:   t.Height,
		`dpi`:      t.Dpi,
	}, db.Cond{`id`: m.Id})
	return
}
