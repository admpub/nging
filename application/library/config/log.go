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

package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/log"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/service"
)

var DefaultLogCategories = []string{`db`, `echo,mock`}

type Log struct {
	Debug        bool     `json:"debug"`
	Colorable    bool     `json:"colorable"`    // for console
	SaveFile     string   `json:"saveFile"`     // for file
	FileMaxBytes int64    `json:"fileMaxBytes"` // for file
	Targets      string   `json:"targets" form_delimiter:","`
	Categories   []string `json:"cagegories"`
}

func (c *Log) LogFilename(category string) (string, error) {
	_, _, timeformat, filename, err := log.DateFormatFilename(c.LogFile())
	if err != nil {
		return ``, err
	}
	var logFile string
	if len(timeformat) > 0 {
		logFile = fmt.Sprintf(filename, time.Now().Format(timeformat))
	} else {
		logFile = filename
	}
	logFile = strings.Replace(logFile, `{category}`, category, -1)
	if com.FileExists(logFile) {
		return logFile, nil
	}
	serviceAppLogFile := service.ServiceLogDir() + echo.FilePathSeparator + service.ServiceAppLogFile
	_, _, timeformat, filename, err = log.DateFormatFilename(serviceAppLogFile)
	if err == nil {
		if len(timeformat) > 0 {
			serviceAppLogFile = fmt.Sprintf(filename, time.Now().Format(timeformat))
		} else {
			serviceAppLogFile = filename
		}
		serviceAppLogFile = strings.Replace(serviceAppLogFile, `{category}`, category, -1)
		if com.FileExists(serviceAppLogFile) {
			logFile = serviceAppLogFile
		}
	}
	return logFile, nil
}

func (c *Log) Show(ctx echo.Context) error {
	category := ctx.Param(`category`, log.DefaultLog.Category)
	if strings.Contains(category, `..`) {
		return ctx.JSON(ctx.Data().SetInfo(ctx.T(`参数错误: %s`, category), 0).SetZone(`category`))
	}
	if category != log.DefaultLog.Category && !log.HasCategory(category) {
		return ctx.JSON(ctx.Data().SetInfo(ctx.T(`不存在日志分类: %s`, category), 0).SetZone(`category`))
	}
	logFile, err := c.LogFilename(category)
	if err != nil {
		return ctx.JSON(ctx.Data().SetError(err))
	}
	return common.LogShow(ctx, logFile)
}

func (c *Log) SetBy(r echo.H, defaults echo.H) *Log {
	if !r.Has(`log`) && defaults != nil {
		r.Set(`log`, defaults.GetStore(`log`))
	}
	loge := r.GetStore(`log`)
	c.Colorable = loge.Bool(`colorable`)
	c.SaveFile = loge.String(`saveFile`)
	switch t := loge.Get(`targets`).(type) {
	case []interface{}:
		for k, v := range t {
			if k > 0 {
				c.Targets += `,`
			}
			c.Targets = fmt.Sprint(v)
		}
	case []string:
		c.Targets = strings.Join(t, `,`)
	case string:
		c.Targets = t
	}
	c.FileMaxBytes = loge.Int64(`fileMaxBytes`)
	return c
}

func (c *Log) LogFile() string {
	if len(c.SaveFile) > 0 {
		return c.SaveFile
	}
	return filepath.Join(echo.Wd(), `data/logs/{category}_{date:20060102}_info.log`)
}

func (c *Log) LogCategories() []string {
	if len(c.Categories) == 0 {
		return DefaultLogCategories
	}
	return c.Categories
}

func (c *Log) Init() {
	//======================================================
	// 配置日志
	//======================================================
	if c.Debug {
		log.DefaultLog.MaxLevel = log.LevelDebug
		//log.DefaultLog.Formatter = log.ShortFileFormatter
	} else {
		log.DefaultLog.MaxLevel = log.LevelInfo
	}
	targets := []log.Target{}
	for _, targetName := range strings.Split(c.Targets, `,`) {
		targetName = strings.TrimSpace(targetName)
		if len(targetName) == 0 {
			continue
		}
		switch targetName {
		case "file":
			//输出到文件
			fileTarget := log.NewFileTarget()
			fileTarget.Filter = &log.Filter{
				MaxLevel: log.DefaultLog.MaxLevel,
			}
			if c.FileMaxBytes > 0 {
				fileTarget.MaxBytes = c.FileMaxBytes
			}
			logFileName := c.LogFile()
			fileTarget.FileName = logFileName
			if strings.Contains(logFileName, `{category}`) {
				fileTarget.FileName = strings.Replace(logFileName, `{category}`, log.DefaultLog.Category, -1)
				fileTarget.Filter.Categories = []string{log.DefaultLog.Category, `websocket`, `watcher`}
				targets = append(targets, fileTarget)

				// subcategory
				for _, category := range c.LogCategories() {
					fileTarget := log.NewFileTarget()
					fileTarget.Filter = &log.Filter{
						MaxLevel:   log.DefaultLog.MaxLevel,
						Categories: strings.Split(category, `,`),
					}
					fileTarget.FileName = logFileName
					if c.FileMaxBytes > 0 {
						fileTarget.MaxBytes = c.FileMaxBytes
					}
					fileTarget.FileName = strings.Replace(logFileName, `{category}`, fileTarget.Filter.Categories[0], -1)
					targets = append(targets, fileTarget)
				}
			} else {
				targets = append(targets, fileTarget)
			}

		case "console":
			fallthrough
		default:
			//输出到命令行
			consoleTarget := log.NewConsoleTarget()
			consoleTarget.ColorMode = c.Colorable
			targets = append(targets, consoleTarget)
		}
	}

	log.SetTarget(targets...)
	log.SetFatalAction(log.ActionExit)
}
