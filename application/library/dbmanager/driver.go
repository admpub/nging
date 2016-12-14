package dbmanager

import (
	//"github.com/webx-top/com"
	"github.com/webx-top/echo"
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
