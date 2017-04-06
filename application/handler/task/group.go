/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package task

import (
	"errors"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func Group(ctx echo.Context) error {
	m := model.NewTaskGroup(ctx)
	page, size := handler.Paging(ctx)
	cnt, err := m.List(nil, nil, page, size)
	ret := handler.Err(ctx, err)
	ctx.SetFunc(`totalRows`, cnt)
	ctx.Set(`listData`, m.Objects())
	return ctx.Render(`task/group`, ret)
}

func GroupAdd(ctx echo.Context) error {
	var err error
	if ctx.IsPost() {
		m := model.NewTaskGroup(ctx)
		err = ctx.MustBind(m.TaskGroup)
		if err == nil {
			if len(m.Name) == 0 {
				err = errors.New(ctx.T(`分组名称不能为空`))
			} else {
				_, err = m.Add()
			}
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(`/task/group`)
		}
	}
	ctx.Set(`activeURL`, `/task/group`)
	return ctx.Render(`task/group_edit`, handler.Err(ctx, err))
}

func GroupEdit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewTaskGroup(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(`/task/group`)
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.TaskGroup)
		if err == nil {
			m.Id = id
			if len(m.Name) == 0 {
				err = errors.New(ctx.T(`分组名称不能为空`))
			} else {
				err = m.Edit(nil, `id`, id)
			}
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`修改成功`))
			return ctx.Redirect(`/task/group`)
		}
	}
	echo.StructToForm(ctx, m.TaskGroup, ``, nil)
	ctx.Set(`activeURL`, `/task/group`)
	return ctx.Render(`task/group_edit`, handler.Err(ctx, err))
}

func GroupDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewTaskGroup(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(`/task/group`)
}
