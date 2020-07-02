package httpClientOptions

import (
	"net/http"

	"github.com/webx-top/com"
)

func NewClient(options ...com.HTTPClientOptions) *http.Client {
	c := &http.Client{}
	return Apply(c, options...)
}

func Apply(client *http.Client, options ...com.HTTPClientOptions) *http.Client {
	for _, option := range options {
		option(client)
	}
	return client
}
