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
package dbmanager

import (
	"errors"

	"github.com/admpub/nging/application/library/dbmanager/driver"
	"github.com/webx-top/echo"
)

func New(ctx echo.Context, auth *driver.DbAuth) *dbManager {
	return &dbManager{
		Context: ctx,
		DbAuth:  auth,
	}
}

type dbManager struct {
	echo.Context
	*driver.DbAuth
	GenURL func(string, ...string) string
}

func (d *dbManager) Driver(typeName string) (driver.Driver, error) {
	dv, ok := driver.Get(typeName)
	if ok {
		dv.Init(d.Context, d.DbAuth)
		return dv, nil
	}
	return nil, errors.New(d.T(`很抱歉，暂时不支持%v`, typeName))
}

func (d *dbManager) Run(typeName string, operation string) error {
	drv, err := d.Driver(typeName)
	if err != nil {
		return err
	}
	if !drv.IsSupported(operation) {
		return errors.New(d.T(`很抱歉，不支持此项操作`))
	}
	defer drv.SaveResults()
	drv.SetURLGenerator(d.GenURL)
	d.Set(`dbType`, drv.Name())
	d.SetFunc(`Results`, drv.SavedResults)
	switch operation {
	case `login`:
		return drv.Login()
	case `logout`:
		return drv.Logout()
	case `processList`:
		return drv.ProcessList()
	case `privileges`:
		return drv.Privileges()
	case `info`:
		return drv.Info()
	case `createDb`:
		return drv.CreateDb()
	case `modifyDb`:
		return drv.ModifyDb()
	case `listDb`:
		return drv.ListDb()
	case `createTable`:
		return drv.CreateTable()
	case `modifyTable`:
		return drv.ModifyTable()
	case `listTable`:
		return drv.ListTable()
	case `viewTable`:
		return drv.ViewTable()
	case `listData`:
		return drv.ListData()
	case `createData`:
		return drv.CreateData()
	case `indexes`:
		return drv.Indexes()
	case `foreign`:
		return drv.Foreign()
	case `trigger`:
		return drv.Trigger()
	case `runCommand`:
		return drv.RunCommand()
	case `import`:
		return drv.Import()
	case `export`:
		return drv.Export()
	case `analysis`:
		return drv.Analysis()
	default:
		return errors.New(d.T(`很抱歉，不支持此项操作`))
	}
}
