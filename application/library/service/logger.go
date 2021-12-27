/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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
	"log"
	"os"
)

func newLogger(writer io.Writer) *consoleLogger {
	if writer == nil {
		writer = os.Stderr
	}
	logger := &consoleLogger{
		info: log.New(writer, "[I] ", log.LstdFlags),
		warn: log.New(writer, "[W] ", log.LstdFlags),
		err:  log.New(writer, "[E] ", log.LstdFlags),
	}
	return logger
}

type consoleLogger struct {
	info, warn, err *log.Logger
}

func (c consoleLogger) Error(v ...interface{}) error {
	c.err.Print(v...)
	return nil
}
func (c consoleLogger) Warning(v ...interface{}) error {
	c.warn.Print(v...)
	return nil
}
func (c consoleLogger) Info(v ...interface{}) error {
	c.info.Print(v...)
	return nil
}
func (c consoleLogger) Errorf(format string, a ...interface{}) error {
	c.err.Printf(format, a...)
	return nil
}
func (c consoleLogger) Warningf(format string, a ...interface{}) error {
	c.warn.Printf(format, a...)
	return nil
}
func (c consoleLogger) Infof(format string, a ...interface{}) error {
	c.info.Printf(format, a...)
	return nil
}
