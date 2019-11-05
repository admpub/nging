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
	"time"

	"github.com/admpub/sonyflake"
)

var sonyFlake *sonyflake.Sonyflake

func NewSonyflake(startDate string, machineIDs ...uint16) (*sonyflake.Sonyflake, error) { // 19位
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
	Init()
}

func Init() {
	err := SetSonyflake(`2018-09-01 08:08:08`)
	if err != nil {
		panic(err)
	}
}

func SetSonyflake(startDate string, machineID ...uint16) (err error) {
	sonyFlake, err = NewSonyflake(startDate, machineID...)
	return err
}

func UniqueID() (string, error) {
	id, err := NextID()
	if err != nil {
		return ``, err
	}
	return fmt.Sprintf(`%d`, id), nil
}

func NextID() (uint64, error) {
	if sonyFlake == nil {
		Init()
	}
	return sonyFlake.NextID()
}
