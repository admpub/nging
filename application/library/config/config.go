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
package config

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"io/ioutil"
	"os"

	"strconv"

	"github.com/admpub/caddyui/application/library/caddy"
	"github.com/admpub/confl"
	"github.com/admpub/log"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/mongo"
	"github.com/webx-top/db/mysql"
)

type CLIConfig struct {
	Port int
	Conf string
	Type string //启动类型: server/manager
	cmd  *exec.Cmd
}

func (c *CLIConfig) InitFlag() {
	flag.IntVar(&c.Port, `p`, 9999, `port`)
	flag.StringVar(&c.Conf, `c`, `config/config.yaml`, `config`)
	flag.StringVar(&c.Type, `t`, `manager`, `operation type`)
}

func dialAddress(address string, timeOut int, args ...func() bool) (err error) {
	seconds := 0
	var fn func() bool
	if len(args) > 0 {
		fn = args[0]
	}
	for {
		select {
		case <-time.After(1 * time.Second):
			conn, err := net.Dial("tcp", address)
			if err == nil {
				conn.Close()
				return err
			}
			if seconds > timeOut {
				return errors.New("Time out")
			}
			seconds++
			if fn != nil && !fn() {
				return nil
			}
		case <-time.After(5 * time.Second):
			fmt.Println("== Waiting for " + address)
			if seconds > timeOut {
				return errors.New("Time out")
			}
			seconds += 5
			if fn != nil && !fn() {
				return
			}
		case <-time.After(time.Duration(timeOut) * time.Second):
			return errors.New("Time out")
		}
	}
	return
}

func (c *CLIConfig) CaddyStopHistory() (err error) {
	if DefaultConfig.Caddy.PidFile == `` {
		return
	}
	b, err := ioutil.ReadFile(DefaultConfig.Caddy.PidFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(b)))
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	procs, err := os.FindProcess(pid)
	if err == nil {
		return procs.Kill()
	}
	log.Error(err.Error())
	return nil
}

func (c *CLIConfig) CaddyStart() (err error) {
	err = c.CaddyStopHistory()
	if err != nil {
		log.Error(err.Error())
	}
	params := []string{`-c`, c.Conf, `-t`, `server`}
	c.cmd = exec.Command(os.Args[0], params...)
	c.cmd.Stdout = os.Stdout
	//cmd.Stderr = StderrCapturer{this}

	var hasError bool
	go func() {
		err := c.cmd.Run()
		if err != nil {
			log.Error(`Caddy Run Error: `, err)
			hasError = true
		}
	}()
	//err = dialAddress("127.0.0.1:80", 60, func() bool {return !hasError})
	return
}

func (c *CLIConfig) CaddyStop() error {
	if c.cmd == nil {
		return nil
	}
	if c.cmd.ProcessState != nil {
		return nil
	}
	return c.cmd.Process.Kill()
}

func (c *CLIConfig) CaddyRestart() error {
	err := c.CaddyStop()
	if err != nil {
		return err
	}
	err = c.CaddyStart()
	return err
}

type Config struct {
	DB struct {
		Type     string            `json:"type"`
		User     string            `json:"user"`
		Password string            `json:"password"`
		Host     string            `json:"host"`
		Database string            `json:"database"`
		Prefix   string            `json:"prefix"`
		Options  map[string]string `json:"options"`
		Debug    bool              `json:"debug"`
	} `json:"db"`

	Log struct {
		Debug        bool   `json:"debug"`
		Colorable    bool   `json:"colorable"`    // for console
		SaveFile     string `json:"saveFile"`     // for file
		FileMaxBytes int64  `json:"fileMaxBytes"` // for file
		Targets      string `json:"targets"`
	} `json:"log"`

	Sys struct {
		VhostsfileDir string            `json:"vhostsfileDir"`
		AllowIP       []string          `json:"allowIP"`
		Accounts      map[string]string `json:"accounts"`
		SSLHosts      []string          `json:"sslHosts"`
		SSLCacheFile  string            `json:"sslCacheFile"`
		SSLKeyFile    string            `json:"sslKeyFile"`
		SSLCertFile   string            `json:"sslCertFile"`
		Debug         bool              `json:"debug"`
	} `json:"sys"`

	Cookie struct {
		Domain   string `json:"domain"`
		MaxAge   int    `json:"maxAge"`
		Path     string `json:"path"`
		HttpOnly bool   `json:"httpOnly"`
		HashKey  string `json:"hashKey"`
		BlockKey string `json:"blockKey"`
	} `json:"cookie"`

	Caddy caddy.Config `json:"caddy"`
}

