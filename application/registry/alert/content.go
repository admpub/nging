package alert

import (
	"github.com/webx-top/echo/param"
)

var DefaultTextContent = &TextContent{}

type TextContent struct {
}

func (c *TextContent) EmailContent(params param.Store) []byte {
	return params.Get(`email-content`).([]byte)
}

func (c *TextContent) MarkdownContent(params param.Store) []byte {
	return params.Get(`markdown-content`).([]byte)
}
