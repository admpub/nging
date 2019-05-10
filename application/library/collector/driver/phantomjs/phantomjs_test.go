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

package phantomjs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = assert.Equal

func TestFectch(t *testing.T) {
	b := New()
	err := b.Start(nil)
	defer b.Close()
	if err != nil {
		panic(err)
	}
	var _b []byte
	_b, err = b.Do(`http://www.admpub.com/`, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(_b))
	//assert.Equal(t, "<h2>安装 Go 第三方包 go-sqlite3</h2>", string(_b))
}
