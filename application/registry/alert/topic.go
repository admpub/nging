package alert

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

// ContentType 消息内容类型
type ContentType interface {
	EmailContent(params param.Store) []byte
	MarkdownContent(params param.Store) []byte
}

type TagValues interface {
	TagValues() map[string]string
}

func SendTopic(c echo.Context, topic string, a *AlertData) error {
	return a.SendTopic(c, topic)
}
