/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

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
package model

import (
	"io"
	"strings"

	colorable "github.com/mattn/go-colorable"
	"github.com/webx-top/echo"
	stdSSH "golang.org/x/crypto/ssh"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/model/base"
	"github.com/admpub/web-terminal/library/ssh"
)

var (
	Decode = func(r string) string { return r }
	_      = stdSSH.CS7
)

type SshUserAndGroup struct {
	*dbschema.NgingSshUser
	Group *dbschema.NgingSshUserGroup
}

func NewSshUser(ctx echo.Context) *SshUser {
	return &SshUser{
		NgingSshUser: &dbschema.NgingSshUser{},
		Base:         base.New(ctx),
	}
}

type SshUser struct {
	*dbschema.NgingSshUser
	*base.Base
}

func (s *SshUser) ExecMultiCMD(writer io.Writer, commands ...string) error {
	multiCMD := strings.Join(commands, "\n")
	reader := strings.NewReader(multiCMD + "\nexit\n")
	return s.Connect(reader, writer)
}

func (s *SshUser) Connect(reader io.Reader, writer io.Writer) error {
	account := &ssh.AccountConfig{
		User:     s.Username,
		Password: Decode(s.Password),
	}
	if len(s.PrivateKey) > 0 {
		account.PrivateKey = []byte(s.PrivateKey)
	}
	if len(s.Passphrase) > 0 {
		account.Passphrase = []byte(Decode(s.Passphrase))
	}
	config, err := ssh.NewSSHConfig(reader, writer, account)
	if err != nil {
		return err
	}
	client := ssh.New(config)
	err = client.Connect(s.Host, s.Port)
	if err != nil {
		return err
	}
	session := client.Session
	defer client.Close()
	// Set up terminal modes
	modes := stdSSH.TerminalModes{
		stdSSH.ECHO:          0,     // enable echoing
		stdSSH.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		stdSSH.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	columns := 120
	rows := 80
	// Request pseudo terminal
	if err = session.RequestPty("xterm", rows, columns, modes); err != nil {
		return err
	}
	writer = colorable.NewNonColorable(writer)
	session.Stdout = writer
	session.Stderr = writer
	session.Stdin = reader
	if err = session.Shell(); nil != err {
		return err
	}
	if err = session.Wait(); nil != err {
		return err
	}
	_ = colorable.NewNonColorable
	return nil
}
