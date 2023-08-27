package handler

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/library/cron"
	"github.com/nging-plugins/dbmanager/application/library/dbmanager"
	"github.com/nging-plugins/dbmanager/application/library/dbmanager/driver"
	"github.com/nging-plugins/dbmanager/application/library/dbmanager/driver/mysql"
	"github.com/nging-plugins/dbmanager/application/library/dbmanager/driver/mysql/utils"
	"github.com/nging-plugins/dbmanager/application/model"
)

func mysqlBackup(id string) cron.Runner {
	parts := strings.SplitN(id, `?`, 2)
	var database string
	keepN := 30
	if len(parts) > 1 {
		params, err := url.ParseQuery(parts[1])
		if err != nil {
			log.Error(err)
		}
		database = params.Get(`database`)
		keepNum := params.Get(`keepN`)
		if len(keepNum) > 0 {
			keepN = param.AsInt(keepNum)
		}
	}
	accountID := param.AsUint(parts[0])
	return func(timeout time.Duration) (out string, runingErr string, onErr error, isTimeout bool) {
		ctx := defaults.NewMockContext()
		m := model.NewDbAccount(ctx)
		err := m.Get(nil, db.Cond{`id`: accountID})
		if err != nil {
			onErr = cron.ErrFailure
			runingErr = err.Error()
			return
		}
		auth := &driver.DbAuth{
			Driver:    `mysql`,
			AccountID: accountID,
			Username:  m.User,
			Host:      m.Host,
			Db:        m.Name,
		}
		if len(database) > 0 {
			m.Name = database
			auth.Db = database
		}
		mgr := dbmanager.New(ctx, auth)
		err, succeed := authentication(mgr, m)
		if err != nil {
			onErr = cron.ErrFailure
			runingErr = err.Error()
			return
		}
		if !succeed {
			onErr = cron.ErrFailure
			runingErr = ctx.T(`登录数据库失败`)
			return
		}
		defer mgr.Run(`logout`)
		ctx.Request().SetMethod(echo.POST)
		ctx.Request().Form().Set(`all`, `1`)
		ctx.Request().Form().Set(`type`, `struct`)
		ctx.Request().Form().Add(`type`, `data`)
		ctx.Request().Header().Set(echo.HeaderAccept, echo.MIMEApplicationJSONCharsetUTF8)
		err = mgr.Run(`export`)
		if err != nil {
			onErr = cron.ErrFailure
			runingErr = err.Error()
			return
		}
		data := ctx.Data()
		if !data.GetCode().Is(code.Success) {
			onErr = cron.ErrFailure
			runingErr = param.AsString(data.GetInfo())
			return
		}

		out = param.AsString(data.GetInfo()) + `，下载地址：` + data.GetURL()

		if keepN > 0 {
			clearHistoryBackup(auth.Db, keepN)
		}

		return
	}
}

func clearHistoryBackup(dbName string, keepN int) {
	saveDir := mysql.TempDir(utils.OpExport)
	dbSaveDir := filepath.Join(saveDir, dbName)
	if !com.FileExists(dbSaveDir) {
		return
	}

	matches, err := filepath.Glob(dbSaveDir + echo.FilePathSeparator + `*.zip`)
	if err != nil {
		log.Error(err.Error())
		return
	}
	num := len(matches)
	if num <= keepN {
		return
	}
	for _, file := range matches[0 : num-keepN] {
		err = os.Remove(file)
		if err != nil {
			log.Errorf(`failed to remove %v: %v`, file, err.Error())
		} else {
			log.Errorf(`successfully remove %v`, file)
		}
		os.Remove(file + `.txt`)
	}
}
