/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/admpub/nging/v5/application/library/msgbox"
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

type userNotices struct {
	users    *OnlineUsers //key: user
	_debug   bool
	_debugMu sync.RWMutex
	onClose  []func(user string)
	onOpen   []func(user string)
}

func NewUserNotices(debug bool) *userNotices {
	return &userNotices{
		users:   NewOnlineUsers(),
		_debug:  debug,
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

func (u *userNotices) SetDebug(on bool) *userNotices {
	u._debugMu.Lock()
	u._debug = on
	u._debugMu.Unlock()
	return u
}

func (u *userNotices) Debug() bool {
	u._debugMu.RLock()
	debug := u._debug
	u._debugMu.RUnlock()
	return debug
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
	oUser, exists := u.users.GetOk(user)
	if !exists {
		return false
	}
	return oUser.Notice.types.Has(types...)
}

func (u *userNotices) Send(user string, message *Message) error {
	debug := u.Debug()
	if debug {
		msgbox.Debug(`[NOTICE]`, `[Send][FindUser]: `+user)
	}
	oUser, exists := u.users.GetOk(user)
	if !exists {
		if debug {
			msgbox.Debug(`[NOTICE]`, `[Send][NotFoundUser]: `+user)
		}
		Stdout(message)
		return ErrUserNotOnline
	}
	if debug {
		msgbox.Debug(`[NOTICE]`, `[Send][CheckRecvType]: `+message.Type+` (for user: `+user+`)`)
	}
	return oUser.Send(message)
}

func (u *userNotices) Recv(user string, clientID string) <-chan *Message {
	oUser, exists := u.users.GetOk(user)
	if !exists {
		oUser = NewOnlineUser(user)
		u.users.Set(user, oUser)
	}
	return oUser.Recv(clientID)
}

func (u *userNotices) CloseClient(user string, clientID string) bool {
	oUser, exists := u.users.GetOk(user)
	if !exists {
		return true
	}
	oUser.CloseClient(clientID)
	if u.Debug() {
		msgbox.Info(`[NOTICE]`, `[CloseClient][ClientID]: `+clientID)
	}
	if oUser.Notice.messages.Size() < 1 {
		oUser.Notice.messages.Clear()
		u.users.Delete(user)
		for _, fn := range u.onClose {
			fn(user)
		}
		return true
	}
	return false
}

func (u *userNotices) IsOnline(user string) bool {
	_, exists := u.users.GetOk(user)
	return exists
}

func (u *userNotices) OpenClient(user string) (oUser *OnlineUser, clientID string) {
	var exists bool
	oUser, exists = u.users.GetOk(user)
	if !exists {
		oUser = NewOnlineUser(user)
		u.users.Set(user, oUser)
		for _, fn := range u.onOpen {
			fn(user)
		}
	}
	clientID = fmt.Sprint(time.Now().UnixMilli())
	oUser.OpenClient(clientID)
	if u.Debug() {
		msgbox.Info(`[NOTICE]`, `[OpenClient][ClientID]: `+clientID)
	}
	return
}

func (u *userNotices) CloseMessage(user string, types ...string) {
	oUser, exists := u.users.GetOk(user)
	if !exists {
		return
	}
	oUser.Notice.types.Clear(types...)
}

func (u *userNotices) OpenMessage(user string, types ...string) {
	oUser, exists := u.users.GetOk(user)
	if !exists {
		oUser = NewOnlineUser(user)
		u.users.Set(user, oUser)
	}
	oUser.Notice.types.Open(types...)
}

func (u *userNotices) Clear() {
	u.users.Clear()
}
