package handler

import (
	"runtime"
	"time"

	"github.com/admpub/go-ttlmap"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/sftpmanager"
	"github.com/nging-plugins/sshmanager/application/dbschema"
	"github.com/webx-top/echo/param"
)

var defaultMaxAge = time.Hour * 3
var cachedSFTPClients = ttlmap.New(&ttlmap.Options{
	InitialCapacity: 15,
	OnWillExpire:    nil,
	OnWillEvict: func(key string, item ttlmap.Item) {
		closeCachedClient(item)
	},
})

func init() {
	runtime.SetFinalizer(cachedSFTPClients, func(t *ttlmap.Map) error {
		cachedSFTPClients.Drain()
		return nil
	})
}

func closeCachedClient(item ttlmap.Item) {
	mgr := item.Value().(*sftpmanager.SftpManager)
	mgr.Close()
}

func getCachedSFTPClient(sshUser *dbschema.NgingSshUser) (mgr *sftpmanager.SftpManager, err error) {
	key := param.AsString(sshUser.Id)
	var item ttlmap.Item
	item, err = cachedSFTPClients.Get(key)
	if err == nil {
		mgr = item.Value().(*sftpmanager.SftpManager)
	} else {
		cfg := sftpConfig(sshUser)
		mgr = sftpmanager.New(sftpmanager.DefaultConnector, &cfg, config.FromFile().Sys.EditableFileMaxBytes())
		err = cachedSFTPClients.Set(key, ttlmap.NewItem(mgr, ttlmap.WithTTL(defaultMaxAge)), nil)
	}
	return
}

func deleteCachedSFTPClient(id uint) (err error) {
	key := param.AsString(id)
	var item ttlmap.Item
	item, err = cachedSFTPClients.Delete(key)
	if err == nil {
		closeCachedClient(item)
	}
	return
}
