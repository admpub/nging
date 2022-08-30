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

func NewFtpUserGroup(ctx echo.Context) *FtpUserGroup {
	return &FtpUserGroup{
		NgingFtpUserGroup: dbschema.NewNgingFtpUserGroup(ctx),
	}
}

type FtpUserGroup struct {
	*dbschema.NgingFtpUserGroup
}

func (f *FtpUserGroup) Exists(name string) (bool, error) {
	return f.NgingFtpUserGroup.Exists(nil, db.Cond{`name`: name})
}

func (f *FtpUserGroup) ExistsOther(name string, id uint) (bool, error) {
	return f.NgingFtpUserGroup.Exists(nil, db.Cond{`name`: name, `id <>`: id})
}

func (f *FtpUserGroup) ListByActive(page int, size int) (func() int64, []*dbschema.NgingFtpUserGroup, error) {
	count, err := f.List(nil, nil, page, size, db.Cond{`disabled`: `N`})
	if err == nil {
		return count, f.Objects(), err
	}
	return count, nil, err
}
