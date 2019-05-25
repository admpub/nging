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

package navigate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigDefaultsAsStore(t *testing.T) {
	v := List{
		{
			Name: `0`,
		},
		{
			Name: `1`,
		},
		{
			Name: `2`,
		},
	}
	v2 := List{
		{
			Name: `00`,
		},
		{
			Name: `01`,
		},
		{
			Name: `02`,
		},
	}

	v.Add(0, v2...)
	assert.Equal(t, List{
		{
			Name: `00`,
		},
		{
			Name: `01`,
		},
		{
			Name: `02`,
		},
		{
			Name: `0`,
		},
		{
			Name: `1`,
		},
		{
			Name: `2`,
		},
	}, v)

	v.Remove(0)
	assert.Equal(t, List{
		{
			Name: `01`,
		},
		{
			Name: `02`,
		},
		{
			Name: `0`,
		},
		{
			Name: `1`,
		},
		{
			Name: `2`,
		},
	}, v)

	v.Remove(4)
	assert.Equal(t, List{
		{
			Name: `01`,
		},
		{
			Name: `02`,
		},
		{
			Name: `0`,
		},
		{
			Name: `1`,
		},
	}, v)

	v.Set(0, List{
		{
			Name: `11`,
		},
		{
			Name: `12`,
		},
	}...)
	assert.Equal(t, List{
		{
			Name: `11`,
		},
		{
			Name: `12`,
		},
		{
			Name: `0`,
		},
		{
			Name: `1`,
		},
	}, v)

}
