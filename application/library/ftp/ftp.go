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
package ftp

import (
	"os"
	"path/filepath"

	"github.com/admpub/log"
	filedriver "github.com/goftp/file-driver"
	ldbauth "github.com/goftp/leveldb-auth"
	ldbperm "github.com/goftp/leveldb-perm"
	qiniudriver "github.com/goftp/qiniu-driver"
	"github.com/goftp/server"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/webx-top/com"
)

var (
	DefaultConfig = &Config{
		PidFile:      `ftp.pid`,
		DBType:       `leveldb`,
		FTPStoreType: `file`,
		FTPDir:       filepath.Join(com.SelfDir(), `data/ftpdir`),
		ServerOpts: server.ServerOpts{
			Name: `TinyFTP`,
		},
	}
	DefaultPidFile = `ftp.pid`
)

func Fixed(c *Config) {
	if c.PidFile == `` {
		c.PidFile = DefaultConfig.PidFile
	}
	if c.ServerOpts.Name == `` {
		c.ServerOpts.Name = DefaultConfig.ServerOpts.Name
	}
	if c.DBType == `` {
		c.DBType = DefaultConfig.DBType
	}
	if c.FTPDir == `` {
		c.FTPDir = DefaultConfig.FTPDir
	}
	if c.FTPStoreType == `` {
		c.FTPStoreType = DefaultConfig.FTPStoreType
	}
}

type Config struct {
	perm         server.Perm
	db           *leveldb.DB
	PidFile      string            `json:"pidFile"`
	DBType       string            `json:"dbType"`
	FTPDir       string            `json:"ftpDir"`
	FTPStoreType string            `json:"ftpStoreType"`
	FTPOptions   map[string]string `json:"ftpOptions"`
	server.ServerOpts
}

func (c *Config) Init() *Config {
	c.SetPermByType(c.DBType, `root`, `root`)
	c.SetAuthByType(c.DBType)
	c.SetStoreByType(c.FTPStoreType, c.FTPDir, c.FTPOptions)
	return c
}

func (c *Config) SetDB(db *leveldb.DB) *Config {
	c.db = db
	return c
}

func (c *Config) DB() *leveldb.DB {
	if c.db == nil {
		db, err := leveldb.OpenFile(filepath.Join(com.SelfDir(), "data/ftpdb"), nil)
		if err != nil {
			panic(err)
		}
		c.db = db
	}
	return c.db
}

func (c *Config) SetPermByType(storeType string, owner string, group string) *Config {
	switch storeType {
	case `leveldb`:
		c.SetPerm(ldbperm.NewLDBPerm(c.DB(), owner, group, os.ModePerm), owner, group)
	default:
		c.SetPerm(NewPerm(owner, group, os.ModePerm), owner, group)
	}
	return c
}

func (c *Config) SetPerm(perm server.Perm, owner string, group string) *Config {
	if perm != nil {
		c.perm = perm
		return c
	}
	c.perm = server.NewSimplePerm(owner, group)
	return c
}

func (c *Config) SetAuthByType(storeType string) *Config {
	switch storeType {
	case `leveldb`:
		c.SetAuth(&ldbauth.LDBAuth{c.DB()})
	default:
		c.SetAuth(NewAuth())
	}
	return c
}

func (c *Config) SetAuth(auth server.Auth) *Config {
	c.ServerOpts.Auth = auth
	return c
}

func (c *Config) SetStoreByType(storeType string, rootPath string, options ...map[string]string) *Config {
	switch storeType {
	case "file":
		_, err := os.Lstat(rootPath)
		if os.IsNotExist(err) {
			os.MkdirAll(rootPath, os.ModePerm)
		} else if err != nil {
			log.Error(err)
			return c
		}
		factory := &filedriver.FileDriverFactory{
			rootPath,
			c.perm,
		}
		c.SetStore(factory)
	case "qiniu":
		var opts map[string]string
		if len(options) > 0 {
			opts = options[0]
		}
		accessKey, _ := opts["accessKey"]
		secretKey, _ := opts["secretKey"]
		bucket := rootPath
		factory := qiniudriver.NewQiniuDriverFactory(accessKey,
			secretKey, bucket)
		c.SetStore(factory)
	default:
		log.Fatal("Unsupported FTP storage type: " + storeType)
	}
	return c
}

func (c *Config) SetStore(store server.DriverFactory) *Config {
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
	ftpServer := server.NewServer(&c.ServerOpts)
	log.Info("Start FTP Server")
	err := ftpServer.ListenAndServe()
	if err != nil {
		if len(c.PidFile) > 0 {
			os.Remove(c.PidFile)
		}
		log.Fatal("Error starting server:", err)
	}
}
