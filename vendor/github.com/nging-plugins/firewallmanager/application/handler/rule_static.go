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
	"errors"
	"strings"
	"time"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/nging-plugins/firewallmanager/application/dbschema"
	"github.com/nging-plugins/firewallmanager/application/library/enums"
	"github.com/nging-plugins/firewallmanager/application/library/firewall"
	"github.com/nging-plugins/firewallmanager/application/model"
)

func ruleStaticSetFormData(c echo.Context) {
	c.Set(`types`, enums.Types.Slice())
	c.Set(`directions`, enums.Directions.Slice())
	c.Set(`ipProtocols`, enums.IPProtocols.Slice())
	c.Set(`netProtocols`, enums.NetProtocols.Slice())
	c.Set(`actions`, enums.Actions.Slice())
	c.Set(`stateList`, enums.StateList)
	c.Set(`tablesChains`, enums.TablesChains)
	c.Set(`chainParams`, enums.ChainParams)
}

func ruleStaticIndex(ctx echo.Context) error {
	m := model.NewRuleStatic(ctx)
	cond := db.NewCompounds()
	sorts := common.Sorts(ctx, m.NgingFirewallRuleStatic, `position`, `id`)
	list, err := m.ListPage(cond, sorts...)
	ctx.Set(`listData`, list)
	ctx.Set(`firewallBackend`, firewall.GetBackend())
	return ctx.Render(`firewall/rule/static`, common.Err(ctx, err))
}

func ruleStaticGetFirewallPosition(m *model.RuleStatic, row *dbschema.NgingFirewallRuleStatic, excludeOther ...uint) (uint, error) {
	next, err := m.NextRow(row.Type, row.Direction, row.IpVersion, row.Position, row.Id, excludeOther...)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			err = nil
		}
		return 0, err
	}
	var pos uint
	pos, err = firewall.FindPositionByID(row.IpVersion, row.Type, row.Direction, next.Id)
	if err != nil {
		return 0, err
	}
	if pos == 0 {
		excludeOther = append(excludeOther, row.Id)
		return ruleStaticGetFirewallPosition(m, next, excludeOther...)
	}
	return pos, err
}

func ruleStaticAdd(ctx echo.Context) error {
	m := model.NewRuleStatic(ctx)
	var err error
	if ctx.IsPost() {
		err = ctx.MustBind(m.NgingFirewallRuleStatic)
		if err != nil {
			goto END
		}
		m.State = strings.Join(param.StringSlice(ctx.FormValues(`state`)).Filter().String(), `,`)
		_, err = m.Add()
		if err != nil {
			goto END
		}
		rule := m.AsRule()
		rule.Number, err = ruleStaticGetFirewallPosition(m, m.NgingFirewallRuleStatic)
		if err != nil {
			goto END
		}
		if rule.Number > 0 {
			err = firewall.Insert(rule)
		} else {
			err = firewall.Append(rule)
		}
		if err != nil {
			goto END
		}
		setStaticRuleLastModifyTime(time.Now())
		return ctx.Redirect(handler.URLFor(`/firewall/rule/static`))
	} else {
		id := ctx.Formx(`copyId`).Uint()
		if id > 0 {
			err = m.Get(nil, db.Cond{`id`: id})
			if err == nil {
				echo.StructToForm(ctx, m.NgingFirewallRuleStatic, ``, echo.LowerCaseFirstLetter)
				ctx.Request().Form().Set(`id`, `0`)
			}
		}
	}

END:
	ctx.Set(`activeURL`, `/firewall/rule/static`)
	ctx.Set(`title`, ctx.T(`添加规则`))
	ctx.Set(`states`, param.StringSlice(strings.Split(m.State, `,`)).Filter().String())
	ruleStaticSetFormData(ctx)
	return ctx.Render(`firewall/rule/static_edit`, common.Err(ctx, err))
}

func ruleStaticEdit(ctx echo.Context) error {
	m := model.NewRuleStatic(ctx)
	id := ctx.Formx(`id`).Uint()
	err := m.Get(nil, `id`, id)
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		old := *m.NgingFirewallRuleStatic
		err = ctx.MustBind(m.NgingFirewallRuleStatic)
		if err != nil {
			goto END
		}
		m.State = strings.Join(param.StringSlice(ctx.FormValues(`state`)).Filter().String(), `,`)
		m.Id = id
		err = m.Edit(nil, `id`, id)
		if err != nil {
			goto END
		}
		oldRule := model.AsRule(&old)
		err = firewall.Delete(oldRule)
		if err != nil {
			goto END
		}
		setStaticRuleLastModifyTime(time.Now())
		if m.Disabled != `Y` {
			rule := m.AsRule()
			rule.Number, err = ruleStaticGetFirewallPosition(m, m.NgingFirewallRuleStatic)
			if err != nil {
				goto END
			}
			if rule.Number > 0 {
				err = firewall.Insert(rule)
			} else {
				err = firewall.Append(rule)
			}
			if err != nil {
				goto END
			}
		}
		return ctx.Redirect(handler.URLFor(`/firewall/rule/static`))
	} else if ctx.IsAjax() {
		disabled := ctx.Query(`disabled`)
		if len(disabled) > 0 {
			m.Disabled = disabled
			data := ctx.Data()
			err = m.UpdateField(nil, `disabled`, disabled, db.Cond{`id`: id})
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			rule := m.AsRule()
			if m.Disabled == `Y` {
				err = firewall.Delete(rule)
			} else {
				rule.Number, err = ruleStaticGetFirewallPosition(m, m.NgingFirewallRuleStatic)
				if err != nil {
					goto END
				}
				if rule.Number > 0 {
					err = firewall.Insert(rule)
				} else {
					err = firewall.Append(rule)
				}
			}
			if err != nil {
				data.SetError(err)
				return ctx.JSON(data)
			}
			data.SetInfo(ctx.T(`操作成功`))
			return ctx.JSON(data)
		}
	}
	echo.StructToForm(ctx, m.NgingFirewallRuleStatic, ``, echo.LowerCaseFirstLetter)

END:
	ctx.Set(`activeURL`, `/firewall/rule/static`)
	ctx.Set(`title`, ctx.T(`修改规则`))
	ctx.Set(`states`, param.StringSlice(strings.Split(m.State, `,`)).Filter().String())
	ruleStaticSetFormData(ctx)
	return ctx.Render(`firewall/rule/static_edit`, common.Err(ctx, err))
}

func ruleStaticDelete(ctx echo.Context) error {
	m := model.NewRuleStatic(ctx)
	id := ctx.Formx(`id`).Uint()
	err := m.Get(nil, `id`, id)
	if err == nil {
		err = m.Delete(nil, `id`, id)
		if err == nil {
			rule := m.AsRule()
			err = firewall.Delete(rule)
		}
	}
	if err == nil {
		setStaticRuleLastModifyTime(time.Now())
		handler.SendOk(ctx, ctx.T(`删除成功`))
	} else {
		handler.SendErr(ctx, err)
	}
	return ctx.Redirect(handler.URLFor(`/firewall/rule/static`))
}

func ruleStaticApply(ctx echo.Context) error {
	firewall.ResetEngine()
	firewall.Clear(`all`)
	err := applyNgingRule(ctx)
	if err == nil {
		err = applyStaticRule(ctx)
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`规则应用成功`))
	} else {
		handler.SendErr(ctx, err)
	}
	return ctx.Redirect(handler.URLFor(`/firewall/rule/static`))
}
