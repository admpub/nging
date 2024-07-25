package bolt

import (
	"path/filepath"
	"reflect"
	"time"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/session/engine"
	"github.com/webx-top/echo/middleware/session/engine/bolt"
	"github.com/webx-top/echo/middleware/session/engine/cookie"
	"github.com/webx-top/echo/param"
)

func init() {
	config.RegisterSessionStore(`bolt`, `BoltDB存储`, initSessionStoreBolt)
}

var sessionStoreBoltOptions *bolt.BoltOptions

func initSessionStoreBolt(_ *config.Config, cookieOptions *cookie.CookieOptions, sessionConfig param.Store) (changed bool, err error) {
	boltOptions := &bolt.BoltOptions{
		File:          sessionConfig.String(`savePath`),
		KeyPairs:      cookieOptions.KeyPairs,
		BucketName:    sessionConfig.String(`bucketName`),
		MaxLength:     sessionConfig.Int(`maxLength`),
		CheckInterval: time.Duration(sessionConfig.Int64(`checkInterval`)) * time.Second,
	}
	if len(boltOptions.BucketName) == 0 {
		boltOptions.BucketName = `sessions`
	}
	if len(boltOptions.File) == 0 {
		boltOptions.File = filepath.Join(echo.Wd(), `data`, `cache`, `sessions`, `bolt`)
	}
	if com.IsDir(boltOptions.File) {
		boltOptions.File = filepath.Join(boltOptions.File, `bolt`)
	}
	if sessionStoreBoltOptions == nil || !engine.Exists(`bolt`) || !reflect.DeepEqual(boltOptions, sessionStoreBoltOptions) {
		bolt.RegWithOptions(boltOptions)
		sessionStoreBoltOptions = boltOptions
		changed = true
	}
	return
}
