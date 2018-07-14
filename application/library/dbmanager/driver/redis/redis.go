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
	"strconv"
	"strings"

	"github.com/admpub/nging/application/library/dbmanager/driver"
	"github.com/gomodule/redigo/redis"
	"github.com/webx-top/echo"
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

// Info 获取redis服务信息
func (r *Redis) Info() ([]*Infos, error) {
	info, err := redis.String(r.conn.Do("INFO"))
	if err != nil {
		return nil, err
	}
	return ParseInfos(info), err
}

// FindKeys 搜索key
func (r *Redis) FindKeys(pattern string) ([]string, error) {
	reply, err := redis.Strings(r.conn.Do("KEYS", pattern))
	if err != nil {
		return nil, err
	}
	return reply, err
}

// DatabaseList 获取数据库列表
func (r *Redis) DatabaseList() ([]int64, error) {
	reply, err := redis.Strings(r.conn.Do("CONFIG", "GET", "databases"))
	if err != nil {
		return nil, err
	}
	var ids []int64
	if len(reply) > 1 {
		num, err := strconv.ParseInt(reply[1], 10, 64)
		if err != nil {
			return nil, err
		}
		var id int64
		for ; id < num; id++ {
			ids = append(ids, id)
		}
	}
	return ids, err
}

// TTL 获取数据有效期
func (r *Redis) TTL(key string) (int64, error) {
	reply, err := redis.Int64(r.conn.Do("TTL", key))
	if err != nil {
		return -3, err
	}
	//reply(-2:key不存在;-1:永不过期;>=0:过期时间)
	return reply, err
}

// ObjectEncoding 获取对象编码方式
func (r *Redis) ObjectEncoding(key string) (string, error) {
	reply, err := redis.String(r.conn.Do("OBJECT", "encoding", key))
	if err != nil {
		return ``, err
	}
	return reply, err
}

func (r *Redis) DataType(key string) (string, error) {
	reply, err := redis.String(r.conn.Do("TYPE", key))
	if err != nil {
		return ``, err
	}
	return reply, err
}

func (r *Redis) Exists(key string) (bool, error) {
	reply, err := redis.Int(r.conn.Do("EXISTS", key))
	if err != nil {
		return false, err
	}
	return reply == 1, err
}

func (r *Redis) ViewValue(key string, typ string, encoding string) {
	switch typ {
	case `string`:
	case `hash`:
	case `list`:
	case `set`:
	case `zset`:
	}
}

func Codec(action string, key string, data string, encoding string) {

}
