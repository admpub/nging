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
	"strings"
	"sync/atomic"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/config/startup"
	"github.com/admpub/nging/v5/application/library/errorslice"
	"github.com/admpub/nging/v5/application/library/route"
	"github.com/nging-plugins/firewallmanager/application/library/cmder"
	"github.com/nging-plugins/firewallmanager/application/library/driver"
	"github.com/nging-plugins/firewallmanager/application/library/enums"
	"github.com/nging-plugins/firewallmanager/application/library/firewall"
	"github.com/nging-plugins/firewallmanager/application/model"
)

func RegisterRoute(r *route.Collection) {
	r.Backend.RegisterToGroup(`/firewall`, registerRoute)
}

var routeRegisters route.Registers

func registerRoute(g echo.RouteRegister) {
	ruleG := g.Group(`/rule`)
	ruleG.Route(`GET,POST`, `/static`, ruleStaticIndex)
	ruleG.Route(`GET,POST`, `/static_add`, ruleStaticAdd)
	ruleG.Route(`GET,POST`, `/static_edit`, ruleStaticEdit)
	ruleG.Route(`GET,POST`, `/static_delete`, ruleStaticDelete)
	ruleG.Route(`GET,POST`, `/static_apply`, ruleStaticApply)
	ruleG.Route(`GET,POST`, `/dynamic`, ruleDynamicIndex)
	ruleG.Route(`GET,POST`, `/dynamic_add`, ruleDynamicAdd)
	ruleG.Route(`GET,POST`, `/dynamic_edit`, ruleDynamicEdit)
	ruleG.Route(`GET,POST`, `/dynamic_delete`, ruleDynamicDelete)

	serviceG := g.Group(`/service`)
	serviceG.Route(`GET,POST`, `/restart`, Restart)
	serviceG.Route(`GET,POST`, `/stop`, Stop)
	serviceG.Route(`GET,POST`, `/log`, Log)

	routeRegisters.Apply(g)
}

var staticRuleLastModifyTs uint64

func setStaticRuleLastModifyTime(t time.Time) {
	atomic.StoreUint64(&staticRuleLastModifyTs, uint64(t.Unix()))
}

func getStaticRuleLastModifyTs() uint64 {
	return atomic.LoadUint64(&staticRuleLastModifyTs)
}
func init() {
	startup.OnAfter(`web.installed`, func() {
		firewall.Clear(`all`)
		ctx := defaults.NewMockContext()
		err := applyNgingRule(ctx)
		if err == nil {
			err = applyStaticRule(ctx)
		}
		if err != nil {
			log.Error(err)
		}
	})
	startup.OnAfter(`web`, func() {
	})
}

// applyNgingRule 添加 Nging 自己的规则。主要用于避免用户设置静态规则不当，拦截了 Nging 自身
func applyNgingRule(ctx echo.Context) error {
	// 放行 Nging 自己的端口
	portStr := param.AsString(config.FromCLI().Port)
	rule := driver.Rule{
		Name:      `NgingPreset`,
		Protocol:  enums.ProtocolTCP,
		Type:      enums.TableFilter,
		Direction: enums.ChainInput,
		Action:    enums.TargetAccept,
		LocalPort: portStr,
		IPVersion: `4`,
	}
	cfg := cmder.GetFirewallConfig()
	if cfg.NgingRule != nil {
		if len(cfg.NgingRule.IPWhitelist) > 0 {
			rule.RemoteIP = cfg.NgingRule.IPWhitelist
		}
		if len(cfg.NgingRule.OtherPort) > 0 {
			otherPortStrs := cfg.NgingRule.OtherPortStrs()
			if len(otherPortStrs) > 0 {
				portStrs := []string{portStr}
				portStrs = append(portStrs, otherPortStrs...)
				portStrs = param.StringSlice(portStrs).Filter().Unique().String()
				rule.LocalPort = strings.Join(portStrs, `,`)
			}
		}
	}
	ngingRules := []driver.Rule{rule}
	if cfg.NgingRule != nil && cfg.NgingRule.RpsLimit > 0 {
		// 限流 Nging 自己的端口
		rateLimitRule := rule
		rateLimitRule.Name = `NgingPresetLimit`
		if cfg.NgingRule.RpsLimit > 0 {
			rateLimitRule.RateLimit = param.AsString(cfg.NgingRule.RpsLimit) + `+/p/s`
		}
		if cfg.NgingRule.RateBurst > 0 {
			rateLimitRule.RateBurst = cfg.NgingRule.RateBurst
		} else {
			rateLimitRule.RateBurst = cfg.NgingRule.RpsLimit
		}
		if cfg.NgingRule.RateExpires > 0 {
			rateLimitRule.RateExpires = cfg.NgingRule.RateExpires
		}
		if rateLimitRule.RateExpires == 0 {
			rateLimitRule.RateExpires = 86400
		}
		rateLimitRule.Action = enums.TargetDrop
		ngingRules = append(ngingRules, rateLimitRule)
	}
	err := firewall.Append(ngingRules...)
	if err != nil {
		err = ctx.NewError(code.Failure, `[firewallManager] failed to applyNgingRule: %v`, err)
	}
	return err
}

// applyStaticRule 添加用户定义的静态规则
func applyStaticRule(ctx echo.Context) error {
	errs := errorslice.New()
	m := model.NewRuleStatic(ctx)
	_, err := m.ListByOffset(nil, nil, 0, -1, `disabled`, `Y`)
	if err == nil {
		rows := m.Objects()
		deleteRules := make([]driver.Rule, len(rows))
		for idx, row := range rows {
			rule := m.AsRule(row)
			deleteRules[idx] = rule
		}
		if len(deleteRules) > 0 {
			err = firewall.Delete(deleteRules...)
			if err != nil {
				errs.Add(err)
			} else {
				setStaticRuleLastModifyTime(time.Now())
			}
		}
	}
	_, err = m.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`position`, `id`)
	}, 0, -1, `disabled`, `N`)
	if err == nil {
		rows := m.Objects()
		createRules := make([]driver.Rule, len(rows))
		for idx, row := range rows {
			rule := m.AsRule(row)
			createRules[idx] = rule
		}
		if len(createRules) > 0 {
			err = firewall.Append(createRules...)
			if err != nil {
				errs.Add(err)
			} else {
				setStaticRuleLastModifyTime(time.Now())
			}
		}
	}
	if err == nil {
		err = errs.ToError()
	}
	return err
}
