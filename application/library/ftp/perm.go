package ftp

import (
	"os"

	ftpserver "github.com/admpub/ftpserver"
)

func NewPerm(owner, group string) *Perm {
	return &Perm{
		SimplePerm: ftpserver.NewSimplePerm(owner, group),
	}
}

type Perm struct {
	*ftpserver.SimplePerm
}

func (s *Perm) GetMode(string) (os.FileMode, error) {
	return 0, nil
}
