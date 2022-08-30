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

package logcool

type Counter struct {
	Views       map[string]uint64  //每个网址的浏览器
	Elapsed     map[string]float64 //每个网址的总耗时(秒)
	MaxElapsed  map[string]float64 //每个网址的最大耗时(秒)
	BodyBytes   map[string]uint64  //每个网址的总字节数
	IPs         map[string]uint64  //每个网址的每分钟的IP数量
	UserAgents  map[string][]int64
	StatusCodes map[string][]int64

	//200/301/302/400/403/404/499/500/502/503/504
	StatusCode map[int]uint64    //状态码次数
	Browsers   map[string]uint64 //浏览器访问次数
}
