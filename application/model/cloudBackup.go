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
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/common"
)

func NewCloudBackup(ctx echo.Context) *CloudBackup {
	m := &CloudBackup{
		NgingCloudBackup: dbschema.NewNgingCloudBackup(ctx),
	}
	return m
}

type CloudBackup struct {
	*dbschema.NgingCloudBackup
}

func (s *CloudBackup) check() error {
	ctx := s.Context()
	if len(s.StorageEngine) == 0 {
		s.StorageEngine = StorageEngineS3
	}
	if !CloudBackupStorageEngines.Has(s.StorageEngine) {
		return ctx.NewError(code.InvalidParameter, `存储引擎无效`).SetZone(`storageEngine`)
	}
	if s.StorageEngine == StorageEngineSFTP || s.StorageEngine == StorageEngineS3 {
		if s.DestStorage < 1 {
			return ctx.NewError(code.InvalidParameter, `请选择目标存储账号`).SetZone(`destStorage`)
		}
	}
	s.SourcePath = strings.TrimSpace(s.SourcePath)
	if len(s.SourcePath) == 0 {
		return ctx.NewError(code.InvalidParameter, `请设置源路径`).SetZone(`sourcePath`)
	}
	s.LogDisabled = common.GetBoolFlag(s.LogDisabled)
	if !CloudBackupLogTypes.Has(s.LogType) {
		return ctx.NewError(code.InvalidParameter, `日志类型无效`).SetZone(`logType`)
	}
	s.Disabled = common.GetBoolFlag(s.Disabled)
	s.WaitFillCompleted = common.GetBoolFlag(s.WaitFillCompleted)
	return nil
}

func (s *CloudBackup) Add() (pk interface{}, err error) {
	if err = s.check(); err != nil {
		return nil, err
	}
	return s.NgingCloudBackup.Insert()
}

func (s *CloudBackup) Edit(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	if err = s.check(); err != nil {
		return err
	}
	return s.NgingCloudBackup.Update(mw, args...)
}

func (s *CloudBackup) ListPage(cond *db.Compounds, sorts ...interface{}) ([]*CloudBackupExt, error) {
	rows := []*CloudBackupExt{}
	_, err := common.NewLister(s.NgingCloudBackup, &rows, func(r db.Result) db.Result {
		return r.OrderBy(sorts...)
	}, cond.And()).Paging(s.Context())
	return rows, err
}
