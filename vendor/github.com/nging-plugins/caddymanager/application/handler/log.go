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
	"encoding/json"
	"net/url"
	"strings"

	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/ip2region"
	"github.com/admpub/tail"
	ua "github.com/admpub/useragent"
	"github.com/nging-plugins/caddymanager/application/library/cmder"
	"github.com/nging-plugins/caddymanager/application/model"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func LogShow(ctx echo.Context) error {
	return common.LogShow(ctx, cmder.GetCaddyConfig().LogFile)
}

func VhostLog(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	if id < 1 {
		return ctx.JSON(ctx.Data().SetError(ctx.E(`id无效`)))
	}
	var err error
	m := model.NewVhost(ctx)
	err = m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.E(`不存在id为%d的网站`)
		}
		return ctx.JSON(ctx.Data().SetError(err))
	}
	var formData url.Values
	err = json.Unmarshal([]byte(m.Setting), &formData)
	if err != nil {
		return ctx.JSON(ctx.Data().SetError(err))
	}
	logFile := formData.Get(`log_file`)
	return common.LogShow(ctx, logFile, echo.H{`title`: m.Name})
}

func ParseTailLine(line *tail.Line) (interface{}, error) {
	logM := model.NewAccessLog(nil)
	err := logM.Parse(line.Text)
	res := logM.ToLite()
	realIP := logM.RemoteAddr
	if len(logM.XForwardFor) > 0 {
		realIP = strings.TrimSpace(strings.SplitN(logM.XForwardFor, ",", 2)[0])
	} else if len(logM.XRealIp) > 0 {
		realIP = logM.XRealIp
	}
	if ipInfo, _err := ip2region.IPInfo(realIP); _err == nil {
		res.Region = ipInfo.Country + " - " + ipInfo.Region + " - " + ipInfo.Province + " - " + ipInfo.City + " " + ipInfo.ISP + " (" + realIP + ")"
	} else {
		res.Region = realIP
	}
	infoUA := ua.Parse(logM.UserAgent)
	res.OS = infoUA.OS
	if len(infoUA.OSVersion) > 0 {
		res.OS += ` (` + infoUA.OSVersion + `)`
	}
	res.Brower = logM.BrowerName
	if len(res.Brower) == 0 {
		res.Brower = infoUA.Name
	}
	if len(infoUA.Version) > 0 {
		res.Brower += ` (` + infoUA.Version + `)`
	}
	res.Type = logM.BrowerType
	return res, err
}
