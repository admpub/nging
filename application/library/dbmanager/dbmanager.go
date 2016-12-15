package dbmanager

import (
	//"github.com/webx-top/com"
	"errors"

	"github.com/webx-top/echo"
)

func New(ctx echo.Context) *dbManager {
	return &dbManager{
		Context: ctx,
	}
}

type dbManager struct {
	echo.Context
}

func (d *dbManager) Driver(typeName string) (Driver, error) {
	driver, ok := Get(typeName)
	if ok {
		driver.Init(d.Context)
		return driver, nil
	}
	return nil, errors.New(d.T(`很抱歉，暂时不支持%v`, typeName))
}

func (d *dbManager) Run(typeName string, operation string) error {
	driver, err := d.Driver(typeName)
	if err != nil {
		return err
	}
	if !driver.IsSupported(operation) {
		return errors.New(d.T(`很抱歉，不支持此项操作`))
	}
	switch operation {
	case `login`:
		return driver.Login()
	case `logout`:
		return driver.Logout()
	case `processList`:
		return driver.ProcessList()
	case `privileges`:
		return driver.Privileges()
	case `info`:
		return driver.Info()
	case `createDb`:
		return driver.CreateDb()
	case `modifyDb`:
		return driver.ModifyDb()
	case `listDb`:
		return driver.ListDb()
	case `createTable`:
		return driver.CreateTable()
	case `modifyTable`:
		return driver.ModifyTable()
	case `listTable`:
		return driver.ListTable()
	case `viewTable`:
		return driver.ViewTable()
	case `listData`:
		return driver.ListData()
	case `createData`:
		return driver.CreateData()
	case `indexes`:
		return driver.Indexes()
	case `foreign`:
		return driver.Foreign()
	case `trigger`:
		return driver.Trigger()
	case `runCommand`:
		return driver.RunCommand()
	case `import`:
		return driver.Import()
	case `export`:
		return driver.Export()
	default:
		return errors.New(d.T(`很抱歉，不支持此项操作`))
	}
}
