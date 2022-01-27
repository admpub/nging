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

package ssystem

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/echo"
)

type System struct {
	Env                     string            `json:"env"` // prod/dev/test
	VhostsfileDir           string            `json:"vhostsfileDir"`
	AllowIP                 []string          `json:"allowIP"`
	SSLAuto                 bool              `json:"sslAuto"`
	SSLEmail                string            `json:"sslEmail"`
	SSLHosts                []string          `json:"sslHosts"`
	SSLCacheDir             string            `json:"sslCacheDir"`
	SSLKeyFile              string            `json:"sslKeyFile"`
	SSLCertFile             string            `json:"sslCertFile"`
	EditableFileExtensions  map[string]string `json:"editableFileExtensions"`
	EditableFileMaxSize     string            `json:"editableFileMaxSize"`
	editableFileMaxBytes    int               `json:"editableFileMaxBytes"`
	PlayableFileExtensions  map[string]string `json:"playableFileExtensions"`
	ErrorPages              map[int]string    `json:"errorPages"`
	CmdTimeout              string            `json:"cmdTimeout"`
	CmdTimeoutDuration      time.Duration     `json:"-"`
	ShowExpirationTime      int64             `json:"showExpirationTime"` //显示过期时间：0为始终显示；大于0为距离剩余到期时间多少秒的时候显示；小于0为不显示
	SessionName             string            `json:"sessionName"`
	SessionEngine           string            `json:"sessionEngine"`
	SessionConfig           echo.H            `json:"sessionConfig"`
	MaxRequestBodySize      string            `json:"maxRequestBodySize"`
	maxRequestBodySizeBytes int
}

func (sys *System) Init() {
	if len(sys.MaxRequestBodySize) == 0 {
		sys.maxRequestBodySizeBytes = 0
		return
	}
	sys.maxRequestBodySizeBytes, _ = ParseBytes(sys.MaxRequestBodySize)
	if sys.editableFileMaxBytes < 1 && len(sys.EditableFileMaxSize) > 0 {
		var err error
		sys.editableFileMaxBytes, err = ParseBytes(sys.EditableFileMaxSize)
		if err != nil {
			log.Error(err.Error())
		}
	}
	sys.CmdTimeoutDuration = ParseTimeDuration(sys.CmdTimeout)
	if sys.CmdTimeoutDuration <= 0 {
		sys.CmdTimeoutDuration = time.Second * 30
	}
	if len(sys.SSLCacheDir) == 0 {
		sys.SSLCacheDir = `data` + echo.FilePathSeparator + `cache` + echo.FilePathSeparator + `autocert`
	}
}

func (sys *System) EditableFileMaxBytes() int {
	return sys.editableFileMaxBytes
}

func (sys *System) MaxRequestBodySizeBytes() int {
	if sys.maxRequestBodySizeBytes <= 0 {
		return sys.maxRequestBodySizeBytes
	}
	return sys.maxRequestBodySizeBytes
}

func (sys *System) Editable(fileName string) (string, bool) {
	if sys.EditableFileExtensions == nil {
		return "", false
	}
	ext := strings.TrimPrefix(filepath.Ext(fileName), `.`)
	ext = strings.ToLower(ext)
	typ, ok := sys.EditableFileExtensions[ext]
	return typ, ok
}

func (sys *System) IsEnv(name string) bool {
	if name == `prod` {
		if len(sys.Env) == 0 {
			return true
		}
	}
	return sys.Env == name
}

func (sys *System) Playable(fileName string) (string, bool) {
	if sys.PlayableFileExtensions == nil {
		sys.PlayableFileExtensions = map[string]string{
			`mp4`:  `video/mp4`,
			`m3u8`: `application/x-mpegURL`,
			//`ts`:   `video/MP2T`,
		}
	}
	ext := strings.TrimPrefix(filepath.Ext(fileName), `.`)
	ext = strings.ToLower(ext)
	typ, ok := sys.PlayableFileExtensions[ext]
	return typ, ok
}