var (
	DefaultConfig         = &Config{}
	DefaultCLIConfig      = &CLIConfig{}
	ErrUnknowDatabaseType = errors.New(`unkown database type`)
)

func ParseConfig() error {
	if false {
		b, err := confl.Marshal(DefaultConfig)
		err = ioutil.WriteFile(DefaultCLIConfig.Conf, b, os.ModePerm)
		if err != nil {
			return err
		}
	}
	_, err := confl.DecodeFile(DefaultCLIConfig.Conf, DefaultConfig)
	if err != nil {
		return err
	}
	confDir := filepath.Dir(DefaultCLIConfig.Conf)
	if len(DefaultConfig.Sys.SSLCacheFile) == 0 {
		DefaultConfig.Sys.SSLCacheFile = filepath.Join(confDir, `letsencrypt.cache`)
	}
	if len(DefaultConfig.Caddy.Caddyfile) == 0 {
		DefaultConfig.Caddy.Caddyfile = `./Caddyfile`
	} else if strings.HasSuffix(DefaultConfig.Caddy.Caddyfile, `/`) || strings.HasSuffix(DefaultConfig.Caddy.Caddyfile, `\`) {
		DefaultConfig.Caddy.Caddyfile = filepath.Join(DefaultConfig.Caddy.Caddyfile, `Caddyfile`)
	}
	if len(DefaultConfig.Sys.VhostsfileDir) == 0 {
		DefaultConfig.Sys.VhostsfileDir = filepath.Join(confDir, `vhosts`)
	}
	caddy.Fixed(&DefaultConfig.Caddy)
	InitLog()
	InitSessionOptions()
	if DefaultConfig.Sys.Debug {
		log.SetLevel(`Debug`)
	} else {
		log.SetLevel(`Info`)
	}
	return ConnectDB()
}

func ConnectDB() error {
	factory.CloseAll()
	switch DefaultConfig.DB.Type {
	case `mysql`:
		settings := mysql.ConnectionURL{
			Host:     DefaultConfig.DB.Host,
			Database: DefaultConfig.DB.Database,
			User:     DefaultConfig.DB.User,
			Password: DefaultConfig.DB.Password,
			Options:  DefaultConfig.DB.Options,
		}
		database, err := mysql.Open(settings)
		if err != nil {
			return err
		}
		factory.SetDebug(DefaultConfig.DB.Debug)
		cluster := factory.NewCluster().AddW(database)
		factory.SetCluster(0, cluster).Cluster(0).SetPrefix(DefaultConfig.DB.Prefix)
	case `mongo`:
		settings := mongo.ConnectionURL{
			Host:     DefaultConfig.DB.Host,
			Database: DefaultConfig.DB.Database,
			User:     DefaultConfig.DB.User,
			Password: DefaultConfig.DB.Password,
			Options:  DefaultConfig.DB.Options,
		}
		database, err := mongo.Open(settings)
		if err != nil {
			return err
		}
		factory.SetDebug(DefaultConfig.DB.Debug)
		cluster := factory.NewCluster().AddW(database)
		factory.SetCluster(0, cluster).Cluster(0).SetPrefix(DefaultConfig.DB.Prefix)
	default:
		return ErrUnknowDatabaseType
	}
	return nil
}

func InitLog() {

	//======================================================
	// 配置日志
	//======================================================
	if DefaultConfig.Log.Debug {
		log.DefaultLog.MaxLevel = log.LevelDebug
	} else {
		log.DefaultLog.MaxLevel = log.LevelInfo
	}
	targets := []log.Target{}

	for _, targetName := range strings.Split(DefaultConfig.Log.Targets, `,`) {
		switch targetName {
		case "console":
			//输出到命令行
			consoleTarget := log.NewConsoleTarget()
			consoleTarget.ColorMode = DefaultConfig.Log.Colorable
			targets = append(targets, consoleTarget)

		case "file":
			//输出到文件
			fileTarget := log.NewFileTarget()
			fileTarget.FileName = DefaultConfig.Log.SaveFile
			fileTarget.Filter.MaxLevel = log.LevelInfo
			fileTarget.MaxBytes = DefaultConfig.Log.FileMaxBytes
			targets = append(targets, fileTarget)
		}
	}

	log.SetTarget(targets...)
	log.SetFatalAction(log.ActionExit)
}

func MustOK(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
