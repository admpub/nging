/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package notice

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/admpub/nging/v3/application/library/msgbox"
)

const (
	Succeed         = 1
	Failed          = 0
	Unauthenticated = -1
	Forbid          = -2
)

func NewProgress() *Progress {
	return &Progress{
		Total:   -1,
		Finish:  -1,
		Percent: 0,
	}
}

type Progress struct {
	Total    int64   `json:"total" xml:"total"`
	Finish   int64   `json:"finish" xml:"finish"`
	Percent  float64 `json:"percent" xml:"percent"`
	Complete bool    `json:"complete" xml:"complete"`
	control  IsExited
}

type Control struct {
	exited bool
}

func (c *Control) IsExited() bool {
	return c.exited
}

func (c *Control) Exited() *Control {
	c.exited = true
	return c
}

func (c *Control) ListenContextAndTimeout(ctx context.Context, timeouts ...time.Duration) *Control {
	timeout := 24 * time.Hour
	if len(timeouts) > 0 && timeouts[0] != 0 {
		timeout = timeouts[0]
	}
	t := time.NewTicker(timeout)
	defer t.Stop()
	go func() {
		for {
			select {
			case <-ctx.Done():
				c.Exited()
				return
			case <-t.C:
				c.Exited()
				return
			}
		}
	}()
	return c
}

type IsExited interface {
	IsExited() bool
}

func (p *Progress) IsExited() bool {
	if p.control == nil {
		return false
	}
	return p.control.IsExited()
}

func (p *Progress) SetControl(control IsExited) *Progress {
	p.control = control
	return p
}

func (p *Progress) CalcPercent() *Progress {
	if p.Total > 0 {
		p.Percent = (float64(p.Finish) / float64(p.Total)) * 100
		if p.Percent < 0 {
			p.Percent = 0
		}
	} else if p.Total == 0 {
		p.Percent = 100
	} else {
		p.Percent = 0
	}
	return p
}

type Message struct {
	ClientID string      `json:"client_id" xml:"client_id"`
	ID       interface{} `json:"id" xml:"id"`
	Type     string      `json:"type" xml:"type"`
	Title    string      `json:"title" xml:"title"`
	Status   int         `json:"status" xml:"status"`
	Content  interface{} `json:"content" xml:"content"`
	Mode     string      `json:"mode" xml:"mode"` //显示模式：notify/element/modal
	Progress *Progress   `json:"progress" xml:"progress"`
}

func (m *Message) SetType(t string) *Message {
	m.Type = t
	return m
}

func (m *Message) SetTitle(title string) *Message {
	m.Title = title
	return m
}

func (m *Message) SetID(id interface{}) *Message {
	m.ID = id
	return m
}

func (m *Message) SetClientID(clientID string) *Message {
	m.ClientID = clientID
	return m
}

func (m *Message) SetStatus(status int) *Message {
	m.Status = status
	return m
}

func (m *Message) SetContent(content interface{}) *Message {
	m.Content = content
	return m
}

func (m *Message) SetMode(mode string) *Message {
	m.Mode = mode
	return m
}

func (m *Message) SetProgress(progress *Progress) *Message {
	m.Progress = progress
	if m.Progress != nil && m.Progress.Percent == 0 {
		m.CalcPercent()
	}
	return m
}

func (m *Message) SetProgressValue(finish int64, total int64) *Message {
	if m.Progress == nil {
		m.Progress = NewProgress()
	}
	m.Progress.Finish = finish
	m.Progress.Total = total
	m.CalcPercent()
	return m
}

func (m *Message) CalcPercent() *Message {
	m.Progress.CalcPercent()
	return m
}

type Notice struct {
	Types    map[string]bool
	Messages map[string]chan *Message `json:"-" xml:"-"`
}

func (a *Notice) CountClient() int {
	return len(a.Messages)
}

func (a *Notice) CloseClient(clientID string) {
	if msg, ok := a.Messages[clientID]; ok {
		close(msg)
		delete(a.Messages, clientID)
	}
}

func (a *Notice) OpenClient(clientID string) {
	if _, ok := a.Messages[clientID]; !ok {
		a.Messages[clientID] = make(chan *Message)
	}
}

func NewMessageWithValue(typ string, title string, content interface{}, status ...int) *Message {
	st := Succeed
	if len(status) > 0 {
		st = status[0]
	}
	return &Message{
		Type:    typ,
		Title:   title,
		Status:  st,
		Content: content,
	}
}

func NewMessage() *Message {
	return &Message{}
}

func NewNotice() *Notice {
	return &Notice{
		Types:    map[string]bool{},
		Messages: map[string](chan *Message){},
	}
}

type OnlineUser struct {
	*Notice
}

func NewOnlineUser() *OnlineUser {
	return &OnlineUser{
		Notice: NewNotice(),
	}
}

type userNotices struct {
	Lock    *sync.RWMutex
	User    map[string]*OnlineUser //key: user
	Debug   bool
	onClose []func(user string)
	onOpen  []func(user string)
}

func NewUserNotices(debug bool) *userNotices {
	return &userNotices{
		Lock:    &sync.RWMutex{},
		User:    map[string]*OnlineUser{},
		Debug:   debug,
		onClose: []func(user string){},
		onOpen:  []func(user string){},
	}
}

func Stdout(message *Message) {
	if message.Status == Succeed {
		os.Stdout.WriteString(fmt.Sprint(message.Content))
	} else {
		os.Stderr.WriteString(fmt.Sprint(message.Content))
	}
}

func (u *userNotices) OnClose(fn ...func(user string)) *userNotices {
	u.onClose = append(u.onClose, fn...)
	return u
}

