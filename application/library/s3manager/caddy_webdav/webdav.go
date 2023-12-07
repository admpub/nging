/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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

package caddy_webdav

import (
	"fmt"
	"time"

	plugin "github.com/caddy-plugins/webdav"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"
	"golang.org/x/net/webdav"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/s3manager"
	"github.com/admpub/nging/v5/application/library/s3manager/s3client"
	"github.com/admpub/nging/v5/application/library/s3manager/s3webdav"
	"github.com/admpub/nging/v5/application/model"
)

func init() {
	plugin.FSGenerator = plugin.LazyloadFS(FSGenerator)
	dbschema.DBI.On(`w+`, func(model factory.Model, editColumns ...string) error {
		m := model.(*dbschema.NgingCloudStorage)
		managers.Delete(param.AsString(m.Id))
		return nil
	}, dbschema.WithPrefix(`nging_cloud_storage`))
}

type MgrCached struct {
	T time.Time
	M *s3manager.S3Manager
}

var managers = param.NewMap()

// webdav / id 1 {}
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
	cached := managers.Get(id)
	mgr, ok := cached.(*MgrCached)
	if !ok || time.Since(mgr.T).Seconds() > 3600 {
		var err error
		mgr, err = genCache(id)
		if err != nil {
			panic(err)
		}
	}
	return s3webdav.New(mgr.M, scope, false, ``)
}

func genCache(id string) (*MgrCached, error) {
	ctx := defaults.NewMockContext()
	m := model.NewCloudStorage(ctx)
	err := m.Get(nil, `id`, id)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = fmt.Errorf(`cannot find the cloud storage account record with ID "%v". `+"\n"+`找不到ID为"%v"的云存储账号记录。`, id, id)
		}
		return nil, err
	}
	mgr := s3client.New(m.NgingCloudStorage, config.FromFile().Sys.EditableFileMaxBytes())
	_, err = mgr.Connect()
	if err != nil {
		return nil, err
	}
	c := &MgrCached{T: time.Now(), M: mgr}
	managers.Set(id, c)
	return c, nil
}
