/*
Nging is a toolbox for webmasters
Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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
package middleware

import (
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/license"
	"github.com/webx-top/echo"
)

func LicenseCheck(h echo.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		//验证授权文件
		if license.Ok(c) {
			return h.Handle(c)
		}

		//需要重新获取授权文件
		return c.Redirect(handler.URLFor(`/license`))
	}
}
