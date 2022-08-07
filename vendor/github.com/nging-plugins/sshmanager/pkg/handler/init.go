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

package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	termConfig "github.com/admpub/web-terminal/config"
	termHandler "github.com/admpub/web-terminal/handler"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	ngingdbschema "github.com/admpub/nging/v4/application/dbschema"
	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/handler/user"
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/registry/perm"

	"github.com/admpub/nging/v4/application/library/route"
	"github.com/nging-plugins/sshmanager/pkg/dbschema"
	"github.com/nging-plugins/sshmanager/pkg/model"
)

func RegisterRoute(r *route.Collection) {
	r.Backend.RegisterToGroup(`/term`, registerRoute)
	user.OnAutoCompletePath(AutoCompletePath)
	perm.AuthRegister(`/term/client/replay`, authTermClient)
	perm.AuthRegister(`/term/client/ssh`, authTermClient)
	perm.AuthRegister(`/term/client/telnet`, authTermClient)
	perm.AuthRegister(`/term/client/cmd`, authTermClient)
	perm.AuthRegister(`/term/client/cmd2`, authTermClient)
	perm.AuthRegister(`/term/client/ssh_exec`, authTermClient)
}

type TerminalParam struct {
	Query url.Values
	User  *dbschema.NgingSshUser
}

type key string

var contextKey key = `param`

func registerRoute(g echo.RouteRegister) {
	g.Route(`GET`, `/account`, AccountIndex)
	g.Route(`GET,POST`, `/account_add`, AccountAdd)
	g.Route(`GET,POST`, `/account_edit`, AccountEdit)
	g.Route(`GET,POST`, `/account_delete`, AccountDelete)
	g.Route(`GET`, `/client`, Client)
	g.Route(`GET,POST`, `/sftp`, Sftp)
	termHandler.ParamGet = func(ctx *termHandler.Context, name string) (value string) {
		/*
			defer func() {
				fmt.Println(`web-terminal: [param]`, name, `->`, value)
			}()
		// */
		var (
			param *TerminalParam
			val   interface{}
			ok    bool
		)
		if val, ok = ctx.Data.Load(contextKey); ok {
			param, ok = val.(*TerminalParam)
		}
		if !ok {
			param = &TerminalParam{
				Query: ctx.Request().URL.Query(),
			}

			id := param.Query.Get(`id`)
			if len(id) > 0 {
				m := model.NewSshUser(nil)
				err := m.Get(nil, `id`, id)
				if err == nil {
					param.User = m.NgingSshUser
				}
			}
			ctx.Data.Store(contextKey, param)
		}
		if param.User != nil {
			switch name {
			case `password`:
				return config.FromFile().Decode(param.User.Password)
			case `user`:
				return param.User.Username
			case `protocol`:
				return param.User.Protocol
			case `hostname`:
				return param.User.Host
			case `port`:
				return fmt.Sprint(param.User.Port)
			case `privateKey`:
				return param.User.PrivateKey
			case `passphrase`:
				return config.FromFile().Decode(param.User.Passphrase)
			case `charset`:
				if len(param.User.Charset) == 0 {
					return `UTF-8`
				}
				return param.User.Charset
			}
		}
		value = param.Query.Get(name)
		if name == `password` {
			value = config.FromFile().Decode(value)
		}
		return value
	}
	termConfig.Default.APPRoot = handler.BackendPrefix + `/client/`
	termConfig.Default.Debug = config.FromFile().Debug
	logDir := filepath.Join(echo.Wd(), `data/logs`)
	err := com.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	termConfig.Default.LogDir = filepath.Join(logDir, `term`)
	termConfig.Default.ResourceDir = `public/xterm`
	termConfig.Default.MIBSDir = filepath.Join(echo.Wd(), `data/mibs`)
	err = com.MkdirAll(termConfig.Default.MIBSDir, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	termConfig.Default.SetDefault()
	termHandler.Register(termConfig.Default.APPRoot, func(path string, h http.Handler) {
		g.Any(path, h)
	})
	g.Route(`GET`, `/group`, GroupIndex)
	g.Route(`GET,POST`, `/group_add`, GroupAdd)
	g.Route(`GET,POST`, `/group_edit`, GroupEdit)
	g.Route(`GET,POST`, `/group_delete`, GroupDelete)
}

func authTermClient(
	h echo.Handler,
	c echo.Context,
	user *ngingdbschema.NgingUser,
	permission *perm.RolePermission,
) (ppath string, returning bool, err error) {
	ppath = `/term/client`
	return
}
