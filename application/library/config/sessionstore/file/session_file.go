package file

import (
	"path/filepath"
	"reflect"
	"time"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/session/engine"
	"github.com/webx-top/echo/middleware/session/engine/cookie"
	"github.com/webx-top/echo/middleware/session/engine/file"
	"github.com/webx-top/echo/param"
)

func init() {
	config.RegisterSessionStore(`file`, `文件存储`, initSessionStoreFile)
}

var sessionStoreFileOptions *file.FileOptions

func initSessionStoreFile(_ *config.Config, cookieOptions *cookie.CookieOptions, sessionConfig param.Store) (changed bool, err error) {
	fileOptions := &file.FileOptions{
		SavePath:      sessionConfig.String(`savePath`),
		KeyPairs:      cookieOptions.KeyPairs,
		MaxAge:        sessionConfig.Int(`maxAge`),
		MaxLength:     sessionConfig.Int(`maxLength`),
		CheckInterval: time.Duration(sessionConfig.Int64(`checkInterval`)) * time.Second,
	}
	if len(fileOptions.SavePath) == 0 {
		fileOptions.SavePath = filepath.Join(echo.Wd(), `data`, `cache`, `sessions`)
	}
	if sessionStoreFileOptions == nil || !engine.Exists(`file`) || !reflect.DeepEqual(fileOptions, sessionStoreFileOptions) {
		file.RegWithOptions(fileOptions)
		sessionStoreFileOptions = fileOptions
		changed = true
	}
	return
}
