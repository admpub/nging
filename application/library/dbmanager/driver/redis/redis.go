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
	"strconv"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v3/application/library/common"
	"github.com/admpub/nging/v3/application/library/dbmanager/driver"
)

func init() {
	driver.Register(`redis`, &Redis{})
}

type Redis struct {
	*driver.BaseDriver
	conn   redis.Conn
	dbName string
	//codecs map[string]map[string]func(string)string
}

func (r *Redis) Name() string {
	return `Redis`
}

func (r *Redis) Init(ctx echo.Context, auth *driver.DbAuth) {
	r.BaseDriver = driver.NewBaseDriver()
	r.BaseDriver.Init(ctx, auth)
}

func (r *Redis) IsSupported(operation string) bool {
	return true
}

func (r *Redis) getTables() (string, []string, error) {
	key := r.Form(`searchkey`, `*`)
	offset := r.Formx(`offset`).Int64()
	size := r.Formx(`size`, `20`).Int()
	return r.ListKeys(size, offset, key)
}

func (r *Redis) Login() (err error) {
	r.dbName = r.Form(`db`, `0`)
	err = r.login()
	if err != nil {
		return err
	}
	r.Set(`dbName`, r.dbName)
	r.Set(`table`, r.Form(`table`))
	return r.baseInfo()
}

func (r *Redis) login() (err error) {
	scheme := `redis`
	host := r.DbAuth.Host
	switch {
	case len(host) == 0:
		host = `127.0.0.1:6379`
	case strings.Contains(host, `://`):
		info := strings.SplitAfterN(host, `://`, 2)
		scheme = info[0]
		host = info[1]
		if !strings.Contains(host, `:`) {
			host += `:6379`
		}
	case !strings.Contains(host, `:`):
		host += `:6379`
	}
	options := []redis.DialOption{}
	if len(r.DbAuth.Username) > 0 {
		options = append(options, redis.DialUsername(r.DbAuth.Username))
	}
	if len(r.DbAuth.Password) > 0 {
		options = append(options, redis.DialPassword(r.DbAuth.Password))
	}
	if len(r.dbName) > 0 && com.IsUnsignedInteger(r.dbName) {
		db, _ := strconv.Atoi(r.dbName)
		if db != 0 {
			options = append(options, redis.DialDatabase(db))
		}
	}
	options = append(options, redis.DialUseTLS(scheme == "rediss"))
	r.conn, err = redis.Dial("tcp", host, options...)
	return
}

func (r *Redis) ListDb() error {
	return r.Info()
}

func (r *Redis) ModifyDb() error {
	delDB := r.Form(`deldb`)
	data := r.Data()
	result, err := r.Flush(delDB)
	if err != nil {
		data.SetError(err)
	} else {
		var msg string
		if delDB == `all` {
			msg = r.T(`已经成功清空全部数据库`)
		} else {
			msg = r.T(`已经成功清空数据库“%v”`, r.dbName)
		}
		data.SetInfo(msg).SetData(result)
	}
	return r.JSON(data)
}

func (r *Redis) ModifyTable() error {
	operate := r.Form(`operate`)
	var err error
	switch operate {
	case `deleteValue`:
		data := r.Data()
		key := r.Form(`key`)
		var dataType string
		dataType, err = r.DataType(key)
		if err != nil {
			return r.JSON(data.SetError(err))
		}
		hkey := r.Form(`hkey`)
		index := r.Form(`index`)
		value := r.Form(`value`)
		err = r.deleteValue(dataType, key, hkey, index, value)
		if err != nil {
			data.SetError(err)
		} else {
			msg := r.T(`已经成功删除`)
			data.SetInfo(msg, 1)
		}
		return r.JSON(data)
	case `delete`:
		data := r.Data()
		key := r.Form(`key`)
		if len(key) == 0 {
			keys := r.FormValues(`key[]`)
			dels := []string{}
			for _, key := range keys {
				_, err = r.DeleteKey(key)
				if err != nil {
					var msg string
					if len(dels) > 0 {
						msg = r.T(`已经成功删除Key“%v”`, strings.Join(dels, `, `))
						msg += `, ` + r.T(`但删除Key“%v”失败: %v`, key, err)
					} else {
						msg = r.T(`删除Key“%v”失败: %v`, key, err)
					}
					return r.JSON(data.SetInfo(msg), 0)
				}
				dels = append(dels, key)
			}
			if len(dels) == len(keys) {
				data.SetInfo(r.T(`已经成功删除`), 1)
			} else {
				data.SetInfo(r.T(`已经成功删除Key“%v”`, strings.Join(dels, `, `)), 1)
			}
			return r.JSON(data)
		}
		_, err = r.DeleteKey(key)
		if err != nil {
			data.SetError(err)
		} else {
			data.SetInfo(r.T(`已经成功删除Key“%v”`, key), 1)
		}
		return r.JSON(data)
	case `editValue`:
		if r.IsPost() {
			data := r.Data()
			err = r.postEditValue()
			if err != nil {
				data.SetError(err)
			} else {
				data.SetInfo(r.T(`修改成功`), 1)
			}
			return r.JSON(data)
		}
		var (
			encoding string
			dataType string
			value    string
		)
		key := r.Form(`key`)
		hkey := r.Form(`hkey`)
		encoding, err = r.ObjectEncoding(key)
		if err != nil {
			goto END
		}
		dataType, err = r.DataType(key)
		if err != nil {
			goto END
		}
		value, err = r.ViewElement(key, hkey, dataType, encoding)
		if err != nil {
			goto END
		}

	END:
		r.Request().Form().Set(`type`, dataType)
		r.Request().Form().Set(`value`, value)
		return r.Render(`db/redis/edit_value`, common.Err(r.Context, err))
	default:
		if r.IsPost() {
			old := r.Form(`old`)
			key := r.Form(`key`)
			ttl := r.Formx(`ttl`).Int64()
			oldTTL := r.Formx(`oldTTL`).Int64()
			data := r.Data()
			if old != key {
				_, err = r.RenameKey(old, key)
				if err != nil {
					return r.JSON(data.SetError(err))
				}
			}
			if oldTTL != ttl {
				if ttl < 0 {
					_, err = redis.Int(r.conn.Do(`PERSIST`, key))
				} else {
					_, err = redis.Int(r.conn.Do(`EXPIRE`, key, ttl))
				}
			}
			if err != nil {
				data.SetError(err)
			} else {
				data.SetInfo(r.T(`修改成功`), 1)
			}
			return r.JSON(data)
		}
	}
	key := r.Form(`table`)
	var ttl int64
	ttl, err = r.TTL(key)
	ret := common.Err(r.Context, err)
	r.Request().Form().Set(`key`, key)
	r.Request().Form().Set(`ttl`, fmt.Sprint(ttl))
	return r.Render(`db/redis/modify_table`, ret)
}

