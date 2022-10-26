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

package handler

import (
	"encoding/json"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"

	"github.com/nging-plugins/collector/application/model"
)

func History(c echo.Context) error {
	pageID := c.Queryx(`pageId`).Uint()
	m := model.NewCollectorHistory(c)
	var cond db.Compound
	if pageID > 0 {
		cond = db.Cond{`page_parent_id`: pageID}
	} else {
		cond = db.Cond{`page_parent_id`: 0}
	}
	_, err := handler.PagingWithLister(c, handler.NewLister(m, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond))
	ret := handler.Err(c, err)
	c.Set(`listData`, m.Objects())
	return c.Render(`collector/history`, ret)
}

func HistoryView(c echo.Context) error {
	ident := c.Form(`ident`)
	data := c.Data()
	if !com.IsAlphaNumericUnderscore(ident) {
		return c.JSON(data.SetInfo(c.T(`无效参数`), 0))
	}
	b, err := common.ReadCache(`colloctor`, ident+`.json`)
	if err != nil {
		return c.JSON(data.SetError(err))
	}
	data.SetData(json.RawMessage(b))
	return c.JSON(data)
}

func HistoryDelete(c echo.Context) error {
	id := c.Form(`id`)
	m := model.NewCollectorHistory(c)
	cond := db.Cond{`id`: id}
	err := m.Delete(nil, cond)
	if err != nil {
		handler.SendFail(c, err.Error())
	} else {
		handler.SendOk(c, c.T(`操作成功`))
	}
	return c.Redirect(handler.URLFor(`/collector/history`))
}
