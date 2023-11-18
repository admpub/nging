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

package compose

import (
	"os"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"

	"github.com/nging-plugins/dockermanager/application/dbschema"
	"github.com/nging-plugins/dockermanager/application/library/compose"
	"github.com/nging-plugins/dockermanager/application/model"
	"github.com/nging-plugins/dockermanager/application/request"
)

func Index(ctx echo.Context) error {
	user := handler.User(ctx)
	name := ctx.Form(`name`, ctx.Form(`q`))
	filters := map[string]string{}
	if len(name) > 0 {
		filters[`name`] = name
	}
	list, err := compose.List(ctx, filters)
	if err != nil {
		return err
	}
	mapR := map[string]int{}
	for index, row := range list {
		mapR[row.Name] = index
	}
	m := model.NewCompose(ctx)
	cond := db.NewCompounds()
	err = m.ListPage(cond, `-id`)
	c := compose.New(user.Username)
	rows := m.Objects()
	for _, row := range rows {
		_, ok := mapR[row.Name]
		var errR error
		if !ok {
			list = append(list, compose.ComposeItem{
				Name:        row.Name,
				Status:      `stopped`,
				ConfigFiles: row.File,
			})
		} else {
			c.SetName(row.Name).SetConfigFile(row.File).SetConfigContent(row.Content)
			errR = c.Up(ctx, true)
			delete(mapR, row.Name)
		}
		if errR != nil {
			log.Error(errR)
		}
	}
	if len(mapR) > 0 {
		for _, idx := range mapR {
			m.Name = list[idx].Name
			m.File = list[idx].ConfigFiles
			_, errR := m.Add()
			if errR != nil {
				log.Error(errR)
			}
			m.Reset()
		}
	}
	ctx.Set(`listData`, list)
	return ctx.Render(`docker/base/compose/index`, handler.Err(ctx, err))
}

func setByModel(c *compose.Compose, m *dbschema.NgingDockerCompose) {
	c.SetName(m.Name).SetConfigFile(m.File).SetConfigContent(m.Content).SetWorkDir(m.WorkDir)
}

func Add(ctx echo.Context) error {
	user := handler.User(ctx)
	var err error
	if ctx.IsPost() {
		req := echo.GetValidated(ctx).(*request.ComposeAdd)

		m := model.NewCompose(ctx)
		m.Content = req.Content
		m.File = req.File
		m.Name = req.Name
		m.WorkDir = req.WorkDir
		_, err = m.Add()
		if err != nil {
			goto END
		}
		c := compose.New(user.Username)
		setByModel(c, m.NgingDockerCompose)
		err = c.Up(ctx, true)
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
		} else {
			handler.SendFail(ctx, err.Error())
		}
		return ctx.Redirect(handler.URLFor(`/docker/base/compose/index`))
	}

END:
	ctx.Set(`activeURL`, `/docker/base/compose/index`)
	ctx.Set(`title`, ctx.T(`新建项目`))
	ctx.Set(`isEdit`, false)
	return ctx.Render(`docker/base/compose/edit`, handler.Err(ctx, err))
}

func getRow(ctx echo.Context, name string, mustRunning bool) (row *model.Compose, running bool, err error) {
	var list []compose.ComposeItem
	list, err = compose.List(ctx, map[string]string{`name`: name})
	if err != nil {
		return
	}

	m := model.NewCompose(ctx)
	dbErr := m.Get(nil, `name`, name)
	if dbErr != nil {
		if dbErr != db.ErrNoMoreRows {
			err = dbErr
			return
		}
	}
	running = len(list) > 0
	if !running {
		if !mustRunning && dbErr == nil {
			return m, running, err
		}
		err = ctx.NewError(code.DataNotFound, `不存在此项目: %s`, name).SetZone(`name`)
		return
	}
	var content string
	if b, err := os.ReadFile(m.File); err == nil {
		content = com.Bytes2str(b)
	}
	if dbErr != nil {
		m.Name = list[0].Name
		m.File = list[0].ConfigFiles
		m.Content = content
		_, err = m.Add()
	} else {
		m.Name = list[0].Name
		set := echo.H{}
		if m.File != list[0].ConfigFiles {
			set.Set(`file`, list[0].ConfigFiles)
		}
		if m.Content != content {
			set.Set(`content`, content)
		}
		if len(set) > 0 {
			m.UpdateFields(nil, set, `id`, m.Id)
		}
		m.File = list[0].ConfigFiles
		m.Content = content
	}
	return m, running, err
}

