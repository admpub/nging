//go:build linux

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

	"github.com/nging-plugins/firewallmanager/application/library/driver"
	"github.com/nging-plugins/firewallmanager/application/library/driver/iptables"
	"github.com/nging-plugins/firewallmanager/application/library/driver/nftables"
)

func initBackend() {
	if nftables.IsSupported() {
		backend = `nftables`
	} else if iptables.IsSupported() {
		backend = `iptables`
	}
}

var engineIPv4 driver.Driver
var engonceIPv4 once.Once
var engineIPv6 driver.Driver
var engonceIPv6 once.Once

func initEngineIPv4() {
	var err error
	if GetBackend() == `nftables` {
		engineIPv4, err = nftables.New(driver.ProtocolIPv4)
	} else {
		engineIPv4, err = iptables.New(driver.ProtocolIPv4, false)
	}
	if err != nil {
		panic(err)
	}
}

func EngineIPv4() driver.Driver {
	engonceIPv4.Do(initEngineIPv4)
	return engineIPv4
}

func initEngineIPv6() {
	var err error
	if GetBackend() == `nftables` {
		engineIPv6, err = nftables.New(driver.ProtocolIPv6)
	} else {
		engineIPv6, err = iptables.New(driver.ProtocolIPv6, false)
	}
	if err != nil {
		panic(err)
	}
}

func EngineIPv6() driver.Driver {
	engonceIPv6.Do(initEngineIPv6)
	return engineIPv6
}

func Engine(ipVersionNumber string) driver.Driver {
	if ipVersionNumber == `6` {
		return EngineIPv6()
	}
	return EngineIPv4()
}

func ResetEngine() {
	engonceIPv4.Reset()
	engonceIPv6.Reset()
}
