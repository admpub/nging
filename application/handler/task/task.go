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
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/dbschema"
	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/cron"
	cronWriter "github.com/admpub/nging/v4/application/library/cron/writer"
	"github.com/admpub/nging/v4/application/model"
)

type extra struct {
	NextTime time.Time
	Running  bool
}

func buidlPattern(subPattern string, extras ...string) string {
	chars := strings.Join(extras, ``)
	return `^([*` + chars + `]|` + subPattern + `(,` + subPattern + `)*)$`
}

func Index(ctx echo.Context) error {
	groupId := ctx.Formx(`groupId`).Uint()
	m := model.NewTask(ctx)
	cond := db.Compounds{}
	if groupId > 0 {
		cond.AddKV(`group_id`, groupId)
	}
	q := ctx.Formx(`q`).String()
	if len(q) > 0 {
		cond.AddKV(`name`, db.Like(`%`+q+`%`))
	}
	var tasks []*model.TaskAndGroup
	_, err := handler.PagingWithLister(ctx, handler.NewLister(m, &tasks, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, cond.And()))
	extraList := make([]*extra, len(tasks))
	for k, u := range tasks {
		ex := &extra{}
		entry := cron.GetEntryById(u.Id)
		if entry != nil {
			ex.NextTime = entry.Next
			ex.Running = true
		} else {
			ex.NextTime = time.Time{}
		}
		extraList[k] = ex
	}

	mg := model.NewTaskGroup(ctx)
	var groupList []*dbschema.NgingTaskGroup
	mg.ListByOffset(&groupList, nil, 0, -1)
	ctx.Set(`listData`, tasks)
	ctx.Set(`extraList`, extraList)
	ctx.Set(`cronRunning`, cron.Running())
	ctx.Set(`histroyRunning`, cron.HistoryJobsRunning())
	ctx.Set(`notRecordPrefixFlag`, cronWriter.NotRecordPrefixFlag)
	ctx.Set(`groupList`, groupList)
	ctx.Set(`groupId`, groupId)
	return ctx.Render(`task/index`, handler.Err(ctx, err))
}

func getCronSpec(ctx echo.Context) string {
	seconds := ctx.Form(`seconds`)
	minutes := ctx.Form(`minutes`)
	hours := ctx.Form(`hours`)
	dayOfMonth := ctx.Form(`dayOfMonth`)
	month := ctx.Form(`month`)
	dayOfWeek := ctx.Form(`dayOfWeek`)
	return seconds + ` ` + minutes + ` ` + hours + ` ` + dayOfMonth + ` ` + month + ` ` + dayOfWeek
}

func checkTaskData(ctx echo.Context, m *dbschema.NgingTask) error {
	var err error
	if len(m.Name) == 0 {
		err = ctx.E(`任务名不能为空`)
	} else if m.EnableNotify > 0 && len(m.NotifyEmail) > 0 {
		for _, email := range strings.Split(m.NotifyEmail, "\n") {
			email = strings.TrimSpace(email)
			if !ctx.Validate(`notifyEmail`, email, `email`).Ok() {
				err = ctx.E(`无效的Email地址：%s`, email)
				break
			}
		}
	} else if err = cron.Parse(m.CronSpec); err != nil {
		err = ctx.E(`无效的Cron时间：%s`, m.CronSpec)
	}
	return err
}

func Add(ctx echo.Context) error {
	var err error
	m := model.NewTask(ctx)
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingTask)
		if err == nil {
			m.NotifyEmail = strings.TrimSpace(m.NotifyEmail)
			m.Command = strings.TrimSpace(m.Command)
			m.CronSpec = getCronSpec(ctx)
			m.Disabled = `Y`
			m.Uid = handler.User(ctx).Id
			err = checkTaskData(ctx, m.NgingTask)
			if err == nil {
				_, err = m.Insert()
			}
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(handler.URLFor(`/task/index`))
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				setFormData(ctx, m)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}
	mg := model.NewTaskGroup(ctx)
	if _, e := mg.ListByOffset(nil, nil, 0, -1); e != nil {
		err = e
	}
	ctx.Set(`groupList`, mg.Objects())
	ctx.SetFunc(`buildPattern`, buidlPattern)
	return ctx.Render(`task/edit`, handler.Err(ctx, err))
}

func setFormData(ctx echo.Context, m *model.Task) {
	specs := strings.Split(m.CronSpec, ` `)
	switch len(specs) {
	case 6:
		ctx.Request().Form().Set(`dayOfWeek`, specs[5])
		fallthrough
	case 5:
		ctx.Request().Form().Set(`month`, specs[4])
		fallthrough
	case 4:
		ctx.Request().Form().Set(`dayOfMonth`, specs[3])
		fallthrough
	case 3:
		ctx.Request().Form().Set(`hours`, specs[2])
		fallthrough
	case 2:
		ctx.Request().Form().Set(`minutes`, specs[1])
		fallthrough
	case 1:
		ctx.Request().Form().Set(`seconds`, specs[0])
	}
	echo.StructToForm(ctx, m.NgingTask, ``, echo.LowerCaseFirstLetter)
}

