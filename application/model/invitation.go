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

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/model/base"
)

func NewInvitation(ctx echo.Context) *Invitation {
	return &Invitation{
		NgingCodeInvitation: &dbschema.NgingCodeInvitation{},
		Base:                base.New(ctx),
	}
}

type Invitation struct {
	*dbschema.NgingCodeInvitation
	*base.Base
}

func (u *Invitation) Exists(code string) (bool, error) {
	return u.NgingCodeInvitation.Exists(nil, db.Cond{`code`: code})
}

func (u *Invitation) Exists2(code string, excludeID uint) (bool, error) {
	return u.NgingCodeInvitation.Exists(nil, db.And(
		db.Cond{`code`: code},
		db.Cond{`id`: db.NotEq(excludeID)},
	))
}
