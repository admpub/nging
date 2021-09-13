package restclient

import (
	"time"

	syncOnce "github.com/admpub/once"
	"github.com/admpub/resty/v2"
)

var (
	// DefaultTimeout 默认超时时间
	DefaultTimeout     = 10 * time.Second
	restyClient        *resty.Client
	restyRetryable     *resty.Client
	restyOnce          syncOnce.Once
	restyRetryableOnce syncOnce.Once
)

func initRestyClient() {
	restyClient = resty.New().SetTimeout(DefaultTimeout)
	InitRestyHook(restyClient)
}

func ResetResty() {
	restyOnce.Reset()
}

func Resty() *resty.Request {
	restyOnce.Do(initRestyClient)
	return restyClient.R()
}

// - retryable -

func initRetryable() {
	restyRetryable = resty.New().SetRetryCount(3).SetTimeout(DefaultTimeout)
	InitRestyHook(restyRetryable)
}

func ResetRestyRetryable() {
	restyRetryableOnce.Reset()
}

func RestyRetryable() *resty.Request {
	restyRetryableOnce.Do(initRetryable)
	return restyRetryable.R()
}
