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
	"testing"
	"time"

	"github.com/admpub/nging/v3/application/library/dbmanager/driver"
	"github.com/webx-top/echo"
)

func connect() *Redis {
	e := echo.New()
	c := echo.NewContext(nil, nil, e)
	r := &Redis{}
	r.Init(c, &driver.DbAuth{
		Driver:   `redis`,
		Username: ``,
		Password: ``,
		Host:     `127.0.0.1`,
		Db:       `0`,
	})
	err := r.login()
	if err != nil {
		panic(err)
	}
	return r
}

func TestInfo(t *testing.T) {
	r := connect()
	defer r.Close()
	info, err := r.info()
	if err != nil {
		panic(err)
	}
	echo.Dump(info)
}

func TestSetString(t *testing.T) {
	r := connect()
	defer r.Close()
	for i := 0; i < 5000; i++ {
		err := r.SetString(fmt.Sprintf(`test_%d`, i), time.Now().Format(`2006-01-02 15:04:05`))
		if err != nil {
			panic(err)
		}
	}
}

/*
func TestFindKeys(t *testing.T) {
	r := connect()
	defer r.Close()
	info, err := r.FindKeys(`*`)
	if err != nil {
		panic(err)
	}
	echo.Dump(info)
}
*/
func TestListKeys(t *testing.T) {
	r := connect()
	defer r.Close()
	offset, keys, err := r.ListKeys(20, 0, `*`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("--------------> cursor: %v\n", offset)
	echo.Dump(keys)
}

func TestDatabaseList(t *testing.T) {
	r := connect()
	defer r.Close()
	info, err := r.DatabaseList()
	if err != nil {
		panic(err)
	}
	echo.Dump(info)
}

func TestDataType(t *testing.T) {
	r := connect()
	defer r.Close()
	keys, err := r.FindKeys(`test`)
	if err != nil {
		panic(err)
	}
	if len(keys) < 1 {
		err = r.SetString(`test`, `2343333`)
		if err != nil {
			panic(err)
		}
		keys = []string{`test`}
	}
	encoding, err := r.ObjectEncoding(keys[0])
	if err != nil {
		panic(err)
	}
	dataType, err := r.DataType(keys[0])
	if err != nil {
		panic(err)
	}
	result, size, err := r.ViewValue(keys[0], dataType, encoding)
	if err != nil {
		panic(err)
	}
	echo.Dump(map[string]interface{}{
		`encoding`: encoding,
		`dataType`: dataType,
		`result`:   result,
		`size`:     size,
	})
}
