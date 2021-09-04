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

func NewData(title string, ct ContentType) *AlertData {
	return &AlertData{
		Title:   title,
		Content: ct,
		Data:    param.Store{},
	}
}

type AlertData struct {
	Title   string
	Content ContentType
	Data    param.Store
}

var SendTopic = func(_ echo.Context, topic string, _ *AlertData) error {
	return nil
}
