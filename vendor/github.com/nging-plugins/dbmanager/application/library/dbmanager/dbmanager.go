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
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/nging-plugins/dbmanager/application/library/dbmanager/driver"
)

type Manager interface {
	Driver() (driver.Driver, error)
	Run(operation string) error
	Context() echo.Context
	Account() *driver.DbAuth
	GenURL() func(string, ...string) string
	SetURLGenerator(fn func(string, ...string) string)
	Logined() (bool, error)
}

func New(ctx echo.Context, auth *driver.DbAuth) Manager {
	return &dbManager{
		context: ctx,
		dbAuth:  auth,
	}
}

type dbManager struct {
	context echo.Context
	dbAuth  *driver.DbAuth
	genURL  func(string, ...string) string
	driver  driver.Driver
}

func (d *dbManager) Context() echo.Context {
	return d.context
}

func (d *dbManager) GenURL() func(string, ...string) string {
	return d.genURL
}

func (d *dbManager) SetURLGenerator(fn func(string, ...string) string) {
	d.genURL = fn
}

func (d *dbManager) Account() *driver.DbAuth {
	return d.dbAuth
}

func (d *dbManager) Driver() (driver.Driver, error) {
	if d.driver != nil {
		return d.driver, nil
	}
	info, ok := driver.Get(d.dbAuth.Driver)
	if ok {
		d.driver = info.New()
		d.driver.Init(d.context, d.dbAuth)
		return d.driver, nil
	}
	return nil, d.context.NewError(code.Unsupported, `很抱歉，暂时不支持%v`, d.dbAuth.Driver)
}

func (d *dbManager) Logined() (bool, error) {
	drv, err := d.Driver()
	if err != nil {
		return false, err
	}
	return drv.Logined(), nil
}

func (d *dbManager) Run(operation string) error {
	drv, err := d.Driver()
	if err != nil {
		return err
	}
	if !drv.IsSupported(operation) {
		return d.context.NewError(code.Unsupported, `很抱歉，不支持此项操作`)
	}
	defer drv.SaveResults()
	drv.SetURLGenerator(d.genURL)
	d.context.Set(`dbType`, drv.Name())
	d.context.SetFunc(`Results`, drv.SavedResults)
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
		return d.context.NewError(code.Unsupported, `很抱歉，不支持此项操作`)
	}
}