func (u *userNotices) OnOpen(fn ...func(user string)) *userNotices {
	u.onOpen = append(u.onOpen, fn...)
	return u
}

func (u *userNotices) Sendable(user string, types ...string) bool {
	oUser, exists := u.User[user]
	if !exists {
		return false
	}
	for _, typ := range types {
		if !oUser.Notice.Types[typ] {
			return false
		}
	}
	return true
}

func (u *userNotices) Send(user string, message *Message) error {
	if u.Debug {
		msgbox.Debug(`[NOTICE]`, `[Send][FindUser]: `+user)
	}
	u.Lock.Lock()
	defer u.Lock.Unlock()
	oUser, exists := u.User[user]
	if !exists {
		if u.Debug {
			msgbox.Debug(`[NOTICE]`, `[Send][NotFoundUser]: `+user)
		}
		Stdout(message)
		return ErrUserNotOnline
	}
	if u.Debug {
		msgbox.Debug(`[NOTICE]`, `[Send][CheckRecvType]: `+message.Type+` (for user: `+user+`)`)
	}
	if !oUser.Notice.Types[message.Type] {
		Stdout(message)
		return ErrMsgTypeNotAccept
	}
	if u.Debug {
		msgbox.Debug(`[NOTICE]`, `[Send][MessageTo]: `+user)
	}
	msg, ok := oUser.Notice.Messages[message.ClientID]
	if ok {
		msg <- message
		return nil
	}
	/*
		for clientID, msg := range oUser.Notice.Messages {
			msg <- message
			if u.Debug {
				msgbox.Debug(`[NOTICE]`, `[Send][MessageTo]: `+user+` [ClientID]: `+fmt.Sprint(clientID))
			}
			return
		}
	*/
	if u.Debug {
		msgbox.Debug(`[NOTICE]`, `[Send][MessageTo]: `+user+` [NotFoundClientID]: `+fmt.Sprint(message.ClientID))
	}
	return ErrClientIDNotOnline
}

func (u *userNotices) Recv(user string, clientID string) <-chan *Message {
	//race...
	//u.Lock.Lock()
	//defer u.Lock.Unlock()
	oUser, exists := u.User[user]
	if !exists {
		oUser = NewOnlineUser()
		u.User[user] = oUser
	}
	msg, ok := oUser.Notice.Messages[clientID]
	if ok {
		return msg
	}
	return nil
}

func (u *userNotices) RecvJSON(user string, clientID string) ([]byte, error) {
	if u.Debug {
		msgbox.Warn(`[NOTICE]`, `[RecvJSON][Waiting]: `+user)
	}
	message := <-u.Recv(user, clientID)
	if message == nil {
		return nil, nil
	}
	message.ClientID = clientID
	b, err := json.Marshal(message)
	if err != nil {
		return b, err
	}
	if u.Debug {
		msgbox.Warn(`[NOTICE]`, `[RecvJSON][Received]: `+user)
	}
	return b, err
}

func (u *userNotices) RecvXML(user string, clientID string) ([]byte, error) {
	if u.Debug {
		msgbox.Warn(`[NOTICE]`, `[RecvXML][Waiting]: `+user)
	}
	message := <-u.Recv(user, clientID)
	if message == nil {
		return nil, nil
	}
	message.ClientID = clientID
	b, err := xml.Marshal(message)
	if err != nil {
		return b, err
	}
	if u.Debug {
		msgbox.Warn(`[NOTICE]`, `[RecvXML][Received]: `+user)
	}
	return b, err
}

func (u *userNotices) CloseClient(user string, clientID string) bool {
	u.Lock.Lock()
	defer u.Lock.Unlock()
	oUser, exists := u.User[user]
	if !exists {
		return true
	}
	oUser.CloseClient(clientID)
	if u.Debug {
		msgbox.Error(`[NOTICE]`, `[CloseClient][ClientID]: `+clientID)
	}
	if len(oUser.Notice.Messages) < 1 {
		for key, msg := range oUser.Notice.Messages {
			close(msg)
			delete(oUser.Notice.Messages, key)
		}
		delete(u.User, user)
		for _, fn := range u.onClose {
			fn(user)
		}
		return true
	}
	return false
}

func (u *userNotices) OpenClient(user string) string {
	u.Lock.Lock()
	defer u.Lock.Unlock()
	oUser, exists := u.User[user]
	if !exists {
		oUser = NewOnlineUser()
		u.User[user] = oUser
		for _, fn := range u.onOpen {
			fn(user)
		}
	}
	clientID := fmt.Sprint(time.Now().Unix())
	oUser.OpenClient(clientID)
	return clientID
}

func (u *userNotices) CloseMessage(user string, types ...string) {
	oUser, exists := u.User[user]
	if !exists {
		return
	}
	if len(types) > 0 {
		for _, typ := range types {
			_, ok := oUser.Notice.Types[typ]
			if !ok {
				continue
			}
			delete(oUser.Notice.Types, typ)
		}
	} else {
		oUser.Notice.Types = map[string]bool{}
	}
}

func (u *userNotices) OpenMessage(user string, types ...string) {
	oUser, exists := u.User[user]
	if !exists {
		oUser = NewOnlineUser()
		u.User[user] = oUser
	}
	if len(types) > 0 {
		for _, typ := range types {
			if opened, _ := oUser.Notice.Types[typ]; !opened {
				oUser.Notice.Types[typ] = true
			}
		}
	} else {
		for key := range oUser.Notice.Types {
			if opened, _ := oUser.Notice.Types[key]; !opened {
				oUser.Notice.Types[key] = true
			}
		}
	}
}

func (u *userNotices) Clear() {
	u.User = map[string]*OnlineUser{}
}
