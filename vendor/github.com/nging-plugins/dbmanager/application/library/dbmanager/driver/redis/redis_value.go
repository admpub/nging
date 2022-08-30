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

package redis

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

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

func (r *Redis) ViewValue(key string, typ string, encoding string) (ret string, siz int, err error) {
	switch typ {
	case `string`:
		ret, err = redis.String(r.conn.Do("GET", key))
		if err != nil {
			return
		}
		ret = r.Codec(`load`, key, ret, encoding)
		siz = len(ret)
	case `hash`:
		var arr map[string]string
		arr, err = redis.StringMap(r.conn.Do("HGETALL", key))
		if err != nil {
			return
		}
		for k, v := range arr {
			arr[k] = r.Codec(`load`, key, v, encoding)
		}
		siz = len(arr)
	case `list`:
		siz, err = redis.Int(r.conn.Do("LLEN", key))
		if err != nil {
			return
		}
	case `set`:
		var arr map[string]string
		arr, err = redis.StringMap(r.conn.Do("SMEMBERS", key))
		if err != nil {
			return
		}
		for k, v := range arr {
			arr[k] = r.Codec(`load`, key, v, encoding)
		}
		siz = len(arr)
	case `zset`:
		var arr map[string]string
		arr, err = redis.StringMap(r.conn.Do("ZRANGE", key, 0, -1))
		if err != nil {
			return
		}
		for k, v := range arr {
			arr[k] = r.Codec(`load`, key, v, encoding)
		}
		siz = len(arr)
	}
	return
}

func (r *Redis) ViewValuePro(key string, typ string, encoding string, size int, offset int64, pattern ...string) (v *Value, err error) {
	v = NewValue(r.Context)
	switch typ {
	case `string`:
		v.Text, err = redis.String(r.conn.Do("GET", key))
		if err != nil {
			return
		}
		v.Text = r.Codec(`load`, key, v.Text, encoding)
		v.TotalRows = len(v.Text)
	case `hash`:
		v.TotalRows, err = redis.Int(r.conn.Do("HLEN", key))
		if err != nil {
			return
		}
		var reply interface{}
		if len(pattern) > 0 {
			reply, err = r.conn.Do("HSCAN", key, offset, "MATCH", pattern[0], "COUNT", size)
		} else {
			reply, err = r.conn.Do("HSCAN", key, offset, "COUNT", size)
		}
		rows := reply.([]interface{})
		list := rows[1].([]interface{})
		var key string
		for index, value := range list {
			if index%2 == 0 {
				key = string(value.([]byte))
				continue
			}
			v.Add(key, string(value.([]byte)))
		}
		v.NextOffset = string(rows[0].([]byte))
	case `list`:
		v.TotalRows, err = redis.Int(r.conn.Do("LLEN", key))
		if err != nil {
			return
		}
		end := offset + int64(size) - 1
		var values []string
		values, err = redis.Strings(r.conn.Do("LRANGE", key, offset, end))
		if err != nil {
			return
		}
		for index, value := range values {
			v.Add(fmt.Sprint(offset+int64(index)), value)
		}
		nextOffset := end + 1
		if nextOffset > int64(v.TotalRows) {
			v.NextOffset = `0`
		} else {
			v.NextOffset = fmt.Sprint(nextOffset)
		}
	case `set`:
		v.TotalRows, err = redis.Int(r.conn.Do("SCARD", key))
		if err != nil {
			return
		}
		var reply interface{}
		if len(pattern) > 0 {
			reply, err = r.conn.Do("SSCAN", key, offset, "MATCH", pattern[0], "COUNT", size)
		} else {
			reply, err = r.conn.Do("SSCAN", key, offset, "COUNT", size)
		}
		rows := reply.([]interface{})
		var values []string
		values, err = redis.Strings(rows[1], err)
		if err != nil {
			return
		}
		for index, value := range values {
			v.Add(fmt.Sprint(index), value)
		}
		v.NextOffset = string(rows[0].([]byte))
	case `zset`:
		v.TotalRows, err = redis.Int(r.conn.Do("ZCARD", key))
		if err != nil {
			return
		}
		var reply interface{}
		if len(pattern) > 0 {
			reply, err = r.conn.Do("ZSCAN", key, offset, "MATCH", pattern[0], "COUNT", size)
		} else {
			reply, err = r.conn.Do("ZSCAN", key, offset, "COUNT", size)
		}
		rows := reply.([]interface{})
		list := rows[1].([]interface{})
		var val string
		for index, value := range list {
			if index%2 == 0 {
				val = string(value.([]byte))
				continue
			}
			key := string(value.([]byte))
			v.Add(key, val)
		}
		v.NextOffset = string(rows[0].([]byte))
	}
	return
}

func (r *Redis) ViewElement(key string, hkey string, typ string, encoding string) (v string, err error) {
	switch typ {
	case `string`:
		v, err = redis.String(r.conn.Do("GET", key))
	case `hash`:
		v, err = redis.String(r.conn.Do("HGET", key, hkey))
	case `list`:
		v, err = redis.String(r.conn.Do("LINDEX", key, r.Form(`index`)))
	case `set`:
		v = r.Form(`value`)
	case `zset`:
		v = r.Form(`value`)
	}
	return
}

func (r *Redis) Codec(action string, key string, data string, encoding string) string {
	if encoding == `raw` {
		return data
	}
	return data
}

func (r *Redis) SetString(key string, value string) error {
	_, err := r.conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	return err
}
