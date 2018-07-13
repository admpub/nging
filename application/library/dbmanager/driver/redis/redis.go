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
	"github.com/webx-top/echo"
	"strings"
)

type Redis struct {
	*driver.BaseDriver
	conn redis.Conn
}

func (r *Redis) Name() string {
	return `Redis`
}

func (r *Redis) Init(ctx echo.Context, auth *driver.DbAuth) {
	r.BaseDriver = driver.NewBaseDriver()
	r.BaseDriver.Init(ctx, auth)
	r.Set(`supportSQL`, false)
}

func (r *Redis) IsSupported(operation string) bool {
	return true
}

func (r *Redis) Login() (err error) {
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
	if len(r.DbAuth.Db) == 0 {
		r.DbAuth.Db = `0`
	}
	scheme := `redis`
	host := r.DbAuth.Host
	if strings.Contains(r.DbAuth.Host, `://`) {
		info := strings.SplitAfterN(r.DbAuth.Host, `://`, 2)
		scheme = info[0]
		host = info[1]
	}
	r.conn, err = redis.DialURL(scheme + `://` + r.DbAuth.Username + `:` + r.DbAuth.Password + `@` + host + `/` + r.DbAuth.Db)
	return
}

func (r *Redis) Info() ([]*Infos, error) {
	info, err := redis.String(r.conn.Do("INFO"))
	if err != nil {
		return nil, err
	}
	return ParseInfos(info), err
}

func (r *Redis) FindKeys(pattern string) ([]string, error) {
	reply, err := redis.Strings(r.conn.Do("KEYS", pattern))
	if err != nil {
		return nil, err
	}
	return reply, err
}
