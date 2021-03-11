package send

import (
	"bytes"

	"github.com/webx-top/echo/param"
)

// ContentType 消息内容类型
type ContentType interface {
	EmailContent(params param.Store) []byte
	MarkdownContent(params param.Store) []byte
}

// NewContent 创建消息内容结构
func NewContent() *Content {
	return &Content{}
}

// Content 消息内容结构
type Content struct {
	emailContent    []byte
	markdownContent []byte
}

// EmailContent 生成E-mail内容
func (c *Content) EmailContent(params param.Store) []byte {
	if c.emailContent == nil {
		b := new(bytes.Buffer)
		MailTpl().Execute(b, params)
		c.emailContent = b.Bytes()
	}
	return c.emailContent
}

// MarkdownContent 生成Markdown格式内容
func (c *Content) MarkdownContent(params param.Store) []byte {
	if c.markdownContent == nil {
		b := new(bytes.Buffer)
		MarkdownTmpl().Execute(b, params)
		c.markdownContent = b.Bytes()
	}
	return c.markdownContent
}
