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

package hosts

import (
	"io/ioutil"
	"os"

	"github.com/webx-top/com"
)

var hostsPath string

func init() {
	Init()
}

func Path() string {
	if len(hostsPath) == 0 {
		Init()
	}
	return hostsPath
}

func Init() {
	var paths []string
	if com.IsWindows {
		winDir := os.Getenv(`WinDir`)
		if len(winDir) == 0 {
			winDir = os.Getenv(`windir`)
		}
		paths = append(
			paths,
			os.Getenv(`SystemRoot`)+`\system32\drivers\etc\hosts`,
		)
		if len(winDir) > 0 {
			paths = append(
				paths,
				winDir+`\hosts`,
			)
		}
	} else {
		paths = append(
			paths,
			`/etc/hosts`,
			`/private/etc/hosts`,
			`/system/etc/hosts`, //Android
		)
	}
	for _, hpath := range paths {
		if com.FileExists(hpath) {
			hostsPath = hpath
			break
		}
	}
}

func ReadFile() ([]byte, error) {
	return ioutil.ReadFile(Path())
}

func WriteFile(content []byte) error {
	return ioutil.WriteFile(Path(), content, os.ModePerm)
}
