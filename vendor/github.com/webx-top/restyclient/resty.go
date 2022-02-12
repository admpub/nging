package restyclient

import (
	"time"

	syncOnce "github.com/admpub/once"
	"github.com/admpub/resty/v2"
)

var (
	// DefaultTimeout 默认超时时间
	DefaultTimeout        = 10 * time.Second
	DefaultMaxRetryCount  = 3
	DefaultRedirectPolicy = resty.FlexibleRedirectPolicy(5)
	classicClient         *resty.Client
	retryableClient       *resty.Client
	classicOnce           syncOnce.Once
	retryableOnce         syncOnce.Once
)

func initClassic() {
	classicClient = resty.New().
		SetTimeout(DefaultTimeout).
		SetRedirectPolicy(DefaultRedirectPolicy)

	InitRestyHook(classicClient)
}

func ResetClassic() {
	classicOnce.Reset()
}

func Classic() *resty.Request {
	classicOnce.Do(initClassic)
	return classicClient.R()
}

// - retryable -

func initRetryable() {
	retryableClient = resty.New().
		SetRetryCount(DefaultMaxRetryCount).
		SetTimeout(DefaultTimeout).
		SetRedirectPolicy(DefaultRedirectPolicy)

	InitRestyHook(retryableClient)
}

func ResetRetryable() {
	retryableOnce.Reset()
}

func Retryable() *resty.Request {
	retryableOnce.Do(initRetryable)
	return retryableClient.R()
}
