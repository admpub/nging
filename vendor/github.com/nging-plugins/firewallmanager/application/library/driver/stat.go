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

package driver

import "net"

// Stat represents a structured statistic entry.
type Stat struct {
	Number      uint64     `json:"num,omitempty"`
	Packets     uint64     `json:"pkts"`
	Bytes       uint64     `json:"bytes"`
	Target      string     `json:"target"`
	Protocol    string     `json:"prot"`
	Opt         string     `json:"opt"`
	Input       string     `json:"in"`
	Output      string     `json:"out"`
	Source      *net.IPNet `json:"source"`
	Destination *net.IPNet `json:"destination"`
	Options     string     `json:"options"`
}
