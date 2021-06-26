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

package perm

import (
	"sync"

	"github.com/admpub/nging/application/registry/navigate"
	"github.com/webx-top/echo"
)

var (
	navTreeCached *Map
	navTreeOnce   sync.Once
)

func initNavTreeCached() {
	if navTreeCached != nil {
		return
	}
	navTreeCached = NewMap()
	for _, project := range navigate.ProjectListAll() {
		if project == nil {
			continue
		}
		navTreeCached.Import(project.NavList)
	}
	navTreeCached.Import(navigate.TopNavigate)
}

func NavTreeCached() *Map {
	navTreeOnce.Do(initNavTreeCached)
	return navTreeCached
}

func init() {
	echo.On(`beforeRun`, func(_ echo.H) error {
		NavTreeCached()
		return nil
	})
}
