package handler

import (
	"github.com/admpub/godownloader/service"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/cron"
	"github.com/nging-plugins/dbmanager/pkg/library/dbmanager/driver/mysql"
	dlconfig "github.com/nging-plugins/dlmanager/pkg/library/config"
)

var downloadDir = func() string {
	if len(dlconfig.Get().SavePath) == 0 {
		return service.GetDownloadPath()
	}
	return dlconfig.Get().SavePath
}

func init() {
	mysql.SQLTempDir = downloadDir //将SQL文件缓存到下载目录里面方便管理
	handler.RegisterToGroup(`/db`, func(g echo.RouteRegister) {
		metaHandler := handler.IRegister().MetaHandler

		// dbmanager
		g.Route(`GET,POST`, ``, Manager)

		// account
		g.Route(`GET,POST`, `/account`, metaHandler(echo.H{`name`: `账号列表`}, AccountIndex))
		g.Route(`GET,POST`, `/account_add`, metaHandler(echo.H{`name`: `添加账号`}, AccountAdd))
		g.Route(`GET,POST`, `/account_edit`, metaHandler(echo.H{`name`: `修改账号`}, AccountEdit))
		g.Route(`GET,POST`, `/account_delete`, metaHandler(echo.H{`name`: `删除账号`}, AccountDelete))

		// schema sync
		g.Route(`GET,POST`, `/schema_sync`, metaHandler(echo.H{`name`: `同步方案列表`}, SchemaSync))
		g.Route(`GET,POST`, `/schema_sync_add`, metaHandler(echo.H{`name`: `添加同步方案`}, SchemaSyncAdd))
		g.Route(`GET,POST`, `/schema_sync_edit`, metaHandler(echo.H{`name`: `编辑同步方案`}, SchemaSyncEdit))
		g.Route(`GET`, `/schema_sync_delete`, metaHandler(echo.H{`name`: `删除同步方案`}, SchemaSyncDelete))
		g.Route(`GET`, `/schema_sync_preview`, metaHandler(echo.H{`name`: `预览要同步的项`}, SchemaSyncPreview))
		g.Route(`GET`, `/schema_sync_run`, metaHandler(echo.H{`name`: `执行同步方案`}, SchemaSyncRun))
		g.Route(`GET`, `/schema_sync_log/:id`, metaHandler(echo.H{`name`: `日志列表`}, SchemaSyncLog))
		g.Route(`GET`, `/schema_sync_log_view/:id`, metaHandler(echo.H{`name`: `日志详情`}, SchemaSyncLogView))
		g.Route(`GET`, `/schema_sync_log_delete`, metaHandler(echo.H{`name`: `删除日志`}, SchemaSyncLogDelete))
	})
	cron.Register(`mysql_schema_sync`, SchemaSyncJob, `>mysql_schema_sync:1`, `同步MySQL数据表结构`)
}
