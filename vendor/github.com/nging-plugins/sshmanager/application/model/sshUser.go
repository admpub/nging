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
	"fmt"
	"io"
	"strings"

	"github.com/admpub/go-sshclient"
	"github.com/webx-top/echo"
	"golang.org/x/crypto/ssh"

	"github.com/nging-plugins/sshmanager/application/dbschema"
)

var (
	Decode = func(r string) string { return r }
)

type SshUserAndGroup struct {
	*dbschema.NgingSshUser
	Group *dbschema.NgingSshUserGroup
}

func NewSshUser(ctx echo.Context) *SshUser {
	return &SshUser{
		NgingSshUser: dbschema.NewNgingSshUser(ctx),
	}
}

type SshUser struct {
	*dbschema.NgingSshUser
}

func (s *SshUser) ExecMultiCMD(writer io.Writer, commands ...string) error {
	if len(commands) == 0 {
		return nil
	}
	client, err := s.Connect()
	if err != nil {
		return err
	}
	defer client.Close()
	err = client.Script(strings.Join(commands, "\r\n")).SetStdio(writer, writer).Run()
	return err
}

func (s *SshUser) Connect() (*sshclient.Client, error) {
	config := &ssh.ClientConfig{
		User:            s.Username,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if len(s.PrivateKey) > 0 {
		var signer ssh.Signer
		var err error
		pemBytes := []byte(s.PrivateKey)
		if len(s.Passphrase) > 0 {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(Decode(s.Passphrase)))
		} else {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		}
		if err != nil {
			return nil, err
		}
		config.Auth = append(config.Auth, ssh.PublicKeys(signer))
	}

	if len(s.Password) > 0 {
		config.Auth = append(config.Auth, ssh.Password(Decode(s.Password)))
	}
	config.SetDefaults()
	if s.Port <= 0 {
		s.Port = 22
	}
	client, err := sshclient.Dial("tcp", fmt.Sprintf(`%s:%d`, s.Host, s.Port), config)
	if err != nil {
		return nil, err
	}
	//defer client.Close()
	return client, nil
}
