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
	"regexp"
	"strings"

	"github.com/webx-top/echo"

	"github.com/nging-plugins/caddymanager/application/dbschema"
)

func NewAccessLog(ctx echo.Context) *AccessLog {
	return &AccessLog{
		NgingAccessLog: dbschema.NewNgingAccessLog(ctx),
	}
}

type AccessLog struct {
	*dbschema.NgingAccessLog
}

// Parse 解析单行字符串到日志对象
func (l *AccessLog) Parse(line string, args ...interface{}) (err error) {
	line = strings.TrimRight(line, " \r\n")
	if len(line) == 0 {
		return nil
	}
	size := len(args)
	if size > 0 {
		var layout string
		if size > 1 {
			layout, _ = args[1].(string)
		}
		switch v := args[0].(type) {
		case *regexp.Regexp:
			err = l.parseWithPattern(line, v)
		case string:
			switch v {
			//case `nginx`:
			default:
				err = l.parseCaddy(line, layout)
			}
		}
	} else {
		err = l.parseCaddy(line, ``)
	}
	return
}

func (m *AccessLog) ToLite() *AccessLogLite {
	a := &AccessLogLite{
		Date:   m.TimeLocal,
		OS:     ``,
		Brower: m.BrowerName,
		Region: ``,
		Type:   m.BrowerType,

		Version: m.Version,
		User:    m.User,
		Method:  m.Method,
		Scheme:  m.Scheme,
		Host:    m.Host,
		URI:     m.Uri,

		Referer: m.Referer,

		BodyBytes:  m.BodyBytes,
		Elapsed:    m.Elapsed,
		StatusCode: m.StatusCode,
		UserAgent:  m.UserAgent,
	}
	return a
}

func (m *AccessLog) ToMap() echo.Store {
	data := echo.Store{}
	data.Set(`TimeLocal`, m.TimeLocal)
	data.Set(`RemoteAddr`, m.RemoteAddr)
	data.Set(`XRealIP`, m.XRealIp)
	data.Set(`XForwardFor`, m.XForwardFor)
	data.Set(`LocalAddr`, m.LocalAddr)
	data.Set(`User`, m.User)
	data.Set(`Version`, m.Version)
	data.Set(`Referer`, m.Referer)
	data.Set(`UserAgent`, m.UserAgent)
	data.Set(`Path`, m.Uri)
	data.Set(`Method`, m.Method)
	data.Set(`Scheme`, m.Scheme)
	data.Set(`Host`, m.Host)
	data.Set(`BrowerName`, m.BrowerName)
	data.Set(`BrowerType`, m.BrowerType)
	data.Set(`BytesSent`, m.BodyBytes)
	data.Set(`StatusCode`, m.StatusCode)
	data.Set(`UpstreamTime`, m.Elapsed)
	data.Set(`RequestTime`, m.Elapsed)

	return data
}
