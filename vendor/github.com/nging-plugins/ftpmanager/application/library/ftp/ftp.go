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

	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/nging-plugins/ftpmanager/application/model"
	ftpserver "goftp.io/server/v2"
)

var (
	DefaultConfig = &Config{
		PidFile:      `ftp.pid`,
		DBType:       `mysql`,
		FTPStoreType: `file`,
		// ftpserver.Options
		Name:           `TinyFTP`,
		PassivePorts:   `6001-7000`,
		Port:           25,
		PublicIP:       `127.0.0.1`,
		WelcomeMessage: `Welcome to the TinyFTP`,
	}
	DefaultPidFile = `ftp.pid`
)

func SetDefaults(c *Config) {
	pidFile := filepath.Join(echo.Wd(), `data/pid/ftp`)
	err := com.MkdirAll(pidFile, os.ModePerm)
	if err != nil {
		log.Error(err)
	}
	pidFile = filepath.Join(pidFile, DefaultPidFile)
	c.PidFile = pidFile
	if c.Name == `` {
		c.Name = DefaultConfig.Name
	}
	if c.PassivePorts == `` {
		c.PassivePorts = DefaultConfig.PassivePorts
	}
	if c.Port == 0 {
		c.Port = DefaultConfig.Port
	}
	if c.PublicIP == `` {
		c.PublicIP = DefaultConfig.PublicIP
	}
	if c.WelcomeMessage == `` {
		c.WelcomeMessage = DefaultConfig.WelcomeMessage
	}
	if c.DBType == `` {
		c.DBType = DefaultConfig.DBType
	}
	if c.FTPStoreType == `` {
		c.FTPStoreType = DefaultConfig.FTPStoreType
	}
}

type Config struct {
	PidFile      string            `json:"-"`
	DBType       string            `json:"dbType"`
	FTPStoreType string            `json:"ftpStoreType"`
	FTPOptions   map[string]string `json:"ftpOptions"`

	// Server Name, Default is Go Ftp Server
	Name string `json:"name"`

	// The hostname that the FTP server should listen on. Optional, defaults to
	// "::", which means all hostnames on ipv4 and ipv6.
	Hostname string `json:"hostName"`

	// Public IP of the server
	PublicIP string `json:"publicIP"`

	// Passive ports
	PassivePorts string `json:"passivePorts"`

	// The port that the FTP should listen on. Optional, defaults to 3000. In
	// a production environment you will probably want to change this to 21.
	Port int `json:"port"`

	// use tls, default is false
	TLS bool `json:"tls"`

	// if tls used, cert file is required
	CertFile string `json:"certFile"`

	// if tls used, key file is required
	KeyFile string `json:"keyFile"`

	// If ture TLS is used in RFC4217 mode
	ExplicitFTPS bool `json:"explicitFTPS"`

	WelcomeMessage string `json:"welcomeMessage"`

	options ftpserver.Options
}

func (c *Config) Init() *Config {
	c.options.Name = c.Name
	c.options.Hostname = c.Hostname
	c.options.PublicIP = c.PublicIP
	c.options.PassivePorts = c.PassivePorts
	c.options.Port = c.Port
	c.options.TLS = c.TLS
	c.options.CertFile = c.CertFile
	c.options.KeyFile = c.KeyFile
	c.options.ExplicitFTPS = c.ExplicitFTPS
	c.options.WelcomeMessage = c.WelcomeMessage

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
		c.options.Perm = perm
		return c
	}
	c.options.Perm = NewPerm(owner, group)
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
	c.options.Auth = auth
	return c
}

func (c *Config) SetStoreByType(storeType string, options ...map[string]string) *Config {
	switch storeType {
	case "file":
		driver := &FileDriver{user: model.NewFtpUser(nil), Perm: c.options.Perm}
		c.SetDriver(driver)
	default:
		log.Fatal("Unsupported FTP storage type: " + storeType)
	}
	return c
}

func (c *Config) SetDriver(driver ftpserver.Driver) *Config {
	c.options.Driver = driver
	return c
}

func (c *Config) SetPort(port int) *Config {
	c.Port = port
	c.options.Port = port
	return c
}

// Start start ftp server
func (c *Config) Start() error {
	if len(c.PidFile) > 0 {
		err := com.WritePidFile(c.PidFile)
		if err != nil {
			log.Error(err.Error())
		}
	}
	ftpServer, err := ftpserver.NewServer(&c.options)
	if err != nil {
		log.Fatal("Error starting server:", err)
		return err
	}
	log.Info("Start FTP Server")
	err = ftpServer.ListenAndServe()
	if err != nil {
		if len(c.PidFile) > 0 {
			os.Remove(c.PidFile)
		}
		log.Fatal("Error starting server:", err)
	}
	return err
}
