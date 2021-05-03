package restclient

import (
	"sync"
	"time"

	"github.com/admpub/resty/v2"
)

// DefaultTimeout 默认超时时间
var (
	DefaultTimeout = 10 * time.Second
	restyClient    *resty.Client
	restyOnce      sync.Once
)

func initRestyClient() {
	restyClient = resty.New().SetTimeout(DefaultTimeout)
}

func ResetResty() {
	restyOnce = sync.Once{}
}

func Resty() *resty.Request {
	restyOnce.Do(initRestyClient)
	return restyClient.R()
}
