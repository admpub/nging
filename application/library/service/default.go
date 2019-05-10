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
	if len(os.Args) > 3 {
		conf.Args = os.Args[3:]
	}
	logDir := filepath.Join(com.SelfDir(), `data`, `logs`)
	if info, err := os.Stat(logDir); err != nil || !info.IsDir() {
		err = os.MkdirAll(logDir, os.ModePerm)
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
	stdLog.SetFlags(stdLog.LstdFlags)
	conf.logger = newLogger(w)
	return New(conf, action)
}

func FileWriter(file string) (io.WriteCloser, error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	return f, err
}
