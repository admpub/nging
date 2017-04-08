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
	"fmt"
	"time"

	"strings"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/cron"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func init() {
	handler.RegisterToGroup(`/task`, func(g *echo.Group) {
		g.Route(`GET,POST`, `/index`, Index)
		g.Route(`GET,POST`, `/add`, Add)
		g.Route(`GET,POST`, `/edit`, Edit)
		g.Route(`GET,POST`, `/delete`, Delete)
		g.Route(`GET,POST`, `/start`, Start)
		g.Route(`GET,POST`, `/pause`, Pause)
		g.Route(`GET,POST`, `/run`, Run)
		g.Route(`GET,POST`, `/exit`, Exit)
		g.Route(`GET,POST`, `/start_history`, StartHistory)
		g.Route(`GET,POST`, `/group`, Group)
		g.Route(`GET,POST`, `/group_add`, GroupAdd)
		g.Route(`GET,POST`, `/group_edit`, GroupEdit)
		g.Route(`GET,POST`, `/group_delete`, GroupDelete)
		g.Route(`GET,POST`, `/log`, Log)
		g.Route(`GET,POST`, `/log_view/:id`, LogView)
		g.Route(`GET,POST`, `/log_delete`, LogDelete)
	})
}

type extra struct {
	NextTime time.Time
	Running  bool
}

func Index(ctx echo.Context) error {
	m := model.NewTask(ctx)
	page, size, totalRows, p := handler.PagingWithPagination(ctx)
	cnt, err := m.List(nil, nil, page, size)
	ret := handler.Err(ctx, err)
	if totalRows <= 0 {
		totalRows = int(cnt())
		p.SetRows(totalRows)
	}
	ctx.Set(`pagination`, p)
	tasks := m.Objects()
	gIds := []uint{}
	tg := make([]*model.TaskAndGroup, len(tasks))
	extraList := make([]*extra, len(tasks))
	for k, u := range tasks {
		tg[k] = &model.TaskAndGroup{
			Task: u,
		}
		ex := &extra{}
		entry := cron.GetEntryById(u.Id)
		if entry != nil {
			ex.NextTime = entry.Next
			ex.Running = true
		} else {
			ex.NextTime = time.Time{}
		}
		extraList[k] = ex

		if u.GroupId < 1 {
			continue
		}
		has := false
		for _, gid := range gIds {
			if gid == u.GroupId {
				has = true
				break
			}
		}
		if !has {
			gIds = append(gIds, u.GroupId)
		}
	}

	mg := model.NewTaskGroup(ctx)
	var groupList []*dbschema.TaskGroup
	_, err = mg.ListByOffset(&groupList, nil, 0, -1, db.Cond{`id IN`: gIds})
	if err != nil {
		if ret == nil {
			ret = err
		}
	} else {
		for k, v := range tg {
			for _, g := range groupList {
				if g.Id == v.GroupId {
					tg[k].Group = g
					break
				}
			}
		}
	}
	ctx.Set(`listData`, tg)
	ctx.Set(`extraList`, extraList)
	ctx.Set(`cronRunning`, cron.Running())
	ctx.Set(`histroyRunning`, cron.HistoryJobsRunning())
	return ctx.Render(`task/index`, ret)
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

func checkTaskData(ctx echo.Context, m *dbschema.Task) error {
	var err error
	if len(m.Name) == 0 {
		err = errors.New(ctx.T(`任务名不能为空`))
	} else if m.EnableNotify > 0 && len(m.NotifyEmail) > 0 {
		for _, email := range strings.Split(m.NotifyEmail, "\n") {
			email = strings.TrimSpace(email)
			if !ctx.ValidateField(`notifyEmail`, email, `email`) {
				err = errors.New(ctx.T(`无效的Email地址：%s`, email))
				break
			}
		}
	} else if err = cron.Parse(m.CronSpec); err != nil {
		err = errors.New(ctx.T(`无效的Cron时间：%s`, m.CronSpec))
	}
	return err
}

func Add(ctx echo.Context) error {
	var err error
	if ctx.IsPost() {
		m := model.NewTask(ctx)
		err = ctx.MustBind(m.Task)
		if err == nil {
			m.NotifyEmail = strings.TrimSpace(m.NotifyEmail)
			m.Command = strings.TrimSpace(m.Command)
			m.CronSpec = getCronSpec(ctx)
			m.Disabled = `Y`
			err = checkTaskData(ctx, m.Task)
			if err == nil {
				_, err = m.Add()
			}
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
			return ctx.Redirect(`/task/index`)
		}
	}
	mg := model.NewTaskGroup(ctx)
	if _, e := mg.ListByOffset(nil, nil, 0, -1); e != nil {
		err = e
	}
	ctx.Set(`groupList`, mg.Objects())
	return ctx.Render(`task/edit`, handler.Err(ctx, err))
}

func Edit(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewTask(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		handler.SendFail(ctx, err.Error())
		return ctx.Redirect(`/task/index`)
	}
	if ctx.IsPost() {
		err = ctx.MustBind(m.Task, func(k string, v []string) (string, []string) {
			if strings.ToLower(k) == `disabled` {
				return ``, nil
			}
			return k, v
		})
		if err == nil {
			m.Id = id
			m.NotifyEmail = strings.TrimSpace(m.NotifyEmail)
			m.Command = strings.TrimSpace(m.Command)
			m.CronSpec = getCronSpec(ctx)
			err = checkTaskData(ctx, m.Task)
			if err == nil {
				err = m.Edit(nil, `id`, id)
			}
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`修改成功`))
			return ctx.Redirect(`/task/index`)
		}
	}
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
	mg := model.NewTaskGroup(ctx)
	if _, e := mg.ListByOffset(nil, nil, 0, -1); e != nil {
		err = e
	}
	ctx.Set(`groupList`, mg.Objects())
	echo.StructToForm(ctx, m.Task, ``, echo.LowerCaseFirstLetter)
	ctx.Set(`activeURL`, `/task/index`)
	return ctx.Render(`task/edit`, handler.Err(ctx, err))
}

