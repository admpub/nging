package caddy_webdav

import (
	"fmt"

	plugin "github.com/caddy-plugins/webdav"
	"github.com/webx-top/db"
	"github.com/webx-top/echo/defaults"
	"golang.org/x/net/webdav"

	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/s3manager/s3client"
	"github.com/admpub/nging/application/library/s3manager/s3webdav"
	"github.com/admpub/nging/application/model"
)

func init() {
	plugin.FSGenerator = plugin.LazyloadFS(FSGenerator)
}

func FSGenerator(scope string, options map[string]string) webdav.FileSystem {
	typ, ok := options[`arg1`]
	if !ok {
		return webdav.Dir(scope)
	}
	if typ != `id` {
		return webdav.Dir(scope)
	}
	id, ok := options[`arg2`]
	if !ok {
		return webdav.Dir(scope)
	}
	ctx := defaults.NewMockContext()
	m := model.NewCloudStorage(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = fmt.Errorf(`cannot find the cloud storage account record with ID "%v". `+"\n"+`找不到ID为"%v"的云存储账号记录。`, id, id)
		}
		panic(err)
	}
	mgr, err := s3client.New(m.NgingCloudStorage, config.DefaultConfig.Sys.EditableFileMaxBytes)
	if err != nil {
		panic(err)
	}
	return s3webdav.New(mgr, false, ``)
}
