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
	"github.com/admpub/nging/v3/application/library/collector/exec"
	"github.com/admpub/nging/v3/application/model/base"
)

func NewCollectorPage(ctx echo.Context) *CollectorPage {
	return &CollectorPage{
		NgingCollectorPage: &dbschema.NgingCollectorPage{},
		Base:               base.New(ctx),
	}
}

type CollectorPage struct {
	*dbschema.NgingCollectorPage
	*base.Base
}

func (c *CollectorPage) FullData() (*exec.Rules, error) {
	data := exec.NewRules()
	data.NgingCollectorPage = c.NgingCollectorPage
	ruleM := NewCollectorRule(c.Base.Context)
	id := c.Id
	_, err := ruleM.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`sort`, `id`)
	}, 0, -1, db.Cond{`page_id`: id})
	if err != nil {
		return data, err
	}
	data.RuleList = ruleM.Objects()
	_, err = c.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`sort`, `id`)
	}, 0, -1, db.And(
		db.Cond{`root_id`: id},
		db.Cond{`parent_id`: db.Gt(0)},
	))
	if err != nil {
		return data, err
	}
	for _, pageRow := range c.Objects() {
		pageForm := &exec.Rule{NgingCollectorPage: pageRow}

		_, err = ruleM.ListByOffset(nil, func(r db.Result) db.Result {
			return r.OrderBy(`sort`, `id`)
		}, 0, -1, db.Cond{`page_id`: pageRow.Id})
		if err != nil {
			return data, err
		}
		pageForm.RuleList = ruleM.Objects()
		data.Extra = append(data.Extra, pageForm)
	}
	return data, nil
}
