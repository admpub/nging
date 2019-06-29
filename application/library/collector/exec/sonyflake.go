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

package exec

import (
	"fmt"
	"strconv"
	"time"

	"github.com/admpub/sonyflake"
)

var SonyFlake *sonyflake.Sonyflake

func init() {
	// 19位
	startTime, _ := time.ParseInLocation(`2006-01-02 15:04:05`, `2018-09-01 08:08:08`, time.Local)
	month, _ := strconv.ParseUint(time.Now().Format(`200601`), 10, 16)
	monthID := uint16(month)
	st := sonyflake.Settings{
		StartTime: startTime,
		MachineID: func() (uint16, error) {
			return monthID, nil
		},
		CheckMachineID: func(id uint16) bool {
			return monthID == id
		},
	}
	SonyFlake = sonyflake.NewSonyflake(st)
}

func UniqueID() (string, error) {
	id, err := SonyFlake.NextID()
	if err != nil {
		return ``, err
	}
	return fmt.Sprintf(`%d`, id), nil
}
