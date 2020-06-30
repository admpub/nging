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

package collector

import (
	"context"
	"sync"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/library/collector/exec"
	"github.com/admpub/nging/application/library/collector/export"
	"github.com/admpub/nging/application/library/cron"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
)

var GoProcs = sync.Map{}

func NewProc(r *exec.Rules) *GoProcess {
	return &GoProcess{
		Rules: r,
	}
}

type GoProcess struct {
	Rules *exec.Rules
	Done  bool
}

func (g *GoProcess) Close(k interface{}) {
	g.Done = true
	GoProcs.Delete(k)
}

func (g *GoProcess) IsExited() bool {
	return g.Done
}

func Exit(k interface{}) (bool, error) {
	proc, exists := GoProcs.Load(k)
	if exists {
		if p, ok := proc.(*GoProcess); ok {
			p.Close(k)
		}
	}
	return exists, nil
}

func Go(k interface{}, r *exec.Rules, f func(), ctx context.Context) (err error) {
	_, err = Exit(k)
	if err != nil {
		return
	}
	process := NewProc(r)
	GoProcs.Store(k, process)
	r.SetExitedFn(process.IsExited)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error(err)
				process.Close(k)
			}
		}()
		select {
		case <-ctx.Done():
			process.Close(k)
			return
		default:
			f()
			process.Close(k)
		}
	}()
	return
}

//CollectPageJob 计划任务调用方式
func CollectPageJob(id string) cron.Runner {
	return func(timeout time.Duration) (out string, runingErr string, onErr error, isTimeout bool) {
		m := model.NewCollectorPage(nil)
		err := m.Get(nil, db.Cond{`id`: id})
		if err != nil {
			onErr = cron.ErrFailure
			runingErr = `获取页面采集参数失败：` + err.Error()
			return
		}
		data, err := m.FullData()
		if err != nil {
			onErr = cron.ErrFailure
			runingErr = `获取页面采集规则失败：` + err.Error()
			return
		}
		data.SetExportFn(export.Export)
		result, err := data.Collect(false, nil, nil)
		if err != nil {
			runingErr = `采集出错：` + err.Error()
			return
		}
		b, err := com.JSONEncode(result, ` `)
		if err != nil {
			runingErr = err.Error()
			return
		}
		out = string(b)
		return
	}
}
