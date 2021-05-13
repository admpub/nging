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

package system

import (
	"strings"

	"github.com/webx-top/echo"
)

func NewSettings() *Settings {
	return &Settings{}
}

type AlarmThreshold struct {
	Memory float64
	CPU    float64
	Temp   float64
}

type Settings struct {
	MonitorOn      bool           // 是否开启监控
	AlarmOn        bool           // 是否开启告警
	AlarmThreshold AlarmThreshold // 告警阀值
	ReportEmail    string         // 如有多个邮箱，则一行一个
}

func (s *Settings) FromStore(h echo.H) *Settings {
	s.MonitorOn = h.Bool(`MonitorOn`)
	s.AlarmOn = h.Bool(`AlarmOn`)
	v := h.GetStore(`AlarmThreshold`)
	s.AlarmThreshold.CPU = v.Float64(`CPU`)
	s.AlarmThreshold.Memory = v.Float64(`Memory`)
	s.AlarmThreshold.Temp = v.Float64(`Temp`)
	s.ReportEmail = strings.TrimSpace(h.String(`ReportEmail`))
	return s
}
