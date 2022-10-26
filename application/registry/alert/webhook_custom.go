package alert

import (
	"strings"

	"github.com/admpub/nging/v5/application/library/webhook"
)

func NewWebhookCustom() *WebhookCustom {
	return &WebhookCustom{
		Headers: []string{},
	}
}

type WebhookCustom struct {
	Name    string
	Method  string
	Url     string
	Content string
	Headers []string
}

func (w *WebhookCustom) ToWebhook() *webhook.Webhook {
	return (&webhook.Webhook{
		Name:    w.Name,
		Method:  w.Method,
		Url:     w.Url,
		Content: w.Content,
	}).AddHeader(w.Headers...)
}

func (w *WebhookCustom) Descriptions() []string {
	descs := []string{
		`<code>Name</code>: 名称（选填）。如不填，错误信息中可能会显示网址。如果网址中包含密码等敏感信息，为了避免泄漏，建议设置名称`,
		`<code>Method</code>: webhook 的提交方式（支持的有：` + strings.Join(webhook.Methods[:], `, `) + `）`,
		`<code>Url</code>: webhook 的网址`,
		`<code>Content</code>: 提交的 body 内容（采用 GET 和 HEAD 方式提交时，此项无效）`,
		`<code>Headers</code>: 提交的 header 键值对（格式为“<code>headerName</code>=<code>headerValue</code>”，例如：["Content-Type=application/x-www-form-urlencoded","Authorization=Bearer guest"]）`,
	}
	return descs
}