func (r *Redis) deleteValue(dataType, key, hkey, index, value string) (err error) {
	switch dataType {
	case `string`:
		_, err = r.DeleteKey(key)
	case `hash`:
		_, err = r.conn.Do("HDEL", key, hkey)
	case `list`:
		_, err = redis.String(r.conn.Do("LSET", key, index, value))
		_, err = redis.Int(r.conn.Do("LREM", key, 1, value))
	case `set`:
		_, err = redis.Int(r.conn.Do("SREM", key, value))
	case `zset`:
		_, err = redis.Int(r.conn.Do("ZREM", key, value))
	}
	return
}

func (r *Redis) postEditValue() (err error) {
	dataType := r.Form(`type`)
	key := r.Form(`key`)
	hkey := r.Form(`hkey`)
	index := r.Formx(`index`)
	score := r.Formx(`score`)
	value := r.Form(`value`)
	oldValue := r.Form(`oldvalue`)
	switch dataType {
	case `string`:
		err = r.SetString(key, value)
	case `hash`:
		var has int
		has, err = redis.Int(r.conn.Do("HEXISTS", key, hkey))
		if err != nil {
			return
		}
		if has > 0 {
			_, err = r.conn.Do("HDEL", key, hkey)
			if err != nil {
				return
			}
		}
		_, err = redis.Int(r.conn.Do("HSET", key, hkey, value))
	case `list`:
		var size int
		size, err = redis.Int(r.conn.Do("LLEN", key))
		if err != nil {
			return
		}
		indexN := index.Int()
		if len(index.String()) == 0 || indexN == size || indexN == -1 {
			_, err = redis.Int(r.conn.Do("RPUSH", key, value))
		} else if indexN >= 0 && indexN < size {
			_, err = redis.String(r.conn.Do("LSET", key, indexN, value))
		}
	case `set`:
		if value != oldValue {
			_, err = redis.Int(r.conn.Do("SREM", key, oldValue))
			if err != nil {
				return
			}
			_, err = redis.Int(r.conn.Do("SADD", key, value))
		}
	case `zset`:
		_, err = redis.Int(r.conn.Do("ZREM", key, oldValue))
		if err != nil {
			return
		}
		_, err = redis.Int(r.conn.Do("ZADD", key, score, value))
	}
	return
}

func (r *Redis) CreateTable() error {
	var err error
	if r.IsPost() {
		data := r.Data()
		err = r.postEditValue()
		if err != nil {
			data.SetError(err)
		} else {
			data.SetInfo(r.T(`添加成功`), 1)
		}
		return r.JSON(data)
	}
	ret := common.Err(r.Context, err)
	return r.Render(`db/redis/edit_value`, ret)
}

// Info 获取redis服务信息
func (r *Redis) Info() error {
	infos, err := r.info()
	if err != nil {
		return err
	}
	r.Set(`infos`, infos)
	section := r.Form(`section`)
	if len(section) == 0 {
		section = `Server`
		r.Request().Form().Set(`section`, section)
	}
	ret := common.Err(r.Context, err)
	return r.Render(`db/redis/info`, ret)
}

func (r *Redis) ListTable() error {
	return r.Render(`db/redis/list_table`, nil)
}

func (r *Redis) ViewTable() error {
	tmpl := `db/redis/view_table`
	key := r.Form(`key`, r.Form(`table`))
	var (
		encoding string
		dataType string
		result   *Value
		ttl      int64
		err      error
		page     = r.Formx(`vpage`, `0`).Int64()
		offset   = r.Formx(`voffset`, `0`).Int64()
		size     = r.Formx(`vsize`, `50`).Int()
	)
	encoding, err = r.ObjectEncoding(key)
	if err != nil {
		goto END
	}
	dataType, err = r.DataType(key)
	if err != nil {
		goto END
	}
	ttl, err = r.TTL(key)
	if err != nil {
		goto END
	}
	if page > 0 {
		offset = int64(com.Offset(uint(page), uint(size)))
	}
	result, err = r.ViewValuePro(key, dataType, encoding, size, offset)
	if err != nil {
		goto END
	}

END:
	r.Set(`encoding`, encoding)
	r.Set(`dataType`, dataType)
	r.Set(`result`, result)
	r.Set(`ttl`, ttl)
	return r.Render(tmpl, common.Err(r.Context, err))
}

func (r *Redis) Close() error {
	if r.conn == nil {
		return nil
	}
	return r.conn.Close()
}

func (r *Redis) Logout() error {
	return r.Close()
}