func Edit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewTask(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(handler.URLFor(`/task/index`))
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingTask, echo.ExcludeFieldName(`disabled`))
		if err == nil {
			m.Id = id
			m.NotifyEmail = strings.TrimSpace(m.NotifyEmail)
			m.Command = strings.TrimSpace(m.Command)
			m.CronSpec = getCronSpec(ctx)
			err = checkTaskData(ctx, m.NgingTask)
			if err == nil {
				err = m.Update(nil, `id`, id)
			}
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`修改成功`))
			return ctx.Redirect(handler.URLFor(`/task/index`))
		}
	}
	setFormData(ctx, m)
	mg := model.NewTaskGroup(ctx)
	if _, e := mg.ListByOffset(nil, nil, 0, -1); e != nil {
		err = e
	}
	ctx.Set(`groupList`, mg.Objects())
	ctx.Set(`activeURL`, `/task/index`)
	ctx.Set(`notRecordPrefixFlag`, cronWriter.NotRecordPrefixFlag)
	ctx.SetFunc(`buildPattern`, buidlPattern)
	return ctx.Render(`task/edit`, handler.Err(ctx, err))
}

func Delete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	next := ctx.Query("next")
	m := model.NewTask(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		cron.RemoveJob(id)
		logM := model.NewTaskLog(ctx)
		err = logM.Delete(nil, db.Cond{`task_id`: id})
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
		} else {
			handler.SendFail(ctx, err.Error())
		}
	}

	if len(next) == 0 {
		next = handler.URLFor(`/task/index`)
	}

	return ctx.Redirect(next)
}

//Start 启动任务
func Start(ctx echo.Context) error {
	id := ctx.Formx("id").Uint()
	next := ctx.Query("next")
	m := model.NewTask(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		return err
	}

	job, err := cron.NewJobFromTask(context.Background(), m.NgingTask)
	if err != nil {
		return err
	}

	if cron.AddJob(m.CronSpec, job) {
		m.Disabled = `N`
		err = m.Update(nil, `id`, id)
		if err != nil {
			return err
		}
	}
	if ctx.Format() == `json` {
		ex := echo.Store{`Running`: false, `Disabled`: m.Disabled}
		entry := cron.GetEntryById(id)
		if entry != nil {
			ex[`NextTime`] = entry.Next.Format(`2006-01-02 15:04:05`)
			ex[`Running`] = true
		} else {
			ex[`NextTime`] = ``
		}
		data := ctx.Data()
		data.SetInfo(ctx.T(`启动成功`)).SetData(ex)
		return ctx.JSON(data)
	}
	if len(next) == 0 {
		next = handler.URLFor(`/task/index`)
	}

	return ctx.Redirect(next)
}

//Pause 暂停任务
func Pause(ctx echo.Context) error {
	id := ctx.Formx("id").Uint()
	next := ctx.Query("next")
	m := model.NewTask(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		return err
	}

	cron.RemoveJob(id)
	m.Disabled = `Y`
	err = m.Update(nil, `id`, id)
	if err != nil {
		return err
	}

	if ctx.Format() == `json` {
		ex := echo.Store{`Running`: false, `Disabled`: m.Disabled}
		entry := cron.GetEntryById(id)
		if entry != nil {
			ex[`NextTime`] = entry.Next.Format(`2006-01-02 15:04:05`)
			ex[`Running`] = true
		} else {
			ex[`NextTime`] = ``
		}
		data := ctx.Data()
		data.SetInfo(ctx.T(`任务已暂停`)).SetData(ex)
		return ctx.JSON(data)
	}
	if len(next) == 0 {
		next = handler.URLFor(`/task/index`)
	}

	return ctx.Redirect(next)
}

//Run 立即执行
func Run(ctx echo.Context) error {
	id := ctx.Formx("id").Uint()
	next := ctx.Query("next")
	m := model.NewTask(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		return err
	}

	job, err := cron.NewJobFromTask(ctx.Request().StdRequest().Context(), m.NgingTask)
	if err != nil {
		return err
	}

	job.Run()

	if len(next) == 0 {
		logID := job.LogID()
		if logID <= 0 {
			taskLog := job.LogData()
			return renderLogViewData(ctx, taskLog, err)
		}
		next = fmt.Sprintf(`/task/log_view/%d`, logID)
	}

	return ctx.Redirect(next)
}

//Exit 关闭所有任务
func Exit(ctx echo.Context) error {
	cron.Close()
	next := ctx.Query("next")
	if len(next) == 0 {
		next = handler.URLFor(`/task/index`)
	}

	return ctx.Redirect(next)
}

//StartHistory 继续历史任务
func StartHistory(ctx echo.Context) error {
	if !cron.HistoryJobsRunning() {
		err := cron.InitJobs(context.Background())
		if err != nil {
			return err
		}
	}
	next := ctx.Query("next")
	if len(next) == 0 {
		next = handler.URLFor(`/task/index`)
	}

	return ctx.Redirect(next)
}
