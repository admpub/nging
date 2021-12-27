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
	"path/filepath"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/webx-top/com"
)

func WatchConfig(fn func(string) error) {
	me := com.MonitorEvent{
		Modify: func(file string) {
			if !strings.HasSuffix(file, `.yaml`) {
				return
			}
			log.Info(`Start reloading configuration file: ` + file)
			err := fn(file)
			if err == nil {
				log.Info(`Succcessfully reload the configuration file: ` + file)
				return
			}
			if err == common.ErrIgnoreConfigChange {
				log.Info(`No need to reload the configuration file: ` + file)
				return
			}
			log.Error(err)
		},
	}
	me.Watch()
	err := me.AddDir(filepath.Dir(DefaultCLIConfig.Conf))
	if err != nil {
		log.Error(err)
	}
}
