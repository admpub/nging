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

package model

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v5/application/dbschema"
)

func NewCloudBackupLog(ctx echo.Context) *CloudBackupLog {
	m := &CloudBackupLog{
		NgingCloudBackupLog: dbschema.NewNgingCloudBackupLog(ctx),
	}
	return m
}

type CloudBackupLog struct {
	*dbschema.NgingCloudBackupLog
}

func (s *CloudBackupLog) check() error {
	ctx := s.Context()
	if !CloudBackupTypes.Has(s.BackupType) {
		return ctx.NewError(code.InvalidParameter, `备份类型无效`).SetZone(`backupType`)
	}
	if !CloudBackupStatuses.Has(s.Status) {
		return ctx.NewError(code.InvalidParameter, `备份状态无效`).SetZone(`status`)
	}
	return nil
}

func (s *CloudBackupLog) Add() (pk interface{}, err error) {
	if err = s.check(); err != nil {
		return nil, err
	}
	return s.NgingCloudBackupLog.Insert()
}

func (s *CloudBackupLog) Edit(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	if err = s.check(); err != nil {
		return err
	}
	return s.NgingCloudBackupLog.Update(mw, args...)
}
