package pkg

import (
	"sync"
	"time"

	"github.com/admpub/resty/v2"
	"github.com/webx-top/restyclient"
)

var (
	client               *resty.Client
	once                 sync.Once
	DefaultMaxRetryCount = 5
)

func initHTTPClient() {
	client = resty.New().
		SetRetryCount(DefaultMaxRetryCount).
		SetTimeout(time.Hour * 24).
		SetRedirectPolicy(restyclient.DefaultRedirectPolicy)

	restyclient.InitRestyHook(client)
}

func Request() *resty.Request {
	once.Do(initHTTPClient)
	return client.R()
}
