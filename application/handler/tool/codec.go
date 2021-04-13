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

package tool

import (
	"net/http"
	"time"

	"github.com/admpub/nging/application/library/common"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func Base64(c echo.Context) (err error) {
	str := c.Formx(`text`).String()
	typ := c.Formx(`type`).String()
	if len(str) > 0 {
		var res string
		if typ == `encode` {
			res = com.Base64Encode(str)
		} else {
			res, err = com.Base64Decode(str)
		}
		if err != nil {
			return
		}
		c.Data().SetData(res)
	}
	c.Set(`title`, c.T(`Base64编码`))
	return c.Render(`/tool/codec`, nil)
}

func URL(c echo.Context) (err error) {
	str := c.Formx(`text`).String()
	typ := c.Formx(`type`).String()
	if len(str) > 0 {
		var res string
		if typ == `encode` {
			res = com.URLEncode(str)
		} else {
			res, err = com.URLDecode(str)
		}
		if err != nil {
			return
		}
		c.Data().SetData(res)
	}
	c.Set(`title`, c.T(`URL编码`))
	return c.Render(`/tool/codec`, nil)
}

func Timestamp(c echo.Context) (err error) {
	str := c.Formx(`text`).String()
	typ := c.Formx(`type`).String()
	layout := c.Formx(`layout`, `2006-01-02 15:04:05`).String()
	if len(str) > 0 {
		var t time.Time
		var res interface{}
		if typ == `encode` {
			t, err = time.Parse(layout, str)
			res = t.Unix()
		} else {
			t = time.Unix(com.Int64(str), 0)
			res = t.Format(layout)
		}
		if err != nil {
			return
		}
		c.Data().SetData(res)
	}
	c.Set(`title`, c.T(`时间戳`))
	c.Set(`timeUTC`, time.Now().UTC())
	c.Set(`timeCST`, time.Now().In(time.FixedZone(`CST`, 8*3600)))
	c.Set(`timeLocal`, time.Now().Local())
	return c.Render(`/tool/timestamp`, nil)
}

func GenPassword(ctx echo.Context) error {
	pwd, err := common.GenPassword()
	if err != nil {
		return ctx.String(err.Error(), http.StatusInternalServerError)
	}
	switch ctx.Format() {
	case echo.ContentTypeHTML, echo.ContentTypeXML:
		pwd = com.HTMLEncode(pwd)
	}
	return ctx.String(pwd)
}
