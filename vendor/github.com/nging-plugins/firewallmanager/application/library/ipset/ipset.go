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

package ipset

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

// Info holds ipset list contents
type Info struct {
	Name         string
	Revision     int
	Header       string
	SizeInMemory int
	References   int
	NumEntries   int
	Entries      []string
}

func ParseInfo(reader io.Reader) (info *Info, err error) {
	info = &Info{}
	s := bufio.NewScanner(reader)

	for s.Scan() {
		t := s.Text()
		parts := strings.SplitN(t, `:`, 2)
		if len(parts) != 2 {
			continue
		}
		switch parts[0] {
		case "Name":
			info.Name = strings.TrimSpace(parts[1])
		case "Revision":
			if info.Revision, err = getNumber(t); err != nil {
				return nil, err
			}
		case "Header":
			info.Header = strings.TrimSpace(parts[1])
		case "Size in memory":
			if info.SizeInMemory, err = getNumber(t); err != nil {
				return nil, err
			}
		case "References":
			if info.References, err = getNumber(t); err != nil {
				return nil, err
			}
		case "Number of entries":
			if info.NumEntries, err = getNumber(t); err != nil {
				return nil, err
			}
		case "Members":
			goto Entries
		}
	}

Entries:
	for s.Scan() {
		info.Entries = append(info.Entries, s.Text())
	}

	return
}

func getNumber(t string) (n int, err error) {
	if i := strings.LastIndexByte(t, ' '); i != -1 {
		return strconv.Atoi(t[i+1:])
	}
	return
}
