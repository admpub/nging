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

package file

import (
	"github.com/admpub/nging/application/registry/navigate"
)

var topNavigate = navigate.List{
	{
		DisplayOnMenu: true,
		Name:          `附件管理`,
		Action:        `file/list`,
	},
	{
		DisplayOnMenu: false,
		Name:          `删除附件`,
		Action:        `file/delete/:id`,
	},
}

func init() {
	navigate.TopNavigate.ChildrenBy(0).Add(-1, topNavigate...)
}
