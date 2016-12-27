/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package driver

import (
	"github.com/admpub/nging/application/library/dbmanager/result"
	"github.com/webx-top/echo"
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
}

func NewBaseDriver() *BaseDriver {
	return &BaseDriver{}
}

type BaseDriver struct {
	echo.Context
	*DbAuth
	results      []result.Resulter
	urlGenerator func(string, ...string) string
}

func (m *BaseDriver) SetURLGenerator(fn func(string, ...string) string) Driver {
	m.urlGenerator = fn
	return m
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
	/*
		if v, y := m.Flash(`dbMgrResults`).([]result.Resulter); y {
			m.results = append(v, m.results...)
		}
	*/
	m.Session().AddFlash(m.results, `dbMgrResults`)
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
