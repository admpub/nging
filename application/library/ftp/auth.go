package ftp

import "github.com/admpub/nging/application/model"

func NewAuth() *Auth {
	return &Auth{
		FtpUser: model.NewFtpUser(nil),
	}
}

type Auth struct {
	*model.FtpUser
}
