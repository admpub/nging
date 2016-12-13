package config

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/admpub/caddyui/application/library/caddy"
	"github.com/admpub/caddyui/application/library/ftp"
	"github.com/admpub/confl"
	"github.com/admpub/log"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/mongo"
	"github.com/webx-top/db/mysql"
	"github.com/webx-top/echo/middleware/bytes"
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
	if DefaultConfig.Sys.EditableFileMaxBytes < 1 && len(DefaultConfig.Sys.EditableFileMaxSize) > 0 {
		DefaultConfig.Sys.EditableFileMaxBytes, err = bytes.Parse(DefaultConfig.Sys.EditableFileMaxSize)
		if err != nil {
			log.Error(err.Error())
		}
	}
	caddy.Fixed(&DefaultConfig.Caddy)
	ftp.Fixed(&DefaultConfig.FTP)
	InitLog()
	InitSessionOptions()
	if DefaultConfig.Sys.Debug {
		log.SetLevel(`Debug`)
	} else {
		log.SetLevel(`Info`)
	}

	err = ConnectDB()
	if err != nil {
		return err
	}

	err = DefaultCLIConfig.Reload()
	return err
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

func CmdIsRunning(cmd *exec.Cmd) bool {
	return cmd != nil && cmd.ProcessState == nil
}
