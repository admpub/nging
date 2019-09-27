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

package manager

import (
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/model/file"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func Finder(ctx echo.Context) error {
	fileM := file.NewFile(ctx)
	cond := db.NewCompounds()
	_, err := common.NewLister(fileM.File, nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()).Paging(ctx)
	if err != nil {
		return err
	}
	ctx.Set(`listData`, fileM.Objects())
	return ctx.Render(`manager/finder`, err)
}
