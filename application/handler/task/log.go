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
	"time"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/cron"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func Log(ctx echo.Context) error {
	taskID := ctx.Formx(`taskId`).Uint()
	totalRows := ctx.Formx(`rows`).Int()
	m := model.NewTaskLog(ctx)
	page, size, totalRows, p := handler.PagingWithPagination(ctx)
	cond := db.Cond{}
	var task *model.Task
	var err error
	if taskID > 0 {
		task = model.NewTask(ctx)
		err = task.Get(nil, `id`, taskID)
		cond[`task_id`] = taskID
	}
	cnt, err2 := m.List(nil, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, page, size, cond)
	if err2 != nil {
		err = err2
	}
	if totalRows <= 0 {
		totalRows = int(cnt())
		p.SetRows(totalRows)
	}
	ctx.Set(`listData`, m.Objects())
	ctx.Set(`pagination`, p)
	ctx.Set(`task`, task)
	ret := handler.Err(ctx, err)
	ctx.Set(`activeURL`, `/task/index`)
	return ctx.Render(`task/log`, ret)
}

func LogView(ctx echo.Context) error {
	id := ctx.Paramx(`id`).Uint()
	m := model.NewTaskLog(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(`/task/log`)
	}
	ctx.Set(`data`, m)
	ctx.Set(`activeURL`, `/task/index`)
	var task *model.Task
	if m.TaskId > 0 {
		task = model.NewTask(ctx)
		err = task.Get(nil, `id`, m.TaskId)
	}
	ex := &extra{}
	entry := cron.GetEntryById(task.Id)
	if entry != nil {
		ex.NextTime = entry.Next
		ex.Running = true
	} else {
		ex.NextTime = time.Time{}
	}
	ctx.Set(`task`, task)
	ctx.Set(`extra`, ex)
	return ctx.Render(`task/log_view`, handler.Err(ctx, err))
}

func LogDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewTaskLog(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(`/task/log`)
}
