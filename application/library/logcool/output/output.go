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

package output

import (
	"github.com/admpub/logcool/utils"
	"github.com/admpub/nging/application/model"
)

const (
	ModuleName = "nging"
)

// OutputConfig Define nging' config.
type OutputConfig struct {
	utils.OutputConfig
	event chan utils.LogEvent
}

func init() {
	utils.RegistOutputHandler(ModuleName, InitHandler)
}

// InitHandler Init outputstdout Handler.
func InitHandler(confraw *utils.ConfigRaw) (retconf utils.TypeOutputConfig, err error) {
	conf := OutputConfig{
		OutputConfig: utils.OutputConfig{
			CommonConfig: utils.CommonConfig{
				Type: ModuleName,
			},
		},
		event: make(chan utils.LogEvent),
	}
	if err = utils.ReflectConfig(confraw, &conf); err != nil {
		return
	}

	go conf.loopEvent()

	retconf = &conf
	return
}

// Event Input's event,and this is the main function of output.
func (oc *OutputConfig) Event(event utils.LogEvent) (err error) {
	oc.event <- event
	return
}

func (oc *OutputConfig) loopEvent() (err error) {
	for {
		event := <-oc.event
		oc.sendEvent(event)
	}
}

func (oc *OutputConfig) sendEvent(event utils.LogEvent) (err error) {
	data := model.NewAccessLog(nil)
	err = data.Parse(event.Message)
	if err != nil {
		return
	}

	_, err = data.Add()
	return
}
