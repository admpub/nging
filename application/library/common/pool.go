package common

import (
	"net/url"
	"sync"
	"time"

	"github.com/admpub/log"
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

func OnErrorRetry(f func() error, maxTimes int, interval time.Duration) error {
	err := f()
	if err == nil {
		return err
	}

	for i := 0; i < maxTimes; i++ {
		log.Errorf(`%v ([%d/%d] retry after %v)`, err, i+1, maxTimes, interval.String())
		time.Sleep(interval)
		err = f()
		if err == nil {
			return err
		}
	}
	return err
}
