/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package redis

import (
	"github.com/admpub/nging/application/library/dbmanager/driver"
	"github.com/gomodule/redigo/redis"
	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/mysql"
	"github.com/webx-top/echo"
	"github.com/webx-top/pagination"
	"strings"
)

type redis struct {
	conn redis.Conn
}

func (r *redis) Name() string {
	return `Redis`
}

func (r *redis) Init(ctx echo.Context, auth *driver.DbAuth) {
	m.BaseDriver = driver.NewBaseDriver()
	m.BaseDriver.Init(ctx, auth)
	m.Set(`supportSQL`, false)
}

func (r *redis) IsSupported(operation string) bool {
	return true
}

func (r *redis) Login() (err error) {
	/*
	  Scheme syntax:
	  Example: redis://user:secret@localhost:6379/0?foo=bar&qux=baz

	  This scheme uses a profile of the RFC 3986 generic URI syntax.
	  All URI fields after the scheme are optional.
	  The "userinfo" field uses the traditional "user:password" format.

	  Expressed using RFC 5234 ABNF, the "path" grammar production from
	  RFC 3986 is overridden as follows:
	    path         = [ path-slashed ]
	                 ; path is optional
	    path-slashed = "/" [ db-number ]
	                 ; exactly zero or one path segments
	    db-number    = "0" / nz-num
	                 ; nonnegative decimal integer with no leading zeros
	    nz-num       = NZDIGIT *DIGIT
	                 ; positive decimal integer with no leading zeros
	    NZDIGIT      = %x31-39
	                 ; the digits 1-9
	*/
	if len(m.DbAuth.Db) == 0 {
		m.DbAuth.Db = `0`
	}
	scheme := `redis`
	host := m.DbAuth.Host
	if strings.Contains(m.DbAuth.Host, `://`) {
		info := strings.SplitAfterN(m.DbAuth.Host, `://`, 1)
		scheme = info[0]
		host = info[1]
	}
	r.conn, err = redis.DialURL(scheme + `://` + m.DbAuth.Username + `:` + m.DbAuth.Password + `@` + host + `/` + m.DbAuth.Db)
	return
}

func (r *redis) Info() error {
	r, err := r.conn.Do("INFO")
	if err != nil {
		return err
	}
	return err
}
