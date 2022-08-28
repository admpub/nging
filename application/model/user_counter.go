package model

import "github.com/webx-top/db"

func (u *User) IncrLoginFails() error {
	return u.NgingUser.UpdateField(nil, `login_fails`, db.Raw(`login_fails+1`), `id`, u.Id)
}

func (u *User) ResetLoginFails() error {
	return u.NgingUser.UpdateField(nil, `login_fails`, 0, `id`, u.Id)
}
