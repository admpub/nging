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

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/library/common"
	"github.com/nging-plugins/dbmanager/pkg/library/dbmanager/result"
)

var (
	drivers       = map[string]Driver{}
	DefaultDriver = &BaseDriver{}
)

type Driver interface {
	Init(echo.Context, *DbAuth)
	SetURLGenerator(func(string, ...string) string) Driver
	GenURL(string, ...string) string
	Results() []result.Resulter
	AddResults(...result.Resulter) Driver
	SetResults(...result.Resulter) Driver
	EnableFlashSession(on ...bool) Driver
	FlashSession() bool
	SaveResults() Driver
	SavedResults() interface{}
	IsSupported(string) bool
	Login() error
	Logout() error
	ProcessList() error
	Privileges() error
	Info() error
	CreateDb() error
	ModifyDb() error
	ListDb() error
	CreateTable() error
	ModifyTable() error
	ListTable() error
	ViewTable() error
	ListData() error
	CreateData() error
	Indexes() error
	Foreign() error
	Trigger() error
	RunCommand() error
	Import() error
	Export() error
	Analysis() error
	Name() string
}

func NewBaseDriver() *BaseDriver {
	return &BaseDriver{}
}

type BaseDriver struct {
	echo.Context
	*DbAuth
	results      []result.Resulter
	urlGenerator func(string, ...string) string
	flashSession bool
}

func (m *BaseDriver) SetURLGenerator(fn func(string, ...string) string) Driver {
	m.urlGenerator = fn
	return m
}

func (m *BaseDriver) EnableFlashSession(on ...bool) Driver {
	if len(on) == 0 || on[0] {
		m.flashSession = true
	} else {
		m.flashSession = false
	}
	return m
}

func (m *BaseDriver) FlashSession() bool {
	return m.flashSession
}

func (m *BaseDriver) GenURL(op string, args ...string) string {
	return m.urlGenerator(op, args...)
}

func (m *BaseDriver) Results() []result.Resulter {
	return m.results
}

func (m *BaseDriver) AddResults(rs ...result.Resulter) Driver {
	if m.results == nil {
		m.results = []result.Resulter{}
	}
	m.results = append(m.results, rs...)
	return m
}

func (m *BaseDriver) SetResults(rs ...result.Resulter) Driver {
	m.results = rs
	return m
}

func (m *BaseDriver) SaveResults() Driver {
	if m.results == nil {
		return m
	}
	if m.flashSession {
		/*
			if v, y := m.Flash(`dbMgrResults`).([]result.Resulter); y {
				m.results = append(v, m.results...)
			}
		*/
		m.Session().AddFlash(m.results, `dbMgrResults`)
	}
	return m
}

func (m *BaseDriver) SavedResults() interface{} {
	if v := m.Flash(`dbMgrResults`); v != nil {
		return v
	}
	return m.results
}

func (m *BaseDriver) Init(ctx echo.Context, auth *DbAuth) {
	m.Context = ctx
	m.DbAuth = auth
}
func (m *BaseDriver) IsSupported(operation string) bool {
	return true
}
func (m *BaseDriver) Login() error {
	return nil
}
func (m *BaseDriver) Logout() error {
	return nil
}
func (m *BaseDriver) ProcessList() error {
	return nil
}
func (m *BaseDriver) Privileges() error {
	return nil
}
func (m *BaseDriver) Info() error {
	return nil
}
func (m *BaseDriver) CreateDb() error {
	return nil
}
func (m *BaseDriver) ModifyDb() error {
	return nil
}
func (m *BaseDriver) ListDb() error {
	return nil
}
func (m *BaseDriver) CreateTable() error {
	return nil
}
func (m *BaseDriver) ModifyTable() error {
	return nil
}
func (m *BaseDriver) ListTable() error {
	return nil
}
func (m *BaseDriver) ViewTable() error {
	return nil
}
func (m *BaseDriver) ListData() error {
	return nil
}
func (m *BaseDriver) CreateData() error {
	return nil
}
func (m *BaseDriver) Indexes() error {
	return nil
}
func (m *BaseDriver) Foreign() error {
	return nil
}
func (m *BaseDriver) Trigger() error {
	return nil
}
func (m *BaseDriver) RunCommand() error {
	return nil
}
func (m *BaseDriver) Import() error {
	return nil
}
func (m *BaseDriver) Export() error {
	return nil
}
func (m *BaseDriver) Analysis() error {
	data := m.Data()
	data.SetInfo(m.T(`不支持此功能`), 0)
	return m.JSON(data)
}
func (m *BaseDriver) Name() string {
	return `Base`
}

//========================================================

func (m *BaseDriver) CheckErr(err error) interface{} {
	return common.Err(m.Context, err)
}

func (m *BaseDriver) SetOk(msg string) {
	common.SendOk(m.Context, msg)
}

func (m *BaseDriver) SetFail(msg string) {
	common.SendFail(m.Context, msg)
}

func (m *BaseDriver) Goto(rets ...string) error {
	next := m.Form(`next`)
	if len(next) == 0 {
		if len(rets) > 0 {
			next = rets[0]
		} else {
			next = m.Request().Referer()
		}
	}
	m.SaveResults()
	return m.Redirect(next)
}

func Register(name string, driver Driver) {
	drivers[name] = driver
}

func Get(name string) (Driver, bool) {
	d, y := drivers[name]
	return d, y
}

func GetForce(name string) Driver {
	d, y := drivers[name]
	if !y {
		d = DefaultDriver
	}
	return d
}

func Has(name string) bool {
	_, y := drivers[name]
	return y
}

func GetAll() map[string]Driver {
	return drivers
}

func Unregister(name string) {
	_, y := drivers[name]
	if y {
		delete(drivers, name)
	}
}
