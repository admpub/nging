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
package model

import (
	"errors"

	"encoding/gob"

	"github.com/admpub/nging/application/dbschema"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func init() {
	gob.Register(&dbschema.User{})
}

func NewUser(ctx echo.Context) *User {
	return &User{
		User: &dbschema.User{},
		Base: &Base{Context: ctx},
	}
}

type User struct {
	*dbschema.User
	*Base
}

func (u *User) Exists(username string) (bool, error) {
	n, e := u.Param().SetArgs(db.Cond{`username`: username}).Count()
	return n > 0, e
}

func (u *User) CheckPasswd(username string, password string) (exists bool, err error) {
	exists = true
	err = u.Get(nil, `username`, username)
	if err != nil {
		if err == db.ErrNoMoreRows {
			exists = false
		}
		return
	}
	if u.User.Disabled == `Y` {
		err = errors.New(u.T(`该用户已被禁用`))
		return
	}
	if u.User.Password != com.MakePassword(password, u.User.Salt) {
		err = errors.New(u.T(`密码不正确`))
	}
	return
}

func (u *User) NeedCheckU2F(uid uint) bool {
	u2f := &dbschema.UserU2f{}
	n, _ := u2f.Count(nil, `uid`, uid)
	return n > 0
}

func (u *User) GetUserAllU2F(uid uint) ([]*dbschema.UserU2f, error) {
	u2f := &dbschema.UserU2f{}
	all := []*dbschema.UserU2f{}
	_, err := u2f.ListByOffset(&all, nil, 0, -1, `uid`, uid)
	return all, err
}

func (u *User) U2F(uid uint, typ string) (u2f *dbschema.UserU2f, err error) {
	u2f = &dbschema.UserU2f{}
	err = u2f.Get(nil, db.And(db.Cond{`uid`: uid}, db.Cond{`type`: typ}))
	return
}

func (u *User) Register(user, pass, email string) error {
	userSchema := &dbschema.User{}
	userSchema.Username = user
	userSchema.Email = email
	userSchema.Salt = com.Salt()
	userSchema.Password = com.MakePassword(pass, userSchema.Salt)
	userSchema.Disabled = `N`
	_, err := userSchema.Add()
	u.User = userSchema
	return err
}

func (u *User) SetSession(users ...*dbschema.User) {
	var user dbschema.User
	if len(users) > 0 {
		user = *(users[0])
	} else {
		user = *(u.User)
	}
	user.Password = ``
	user.Salt = ``
	u.Context.Session().Set(`user`, &user)
}
