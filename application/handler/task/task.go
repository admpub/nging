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
package task

import (
	"github.com/admpub/nging/application/handler"
	"github.com/webx-top/echo"
)

func init() {
	handler.RegisterToGroup(`/task`, func(g *echo.Group) {
		g.Route(`GET,POST`, `/index`, Index)
		g.Route(`GET,POST`, `/add`, Add)
		g.Route(`GET,POST`, `/edit`, Edit)
		g.Route(`GET,POST`, `/delete`, Delete)
		g.Route(`GET,POST`, `/group`, Group)
		g.Route(`GET,POST`, `/group_add`, GroupAdd)
		g.Route(`GET,POST`, `/group_edit`, GroupEdit)
		g.Route(`GET,POST`, `/group_delete`, GroupDelete)
	})
}

func Index(c echo.Context) error {
	var err error
	return c.Render(`task/index`, handler.Err(c, err))
}

func Add(c echo.Context) error {
	var err error
	return c.Render(`task/edit`, handler.Err(c, err))
}

func Edit(c echo.Context) error {
	var err error
	c.Set(`activeURL`, `/task/index`)
	return c.Render(`task/edit`, handler.Err(c, err))
}

func Delete(c echo.Context) error {
	return c.Redirect(`/task/index`)
}
