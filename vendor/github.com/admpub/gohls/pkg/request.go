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
	UserAgent            string // "Mozilla/5.0 (X11; Linux x86_64; rv:38.0) Gecko/38.0 Firefox/38.0"
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
	if len(UserAgent) > 0 {
		return client.R().SetHeader("User-Agent", UserAgent)
	}
	return client.R()
}