func Delete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	returnTo := ctx.Query("returnTo")
	m := model.NewTask(ctx)
	err := m.Delete(nil, db.Cond{`id`: id})
	if err == nil {
		cron.RemoveJob(id)
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	if len(returnTo) == 0 {
		returnTo = `/task/index`
	}

	return ctx.Redirect(returnTo)
}

//Start 启动任务
func Start(ctx echo.Context) error {
	id := ctx.Formx("id").Uint()
	returnTo := ctx.Query("returnTo")
	m := model.NewTask(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		return err
	}

	job, err := cron.NewJobFromTask(m.Task)
	if err != nil {
		return err
	}

	if cron.AddJob(m.CronSpec, job) {
		m.Disabled = `N`
		m.Edit(nil, `id`, id)
	}

	if len(returnTo) == 0 {
		returnTo = `/task/index`
	}

	return ctx.Redirect(returnTo)
}

//Pause 暂停任务
func Pause(ctx echo.Context) error {
	id := ctx.Formx("id").Uint()
	returnTo := ctx.Query("returnTo")
	m := model.NewTask(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		return err
	}

	cron.RemoveJob(id)
	m.Disabled = `Y`
	m.Edit(nil, `id`, id)

	if len(returnTo) == 0 {
		returnTo = `/task/index`
	}

	return ctx.Redirect(returnTo)
}

//Run 立即执行
func Run(ctx echo.Context) error {
	id := ctx.Formx("id").Uint()
	returnTo := ctx.Query("returnTo")
	m := model.NewTask(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		return err
	}

	job, err := cron.NewJobFromTask(m.Task)
	if err != nil {
		return err
	}

	job.Run()

	if len(returnTo) == 0 {
		returnTo = fmt.Sprintf(`/task/log_view/%d`, job.LogId())
	}

	return ctx.Redirect(returnTo)
}

//Exit 关闭所有任务
func Exit(ctx echo.Context) error {
	cron.Close()
	returnTo := ctx.Query("returnTo")
	if len(returnTo) == 0 {
		returnTo = `/task/index`
	}

	return ctx.Redirect(returnTo)
}

//StartHistory 继续历史任务
func StartHistory(ctx echo.Context) error {
	if !cron.HistoryJobsRunning() {
		err := cron.InitJobs()
		if err != nil {
			return err
		}
	}
	returnTo := ctx.Query("returnTo")
	if len(returnTo) == 0 {
		returnTo = `/task/index`
	}

	return ctx.Redirect(returnTo)
}
