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

package collector

import (
	"sort"

	"github.com/webx-top/echo"
)

type Browser interface {
	Start(echo.Store) error
	Close() error
	Name() string
	Description() string
	Transcoded() bool
	Do(pageURL string, data echo.Store) ([]byte, error)
}

var (
	Browsers = map[string]Browser{}
)

func BrowserKeys() []string {
	keys := make([]string, 0)
	for key := range Browsers {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
