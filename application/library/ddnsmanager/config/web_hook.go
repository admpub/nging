package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/admpub/nging/v3/application/library/ddnsmanager/ddnserrors"
	"github.com/admpub/nging/v3/application/library/ddnsmanager/domain/dnsdomain"
	"github.com/admpub/nging/v3/application/library/restclient"
)

var Methods = [6]string{`GET`, `POST`, `PUT`, `HEAD`, `PATCH`, `DELETE`}
var ErrInvalidHTTPMethod = errors.New(`invalid http method`)
var ErrWebhookUrlRequired = errors.New(`webhook 网址不能为空`)

type Webhook struct {
	Name    string
	Method  string
	Url     string
	Content string
	Header  string
}

func (w *Webhook) validate() error {
	if len(w.Url) == 0 {
		return ErrWebhookUrlRequired
	}
	for _, method := range Methods {
		if method == w.Method {
			return nil
		}
	}
	return ErrInvalidHTTPMethod
}

func (w *Webhook) Exec(tagValues *dnsdomain.TagValues) error {
	client := restclient.RestyRetryable()
	if len(w.Content) > 0 {
		content := w.Content
		content = tagValues.Parse(content)
		client.SetBody(content)
	}
	header := strings.TrimSpace(w.Header)
	if len(header) > 0 {
		for _, h := range strings.Split(header, "\n") {
			h = strings.TrimSpace(h)
			if len(h) == 0 {
				continue
			}
			parts := strings.SplitN(h, `=`, 2)
			if len(parts) != 2 {
				continue
			}
			for k, v := range parts {
				parts[k] = strings.TrimSpace(v)
			}
			client.SetHeader(parts[0], parts[1])
		}
	}
	name := w.Name
	if len(name) == 0 {
		name = w.Url
	}
	resp, err := client.Execute(w.Method, w.Url)
	if err != nil {
		return fmt.Errorf(`%s: %w`, name, err)
	}
	if resp.IsError() {
		return fmt.Errorf(`%s: %w: %d: %s`, name, ddnserrors.ErrBadHTTPStatus, resp.StatusCode(), resp.String())
	}
	return err
}
