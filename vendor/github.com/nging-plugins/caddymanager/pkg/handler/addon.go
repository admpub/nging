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
	"net/http"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func ValidAddonName(addon string) bool {
	return com.IsAlphaNumericUnderscore(addon)
}

func AddonIndex(ctx echo.Context) error {
	return ctx.Render(`caddy/addon/index`, nil)
}

func AddonForm(ctx echo.Context) error {
	addon := ctx.Query(`addon`)
	if len(addon) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, ctx.T("参数 addon 的值不能为空"))
	}
	if !ValidAddonName(addon) {
		return echo.NewHTTPError(http.StatusBadRequest, ctx.T("参数 addon 的值包含非法字符"))
	}
	ctx.SetFunc(`Val`, func(name, defaultValue string) string {
		return defaultValue
	})
	index := ctx.Queryx(`index`, `0`).Int()
	return ctx.Render(`caddy/addon/form/`+addon, index)
}
