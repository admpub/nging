package service

import (
	stdLog "log"
	"os"
	"path/filepath"

	"github.com/admpub/log"
	"github.com/webx-top/com"
)

func Run(action string) error {
	conf := &Config{}
	conf.Name = `Nging`
	conf.DisplayName = `Nging`
	conf.Description = `Nging Server Manager`
	conf.Dir = com.SelfDir()
	conf.Exec = os.Args[0]
	if len(os.Args) > 2 {
		conf.Args = os.Args[2:]
	}
	logDir := filepath.Join(com.SelfDir(), `logs`)
	if info, err := os.Stat(logDir); err != nil || !info.IsDir() {
		err = os.Mkdir(logDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	fileTarget := log.NewFileTarget()
	fileTarget.FileName = filepath.Join(logDir, `app_{date:20060102}.log`) //按天分割日志
	fileTarget.MaxBytes = 10 * 1024 * 1024
	log.SetTarget(fileTarget)
	conf.Stderr = log.Writer(log.LevelError)
	conf.Stdout = log.Writer(log.LevelInfo)

	w, err := FileWriter(filepath.Join(logDir, `service.log`))
	if err != nil {
		return err
	}
	defer w.Close()
	stdLog.SetOutput(w)
	return New(conf, action)
}
