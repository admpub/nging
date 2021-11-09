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

package notice

import (
	"sync"
)

func newNoticeTypes() *noticeTypes {
	return &noticeTypes{types: map[string]bool{}}
}

type noticeTypes struct {
	types map[string]bool
	lock  sync.RWMutex
}

func (n *noticeTypes) Has(types ...string) bool {
	n.lock.RLock()
	for _, typ := range types {
		if !n.types[typ] {
			return false
		}
	}
	n.lock.RUnlock()
	return true
}

func (n *noticeTypes) Clear(types ...string) {
	n.lock.Lock()
	if len(types) > 0 {
		for _, typ := range types {
			_, ok := n.types[typ]
			if !ok {
				continue
			}
			delete(n.types, typ)
		}
	} else {
		for typ := range n.types {
			delete(n.types, typ)
		}
	}
	n.lock.Unlock()
}

func (n *noticeTypes) Open(types ...string) {
	n.lock.Lock()
	if len(types) > 0 {
		for _, typ := range types {
			if opened := n.types[typ]; !opened {
				n.types[typ] = true
			}
		}
	} else {
		for key := range n.types {
			if opened := n.types[key]; !opened {
				n.types[key] = true
			}
		}
	}
	n.lock.Unlock()
}
