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
	"strings"
	"time"

	"github.com/admpub/sonyflake"
)

var (
	sonyFlakes         = map[uint16]*sonyflake.Sonyflake{}
	SonyflakeStartDate = `2018-09-01 08:08:08`
)

//NewSonyflake 19位
func NewSonyflake(startDate string, machineIDs ...uint16) (*sonyflake.Sonyflake, error) {
	if !strings.Contains(startDate, ` `) {
		startDate += ` 00:00:00`
	}
	startTime, err := time.ParseInLocation(`2006-01-02 15:04:05`, startDate, time.Local)
	if err != nil {
		return nil, err
	}
	var machineID uint16
	if len(machineIDs) > 0 {
		machineID = machineIDs[0]
	}
	st := sonyflake.Settings{
		StartTime: startTime,
		MachineID: func() (uint16, error) {
			return machineID, nil
		},
		CheckMachineID: func(id uint16) bool {
			return machineID == id
		},
	}
	return sonyflake.NewSonyflake(st), err
}

func init() {
	SonyflakeInit(0)
}

func SonyflakeInit(machineIDs ...uint16) *sonyflake.Sonyflake {
	sonyFlake, err := SetSonyflake(SonyflakeStartDate, machineIDs...)
	if err != nil {
		panic(err)
	}
	return sonyFlake
}

func SetSonyflake(startDate string, machineIDs ...uint16) (sonyFlake *sonyflake.Sonyflake, err error) {
	sonyFlake, err = NewSonyflake(startDate, machineIDs...)
	if err != nil {
		return nil, err
	}
	var machineID uint16
	if len(machineIDs) > 0 {
		machineID = machineIDs[0]
	}
	sonyFlakes[machineID] = sonyFlake
	return sonyFlake, err
}

func UniqueID(machineIDs ...uint16) (string, error) {
	id, err := NextID(machineIDs...)
	if err != nil {
		return ``, err
	}
	return fmt.Sprintf(`%d`, id), nil
}

func NextID(machineIDs ...uint16) (uint64, error) {
	var machineID uint16
	if len(machineIDs) > 0 {
		machineID = machineIDs[0]
	}
	sonyFlake, ok := sonyFlakes[machineID]
	if !ok || sonyFlake == nil {
		var err error
		sonyFlake, err = SetSonyflake(SonyflakeStartDate, machineIDs...)
		if err != nil {
			return 0, err
		}
	}
	return sonyFlake.NextID()
}
