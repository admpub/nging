package restclient

import (
	"github.com/webx-top/restyclient"
)

var (
	ResetResty           = restyclient.ResetResty
	Resty                = restyclient.Resty
	ResetRestyRetryable  = restyclient.ResetRestyRetryable
	RestyRetryable       = restyclient.RestyRetryable
	SetProxy             = restyclient.SetProxy
	NewClient            = restyclient.NewClient
	NewCookiejar         = restyclient.NewCookiejar
	ProxyURL             = restyclient.ProxyURL
	NewClientWithOptions = restyclient.NewClientWithOptions
	InitRestyHook        = restyclient.InitRestyHook
	OutputMaps           = restyclient.OutputMaps
)
