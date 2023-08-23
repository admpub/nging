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

package pagination

import (
	"reflect"
	"strconv"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	clientPagination "github.com/webx-top/db/lib/factory/pagination/client"
	"github.com/webx-top/echo"
	"github.com/webx-top/pagination"
)

var (
	// PageMaxSize 每页最大数据量
	PageMaxSize = 1000
	// PageDefaultSize 默认每页数据量
	PageDefaultSize = 50

	// Sorts 获取数据查询时的排序方式
	Sorts = clientPagination.Sorts
)

func SetPageDefaultSize(ctx echo.Context, pageSize int) {
	ctx.Internal().Set(`paging.pageDefaultSize`, pageSize)
}

// Paging 获取当前页码和每页数据量
func Paging(ctx echo.Context) (page int, size int) {
	page = ctx.Formx(`page`, ctx.Form(`pageNumber`)).Int()
	size = ctx.Formx(`size`, ctx.Form(`pageSize`)).Int()
	if page < 1 {
		page = 1
	}
	if size < 1 || size > PageMaxSize {
		size = ctx.Internal().Int(`paging.pageDefaultSize`)
		if size < 1 {
			size = PageDefaultSize
		}
	}
	return
}

// PagingPosition 获取偏移值
func PagingPosition(ctx echo.Context) (offset int, size int) {
	offset = ctx.Formx(`offset`).Int()
	if offset < 0 {
		offset = 0
	}
	size = ctx.Formx(`size`, ctx.Form(`pageSize`)).Int()
	if size < 1 || size > PageMaxSize {
		size = ctx.Internal().Int(`paging.pageDefaultSize`)
		if size < 1 {
			size = PageDefaultSize
		}
	}
	return
}

// PagingWithPagination 获取分页信息
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

// PagingWithPosition 获取分页信息
func PagingWithPosition(ctx echo.Context, delKeys ...string) (offset int, size int, p *pagination.Pagination) {
	offset, size = PagingPosition(ctx)
	pjax := ctx.PjaxContainer()
	if len(pjax) > 0 {
		delKeys = append(delKeys, `_pjax`)
	}
	next := strconv.FormatInt(int64(offset+size), 10)
	curr := strconv.FormatInt(int64(offset), 10)
	var prev string
	prevN := offset - size
	if prevN > 0 {
		prev = strconv.FormatInt(int64(prevN), 10)
	}
	p = pagination.New(ctx).SetPosition(prev, next, curr).SetURL(map[string]string{
		//`prev`: `prev`,
		//`curr`: `curr`,
		`next`: `offset`,
		`size`: `size`,
	}, delKeys...)
	return
}

// PagingWithLister 通过分页查询接口获取分页信息
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

type ListSizer interface {
	ListSize() int
}

// PagingWithOffsetLister 通过分页查询接口获取分页信息
func PagingWithOffsetLister(ctx echo.Context, m OffsetLister, varSuffix ...string) (*pagination.Pagination, error) {
	offset, size, p := PagingWithPosition(ctx)
	_, err := m.ListByOffset(nil, nil, offset, size)
	if len(varSuffix) > 0 {
		ctx.Set(`pagination`+varSuffix[0], p)
	} else {
		ctx.Set(`pagination`, p)
	}
	if sz, ok := m.(ListSizer); ok {
		if sz.ListSize() < size {
			p.SetPosition(p.PrevPosition(), ``, p.Position())
		}
	} else {
		if ObjectsSize(m) < size {
			p.SetPosition(p.PrevPosition(), ``, p.Position())
		}
	}
	return p, err
}

func ObjectsSize(m interface{}) int {
	rv := reflect.ValueOf(m)
	if !rv.IsValid() {
		return 0
	}
	rve := reflect.Indirect(rv)
	if rve.Kind() == reflect.Slice {
		return rve.Len()
	}
	rv = rv.MethodByName("Objects")
	if rv.IsValid() {
		rv = rv.Call(nil)[0]
		rv = reflect.Indirect(rv)
		if rv.Kind() == reflect.Slice {
			return rv.Len()
		}
	}
	return 0
}

// PagingWithListerCond 通过分页查询接口和附加条件获取分页信息
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

// PagingWithSelectList 通过Select查询参数获取分页信息
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
