package dbmanager

import "github.com/webx-top/echo"

type mysql struct {
	echo.Context
}

func (m *mysql) Init(ctx echo.Context) {
	m.Context = ctx
}
func (m *mysql) IsSupported(operation string) bool {
	return true
}
func (m *mysql) Login() error {
	return nil
}
func (m *mysql) Logout() error {
	return nil
}
func (m *mysql) ProcessList() error {
	return nil
}
func (m *mysql) Privileges() error {
	return nil
}
func (m *mysql) Info() error {
	return nil
}
func (m *mysql) CreateDb() error {
	return nil
}
func (m *mysql) ModifyDb() error {
	return nil
}
func (m *mysql) ListDb() error {
	return nil
}
func (m *mysql) CreateTable() error {
	return nil
}
func (m *mysql) ModifyTable() error {
	return nil
}
func (m *mysql) ListTable() error {
	return nil
}
func (m *mysql) ViewTable() error {
	return nil
}
func (m *mysql) ListData() error {
	return nil
}
func (m *mysql) CreateData() error {
	return nil
}
func (m *mysql) Indexes() error {
	return nil
}
func (m *mysql) Foreign() error {
	return nil
}
func (m *mysql) Trigger() error {
	return nil
}
func (m *mysql) RunCommand() error {
	return nil
}
func (m *mysql) Import() error {
	return nil
}
func (m *mysql) Export() error {
	return nil
}
