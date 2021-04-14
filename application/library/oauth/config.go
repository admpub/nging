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

package oauth

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/handler/oauth2"
)

func NewConfig() *Config {
	return &Config{
		Extra: map[string]interface{}{
			`iconImage`: ``, //图标图片路径
			`iconClass`: ``, //图标类标识
			`title`:     ``, //显示标题
		},
	}
}

type Config struct {
	On     bool                   `json:"-"`      //开关
	Key    string                 `json:"key"`    //App ID
	Secret string                 `json:"secret"` //Secret Key
	Name   string                 `json:"-"`      //标识，如：github,wechat,alipay等
	Extra  map[string]interface{} `json:"extra"`  //其它扩展数据
}

func (c *Config) ToAccount(name string) *oauth2.Account {
	c.Name = name
	a := &oauth2.Account{
		On:     c.On,
		Name:   c.Name,
		Secret: c.Secret,
		Key:    c.Key,
		Extra:  c.Extra,
	}
	if c.Extra == nil {
		a.Extra = map[string]interface{}{}
	}
	return a
}

func (c *Config) FromStore(name string, v echo.Store) *Config {
	c.Name = name
	c.Key = v.String(`key`)
	c.Secret = v.String(`secret`)
	c.Extra = v.GetStore(`extra`)
	return c
}
