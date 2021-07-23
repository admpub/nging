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

package export

import (
	"fmt"
	"sync"

	"github.com/admpub/log"
	"github.com/admpub/nging/v3/application/dbschema"
	"github.com/admpub/nging/v3/application/library/collector/exec"
	"github.com/admpub/nging/v3/application/library/collector/sender"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/db/mysql"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
)

var connections = sync.Map{}

func CloseAllDBConn() {
	connections.Range(func(key, conn interface{}) bool {
		conn.(sqlbuilder.Database).Close()
		connections.Delete(key)
		return true
	})
}

func DBConn(dsn interface{}) (db sqlbuilder.Database, err error) {
	key := fmt.Sprintf(`%#v`, dsn)
	conn, ok := connections.Load(key)
	if ok {
		db, ok = conn.(sqlbuilder.Database)
	}
	if !ok {
		var settings mysql.ConnectionURL
		if ds, ok := dsn.(string); ok {
			settings, err = mysql.ParseURL(ds)
			if err != nil {
				return
			}
		} else {
			settings = dsn.(mysql.ConnectionURL)
		}
		db, err = mysql.Open(settings)
		if err != nil {
			return db, err
		}
		connections.Store(key, db)
	}
	return
}

func Export(pageID uint, result *exec.Recv, collected echo.Store, noticeSender sender.Notice) error {
	// 导出数据
	exportM := &dbschema.NgingCollectorExport{}
	_cnt, _err := exportM.ListByOffset(nil, nil, 0, -1, db.And(
		db.Cond{`page_id`: pageID},
		db.Cond{`disabled`: `N`},
		db.Cond{`mapping`: db.NotEq(``)},
		db.Cond{`dest`: db.NotEq(``)},
	))
	if _err == nil && _cnt() > 0 {
		for _, expc := range exportM.Objects() {
			mappings := NewMappings()
			_err = com.JSONDecode(engine.Str2bytes(expc.Mapping), mappings)
			if _err == nil {
				_err = mappings.Export(result, collected, expc, noticeSender)
			}
			if _err != nil {
				if sendErr := noticeSender(_err.Error(), 0); sendErr != nil {
					return sendErr
				}
				log.Error(_err)
				continue
			}
		}
	}
	if _err != nil {
		if sendErr := noticeSender(_err.Error(), 0); sendErr != nil {
			return sendErr
		}
	}
	return _err
}
