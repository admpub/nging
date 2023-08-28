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

package charset

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	str := `历史不止一次的告诉我们`
	abc := strings.Repeat(`.`, With(str))
	fmt.Println(str)
	fmt.Println(abc)
	assert.Equal(t, `历史不止`, Truncate(str, 9))
}

func TestConvert(t *testing.T) {
	b, err := os.ReadFile(`gbktest.txt`)
	assert.NoError(t, err)
	r, err := Convert(`gbk`, `utf8`, b)
	assert.NoError(t, err)
	str := `炎黄子孙
华夏民族`
	assert.Equal(t, str, string(r))

	r, err = Convert(`utf8`, `gbk`, r)
	assert.NoError(t, err)
	assert.Equal(t, b, r)
}
