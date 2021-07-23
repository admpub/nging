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

package model

import (
	"strconv"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v3/application/dbschema"
	"github.com/admpub/nging/v3/application/library/common"
	"github.com/admpub/nging/v3/application/library/ip2region"
	"github.com/admpub/nging/v3/application/model/base"
)

func NewLoginLog(ctx echo.Context) *LoginLog {
	m := &LoginLog{
		NgingLoginLog: &dbschema.NgingLoginLog{},
		base:          base.New(ctx),
	}
	m.SetContext(ctx)
	return m
}

type LoginLog struct {
	*dbschema.NgingLoginLog
	base *base.Base
}

func (s *LoginLog) check() error {
	s.IpAddress = s.base.RealIP()
	ipInfo, err := ip2region.IPInfo(s.IpAddress)
	if err != nil {
		return err
	}
	s.IpLocation = ip2region.Stringify(ipInfo)
	s.Success = common.GetBoolFlag(s.Success)
	s.UserAgent = s.base.Request().UserAgent()
	s.Errpwd = com.MaskString(s.Errpwd)
	day, _ := strconv.Atoi(time.Now().Local().Format(`20060102`))
	s.Day = uint(day)
	return nil
}

func (s *LoginLog) Add() (pk interface{}, err error) {
	if err = s.check(); err != nil {
		return nil, err
	}
	return s.NgingLoginLog.Add()
}

func (s *LoginLog) Edit(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	if err = s.check(); err != nil {
		return err
	}
	return s.NgingLoginLog.Edit(mw, args...)
}

func (s *LoginLog) ListPage(cond *db.Compounds, sorts ...interface{}) ([]*dbschema.NgingLoginLog, error) {
	_, err := common.NewLister(s, nil, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(s.base.Context)
	return s.Objects(), err
}
