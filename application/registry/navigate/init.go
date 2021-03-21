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

package navigate

import "github.com/webx-top/echo"

var (
	topNavURLs = map[string]int{}

	// EmptyList 空菜单列表
	EmptyList = List{}

	//LeftNavigate 左边导航菜单
	LeftNavigate = &List{}

	//TopNavigate 顶部导航菜单
	TopNavigate = &List{}
)

func TopNavURLs() map[string]int {
	return topNavURLs
}

func init() {
	echo.On(`beforeRun`, func(_ echo.H) error {
		ProjectInitURLsIdent()
		for index, urlPath := range TopNavigate.FullPath(``) {
			topNavURLs[urlPath] = index
		}
		return nil
	})
}
