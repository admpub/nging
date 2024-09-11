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
package task

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/model"
)

func Group(ctx echo.Context) error {
	m := model.NewTaskGroup(ctx)
	_, err := common.PagingWithLister(ctx, m)
	ret := common.Err(ctx, err)
	ctx.Set(`listData`, m.Objects())
	return ctx.Render(`task/group`, ret)
}

func GroupAdd(ctx echo.Context) error {
	var err error
	if ctx.IsPost() {
		m := model.NewTaskGroup(ctx)
		err = ctx.MustBind(m.NgingTaskGroup)
		if err == nil {
			if len(m.Name) == 0 {
				err = ctx.E(`分组名称不能为空`)
			} else {
				_, err = m.Insert()
			}
		}
		if err == nil {
			common.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(backend.URLFor(`/task/group`))
		}
	}
	ctx.Set(`activeURL`, `/task/group`)
	return ctx.Render(`task/group_edit`, common.Err(ctx, err))
}

func GroupEdit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewTaskGroup(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		common.SendFail(ctx, err.Error())
		return ctx.Redirect(backend.URLFor(`/task/group`))
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingTaskGroup)
		if err == nil {
			m.Id = id
			if len(m.Name) == 0 {
				err = ctx.E(`分组名称不能为空`)
			} else {
				err = m.Update(nil, `id`, id)
			}
		}
		if err == nil {
			common.SendOk(ctx, ctx.T(`修改成功`))
			return ctx.Redirect(backend.URLFor(`/task/group`))
		}
	}
	echo.StructToForm(ctx, m.NgingTaskGroup, ``, echo.LowerCaseFirstLetter)
	ctx.Set(`activeURL`, `/task/group`)
	return ctx.Render(`task/group_edit`, common.Err(ctx, err))
}

func GroupDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewTaskGroup(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		common.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		common.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(backend.URLFor(`/task/group`))
}
