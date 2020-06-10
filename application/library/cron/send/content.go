package send

import (
	"bytes"

	"github.com/webx-top/echo/param"
)

type ContentType interface{
	EmailContent(params param.Store) []byte
	MarkdownContent(params param.Store) []byte
}

func NewContent() *Content {
	return &Content{}
}

type Content struct {
	emailContent []byte
	markdownContent []byte
}

func (c *Content) EmailContent(params param.Store) []byte {
	if c.emailContent == nil {
		b := new(bytes.Buffer)
		MailTpl().Execute(b, params)
		c.emailContent = b.Bytes()
	}
	return c.emailContent
}

func (c *Content) MarkdownContent(params param.Store) []byte {
	if c.markdownContent == nil {
		b := new(bytes.Buffer)
		MarkdownTmpl().Execute(b, params)
		c.markdownContent = b.Bytes()
	}
	return c.markdownContent
}
