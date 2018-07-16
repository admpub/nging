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
	defer r.Close()
	info, err := r.Info()
	if err != nil {
		panic(err)
	}
	echo.Dump(info)
}

func TestFindKeys(t *testing.T) {
	r := connect()
	defer r.Close()
	info, err := r.FindKeys(`*`)
	if err != nil {
		panic(err)
	}
	echo.Dump(info)
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
