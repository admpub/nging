package notice

import (
	"fmt"
	"sync"

	"github.com/admpub/nging/v5/application/library/msgbox"
)

type OnlineUser struct {
	User string
	*Notice
}

func (oUser *OnlineUser) Send(message *Message) error {
	if !oUser.Notice.types.Has(message.Type) {
		Stdout(message)
		return ErrMsgTypeNotAccept
	}
	if debug {
		msgbox.Debug(`[NOTICE]`, `[Send][MessageTo]: `+oUser.User)
	}
	err := oUser.Notice.messages.Send(message)
	if err != nil && debug {
		msgbox.Debug(`[NOTICE]`, `[Send][MessageTo]: `+oUser.User+` [NotFoundClientID]: `+fmt.Sprint(message.ClientID))
	}
	return err
}

func (oUser *OnlineUser) Recv(clientID string) <-chan *Message {
	return oUser.Notice.messages.Recv(clientID)
}

func NewOnlineUser(user string) *OnlineUser {
	return &OnlineUser{
		User:   user,
		Notice: NewNotice(),
	}
}

func NewOnlineUsers() *OnlineUsers {
	return &OnlineUsers{
		user: map[string]*OnlineUser{},
	}
}

type OnlineUsers struct {
	lock sync.RWMutex
	user map[string]*OnlineUser //key: user
}

func (o *OnlineUsers) GetOk(user string, noLock ...bool) (*OnlineUser, bool) {
	if len(noLock) > 0 && noLock[0] {
		oUser, exists := o.user[user]
		return oUser, exists
	}
	o.lock.RLock()
	oUser, exists := o.user[user]
	o.lock.RUnlock()
	return oUser, exists
}

func (o *OnlineUsers) Set(user string, oUser *OnlineUser) {
	o.lock.Lock()
	o.user[user] = oUser
	o.lock.Unlock()
}

func (o *OnlineUsers) Delete(user string) {
	o.lock.Lock()
	delete(o.user, user)
	o.lock.Unlock()
}

func (o *OnlineUsers) Clear() {
	o.lock.Lock()
	o.user = map[string]*OnlineUser{}
	o.lock.Unlock()
}
