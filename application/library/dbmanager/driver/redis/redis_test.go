package redis

import (
	"github.com/admpub/nging/application/library/dbmanager/driver"
	"github.com/webx-top/echo"
	"testing"
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
	err := r.Login()
	if err != nil {
		panic(err)
	}
	return r
}

func TestInfo(t *testing.T) {
	r := connect()
	info, err := r.Info()
	if err != nil {
		panic(err)
	}
	echo.Dump(info)
}

func TestFindKeys(t *testing.T) {
	r := connect()
	info, err := r.FindKeys(`*`)
	if err != nil {
		panic(err)
	}
	echo.Dump(info)
}

func TestDatabaseList(t *testing.T) {
	r := connect()
	info, err := r.DatabaseList()
	if err != nil {
		panic(err)
	}
	echo.Dump(info)
}
