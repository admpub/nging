package notice

import "sync"

func NewOnlineUsers() *OnlineUsers {
	return &OnlineUsers{
		user: map[string]*OnlineUser{},
	}
}

type OnlineUsers struct {
	lock sync.RWMutex
	user map[string]*OnlineUser //key: user
}

func (o *OnlineUsers) GetOk(user string) (*OnlineUser, bool) {
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
