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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/com"
)

func TestCaddy(t *testing.T) {
	l := NewAccessLog(nil)
	err := l.parseCaddy(`::1 - - [06/Nov/2018:00:10:59 +0800] "GET /assets/img/bg-middle.jpg HTTP/1.1" 304 0 "http://nging.coscms.com:2008/assets/css/index.css" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:63.0) Gecko/20100101 Firefox/63.0" 172µs`, ``)
	if err != nil {
		panic(err)
	}
	com.Dump(l)

	l2 := NewAccessLog(nil)
	err = l2.parseCaddy(`::1 - - [06/Nov/2018:00:10:59 +0800] "GET /assets/img/bg-middle.jpg HTTP/1.1" 304 0 "http://nging.coscms.com:2008/assets/css/index.css" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:63.0) Gecko/20100101 Firefox/63.0" 172µs `, `{remote} - {user} [{when}] "{method} {uri} {proto}" {status} {size} "{>Referer}" "{>User-Agent}" {latency} `)
	if err != nil {
		panic(err)
	}
	com.Dump(l2)
	assert.Equal(t, l, l2)
}
