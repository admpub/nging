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

package middleware

import (
	"fmt"
	"html/template"
	"strconv"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo/middleware/tplfunc"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/license"
	"github.com/admpub/nging/application/library/modal"
	"github.com/admpub/nging/application/model"
	"github.com/admpub/nging/application/registry/navigate"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/subdomains"
)

var DefaultAvatarURL = `/public/assets/backend/images/user_128.png`

func ErrorPageFunc(c echo.Context) error {
	c.SetFunc(`URLFor`, subdomains.Default.URL)
	c.SetFunc(`IsMessage`, common.IsMessage)
	c.SetFunc(`Stored`, c.Stored)
	c.SetFunc(`Languages`, func() []string {
		return config.DefaultConfig.Language.AllList
	})
	c.SetFunc(`IsError`, common.IsError)
	c.SetFunc(`IsOk`, common.IsOk)
	c.SetFunc(`Message`, common.Message)
	c.SetFunc(`Ok`, common.OkString)
	c.SetFunc(`Version`, func() *config.VersionInfo { return config.Version })
	c.SetFunc(`VersionNumber`, func() string { return config.Version.Number })
	c.SetFunc(`CommitID`, func() string { return config.Version.CommitID })
	c.SetFunc(`BuildTime`, func() string { return config.Version.BuildTime })
	c.SetFunc(`TrackerURL`, license.TrackerURL)
	c.SetFunc(`Fetch`, func(tmpl string, data interface{}) template.HTML {
		b, e := c.Fetch(tmpl, data)
		if e != nil {
			return template.HTML(e.Error())
		}
		return template.HTML(string(b))
	})
	c.SetFunc(`Prefix`, func() string {
		return c.Route().Prefix
	})
	c.SetFunc(`Path`, c.Path)
	c.SetFunc(`Queries`, c.Queries)
	configs := config.InDB()
	c.SetFunc(`Config`, func() echo.H {
		return configs
	})
	return nil
}

func FuncMap() echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			c.SetFunc(`Now`, time.Now)
			c.SetFunc(`UnixTime`, time.Now().Local().Unix)
			c.SetFunc(`HasString`, hasString)
			c.SetFunc(`Date`, date)
			c.SetFunc(`Token`, Token)
			c.SetFunc(`Modal`, func(data interface{}) template.HTML {
				return modal.Render(c, data)
			})
			ErrorPageFunc(c)
			c.SetFunc(`IndexStrSlice`, indexStrSlice)

			if !config.DefaultConfig.ConnectedDB(false) {
				return h.Handle(c)
			}
			//用户相关函数
			user, _ := c.Session().Get(`user`).(*dbschema.User)
			roleM := model.NewUserRole(c)
			var roleList []*dbschema.UserRole
			if user != nil {
				c.Set(`user`, user)
				c.SetFunc(`Username`, func() string { return user.Username })
				roleList = roleM.ListByUser(user)
				c.Set(`roleList`, roleList)
			}
			c.SetFunc(`Avatar`, func(avatar string, defaults ...string) string {
				if len(avatar) > 0 {
					return tplfunc.AddSuffix(avatar, `_200_200`)
				}
				if len(defaults) > 0 && len(defaults[0]) > 0 {
					return defaults[0]
				}
				return DefaultAvatarURL
			})
			c.SetFunc(`Project`, func(ident string) *navigate.ProjectItem {
				return navigate.ProjectGet(ident)
			})
			c.SetFunc(`Projects`, func() navigate.ProjectList {
				return navigate.ProjectListAll()
			})
			c.SetFunc(`Navigate`, func(side string) navigate.List {
				switch side {
				case `top`:
					if user != nil && user.Id == 1 {
						if navigate.TopNavigate == nil {
							return navigate.EmptyList
						}
						return *navigate.TopNavigate
					}
					return roleM.FilterNavigate(roleList, navigate.TopNavigate)
				case `left`:
					fallthrough
				default:
					if user != nil && user.Id == 1 {
						if navigate.LeftNavigate == nil {
							return navigate.EmptyList
						}
						return *navigate.LeftNavigate
					}
					return roleM.FilterNavigate(roleList, navigate.LeftNavigate)
				}
			})
			return h.Handle(c)
		})
	}
}

func indexStrSlice(slice []string, index int) string {
	if slice == nil {
		return ``
	}
	if index >= len(slice) {
		return ``
	}
	return slice[index]
}

func hasString(slice []string, str string) bool {
	if slice == nil {
		return false
	}
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func date(timestamp interface{}) time.Time {
	if v, y := timestamp.(int64); y {
		return time.Unix(v, 0)
	}
	if v, y := timestamp.(uint); y {
		return time.Unix(int64(v), 0)
	}
	v, _ := strconv.ParseInt(fmt.Sprint(timestamp), 10, 64)
	return time.Unix(v, 0)
}

func Token(values ...interface{}) string {
	urlValues := tplfunc.URLValues(values)
	return com.SafeBase64Encode(com.Token(config.DefaultConfig.APIKey, com.Str2bytes(urlValues.Encode())))
}
