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

package cloud

import (
	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/config/startup"
	"github.com/admpub/nging/v5/application/model"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
)

func init() {
	handler.RegisterToGroup(`/cloud`, func(g echo.RouteRegister) {
		g.Route(`GET,POST`, `/storage`, StorageIndex)
		g.Route(`GET,POST`, `/storage_add`, StorageAdd)
		g.Route(`GET,POST`, `/storage_edit`, StorageEdit)
		g.Route(`GET,POST`, `/storage_delete`, StorageDelete)
		g.Route(`GET,POST`, `/storage_file`, StorageFile)

		g.Route(`GET,POST`, `/backup`, BackupConfigList)
		g.Route(`GET,POST`, `/backup_add`, BackupConfigAdd)
		g.Route(`GET,POST`, `/backup_edit`, BackupConfigEdit)
		g.Route(`GET,POST`, `/backup_delete`, BackupConfigDelete)
		g.Route(`GET,POST`, `/backup_start`, BackupStart)
		g.Route(`GET,POST`, `/backup_stop`, BackupStop)
		g.Route(`GET,POST`, `/backup_log`, Log)
		g.Route(`GET,POST`, `/backup_log_delete`, LogDelete)
	})

	startup.OnBefore(`web`, func() {
		if !config.IsInstalled() {
			return
		}
		ctx := defaults.NewMockContext()
		m := model.NewCloudBackup(ctx)
		_, err := m.ListByOffset(nil, nil, 0, -1, db.Cond{`disabled`: `N`})
		if err != nil {
			log.Error(err)
			return
		}
		for _, row := range m.Objects() {
			err = monitorBackupStart(*row)
			if err != nil {
				log.Error(err)
			}
		}
	})
}
