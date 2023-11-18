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

package utils

import (
	"strings"

	"github.com/docker/docker/pkg/stdcopy"
)

func ShortenID(imageID string) string {
	cleanedID := strings.TrimPrefix(imageID, `sha256:`)
	if len(cleanedID) > 12 {
		return cleanedID[0:12]
	}
	return imageID
}

func TrimHeader(message string) string {
	//fmt.Printf(`%d|%d|%d|%d|%d|%d|%d|%d|%d`+"\n", message[0], message[1], message[2], message[3], message[4], message[5], message[6], message[7], message[8])
	if len(message) > 9 && (message[0] == 0 /*stdin*/ ||
		message[0] == 1 /*stdout*/ ||
		message[0] == 2 /*stderr*/) {
		message = message[9:]
	}
	return message
}

var StdCopy = stdcopy.StdCopy
