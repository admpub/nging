package dbmanager

import (
	//"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var (
	drivers       = map[string]Driver{}
	DefaultDriver = &BaseDriver{}
)

type Driver interface {
	Init(echo.Context)
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
}

func (m *BaseDriver) Init(ctx echo.Context) {
	m.Context = ctx
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
