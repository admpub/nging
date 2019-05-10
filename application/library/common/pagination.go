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

package common

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
	"github.com/webx-top/pagination"
)

var PageMaxSize = 1000

func Paging(ctx echo.Context) (page int, size int) {
	page = ctx.Formx(`page`, ctx.Form(`pageNumber`)).Int()
	size = ctx.Formx(`size`, ctx.Form(`pageSize`)).Int()
	if page < 1 {
		page = 1
	}
	if size < 1 || size > PageMaxSize {
		size = 50
	}
	return
}

func PagingWithPagination(ctx echo.Context, delKeys ...string) (page int, size int, totalRows int, p *pagination.Pagination) {
	page, size = Paging(ctx)
	totalRows = ctx.Formx(`rows`).Int()
	pjax := ctx.PjaxContainer()
	if len(pjax) > 0 {
		delKeys = append(delKeys, `_pjax`)
	}
	p = pagination.New(ctx).SetAll(``, totalRows, page, 10, size).SetURL(map[string]string{
		`rows`: `rows`,
		`page`: `page`,
		`size`: `size`,
	}, delKeys...)
	return
}

func PagingWithLister(ctx echo.Context, m Lister, varSuffix ...string) (*pagination.Pagination, error) {
	page, size, totalRows, p := PagingWithPagination(ctx)
	cnt, err := m.List(nil, nil, page, size)
	if totalRows <= 0 {
		totalRows = int(cnt())
		p.SetRows(totalRows)
	}
	if len(varSuffix) > 0 {
		ctx.Set(`pagination`+varSuffix[0], p)
	} else {
		ctx.Set(`pagination`, p)
	}
	return p, err
}

func PagingWithListerCond(ctx echo.Context, m Lister, cond db.Compound, varSuffix ...string) (*pagination.Pagination, error) {
	page, size, totalRows, p := PagingWithPagination(ctx)
	cnt, err := m.List(nil, nil, page, size, cond)
	if totalRows <= 0 {
		totalRows = int(cnt())
		p.SetRows(totalRows)
	}
	if len(varSuffix) > 0 {
		ctx.Set(`pagination`+varSuffix[0], p)
	} else {
		ctx.Set(`pagination`, p)
	}
	return p, err
}

func PagingWithSelectList(ctx echo.Context, param *factory.Param, varSuffix ...string) (*pagination.Pagination, error) {
	page, size, totalRows, p := PagingWithPagination(ctx)
	cnt, err := param.SetPage(page).SetSize(size).SelectList()
	if totalRows <= 0 {
		totalRows = int(cnt())
		p.SetRows(totalRows)
	}
	if len(varSuffix) > 0 {
		ctx.Set(`pagination`+varSuffix[0], p)
	} else {
		ctx.Set(`pagination`, p)
	}
	return p, err
}
