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

package middleware

import (
	"html/template"
	"net/url"
	"time"

	"github.com/admpub/timeago"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"
	"github.com/webx-top/echo/param"
	"github.com/webx-top/echo/subdomains"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/license"
	"github.com/admpub/nging/v5/application/library/modal"
	"github.com/admpub/nging/v5/application/library/perm"
	"github.com/admpub/nging/v5/application/library/role"
	"github.com/admpub/nging/v5/application/library/role/roleutils"
	uploadLibrary "github.com/admpub/nging/v5/application/library/upload"
	"github.com/admpub/nging/v5/application/registry/dashboard"
	"github.com/admpub/nging/v5/application/registry/navigate"
	"github.com/admpub/nging/v5/application/registry/settings"
	"github.com/admpub/nging/v5/application/registry/upload/checker"
)

var (
	DefaultAvatarURL = `/public/assets/backend/images/user_128.png`
	EmptyURL         = &url.URL{}
)

func init() {
	timeago.Set(`language`, `zh-cn`)
	tplfunc.TplFuncMap[`Languages`] = languages
	tplfunc.TplFuncMap[`URLFor`] = subdomains.Default.URL
	tplfunc.TplFuncMap[`URLByName`] = subdomains.Default.URLByName
	tplfunc.TplFuncMap[`IsMessage`] = common.IsMessage
	tplfunc.TplFuncMap[`IsError`] = common.IsError
	tplfunc.TplFuncMap[`IsOk`] = common.IsOk
	tplfunc.TplFuncMap[`Message`] = common.Message
	tplfunc.TplFuncMap[`Ok`] = common.OkString
	tplfunc.TplFuncMap[`Version`] = func() *config.VersionInfo { return config.Version }
	tplfunc.TplFuncMap[`VersionNumber`] = func() string { return config.Version.Number }
	tplfunc.TplFuncMap[`CommitID`] = func() string { return config.Version.CommitID }
	tplfunc.TplFuncMap[`BuildTime`] = func() string { return config.Version.BuildTime }
	tplfunc.TplFuncMap[`TrackerURL`] = license.TrackerURL
	tplfunc.TplFuncMap[`TrackerHTML`] = license.TrackerHTML
	tplfunc.TplFuncMap[`Config`] = getConfig
	tplfunc.TplFuncMap[`WithURLParams`] = common.WithURLParams
	tplfunc.TplFuncMap[`FullURL`] = common.FullURL
	tplfunc.TplFuncMap[`MaxRequestBodySize`] = func() int { return config.FromFile().GetMaxRequestBodySize() }
	tplfunc.TplFuncMap[`IndexStrSlice`] = indexStrSlice
	tplfunc.TplFuncMap[`HasString`] = hasString
	tplfunc.TplFuncMap[`Date`] = date
	tplfunc.TplFuncMap[`Token`] = checker.Token
	tplfunc.TplFuncMap[`BackendUploadURL`] = checker.BackendUploadURL
	tplfunc.TplFuncMap[`FrontendUploadURL`] = checker.FrontendUploadURL
	tplfunc.TplFuncMap[`Avatar`] = getAvatar
}

func languages() []string {
	return config.FromFile().Language.AllList
}

func getConfig(args ...string) echo.H {
	if len(args) > 0 {
		return config.Setting(args...)
	}
	return config.Setting()
}

func getAvatar(avatar string, defaults ...string) string {
	if len(avatar) > 0 {
		return tplfunc.AddSuffix(avatar, `_200_200`)
	}
	if len(defaults) > 0 && len(defaults[0]) > 0 {
		return defaults[0]
	}
	return DefaultAvatarURL
}

func ErrorPageFunc(c echo.Context) error {
	c.SetFunc(`Ext`, c.DefaultExtension)
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
	c.SetFunc(`Domain`, c.Domain)
	c.SetFunc(`Port`, c.Port)
	c.SetFunc(`Scheme`, c.Scheme)
	c.SetFunc(`Site`, c.Site)
	var siteURI *url.URL
	siteURL := config.Setting(`base`).String(`siteURL`)
	if len(siteURL) > 0 {
		siteURI, _ = url.Parse(siteURL)
	}
	c.Internal().Set(`siteURI`, siteURI)
	c.SetFunc(`SiteURI`, func() *url.URL {
		if siteURI == nil {
			return EmptyURL
		}
		return siteURI
	})
	c.SetFunc(`RequestURI`, c.RequestURI)
	c.SetFunc(`GetNextURL`, func(varNames ...string) string {
		return common.GetNextURL(c, varNames...)
	})
	c.SetFunc(`ReturnToCurrentURL`, func(varNames ...string) string {
		return common.ReturnToCurrentURL(c, varNames...)
	})
	c.SetFunc(`WithNextURL`, func(urlStr string, varNames ...string) string {
		return common.WithNextURL(c, urlStr, varNames...)
	})
	c.SetFunc(`CaptchaForm`, func(args ...interface{}) template.HTML {
		options := tplfunc.MakeMap(args)
		options.Set("captchaId", common.GetHistoryOrNewCaptchaID(c))
		return tplfunc.CaptchaFormWithURLPrefix(c.Echo().Prefix(), options)
	})
	c.SetFunc(`SQLQuery`, func() *common.SQLQuery {
		return common.NewSQLQuery(c)
	})
	c.SetFunc(`SQLQueryLimit`, func(offset int, limit int, linkID ...int) *common.SQLQuery {
		return common.NewSQLQueryLimit(c, offset, limit, linkID...)
	})
	c.SetFunc(`TimeAgo`, func(v interface{}, options ...string) string {
		if datetime, ok := v.(string); ok {
			return timeago.Take(datetime, c.Lang().Format(false, `-`))
		}
		var option string
		if len(options) > 0 {
			option = options[0]
		}
		return timeago.Timestamp(param.AsInt64(v), c.Lang().Format(false, `-`), option)
	})
	return nil
}

