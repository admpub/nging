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

package filter

import (
	"github.com/admpub/logcool/utils"
	"github.com/admpub/nging/v3/application/model"
)

const (
	ModuleName = "nging"
)

// FilterConfig define nging' config.
type FilterConfig struct {
	utils.FilterConfig
}

func init() {
	utils.RegistFilterHandler(ModuleName, InitHandler)
}

// InitHandler Init nging Handler.
func InitHandler(confraw *utils.ConfigRaw) (tfc utils.TypeFilterConfig, err error) {
	conf := FilterConfig{
		FilterConfig: utils.FilterConfig{
			CommonConfig: utils.CommonConfig{
				Type: ModuleName,
			},
		},
	}
	// Reflect config from configraw.
	if err = utils.ReflectConfig(confraw, &conf); err != nil {
		return
	}

	tfc = &conf
	return
}

// Event Filter's event,and this is the main function of filter.
func (fc *FilterConfig) Event(event utils.LogEvent) utils.LogEvent {
	data := model.NewAccessLog(nil)
	err := data.Parse(event.Message)
	if err != nil {
		return event
	}
	event.Extra = data.ToMap()
	return event
}
