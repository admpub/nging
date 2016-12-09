package ftp

import "github.com/admpub/caddyui/application/model"

func NewAuth() *Auth {
	return &Auth{
		FtpUser: model.NewFtpUser(nil),
	}
}

type Auth struct {
	*model.FtpUser
}
