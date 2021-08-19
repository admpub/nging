package restclient

import (
	"sync"
	"time"

	"github.com/admpub/resty/v2"
)

// DefaultTimeout 默认超时时间
var (
	DefaultTimeout     = 10 * time.Second
	restyClient        *resty.Client
	restyRetryable     *resty.Client
	restyOnce          sync.Once
	restyRetryableOnce sync.Once
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

// - retryable -

func initRetryable() {
	restyRetryable = resty.New().SetRetryCount(3).SetTimeout(DefaultTimeout)
}

func ResetRestyRetryable() {
	restyRetryableOnce = sync.Once{}
}

func RestyRetryable() *resty.Request {
	restyRetryableOnce.Do(initRetryable)
	return restyRetryable.R()
}
