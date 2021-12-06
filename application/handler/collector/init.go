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

package collector

import (
	"github.com/admpub/nging/v3/application/handler"
	"github.com/admpub/nging/v3/application/library/cron"
	"github.com/webx-top/echo"
)

func init() {
	handler.RegisterToGroup(`/collector`, func(g echo.RouteRegister) {
		metaHandler := handler.IRegister().MetaHandler
		g.Route(`GET,POST`, `/export`, metaHandler(echo.H{`name`: `导出管理`}, Export))
		g.Route(`GET,POST`, `/export_log`, metaHandler(echo.H{`name`: `日子列表`}, ExportLog))
		g.Route(`GET,POST`, `/export_log_view/:id`, metaHandler(echo.H{`name`: `日志详情`}, ExportLogView))
		g.Route(`GET,POST`, `/export_log_delete`, metaHandler(echo.H{`name`: `删除日志`}, ExportLogDelete))
		g.Route(`GET,POST`, `/export_add`, metaHandler(echo.H{`name`: `添加导出规则`}, ExportAdd))
		g.Route(`GET,POST`, `/export_edit`, metaHandler(echo.H{`name`: `修改导出规则`}, ExportEdit))
		g.Route(`GET,POST`, `/export_edit_status`, metaHandler(echo.H{`name`: `修改导出规则`}, ExportEditStatus))
		g.Route(`GET,POST`, `/export_delete`, metaHandler(echo.H{`name`: `删除导出规则`}, ExportDelete))
		g.Route(`GET,POST`, `/history`, metaHandler(echo.H{`name`: `历史记录`}, History))
		g.Route(`GET,POST`, `/history_view`, metaHandler(echo.H{`name`: `查看历史内容`}, HistoryView))
		g.Route(`GET,POST`, `/history_delete`, metaHandler(echo.H{`name`: `删除历史记录`}, HistoryDelete))
		g.Route(`GET,POST`, `/rule`, metaHandler(echo.H{`name`: `规则列表`}, Rule))
		g.Route(`GET,POST`, `/rule_add`, metaHandler(echo.H{`name`: `添加规则`}, RuleAdd))
		g.Route(`GET,POST`, `/rule_edit`, metaHandler(echo.H{`name`: `修改规则`}, RuleEdit))
		g.Route(`GET,POST`, `/rule_delete`, metaHandler(echo.H{`name`: `删除规则`}, RuleDelete))
		g.Route(`GET,POST`, `/rule_collect`, metaHandler(echo.H{`name`: `采集`}, RuleCollect))
		g.Route(`GET,POST`, `/group`, metaHandler(echo.H{`name`: `任务分组列表`}, Group))
		g.Route(`GET,POST`, `/group_add`, metaHandler(echo.H{`name`: `添加分组`}, GroupAdd))
		g.Route(`GET,POST`, `/group_edit`, metaHandler(echo.H{`name`: `修改分组`}, GroupEdit))
		g.Route(`GET,POST`, `/group_delete`, metaHandler(echo.H{`name`: `删除分组`}, GroupDelete))
		g.Route(`GET,POST`, `/regexp_test`, metaHandler(echo.H{`name`: `测试正则表达式`}, RegexpTest))
	})
	cron.Register(`collect_page`, CollectPageJob, `>collect_page:1`, `网页采集`)
}
