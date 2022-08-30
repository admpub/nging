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
	"errors"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/gomodule/redigo/redis"
)

func (r *Redis) baseInfo() error {
	if r.Get(`dbList`) == nil {
		dbList, err := r.DatabaseList()
		if err != nil {
			r.SetFail(err.Error())
			return r.Goto(`/db`)
		}
		r.Set(`dbList`, dbList)
	}
	if len(r.dbName) > 0 {
		nextOffset, tableList, err := r.getTables()
		if err != nil {
			r.SetFail(err.Error())
			return r.Goto(r.GenURL(`listDb`))
		}
		r.Set(`tableList`, tableList)
		r.Set(`nextOffset`, nextOffset)
		_, _, _, pagination := handler.PagingWithPagination(r.Context)
		prevOffset := r.Form(`offset`, `0`)
		q := r.Request().URL().Query()
		q.Del(`offset`)
		q.Del(`prev`)
		pagination.SetURL(`/db?`+q.Encode()+`&offset={curr}&prev={prev}`).SetPosition(prevOffset, nextOffset, nextOffset)
		r.Set(`pagination`, pagination)
	}

	r.Set(`dbVersion`, r.getVersion())
	return nil
}

func (r *Redis) getVersion() string {
	info, err := redis.String(r.conn.Do("INFO", "server"))
	if err != nil {
		return err.Error()
	}
	infos := ParseInfos(info)
	if len(infos) > 0 {
		for _, attr := range infos[0].Attrs {
			if attr.Name == `redis_version` {
				return attr.Value
			}
		}
	}
	return info
}

func (r *Redis) info(args ...interface{}) ([]*Infos, error) {
	if r.conn == nil {
		return nil, errors.New(`Redis connection failed`)
	}
	info, err := redis.String(r.conn.Do(`INFO`, args...))
	if err != nil {
		return nil, err
	}
	return ParseInfos(info), err
}
