//go:build !linux && !windows

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
	"errors"

	"github.com/nging-plugins/firewallmanager/application/library/driver"
)

func initBackend() {
}

func Engine(ipVersionNumber string) driver.Driver {
	return defaultUnsupportedDriver
}

func ResetEngine() {
}

var ErrUnsupportedOperatingSystem = errors.New(`This feature is not supported in the current operating system`) //此功能在当前操作系统里不支持

var defaultUnsupportedDriver = &unsupportedDriver{}

type unsupportedDriver struct {
}

func (unsupportedDriver) RuleFrom(rule *driver.Rule) []string {
	return nil
}

func (unsupportedDriver) Enabled(on bool) error {
	return ErrUnsupportedOperatingSystem
}

func (unsupportedDriver) Reset() error {
	return ErrUnsupportedOperatingSystem
}

func (unsupportedDriver) Clear() error {
	return ErrUnsupportedOperatingSystem
}

func (unsupportedDriver) Import(wfwFile string) error {
	return ErrUnsupportedOperatingSystem
}

func (unsupportedDriver) Export(wfwFile string) error {
	return ErrUnsupportedOperatingSystem
}

func (unsupportedDriver) Insert(rules ...driver.Rule) error {
	return ErrUnsupportedOperatingSystem
}

func (unsupportedDriver) AsWhitelist(table, chain string) error {
	return nil
}

func (unsupportedDriver) Append(rules ...driver.Rule) error {
	return ErrUnsupportedOperatingSystem
}

func (unsupportedDriver) Update(rule driver.Rule) error {
	return ErrUnsupportedOperatingSystem
}

func (unsupportedDriver) Delete(rules ...driver.Rule) error {
	return ErrUnsupportedOperatingSystem
}

func (unsupportedDriver) Exists(rule driver.Rule) (bool, error) {
	return false, ErrUnsupportedOperatingSystem
}

func (unsupportedDriver) Stats(table, chain string) ([]map[string]string, error) {
	return nil, nil
}

func (unsupportedDriver) List(table, chain string) ([]*driver.Rule, error) {
	return nil, ErrUnsupportedOperatingSystem
}

func (unsupportedDriver) FindPositionByID(table, chain string, id uint) (uint64, error) {
	return 0, ErrUnsupportedOperatingSystem
}
