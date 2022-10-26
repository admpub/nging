package webhook

import (
	"errors"
	"fmt"
	"strings"

	"github.com/admpub/nging/v5/application/library/restclient"
)

var Methods = [6]string{`GET`, `POST`, `PUT`, `HEAD`, `PATCH`, `DELETE`}
var (
	ErrInvalidHTTPMethod  = errors.New(`invalid http method`)
	ErrWebhookUrlRequired = errors.New(`webhook URL cannot be empty`)
	ErrBadHTTPStatus      = errors.New("bad HTTP status")
)

func New() *Webhook {
	return &Webhook{}
}

type Webhook struct {
	Name    string
	Method  string
	Url     string
	Content string
	Header  string
	headers []string
}

// AddHeader(`headerName=headerValue`)
func (w *Webhook) AddHeader(kvs ...string) *Webhook {
	w.headers = append(w.headers, kvs...)
	return w
}

func (w *Webhook) Validate() error {
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

func (w *Webhook) Exec(contentTagReplacer func(string) string, urlTagReplacer func(string) string) error {
	client := restclient.RestyRetryable()
	if len(w.Content) > 0 {
		content := w.Content
		if contentTagReplacer != nil {
			content = contentTagReplacer(content)
		}
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
	if len(w.headers) > 0 {
		for _, h := range w.headers {
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
	url := w.Url
	if urlTagReplacer != nil {
		url = urlTagReplacer(url)
	}
	method := w.Method
	if len(method) == 0 {
		if len(w.Content) > 0 {
			method = `POST`
		} else {
			method = `GET`
		}
	}
	resp, err := client.Execute(method, url)
	if err != nil {
		return fmt.Errorf(`%s: %w`, name, err)
	}
	if resp.IsError() {
		return fmt.Errorf(`%s: %w: %d: %s`, name, ErrBadHTTPStatus, resp.StatusCode(), resp.String())
	}
	return err
}
