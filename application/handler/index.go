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
package handler

import (
	"github.com/admpub/caddyui/application/middleware"
	"github.com/webx-top/echo"
)

func Index(ctx echo.Context) error {
	return ctx.Redirect(`/manage`)
}

func Login(ctx echo.Context) error {
	var err error
	if ctx.Request().Method() == echo.POST {
		err = middleware.Auth(ctx, true)
		if err == nil {
			returnTo := ctx.Form(`return_to`)
			if len(returnTo) == 0 {
				returnTo = `/manage`
			}
			return ctx.Redirect(returnTo)
		}
	}
	return ctx.Render(`login`, err)
}
