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

package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/admpub/goforever"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/config"

	conf "github.com/nging-plugins/servermanager/application/library/config"
	"github.com/nging-plugins/servermanager/application/model"
)

func DaemonIndex(ctx echo.Context) error {
	m := model.NewForeverProcess(ctx)
	cond := db.Cond{}
	_, err := handler.PagingWithListerCond(ctx, m, cond)
	ret := handler.Err(ctx, err)
	configs := m.Objects()
	for _, c := range configs {
		if c.Disabled == `N` {
			p := conf.Daemon.Child(fmt.Sprint(c.Id))
			if p != nil {
				c.Status = p.Status()
				if len(c.Error) == 0 && p.Error() != nil {
					c.Error = p.Error().Error()
				}
			}
		}
	}
	ctx.Set(`listData`, configs)
	return ctx.Render(`server/daemon_index`, ret)
}

func ignoreOptionsPrefix(k string, v []string) (string, []string) {
	if strings.HasPrefix(k, `options.`) {
		return ``, v
	}
	return k, v
}

func DaemonAdd(ctx echo.Context) error {
	user := handler.User(ctx)
	var err error
	m := model.NewForeverProcess(ctx)
	if ctx.IsPost() {
		options := echo.H{}
		err = echo.NamedStructMap(ctx.Echo(), &options, ctx.Forms(), `options`)
		if err != nil {
			goto END
		}
		err = ctx.MustBind(m.NgingForeverProcess, ignoreOptionsPrefix)
		if err != nil {
			goto END
		}
		if options.Has(`password`) && len(options.String(`password`)) > 0 {
			options.Set(`password`, config.FromFile().Encode(options.String(`password`)))
		}
		b, _ := com.JSONEncode(options)
		m.Options = com.Bytes2str(b)
		m.Uid = user.Id
		_, err = m.Add()
		if err != nil {
			goto END
		}
		if m.Disabled == `N` {
			conf.AddDaemon(m.NgingForeverProcess, true)
		}
		handler.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(handler.URLFor(`/server/daemon_index`))
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, `id`, id)
			if err == nil {
				echo.StructToForm(ctx, m.NgingForeverProcess, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
				if len(m.Options) > 0 {
					options := echo.H{}
					com.JSONDecode(com.Str2bytes(m.Options), &options)
					echo.StructToForm(ctx, options, `options`, echo.LowerCaseFirstLetter)
				}
			}
		}
		logRandName := time.Now().Format(`20060102`) + `-` + com.RandomAlphanumeric(8)
		if len(ctx.Form(`logfile`)) == 0 {
			ctx.Request().Form().Set(`logfile`, `./data/logs/forever/`+logRandName+`.info.log`)
		}
		if len(ctx.Form(`errfile`)) == 0 {
			ctx.Request().Form().Set(`errfile`, `./data/logs/forever/`+logRandName+`.err.log`)
		}
	}

END:
	ctx.Set(`activeURL`, `/server/daemon_index`)
	ctx.Set(`isWindows`, com.IsWindows)
	//ctx.Set(`isWindows`, true)
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
		oldOptions := echo.H{}
		if len(m.Options) > 0 {
			com.JSONDecode(com.Str2bytes(m.Options), &oldOptions)
		}
		options := echo.H{}
		err = echo.NamedStructMap(ctx.Echo(), &options, ctx.Forms(), `options`)
		if err != nil {
			goto END
		}
		err = ctx.MustBind(m.NgingForeverProcess, ignoreOptionsPrefix, echo.ExcludeFieldName(`created`, `uid`, `lastrun`))
		if err != nil {
			goto END
		}
		if options.Has(`password`) && len(options.String(`password`)) > 0 {
			options.Set(`password`, config.FromFile().Encode(options.String(`password`)))
		} else {
			options.Set(`password`, oldOptions.String(`password`))
		}
		b, _ := com.JSONEncode(options)
		m.Options = com.Bytes2str(b)
		m.Id = id
		err = m.Edit(nil, db.Cond{`id`: id})
		if err != nil {
			goto END
		}
		if oldName != m.Name {
			conf.Daemon.StopChild(m.Name)
			conf.AddDaemon(m.NgingForeverProcess, true)
		} else if disabled != m.Disabled {
			if m.Disabled == `N` {
				conf.AddDaemon(m.NgingForeverProcess, true)
			} else {
				conf.Daemon.StopChild(fmt.Sprint(m.Id))
			}
		}
		handler.SendOk(ctx, ctx.T(`操作成功`))
		return ctx.Redirect(handler.URLFor(`/server/daemon_index`))
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
					procs := conf.AddDaemon(m.NgingForeverProcess)
					<-goforever.RunProcess(procsName, procs)
					err = procs.Error()
					if err != nil {
						return ctx.JSON(data.SetError(err))
					}
				} else {
					conf.Daemon.StopChild(procsName)
				}
			}
			data.SetInfo(ctx.T(`操作成功`))
			procs := conf.Daemon.Child(procsName)
			if procs != nil {
				data.SetData(procs.Status())
			} else {
				data.SetData(`idle`)
			}
			return ctx.JSON(data)
		}
	}
	if err == nil {
		echo.StructToForm(ctx, m.NgingForeverProcess, ``, echo.LowerCaseFirstLetter)
		if len(m.Options) > 0 {
			options := echo.H{}
			com.JSONDecode(com.Str2bytes(m.Options), &options)
			password := options.String(`password`)
			if len(password) > 0 {
				options.Set(`password`, config.FromFile().Decode(password))
			}
			echo.StructToForm(ctx, options, `options`, echo.LowerCaseFirstLetter)
		}
	}

END:
	ctx.Set(`activeURL`, `/server/daemon_index`)
	ctx.Set(`isWindows`, com.IsWindows)
	//ctx.Set(`isWindows`, true)
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
		conf.Daemon.StopChild(fmt.Sprint(m.Id))
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/server/daemon_index`))
}

func DaemonRestart(ctx echo.Context) error {
	conf.RestartDaemon()
	data := ctx.Data()
	data.SetInfo(ctx.T(`操作成功`))
	return ctx.JSON(data)
}

func DaemonLog(ctx echo.Context) error {
	typ := ctx.Form(`type`)
	id := ctx.Formx(`id`).Uint()
	if id < 1 {
		return ctx.JSON(ctx.Data().SetError(ctx.E(`id无效`)))
	}
	var err error
	m := model.NewForeverProcess(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.E(`不存在id为%d的任务`)
		}
		return ctx.JSON(ctx.Data().SetError(err))
	}
	var logFile string
	switch typ {
	case `error`:
		logFile = m.Errfile
	default:
		logFile = m.Logfile
	}
	return common.LogShow(ctx, logFile, echo.H{`title`: m.Name, `charset`: m.LogCharset})
}
