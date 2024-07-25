package mysql

import (
	"reflect"
	"time"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/encoding/dbconfig"
	"github.com/webx-top/echo/middleware/session/engine"
	"github.com/webx-top/echo/middleware/session/engine/cookie"
	"github.com/webx-top/echo/middleware/session/engine/mysql"
	"github.com/webx-top/echo/param"
)

func init() {
	config.RegisterSessionStore(`mysql`, `MySQL存储`, initSessionStoreMySQL)
}

var sessionStoreMySQLOptions *mysql.Options

func initSessionStoreMySQL(c *config.Config, cookieOptions *cookie.CookieOptions, sessionConfig param.Store) (changed bool, err error) {
	mysqlOptions := &mysql.Options{
		Config: dbconfig.Config{
			User:    sessionConfig.String(`user`),
			Pass:    sessionConfig.String(`password`),
			Name:    sessionConfig.String(`database`),
			Host:    sessionConfig.String(`host`),
			Port:    sessionConfig.String(`port`),
			Charset: sessionConfig.String(`charset`),
			Prefix:  sessionConfig.String(`prefix`),
			Options: map[string]string{},
		},
		Table:         sessionConfig.String(`table`),
		KeyPairs:      cookieOptions.KeyPairs,
		MaxAge:        sessionConfig.Int(`maxAge`),
		MaxLength:     sessionConfig.Int(`maxLength`),
		CheckInterval: time.Duration(sessionConfig.Int64(`checkInterval`)) * time.Second,
		MaxReconnect:  sessionConfig.Int(`maxReconnect`),
	}
	for k, v := range c.DB.Options {
		mysqlOptions.Config.Options[k] = v
	}
	for k, v := range sessionConfig.GetStore(`options`) {
		mysqlOptions.Config.Options[k] = param.AsString(v)
	}
	if len(mysqlOptions.Config.User) == 0 {
		mysqlOptions.Config.User = c.DB.User
		mysqlOptions.Config.Pass = c.DB.Password
	}
	if len(mysqlOptions.Config.Name) == 0 {
		mysqlOptions.Config.Name = c.DB.Database
	}
	if len(mysqlOptions.Config.Host) == 0 {
		mysqlOptions.Config.Host = c.DB.Host
	}
	if len(mysqlOptions.Config.Charset) == 0 {
		mysqlOptions.Config.Charset = c.DB.Charset()
	}
	if len(mysqlOptions.Config.Prefix) == 0 {
		mysqlOptions.Config.Prefix = c.DB.Prefix
	}
	if len(mysqlOptions.Config.Port) == 0 {
		_, port := com.SplitHostPort(mysqlOptions.Config.Host)
		if len(port) == 0 {
			mysqlOptions.Config.Port = `3306`
		}
	}
	if len(mysqlOptions.Table) == 0 {
		mysqlOptions.Table = mysqlOptions.Config.Prefix + `sessions`
	} else {
		mysqlOptions.Table = mysqlOptions.Config.Prefix + mysqlOptions.Table
	}
	if mysqlOptions.MaxReconnect <= 0 {
		mysqlOptions.MaxReconnect = 30
	}
	if sessionStoreMySQLOptions == nil || !engine.Exists(`mysql`) || !reflect.DeepEqual(mysqlOptions, sessionStoreMySQLOptions) {
		mysql.RegWithOptions(mysqlOptions)
		sessionStoreMySQLOptions = mysqlOptions
		changed = true
	}
	return
}
