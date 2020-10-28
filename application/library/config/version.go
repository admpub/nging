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

package config

import (
	"time"

	"github.com/webx-top/com"
)

// Version 版本信息
var Version = &VersionInfo{Name: `Nging`}
var versionLabelWeight = map[string]int{
	`stable`: 0,
	`beta`:   1,
	`alpha`:  2,
}

type VersionInfo struct {
	Name      string    //软件名称
	Number    string    //版本号 1.0.1
	Package   string    //套餐
	Label     string    //版本标签 beta/alpha/stable
	DBSchema  float64   //数据库表版本 例如：1.2
	BuildTime string    //构建时间
	CommitID  string    //GIT提交ID
	Licensed  bool      //是否已授权
	Expired   time.Time //过期时间
}

func (v *VersionInfo) IsExpired() bool {
	if v.Expired.IsZero() {
		return false
	}
	return v.Expired.Before(time.Now())
}

func (v *VersionInfo) String() string {
	return v.Name + ` ` + v.VString()
}

func (v *VersionInfo) VString() string {
	var licenseTag string
	if v.Licensed {
		licenseTag = `licensed`
	} else {
		licenseTag = `unlicensed`
	}
	var label string
	if v.Label != `stable` {
		label = `(` + v.Label + `)`
	}
	return `v` + v.Number + label + ` ` + licenseTag
}

func (v *VersionInfo) IsNew(number string, label string) bool {
	var hasNew bool
	compared := com.VersionCompare(number, v.Number)
	if compared == com.VersionCompareGt {
		hasNew = true
	} else if compared == com.VersionCompareEq {
		if weight, ok := versionLabelWeight[label]; ok {
			if currWeight, ok := versionLabelWeight[v.Label]; ok {
				hasNew = weight < currWeight
			}
		}
	}
	return hasNew
}