func Edit(ctx echo.Context) error {
	user := handler.User(ctx)
	var err error
	name := ctx.Param(`id`)
	var b []byte
	row, running, err := getRow(ctx, name, false)
	if err != nil {
		return err
	}
	if !running {
		row.File = compose.ConfigPath(row.Name)
	}
	if ctx.IsPost() {
		req := echo.GetValidated(ctx).(*request.ComposeEdit)

		row.Content = req.Content
		row.WorkDir = req.WorkDir
		_, err = row.Add()
		if err != nil {
			goto END
		}
		b = com.Str2bytes(req.Content)
		err = os.WriteFile(row.File, b, os.ModePerm)
		if err != nil {
			goto END
		}
		c := compose.New(user.Username)
		setByModel(c, row.NgingDockerCompose)
		err = c.Reload(ctx)
		if err == nil {
			handler.SendOk(ctx, ctx.T(`操作成功`))
		} else {
			handler.SendFail(ctx, err.Error())
		}
		return ctx.Redirect(handler.URLFor(`/docker/base/compose/index`))
	}
	echo.StructToForm(ctx, row.NgingDockerCompose, ``, echo.LowerCaseFirstLetter)

END:
	ctx.Set(`name`, name)
	ctx.Set(`configFile`, row.File)
	ctx.Set(`activeURL`, `/docker/base/compose/index`)
	ctx.Set(`title`, ctx.T(`修改项目`))
	ctx.Set(`isEdit`, true)
	return ctx.Render(`docker/base/compose/edit`, handler.Err(ctx, err))
}

func Detail(ctx echo.Context) error {
	var err error
	name := ctx.Param(`id`)
	row, _, err := getRow(ctx, name, false)
	if err != nil {
		return err
	}
	detail := row.Content
	ctx.Set(`name`, name)
	ctx.Set(`configFile`, row.File)
	ctx.Set(`activeURL`, `/docker/base/compose/index`)
	ctx.Set(`title`, ctx.T(`配置详情`))
	ctx.Set(`detail`, detail)
	ctx.Set(`readOnly`, true)
	return ctx.Render(`docker/base/compose/detail`, handler.Err(ctx, err))
}

func Scale(ctx echo.Context) error {
	user := handler.User(ctx)
	name := ctx.Param(`id`)
	service := ctx.Param(`service`)
	replicas := ctx.Formx(`replicas`).Uint()
	if replicas < 1 {
		return ctx.NewError(code.InvalidParameter, `副本数量不能小于1`).SetZone(`replicas`)
	}
	row, _, err := getRow(ctx, name, true)
	if err != nil {
		return err
	}
	c := compose.New(user.Username)
	setByModel(c, row.NgingDockerCompose)
	err = c.ScaleService(ctx, service, replicas)
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}
	return ctx.Redirect(handler.URLFor(`/docker/base/compose/index`))
}

func ListContainers(ctx echo.Context) error {
	name := ctx.Param(`id`)
	row, _, err := getRow(ctx, name, true)
	if err != nil {
		return err
	}
	c := compose.New()
	setByModel(c, row.NgingDockerCompose)
	var containers []compose.ContainerItem
	containers, err = c.ListContainers(ctx)
	ctx.Set(`name`, name)
	ctx.Set(`listData`, containers)
	ctx.Set(`activeURL`, `/docker/base/compose/index`)
	return ctx.Render(`docker/base/compose/list_containers`, handler.Err(ctx, err))
}

func start(ctx echo.Context, names ...string) error {
	user := handler.User(ctx)
	c := compose.New(user.Username)
	errs := common.NewErrors()
	m := model.NewCompose(ctx)
	cond := db.NewCompounds()
	cond.AddKV(`type`, model.TypeCompose)
	cond.AddKV(`name`, db.In(names))
	_, err := m.ListByOffset(nil, nil, 0, -1, cond.And())
	if err != nil {
		return err
	}
	for _, row := range m.Objects() {
		setByModel(c, row)
		err := c.Up(ctx, true)
		if err != nil {
			errs.Add(err)
			continue
		}
	}
	return errs.ToError()
}

func stop(ctx echo.Context, names ...string) error {
	user := handler.User(ctx)
	c := compose.New(user.Username)
	errs := common.NewErrors()
	for _, name := range names {
		if len(name) == 0 {
			continue
		}
		list, err := compose.List(ctx, map[string]string{`name`: name})
		if err != nil {
			return err
		}
		if len(list) == 0 {
			continue
		}
		err = c.SetConfigFile(list[0].ConfigFiles).Down(ctx)
		if err != nil {
			errs.Add(err)
			continue
		}
	}
	return errs.ToError()
}

func Stop(ctx echo.Context) error {
	var err error
	name := ctx.Param(`id`)
	if name == `0` {
		err = stop(ctx, ctx.FormValues(`id[]`)...)
	} else {
		err = stop(ctx, name)
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/base/compose/index`))
}

func Start(ctx echo.Context) error {
	var err error
	name := ctx.Param(`id`)
	if name == `0` {
		err = start(ctx, ctx.FormValues(`id[]`)...)
	} else {
		err = start(ctx, name)
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/base/compose/index`))
}

func Delete(ctx echo.Context) error {
	var err error
	name := ctx.Param(`id`)
	var names []string
	if name == `0` {
		names = ctx.FormValues(`id[]`)
		err = stop(ctx, names...)
	} else {
		names = []string{name}
		err = stop(ctx, name)
	}
	if err == nil {
		m := model.NewCompose(ctx)
		m.Delete(nil, `name`, db.In(names))
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/base/compose/index`))
}
