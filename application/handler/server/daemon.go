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

package server

import (
	"fmt"
	"strings"

	"github.com/admpub/goforever"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func DaemonIndex(ctx echo.Context) error {
	m := model.NewForeverProcess(ctx)
	cond := db.Cond{}
	_, err := handler.PagingWithListerCond(ctx, m, cond)
	ret := handler.Err(ctx, err)
	configs := m.Objects()
	for _, c := range configs {
		if c.Disabled == `N` {
			p := config.Daemon.Child(fmt.Sprint(c.Id))
			if p != nil {
				c.Status = p.Status
				if len(c.Error) == 0 && p.Error() != nil {
					c.Error = p.Error().Error()
				}
			}
		}
	}
	ctx.Set(`listData`, configs)
	return ctx.Render(`server/daemon_index`, ret)
}

func DaemonAdd(ctx echo.Context) error {
	var err error
	m := model.NewForeverProcess(ctx)
	if ctx.IsPost() {
		name := ctx.Form(`name`)
		if len(name) == 0 {
			err = ctx.E(`名称不能为空`)
		} else if y, e := m.Exists(name); e != nil {
			err = e
		} else if y {
			err = ctx.E(`名称已经存在`)
		} else {
			err = ctx.MustBind(m.ForeverProcess)
		}
		if err == nil {
			_, err = m.Add()
			if err == nil {
				if m.Disabled == `N` {
					config.AddDaemon(m.ForeverProcess, true)
				}
			}
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/server/daemon_index`))
			}
		}
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.ForeverProcess, ``, func(topName, fieldName string) string {
					return echo.LowerCaseFirstLetter(topName, fieldName)
				})
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}
	ctx.Set(`activeURL`, `/server/daemon_index`)
	return ctx.Render(`server/daemon_edit`, err)
}

func DaemonEdit(ctx echo.Context) error {
	var err error
	id := ctx.Formx(`id`).Uint()
	m := model.NewForeverProcess(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	disabled := m.Disabled
	oldName := m.Name
	if ctx.IsPost() {
		name := ctx.Form(`name`)
		if len(name) == 0 {
			err = ctx.E(`名称不能为空`)
		} else if y, e := m.Exists2(name, id); e != nil {
			err = e
		} else if y {
			err = ctx.E(`名称已经存在`)
		} else {
			err = ctx.MustBind(m.ForeverProcess, func(k string, v []string) (string, []string) {
				switch strings.ToLower(k) {
				case `created`: //禁止修改创建时间和用户名
					return ``, v
				}
				return k, v
			})
		}
		if err == nil {
			m.Id = id
			err = m.Edit(nil, db.Cond{`id`: id})
			if err == nil {
				if oldName != m.Name {
					config.Daemon.StopChild(m.Name)
					config.AddDaemon(m.ForeverProcess, true)
				} else if disabled != m.Disabled {
					if m.Disabled == `N` {
						config.AddDaemon(m.ForeverProcess, true)
					} else {
						config.Daemon.StopChild(fmt.Sprint(m.Id))
					}
				}
			}
			if err == nil {
				handler.SendOk(ctx, ctx.T(`操作成功`))
				return ctx.Redirect(handler.URLFor(`/server/daemon_index`))
			}
		}
	} else if ctx.IsAjax() {
		setDisabled := ctx.Query(`disabled`)
		if len(setDisabled) > 0 {
			m.Disabled = setDisabled
			data := ctx.Data()
			err = m.Edit(nil, db.Cond{`id`: id})
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			procsName := fmt.Sprint(m.Id)
			if disabled != m.Disabled {
				if m.Disabled == `N` {
					procs := config.AddDaemon(m.ForeverProcess)
					<-goforever.RunProcess(procsName, procs)
					err = procs.Error()
					if err != nil {
						return ctx.JSON(data.SetError(err))
					}
				} else {
					config.Daemon.StopChild(procsName)
				}
			}
			data.SetInfo(ctx.T(`操作成功`))
			procs := config.Daemon.Child(procsName)
			if procs != nil {
				data.SetData(procs.Status)
			} else {
				data.SetData(`idle`)
			}
			return ctx.JSON(data)
		}
	}
	if err == nil {
		echo.StructToForm(ctx, m.ForeverProcess, ``, func(topName, fieldName string) string {
			return echo.LowerCaseFirstLetter(topName, fieldName)
		})
	}

	ctx.Set(`activeURL`, `/server/daemon_index`)
	return ctx.Render(`server/daemon_edit`, err)
}

func DaemonDelete(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewForeverProcess(ctx)
	err := m.Get(nil, db.Cond{`id`: id})
	if err == nil {
		err = m.Delete(nil, db.Cond{`id`: id})
	}
	if err == nil {
		config.Daemon.StopChild(fmt.Sprint(m.Id))
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/server/daemon_index`))
}
