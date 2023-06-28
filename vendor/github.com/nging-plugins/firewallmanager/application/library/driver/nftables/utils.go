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

package nftables

import (
	"strconv"
	"strings"

	"github.com/nging-plugins/firewallmanager/application/library/cmdutils"
)

func LineParser(i uint, t string) (rowInfo cmdutils.RowInfo, err error) {
	if strings.HasSuffix(t, `{`) || t == `}` {
		return
	}
	parts := strings.SplitN(t, `# handle `, 2)
	if len(parts) == 2 {
		parts[0] = strings.TrimSpace(parts[0])
		if strings.HasSuffix(parts[0], `{`) {
			return
		}
		var handleID uint64
		handleID, err = strconv.ParseUint(parts[1], 10, 0)
		if err != nil {
			return
		}
		rowInfo = cmdutils.RowInfo{
			RowNo: i,
			Row:   parts[0],
		}
		rowInfo.Handle.SetValid(uint(handleID))
	} else {
		rowInfo = cmdutils.RowInfo{
			RowNo: i,
			Row:   t,
		}
	}
	return
}
