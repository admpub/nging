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
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/admpub/gerberos"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/nging-plugins/firewallmanager/application/dbschema"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"
)

func RegisterDynamicRuleBackend(k string, v string) {
	DynamicRuleBackends.Add(k, v)
}

func RegisterDynamicRuleSource(k string, v string, formElements echo.KVList) {
	DynamicRuleSources.Add(k, v, echo.KVOptHKV(`formElements`, formElements))
}

func RegisterDynamicRuleAction(k string, v string, formElements echo.KVList) {
	DynamicRuleActions.Add(k, v, echo.KVOptHKV(`formElements`, formElements))
}

func setFromValues(c echo.Context, key string, values ...string) {
	for i, v := range values {
		if i == 0 {
			c.Request().Form().Set(key, v)
		} else {
			c.Request().Form().Add(key, v)
		}
	}
}

func SetDynamicRuleForm(c echo.Context, rule *dbschema.NgingFirewallRuleDynamic) error {
	var args []string
	if len(rule.SourceArgs) > 0 {
		err := json.Unmarshal([]byte(rule.SourceArgs), &args)
		if err == nil {
			setFromValues(c, `sourceArgs`, args...)
		}
	}

	args = []string{}
	if len(rule.Regexp) > 0 {
		err := json.Unmarshal([]byte(rule.Regexp), &args)
		if err == nil {
			setFromValues(c, `regexp`, strings.Join(args, "\n"))
		}
	}
	if len(rule.AggregateRegexp) > 0 {
		args = []string{}
		err := json.Unmarshal([]byte(rule.AggregateRegexp), &args)
		if err == nil {
			setFromValues(c, `aggregateRegexp`, strings.Join(args, "\n"))
		}
	}
	return nil
}

func DynamicRuleParseForm(c echo.Context, rule *dbschema.NgingFirewallRuleDynamic) error {
	rule.Name = c.Form(`name`)

	// source
	rule.SourceType = c.Form(`sourceType`)
	sourceArgs := c.FormxValues(`sourceArgs`).Filter()
	b, _ := json.Marshal(sourceArgs)
	rule.SourceArgs = string(b)

	// action
	rule.ActionType = c.Form(`actionType`)
	rule.ActionArg = c.Form(`actionArg`)

	// aggregate
	rule.AggregateDuration = c.Form(`aggregateDuration`)
	aggregateRegexp := c.Form(`aggregateRegexp`)
	var aggregateRegexpList []string
	for idx, re := range strings.Split(aggregateRegexp, "\n") {
		re = strings.TrimSpace(re)
		if len(re) == 0 {
			continue
		}
		if !strings.Contains(re, `%id%`) {
			return c.NewError(code.InvalidParameter, `必须在“聚合规则”的每一行规则里包含“%%id%%”，在第 %d 行规则中没有找到“%%id%%”，请添加`, idx+1).SetZone(`aggregateRegexp`)
		}
		if _, err := regexp.Compile(re); err != nil {
			return c.NewError(code.InvalidParameter, `“聚合规则”的第 %d 行规则有误：%v`, idx+1, err.Error()).SetZone(`regexp`)
		}
		aggregateRegexpList = append(aggregateRegexpList, re)
	}
	b, _ = json.Marshal(aggregateRegexpList)
	rule.AggregateRegexp = string(b)
	hasAggregate := len(rule.AggregateDuration) > 0 && len(aggregateRegexpList) > 0

	// occurrence
	rule.OccurrenceNum = c.Formx(`occurrenceNum`).Uint()
	rule.OccurrenceDuration = c.Form(`occurrenceDuration`)

	// regexp
	regexpRule := c.Form(`regexp`)
	var regexpList []string
	for idx, re := range strings.Split(regexpRule, "\n") {
		re = strings.TrimSpace(re)
		if len(re) == 0 {
			continue
		}
		if !strings.Contains(re, `%ip%`) {
			return c.NewError(code.InvalidParameter, `必须在“匹配规则”的每一行里包含“%%ip%%”，在第 %d 行规则中没有找到“%%ip%%”，请添加`, idx+1).SetZone(`regexp`)
		}
		if hasAggregate && !strings.Contains(re, `%id%`) {
			return c.NewError(code.InvalidParameter, `在设置“聚合规则”的情况下，必须同时在“匹配规则”的每一行规则里包含“%%id%%”，在第 %d 行规则中没有找到“%%id%%”，请添加`, idx+1).SetZone(`regexp`)
		}
		if _, err := regexp.Compile(re); err != nil {
			return c.NewError(code.InvalidParameter, `“匹配规则”的第 %d 行规则有误：%v`, idx+1, err.Error()).SetZone(`regexp`)
		}
		regexpList = append(regexpList, re)
	}
	b, _ = json.Marshal(regexpList)
	rule.Regexp = string(b)

	// status
	rule.Disabled = c.Form(`disabled`)
	return nil
}

func DynamicRuleFromDB(c echo.Context, row *dbschema.NgingFirewallRuleDynamic) (rule gerberos.Rule, err error) {
	var args []string
	if len(row.SourceArgs) > 0 {
		err = json.Unmarshal([]byte(row.SourceArgs), &args)
		if err != nil {
			err = common.JSONBytesParseError(err, []byte(row.SourceArgs))
			err = fmt.Errorf(`failed to parse SourceArgs: %w`, err)
			return
		}
	}
	rule.Source = []string{row.SourceType}
	rule.Source = append(rule.Source, args...)

	args = []string{}
	if len(row.Regexp) > 0 {
		err = json.Unmarshal([]byte(row.Regexp), &args)
		if err != nil {
			err = common.JSONBytesParseError(err, []byte(row.Regexp))
			err = fmt.Errorf(`failed to parse Regexp: %w`, err)
			return
		}
	}
	rule.Regexp = args
	rule.Action = []string{row.ActionType}
	if len(row.ActionArg) > 0 {
		rule.Action = append(rule.Action, row.ActionArg)
	}
	if len(row.AggregateDuration) > 0 && len(row.AggregateRegexp) > 0 {
		rule.Aggregate = []string{}
		rule.Aggregate = append(rule.Aggregate, row.AggregateDuration)
		args = []string{}
		err = json.Unmarshal([]byte(row.AggregateRegexp), &args)
		if err != nil {
			err = common.JSONBytesParseError(err, []byte(row.AggregateRegexp))
			err = fmt.Errorf(`failed to parse AggregateRegexp: %w`, err)
			return
		}
		rule.Aggregate = append(rule.Aggregate, args...)
	}
	if row.OccurrenceNum > 0 && len(row.OccurrenceDuration) > 0 {
		rule.Occurrences = []string{}
		rule.Occurrences = append(rule.Occurrences, param.AsString(row.OccurrenceNum))
		rule.Occurrences = append(rule.Occurrences, row.OccurrenceDuration)
	}

	return
}
