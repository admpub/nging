/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package notice

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"sync"
)

const (
	Succeed         = 1
	Failed          = 0
	Unauthenticated = -1
	Forbid          = -2
)

type Message struct {
	Type    string      `json:"type" xml:"type"`
	Title   string      `json:"title" xml:"title"`
	Status  int         `json:"status" xml:"status"`
	Content interface{} `json:"content" xml:"content"`
}

type Notice struct {
	Types   map[string]bool
	Message chan *Message `json:"-" xml:"-"`
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
		Types:   map[string]bool{},
		Message: make(chan *Message),
	}
}

type OnlineUser struct {
	Notice  *Notice
	Clients uint
}

func NewOnlineUser() *OnlineUser {
	return &OnlineUser{
		Notice:  NewNotice(),
		Clients: 1,
	}
}

type userNotices struct {
	Lock *sync.RWMutex
	User map[string]*OnlineUser //key: user
}

func NewUserNotices() *userNotices {
	return &userNotices{
		Lock: &sync.RWMutex{},
		User: map[string]*OnlineUser{},
	}
}

func Stdout(message *Message) {
	if message.Status == Succeed {
		os.Stdout.WriteString(fmt.Sprint(message.Content))
	} else {
		os.Stderr.WriteString(fmt.Sprint(message.Content))
	}
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

func (u *userNotices) Send(user string, message *Message) {
	u.Lock.Lock()
	defer u.Lock.Unlock()
	oUser, exists := u.User[user]
	if !exists {
		Stdout(message)
		return
	}
	if !oUser.Notice.Types[message.Type] {
		Stdout(message)
		return
	}
	oUser.Notice.Message <- message
}

func (u *userNotices) Recv(user string) <-chan *Message {
	//race...
	//u.Lock.Lock()
	//defer u.Lock.Unlock()
	oUser, exists := u.User[user]
	if !exists {
		oUser = NewOnlineUser()
		u.User[user] = oUser
	}
	return oUser.Notice.Message
}

func (u *userNotices) RecvJSON(user string) []byte {
	message := <-u.Recv(user)
	b, err := json.Marshal(message)
	if err != nil {
		return []byte(err.Error())
	}
	return b
}

func (u *userNotices) RecvXML(user string) []byte {
	message := <-u.Recv(user)
	b, err := xml.Marshal(message)
	if err != nil {
		return []byte(err.Error())
	}
	return b
}

func (u *userNotices) CloseClient(user string) bool {
	u.Lock.Lock()
	defer u.Lock.Unlock()
	oUser, exists := u.User[user]
	if !exists {
		return true
	}
	oUser.Clients--
	if oUser.Clients <= 0 {
		delete(u.User, user)
		return true
	}
	return false
}

func (u *userNotices) OpenClient(user string) {
	u.Lock.Lock()
	defer u.Lock.Unlock()
	oUser, exists := u.User[user]
	if !exists {
		oUser = NewOnlineUser()
		u.User[user] = oUser
	}
	oUser.Clients++
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
			oUser.Notice.Types[typ] = true
		}
	} else {
		for key := range oUser.Notice.Types {
			oUser.Notice.Types[key] = true
		}
	}
}

func (u *userNotices) Clear() {
	u.User = map[string]*OnlineUser{}
}