func FuncMap() echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			now := com.NewTime(time.Now())
			c.SetFunc(`Now`, func() *com.Time {
				return now
			})
			c.SetFunc(`UnixTime`, now.Local().Unix)
			c.SetFunc(`Modal`, func(data interface{}) template.HTML {
				return modal.Render(c, data)
			})
			ErrorPageFunc(c)
			if !config.FromFile().ConnectedDB(false) {
				return h.Handle(c)
			}
			c.SetFunc(`FileTypeByName`, uploadLibrary.FileTypeByName)
			var uploadCfg *uploadLibrary.Config
			c.SetFunc(`FileTypeIcon`, func(typ string) string {
				if uploadCfg == nil {
					uploadCfg = uploadLibrary.Get()
				}
				return uploadCfg.FileIcon(typ)
			})
			c.SetFunc(`Project`, func(ident string) *navigate.ProjectItem {
				return navigate.ProjectGet(ident)
			})
			c.SetFunc(`ProjectSearchIdent`, func(ident string) int {
				return navigate.ProjectSearchIdent(ident)
			})
			c.SetFunc(`Projects`, func() navigate.ProjectList {
				return navigate.ProjectListAll()
			})
			return h.Handle(c)
		})
	}
}

func BackendFuncMap() echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {

			//用户相关函数
			user := handler.User(c)
			if user != nil {
				c.Set(`user`, user)
				c.SetFunc(`Username`, func() string { return user.Username })
				c.Set(`roleList`, roleutils.UserRoles(c))
			}
			c.SetFunc(`ProjectIdent`, func() string {
				return GetProjectIdent(c)
			})
			c.SetFunc(`TopButtons`, func() dashboard.Buttons {
				buttons := dashboard.TopButtonAll(c)
				buttons.Ready(c)
				return buttons
			})
			c.SetFunc(`GlobalHeads`, func() dashboard.GlobalHeads {
				heads := dashboard.GlobalHeadAll(c)
				heads.Ready(c)
				return heads
			})
			c.SetFunc(`GlobalFooters`, func() dashboard.GlobalFooters {
				footers := dashboard.GlobalFooterAll(c)
				footers.Ready(c)
				return footers
			})
			c.SetFunc(`IsHiddenCard`, func(card *dashboard.Card) bool {
				return card.IsHidden(c)
			})
			c.SetFunc(`IsHiddenBlock`, func(block *dashboard.Block) bool {
				return block.IsHidden(c)
			})
			c.SetFunc(`IsValidPermHandler`, func(h perm.Handler) interface{} {
				return h.IsValid(c)
			})
			c.SetFunc(`SettingFormRender`, func(s *settings.SettingForm) interface{} {
				return s.Render(c)
			})
			c.SetFunc(`PermissionCheckByType`, func(permission *role.CommonPermission, typ string, permPath string) interface{} {
				return permission.CheckByType(c, typ, permPath)
			})
			c.SetFunc(`Navigate`, func(side string) navigate.List {
				return GetBackendNavigate(c, side)
			})
			return h.Handle(c)
		})
	}
}

func UserPermission(c echo.Context) *role.RolePermission {
	permission, ok := c.Internal().Get(`userPermission`).(*role.RolePermission)
	if !ok || permission == nil {
		permission = role.NewRolePermission().Init(roleutils.UserRoles(c))
		c.Internal().Set(`userPermission`, permission)
	}
	return permission
}

func GetProjectIdent(c echo.Context) string {
	projectIdent := c.Internal().String(`projectIdent`)
	if len(projectIdent) == 0 {
		projectIdent = navigate.ProjectIdent(c.Path())
		if len(projectIdent) == 0 {
			if proj := navigate.ProjectFirst(true); proj != nil {
				projectIdent = proj.Ident
			}
		}
		c.Internal().Set(`projectIdent`, projectIdent)
	}
	return projectIdent
}

func GetBackendNavigate(c echo.Context, side string) navigate.List {
	switch side {
	case `top`:
		navList, ok := c.Internal().Get(`navigate.top`).(navigate.List)
		if ok {
			return navList
		}
		user := handler.User(c)
		if user != nil && user.Id == 1 {
			if navigate.TopNavigate == nil {
				return navigate.EmptyList
			}
			return *navigate.TopNavigate
		}
		permission := UserPermission(c)
		navList = permission.FilterNavigate(c, navigate.TopNavigate)
		c.Internal().Set(`navigate.top`, navList)
		return navList
	case `left`:
		fallthrough
	default:
		navList, ok := c.Internal().Get(`navigate.left`).(navigate.List)
		if ok {
			return navList
		}
		user := handler.User(c)
		var leftNav *navigate.List
		ident := GetProjectIdent(c)
		if len(ident) > 0 {
			if proj := navigate.ProjectGet(ident); proj != nil {
				leftNav = proj.NavList
			}
		}
		if user != nil && user.Id == 1 {
			if leftNav == nil {
				return navigate.EmptyList
			}
			return *leftNav
		}
		permission := UserPermission(c)
		navList = permission.FilterNavigate(c, leftNav)
		c.Internal().Set(`navigate.left`, navList)
		return navList
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
	v := param.AsInt64(timestamp)
	return time.Unix(v, 0)
}
