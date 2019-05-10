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

package settings

import "github.com/webx-top/echo"

type SettingForm struct {
	Short    string                     //简短标签
	Label    string                     //标签文本
	Group    string                     //组标识
	Tmpl     []string                   //输入表单模板路径
	hookPost []func(echo.Context) error //数据提交逻辑处理
	hookGet  []func(echo.Context) error //数据读取逻辑处理
}

func (s *SettingForm) AddTmpl(tmpl string) *SettingForm {
	s.Tmpl = append(s.Tmpl, tmpl)
	return s
}

func (s *SettingForm) AddHookPost(hook func(echo.Context) error) *SettingForm {
	s.hookPost = append(s.hookPost, hook)
	return s
}

func (s *SettingForm) AddHookGet(hook func(echo.Context) error) *SettingForm {
	s.hookGet = append(s.hookGet, hook)
	return s
}

func (s *SettingForm) RunHookPost(ctx echo.Context) error {
	if s.hookPost == nil {
		return nil
	}
	for _, hook := range s.hookPost {
		err := hook(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SettingForm) RunHookGet(ctx echo.Context) error {
	if s.hookGet == nil {
		return nil
	}
	for _, hook := range s.hookGet {
		err := hook(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

var settings = []*SettingForm{
	&SettingForm{
		Short: `系统`,
		Label: `系统设置`,
		Group: `base`,
		Tmpl:  []string{`manager/settings/base`},
	},
	&SettingForm{
		Short: `SMTP`,
		Label: `SMTP服务器设置`,
		Group: `smtp`,
		Tmpl:  []string{`manager/settings/smtp`},
	},
	&SettingForm{
		Short: `日志`,
		Label: `日志设置`,
		Group: `log`,
		Tmpl:  []string{`manager/settings/log`},
	},
}

func Settings() []*SettingForm {
	return settings
}

func Register(sf ...*SettingForm) {
	settings = append(settings, sf...)
}

func Get(group string) (int, *SettingForm) {
	for index, setting := range settings {
		if setting.Group == group {
			return index, setting
		}
	}
	return -1, nil
}

func RunHookPost(ctx echo.Context) error {
	for _, setting := range settings {
		err := setting.RunHookPost(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func RunHookGet(ctx echo.Context) error {
	for _, setting := range settings {
		err := setting.RunHookGet(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
