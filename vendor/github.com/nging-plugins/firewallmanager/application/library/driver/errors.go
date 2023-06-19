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

package driver

import (
	"errors"
	"fmt"
)

var ErrUnsupported = errors.New(`unsupported`)
var ErrUnsupportedTable = fmt.Errorf(`%w table`, ErrUnsupported)
var ErrUnsupportedChain = fmt.Errorf(`%w chain`, ErrUnsupported)
var ErrNatIPOrNatPortRequired = errors.New(`NAT IP 和 NAT 端口 不能同时为空`)
var ErrInvalidIPv4 = errors.New(`invalid IPv4 address`)
var ErrInvalidIPv6 = errors.New(`invalid IPv6 address`)
var ErrInvalidRuleNumber = errors.New(`invalid rule number`)
