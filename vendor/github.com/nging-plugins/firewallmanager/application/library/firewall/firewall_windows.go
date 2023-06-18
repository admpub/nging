//go:build windows

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

package firewall

import (
	"github.com/admpub/once"
)

func initBackend() {
}

var engine driver.Driver
var engonce once.Once

func initEngine() {
	var err error
	engine, err = netsh.New()
	if err != nil {
		panic(err)
	}
}

func Engine(_ string) driver.Driver {
	engonce.Do(initEngine)
	return engine
}

func ResetEngine() {
	engonce.Reset()
}
