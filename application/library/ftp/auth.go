package ftp

import (
	"github.com/admpub/caddyui/application/dbschema"
	"github.com/admpub/caddyui/application/model"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
)

func NewAuth() *Auth {
	return &Auth{
		FtpUser: &dbschema.FtpUser{},
	}
}

type Auth struct {
	*dbschema.FtpUser
}

func (f *Auth) CheckPasswd(username string, password string) (bool, error) {
	n, e := f.Param().SetArgs(db.Cond{`username`: username, `password`: com.MakePassword(password, model.DefaultSalt)}).Count()
	return n > 0, e
}
