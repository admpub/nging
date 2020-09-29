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

package model

import (
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/model/base"
)

func NewCloudBackup(ctx echo.Context) *CloudBackup {
	m := &CloudBackup{
		NgingCloudBackup: &dbschema.NgingCloudBackup{},
		base:             base.New(ctx),
	}
	m.SetContext(ctx)
	return m
}

type CloudBackup struct {
	*dbschema.NgingCloudBackup
	base *base.Base
}

func (s *CloudBackup) check() error {
	ctx := s.base.Context
	s.SourcePath = strings.TrimSpace(s.SourcePath)
	if len(s.SourcePath) == 0 {
		return ctx.NewError(code.InvalidParameter, ctx.T(`请设置源路径`))
	}
	if s.DestStorage < 1 {
		return ctx.NewError(code.InvalidParameter, ctx.T(`请选择目标存储账号`))
	}
	s.Disabled = common.GetBoolFlag(s.Disabled)
	return nil
}

func (s *CloudBackup) Add() (pk interface{}, err error) {
	if err = s.check(); err != nil {
		return nil, err
	}
	return s.NgingCloudBackup.Add()
}

func (s *CloudBackup) Edit(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	if err = s.check(); err != nil {
		return err
	}
	return s.NgingCloudBackup.Edit(mw, args...)
}

func (s *CloudBackup) ListPage(cond *db.Compounds, sorts ...interface{}) ([]*CloudBackupExt, error) {
	rows := []*CloudBackupExt{}
	_, err := common.NewLister(s.NgingCloudBackup, &rows, func(r db.Result) db.Result {
		return r.Relation(`Storage`, func(sel sqlbuilder.Selector) sqlbuilder.Selector {
			return sel.Columns(`id`, `name`)
		}).OrderBy(sorts...)
	}, cond.And()).Paging(s.base.Context)
	return rows, err
}
