package handler

import (
	"runtime"
	"time"

	"github.com/admpub/go-ttlmap"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/sftpmanager"
	"github.com/nging-plugins/sshmanager/application/dbschema"
	"github.com/webx-top/echo/param"
	"golang.org/x/sync/singleflight"
)

var defaultMaxAge = time.Hour * 3
var cachedSFTPClients = ttlmap.New(&ttlmap.Options{
	InitialCapacity: 15,
	OnWillExpire:    nil,
	OnWillEvict: func(key string, item ttlmap.Item) {
		closeCachedClient(item)
	},
})
var sg singleflight.Group

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
	var v interface{}
	v, err, _ = sg.Do(key, func() (interface{}, error) {
		item, err := cachedSFTPClients.Get(key)
		if err == nil {
			mgr = item.Value().(*sftpmanager.SftpManager)
			return mgr, nil
		}
		cfg := sftpConfig(sshUser)
		mgr = sftpmanager.New(sftpmanager.DefaultConnector, &cfg, config.FromFile().Sys.EditableFileMaxBytes())
		mgr.Client()
		if mgr.ConnError() != nil {
			err = mgr.ConnError()
			return mgr, err
		}
		err = cachedSFTPClients.Set(key, ttlmap.NewItem(mgr, ttlmap.WithTTL(defaultMaxAge)), nil)
		return mgr, err
	})
	if err != nil {
		return
	}
	mgr = v.(*sftpmanager.SftpManager)
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
