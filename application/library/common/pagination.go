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

package common

import (
	dbPagination "github.com/webx-top/db/lib/factory/pagination"
	clientPagination "github.com/webx-top/db/lib/factory/pagination/client"
)

var (
	// Sorts 获取数据查询时的排序方式
	Sorts                  = clientPagination.Sorts
	Paging                 = dbPagination.Paging
	PagingWithPagination   = dbPagination.PagingWithPagination
	PagingWithPosition     = dbPagination.PagingWithPosition
	PagingWithLister       = dbPagination.PagingWithLister
	PagingWithOffsetLister = dbPagination.PagingWithOffsetLister
	PagingWithListerCond   = dbPagination.PagingWithListerCond
	PagingWithSelectList   = dbPagination.PagingWithSelectList
	NewLister              = dbPagination.NewLister
	NewListParam           = dbPagination.NewListParam
	NewOffsetLister        = dbPagination.NewOffsetLister
)

type (
	PageListFunc   = dbPagination.PageListFunc
	OffsetListFunc = dbPagination.OffsetListFunc
	OffsetLister   = dbPagination.OffsetLister
	List           = dbPagination.List
	Lister         = dbPagination.Lister
)

// FloorNumber 楼层号
func FloorNumber(page int, pageSize int, index int) int {
	if page < 1 {
		page = 1
	}
	return (page-1)*pageSize + (index + 1)
}
