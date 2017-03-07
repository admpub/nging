/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package caddy

import (
	"net/http"
	"regexp"

	"github.com/webx-top/echo"
)

var validAddonName = regexp.MustCompile(`^[a-z0-9]+$`)

func ValidAddonName(addon string) bool {
	return validAddonName.MatchString(addon)
}

func AddonIndex(ctx echo.Context) error {
	return ctx.Render(`addon/index`, nil)
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
	return ctx.Render(`addon/form/`+addon, nil)
}
