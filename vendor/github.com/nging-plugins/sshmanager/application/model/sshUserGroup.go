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

	"github.com/nging-plugins/sshmanager/application/dbschema"
)

func NewSshUserGroup(ctx echo.Context) *SshUserGroup {
	return &SshUserGroup{
		NgingSshUserGroup: dbschema.NewNgingSshUserGroup(ctx),
	}
}

type SshUserGroup struct {
	*dbschema.NgingSshUserGroup
}

func (f *SshUserGroup) Exists(name string) (bool, error) {
	return f.NgingSshUserGroup.Exists(nil, db.Cond{`name`: name})
}

func (f *SshUserGroup) ExistsOther(name string, id uint) (bool, error) {
	return f.NgingSshUserGroup.Exists(nil, db.Cond{`name`: name, `id <>`: id})
}
