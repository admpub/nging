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

package service

import (
	"io"
	stdLog "log"
	"os"
	"path/filepath"

	"github.com/webx-top/com"
	"github.com/admpub/log"
)

func Run(options *Options, action string) error {
	conf := &Config{}
	conf.CopyFromOptions(options)
	conf.Dir = com.SelfDir()
	conf.Exec = os.Args[0]
	if len(os.Args) > 3 {
		conf.Args = os.Args[3:]
	}
	if err := initServiceLog(conf); err != nil {
		return err
	}
	return New(conf, action)
}

func FileWriter(file string) (io.WriteCloser, error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	return f, err
}

func initServiceLog(conf *Config) error {
	logDir := filepath.Join(com.SelfDir(), `data`, `logs`)
	err := com.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		return err
	}
	// 保存子进程在控制台输出的日志
	serviceLog := log.New()
	serviceLog.SetFormatter(log.EmptyFormatter)
	fileTarget := log.NewFileTarget()
	fileTarget.FileName = filepath.Join(logDir, `service_app_{date:20060102}.log`) //按天分割日志
	fileTarget.MaxBytes = 100 * 1024 * 1024
	serviceLog.SetTarget(fileTarget)
	conf.Stderr = serviceLog.Writer(log.LevelError)
	conf.Stdout = serviceLog.Writer(log.LevelInfo)

	// 保存service程序中输出的日志
	w, err := FileWriter(filepath.Join(logDir, `service.log`))
	if err != nil {
		return err
	}
	conf.OnExited = func() error {
		serviceLog.Close()
		if w != nil {
			return w.Close()
		}
		return nil
	}
	stdLog.SetOutput(w)
	stdLog.SetFlags(stdLog.Lshortfile)
	conf.logger = newLogger(w)
	return err
}
