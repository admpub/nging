//go:build linux

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

package firewall

import (
	"github.com/webx-top/echo"
)

var DynamicRuleBackends = echo.NewKVData().
	Add(`ipset`, `ipset`).
	Add(`nft`, `nftables`)

var dynamicRuleSourceFormElements = map[string]echo.KVList{
	`file`: {
		echo.NewKV(`path`, `文件路径`).SetHKV(`inputType`, `text`).SetHKV(`required`, true),
	},
	`systemd`: {
		echo.NewKV(`service`, `服务名`).SetHKV(`inputType`, `text`).SetHKV(`required`, true),
	},
	`process`: {
		echo.NewKV(`name`, `可执行文件`).SetHKV(`inputType`, `text`).SetHKV(`required`, true),
		echo.NewKV(`args`, `命令参数`).SetHKV(`inputType`, `text`),
	},
}

var DynamicRuleSources = echo.NewKVData().
	Add(`file`, `日志文件(使用tail)`, echo.KVOptHKV(`formElements`, dynamicRuleSourceFormElements[`file`])).
	Add(`systemd`, `服务日志(使用journalctl)`, echo.KVOptHKV(`formElements`, dynamicRuleSourceFormElements[`systemd`])).
	Add(`kernel`, `系统日志(使用journalctl)`).
	Add(`process`, `命令`, echo.KVOptHKV(`formElements`, dynamicRuleSourceFormElements[`process`]))

var dynamicRuleActionFormElements = map[string]echo.KVList{
	`ban`: {
		echo.NewKV(`duration`, `时长`).SetHKV(`inputType`, `text`).SetHKV(`required`, true).SetHKV(`pattern`, `^([\d]+[smh])*$`).SetHKV(`helpBlock`, `指定时长，由数字和单位字母(h-时/m-分/s-秒)组成`).SetHKV(`placeholder`, `例如: 24h(24小时)`),
	},
	`log`: {
		echo.NewKV(`extended`, `日志类型`).SetHKV(`inputType`, `select`).SetHKV(`options`, echo.KVList{
			echo.NewKV(`simple`, `基础版日志(记录时间和IP)`),
			echo.NewKV(`extended`, `扩展版日志(基础版日志+匹配到的行和规则)`),
		}).SetHKV(`required`, true),
	},
}

var DynamicRuleActions = echo.NewKVData().
	Add(`ban`, `禁止访问`, echo.KVOptHKV(`formElements`, dynamicRuleActionFormElements[`ban`])).
	Add(`log`, `记录日志`, echo.KVOptHKV(`formElements`, dynamicRuleActionFormElements[`log`]))
