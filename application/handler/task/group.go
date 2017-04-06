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

func Group(c echo.Context) error {
	var err error
	return c.Render(`task/group`, handler.Err(c, err))
}

func GroupAdd(c echo.Context) error {
	var err error

	c.Set(`activeURL`, `/task/group`)
	return c.Render(`task/group_edit`, handler.Err(c, err))
}

func GroupEdit(c echo.Context) error {
	var err error
	c.Set(`activeURL`, `/task/group`)
	return c.Render(`task/group_edit`, handler.Err(c, err))
}

func GroupDelete(c echo.Context) error {
	return c.Redirect(`/task/group`)
}
