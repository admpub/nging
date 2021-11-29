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
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v3/application/dbschema"
)

func NewCollectorGroup(ctx echo.Context) *CollectorGroup {
	return &CollectorGroup{
		NgingCollectorGroup: dbschema.NewNgingCollectorGroup(ctx),
	}
}

type CollectorPageAndGroup struct {
	*dbschema.NgingCollectorPage
	Group *dbschema.NgingCollectorGroup `db:"-,relation=id:group_id|gtZero"`
}

type CollectorExportAndGroup struct {
	*dbschema.NgingCollectorExport
	Group *dbschema.NgingCollectorGroup
}

type CollectorGroup struct {
	*dbschema.NgingCollectorGroup
}
