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
package collector

import (
	"github.com/admpub/nging/application/handler"
	"github.com/webx-top/echo"
)

func init() {
	handler.RegisterToGroup(`/collector`, func(g *echo.Group) {
		g.Route(`GET,POST`, `/task`, Task)
		g.Route(`GET,POST`, `/rule`, Rule)
		g.Route(`GET,POST`, `/rule_add`, RuleAdd)
		g.Route(`GET,POST`, `/rule_edit`, RuleEdit)
		g.Route(`GET,POST`, `/rule_delete`, RuleDelete)
	})
}

func Task(c echo.Context) error {
	var err error
	return c.Render(`collector/task`, handler.Err(c, err))
}

func Rule(c echo.Context) error {
	var err error
	return c.Render(`collector/rule`, handler.Err(c, err))
}

func RuleAdd(c echo.Context) error {
	var err error
	return c.Render(`collector/rule_edit`, handler.Err(c, err))
}

func RuleEdit(c echo.Context) error {
	var err error
	return c.Render(`collector/rule_edit`, handler.Err(c, err))
}

func RuleDelete(c echo.Context) error {
	return c.Redirect(`/collector/rule`)
}
