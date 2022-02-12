package restclient

import (
	"github.com/webx-top/restyclient"
)

var (
	ResetResty           = restyclient.ResetClassic
	Resty                = restyclient.Classic
	ResetRestyRetryable  = restyclient.ResetRetryable
	RestyRetryable       = restyclient.Retryable
	SetProxy             = restyclient.SetProxy
	NewClient            = restyclient.New
	NewCookiejar         = restyclient.NewCookiejar
	ProxyURL             = restyclient.ProxyURL
	NewClientWithOptions = restyclient.NewWithOptions
	InitRestyHook        = restyclient.InitRestyHook
	OutputMaps           = restyclient.OutputMaps
)
