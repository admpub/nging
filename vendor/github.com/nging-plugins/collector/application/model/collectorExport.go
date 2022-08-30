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
	"github.com/webx-top/echo/engine"

	"github.com/nging-plugins/collector/application/dbschema"
	"github.com/nging-plugins/collector/application/library/collector/export"
)

func NewCollectorExport(ctx echo.Context) *CollectorExport {
	return &CollectorExport{
		NgingCollectorExport: dbschema.NewNgingCollectorExport(ctx),
	}
}

type CollectorExport struct {
	*dbschema.NgingCollectorExport
}

func (c *CollectorExport) Add() (pk interface{}, err error) {
	err = c.check()
	if err != nil {
		return
	}
	return c.NgingCollectorExport.Insert()
}

func (c *CollectorExport) check() error {
	var err error
	if c.PageRoot < 1 {
		err = errors.New(c.Context().T(`请选择页面规则`))
	} else if len(c.Dest) == 0 {
		err = errors.New(c.Context().T(`请设置导出到哪儿`))
	} else if c.PageId < 1 {
		c.PageId = c.PageRoot
	}
	return err
}

func (c *CollectorExport) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	err := c.check()
	if err != nil {
		return err
	}
	return c.NgingCollectorExport.Update(mw, args...)
}

func (c *CollectorExport) Export() (int64, error) {
	rows := []*dbschema.NgingCollectorExport{}
	cnt, err := c.ListByOffset(&rows, nil, 0, -1, db.Cond{`disabled`: `N`})
	if err != nil {
		return 0, err
	}
	for _, row := range rows {
		if len(row.Mapping) == 0 {
			continue
		}
		mapping := export.NewMappings()
		err = com.JSONDecode(engine.Str2bytes(row.Mapping), mapping)
		if err != nil {
			return 0, err
		}
	}
	return cnt(), err
}
