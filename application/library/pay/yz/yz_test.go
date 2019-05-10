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

package yz

import (
	"testing"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func TestCreateQROrder(t *testing.T) {
	client := New(`d4d22e11de5466ca90`, ``, `42135241`)
	data := echo.Store{
		`price`: 0.01,
		`name`:  `test`,
		`ids`:   []string{"1", "2", "3"},
	}
	_, err := client.CreateQROrder(data)
	if err != nil {
		t.Error(err)
	}
	echo.Dump(data)
	var r = echo.Store{}
	com.JSONDecode([]byte(`{
		"a":"b",
		"c":{"d":1,"e":"test"}
	}`), &r)
	v := r.Store("c").Int64("d")
	echo.Dump(v)
}
