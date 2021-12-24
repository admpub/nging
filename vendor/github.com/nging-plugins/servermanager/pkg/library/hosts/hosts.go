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
	"database/sql"
	"os"

	"github.com/webx-top/com"
)

var hostsPath sql.NullString

func init() {
	Init()
}

func Path() string {
	if len(hostsPath.String) == 0 && !hostsPath.Valid {
		Init()
	}
	return hostsPath.String
}

func detectFile(paths []string) {
	hostsPath.Valid = true
	for _, hpath := range paths {
		if !com.FileExists(hpath) {
			continue
		}
		hostsPath.String = hpath
		break
	}
}

func ReadFile() ([]byte, error) {
	return os.ReadFile(Path())
}

func WriteFile(content []byte) error {
	return os.WriteFile(Path(), content, os.ModePerm)
}
