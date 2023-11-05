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
	"github.com/nging-plugins/collector/application/dbschema"
	"github.com/nging-plugins/collector/application/model"
	"github.com/webx-top/echo"
)

func parseFormToDb(c echo.Context, pageM *dbschema.NgingCollectorPage, formPrefix string, update bool) (rules []*dbschema.NgingCollectorRule, err error) {
	if pageM.Id > 0 {
		err = pageM.Update(nil, `id`, pageM.Id)
	} else {
		if len(pageM.HasChild) == 0 {
			pageM.HasChild = `N`
		}
		_, err = pageM.Insert()
	}
	if err != nil {
		return
	}
	if pageM.Id < 1 {
		err = c.E(`Id无效`)
		return
	}
	ids := c.FormxValues(formPrefix + `[id][]`).Uint()
	idCount := len(ids)
	vars := c.FormValues(formPrefix + `[var][]`)
	varCount := len(vars)
	ruleList := c.FormValues(formPrefix + `[rule][]`)
	ruleCount := len(ruleList)
	types := c.FormValues(formPrefix + `[type][]`)
	typeCount := len(types)
	filters := c.FormValues(formPrefix + `[filter][]`)
	filterCount := len(filters)
	for i := 0; i < varCount; i++ {
		rule := dbschema.NgingCollectorRule{
			PageId: pageM.Id,
			Name:   vars[i],
			Type:   ``,
			Filter: ``,
			Rule:   ``,
			Sort:   i,
		}
		if update {
			if i >= idCount {
				break
			}
			rule.Id = ids[i]
		}
		if i >= ruleCount {
			break
		}
		rule.Rule = ruleList[i]

		if i >= typeCount {
			break
		}
		rule.Type = types[i]

		if i >= filterCount {
			break
		}
		rule.Filter = filters[i]

		rule.Use(pageM.Trans())
		if rule.Id > 0 {
			err = rule.Update(nil, `id`, rule.Id)
		} else {
			_, err = rule.Insert()
		}
		if err != nil {
			return
		}
		if update {
			rules = append(rules, &rule)
		}
	}
	return
}

var DataTypes = echo.NewKVData().Add("int", "整数").Add("float", "小数").Add("bool", "布尔值").Add("string", "字符串").
	Add("int-array", "整数数组").Add("float-array", "小数数组").Add("bool-array", "布尔值数组").Add("string-array", "字符串数组").
	Add("map", "字典").Add("array", "数组").
	Add("json", "JSON数据作为结果值").Add("jsonparse", "从JSON数据中解析").

	// begin compatible old version
	Add("href", "href属性值").Add("src", "src属性值").Add("alt", "alt属性值").
	Add("text-array", "纯文本数组").Add("href-array", "href属性值数组").
	Add("html-array", "HTML内容数组").Add("html", "innerHTML内容").Add("outhtml", "父级HTML").
	Add("raw", "直接用规则作为结果值")

func setFormData(c echo.Context, pageM *model.CollectorPage) error {
	data, err := pageM.FullData()
	c.Set(`data`, data)
	return err
}
