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
package application

import (
	"path/filepath"
	"strings"

	"github.com/admpub/log"
	"github.com/webx-top/echo"

	. "github.com/admpub/caddyui/application/handler"
	"github.com/admpub/caddyui/application/library/config"
	"github.com/admpub/caddyui/application/library/monitor"
	"github.com/admpub/caddyui/application/middleware"
)

func Initialize(e *echo.Echo) {
	addRouter(e)
	config.MustOK(config.ParseConfig())
	me := monitor.MonitorEvent{
		Modify: func(file string) {
			if strings.HasSuffix(file, `.yaml`) {
				log.Info(`reload config from ` + file)
				config.MustOK(config.ParseConfig())
			}
		},
	}
	me.Watch(filepath.Dir(config.DefaultCLIConfig.Conf))
}

func addRouter(e *echo.Echo) {
	e.Get(`/`, Index)
	e.Get(`/login`, Login)

	g := e.Group(`/manage`, middleware.AuthCheck)
	{
		g.Get(``, ManageIndex)
	}

}
