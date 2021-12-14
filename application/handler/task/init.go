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

package task

import (
	"github.com/admpub/nging/v4/application/handler"
	"github.com/webx-top/echo"
)

func init() {
	handler.RegisterToGroup(`/task`, func(g echo.RouteRegister) {
		metaHandler := handler.IRegister().MetaHandler
		g.Route(`GET,POST`, `/index`, metaHandler(echo.H{`name`: `任务列表`}, Index))
		g.Route(`GET,POST`, `/add`, metaHandler(echo.H{`name`: `添加任务`}, Add))
		g.Route(`GET,POST`, `/edit`, metaHandler(echo.H{`name`: `修改任务`}, Edit))
		g.Route(`GET,POST`, `/delete`, metaHandler(echo.H{`name`: `删除任务`}, Delete))
		g.Route(`GET,POST`, `/start`, metaHandler(echo.H{`name`: `启动任务`}, Start))
		g.Route(`GET,POST`, `/pause`, metaHandler(echo.H{`name`: `暂停任务`}, Pause))
		g.Route(`GET,POST`, `/run`, metaHandler(echo.H{`name`: `立即执行`}, Run))
		g.Route(`GET,POST`, `/exit`, metaHandler(echo.H{`name`: `退出任务`}, Exit))
		g.Route(`GET,POST`, `/start_history`, metaHandler(echo.H{`name`: `启动历史任务`}, StartHistory))
		g.Route(`GET,POST`, `/group`, metaHandler(echo.H{`name`: `任务分组列表`}, Group))
		g.Route(`GET,POST`, `/group_add`, metaHandler(echo.H{`name`: `添加分组`}, GroupAdd))
		g.Route(`GET,POST`, `/group_edit`, metaHandler(echo.H{`name`: `修改分组`}, GroupEdit))
		g.Route(`GET,POST`, `/group_delete`, metaHandler(echo.H{`name`: `删除分组`}, GroupDelete))
		g.Route(`GET,POST`, `/log`, metaHandler(echo.H{`name`: `任务日志列表`}, Log))
		g.Route(`GET,POST`, `/log_view/:id`, metaHandler(echo.H{`name`: `任务日志详情`}, LogView))
		g.Route(`GET,POST`, `/log_delete`, metaHandler(echo.H{`name`: `删除任务日志`}, LogDelete))
		g.Route(`GET,POST`, `/email_test`, metaHandler(echo.H{`name`: `测试E-mail`}, EmailTest))
	})
}
