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

	"github.com/nging-plugins/frpmanager/application/dbschema"
)

func NewFrpGroup(ctx echo.Context) *FrpGroup {
	return &FrpGroup{
		NgingFrpGroup: dbschema.NewNgingFrpGroup(ctx),
	}
}

type FrpGroup struct {
	*dbschema.NgingFrpGroup
}

func (f *FrpGroup) Exists(name string) (bool, error) {
	return f.NgingFrpGroup.Exists(nil, db.Cond{`name`: name})
}

func (f *FrpGroup) ExistsOther(name string, id uint) (bool, error) {
	return f.NgingFrpGroup.Exists(nil, db.Cond{`name`: name, `id <>`: id})
}
