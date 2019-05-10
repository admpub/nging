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

package ftp

import (
	"os"
	"path/filepath"

	ftpserver "github.com/admpub/ftpserver"
	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var (
	DefaultConfig = &Config{
		PidFile:      `ftp.pid`,
		DBType:       `mysql`,
		FTPStoreType: `file`,
		ServerOpts: ftpserver.ServerOpts{
			Name:           `TinyFTP`,
			PassivePorts:   `6001-7000`,
			Port:           25,
			PublicIp:       `127.0.0.1`,
			WelcomeMessage: `Welcome to the TinyFTP`,
		},
	}
	DefaultPidFile = `ftp.pid`
)

func Fixed(c *Config) {
	//if c.PidFile == `` {
	pidFile := filepath.Join(echo.Wd(), `data/pid/ftp`)
	if !com.IsDir(pidFile) {
		err := os.MkdirAll(pidFile, os.ModePerm)
		if err != nil {
			log.Error(err)
		}
	}
	pidFile = filepath.Join(pidFile, DefaultPidFile)
	c.PidFile = pidFile
	//}
	if c.ServerOpts.Name == `` {
		c.ServerOpts.Name = DefaultConfig.ServerOpts.Name
	}
	if c.ServerOpts.PassivePorts == `` {
		c.ServerOpts.PassivePorts = DefaultConfig.ServerOpts.PassivePorts
	}
	if c.ServerOpts.Port == 0 {
		c.ServerOpts.Port = DefaultConfig.ServerOpts.Port
	}
	if c.ServerOpts.PublicIp == `` {
		c.ServerOpts.PublicIp = DefaultConfig.ServerOpts.PublicIp
	}
	if c.ServerOpts.WelcomeMessage == `` {
		c.ServerOpts.WelcomeMessage = DefaultConfig.ServerOpts.WelcomeMessage
	}
	if c.DBType == `` {
		c.DBType = DefaultConfig.DBType
	}
	if c.FTPStoreType == `` {
		c.FTPStoreType = DefaultConfig.FTPStoreType
	}
}

type Config struct {
	perm         ftpserver.Perm
	PidFile      string            `json:"-"`
	DBType       string            `json:"dbType"`
	FTPStoreType string            `json:"ftpStoreType"`
	FTPOptions   map[string]string `json:"ftpOptions"`
	ftpserver.ServerOpts
}

func (c *Config) Init() *Config {
	c.SetPermByType(c.DBType, `root`, `root`)
	c.SetAuthByType(c.DBType)
	c.SetStoreByType(c.FTPStoreType, c.FTPOptions)
	return c
}

func (c *Config) SetPermByType(storeType string, owner string, group string) *Config {
	switch storeType {
	default:
		c.SetPerm(nil, owner, group)
	}
	return c
}

func (c *Config) SetPerm(perm ftpserver.Perm, owner string, group string) *Config {
	if perm != nil {
		c.perm = perm
		return c
	}
	c.perm = NewPerm(owner, group)
	return c
}

func (c *Config) SetAuthByType(storeType string) *Config {
	switch storeType {
	default:
		c.SetAuth(NewAuth())
	}
	return c
}

func (c *Config) SetAuth(auth ftpserver.Auth) *Config {
	c.ServerOpts.Auth = auth
	return c
}

func (c *Config) SetStoreByType(storeType string, options ...map[string]string) *Config {
	switch storeType {
	case "file":
		factory := &FileDriverFactory{c.perm}
		c.SetStore(factory)
	default:
		log.Fatal("Unsupported FTP storage type: " + storeType)
	}
	return c
}

func (c *Config) SetStore(store ftpserver.DriverFactory) *Config {
	c.ServerOpts.Factory = store
	return c
}

func (c *Config) SetPort(port int) *Config {
	c.ServerOpts.Port = port
	return c
}

// Start start ftp server
func (c *Config) Start() {
	if len(c.PidFile) > 0 {
		err := com.WritePidFile(c.PidFile)
		if err != nil {
			log.Error(err.Error())
		}
	}
	ftpServer := ftpserver.NewServer(&c.ServerOpts)
	log.Info("Start FTP Server")
	err := ftpServer.ListenAndServe()
	if err != nil {
		if len(c.PidFile) > 0 {
			os.Remove(c.PidFile)
		}
		log.Fatal("Error starting server:", err)
	}
}
