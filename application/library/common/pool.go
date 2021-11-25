package common

import (
	"net/url"
	"sync"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

var (
	hpool = sync.Pool{
		New: func() interface{} {
			return echo.H{}
		},
	}

	urlValuesPool = sync.Pool{
		New: func() interface{} {
			return url.Values{}
		},
	}

	stringMapPool = sync.Pool{
		New: func() interface{} {
			return param.StringMap{}
		},
	}
)

func HPoolGet() echo.H {
	return hpool.Get().(echo.H)
}

func HPoolRelease(m echo.H) {
	for k := range m {
		delete(m, k)
	}

	hpool.Put(m)
}

func URLValuesPoolGet() url.Values {
	return urlValuesPool.Get().(url.Values)
}

func URLValuesPoolRelease(m url.Values) {
	for k := range m {
		delete(m, k)
	}

	urlValuesPool.Put(m)
}

func StringMapPoolGet() param.StringMap {
	return stringMapPool.Get().(param.StringMap)
}

func StringMapPoolRelease(m param.StringMap) {
	for k := range m {
		delete(m, k)
	}

	stringMapPool.Put(m)
}
