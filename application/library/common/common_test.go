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

package common

import (
	"fmt"
	"testing"
)

func TestA(t *testing.T) {
	fmt.Println(DirSharding(10000))
	fmt.Println(ModifyAsThumbnailName(`/upload/news/1/123233232.png`, `/upload/news/0/123233232_200_200.png`))
	fmt.Println(ParseTemporaryFileName(`<p><img style="" src="/public/upload/news/0/29632008003059711.jpg"><img style="" src="/public/upload/news/0/29632023740088319.jpg"><br></p>`))
}
