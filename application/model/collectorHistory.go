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
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/model/base"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func NewCollectorHistory(ctx echo.Context) *CollectorHistory {
	return &CollectorHistory{
		CollectorHistory: &dbschema.CollectorHistory{},
		Base:             base.New(ctx),
	}
}

type CollectorHistory struct {
	*dbschema.CollectorHistory
	*base.Base
}

func (c *CollectorHistory) Reset() *CollectorHistory {
	c.CollectorHistory.Created = 0
	c.CollectorHistory.Data = ``
	c.CollectorHistory.Content = ``
	c.CollectorHistory.ParentId = 0
	c.CollectorHistory.Exported = 0
	c.CollectorHistory.Id = 0
	c.CollectorHistory.PageId = 0
	c.CollectorHistory.PageParentId = 0
	c.CollectorHistory.PageRootId = 0
	c.CollectorHistory.Url = ``
	c.CollectorHistory.UrlMd5 = ``
	c.CollectorHistory.HasChild = ``
	c.CollectorHistory.RuleMd5 = ``
	return c
}

func (this *CollectorHistory) Delete(mw func(db.Result) db.Result, args ...interface{}) error {
	err := this.Get(mw, args...)
	if err != nil {
		return err
	}
	err = common.RemoveCache(`colloctor`, this.UrlMd5+`.json`)
	if err != nil {
		return err
	}
	return this.CollectorHistory.Delete(mw, args...)
}
