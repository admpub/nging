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
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/nging-plugins/ftpmanager/application/dbschema"
)

func NewFtpPermission(ctx echo.Context) *FtpPermission {
	return &FtpPermission{
		NgingFtpPermission: dbschema.NewNgingFtpPermission(ctx),
	}
}

type FtpPermission struct {
	*dbschema.NgingFtpPermission
}

func (f *FtpPermission) Exists(targetType string, targetID uint) (bool, error) {
	return f.NgingFtpPermission.Exists(nil, db.And(
		db.Cond{`target_type`: targetType},
		db.Cond{`target_id`: targetID},
	))
}

func (f *FtpPermission) GetByTarget(targetType string, targetID uint) error {
	err := f.NgingFtpPermission.Get(nil, db.And(
		db.Cond{`target_type`: targetType},
		db.Cond{`target_id`: targetID},
	))
	return err
}

func (f *FtpPermission) DeleteByTarget(targetType string, targetID uint) error {
	err := f.NgingFtpPermission.Delete(nil, db.And(
		db.Cond{`target_type`: targetType},
		db.Cond{`target_id`: targetID},
	))
	return err
}

func (f *FtpPermission) Save() (pk interface{}, err error) {
	row := dbschema.NewNgingFtpPermission(f.Context())
	err = row.Get(nil, db.And(
		db.Cond{`target_type`: f.TargetType},
		db.Cond{`target_id`: f.TargetId},
	))
	if err != nil {
		if err != db.ErrNoMoreRows {
			return
		}
		pk, err = f.NgingFtpPermission.Insert()
		return
	}
	kvset := echo.H{
		`permission`: f.Permission,
	}
	err = f.NgingFtpPermission.UpdateFields(nil, kvset, `id`, row.Id)
	return
}
