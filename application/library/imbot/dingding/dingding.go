//Package dingding 钉钉机器人
package dingding

import (
	"github.com/admpub/nging/application/library/imbot/http"
	"github.com/admpub/nging/application/library/imbot"
)

// document https://ding-doc.dingtalk.com/doc#/serverapi3/iydd5h

// Message 消息实体 每个机器人发送的消息不能超过20条/分钟。
type Message struct {
	MsgType string `json:"msgtype"` //消息类型 text/markdown/link/actionCard/feedCard
	Text *Text `json:"text,omitempty"` // MsgType=text 时有效
	Markdown *Markdown `json:"markdown,omitempty"` // MsgType=markdown 时有效
	Link *Link `json:"link,omitempty"` // MsgType=link 时有效
	ActionCard *ActionCard `json:"actionCard,omitempty"` // MsgType=actionCard 时有效
	FeedCard *FeedCard `json:"feedCard,omitempty"` // MsgType=feedCard 时有效
	At *At `json:"at,omitempty"` // MsgType=text 和 MsgType=markdown  时有效
}

func (m *Message) Reset() imbot.Messager {
	m.MsgType = ``
	m.Text = nil
	m.Markdown = nil
	m.Link = nil
	m.ActionCard = nil
	m.FeedCard = nil
	m.At = nil
	return m
}

func (m *Message) SendText(url, text string) error {
	m.MsgType = `text`
	m.Text = &Text{Content:text}
	_, err := http.Send(url, m)
	return err
}

func (m *Message) SendMarkdown(url, title, markdown string) error {
	m.MsgType = `markdown`
	m.Markdown = &Markdown{Title:title,Text:markdown}
	_, err := http.Send(url, m)
	return err
}

type Text struct {
	Content string `json:"content"` //文本内容
}

type At struct {
	AtMobiles []string `json:"atMobiles"` //被@人的手机号(在content里添加@人的手机号)
	IsAtAll bool `json:"isAtAll"` //@所有人时：true，否则为：false
}

type Markdown struct {
	Title string `json:"title"` //消息标题
	Text string `json:"text"` //markdown内容
}

type Image struct {
	Base64 string `json:"base64"` //图片内容的base64编码
	Md5 string `json:"md5"` //图片内容（base64编码前）的md5值
}

type Link struct{
	Title string `json:"title"` //消息标题
	Text string `json:"text,omitempty"` //消息内容。如果太长只会部分展示
	MessageURL string `json:"messageUrl"` //点击后跳转的链接。
	PicURL string `json:"picUrl"` //图片URL
}

type FeedCard struct{
	Links []*Link `json:"links"` //图文消息，一个图文消息支持1到8条图文
}

type ActionCard struct{
	Title string `json:"title"` //消息标题
	Text string `json:"text"` //markdown格式的消息
	SingleTitle string `json:"singleTitle"` //单个按钮的方案。(设置此项和singleURL后btns无效)
	SingleURL string `json:"singleURL"` //点击singleTitle按钮触发的URL
	BtnOrientation string `json:"btnOrientation"` //0-按钮竖直排列，1-按钮横向排列
	Buttons []*Button `json:"btns,omitempty"` //按钮组
}

type Button struct {
	Title string `json:"title"` //按钮标题
	ActionURL string `json:"actionURL"` //按钮URL
}
