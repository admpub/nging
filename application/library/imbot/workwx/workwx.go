// Package workwechat 企业微信机器人
package workwechat

import (
	"github.com/admpub/nging/application/library/imbot/http"
	"github.com/admpub/nging/application/library/imbot"
)

// document https://work.weixin.qq.com/help?person_id=1&doc_id=13376

var (
	KEY = ``
	API_URL = `https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%v`
)

func init() {
	imbot.Register(`workwx`, `企业微信机器人`, &Message{})
}

// Message 企业微信消息实体 每个机器人发送的消息不能超过20条/分钟。
type Message struct {
	MsgType string `json:"msgtype"` //消息类型 text/markdown/image/news
	Text *Text `json:"text,omitempty"` // MsgType=text 时有效
	Markdown *Markdown `json:"markdown,omitempty"` // MsgType=markdown 时有效
	Image *Image `json:"image,omitempty"` // MsgType=image 时有效
	News *News `json:"news,omitempty"` // MsgType=news 时有效
}

func (m *Message) BuildURL(args ...string) string {
	key := KEY
	if len(args) > 0 {
		key = args[0]
	}
	return fmt.Sprintf(API_URL, key)
}

func (m *Message) Reset() imbot.Messager {
	m.MsgType = ``
	m.Text = nil
	m.Markdown = nil
	m.Image = nil
	m.News = nil
	return m
}

func (m *Message) SendText(url, text string, atMobiles ...string) error {
	m.MsgType = `text`
	m.Text = &Text{Content:text}
	if len(atMobiles) > 0 {
		m.Text.MentionedMobileList = atMobiles
	}
	_, err := http.Send(url, m)
	return err
}

func (m *Message) SendMarkdown(url, title, markdown string, atMobiles ...string) error {
	m.MsgType = `markdown`
	m.Markdown = &Markdown{Content:markdown}
	if len(atMobiles) > 0 {
		m.Markdown.MentionedMobileList = atMobiles
	}
	_, err := http.Send(url, m)
	return err
}

type Text struct {
	Content string `json:"content"` //文本内容，最长不超过2048个字节，必须是utf8编码
	MentionedList []string `json:"mentioned_list"` //userid的列表，提醒群中的指定成员(@某个成员)，@all表示提醒所有人，如果开发者获取不到userid，可以使用mentioned_mobile_list
	MentionedMobileList []string `json:"mentioned_mobile_list"` //手机号列表，提醒手机号对应的群成员(@某个成员)，@all表示提醒所有人
}

type Markdown struct {
	Content string `json:"content"` //markdown内容，最长不超过4096个字节，必须是utf8编码
	MentionedMobileList []string `json:"mentioned_mobile_list"` //手机号列表，提醒手机号对应的群成员(@某个成员)，@all表示提醒所有人
}

type Image struct {
	Base64 string `json:"base64"` //图片内容的base64编码
	Md5 string `json:"md5"` //图片内容（base64编码前）的md5值
}

type News struct{
	Articles []*Article `json:"articles"` //图文消息，一个图文消息支持1到8条图文
}

type Article struct{
	Title string `json:"title"` //标题，不超过128个字节，超过会自动截断
	Description string `json:"description"` //描述，不超过512个字节，超过会自动截断
	URL string `json:"url"` //点击后跳转的链接。
	PicURL string `json:"picurl"` //图文消息的图片链接，支持JPG、PNG格式，较好的效果为大图 1068*455，小图150*150。
}
