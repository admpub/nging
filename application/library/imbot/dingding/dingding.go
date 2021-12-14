//Package dingding 钉钉机器人
package dingding

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"time"

	"github.com/admpub/nging/v4/application/library/imbot"
	"github.com/admpub/nging/v4/application/library/imbot/http"
)

// document https://ding-doc.dingtalk.com/doc#/serverapi3/iydd5h

var (
	TOKEN   = ``
	SECRET  = ``
	API_URL = `https://oapi.dingtalk.com/robot/send?access_token=%v&timestamp=%v&sign=%v`
)

func init() {
	imbot.Register(`dingding`, `钉钉机器人`, &Message{})
}

// Message 消息实体
type Message struct {
	MsgType    string      `json:"msgtype"`              //消息类型 text/markdown/link/actionCard/feedCard
	Text       *Text       `json:"text,omitempty"`       // MsgType=text 时有效
	Markdown   *Markdown   `json:"markdown,omitempty"`   // MsgType=markdown 时有效
	Link       *Link       `json:"link,omitempty"`       // MsgType=link 时有效
	ActionCard *ActionCard `json:"actionCard,omitempty"` // MsgType=actionCard 时有效
	FeedCard   *FeedCard   `json:"feedCard,omitempty"`   // MsgType=feedCard 时有效
	At         *At         `json:"at,omitempty"`         // MsgType=text 和 MsgType=markdown  时有效
}

func (m *Message) BuildURL(args ...string) string {
	timestamp := time.Now().UnixNano() / 1e6
	secret := SECRET
	if len(args) > 0 {
		secret = args[0]
	}
	sign := fmt.Sprintf("%v\n%v", timestamp, secret)
	sign = HmacSHA256(sign, secret)
	return fmt.Sprintf(API_URL, TOKEN, timestamp, sign)
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

func (m *Message) SendText(url, text string, atMobiles ...string) error {
	m.MsgType = `text`
	m.Text = &Text{Content: text}
	m.at(atMobiles...)
	_, err := http.Send(url, m)
	return err
}

func (m *Message) at(atMobiles ...string) {
	if len(atMobiles) > 0 {
		m.At = &At{}
		if len(atMobiles) == 1 && atMobiles[0] == `@all` {
			m.At.IsAtAll = true
		} else {
			m.At.AtMobiles = atMobiles
		}
	}
}

func (m *Message) SendMarkdown(url, title, markdown string, atMobiles ...string) error {
	m.MsgType = `markdown`
	m.Markdown = &Markdown{Title: title, Text: markdown}
	m.at(atMobiles...)
	_, err := http.Send(url, m)
	return err
}

func HmacSHA256(s string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

type Text struct {
	Content string `json:"content"` //文本内容
}

type At struct {
	AtMobiles []string `json:"atMobiles"` //被@人的手机号(在content里添加@人的手机号)
	IsAtAll   bool     `json:"isAtAll"`   //@所有人时：true，否则为：false
}

type Markdown struct {
	Title string `json:"title"` //消息标题
	Text  string `json:"text"`  //markdown内容
}

type Image struct {
	Base64 string `json:"base64"` //图片内容的base64编码
	Md5    string `json:"md5"`    //图片内容（base64编码前）的md5值
}

type Link struct {
	Title      string `json:"title"`          //消息标题
	Text       string `json:"text,omitempty"` //消息内容。如果太长只会部分展示
	MessageURL string `json:"messageUrl"`     //点击后跳转的链接。
	PicURL     string `json:"picUrl"`         //图片URL
}

type FeedCard struct {
	Links []*Link `json:"links"` //图文消息，一个图文消息支持1到8条图文
}

type ActionCard struct {
	Title          string    `json:"title"`          //消息标题
	Text           string    `json:"text"`           //markdown格式的消息
	SingleTitle    string    `json:"singleTitle"`    //单个按钮的方案。(设置此项和singleURL后btns无效)
	SingleURL      string    `json:"singleURL"`      //点击singleTitle按钮触发的URL
	BtnOrientation string    `json:"btnOrientation"` //0-按钮竖直排列，1-按钮横向排列
	Buttons        []*Button `json:"btns,omitempty"` //按钮组
}

type Button struct {
	Title     string `json:"title"`     //按钮标题
	ActionURL string `json:"actionURL"` //按钮URL
}
