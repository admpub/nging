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

type Notice struct {
	types    *noticeTypes
	messages *noticeMessages
}

func (a *Notice) CountClient() int {
	return a.messages.Size()
}

func (a *Notice) CloseClient(clientID string) {
	a.messages.Delete(clientID)
}

func (a *Notice) OpenClient(clientID string) {
	a.messages.Add(clientID)
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
		types:    newNoticeTypes(),
		messages: newNoticeMessages(),
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
	lock    sync.RWMutex
	user    map[string]*OnlineUser //key: user
	debug   bool
	onClose []func(user string)
	onOpen  []func(user string)
}

func NewUserNotices(debug bool) *userNotices {
	return &userNotices{
		user:    map[string]*OnlineUser{},
		debug:   debug,
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
	u.lock.RLock()
	oUser, exists := u.user[user]
	u.lock.RUnlock()
	if !exists {
		return false
	}
	return oUser.Notice.types.Has(types...)
}

func (u *userNotices) Send(user string, message *Message) error {
	if u.debug {
		msgbox.Debug(`[NOTICE]`, `[Send][FindUser]: `+user)
	}
	u.lock.RLock()
	oUser, exists := u.user[user]
	u.lock.RUnlock()
	if !exists {
		if u.debug {
			msgbox.Debug(`[NOTICE]`, `[Send][NotFoundUser]: `+user)
		}
		Stdout(message)
		return ErrUserNotOnline
	}
	if u.debug {
		msgbox.Debug(`[NOTICE]`, `[Send][CheckRecvType]: `+message.Type+` (for user: `+user+`)`)
	}
	if !oUser.Notice.types.Has(message.Type) {
		Stdout(message)
		return ErrMsgTypeNotAccept
	}
	if u.debug {
		msgbox.Debug(`[NOTICE]`, `[Send][MessageTo]: `+user)
	}
	err := oUser.Notice.messages.Send(message)
	if err != nil && u.debug {
		msgbox.Debug(`[NOTICE]`, `[Send][MessageTo]: `+user+` [NotFoundClientID]: `+fmt.Sprint(message.ClientID))
	}
	return err
}

func (u *userNotices) Recv(user string, clientID string) chan *Message {
	u.lock.RLock()
	oUser, exists := u.user[user]
	u.lock.RUnlock()
	if !exists {
		oUser = NewOnlineUser()
		u.lock.Lock()
		u.user[user] = oUser
		u.lock.Unlock()
	}
	return oUser.Notice.messages.Recv(clientID)
}

func (u *userNotices) RecvJSON(user string, clientID string) ([]byte, error) {
	if u.debug {
		msgbox.Warn(`[NOTICE]`, `[RecvJSON][Waiting]: `+user)
	}
	msgChan := u.Recv(user, clientID)
	if msgChan == nil {
		return nil, nil
	}
	message := <-msgChan
	if message == nil {
		return nil, nil
	}
	message.ClientID = clientID
	b, err := json.Marshal(message)
	if err != nil {
		return b, err
	}
	if u.debug {
		msgbox.Warn(`[NOTICE]`, `[RecvJSON][Received]: `+user)
	}
	return b, err
}

func (u *userNotices) RecvXML(user string, clientID string) ([]byte, error) {
	if u.debug {
		msgbox.Warn(`[NOTICE]`, `[RecvXML][Waiting]: `+user)
	}
	msgChan := u.Recv(user, clientID)
	if msgChan == nil {
		return nil, nil
	}
	message := <-msgChan
	if message == nil {
		return nil, nil
	}
	message.ClientID = clientID
	b, err := xml.Marshal(message)
	if err != nil {
		return b, err
	}
	if u.debug {
		msgbox.Warn(`[NOTICE]`, `[RecvXML][Received]: `+user)
	}
	return b, err
}

func (u *userNotices) CloseClient(user string, clientID string) bool {
	u.lock.RLock()
	oUser, exists := u.user[user]
	u.lock.RUnlock()
	if !exists {
		return true
	}
	oUser.CloseClient(clientID)
	if u.debug {
		msgbox.Info(`[NOTICE]`, `[CloseClient][ClientID]: `+clientID)
	}
	if oUser.Notice.messages.Size() < 1 {
		oUser.Notice.messages.Clear()
		u.lock.Lock()
		delete(u.user, user)
		u.lock.Unlock()
		for _, fn := range u.onClose {
			fn(user)
		}
		return true
	}
	return false
}

func (u *userNotices) OpenClient(user string) string {
	u.lock.RLock()
	oUser, exists := u.user[user]
	u.lock.RUnlock()
	if !exists {
		oUser = NewOnlineUser()
		u.lock.Lock()
		u.user[user] = oUser
		u.lock.Unlock()
		for _, fn := range u.onOpen {
			fn(user)
		}
	}
	clientID := fmt.Sprint(time.Now().UnixMilli())
	oUser.OpenClient(clientID)
	if u.debug {
		msgbox.Info(`[NOTICE]`, `[OpenClient][ClientID]: `+clientID)
	}
	return clientID
}

func (u *userNotices) CloseMessage(user string, types ...string) {
	u.lock.RLock()
	oUser, exists := u.user[user]
	u.lock.RUnlock()
	if !exists {
		return
	}
	oUser.Notice.types.Clear(types...)
}

func (u *userNotices) OpenMessage(user string, types ...string) {
	u.lock.RLock()
	oUser, exists := u.user[user]
	u.lock.RUnlock()
	if !exists {
		oUser = NewOnlineUser()
		u.lock.Lock()
		u.user[user] = oUser
		u.lock.Unlock()
	}
	oUser.Notice.types.Open(types...)
}

func (u *userNotices) Clear() {
	u.lock.Lock()
	u.user = map[string]*OnlineUser{}
	u.lock.Unlock()
}
