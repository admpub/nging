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

package settings

import (
	"testing"

	"github.com/admpub/nging/v4/application/dbschema"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
)

var expect = `{
  "base": {
    "debug": "test"
  }
}`

func TestConfigDefaultsAsStore(t *testing.T) {
	actual := echo.Dump(configAsStore(map[string]map[string]*dbschema.NgingConfig{
		`base`: {
			`debug`: {
				Key:         `debug`,
				Label:       `调试模式`,
				Description: ``,
				Value:       `test`,
				Group:       `base`,
				Type:        `text`,
				Sort:        0,
				Disabled:    `N`,
			},
		},
	}), false)
	assert.Equal(t, expect, actual)
}
