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

import "github.com/gomodule/redigo/redis"

func (r *Redis) Exists(key string) (bool, error) {
	reply, err := redis.Int(r.conn.Do("EXISTS", key))
	if err != nil {
		return false, err
	}
	return reply == 1, err
}

// ListKeys 搜索key
func (r *Redis) ListKeys(size int, offset int64, pattern ...string) (string, []string, error) {
	var (
		reply interface{}
		err   error
	)
	if len(pattern) > 0 {
		reply, err = r.conn.Do("SCAN", offset, "MATCH", pattern[0], "COUNT", size)
	} else {
		reply, err = r.conn.Do("SCAN", offset, "COUNT", size)
	}
	if err != nil {
		return "", nil, err
	}
	rows := reply.([]interface{})
	list := rows[1].([]interface{})
	keys := make([]string, len(list))
	for index, key := range list {
		keys[index] = string(key.([]byte))
	}
	nextOffsetN := string(rows[0].([]byte))
	return nextOffsetN, keys, err
}

// FindKeys 搜索key
func (r *Redis) FindKeys(pattern string) ([]string, error) {
	reply, err := redis.Strings(r.conn.Do("KEYS", pattern))
	if err != nil {
		return nil, err
	}
	return reply, err
}

func (r *Redis) DeleteKey(key string) (int, error) {
	return redis.Int(r.conn.Do("DEL", key))
}

func (r *Redis) RenameKey(key string, newKey string) (string, error) {
	return redis.String(r.conn.Do("RENAME", key, newKey))
}

func (r *Redis) MoveKey(key string, destDB string) (string, error) {
	return redis.String(r.conn.Do("MOVE", key, destDB))
}
