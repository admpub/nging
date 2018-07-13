package redis

import (
	"fmt"
	"github.com/admpub/nging/application/library/dbmanager/driver"
	"github.com/webx-top/echo"
	"testing"
)

func TestInfo(t *testing.T) {
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
	info, err := r.Info()
	if err != nil {
		panic(err)
	}
	echo.Dump(info)
}
