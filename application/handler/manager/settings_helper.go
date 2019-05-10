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

package manager

import (
	"github.com/admpub/nging/application/registry/settings"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func configPost(c echo.Context) error {
	m := model.NewConfig(c)
	formValues := c.Forms()
	mapx := echo.NewMapx(formValues)
	for group, configs := range settings.ConfigDefaults() {
		_, err := m.ListByGroup(group)
		if err != nil {
			if err == db.ErrNoMoreRows {
				for _, cfg := range configs {
					_, err := cfg.Add()
					if err != nil {
						return err
					}
				}
				continue
			}
			return err
		}
		encoder := settings.GetEncoder(group)
		gm := mapx.Get(group)
		if gm == nil {
			continue
		}
		added := map[string]int{}
		for k, v := range m.Objects() {
			added[v.Key] = k
			setting := gm.Get(v.Key)
			if setting == nil {
				continue
			}
			disabled := setting.Value(`disabled`)
			_v := setting.Get(`value`)
			if _v == nil {
				continue
			}
			value, err := settings.EncodeConfigValue(_v, v, encoder)
			if err != nil {
				return err
			}
			var n int64
			condition := db.And(
				db.Cond{`key`: v.Key},
				db.Cond{`group`: group},
			)
			n, err = m.Count(nil, condition)
			if err != nil {
				return err
			}
			if n < 1 {
				err = settings.InsertBy(configs, v.Key, value, disabled)
				if err != nil {
					return err
				}
			}
			set := echo.H{
				`value`: value,
			}
			if _v.IsMap() {
				set[`type`] = `json`
			} else if _v.IsSlice() {
				set[`type`] = `list`
			} else {
				set[`type`] = `text`
			}
			if len(disabled) > 0 {
				set[`disabled`] = disabled
			}
			err = m.NewParam().SetArgs(condition).SetSend(set).Update()
			if err != nil {
				return err
			}
		}
		err = settings.InsertMissing(gm, added, configs, encoder)
		if err != nil {
			return err
		}
	}
	return nil
}

func configGet(c echo.Context) error {
	m := model.NewConfig(c)
	_, err := m.ListByOffset(nil, nil, 0, -1)
	if err != nil {
		return err
	}
	for _, setting := range m.Objects() {
		group := setting.Group
		cfg, err := m.ListMapByGroup(group)
		if err != nil {
			return err
		}
		c.Set(group, cfg) //Stored.base.siteName
	}
	return nil
}
