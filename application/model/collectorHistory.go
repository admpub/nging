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
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v3/application/dbschema"
	"github.com/admpub/nging/v3/application/library/common"
)

func NewCollectorHistory(ctx echo.Context) *CollectorHistory {
	return &CollectorHistory{
		NgingCollectorHistory: dbschema.NewNgingCollectorHistory(ctx),
	}
}

type CollectorHistory struct {
	*dbschema.NgingCollectorHistory
}

func (c *CollectorHistory) Reset() *CollectorHistory {
	c.NgingCollectorHistory.Created = 0
	c.NgingCollectorHistory.Data = ``
	c.NgingCollectorHistory.Content = ``
	c.NgingCollectorHistory.ParentId = 0
	c.NgingCollectorHistory.Exported = 0
	c.NgingCollectorHistory.Id = 0
	c.NgingCollectorHistory.PageId = 0
	c.NgingCollectorHistory.PageParentId = 0
	c.NgingCollectorHistory.PageRootId = 0
	c.NgingCollectorHistory.Url = ``
	c.NgingCollectorHistory.UrlMd5 = ``
	c.NgingCollectorHistory.HasChild = ``
	c.NgingCollectorHistory.RuleMd5 = ``
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
	return this.NgingCollectorHistory.Delete(mw, args...)
}
