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

	"github.com/gomodule/redigo/redis"
)

// DatabaseList 获取数据库列表
func (r *Redis) DatabaseList() ([]string, error) {
	reply, err := redis.Strings(r.conn.Do("CONFIG", "GET", "databases"))
	if err != nil {
		return nil, err
	}
	var ids []string
	if len(reply) > 1 {
		num, err := strconv.ParseInt(reply[1], 10, 64)
		if err != nil {
			return nil, err
		}
		var id int64
		for ; id < num; id++ {
			ids = append(ids, fmt.Sprint(id))
		}
	}
	return ids, err
}

func (r *Redis) Flush(db string) (string, error) {
	var (
		rep interface{}
		err error
	)
	if db != `all` {
		rep, err = r.conn.Do("FLUSHDB")
	} else {
		rep, err = r.conn.Do("FLUSHALL")
	}
	return redis.String(rep, err)
}
