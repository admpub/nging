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

package swarm

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"

	"github.com/nging-plugins/dockermanager/application/dbschema"
	"github.com/nging-plugins/dockermanager/application/library/stack"
	"github.com/nging-plugins/dockermanager/application/model"
	"github.com/nging-plugins/dockermanager/application/request"
)

func StackIndex(ctx echo.Context) error {
	name := ctx.Form(`name`, ctx.Form(`q`))
	filters := map[string]string{}
	if len(name) > 0 {
		filters[`name`] = name
	}
	list, err := stack.List(ctx, filters)
	if err != nil {
		return err
	}
	mapR := map[string]int{}
	for index, row := range list {
		mapR[row.Name] = index
	}
	m := model.NewStack(ctx)
	cond := db.NewCompounds()
	err = m.ListPage(cond, `-id`)
	for _, row := range m.Objects() {
		_, ok := mapR[row.Name]
		if !ok {
			list = append(list, stack.Item{
				Name: row.Name,
			})
		}
	}
	ctx.Set(`listData`, list)
	return ctx.Render(`docker/swarm/stack/index`, handler.Err(ctx, err))
}

func setStackByModel(c *stack.Stack, m *dbschema.NgingDockerCompose) {
	c.SetName(m.Name).SetConfigFile(m.File).SetConfigContent(m.Content).SetWorkDir(m.WorkDir)
}

func StackAdd(ctx echo.Context) error {
	user := handler.User(ctx)
	var err error
	if ctx.IsPost() {
		req := echo.GetValidated(ctx).(*request.StackAdd)

		m := model.NewStack(ctx)
		m.Content = req.Content
		m.File = req.File
		m.Name = req.Name
		m.WorkDir = req.WorkDir
		_, err = m.Add()
		if err != nil {
			goto END
		}

		c := stack.New(user.Username)
		setStackByModel(c, m.NgingDockerCompose)
		err = c.Up(ctx)
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
		} else {
			handler.SendFail(ctx, err.Error())
		}
		return ctx.Redirect(handler.URLFor(`/docker/swarm/stack/index`))
	}

END:
	ctx.Set(`activeURL`, `/docker/swarm/stack/index`)
	ctx.Set(`title`, ctx.T(`新建项目`))
	ctx.Set(`isEdit`, false)
	return ctx.Render(`docker/swarm/stack/edit`, handler.Err(ctx, err))
}

func StackEdit(ctx echo.Context) error {
	user := handler.User(ctx)
	name := ctx.Param(`id`)
	m := model.NewStack(ctx)
	err := m.Get(nil, `name`, name)
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		req := echo.GetValidated(ctx).(*request.StackEdit)

		m.Content = req.Content
		m.WorkDir = req.WorkDir
		_, err = m.Add()
		if err != nil {
			goto END
		}

		c := stack.New(user.Username)
		setStackByModel(c, m.NgingDockerCompose)
		err = c.Reload(ctx)
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
		} else {
			handler.SendFail(ctx, err.Error())
		}
		return ctx.Redirect(handler.URLFor(`/docker/swarm/stack/index`))
	}
	echo.StructToForm(ctx, m.NgingDockerCompose, ``, echo.LowerCaseFirstLetter)

END:
	ctx.Set(`name`, name)
	ctx.Set(`activeURL`, `/docker/swarm/stack/index`)
	ctx.Set(`title`, ctx.T(`修改项目`))
	ctx.Set(`isEdit`, true)
	return ctx.Render(`docker/swarm/stack/edit`, handler.Err(ctx, err))
}

func StackDetail(ctx echo.Context) error {
	name := ctx.Param(`id`)
	m := model.NewStack(ctx)
	err := m.Get(nil, `name`, name)
	if err != nil {
		return err
	}
	detail := m.Content
	ctx.Set(`name`, name)
	ctx.Set(`activeURL`, `/docker/swarm/stack/index`)
	ctx.Set(`title`, ctx.T(`配置详情`))
	ctx.Set(`detail`, detail)
	ctx.Set(`readOnly`, true)
	return ctx.Render(`docker/swarm/stack/detail`, handler.Err(ctx, err))
}

func StackListTasks(ctx echo.Context) error {
	name := ctx.Param(`id`)
	c := stack.New()
	c.SetName(name)
	var tasks []stack.TaskItem
	var err error
	tasks, err = c.ListTasks(ctx)
	ctx.Set(`name`, name)
	ctx.Set(`listData`, tasks)
	ctx.Set(`activeURL`, `/docker/swarm/stack/index`)
	return ctx.Render(`docker/swarm/stack/list_tasks`, handler.Err(ctx, err))
}

func StackListServices(ctx echo.Context) error {
	name := ctx.Param(`id`)
	c := stack.New()
	c.SetName(name)
	var tasks []stack.ServiceItem
	var err error
	tasks, err = c.ListServices(ctx)
	ctx.Set(`name`, name)
	ctx.Set(`listData`, tasks)
	ctx.Set(`activeURL`, `/docker/swarm/stack/index`)
	return ctx.Render(`docker/swarm/stack/list_services`, handler.Err(ctx, err))
}

func stackStop(ctx echo.Context, names ...string) error {
	user := handler.User(ctx)
	c := stack.New(user.Username)
	errs := common.NewErrors()
	for _, name := range names {
		err := c.SetName(name).Down(ctx)
		if err != nil {
			errs.Add(err)
			continue
		}
	}
	return errs.ToError()
}

func stackStart(ctx echo.Context, names ...string) error {
	user := handler.User(ctx)
	c := stack.New(user.Username)
	errs := common.NewErrors()
	m := model.NewStack(ctx)
	cond := db.NewCompounds()
	cond.AddKV(`type`, model.TypeStack)
	cond.AddKV(`name`, db.In(names))
	_, err := m.ListByOffset(nil, nil, 0, -1, cond.And())
	if err != nil {
		return err
	}
	for _, row := range m.Objects() {
		setStackByModel(c, row)
		err := c.Up(ctx)
		if err != nil {
			errs.Add(err)
			continue
		}
	}
	return errs.ToError()
}

func StackStart(ctx echo.Context) error {
	var err error
	name := ctx.Param(`id`)
	if name == `0` {
		err = stackStart(ctx, ctx.FormValues(`id[]`)...)
	} else {
		err = stackStart(ctx, name)
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/swarm/stack/index`))
}

func StackStop(ctx echo.Context) error {
	var err error
	name := ctx.Param(`id`)
	if name == `0` {
		err = stackStop(ctx, ctx.FormValues(`id[]`)...)
	} else {
		err = stackStop(ctx, name)
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/swarm/stack/index`))
}

func StackDelete(ctx echo.Context) error {
	var err error
	name := ctx.Param(`id`)
	if name == `0` {
		err = stackStop(ctx, ctx.FormValues(`id[]`)...)
	} else {
		err = stackStop(ctx, name)
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}
	return ctx.Redirect(handler.URLFor(`/docker/swarm/stack/index`))
}
